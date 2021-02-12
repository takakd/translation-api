package translator

import (
	"api/internal/app/grpc/translator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToGrpcLang(t *testing.T) {
	assert.Equal(t, translator.LangType_EN, toGrpcLang(EN))
	assert.Equal(t, translator.LangType_JP, toGrpcLang(JP))
}

func TestToTranslatorLang(t *testing.T) {
	assert.Equal(t, JP, toTranslatorLang(translator.LangType_JP))
	assert.Equal(t, EN, toTranslatorLang(translator.LangType_EN))
}
