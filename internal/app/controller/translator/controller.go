package translator

import (
	"api/internal/app/grpc/translator"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"time"

	"api/internal/app/driver/aws"

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

//func googleTranslationAPIRequestLang(langType translator.LangType) google.LanguageType {
//	l := google.EN
//	if langType == translator.LangType_JP {
//		l = google.JP
//	}
//	return l
//}

func awsTranslationAPIRequestLang(langType translator.LangType) aws.LanguageType {
	l := aws.EN
	if langType == translator.LangType_JP {
		l = aws.JP
	}
	return l
}

// Translate processes a method of Translator gRPC service.
func (Controller) Translate(ctx context.Context, r *translator.TranslateRequest) (*translator.TranslateResponse, error) {
	appCtx := log.WithLogContextValue(ctx, uuid.New().String())

	// Access log
	now := time.Now()
	log.Info(appCtx, log.Value{
		"request": r,
		"date":    now.Format(time.RFC3339),
	})

	awsReqBody := aws.TranslateInput{
		Text:       r.GetText(),
		SourceLang: awsTranslationAPIRequestLang(r.GetSrcLang()),
		TargetLang: awsTranslationAPIRequestLang(r.GetTargetLang()),
	}

	awsAPI, err := aws.NewTranslationAPI()
	if err != nil {
		return nil, fmt.Errorf("translation api error: %w", err)
	}

	apiOutput, err := awsAPI.Translate(ctx, awsReqBody)
	if err != nil {
		return nil, fmt.Errorf("translate error: %w", err)
	}

	// Create a response
	resp := &translator.TranslateResponse{
		Text:           r.GetText(),
		SrcLang:        r.GetSrcLang(),
		TargetLang:     r.GetTargetLang(),
		TranslatedText: apiOutput.TranslatedText,
	}

	//// Google
	//apiReqBody := google.TranslationAPIRequestBody{
	//	Text:       r.GetText(),
	//	SourceLang: googleTranslationAPIRequestLang(r.GetSrcLang()),
	//	TargetLang: googleTranslationAPIRequestLang(r.GetTargetLang()),
	//}
	//
	//// Call Google translation API
	//googleAPI, err := google.NewTranslationAPI()
	//if err != nil {
	//	return nil, fmt.Errorf("translation api error: %w", err)
	//}
	//
	//apiResp, err := googleAPI.Translate(appCtx, apiReqBody)
	//if err != nil {
	//	return nil, fmt.Errorf("translate error: %w", err)
	//}
	//
	//// Create a response
	//resp := &translator.TranslateResponse{
	//	Text:           r.GetText(),
	//	SrcLang:        r.GetSrcLang(),
	//	TargetLang:     r.GetTargetLang(),
	//	TranslatedText: apiResp.TranslatedText,
	//}
	return resp, nil
}
