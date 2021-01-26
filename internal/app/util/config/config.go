// Package config provides config values used in the app.
package config

import (
	"fmt"
)

var (
	// Use this interface for managing config.
	config Config
)

// Config defines methods that returns config values used in the application.
type Config interface {
	Get(name string) (string, error)
}

// SetConfig sets config used "config.Get".
func SetConfig(c Config) {
	config = c
}

// Get returns the config value.
func Get(name string) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config is null")
	}
	return config.Get(name)
}
