package google

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "cloud.google.com/go/translate/apiv3"
)

func TestNewClient(t *testing.T) {
    t.Run("nil", func(t *testing.T) {
        got, err := NewClient(nil)
        assert.Nil(t, got)
        assert.Error(t, err)
    })

    t.Run("ok", func(t *testing.T) {
        // Only check whether NewClient can create new struct.
        client := &translate.TranslationClient{}
        got, err := NewClient(client)
        assert.NoError(t, err)
        assert.NotNil(t, got)
    })
}
