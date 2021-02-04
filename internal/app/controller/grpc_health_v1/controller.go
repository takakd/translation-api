package grpc_health_v1

import (
	"api/internal/app/util/log"
	"context"
	"time"

	"api/internal/app/grpc/health/grpc_health_v1"

	"github.com/google/uuid"
)

// Controller handles grpc_health_v1 API.
type Controller struct {
	grpc_health_v1.UnimplementedHealthServer
}

var _ grpc_health_v1.HealthServer = (*Controller)(nil)

// NewController creates new struct.
func NewController() (*Controller, error) {
	return &Controller{}, nil
}

// Check processes a method of grpc_health_v1 gRPC service.
func (c *Controller) Check(ctx context.Context, r *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	appCtx := log.WithLogContextValue(ctx, uuid.New().String())

	// Access log
	now := time.Now()
	log.Info(appCtx, log.Value{
		"request": r,
		"date":    now.Format(time.RFC3339),
	})

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}
