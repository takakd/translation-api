// Package grpcserver implements a server for Greeter service.
package grpcserver

import (
	"net"

	"api/internal/app/controller/translator"
	pb "api/internal/app/grpc/translator"
	"api/internal/app/util/appcontext"
	"api/internal/app/util/log"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// Server is used to implement translator.TranslatorServer.
type Server struct {
	pb.UnimplementedTranslatorServer
}

// NewServer creates new struct.
func NewServer() *Server {
	return &Server{}
}

// Run start web server.
func (Server) Run() {
	ctx := appcontext.NewContext(context.Background(), "")

	// Create gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to listen: %v", err))
	}
	gs := grpc.NewServer()

	// Set method handlers.
	translatorCtrl := translator.NewController()
	pb.RegisterTranslatorServer(gs, translatorCtrl)

	// Run.
	if err := gs.Serve(lis); err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to serve: %v", err))
	}
}
