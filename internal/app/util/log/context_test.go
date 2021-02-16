package log

import (
	"context"
	"testing"

	"net/http"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestWithLogContextValue(t *testing.T) {
	tReqID := "req123"
	tReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	tReq.Header.Add("Key", "value1")
	tDate := "2021-02-16T20:56:21Z"

	ctx := WithLogContextValue(context.Background(), tReqID, tReq, tDate)

	got := GetContextValue(ctx)
	assert.Equal(t, tReqID, got.RequestID)
	assert.Equal(t, tDate, got.Date)
	assert.Equal(t, tReq.Header, got.Header)
	assert.Equal(t, tReq.Host, got.Host)
	assert.Equal(t, tReq.URL.Path, got.Path)
	assert.Equal(t, tReq.Method, got.Method)
}
