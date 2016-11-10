package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fuserobotics/reporter"
	"github.com/fuserobotics/reporter/api"
	"github.com/fuserobotics/reporter/service"
	"github.com/fuserobotics/reporter/view"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var RuntimeArgs struct {
	GrpcPort       int
	HttpPort       int
	DbPath         string
	HostIdentifier string
	DisableRemotes bool
}

var reporterInstance *reporter.Reporter

func bindFlags() {
	flag.IntVar(&RuntimeArgs.GrpcPort, "grpcport", 5000, "GRPC port to bind")
	flag.IntVar(&RuntimeArgs.HttpPort, "httpport", 8085, "HTTP port to bind")
	flag.StringVar(&RuntimeArgs.DbPath, "dbpath", "reporter.db", "Database path")
	flag.StringVar(&RuntimeArgs.HostIdentifier, "ident", "", "Host identifier")
	flag.BoolVar(&RuntimeArgs.DisableRemotes, "noremotes", false, "Flag to disable pushing to remotes.")
	flag.CommandLine.Usage = func() {
		fmt.Println(`reporter
Starts the API at the ports specified.
Flags:`)
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()
}

func initReporter() error {
	ri, err := reporter.NewReporter(RuntimeArgs.HostIdentifier, RuntimeArgs.DbPath, !RuntimeArgs.DisableRemotes)
	if err != nil {
		return err
	}
	reporterInstance = ri
	return nil
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
	if RuntimeArgs.HostIdentifier == "" {
		return fmt.Errorf("Host identifier must be specified.")
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
	err = view.RegisterReporterServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
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
	if err := initReporter(); err != nil {
		glog.Fatalf("Error initing reporter: %v\n", err)
	}

	glog.Info("Registering services...")

	grpcServer := grpc.NewServer()
	service.RegisterServer(grpcServer, reporterInstance)

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Setup HTTP service
		if err := runHttpService(httpEndpoint, listenStr, ctx); err != nil {
			glog.Fatal(err)
		}
		defer func() {
			glog.Info("Http service exiting...")
		}()
	}()

	go func() {
		// Start GRPC service
		glog.Infof("grpc listening on %s", listenStr)
		grpcServer.Serve(lis)
	}()

	<-sigs

	glog.Info("Exiting...")
	grpcServer.GracefulStop()
}
