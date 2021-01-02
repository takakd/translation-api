package translate

import (
	"net/http"
	"encoding/json"
	"api/internal/app/driver/api"
	"context"
	"api/internal/pkg/util"
)

// Language type
type LanguageType string

const (
	JP LanguageType = "JP"
	EN LanguageType = "EN"
)

// RequestBody is a parameter of controller action.
type RequestBody struct {
	// Source text for translation
	Text string `json:"text"`
	// Language to translate into
	//		e.g. Set EN if text is to be in English.
	Lang LanguageType `json:"lang"`
}

// googleTranslationAPIRequestLang converts Lang to api.LanguageType and returns it.
func (r RequestBody) googleTranslationAPIRequestLang() api.LanguageType {
	var lang api.LanguageType
	switch r.Lang {
	case EN:
		lang = api.EN
	default:
		lang = api.JP
	}
	return lang
}

// ResponseBody is a response of controller action.
type ResponseBody  struct {
	Text string `json:"text"`
	Lang LanguageType `json:"lang"`
	TranslatedText string `json:"translatedtext"`
}

// Controller handles translate API.
type Controller struct {
}

// NewController creates new struct.
func NewController() *Controller {
	return &Controller{}
}

// Handle processes translate API.
func (c *Controller)Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiReqBody := api.GoogleTranslationAPIRequestBody{
		Text: body.Text,
		Lang: body.googleTranslationAPIRequestLang(),
	}
	api := api.NewGoogleTranslationAPI()

	apiResp, err := api.Translate(ctx, apiReqBody);
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ResponseBody{
		Text: body.Text,
		Lang: body.Lang,
		TranslatedText: apiResp.TranslatedText,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSONResponse(w, nil, http.StatusOK, respBody)
	return
}
