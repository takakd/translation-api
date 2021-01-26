package google

import (
	"cloud.google.com/go/translate/apiv3"
	"github.com/googleapis/gax-go/v2"
	"context"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// Client dispatches to TranslationClient methods.
type Client struct {
	client *translate.TranslationClient
}

// TranslateText dispatches TranslationClient method simply.
func (c *Client) TranslateText(ctx context.Context, req *translatepb.TranslateTextRequest, opts ...gax.CallOption) (*translatepb.TranslateTextResponse, error) {
	return c.client.TranslateText(ctx, req, opts...)
}

// Close dispatches TranslationClient method simply.
func (c *Client) Close() error {
	return c.client.Close()
}
