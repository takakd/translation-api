// Package log provides logging feature.
package log

import "context"

// ContextKey represents this context value's key.
type ContextKey struct{}

// ContextValue is set in Context.
type ContextValue struct {
	// Current request ID.
	requestID string
}

// WithLogContextValue creates context with log properties.
func WithLogContextValue(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKey{}, &ContextValue{
		requestID: requestID,
	})
}

// getRequestID returns request ID from ctx.
func getRequestID(ctx context.Context) string {
	if cv, ok := ctx.Value(ContextKey{}).(*ContextValue); ok {
		return cv.requestID
	}
	return ""
}
