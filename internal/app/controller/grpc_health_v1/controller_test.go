package grpc_health_v1

import (
	"context"
	"testing"

	"api/internal/app/grpc/health/grpc_health_v1"
	"api/internal/app/util/log"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewController(t *testing.T) {
	got, err := NewController()
	assert.NotNil(t, got)
	assert.NoError(t, err)
}

func TestController_Check(t *testing.T) {
	tests := []struct {
		name string
		req  *grpc_health_v1.HealthCheckRequest
	}{
		{
			name: "nil",
		},
		{
			name: "empty",
			req:  &grpc_health_v1.HealthCheckRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ml := log.NewMockLogger(ctrl)
			ml.EXPECT().Info(gomock.Any(), gomock.Any())
			log.SetLogger(ml)

			c, err := NewController()
			require.NoError(t, err)

			got, err := c.Check(context.TODO(), tt.req)
			assert.NoError(t, err)
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, got.Status)
		})
	}
}
