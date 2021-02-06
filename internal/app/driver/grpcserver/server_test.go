package grpcserver

import (
	"testing"

	"api/internal/app/grpc/translator"
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"errors"
	"sync"
	"time"

	"api/internal/app/util/log"

	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testdata = struct {
	cert        string
	invalidCert string
	key         string
}{
	cert:        "testdata/server.crt",
	invalidCert: "testdata/server.invalid.crt",
	key:         "testdata/server.key",
}

func TestNewServer(t *testing.T) {
	t.Run("port", func(t *testing.T) {
		tests := []struct {
			name     string
			port     string
			wantPort string
			portErr  error
		}{
			{name: "ok", port: "50123", wantPort: "50123"},
			{name: "default", port: "", wantPort: DefaultPort},
			{name: "error", port: "50123", portErr: errors.New("error")},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return(tt.port, tt.portErr)
				if tt.portErr == nil {
					mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.cert, nil)
					mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(testdata.key, nil)
					mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health", nil)
				}

				config.SetConfig(mc)

				s, err := NewServer()

				if tt.portErr != nil {
					require.Nil(t, s)
					assert.Error(t, err)
				} else {
					require.NotNil(t, s)
					assert.Equal(t, tt.wantPort, s.port)
				}
			})
		}
	})

	t.Run("files", func(t *testing.T) {
		testErr := errors.New("error")

		tests := []struct {
			name    string
			cert    string
			certErr error
			key     string
			keyErr  error
			err     error
		}{
			{name: "ok", cert: testdata.cert, key: testdata.key},
			{name: "cert empty", cert: "", key: testdata.key, certErr: testErr, err: testErr},
			{name: "cert not exists", cert: "dummy", key: testdata.key, certErr: testErr, err: testErr},
			{name: "key empty", cert: testdata.cert, key: "", keyErr: testErr, err: testErr},
			{name: "key not exists", cert: testdata.cert, key: "dummy", keyErr: testErr, err: testErr},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("50051", nil)
				mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(tt.cert, tt.certErr)
				if tt.certErr == nil {
					mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(tt.key, tt.keyErr)
				}
				if tt.certErr == nil && tt.keyErr == nil {
					mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("ok", nil)
				}

				config.SetConfig(mc)

				s, err := NewServer()
				if tt.err == nil {
					assert.NotNil(t, s)
					assert.NoError(t, err)
				} else {
					assert.Nil(t, s)
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("health", func(t *testing.T) {
		tests := []struct {
			name string
			path string
		}{
			{name: "ok", path: "/health"},
			{name: "health check error"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				err := errors.New("config error")

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("50051", nil)
				mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.cert, nil)
				mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(testdata.key, nil)
				mc.EXPECT().Get("HEALTH_CHECK_PATH").Return(tt.path, nil)

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
		mc.EXPECT().Get("GRPC_PORT").Return("", nil)
		mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.cert, nil)
		mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(testdata.key, nil)
		mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health", nil)
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
		mc.EXPECT().Get("GRPC_PORT").Return("50051", nil)
		mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.cert, nil)
		mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(testdata.key, nil)
		mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health", nil)
		config.SetConfig(mc)

		s, _ := NewServer()
		err := s.Run()
		assert.Error(t, err)
	})

	t.Run("run translator", func(t *testing.T) {
		tests := []struct {
			name string
			err  bool
		}{
			{name: "ok"},
			{name: "error", err: true},
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
				ml.EXPECT().Info(gomock.Any(), log.StringValue("server start"))
				if !tt.err {
					ml.EXPECT().Info(gomock.Any(), log.StringValue("server end"))
				}
				log.SetLogger(ml)

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("", nil)
				mc.EXPECT().Get("HEALTH_CHECK_PATH").Return("/health", nil)
				if tt.err {
					mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.invalidCert, nil)
				} else {
					mc.EXPECT().Get("SERVER_CERT_FILE_PATH").Return(testdata.cert, nil)
				}
				mc.EXPECT().Get("SERVER_KEY_FILE_PATH").Return(testdata.key, nil)
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

				if !tt.err {
					s.srv.Close()
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
