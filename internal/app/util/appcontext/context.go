package appcontext

import (
	"context"
)

// Context defines the context interface used in the app.
type Context struct {
	ctx       context.Context
	requestID string
}

// NewContext creates new struct.
func NewContext(ctx context.Context, requestID string) (Context, error) {
	return Context{
		ctx:       ctx,
		requestID: requestID,
	}, nil
}

// TODO returns empty context, which is only used in unit test.
func TODO() *Context {
	return &Context{}
}

// RequestID returns current request ID.
func (c *Context) RequestID() string {
	return c.requestID
}
