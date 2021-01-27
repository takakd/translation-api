package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTranslate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, err := NewTranslate()
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}
