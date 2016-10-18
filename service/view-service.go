package service

import (
	"encoding/json"
	"golang.org/x/net/context"

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
	if err := req.Validate(); err != nil {
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
