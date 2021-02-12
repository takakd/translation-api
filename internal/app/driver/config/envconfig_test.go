package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		NewEnvConfig()
	})

	_, err := NewEnvConfig("")
	assert.Error(t, err)
}

func TestEnvConfig_Get(t *testing.T) {
	tests := []struct {
		name string
		env  string
		key  string
		want string
	}{
		{name: "value", env: ".env.test", key: "NAME1", want: "value1"},
		{name: "empty value", env: ".env.test", key: "NAME2", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, filePath, _, _ := runtime.Caller(0)
			// e.g. $(pwd)/testdata/.env.test
			envPath := filepath.Join(filepath.Dir(filePath), "./testdata/", tt.env)
			config, err := NewEnvConfig(envPath)
			require.NoError(t, err)

			got := config.Get(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("no env", func(t *testing.T) {
		want := "value3"

		os.Setenv("NAME3", want)

		_, filePath, _, _ := runtime.Caller(0)
		envPath := filepath.Join(filepath.Dir(filePath), "./testdata/.env.test")
		config, err := NewEnvConfig(envPath)
		require.NoError(t, err)

		got := config.Get("NAME3")
		assert.Equal(t, want, got)
	})

	t.Run("overwrite", func(t *testing.T) {
		_, filePath, _, _ := runtime.Caller(0)
		envPath := filepath.Join(filepath.Dir(filePath), "./testdata/.env.test")
		config, err := NewEnvConfig(envPath)
		require.NoError(t, err)

		want := "override"
		os.Setenv("NAME1", want)

		got := config.Get("NAME1")
		assert.Equal(t, want, got)
	})
}
