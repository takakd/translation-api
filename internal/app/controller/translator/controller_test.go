package translator

import (
	"api/internal/app/grpc/translator"
	"api/internal/app/util/di"
	"api/internal/app/util/log"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type TestTranslatorTypeError struct{}
type TestTranslator struct{}

func (t *TestTranslator) Translate(ctx context.Context, text string, srcLang LanguageType, targetLang LanguageType) (*Result, error) {
	return nil, nil
}

var _ TextTranslator = (*TestTranslator)(nil)

func TestNewController(t *testing.T) {
	tError := errors.New("error")
	var tInf TextTranslator = &TestTranslator{}
	tests := []struct {
		name             string
		awsErr           error
		awsTranslator    interface{}
		googleErr        error
		googleTranslator interface{}
	}{
		{name: "aws error", awsErr: tError},
		{name: "aws type error", awsTranslator: TestTranslatorTypeError{}},
		{name: "google error", awsTranslator: tInf, googleErr: tError},
		{name: "google type error", awsTranslator: tInf, googleTranslator: TestTranslatorTypeError{}},
		{name: "ok", awsTranslator: tInf, googleTranslator: tInf},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			md := di.NewMockDI(ctrl)

			md.EXPECT().Get("translator.awsTextTranslator").Return(tt.awsTranslator, tt.awsErr)

			if tt.awsTranslator == tInf && tt.awsErr == nil {
				md.EXPECT().Get("translator.googleTextTranslator").Return(tt.googleTranslator, tt.googleErr)
			}

			di.SetDi(md)

			got, err := NewController()
			if tt.name == "ok" {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tInf, got.translatorList[0])
				assert.Equal(t, tInf, got.translatorList[1])
			} else {
				assert.Error(t, err)
				assert.Nil(t, got)
			}
		})
	}
}

func TestController_TranslateParallel(t *testing.T) {
	tests := []struct {
		name       string
		awsRes     *Result
		awsErr     error
		googleRes  *Result
		googleErr  error
		wantAws    *TranslateParallelResult
		wantGoogle *TranslateParallelResult
	}{
		{
			name: "aws error",
			googleRes: &Result{
				Service: Google,
				Text:    "text",
				Lang:    EN,
			},
			wantGoogle: &TranslateParallelResult{
				serviceType: Google,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
		},
		{
			name: "google error",
			awsRes: &Result{
				Service: AWS,
				Text:    "text",
				Lang:    EN,
			},
			wantAws: &TranslateParallelResult{
				serviceType: AWS,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
		},
		{
			name: "both error",
		},
		{
			name: "ok",
			awsRes: &Result{
				Service: AWS,
				Text:    "text",
				Lang:    EN,
			},
			googleRes: &Result{
				Service: Google,
				Text:    "text",
				Lang:    EN,
			},
			wantAws: &TranslateParallelResult{
				serviceType: AWS,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
			wantGoogle: &TranslateParallelResult{
				serviceType: Google,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			text := "text"
			srcLang := translator.LangType_JP
			srcLangType := JP
			targetLang := translator.LangType_EN
			targetLangType := EN

			req := &translator.TranslateRequest{
				Text:       text,
				SrcLang:    srcLang,
				TargetLang: targetLang,
			}

			mat := NewMockTextTranslator(ctrl)
			mat.EXPECT().Translate(ctx, text, srcLangType, targetLangType).Return(tt.awsRes, tt.awsErr)

			mgt := NewMockTextTranslator(ctrl)
			mgt.EXPECT().Translate(ctx, text, srcLangType, targetLangType).Return(tt.googleRes, tt.googleErr)

			c := &Controller{}
			c.translatorList = []TextTranslator{
				mat,
				mgt,
			}

			ch := make(chan *TranslateParallelResult)

			c.TranslateParallel(ctx, ch, req)

			// Wait for API response.
			for c := range ch {
				if c.serviceType == Google {
					if tt.wantGoogle != nil {
						assert.NoError(t, c.err)
						assert.Equal(t, c.translated, tt.wantGoogle.translated)
					} else {
						assert.Error(t, c.err)
					}

				} else if c.serviceType == AWS {
					if tt.wantAws != nil {
						assert.NoError(t, c.err)
						assert.Equal(t, c.translated, tt.wantAws.translated)
					} else {
						assert.Error(t, c.err)
					}
				}
			}
		})
	}
}

func TestController_Translate(t *testing.T) {
	tests := []struct {
		name       string
		awsRes     *Result
		awsErr     error
		googleRes  *Result
		googleErr  error
		wantAws    *TranslateParallelResult
		wantGoogle *TranslateParallelResult
	}{
		{
			name: "error",
			googleRes: &Result{
				Service: Google,
				Text:    "text",
				Lang:    EN,
			},
			wantGoogle: &TranslateParallelResult{
				serviceType: Google,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
		},
		{
			name: "ok",
			awsRes: &Result{
				Service: AWS,
				Text:    "text",
				Lang:    EN,
			},
			googleRes: &Result{
				Service: Google,
				Text:    "text",
				Lang:    EN,
			},
			wantAws: &TranslateParallelResult{
				serviceType: AWS,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
			wantGoogle: &TranslateParallelResult{
				serviceType: Google,
				translated: &translator.TranslatedText{
					Text: "text",
					Lang: translator.LangType_EN,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			text := "text"
			srcLang := translator.LangType_JP
			srcLangType := JP
			targetLang := translator.LangType_EN
			targetLangType := EN

			req := &translator.TranslateRequest{
				Text:       text,
				SrcLang:    srcLang,
				TargetLang: targetLang,
			}

			mat := NewMockTextTranslator(ctrl)
			mat.EXPECT().Translate(ctx, text, srcLangType, targetLangType).Return(tt.awsRes, tt.awsErr)

			mgt := NewMockTextTranslator(ctrl)
			mgt.EXPECT().Translate(ctx, text, srcLangType, targetLangType).Return(tt.googleRes, tt.googleErr)

			ml := log.NewMockLogger(ctrl)
			ml.EXPECT().Info(gomock.Any(), gomock.Any())
			log.SetLogger(ml)

			c := &Controller{}
			c.translatorList = []TextTranslator{
				mat,
				mgt,
			}

			resp, err := c.Translate(ctx, req)

			if tt.name == "ok" {
				assert.NotNil(t, resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAws.translated, resp.TranslatedTextList[string(AWS)])
				assert.Equal(t, tt.wantGoogle.translated, resp.TranslatedTextList[string(Google)])
			} else {
				assert.Nil(t, resp)
				assert.Error(t, err)
			}
		})
	}
}
