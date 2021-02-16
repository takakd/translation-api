// Package log provides logging feature.
package log

import (
	"context"
	"net/http"
)

// ContextKey represents this context value's key.
type ContextKey struct{}

// ContextValue is set in Context.
type ContextValue struct {
	// Current request ID.
	RequestID string
	// Request date
	Date   string
	Header http.Header
	Host   string
	Path   string
	Method string
}

// WithLogContextValue creates context with log properties.
func WithLogContextValue(ctx context.Context, requestID string, req *http.Request, date string) context.Context {
	ctxValue := &ContextValue{
		RequestID: requestID,
		Date:      date,
	}
	if req != nil {
		ctxValue.Header = req.Header
		ctxValue.Path = req.URL.Path
		ctxValue.Host = req.Host
		ctxValue.Method = req.Method
	}
	return context.WithValue(ctx, ContextKey{}, ctxValue)
}

// GetContextValue returns logging context values.
func GetContextValue(ctx context.Context) *ContextValue {
	if cv, ok := ctx.Value(ContextKey{}).(*ContextValue); ok {
		return cv
	}
	return nil
}
