package google

import (
	"fmt"

	translatorapp "api/internal/app/controller/translator"
	"context"

	"api/internal/app/util/config"
	"api/internal/app/util/di"

	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// TranslationAPI serves Google Translate API handlers.
type TranslationAPI struct {
	projectID      string
	apiKeyFilePath string
}

var _ translatorapp.TextTranslator = (*TranslationAPI)(nil)

// NewTranslationAPI creates new struct.
func NewTranslationAPI() (*TranslationAPI, error) {
	var err error

	a := &TranslationAPI{}

	a.projectID, err = config.Get("GOOGLE_PROJECT_ID")
	if a.projectID == "" || err != nil {
		return nil, fmt.Errorf("config error name=GOOGLE_PROJECT_ID: %w", err)
	}

	a.apiKeyFilePath, err = config.Get("GOOGLE_KEY_FILE_PATH")
	if a.apiKeyFilePath == "" || err != nil {
		return nil, fmt.Errorf("config error name=GOOGLE_KEY_FILE_PATH: %w", err)
	}

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

	if text == "" {
		return nil, fmt.Errorf("text empty error")
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
func (a *TranslationAPI) Translate(ctx context.Context, text string, srcLang translatorapp.LanguageType, targetLang translatorapp.LanguageType) (*translatorapp.Result, error) {
	apiReq, err := translateRequest(a.projectID, text, srcLang, targetLang)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	clientInf, err := di.Get("driver.google.Client", ctx, a.apiKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("api initialize error: %w", err)
	}

	client := clientInf.(ClientWrapper)
	defer client.Close()

	// Call Google Translate API
	apiResp, err := client.TranslateText(ctx, apiReq)
	if err != nil {
		return nil, fmt.Errorf("api request error: %w", err)
	}

	result := &translatorapp.Result{
		Lang:    targetLang,
		Service: translatorapp.Google,
	}
	for _, translation := range apiResp.GetTranslations() {
		result.Text = translation.TranslatedText
		// One text Should be returned.
		break
	}

	return result, nil
}
