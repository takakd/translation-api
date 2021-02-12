package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLogContextValue(t *testing.T) {
	tt := "req123"
	ctx := WithLogContextValue(context.Background(), tt)
	assert.Equal(t, tt, GetRequestID(ctx))
}
