package grpcserver

import (
	"testing"

	"api/internal/app/util/config"
	"api/internal/app/util/log"
	"errors"

	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"api/internal/app/grpc/translator"
	"api/internal/app/util/di"
	"context"
	"sync"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testdata = struct {
	cert        string
	invalidCert string
	key         string
}{
	cert:        "testdata/server.example.crt",
	invalidCert: "testdata/server.invalid.crt",
	key:         "testdata/server.example.key",
}

func TestNewServer(t *testing.T) {
	t.Run("port", func(t *testing.T) {
		tests := []struct {
			name     string
			port     string
			wantPort string
		}{
			{name: "ok", port: "50123", wantPort: "50123"},
			{name: "default", port: "", wantPort: DefaultPort},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return(tt.port)
				mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health")
				mc.EXPECT().Get("TLS").Return("")
				config.SetConfig(mc)

				s, err := NewServer()
				require.NoError(t, err)
				require.NotNil(t, s)
				assert.Equal(t, tt.wantPort, s.port)
			})
		}
	})

	t.Run("tls", func(t *testing.T) {
		testErr := errors.New("error")

		tests := []struct {
			name string
			tls  string
			cert string
			key  string
			err  error
		}{
			{name: "enabled", tls: "true", cert: testdata.cert, key: testdata.key},
			{name: "disabled empty", tls: "false", cert: testdata.cert, key: testdata.key},
			{name: "cert empty", tls: "true", cert: "", key: testdata.key, err: testErr},
			{name: "cert not exists", tls: "true", cert: "dummy", key: testdata.key, err: testErr},
			{name: "key empty", tls: "true", cert: testdata.cert, key: "", err: testErr},
			{name: "key not exists", tls: "true", cert: testdata.cert, key: "dummy", err: testErr},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("50051")
				mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health")

				mc.EXPECT().Get("TLS").Return(tt.tls)
				if tt.tls == "true" {
					mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(tt.cert)
					mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(tt.key)
				}
				config.SetConfig(mc)

				s, err := NewServer()
				if tt.err != nil {
					assert.Nil(t, s)
					assert.Error(t, err)
				} else {
					assert.NotNil(t, s)
					assert.NoError(t, err)

					if tt.tls == "true" {
						assert.Equal(t, true, s.tlsEnabled)
					} else {
						assert.Equal(t, false, s.tlsEnabled)
					}
				}
			})
		}
	})

	t.Run("health", func(t *testing.T) {
		tests := []struct {
			name string
			path string
			err  error
		}{
			{name: "ok", path: "/health"},
			{name: "empty", path: ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("50051")
				mc.EXPECT().Get("HEALTH_CHECK_PATH").Return(tt.path)
				if tt.path != "" {
					mc.EXPECT().Get("TLS").Return("false")
				}
				config.SetConfig(mc)

				s, err := NewServer()
				if tt.path != "" {
					assert.NotNil(t, s)
					assert.NoError(t, err)
				} else {
					assert.Nil(t, s)
					assert.Error(t, err)
				}
			})
		}
	})
}

func TestServer_ServeHTTP(t *testing.T) {
	t.Run("health check", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().Info(gomock.Any(), gomock.Any())
		log.SetLogger(ml)

		s := &Server{
			healthCheckPath: "/health",
		}

		r := httptest.NewRequest(http.MethodGet, s.healthCheckPath, bytes.NewReader(nil))
		w := httptest.NewRecorder()
		s.ServeHTTP(w, r)

		ret := w.Result()
		assert.Equal(t, http.StatusOK, ret.StatusCode)

		body, err := ioutil.ReadAll(ret.Body)
		assert.NoError(t, err)
		assert.Equal(t, "OK", string(body))
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("tcp error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GRPC_PORT").Return("")
		mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health")
		mc.EXPECT().Get("TLS").Return("false")
		config.SetConfig(mc)

		s, _ := NewServer()
		s.port = "invalid port"
		err := s.Run()
		assert.Error(t, err)
	})

	t.Run("translator controller error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("controller.translator.Controller").Return(nil, errors.New("controller error"))
		di.SetDi(md)

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GRPC_PORT").Return("50051")
		mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health")
		mc.EXPECT().Get("TLS").Return("false")
		config.SetConfig(mc)

		s, _ := NewServer()
		err := s.Run()
		assert.Error(t, err)
	})

	t.Run("run translator", func(t *testing.T) {
		tests := []struct {
			name string
			tls  string
			err  bool
		}{
			{name: "ok", tls: "false"},
			{name: "ok tls", tls: "true"},
			{name: "error", tls: "false", err: true},
			{name: "error tls", tls: "true", err: true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mt := translator.NewMockTranslatorServer(ctrl)

				md := di.NewMockDI(ctrl)
				md.EXPECT().Get("controller.translator.Controller").Return(mt, nil)
				di.SetDi(md)

				ml := log.NewMockLogger(ctrl)
				if tt.tls == "true" {
					ml.EXPECT().Info(gomock.Any(), log.StringValue("server start: tls=on port=50051"))
				} else {
					ml.EXPECT().Info(gomock.Any(), log.StringValue("server start: tls=off port=50051"))
				}
				if !tt.err {
					ml.EXPECT().Info(gomock.Any(), log.StringValue("server end"))
				}
				log.SetLogger(ml)

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("")
				mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health")
				mc.EXPECT().Get("TLS").Return(tt.tls)

				if tt.tls == "true" {
					if tt.err {
						mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.invalidCert)
					} else {
						mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.cert)
					}
					mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(testdata.key)
				}

				config.SetConfig(mc)

				s, _ := NewServer()

				// Run server in Go routine and force to stop.
				var wg sync.WaitGroup
				chErr := make(chan error)

				wg.Add(1)
				go func(s *Server, chErr chan<- error) {
					wg.Done()
					err := s.Run()
					chErr <- err
					close(chErr)
				}(s, chErr)

				wg.Wait()

				// Wait to run server.
				time.Sleep(1 * time.Second)

				if tt.err {
					s.ln.Close()
				} else {
					s.srv.Shutdown(context.TODO())
				}

				// Wait to receive error in chErr.
				time.Sleep(1 * time.Second)

				for err := range chErr {
					if tt.err {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				}
			})
		}
	})
}
