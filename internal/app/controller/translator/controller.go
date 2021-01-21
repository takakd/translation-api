package translator

import (
	"api/internal/app/grpc/translator"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"time"

	"api/internal/app/util/di"

	"github.com/google/uuid"
)

// Controller handles translate API.
type Controller struct {
	translator.UnimplementedTranslatorServer

	translatorList []TextTranslator
}

// NewController creates new struct.
func NewController() (*Controller, error) {
	var (
		tmp interface{}
		err error
	)

	tmp, err = di.Get("translator.awsTextTranslator")
	if err != nil {
		return nil, fmt.Errorf("nil error: awsTextTranslator")
	}
	awsTextTranslator, ok := tmp.(TextTranslator)
	if !ok {
		return nil, fmt.Errorf("type error: awsTextTranslator")
	}

	tmp, err = di.Get("translator.googleTextTranslator")
	if err != nil {
		return nil, fmt.Errorf("nil error: googleTextTranslator")
	}
	googleTextTranslator, ok := tmp.(TextTranslator)
	if !ok {
		return nil, fmt.Errorf("type error: googleTextTranslator")
	}

	c := &Controller{}
	c.translatorList = []TextTranslator{
		awsTextTranslator,
		googleTextTranslator,
	}
	return c, nil
}

// TranslateParallelResult has TranslateParallel result.
type TranslateParallelResult struct {
	translated  *translator.TranslatedText
	serviceType ServiceType
	err         error
}

// TranslateParallel translate text in parallel.
func (c *Controller) TranslateParallel(ctx context.Context, ch chan<- *TranslateParallelResult, r *translator.TranslateRequest) {

	for _, t := range c.translatorList {
		t := t
		go func(ctx context.Context, ch chan<- *TranslateParallelResult, r *translator.TranslateRequest, t TextTranslator) {
			result := &TranslateParallelResult{
			}

			// Call API
			if apiResult, err := t.Translate(ctx, r.GetText(), toTranslatorLang(r.GetSrcLang()), toTranslatorLang(r.GetTargetLang())); err != nil {
				result.err = err
				ch <- result
			} else {
				result.translated = &translator.TranslatedText{
					Text: apiResult.Text,
					Lang: toGrpcLang(apiResult.Lang),
				}
				result.serviceType = apiResult.Service
			}

			// Send result
			ch <- result

			fmt.Println("done")
		}(ctx, ch, r, t)
	}
}

// Translate processes a method of Translator gRPC service.
func (c *Controller) Translate(ctx context.Context, r *translator.TranslateRequest) (*translator.TranslateResponse, error) {
	appCtx := log.WithLogContextValue(ctx, uuid.New().String())
	appCtx, cancel := context.WithTimeout(appCtx, time.Minute)
	defer cancel()

	// Access log
	now := time.Now()
	log.Info(appCtx, log.Value{
		"request": r,
		"date":    now.Format(time.RFC3339),
	})

	// Prepare result
	resp := &translator.TranslateResponse{
		Text:    r.GetText(),
		SrcLang: r.GetSrcLang(),
	}
	resp.TranslatedTextList = make(map[string]*translator.TranslatedText)

	// Call translation in parallel.
	ch := make(chan *TranslateParallelResult)
	c.TranslateParallel(ctx, ch, r)
	//if err != nil {
	//	return nil, fmt.Errorf("translation error: %w", err)
	//}

	fmt.Println("hoge")

	// Set translation result.
	// TOOD: https://hori-ryota.com/blog/golang-channel-pattern/
	// Wait with other go routine.
	for c := range ch {
		fmt.Println(c)
		resp.TranslatedTextList[string(c.serviceType)] = c.translated
	}

	fmt.Print("hogehoge")

	return resp, nil
}
