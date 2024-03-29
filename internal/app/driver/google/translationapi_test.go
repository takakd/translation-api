package google

import (
	"fmt"

	translatorapp "api/internal/app/controller/translator"
	"context"

	"api/internal/app/util/config"
	"api/internal/app/util/di"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

func TestNewTranslationAPI(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		apiKeyPath string
	}{
		{name: "project id empty", projectID: "", apiKeyPath: "value"},
		{name: "api key path empty", projectID: "value", apiKeyPath: ""},
		{name: "error", projectID: "value", apiKeyPath: "value"},
		{name: "ok", projectID: "value", apiKeyPath: "testdata/google.key.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mc := config.NewMockConfig(ctrl)
			mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return(tt.projectID)
			mc.EXPECT().Get("GOOGLE_KEY_FILE_PATH").Return(tt.apiKeyPath)
			config.SetConfig(mc)

			_, err := NewTranslationAPI()
			if tt.projectID != "" && tt.apiKeyPath != "" {
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
		lang translatorapp.LanguageType
		want string
	}{
		{name: "jp", lang: translatorapp.JP, want: "ja"},
		{name: "en", lang: translatorapp.EN, want: "en"},
		{name: "error", lang: translatorapp.LanguageType(""), want: ""},
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

func TestTranslateRequest(t *testing.T) {
	tests := []struct {
		name       string
		srcLang    translatorapp.LanguageType
		targetLang translatorapp.LanguageType
		text       string
	}{
		{name: "src error", srcLang: translatorapp.LanguageType(""), targetLang: translatorapp.JP, text: "text"},
		{name: "target error", srcLang: translatorapp.JP, targetLang: translatorapp.LanguageType(""), text: "text"},
		{name: "text error", srcLang: translatorapp.JP, targetLang: translatorapp.EN, text: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := translateRequest("", "", tt.srcLang, tt.targetLang)
			assert.Nil(t, req)
			assert.Error(t, err)
		})
	}

	t.Run("ok", func(t *testing.T) {
		projectID := "test"
		text := "dummy"
		srcLang := translatorapp.JP
		targetLang := translatorapp.EN
		want := &translatepb.TranslateTextRequest{
			Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
			SourceLanguageCode: "ja",
			TargetLanguageCode: "en",
			MimeType:           "text/plain",
			Contents:           []string{text},
		}

		got, err := translateRequest(projectID, text, srcLang, targetLang)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestTranslationAPI_Translate(t *testing.T) {
	t.Run("request value error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return("id")
		mc.EXPECT().Get("GOOGLE_KEY_FILE_PATH").Return("testdata/google.key.json")
		config.SetConfig(mc)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		resp, err := s.Translate(context.TODO(), "", translatorapp.LanguageType(""), translatorapp.JP)
		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("client error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return("id")
		mc.EXPECT().Get("GOOGLE_KEY_FILE_PATH").Return("testdata/google.key.json")
		config.SetConfig(mc)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("driver.google.Client", gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		resp, err := s.Translate(context.TODO(), "test", translatorapp.EN, translatorapp.JP)
		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("request error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		text := "text"
		srcLang := translatorapp.JP
		targetLang := translatorapp.EN
		projectID := "id"

		req := &translatepb.TranslateTextRequest{
			Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
			SourceLanguageCode: "ja",
			TargetLanguageCode: "en",
			MimeType:           "text/plain",
			Contents:           []string{text},
		}

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return(projectID)
		mc.EXPECT().Get("GOOGLE_KEY_FILE_PATH").Return("testdata/google.key.json")
		config.SetConfig(mc)

		mt := NewMockClientWrapper(ctrl)
		mt.EXPECT().TranslateText(ctx, req).Return(nil, errors.New("error"))
		mt.EXPECT().Close()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("driver.google.Client", gomock.Any()).Return(mt, nil)

		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		resp, err := s.Translate(ctx, text, srcLang, targetLang)
		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		text := "text"
		srcLang := translatorapp.JP
		targetLang := translatorapp.EN
		projectID := "id"

		req := &translatepb.TranslateTextRequest{
			Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
			SourceLanguageCode: "ja",
			TargetLanguageCode: "en",
			MimeType:           "text/plain",
			Contents:           []string{text},
		}

		clientResp := &translatepb.TranslateTextResponse{
			Translations: []*translatepb.Translation{
				{
					TranslatedText: "translated",
				},
			},
		}
		wantResp := &translatorapp.Result{
			Lang:    targetLang,
			Service: translatorapp.Google,
			Text:    clientResp.Translations[0].TranslatedText,
		}

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return(projectID)
		mc.EXPECT().Get("GOOGLE_KEY_FILE_PATH").Return("testdata/google.key.json")
		config.SetConfig(mc)

		mt := NewMockClientWrapper(ctrl)
		mt.EXPECT().TranslateText(ctx, req).Return(clientResp, nil)
		mt.EXPECT().Close()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("driver.google.Client", gomock.Any(), gomock.Any()).Return(mt, nil)

		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		got, err := s.Translate(ctx, text, srcLang, targetLang)
		assert.NoError(t, err)
		assert.Equal(t, wantResp, got)
	})
}
