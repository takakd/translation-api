package google

import (
	"fmt"

	"context"
	"os"

	translatorapp "api/internal/app/controller/translator"

	translate "cloud.google.com/go/translate/apiv3"
	"google.golang.org/api/option"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// TranslationAPI serves Google Translate API handlers.
type TranslationAPI struct {
}

// NewTranslationAPI creates new struct.
func NewTranslationAPI() (*TranslationAPI, error) {
	for _, v := range []string{"GOOGLE_PROJECT_ID", "GOOGLE_API_KEY"} {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf(`environment error: %s`, v)
		}
	}
	a := &TranslationAPI{}
	return a, nil
}

// Returns language code of Google Translate API.
// Code format is BCP-47
// Ref:
// https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
// https://tools.ietf.org/html/bcp47
func languageCode(lang translatorapp.LanguageType) (string, error) {
	code := ""
	switch lang {
	case translatorapp.JP:
		code = "ja"
	case translatorapp.EN:
		code = "en"
	}

	if code == "" {
		return "", fmt.Errorf("unkown type: %v", lang)
	}
	return code, nil
}

// Returns TranslateTextRequest with Google API projectID.
func translateRequest(projectID string, text string, srcLang, targetLang translatorapp.LanguageType) (*translatepb.TranslateTextRequest, error) {
	sourceCode, err := languageCode(srcLang)
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	targetCode, err := languageCode(targetLang)
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		SourceLanguageCode: sourceCode,
		TargetLanguageCode: targetCode,
		MimeType:           "text/plain",
		Contents:           []string{text},
	}

	return req, nil
}

// Translate returns translated text with Google Translate API.
func (a TranslationAPI) Translate(ctx context.Context, text string, srcLang translatorapp.LanguageType, targetLang translatorapp.LanguageType) (*translatorapp.Result, error) {

	apiReq, err := translateRequest(os.Getenv("GOOGLE_PROJECT_ID"), text, srcLang, targetLang)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	// Ref: https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_API_KEY"))))
	if err != nil {
		return nil, fmt.Errorf("api initialize error: %w", err)
	}
	defer client.Close()

	// Call Google Translate API
	apiResp, err := client.TranslateText(ctx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("api request error: %w", err)
	}

	result := &translatorapp.Result{
		Lang:        targetLang,
		ServiceName: translatorapp.Google,
	}
	for _, translation := range apiResp.GetTranslations() {
		result.Text = translation.TranslatedText
		// One text Should be returned.
		break
	}

	return result, nil
}
