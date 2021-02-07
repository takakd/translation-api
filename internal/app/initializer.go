package app

import (
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"api/internal/app/util/di/container/local"
	"api/internal/app/util/di/container/prod"
	"api/internal/app/util/di/container/test"
	"api/internal/app/util/log"
	"errors"
	"fmt"
	"os"
)

// InitDI sets up DI container.
func InitDI() error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return errors.New("APP_ENV empty error")
	}

	var container di.DI

	switch appEnv {
	case "prod":
		container = &prod.Container{}
	case "local":
		container = &local.Container{}
	case "test":
		container = &test.Container{}
	}

	if container == nil {
		return fmt.Errorf("invalid APP_ENV=%s", appEnv)
	}

	di.SetDi(container)

	return nil
}

// InitConfig sets up Config.
// Need to call InitDI before InitConfig.
func InitConfig() error {
	var (
		err error
	)

	var cnf interface{}
	name := "util.config.Config"
	if envFile := os.Getenv("ENV_FILE"); envFile != "" {
		fmt.Printf("env file: %s\n", envFile)
		cnf, err = di.Get(name, envFile)
	} else {
		cnf, err = di.Get(name)
	}
	if err != nil {
		return fmt.Errorf("nil error: %s: %w", name, err)
	}

	config.SetConfig(cnf.(config.Config))

	return nil
}

// InitLogger sets up Logger.
// Need to call InitDI and InitConfig before InitLogger.
func InitLogger() error {
	logger, err := di.Get("util.log.Logger")
	if err != nil {
		return fmt.Errorf("nil error: log.logger")
	}
	log.SetLogger(logger.(log.Logger))

	level := log.LevelDebug
	switch config.Get("DEBUG_LEVEL") {
	case "INFO":
		level = log.LevelInfo
	case "ERROR":
		level = log.LevelError
	}
	log.SetLevel(level)

	return nil
}
