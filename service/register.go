package service

import (
	"google.golang.org/grpc"

	"github.com/fuserobotics/reporter"
	"github.com/fuserobotics/reporter/api"
	"github.com/fuserobotics/reporter/view"
)

func RegisterServer(server *grpc.Server, r *reporter.Reporter) {
	api.RegisterReporterServiceServer(server, &ReporterServiceServer{Reporter: r})
	view.RegisterReporterServiceServer(server, &ViewServiceServer{Reporter: r})
}
