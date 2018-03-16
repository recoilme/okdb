package api

import (
	"log"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) SayOk(ctx context.Context, in *Ok) (*Ok, error) {
	log.Printf("Receive message %s", in.Message)
	return &Ok{Message: "ok"}, nil
}
