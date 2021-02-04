package grpcserver

import (
	"testing"

	"api/internal/app/grpc/translator"
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"errors"
	"sync"
	"time"

	"api/internal/app/grpc/health/grpc_health_v1"
	"api/internal/app/util/log"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
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
			config.SetConfig(mc)

			s, err := NewServer()

			if tt.portErr != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, ":"+tt.wantPort, s.port)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	t.Run("tcp error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GRPC_PORT").Return("", nil)
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
		mc.EXPECT().Get("GRPC_PORT").Return("", nil)
		config.SetConfig(mc)

		s, _ := NewServer()
		err := s.Run()
		assert.Error(t, err)
	})

	t.Run("health check controller error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mt := translator.NewMockTranslatorServer(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("controller.translator.Controller").Return(mt, nil)
		md.EXPECT().Get("controller.grpc_health_v1.Controller").Return(nil, errors.New("controller error"))
		di.SetDi(md)

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GRPC_PORT").Return("", nil)
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
				mh := grpc_health_v1.NewMockHealthServer(ctrl)

				md := di.NewMockDI(ctrl)
				md.EXPECT().Get("controller.translator.Controller").Return(mt, nil)
				md.EXPECT().Get("controller.grpc_health_v1.Controller").Return(mh, nil)
				di.SetDi(md)

				ml := log.NewMockLogger(ctrl)
				ml.EXPECT().Info(gomock.Any(), log.StringValue("server start"))
				log.SetLogger(ml)

				mc := config.NewMockConfig(ctrl)
				mc.EXPECT().Get("GRPC_PORT").Return("", nil)
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
