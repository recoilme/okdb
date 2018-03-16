package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/recoilme/okdb/api"
	"google.golang.org/grpc"
)

// main start a gRPC server and waits for connection
// run from root: protoc -I api/ -I${GOPATH}/src --go_out=plugins=grpc:api api/api.proto
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := api.Server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	api.RegisterOkdbServer(grpcServer, &s)

	// start the server
	//if err := grpcServer.Serve(lis); err != nil {
	//log.Fatalf("failed to serve: %s", err)
	//}

	go func() {
		//grpcServer.Serve(lis)
		log.Fatal(grpcServer.Serve(lis))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	//grpcServer.Shutdown(context.Background())

}
