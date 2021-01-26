package grpcserver

import (
	"api/internal/app/util/di"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"api/internal/app/util/config"
	"errors"
	"api/internal/app/util/log"
	"api/internal/app/grpc/translator"
	"time"
	"sync"
)

func TestNewServer(t *testing.T) {
	tt := NewServer()
	assert.Equal(t, ":"+DefaultPort, tt.port)
}

func TestServer_setupConfig(t *testing.T) {
	tests := []struct {
		name string
		env  string
		err  error
	}{
		{name: "env", env: "test"},
		{name: "env empty", env: ""},
		{name: "error", env: "", err: errors.New("config error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Test value after setupConfig
			cnfName := "did_set"
			cnfValue := "ok"

			mc := config.NewMockConfig(ctrl)
			if tt.err == nil {
				mc.EXPECT().Get(cnfName).Return(cnfValue, nil)
			}

			md := di.NewMockDI(ctrl)
			if tt.env == "" {
				md.EXPECT().Get("config.config").Return(mc, tt.err)
			} else {
				envFilename := ".env." + tt.env
				md.EXPECT().Get("config.config", []string{envFilename}).Return(mc, tt.err)
			}

			di.SetDi(md)
			os.Setenv("DOT_ENV", tt.env)

			s := NewServer()
			err := s.setupConfig()

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check config
				v, _ := config.Get(cnfName)
				assert.Equal(t, cnfValue, v)
			}
		})
	}
}

func TestServer_setupLogger(t *testing.T) {
	tests := []struct {
		name       string
		levelLabel string
		level      log.Level
		err        error
	}{
		{name: "debug", levelLabel: "", level: log.LevelDebug},
		{name: "error", levelLabel: "ERROR", level: log.LevelError},
		{name: "info", levelLabel: "INFO", level: log.LevelInfo},
		{name: "err", levelLabel: "DEBUG", level: log.LevelDebug, err: errors.New("log error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ml := log.NewMockLogger(ctrl)
			if tt.err == nil {
				ml.EXPECT().SetLevel(tt.level)
			}

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("log.logger").Return(ml, tt.err)

			di.SetDi(md)
			os.Setenv("DEBUG_LEVEL", tt.levelLabel)

			s := NewServer()
			err := s.setupLogger()
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServer_Setup(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		wantPort string
	}{
		{name: "port", port: "50011", wantPort: ":50011"},
		{name: "default port", wantPort: ":" + DefaultPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mc := config.NewMockConfig(ctrl)

			ml := log.NewMockLogger(ctrl)
			ml.EXPECT().SetLevel(gomock.Any())

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("config.config").Return(mc, nil)
			md.EXPECT().Get("log.logger").Return(ml, nil)

			di.SetDi(md)
			os.Setenv("GRPC_PORT", tt.port)

			s := NewServer()
			err := s.Setup()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantPort, s.port)
		})
	}

	t.Run("config error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("config.config").Return(nil, errors.New("log error"))

		di.SetDi(md)

		s := NewServer()
		err := s.Setup()
		assert.Error(t, err)
	})

	t.Run("log error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("config.config").Return(nil, errors.New("log error"))

		di.SetDi(md)

		s := NewServer()
		err := s.Setup()
		assert.Error(t, err)
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("tcp error", func(t *testing.T) {
		s := NewServer()
		s.port = "invalid port"
		err := s.Run()
		assert.Error(t, err)
	})

	t.Run("controller error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("translator.NewController").Return(nil, errors.New("controller error"))
		di.SetDi(md)

		s := NewServer()
		err := s.Run()
		assert.Error(t, err)
	})

	t.Run("run", func(t *testing.T) {
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

				mc := translator.NewMockTranslatorServer(ctrl)
				md := di.NewMockDI(ctrl)
				md.EXPECT().Get("translator.NewController").Return(mc, nil)
				di.SetDi(md)

				s := NewServer()

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
					// Close listener to raise error.
					s.lis.Close()
				} else {
					s.gs.Stop()
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
