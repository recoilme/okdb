package api

import (
	"github.com/recoilme/slowpoke"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server
type Server struct {
}

// SayOk generates response to a Ping request
func (s *Server) SayOk(ctx context.Context, in *Empty) (*Ok, error) {
	//log.Printf("Receive message")
	return &Ok{Message: "ok"}, nil
}

// Set store key and value in file
func (s *Server) Set(ctx context.Context, cmdSet *CmdSet) (*Empty, error) {
	err := slowpoke.Set(cmdSet.File, cmdSet.Key, cmdSet.Val)
	if err != nil {
		return &Empty{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &Empty{}, nil
}

// Get get value by key
func (s *Server) Get(ctx context.Context, cmdGet *CmdGet) (*ResBytes, error) {
	bytes, err := slowpoke.Get(cmdGet.File, cmdGet.Key)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &ResBytes{Bytes: bytes}, nil
}

// Sets - write key/value pairs -  return error if any
func (s *Server) Sets(ctx context.Context, cmdSets *CmdSets) (*Empty, error) {
	err := slowpoke.Sets(cmdSets.File, cmdSets.Keys)
	if err != nil {
		return &Empty{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &Empty{}, nil
}

// Keys return keys from file
func (s *Server) Keys(ctx context.Context, cmdKeys *CmdKeys) (*ResKeys, error) {
	b, err := slowpoke.Keys(cmdKeys.File, cmdKeys.From, cmdKeys.Limit, cmdKeys.Offset, cmdKeys.Asc)
	if err != nil {
		return &ResKeys{}, status.Errorf(codes.Unknown, err.Error())
	}
	return &ResKeys{Keys: b}, nil
}

// Gets return key/value pairs
func (s *Server) Gets(ctx context.Context, cmdGets *CmdGets) (*ResPairs, error) {
	b := slowpoke.Gets(cmdGets.File, cmdGets.Keys)
	//	slowpoke.Delete()
	return &ResPairs{Pairs: b}, nil
}

// Delete key and val by key
func (s *Server) Delete(ctx context.Context, cmdDel *CmdDel) (*ResDel, error) {
	res, err := slowpoke.Delete(cmdDel.File, cmdDel.Key)
	if err != nil {
		return &ResDel{Deleted: res}, status.Errorf(codes.Unknown, err.Error())
	}
	return &ResDel{Deleted: res}, nil
}
