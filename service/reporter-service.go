package service

import (
	"errors"
	"time"

	"encoding/json"
	"golang.org/x/net/context"

	"github.com/fuserobotics/reporter"
	"github.com/fuserobotics/reporter/api"
	"github.com/fuserobotics/reporter/util"
	"github.com/fuserobotics/reporter/view"
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
		reportTime = util.NumberToTime(req.Report.Timestamp)
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

func (s *ReporterServiceServer) ListRemotes(c context.Context, req *api.ListRemotesRequest) (*api.ListRemotesResponse, error) {
	res := &api.ListRemotesResponse{
		List: &api.RemoteList{
			Remotes: make(map[string]*view.StateList),
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
