// Package grpcserver implements a server for Greeter service.
package grpcserver

import (
	pbt "api/internal/app/grpc/translator"
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"net/http"

	"api/internal/pkg/util"

	"google.golang.org/grpc"
)

const (
	// DefaultPort is used as server port if environment variables does not exists.
	DefaultPort = "50051"
)

// Server is used to implement translator.TranslatorServer.
type Server struct {
	pbt.UnimplementedTranslatorServer
	gs              *grpc.Server
	srv             *http.Server
	port            string
	certFilePath    string
	keyFilePath     string
	healthCheckPath string
}

// NewServer creates new struct.
func NewServer() (*Server, error) {
	s := &Server{}

	var err error

	s.port, err = config.Get("GRPC_PORT")
	if err != nil {
		return nil, fmt.Errorf("port error: %w", err)
	}
	if s.port == "" {
		s.port = DefaultPort
	}

	s.certFilePath, err = config.Get("SERVER_CERT_FILE_PATH")
	if err != nil {
		return nil, fmt.Errorf("server cert file error: %w", err)
	} else if !util.FileExists(s.certFilePath) {
		return nil, fmt.Errorf("server cert file not exists: path=%s", s.certFilePath)
	}

	s.keyFilePath, err = config.Get("SERVER_KEY_FILE_PATH")
	if err != nil {
		return nil, fmt.Errorf("server key file error: %w", err)
	} else if !util.FileExists(s.keyFilePath) {
		return nil, fmt.Errorf("server key file not exists: path=%s", s.certFilePath)
	}

	s.healthCheckPath, err = config.Get("HEALTH_CHECK_PATH")
	if err != nil {
		return nil, fmt.Errorf("health check path error: %w", err)
	} else if s.healthCheckPath == "" {
		return nil, fmt.Errorf("health check path empty error")
	}

	return s, nil
}

// ServeHTTP handles requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handles health check.
	if r.URL.Path == s.healthCheckPath && r.Method == http.MethodGet {
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handles gRPC response.
	s.gs.ServeHTTP(w, r)
}

// Run starts grpc server.
func (s *Server) Run() error {
	var err error

	// Setup gRPC handlers.
	s.gs = grpc.NewServer()
	translatorCtrl, err := di.Get("controller.translator.Controller")
	if err != nil {
		return fmt.Errorf("failed to setup translator handler: %v", err)
	}
	pbt.RegisterTranslatorServer(s.gs, translatorCtrl.(pbt.TranslatorServer))

	// Run
	log.Info(context.Background(), log.StringValue("server start"))
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: s,
	}
	if err := s.srv.ListenAndServeTLS(s.certFilePath, s.keyFilePath); err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %v", err)
	}

	log.Info(context.Background(), log.StringValue("server end"))
	return nil
}
