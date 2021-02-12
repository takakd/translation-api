package google

import (
	"context"
	"errors"
	"fmt"

	"api/internal/pkg/util"

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
func NewClient(ctx context.Context, apiKeyFilePath string) (*Client, error) {

	if apiKeyFilePath == "" {
		return nil, errors.New("key file path empty error")
	} else if !util.FileExists(apiKeyFilePath) {
		return nil, fmt.Errorf("key file not exists: path=%s", apiKeyFilePath)
	}

	// Ref: https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile(apiKeyFilePath))
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
