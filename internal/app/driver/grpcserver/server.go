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
	gs              *grpc.Server
	srv             *http.Server
	port            string
	healthCheckPath string
	ln              net.Listener

	tlsEnabled   bool
	certFilePath string
	keyFilePath  string
}

// NewServer creates new server.
func NewServer() (*Server, error) {
	s := &Server{}

	s.port = config.Get("GRPC_PORT")
	if s.port == "" {
		s.port = DefaultPort
	}

	s.healthCheckPath = config.Get("HEALTH_CHECK_PATH")
	if s.healthCheckPath == "" {
		return nil, fmt.Errorf("health check path empty error")
	}

	if tls := config.Get("TLS"); tls == "true" {
		s.tlsEnabled = true
	} else {
		s.tlsEnabled = false
	}

	if s.tlsEnabled {
		s.certFilePath = config.Get("SERVER_CERT_FILE_PATH")
		s.keyFilePath = config.Get("SERVER_KEY_FILE_PATH")
		if !util.FileExists(s.certFilePath) || !util.FileExists(s.keyFilePath) {
			return nil, fmt.Errorf("server cert files not exists: cert=%s, key=%s", s.certFilePath, s.keyFilePath)
		}
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
			"tag":    "health-check",
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
	if s.tlsEnabled {
		log.Info(context.Background(), log.StringValue(fmt.Sprintf("server start: tls=on port=%s", s.port)))
		if err := s.srv.ServeTLS(s.ln, s.certFilePath, s.keyFilePath); err != http.ErrServerClosed {
			return fmt.Errorf("failed to serve: %v", err)
		}
	} else {
		log.Info(context.Background(), log.StringValue(fmt.Sprintf("server start: tls=off port=%s", s.port)))
		if err := s.srv.Serve(s.ln); err != http.ErrServerClosed {
			return fmt.Errorf("failed to serve: %v", err)
		}
	}

	log.Info(context.Background(), log.StringValue("server end"))
	return nil
}
