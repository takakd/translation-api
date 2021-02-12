package google

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("api key error", func(t *testing.T) {
		tests := []struct {
			name string
			path string
		}{
			{name: "empty", path: ""},
			{name: "not exists", path: "not_exists.json"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := NewClient(context.TODO(), tt.path)
				assert.Nil(t, got)
				assert.Error(t, err)
			})
		}
	})

	t.Run("ok", func(t *testing.T) {
		fmtOkButDisabledKeyFile := "testdata/google.key.json"
		got, err := NewClient(context.TODO(), fmtOkButDisabledKeyFile)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}
