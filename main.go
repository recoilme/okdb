package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/recoilme/slowpoke"

	"github.com/recoilme/okdb/api"
	"google.golang.org/grpc"
)

// main start a gRPC server and waits for connection
// run from root: protoc -I api/ -I${GOPATH}/src --go_out=plugins=grpc:api api/api.proto
func main() {

	grpcPort := flag.Int("grpc", 7777, "grpc port number, default 7777")
	httpAddr := flag.String("http", ":5000", "http api address, default :5000")
	flag.Parse()
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
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
	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	// start http server
	Serve(*httpAddr)

	// handle kill
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	closeErr := slowpoke.CloseAll()
	if closeErr != nil {
		log.Fatal(closeErr.Error())
	}

}

// Serve run server
// example addr: ":5000"
// example usage ./simpleserver :5000>>simpleserver.log &
func Serve(addr string) {
	http.HandleFunc("/slowpoke/", handlerSlowPoke)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handlerSlowPoke(w http.ResponseWriter, r *http.Request) {
	api.Parser(w, r)
}
