package api

import (
	"api/internal/app/util/appcontext"
	"fmt"

	"cloud.google.com/go/translate/apiv3"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
	"os"
)

// GoogleTranslationAPI serves Google Translate API handlers used by the app.
type GoogleTranslationAPI struct {
	projectID string
}

// NewGoogleTranslationAPI creates new struct.
func NewGoogleTranslationAPI() *GoogleTranslationAPI {
	a := &GoogleTranslationAPI{}
	a.projectID = os.Getenv("GOOGLE_PROJECT_ID")
	return a
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
		// TODO: check, en-US?
		code = "en"
	}

	if code == "" {
		return "", fmt.Errorf("unkown type: %v", l)
	}
	return code, nil
}

// GoogleTranslationAPIRequestBody is a parameter of Translate method.
type GoogleTranslationAPIRequestBody struct {
	Text       string
	SourceLang LanguageType
	TargetLang LanguageType
}

// Returns TranslateTextRequest with Google API projectID.
func (r *GoogleTranslationAPIRequestBody) translateRequest(projectID string) (*translatepb.TranslateTextRequest, error) {
	sourceLang, err := r.SourceLang.languageCode()
	if err != nil {
		return nil, fmt.Errorf("language code error: %w", err)
	}

	targetLang, err := r.SourceLang.languageCode()
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

// GoogleTranslationAPIResponseBody is a response of Translate method.
type GoogleTranslationAPIResponseBody struct {
	Text           string
	Lang           LanguageType
	TranslatedText string
}

// Translate returns translated text with Google Translate API.
func (a GoogleTranslationAPI) Translate(ctx appcontext.Context, body GoogleTranslationAPIRequestBody) (GoogleTranslationAPIResponseBody, error) {
	resp := GoogleTranslationAPIResponseBody{
		Text:           body.Text,
		Lang:           body.TargetLang,
		TranslatedText: "",
	}

	apiReq, err := body.translateRequest(a.projectID)
	if err != nil {
		return resp, fmt.Errorf("request creation error: %w", err)
	}

	// Ref: https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
	client, err := translate.NewTranslationClient(ctx.Context())
	if err != nil {
		return resp, fmt.Errorf("api initialize error: %w", err)
	}
	defer client.Close()

	// Call Google Translate API
	apiResp, err := client.TranslateText(ctx.Context(), apiReq)
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
