package google

import (
    "cloud.google.com/go/translate/apiv3"
    "github.com/googleapis/gax-go/v2"
    "context"
    translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
    "errors"
)

// Client dispatches to TranslationClient methods.
type Client struct {
    client *translate.TranslationClient
}

var _ ClientWrapper = (*Client)(nil)

// NewClient creates new struct.
func NewClient(client *translate.TranslationClient) (*Client, error) {
    if client == nil {
        return nil, errors.New("nil error")
    }
    return &Client{
        client: client,
    }, nil
}

// TranslateText dispatches TranslationClient method simply.
func (c *Client) TranslateText(ctx context.Context, req *translatepb.TranslateTextRequest, opts ...gax.CallOption) (*translatepb.TranslateTextResponse, error) {
    return c.client.TranslateText(ctx, req, opts...)
}

// Close dispatches TranslationClient method simply.
func (c *Client) Close() error {
    return c.client.Close()
}
