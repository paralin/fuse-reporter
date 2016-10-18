package service

import (
	"encoding/json"
	"golang.org/x/net/context"
	"time"

	"github.com/fuserobotics/reporter"
	"github.com/fuserobotics/reporter/util"
	"github.com/fuserobotics/reporter/view"
	statestream "github.com/fuserobotics/statestream"
)

type ViewServiceServer struct {
	Reporter *reporter.Reporter
}

func (s *ViewServiceServer) ListStates(c context.Context, req *view.ListStatesRequest) (*view.ListStatesResponse, error) {
	list, err := buildStateList(s.Reporter.LocalTree)
	if err != nil {
		return nil, err
	}
	return &view.ListStatesResponse{
		List: list,
	}, nil
}

func (s *ViewServiceServer) GetState(c context.Context, req *view.GetStateRequest) (*view.GetStateResponse, error) {
	if err := req.Validate(false); err != nil {
		return nil, err
	}

	component, err := s.Reporter.LocalTree.GetComponent(req.Context.Component)
	if err != nil {
		return nil, err
	}
	state, err := component.GetState(req.Context.StateId)
	if err != nil {
		return nil, err
	}

	stream := state.Stream()
	var cursor *statestream.Cursor

	if req.Query.Time <= 0 {
		writeCursor, err := stream.WriteCursor()
		if err != nil {
			return nil, err
		}
		cursor = writeCursor
	} else {
		cursor = stream.BuildCursor(statestream.ReadForwardCursor)
		if err := cursor.Init(util.NumberToTime(req.Query.Time)); err != nil {
			return nil, err
		}
	}

	stateData, err := cursor.State()
	if err != nil {
		return nil, err
	}
	stateMap := map[string]interface{}(stateData)

	respData, err := json.Marshal(&stateMap)
	if err != nil {
		return nil, err
	}

	return &view.GetStateResponse{
		State: &view.StateReport{
			JsonState: string(respData),
			Timestamp: util.TimeToNumber(cursor.ComputedTimestamp()),
		},
	}, nil
}

func buildStateList(tree *reporter.ComponentTree) (*view.StateList, error) {
	res := &view.StateList{
		Components: make([]*view.StateListComponent, len(tree.ComponentList.Data.ComponentName)),
	}

	for i, componentName := range tree.ComponentList.Data.ComponentName {
		cmp, err := tree.GetComponent(componentName)
		if err != nil {
			return nil, err
		}
		slc := &view.StateListComponent{Name: componentName, States: make([]*view.StateListState, len(cmp.Data.StateName))}
		res.Components[i] = slc
		for is, stateName := range cmp.Data.StateName {
			state, err := cmp.GetState(stateName)
			if err != nil {
				return nil, err
			}
			slc.States[is] = &view.StateListState{
				Name:   stateName,
				Config: state.Data.StreamConfig,
			}
		}
	}
	return res, nil
}

func (s *ViewServiceServer) GetStateHistory(req *view.StateHistoryRequest, srvstream view.ReporterService_GetStateHistoryServer) error {
	if err := req.Context.Validate(false); err != nil {
		return err
	}
	if err := req.Query.Validate(); err != nil {
		return err
	}

	component, err := s.Reporter.LocalTree.GetComponent(req.Context.Component)
	if err != nil {
		return err
	}
	state, err := component.GetState(req.Context.StateId)
	if err != nil {
		return err
	}

	stream := state.Stream()
	var cursor *statestream.Cursor

	// first send the history
	earlyBound := util.NumberToTime(req.Query.BeginTime)
	lateBound := util.NumberToTime(req.Query.EndTime)
	if req.Query.EndTime == 0 {
		lateBound = time.Now()
	}

	cursor = stream.BuildCursor(statestream.ReadForwardCursor)
	ch := make(chan *statestream.StreamEntry)
	sub := cursor.SubscribeEntries(ch)
	// spawn a goroutine to send the data, while simultaneously fast forwarding
	go func() {
		for {
			select {
			case <-srvstream.Context().Done():
				return
			case entry, ok := <-ch:
				if !ok {
					return
				}
				jsonData, err := json.Marshal(entry.Data)
				if err != nil {
					break
				}
				// not sure what to do if we get an error here.
				srvstream.Send(&view.StateHistoryResponse{
					Status: view.StateHistoryResponse_HISTORY_INITIAL_SET,
					State: &view.StateEntry{
						JsonState: string(jsonData),
						Timestamp: util.TimeToNumber(entry.Timestamp),
						Type:      int32(entry.Type),
					},
				})
			}
		}
	}()

	if err := cursor.Init(earlyBound); err != nil {
		if err == statestream.NoDataError {
			// Check if we can get something after it.
			entry, err := state.GetEntryAfter(earlyBound, statestream.StreamEntrySnapshot)
			if err != nil {
				return err
			}
			if entry == nil || entry.Timestamp.After(lateBound) {
				return statestream.NoDataError
			}
			earlyBound = entry.Timestamp
			if err := cursor.InitWithSnapshot(entry); err != nil {
				return err
			}
			// Send the initial snapshot entry
			ch <- entry
		} else {
			return err
		}
	}

	cursor.SetTimestamp(lateBound)
	if err := cursor.ComputeState(); err != nil {
		return err
	}

	// close the history channel
	sub.Unsubscribe()
	close(ch)

	// if we need to tail, grab the write cursor.
	if !req.Query.Tail {
		return nil
	}

	cursor, err = state.Stream().WriteCursor()
	if err != nil {
		return err
	}
	if !cursor.Ready() {
		if err := cursor.Init(time.Now()); err != nil {
			return err
		}
	}

	// Signal to the client the initial set is done
	if err := srvstream.Send(&view.StateHistoryResponse{Status: view.StateHistoryResponse_HISTORY_TAIL}); err != nil {
		return err
	}

	// Buffer this one, to make sure we don't slow down writing.
	wch := make(chan *statestream.StreamEntry, 100)
	sub = cursor.SubscribeEntries(wch)
	defer func() {
		sub.Unsubscribe()
	}()

	for {
		select {
		case ent := <-wch:
			jsonData, err := json.Marshal(ent.Data)
			if err != nil {
				return err
			}
			srvstream.Send(&view.StateHistoryResponse{
				Status: view.StateHistoryResponse_HISTORY_TAIL,
				State: &view.StateEntry{
					JsonState: string(jsonData),
					Timestamp: util.TimeToNumber(ent.Timestamp),
					Type:      int32(ent.Type),
				},
			})
		case <-srvstream.Context().Done():
			return nil
		}
	}
}
