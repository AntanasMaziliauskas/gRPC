package handler

import (
	"log"

	"github.com/AntanasMaziliauskas/grpc/api"
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) FindData(ctx context.Context, in *api.LookFor) (*api.Found, error) {
	log.Printf("Received request to look for: %s", in.Id)
	return &api.Found{Id: "bar", Name: "Jonas", Age: 20}, nil
}
