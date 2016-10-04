package service

import (
	"golang.org/x/net/context"

	"github.com/fuserobotics/reporter/api"
	"google.golang.org/grpc"
)

type ReporterServiceServer struct {
}

func (s *ReporterServiceServer) RecordState(c context.Context, req *api.RecordStateRequest) (*api.RecordStateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &api.RecordStateResponse{}, nil
}

func RegisterServer(server *grpc.Server) {
	api.RegisterReporterServiceServer(server, &ReporterServiceServer{})
}
