package service

import (
	"errors"
	"time"

	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/fuserobotics/reporter"
	"github.com/fuserobotics/reporter/api"
	statestream "github.com/fuserobotics/statestream"
)

type ReporterServiceServer struct {
	Reporter *reporter.Reporter
}

func (s *ReporterServiceServer) RecordState(c context.Context, req *api.RecordStateRequest) (*api.RecordStateResponse, error) {
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

	var reportTime time.Time
	if req.Report.Timestamp > 0 {
		reportTime = reporter.NumberToTime(req.Report.Timestamp)
	} else {
		reportTime = time.Now()
	}

	var dataMap map[string]interface{}
	var data statestream.StateData

	if req.Report.JsonState != "" {
		if err := json.Unmarshal([]byte(req.Report.JsonState), &dataMap); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("No data specified in the report.")
	}

	data = statestream.StateData(dataMap)

	err = state.WriteState(reportTime, data)
	if err != nil {
		return nil, err
	}

	return &api.RecordStateResponse{}, nil
}

func (s *ReporterServiceServer) RegisterState(c context.Context, req *api.RegisterStateRequest) (*api.RegisterStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	component, err := s.Reporter.LocalTree.CreateComponentIfNotExists(req.Context.Component)
	if err != nil {
		return nil, err
	}
	_, err = component.CreateStateIfNotExists(req.Context.StateId, req.StreamConfig)
	if err != nil {
		return nil, err
	}
	return &api.RegisterStateResponse{}, nil
}

func (s *ReporterServiceServer) ListStates(c context.Context, req *api.ListStatesRequest) (*api.ListStatesResponse, error) {
	list, err := buildStateList(s.Reporter.LocalTree)
	if err != nil {
		return nil, err
	}
	return &api.ListStatesResponse{
		List: list,
	}, nil
}

func (s *ReporterServiceServer) GetState(c context.Context, req *api.GetStateRequest) (*api.GetStateResponse, error) {
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
		if err := cursor.Init(reporter.NumberToTime(req.Query.Time)); err != nil {
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

	return &api.GetStateResponse{
		State: &api.StateReport{
			JsonState: string(respData),
			Timestamp: reporter.TimeToNumber(cursor.ComputedTimestamp()),
		},
	}, nil
}

func buildStateList(tree *reporter.ComponentTree) (*api.StateList, error) {
	res := &api.StateList{
		Components: make([]*api.StateListComponent, len(tree.ComponentList.Data.ComponentName)),
	}

	for i, componentName := range tree.ComponentList.Data.ComponentName {
		cmp, err := tree.GetComponent(componentName)
		if err != nil {
			return nil, err
		}
		slc := &api.StateListComponent{Name: componentName, States: make([]*api.StateListState, len(cmp.Data.StateName))}
		res.Components[i] = slc
		for is, stateName := range cmp.Data.StateName {
			state, err := cmp.GetState(stateName)
			if err != nil {
				return nil, err
			}
			slc.States[is] = &api.StateListState{
				Name:   stateName,
				Config: state.Data.StreamConfig,
			}
		}
	}
	return res, nil
}

func (s *ReporterServiceServer) ListRemotes(c context.Context, req *api.ListRemotesRequest) (*api.ListRemotesResponse, error) {
	res := &api.ListRemotesResponse{
		List: &api.RemoteList{
			Remotes: make(map[string]*api.StateList),
		},
	}
	for id, remote := range s.Reporter.RemoteList.Remotes {
		sl, err := buildStateList(remote.ComponentTree)
		if err != nil {
			return nil, err
		}
		res.List.Remotes[id] = sl
	}
	return res, nil
}

func (s *ReporterServiceServer) CreateRemote(c context.Context, req *api.CreateRemoteRequest) (*api.CreateRemoteResponse, error) {
	if req.Context == nil || req.Context.RemoteId == "" {
		return nil, errors.New("Must supply context in request.")
	}

	_, err := s.Reporter.RemoteList.CreateOrUpdateRemote(req.Context.RemoteId, req.Endpoint)
	if err != nil {
		return nil, err
	}
	return &api.CreateRemoteResponse{}, nil
}

func RegisterServer(server *grpc.Server, r *reporter.Reporter) {
	api.RegisterReporterServiceServer(server, &ReporterServiceServer{Reporter: r})
}
