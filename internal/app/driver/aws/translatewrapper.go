package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/translate"
)

// TranslateWrapper wraps AWS Translate API interface for Unit test.
//		Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/translate/
type TranslateWrapper interface {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/translate/#Translate.TextWithContext
	TextWithContext(ctx aws.Context, input *translate.TextInput, opts ...request.Option) (*translate.TextOutput, error)
}
