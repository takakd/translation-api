package aws

import (
	"api/internal/app/controller/translator"
	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranslationAPI(t *testing.T) {
	tests := []struct {
		name      string
		region    string
		keyID     string
		accessKey string
	}{
		{name: "region empty", region: "", keyID: "key", accessKey: "access"},
		{name: "keyID empty", region: "region", keyID: "", accessKey: "access"},
		{name: "accessKey empty", region: "region", keyID: "key", accessKey: ""},
		{name: "ok", region: "region", keyID: "key", accessKey: "access"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mc := config.NewMockConfig(ctrl)
			mc.EXPECT().Get("AWS_REGION").Return(tt.region)
			if tt.region != "" {
				mc.EXPECT().Get("AWS_ACCESS_KEY_ID").Return(tt.keyID)
				if tt.keyID != "" {
					mc.EXPECT().Get("AWS_SECRET_ACCESS_KEY").Return(tt.accessKey)
				}
			}
			config.SetConfig(mc)

			_, err := NewTranslationAPI()
			if tt.name == "ok" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestLanguageCode(t *testing.T) {
	tests := []struct {
		name string
		lang translator.LanguageType
		want string
	}{
		{name: "ja", lang: translator.JP, want: "ja"},
		{name: "en", lang: translator.EN, want: "en"},
		{name: "error", lang: translator.LanguageType(""), want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := languageCode(tt.lang)
			if tt.want != "" {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestTextInput(t *testing.T) {
	tests := []struct {
		name       string
		srcLang    translator.LanguageType
		targetLang translator.LanguageType
		text       string
	}{
		{name: "src error", srcLang: translator.LanguageType(""), targetLang: translator.JP, text: "text"},
		{name: "target error", srcLang: translator.JP, targetLang: translator.LanguageType(""), text: "text"},
		{name: "text error", srcLang: translator.LanguageType(""), targetLang: translator.JP, text: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := textInput(tt.text, tt.srcLang, tt.targetLang)
			assert.Nil(t, got)
			assert.Error(t, err)
		})
	}

	t.Run("ok", func(t *testing.T) {
		srcLang := translator.JP
		targetLang := translator.EN
		srcCode := "ja"
		targetCode := "en"
		text := "text"
		want := &translate.TextInput{
			Text:               &text,
			SourceLanguageCode: &srcCode,
			TargetLanguageCode: &targetCode,
		}
		got, err := textInput(text, srcLang, targetLang)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestTranslationAPI_Translate(t *testing.T) {
	t.Run("input error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("AWS_REGION").Return("region")
		mc.EXPECT().Get("AWS_ACCESS_KEY_ID").Return("key")
		mc.EXPECT().Get("AWS_SECRET_ACCESS_KEY").Return("access")
		config.SetConfig(mc)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		got, err := s.Translate(context.TODO(), "", translator.JP, translator.EN)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("translate service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("AWS_REGION").Return("region")
		mc.EXPECT().Get("AWS_ACCESS_KEY_ID").Return("key")
		mc.EXPECT().Get("AWS_SECRET_ACCESS_KEY").Return("access")
		config.SetConfig(mc)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("driver.aws.Translate").Return(nil, errors.New("error"))
		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		got, err := s.Translate(context.TODO(), "text", translator.JP, translator.EN)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("request error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("AWS_REGION").Return("region")
		mc.EXPECT().Get("AWS_ACCESS_KEY_ID").Return("key")
		mc.EXPECT().Get("AWS_SECRET_ACCESS_KEY").Return("access")
		config.SetConfig(mc)

		ms := NewMockTranslateWrapper(ctrl)
		ms.EXPECT().TextWithContext(ctx, gomock.Any()).Return(nil, errors.New("error"))

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("driver.aws.Translate").Return(ms, nil)
		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		got, err := s.Translate(ctx, "text", translator.JP, translator.EN)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		text := "text"
		srcLang := translator.JP
		srcLangCode := "ja"
		targetLang := translator.EN
		targetLangCode := "en"

		input := &translate.TextInput{
			Text:               &text,
			SourceLanguageCode: &srcLangCode,
			TargetLanguageCode: &targetLangCode,
		}

		ret := &translate.TextOutput{
			TranslatedText: &text,
		}

		want := &translator.Result{
			Text:    *ret.TranslatedText,
			Lang:    targetLang,
			Service: translator.AWS,
		}

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("AWS_REGION").Return("region")
		mc.EXPECT().Get("AWS_ACCESS_KEY_ID").Return("key")
		mc.EXPECT().Get("AWS_SECRET_ACCESS_KEY").Return("access")
		config.SetConfig(mc)

		ms := NewMockTranslateWrapper(ctrl)
		ms.EXPECT().TextWithContext(ctx, input).Return(ret, nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("driver.aws.Translate").Return(ms, nil)
		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		got, err := s.Translate(ctx, text, srcLang, targetLang)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
