package translator

import (
	"api/internal/app/grpc/translator"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"time"

	"api/internal/app/util/di"

	"sync"
)

// Controller handles translate API.
type Controller struct {
	translator.UnimplementedTranslatorServer

	translatorList []TextTranslator
}

var _ translator.TranslatorServer = (*Controller)(nil)

// NewController creates new struct.
func NewController() (*Controller, error) {
	var (
		tmp  interface{}
		name string
		err  error
	)

	name = "driver.aws.TranslationAPI"
	tmp, err = di.Get(name)
	if err != nil {
		return nil, fmt.Errorf("nil error name=%s: %w", name, err)
	}
	awsTextTranslator, ok := tmp.(TextTranslator)
	if !ok {
		return nil, fmt.Errorf("type error name=%s: %w", name, err)
	}

	name = "driver.google.TranslationAPI"
	tmp, err = di.Get(name)
	if err != nil {
		return nil, fmt.Errorf("nil error name=%s: %w", name, err)
	}
	googleTextTranslator, ok := tmp.(TextTranslator)
	if !ok {
		return nil, fmt.Errorf("type error name=%s: %w", name, err)
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

	// Wait for API response.
	var wg sync.WaitGroup

	for _, t := range c.translatorList {
		// Increment counter
		wg.Add(1)

		// Copy as Value
		t := t

		// Call each API
		go func(ctx context.Context, ch chan<- *TranslateParallelResult, r *translator.TranslateRequest, t TextTranslator) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			result := &TranslateParallelResult{}

			// Call API
			if apiResult, err := t.Translate(ctx, r.GetText(), toTranslatorLang(r.GetSrcLang()), toTranslatorLang(r.GetTargetLang())); apiResult == nil || err != nil {
				if err != nil {
					result.err = err
				} else {
					result.err = fmt.Errorf("result is empty %v", t)
				}
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
		}(ctx, ch, r, t)
	}

	// Wait go routines and close channel.
	go func(ch chan<- *TranslateParallelResult) {
		wg.Wait()
		close(ch)
	}(ch)
}

// Translate processes a method of Translator gRPC service.
func (c *Controller) Translate(ctx context.Context, r *translator.TranslateRequest) (*translator.TranslateResponse, error) {
	// NOTE: grpcserver.Server already sets a request ID in ServeHTTP method.
	appCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// Access log
	now := time.Now()
	log.Info(appCtx, log.Value{
		"request": r,
		"date":    now.Format(time.RFC3339),
		"tag":     "Translate",
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

	// Wait for API response.
	var err error
	for c := range ch {
		if c.err != nil {
			err = fmt.Errorf("translation error: %s, %w", c.serviceType, c.err)
			continue
		}
		resp.TranslatedTextList[string(c.serviceType)] = c.translated
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}
