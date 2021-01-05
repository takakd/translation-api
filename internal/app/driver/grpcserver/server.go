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
	"os"
	"github.com/joho/godotenv"
)

const (
	DefaultPort = "50051"
)

// Server is used to implement translator.TranslatorServer.
type Server struct {
	pb.UnimplementedTranslatorServer
	port string
}

// NewServer creates new struct.
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Setup() {
	ctx := appcontext.NewContext(context.Background(), "Server.Init")

	if env := os.Getenv("DOT_ENV"); env != "" {
		dotEnvFilename := ".env." + env
		if err := godotenv.Load(dotEnvFilename); err != nil {
			log.Fatal(ctx, fmt.Errorf(".env loading error: %w", err))
		}
		fmt.Printf("%s loaded\n", dotEnvFilename)
	}

	level := log.LevelDebug
	switch os.Getenv("DEBUG_LEVEL") {
	case "INFO":
		level = log.LevelInfo
	case "ERROR":
		level = log.LevelError
	}
	log.SetLevel(level)

	var port string
	if port = os.Getenv("GRPC_PORT"); port == "" {
		port = DefaultPort
	}
	s.port = fmt.Sprintf(":%s", port)
}

// Run start web server.
func (s *Server) Run() {
	ctx := appcontext.NewContext(context.Background(), "Server.Run")

	fmt.Println(s.port)

	// Create gRPC server.
	lis, err := net.Listen("tcp", s.port)
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
	fmt.Println("test")
}
