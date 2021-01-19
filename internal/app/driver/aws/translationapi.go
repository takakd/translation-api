package aws

import (
	"context"
	"fmt"
	"os"

	"api/internal/app/controller/translator"

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

// Returns language code of AWS Translate API.
// Ref:
// 		https://docs.aws.amazon.com/translate/latest/dg/what-is.html#what-is-languages
func languageCode(lang translator.LanguageType) (string, error) {
	code := ""
	switch lang {
	case translator.JP:
		code = "ja"
	case translator.EN:
		code = "en"
	}

	if code == "" {
		return "", fmt.Errorf("unkown type: %v", lang)
	}
	return code, nil
}

func textInput(text string, srcLang, targetLang translator.LanguageType) (*translate.TextInput, error) {
	sourceCode, err := languageCode(srcLang)
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	targetCode, err := languageCode(targetLang)
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	input := &translate.TextInput{
		Text:               &text,
		SourceLanguageCode: &sourceCode,
		TargetLanguageCode: &targetCode,
	}
	return input, nil
}

// Translate translates text with AWS Translate service.
func (a TranslationAPI) Translate(ctx context.Context, text string, srcLang translator.LanguageType, targetLang translator.LanguageType) (*translator.Result, error) {

	svcInput, err := textInput(text, srcLang, targetLang)
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

	return &translator.Result{
		Text:        *svcOutput.TranslatedText,
		Lang:        targetLang,
		ServiceName: translator.AWS,
	}, nil
}
