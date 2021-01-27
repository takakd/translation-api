// Package config provides the application config.
package config

import (
	"api/internal/app/util/config"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// EnvConfig provides the implementation of config.Config based on .env file.
type EnvConfig struct {
}

var _ config.Config = (*EnvConfig)(nil)

// NewEnvConfig creates new struct.
func NewEnvConfig(filepathList ...string) (*EnvConfig, error) {
	if len(filepathList) > 0 {
		if err := godotenv.Load(filepathList...); err != nil {
			return nil, fmt.Errorf(".env loading error %v: %w", filepathList, err)
		}
	}
	return &EnvConfig{}, nil
}

// Get returns value corresponding name.
func (e EnvConfig) Get(name string) (string, error) {
	return os.Getenv(strings.ToUpper(name)), nil
}
