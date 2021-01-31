// Package grpcserver implements a server for Greeter service.
package grpcserver

import (
	"net"

	pb "api/internal/app/grpc/translator"
	"fmt"

	"api/internal/app/util/di"

	"api/internal/app/util/config"

	"google.golang.org/grpc"
)

const (
	// DefaultPort is used as server port if environment variables does not exists.
	DefaultPort = "50051"
)

// Server is used to implement translator.TranslatorServer.
type Server struct {
	pb.UnimplementedTranslatorServer
	port string
	gs   *grpc.Server
	lis  net.Listener
}

// NewServer creates new struct.
func NewServer() (*Server, error) {
	// Set port number.
	port, err := config.Get("GRPC_PORT")
	if err != nil {
		return nil, fmt.Errorf("port error: %w", err)
	}
	if port == "" {
		port = DefaultPort
	}

	return &Server{
		port: fmt.Sprintf(":%s", port),
	}, nil
}

// Run starts grpc server.
func (s *Server) Run() error {
	var err error

	// Create gRPC server.
	s.lis, err = net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	defer s.lis.Close()

	s.gs = grpc.NewServer()

	// Set method handlers.
	translatorCtrl, err := di.Get("controller.translator.Controller")
	if err != nil {
		return fmt.Errorf("failed to setup api: %v", err)
	}
	pb.RegisterTranslatorServer(s.gs, translatorCtrl.(pb.TranslatorServer))

	// Run.
	if err := s.gs.Serve(s.lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	fmt.Println("serve end")
	return nil
}
