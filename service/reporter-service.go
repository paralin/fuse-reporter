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

func RegisterServer(server *grpc.Server, r *reporter.Reporter) {
	api.RegisterReporterServiceServer(server, &ReporterServiceServer{Reporter: r})
}
