package aws

import (
	"api/internal/app/controller/translator"
	"context"
	"fmt"

	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"errors"

	"github.com/aws/aws-sdk-go/service/translate"
)

// TranslationAPI serves AWS Translate API handlers.
type TranslationAPI struct {
}

var _ translator.TextTranslator = (*TranslationAPI)(nil)

// NewTranslationAPI creates new struct.
func NewTranslationAPI() (*TranslationAPI, error) {
	// Check to exist environment variables because AWS SDK reads credentials through environment variables.
	for _, v := range []string{"AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"} {
		if cv := config.Get(v); cv == "" {
			return nil, fmt.Errorf(`config error name=%s`, v)
		}
	}
	return &TranslationAPI{}, nil
}

// Returns language code of AWS Translate API.
// Ref: https://docs.aws.amazon.com/translate/latest/dg/what-is.html#what-is-languages
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

	if text == "" {
		return nil, errors.New("text empty error")
	}

	input := &translate.TextInput{
		Text:               &text,
		SourceLanguageCode: &sourceCode,
		TargetLanguageCode: &targetCode,
	}
	return input, nil
}

// Translate translates text with AWS Translate service.
func (a *TranslationAPI) Translate(ctx context.Context, text string, srcLang translator.LanguageType, targetLang translator.LanguageType) (*translator.Result, error) {

	svcInput, err := textInput(text, srcLang, targetLang)
	if err != nil {
		return nil, err
	}

	svcInf, err := di.Get("driver.aws.Translate")
	if err != nil {
		return nil, fmt.Errorf("translate initialize error: %w", err)
	}

	svc := svcInf.(TranslateWrapper)

	// Call AWS Translate API.
	svcOutput, err := svc.TextWithContext(ctx, svcInput)
	if err != nil {
		return nil, fmt.Errorf("translate error: %w", err)
	}

	return &translator.Result{
		Text:    *svcOutput.TranslatedText,
		Lang:    targetLang,
		Service: translator.AWS,
	}, nil
}
