// Package config provides config values used in the app.
package config

var (
	// Use this interface for managing config.
	config Config
)

// Config defines methods that returns config values used in the application.
type Config interface {
	Get(name string) string
}

// SetConfig sets config used "config.Get".
func SetConfig(c Config) {
	config = c
}

// Get returns the config value.
func Get(name string) string {
	if config == nil {
		panic("config driver is not set")
	}
	return config.Get(name)
}
