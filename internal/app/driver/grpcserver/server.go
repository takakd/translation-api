// Package grpcserver implements a server for Greeter service.
package grpcserver

import (
	pbt "api/internal/app/grpc/translator"
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"api/internal/app/util/log"
	"api/internal/pkg/util"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	// DefaultPort is used as server port if environment variables does not exists.
	DefaultPort = "50051"
)

// Server is used to implement translator.TranslatorServer.
type Server struct {
	pbt.UnimplementedTranslatorServer
	gs           *grpc.Server
	srv          *http.Server
	port         string
	certFilePath string
	keyFilePath  string

	healthCheckPath string
	ln              net.Listener
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
	// Set request ID.
	ctx := log.WithLogContextValue(context.Background(), uuid.New().String())

	// Set context to the request.
	r = r.WithContext(ctx)

	// Handles health check.
	if r.URL.Path == s.healthCheckPath && r.Method == http.MethodGet {
		// Access log
		now := time.Now()
		log.Info(ctx, log.Value{
			"header": r.Header,
			"host":   r.Host,
			"date":   now.Format(time.RFC3339),
		})

		w.Write([]byte("OK"))
		return
	}

	// Handles gRPC response.
	s.gs.ServeHTTP(w, r)
}

// Run starts grpc server.
func (s *Server) Run() error {
	var err error

	s.gs = grpc.NewServer()

	// Setup gRPC handlers.
	translatorCtrl, err := di.Get("controller.translator.Controller")
	if err != nil {
		return fmt.Errorf("failed to setup translator handler: %v", err)
	}
	pbt.RegisterTranslatorServer(s.gs, translatorCtrl.(pbt.TranslatorServer))

	log.Info(context.Background(), log.StringValue(fmt.Sprintf("server start: port=%s", s.port)))

	// Create TCP listener
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: s,
	}
	s.ln, err = net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("tcp listener creation error: %w", err)
	}
	defer s.ln.Close()

	// Run
	if err := s.srv.ServeTLS(s.ln, s.certFilePath, s.keyFilePath); err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %v", err)
	}

	log.Info(context.Background(), log.StringValue("server end"))
	return nil
}
