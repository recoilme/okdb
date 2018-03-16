package api

import (
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// SayOk generates response to a Ping request
func (s *Server) SayOk(ctx context.Context, in *Empty) (*Ok, error) {
	//log.Printf("Receive message")
	return &Ok{Message: "ok"}, nil
}
