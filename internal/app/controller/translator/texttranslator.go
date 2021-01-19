package translator

import (
	"api/internal/app/grpc/translator"
	"context"
)

// TextTranslator is the interface that translate text.
type TextTranslator interface {
	Translate(ctx context.Context, text string, srcLang LanguageType, targetLang LanguageType) (*Result, error)
}

// LanguageType is type of translation language.
type LanguageType string

const (
	// JP language type
	JP LanguageType = "JP"
	// EN language type
	EN LanguageType = "EN"
)

func toGrpcLang(lang LanguageType) translator.LangType {
	l := translator.LangType_EN
	if lang == JP {
		l = translator.LangType_JP
	}
	return l
}

func toTranslatorLang(lang translator.LangType) LanguageType {
	l := EN
	if lang == translator.LangType_JP {
		l = JP
	}
	return l
}

// ServiceType is a translation API service name.
type ServiceType string

const (
	// AWS is AWS Translate service
	AWS ServiceType = "aws"
	// Google is Google Translation API service
	Google ServiceType = "google"
)

// Result has translated result.
type Result struct {
	Text    string
	Lang    LanguageType
	Service ServiceType
}
