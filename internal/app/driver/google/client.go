package google

import (
	"context"
	"errors"
	"fmt"

	translate "cloud.google.com/go/translate/apiv3"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// Client dispatches to TranslationClient methods.
type Client struct {
	client *translate.TranslationClient
}

var _ ClientWrapper = (*Client)(nil)

// NewClient creates new struct.
func NewClient(ctx context.Context, apiKey string) (*Client, error) {

	if apiKey == "" {
		return nil, errors.New("apiKey empty error")
	}

	// Ref: https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsJSON([]byte(apiKey)))
	if err != nil {
		err = fmt.Errorf("api initialize error: %w", err)
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
