package google

import (
	"fmt"

	"context"
	"os"

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

// LanguageType is type of translation language.
type LanguageType string

const (
	// JP language type
	JP LanguageType = "JP"
	// EN language type
	EN LanguageType = "EN"
)

// Returns language code of Google Translate API.
// Code format is BCP-47
// Ref:
// https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
// https://tools.ietf.org/html/bcp47
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

// TranslationAPIRequestBody is a parameter of Translate method.
type TranslationAPIRequestBody struct {
	Text       string
	SourceLang LanguageType
	TargetLang LanguageType
}

// Returns TranslateTextRequest with Google API projectID.
func (r *TranslationAPIRequestBody) translateRequest(projectID string) (*translatepb.TranslateTextRequest, error) {
	sourceLang, err := r.SourceLang.languageCode()
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	targetLang, err := r.TargetLang.languageCode()
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		SourceLanguageCode: sourceLang,
		TargetLanguageCode: targetLang,
		MimeType:           "text/plain",
		Contents:           []string{r.Text},
	}

	return req, nil
}

// TranslationAPIResponseBody is a response of Translate method.
type TranslationAPIResponseBody struct {
	Text           string
	Lang           LanguageType
	TranslatedText string
}

// Translate returns translated text with Google Translate API.
func (a TranslationAPI) Translate(ctx context.Context, body TranslationAPIRequestBody) (TranslationAPIResponseBody, error) {
	resp := TranslationAPIResponseBody{
		Text:           body.Text,
		Lang:           body.TargetLang,
		TranslatedText: "",
	}

	apiReq, err := body.translateRequest(os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return resp, fmt.Errorf("request creation error: %w", err)
	}

	// Ref: https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_API_KEY"))))
	if err != nil {
		return resp, fmt.Errorf("api initialize error: %w", err)
	}
	defer client.Close()

	// Call Google Translate API
	apiResp, err := client.TranslateText(ctx, apiReq)
	if err != nil {
		return resp, fmt.Errorf("api request error: %w", err)
	}

	for _, translation := range apiResp.GetTranslations() {
		resp.TranslatedText = translation.TranslatedText
		// One text Should be returned.
		break
	}

	return resp, nil
}
