package app

import (
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"api/internal/app/util/log"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInitDI(t *testing.T) {
	tests := []struct {
		name      string
		appEnv    string
		container di.DI
		err       bool
	}{
		{name: "empty", appEnv: "", err: true},
		{name: "prod", appEnv: "prod", err: false},
		{name: "test", appEnv: "test", err: false},
		{name: "local", appEnv: "local", err: false},
		{name: "env invalid", appEnv: "invalid", err: true},
	}
	for _, tt := range tests {
		os.Setenv("APP_ENV", tt.appEnv)
		err := InitDI()

		if tt.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		err     error
	}{
		{name: "env", envFile: "/dummy/.env.test"},
		{name: "env empty", envFile: ""},
		{name: "error", envFile: "", err: errors.New("config error")},
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
			if tt.envFile == "" {
				md.EXPECT().Get("util.config.Config").Return(mc, tt.err)
			} else {
				md.EXPECT().Get("util.config.Config", []string{tt.envFile}).Return(mc, tt.err)
			}

			di.SetDi(md)
			os.Setenv("ENV_FILE", tt.envFile)

			err := InitConfig()

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

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name       string
		levelLabel string
		level      log.Level
		err        error
	}{
		{name: "debug", levelLabel: "", level: log.LevelDebug},
		{name: "error", levelLabel: "ERROR", level: log.LevelError},
		{name: "info", levelLabel: "INFO", level: log.LevelInfo},
		{name: "set error", levelLabel: "INFO", level: log.LevelInfo, err: errors.New("error")},
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
			md.EXPECT().Get("util.log.Logger").Return(ml, nil)
			di.SetDi(md)

			mc := config.NewMockConfig(ctrl)
			mc.EXPECT().Get("DEBUG_LEVEL").Return(tt.levelLabel, tt.err)
			config.SetConfig(mc)

			err := InitLogger()
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("logger error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("util.log.Logger").Return(nil, errors.New("error"))
		di.SetDi(md)

		err := InitLogger()
		assert.Error(t, err)
	})
}
