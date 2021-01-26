// Package grpcserver implements a server for Greeter service.
package grpcserver

import (
	"net"

	pb "api/internal/app/grpc/translator"
	"api/internal/app/util/log"
	"fmt"

	"os"

	"api/internal/app/util/di"

	"google.golang.org/grpc"
	"api/internal/app/util/config"
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
func NewServer() *Server {
	return &Server{
		port: fmt.Sprintf(":%s", DefaultPort),
	}
}

// Setup initialize server.
func (s *Server) Setup() error {
	if err := s.setupConfig(); err != nil {
		return fmt.Errorf("config setup error: %w", err)
	}

	if err := s.setupLogger(); err != nil {
		return fmt.Errorf("logger setup error: %w", err)
	}

	// Set port number
	var port string
	if port = os.Getenv("GRPC_PORT"); port != "" {
		s.port = fmt.Sprintf(":%s", port)
	}

	return nil
}

func (s *Server) setupConfig() error {
	var (
		err         error
		envFilename string
	)

	var cnf interface{}
	if env := os.Getenv("DOT_ENV"); env != "" {
		envFilename = ".env." + env
		//if err := godotenv.Load(dotEnvFilename); err != nil {
		//	return fmt.Errorf(".env loading error: %w", err)
		//}
		//fmt.Printf("%s loaded\n", dotEnvFilename)

		cnf, err = di.Get("config.config", []string{envFilename})
	} else {
		cnf, err = di.Get("config.config")
	}
	if err != nil {
		return fmt.Errorf("nil error: config.config: %w", err)
	}

	config.SetConfig(cnf.(config.Config))

	return nil
}

func (s *Server) setupLogger() error {
	logger, err := di.Get("log.logger")
	if err != nil {
		return fmt.Errorf("nil error: log.logger")
	}
	log.SetLogger(logger.(log.Logger))

	level := log.LevelDebug
	switch os.Getenv("DEBUG_LEVEL") {
	case "INFO":
		level = log.LevelInfo
	case "ERROR":
		level = log.LevelError
	}
	log.SetLevel(level)

	return nil
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
	translatorCtrl, err := di.Get("translator.NewController")
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
