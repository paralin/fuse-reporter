package service

import (
	"golang.org/x/net/context"

	"github.com/fuserobotics/reporter"
	"github.com/fuserobotics/reporter/api"

	"google.golang.org/grpc"
)

type ReporterServiceServer struct {
	Reporter *reporter.Reporter
}

func (s *ReporterServiceServer) RecordState(c context.Context, req *api.RecordStateRequest) (*api.RecordStateResponse, error) {
	if err := req.Validate(); err != nil {
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

func RegisterServer(server *grpc.Server, r *reporter.Reporter) {
	api.RegisterReporterServiceServer(server, &ReporterServiceServer{Reporter: r})
}
