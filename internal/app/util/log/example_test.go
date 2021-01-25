package log_test

import (
	"api/internal/app/util/log"
	"context"
)

func ExampleDebug() {
	log.SetLevel(log.LevelDebug)

	ctx := context.Background()
	log.Info(ctx, log.Value{"message": "test"})

	ctx = log.WithLogContextValue(ctx, "req123456")
	log.Debug(ctx, log.Value{"message": "test"})

	// Output:
	// {"level":"INFO","message":"test","rid":""}
	// {"level":"DEBUG","message":"test","rid":"req123456"}
}
