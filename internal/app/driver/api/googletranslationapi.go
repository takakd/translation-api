package api

import "context"

// GoogleTranslationAPI serves Google Translate API handlers used by the app.
type GoogleTranslationAPI struct {
}

// NewGoogleTranslationAPI creates new struct.
func NewGoogleTranslationAPI() *GoogleTranslationAPI {
	return &GoogleTranslationAPI{}
}

// LanguageType is type of translation language.
type LanguageType string

const (
	// JP language type
	JP LanguageType = "JP"
	// EN language type
	EN LanguageType = "EN"
)

// GoogleTranslationAPIRequestBody is a parameter of Translate method.
type GoogleTranslationAPIRequestBody struct {
	Text string
	Lang LanguageType
}

// GoogleTranslationAPIResponseBody is a response of Translate method.
type GoogleTranslationAPIResponseBody struct {
	Text           string
	Lang           LanguageType
	TranslatedText string
}

// Translate returns translated text with Google Translate API.
func (a GoogleTranslationAPI) Translate(ctx context.Context, body GoogleTranslationAPIRequestBody) (resp GoogleTranslationAPIResponseBody, err error) {
	// TODO: channel
	resp = GoogleTranslationAPIResponseBody{
		Text:           body.Text,
		Lang:           body.Lang,
		TranslatedText: "hoge",
	}
	return
}
