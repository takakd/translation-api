package translator

import (
	"api/internal/app/grpc/translator"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Controller handles translate API.
type Controller struct {
	translator.UnimplementedTranslatorServer

	translatorList []TextTranslator
}

// NewController creates new struct.
func NewController() (*Controller, error) {
	c := &Controller{}

	var (
		err error
		t   TextTranslator
	)

	//t, err = aws.NewTranslationAPI()
	//if err != nil {
	//	return nil, fmt.Errorf("NewController error: %w", err)
	//}
	//c.translatorList = append(c.translatorList, t)
	//
	//t, err = google.NewTranslationAPI()
	//if err != nil {
	//	return nil, fmt.Errorf("NewController error: %w", err)
	//}
	c.translatorList = append(c.translatorList, t)

	return c, nil
}

// TranslateParallelResult has TranslateParallel result.
type TranslateParallelResult struct {
	translated  *translator.TranslatedText
	serviceType ServiceType
}

// TranslateParallel translate text in parallel.
func (c Controller) TranslateParallel(ctx context.Context, ch chan<- *TranslateParallelResult, r *translator.TranslateRequest, eg errgroup.Group) {
	for _, t := range c.translatorList {
		eg.Go(func() error {
			// Call API
			result, err := t.Translate(ctx, r.GetText(), toTranslatorLang(r.GetSrcLang()), toTranslatorLang(r.GetTargetLang()))
			if err != nil {
				return fmt.Errorf("translate error: %w", err)
			}

			// Send result
			ch <- &TranslateParallelResult{
				translated: &translator.TranslatedText{
					Text: result.Text,
					Lang: toGrpcLang(result.Lang),
				},
				serviceType: result.Service,
			}

			return nil
		})
	}
}

// Translate processes a method of Translator gRPC service.
func (c Controller) Translate(ctx context.Context, r *translator.TranslateRequest) (*translator.TranslateResponse, error) {
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
	var eg errgroup.Group
	ch := make(chan *TranslateParallelResult)
	c.TranslateParallel(ctx, ch, r, eg)

	// Wait
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	// Set translation result.
	for c := range ch {
		resp.TranslatedTextList[string(c.serviceType)] = c.translated
	}

	return resp, nil
}
