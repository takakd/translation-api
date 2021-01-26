package google

import (
	"fmt"

	"context"
	translatorapp "api/internal/app/controller/translator"

	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"api/internal/app/util/config"
	"gobotskelton.org/appgopath/src/github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"api/internal/app/util/di"
)

func TestNewTranslationAPI(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		apiKey    string
		err       error
	}{
		{name: "project id empty", projectID: "", apiKey: "value"},
		{name: "api key empty", projectID: "value", apiKey: ""},
		{name: "error", projectID: "value", apiKey: "value", err: errors.New("error")},
		{name: "ok", projectID: "value", apiKey: "value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mc := config.NewMockConfig(ctrl)

			mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return(tt.projectID, tt.err)

			if tt.projectID != "" && tt.err == nil {
				mc.EXPECT().Get("GOOGLE_API_KEY").Return(tt.apiKey, tt.err)
			}

			_, err := NewTranslationAPI()
			if tt.projectID != "" && tt.apiKey != "" && tt.err == nil {
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
		{name: "target error", srcLang: translatorapp.JP, targetLang: translatorapp.EN, text: ""},
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
			Parent:             fmt.Sprintf("project/%s/locations/global", projectID),
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
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return("id", nil)
		mc.EXPECT().Get("GOOGLE_API_KEY").Return("key", nil)

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
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return("id", nil)
		mc.EXPECT().Get("GOOGLE_API_KEY").Return("key", nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("translate.NewTranslationClient", gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

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

		mc := config.NewMockConfig(ctrl)
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return("id", nil)
		mc.EXPECT().Get("GOOGLE_API_KEY").Return("key", nil)

		mt := NewMockClientWrapper(ctrl)
		mt.EXPECT().TranslateText(ctx, text, srcLang, targetLang).Return(nil, errors.New("error"))

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("translate.NewTranslationClient", gomock.Any(), gomock.Any()).Return(mt, nil)

		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		resp, err := s.Translate(context.TODO(), text, srcLang, targetLang)
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
		mc.EXPECT().Get("GOOGLE_PROJECT_ID").Return("id", nil)
		mc.EXPECT().Get("GOOGLE_API_KEY").Return("key", nil)

		mt := NewMockClientWrapper(ctrl)
		mt.EXPECT().TranslateText(ctx, text, srcLang, targetLang).Return(clientResp, nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("translate.NewTranslationClient", gomock.Any(), gomock.Any()).Return(mt, nil)

		di.SetDi(md)

		s, err := NewTranslationAPI()
		require.NoError(t, err)

		got, err := s.Translate(context.TODO(), text, srcLang, targetLang)
		assert.NoError(t, err)
		assert.Equal(t, wantResp, got)
	})
}
