package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/fuserobotics/reporter/api"
	"github.com/fuserobotics/reporter/service"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var RuntimeArgs struct {
	GrpcPort int
	HttpPort int

	ConfigPath string
}

func bindFlags() {
	flag.IntVar(&RuntimeArgs.GrpcPort, "grpcport", 5000, "GRPC port to bind")
	flag.IntVar(&RuntimeArgs.HttpPort, "httpport", 8085, "HTTP port to bind")
	flag.StringVar(&RuntimeArgs.ConfigPath, "config", "", "path to config file")
	flag.CommandLine.Usage = func() {
		fmt.Println(`reporter
Starts the API at the ports specified.
Flags:`)
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()
}

func bindEnv() {
	if ev := os.Getenv("GRPC_PORT"); ev != "" {
		port, err := strconv.Atoi(ev)
		if err != nil {
			fmt.Printf("Couldn't parse env GRPC_PORT (%s), error %v\n", ev, err)
		} else {
			RuntimeArgs.GrpcPort = port
		}
	}
	if ev := os.Getenv("PORT"); ev != "" {
		port, err := strconv.Atoi(ev)
		if err != nil {
			fmt.Printf("Couldn't parse env PORT (%s), error %v\n", ev, err)
		} else {
			RuntimeArgs.HttpPort = port
		}
	}
}

func verifyPort(port int) error {
	if port < 50 || port > 65535 {
		return fmt.Errorf("Port number %d invalid.", port)
	}
	return nil
}

func verifyArgs() error {
	if err := verifyPort(RuntimeArgs.GrpcPort); err != nil {
		return fmt.Errorf("GRPC port invalid: %v", err)
	}
	if err := verifyPort(RuntimeArgs.HttpPort); err != nil {
		return fmt.Errorf("HTTP port invalid: %v", err)
	}

	return nil
}

func runHttpService(endpoint, grpcEndpoint string, ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterReporterServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	glog.Infof("GRPC-Proxy listening on %s", endpoint)
	http.ListenAndServe(endpoint, mux)
	return nil
}

func main() {
	// Log to stdout
	flag.Lookup("logtostderr").Value.Set("true")

	defer func() {
		glog.Info("Exiting...")
	}()
	defer glog.Flush()

	bindFlags()
	bindEnv()
	if err := verifyArgs(); err != nil {
		glog.Fatalf("Error with args: %v\n", err)
	}

	glog.Info("Registering services...")

	grpcServer := grpc.NewServer()
	service.RegisterServer(grpcServer)

	glog.Info("Starting up services...")
	httpEndpoint := fmt.Sprintf("0.0.0.0:%d", RuntimeArgs.HttpPort)
	listenStr := fmt.Sprintf("0.0.0.0:%d", RuntimeArgs.GrpcPort)
	lis, err := net.Listen("tcp", listenStr)
	if err != nil {
		glog.Fatalf("Error listening: %v\n", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		// Setup HTTP service
		if err := runHttpService(httpEndpoint, listenStr, ctx); err != nil {
			glog.Fatal(err)
		}
		defer func() {
			glog.Info("Http service exiting...")
		}()
	}()

	// Start GRPC service
	glog.Infof("grpc listening on %s", listenStr)
	grpcServer.Serve(lis)
}
