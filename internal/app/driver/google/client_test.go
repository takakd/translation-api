package google

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("api key empty", func(t *testing.T) {
		got, err := NewClient(context.TODO(), "")
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		dummyAPIKey :=
			`{
  "type": "service_account",
  "project_id": "example",
  "private_key_id": "0000000000",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDeferB78h4B47S\nHbvqsrd+xPQ5ud6uqMi8fIZ8lFH7CdA==\n-----END PRIVATE KEY-----\n",
  "client_email": "example@example.iam.gserviceaccount.com",
  "client_id": "00000000000",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/example%40example.iam.gserviceaccount.com"
}`
		got, err := NewClient(context.TODO(), dummyAPIKey)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}
