package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

// TranslationAPI serves AWS Translate API handlers.
type TranslationAPI struct {
}

// NewTranslationAPI creates new struct.
func NewTranslationAPI() (*TranslationAPI, error) {
	for _, v := range []string{"AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"} {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf(`environment error: %s`, v)
		}
	}
	a := &TranslationAPI{}
	return a, nil
}

// LanguageType is type of translation language.
type LanguageType string

const (
	// JP language type
	JP LanguageType = "JP"
	// EN language type
	EN LanguageType = "EN"
)

// Returns language code of AWS Translate API.
// Ref:
// https://docs.aws.amazon.com/translate/latest/dg/what-is.html#what-is-languages
func (l LanguageType) languageCode() (string, error) {
	code := ""
	switch l {
	case JP:
		code = "ja"
	case EN:
		code = "en"
	}

	if code == "" {
		return "", fmt.Errorf("unkown type: %v", l)
	}
	return code, nil
}

// TranslateInput is a parameter of Translate method.
type TranslateInput struct {
	Text       string
	SourceLang LanguageType
	TargetLang LanguageType
}

func (t *TranslateInput) textInput() (*translate.TextInput, error) {
	sourceLang, err := t.SourceLang.languageCode()
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	targetLang, err := t.TargetLang.languageCode()
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	input := &translate.TextInput{
		Text:               &t.Text,
		SourceLanguageCode: &sourceLang,
		TargetLanguageCode: &targetLang,
	}
	return input, nil
}

// TranslateOutput is a response of Translate method.
type TranslateOutput struct {
	Text           string
	Lang           LanguageType
	TranslatedText string
}

// Translate translates text with AWS Translate service.
func (a TranslationAPI) Translate(ctx context.Context, input TranslateInput) (*TranslateOutput, error) {
	svcInput, err := input.textInput()
	if err != nil {
		return nil, err
	}

	// Ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
	appSession, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("aws session error: %w", err)
	}

	svc := translate.New(appSession)

	// Call AWS Translate API.
	svcOutput, err := svc.Text(svcInput)
	if err != nil {
		return nil, fmt.Errorf("translate error: %w", err)
	}

	return &TranslateOutput{
		Text:           input.Text,
		Lang:           input.TargetLang,
		TranslatedText: *svcOutput.TranslatedText,
	}, nil
}
