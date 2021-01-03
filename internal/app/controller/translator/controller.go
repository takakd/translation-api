package translator

import (
	"api/internal/app/driver/api"
	"api/internal/app/grpc/translator"
	"api/internal/app/util/appcontext"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Controller handles translate API.
type Controller struct {
	translator.UnimplementedTranslatorServer
}

// NewController creates new struct.
func NewController() *Controller {
	return &Controller{}
}

func googleTranslationAPIRequestLang(langType translator.LangType) api.LanguageType {
	l := api.EN
	if langType == translator.LangType_JP {
		l = api.JP
	}
	return l
}

// Translate processes a method of Translator gRPC service.
func (Controller) Translate(ctx context.Context, r *translator.TranslateRequest) (*translator.TranslateResponse, error) {
	appCtx := appcontext.NewContext(ctx, uuid.New().String())

	// Access log
	now := time.Now()
	log.Info(appCtx, map[string]interface{}{
		"request": r,
		"date":    now.Format(time.RFC3339),
	})

	// Call Google translation API
	googleAPI := api.NewGoogleTranslationAPI()
	apiReqBody := api.GoogleTranslationAPIRequestBody{
		Text: r.Text,
		Lang: googleTranslationAPIRequestLang(r.Lang),
	}
	apiResp, err := googleAPI.Translate(appCtx, apiReqBody)
	if err != nil {
		return nil, fmt.Errorf("translate error: %w", err)
	}

	// Create a response
	resp := &translator.TranslateResponse{
		Text:           r.GetText(),
		Lang:           r.GetLang(),
		TranslatedText: apiResp.TranslatedText,
	}
	return resp, nil
}
