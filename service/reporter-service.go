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

	component, err := s.Reporter.GetComponent(req.Context.Component)
	if err != nil {
		return nil, err
	}
	state, err := component.GetState(req.Context.StateId)
	if err != nil {
		return nil, err
	}
	stream := state.Stream()
	writeCursor, err := stream.WriteCursor()
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
	err = writeCursor.WriteState(reportTime, data, state.Data.StreamConfig.RecordRate)
	if err != nil {
		return nil, err
	}

	return &api.RecordStateResponse{}, nil
}

func (s *ReporterServiceServer) RegisterState(c context.Context, req *api.RegisterStateRequest) (*api.RegisterStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	component, err := s.Reporter.CreateComponentIfNotExists(req.Context.Component)
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
	res := &api.ListStatesResponse{
		List: &api.StateList{
			Components: make([]*api.StateListComponent, len(s.Reporter.ComponentList.Data.ComponentName)),
		},
	}
	for i, componentName := range s.Reporter.ComponentList.Data.ComponentName {
		cmp, err := s.Reporter.GetComponent(componentName)
		if err != nil {
			return nil, err
		}
		slc := &api.StateListComponent{Name: componentName, States: make([]*api.StateListState, len(cmp.Data.StateName))}
		res.List.Components[i] = slc
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

func (s *ReporterServiceServer) GetState(c context.Context, req *api.GetStateRequest) (*api.GetStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	component, err := s.Reporter.GetComponent(req.Context.Component)
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

func RegisterServer(server *grpc.Server, r *reporter.Reporter) {
	api.RegisterReporterServiceServer(server, &ReporterServiceServer{Reporter: r})
}
