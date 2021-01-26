package google

import (
	"github.com/googleapis/gax-go/v2"
	"context"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// ClientWrapper wraps Google Translation gRPC interface for Unit test.
type ClientWrapper interface {
	// https://pkg.go.dev/cloud.google.com/go/translate/apiv3#TranslationClient.TranslateText
	TranslateText(ctx context.Context, req *translatepb.TranslateTextRequest, opts ...gax.CallOption) (*translatepb.TranslateTextResponse, error)

	// https://pkg.go.dev/cloud.google.com/go/translate/apiv3#TranslationClient.Close
	Close() error
}
