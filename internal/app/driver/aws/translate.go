package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

// Translate dispatches to translate.Translate methods.
type Translate struct {
	svc *translate.Translate
}

var _ TranslateWrapper = (*Translate)(nil)

// NewTranslate creates new struct.
func NewTranslate() (*Translate, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("aws session error: %w", err)
	}

	return &Translate{
		svc: translate.New(sess),
	}, nil
}

// TextWithContext dispatches translate.Translate method simply.
func (t *Translate) TextWithContext(ctx aws.Context, input *translate.TextInput, opts ...request.Option) (*translate.TextOutput, error) {
	return t.svc.TextWithContext(ctx, input)
}
