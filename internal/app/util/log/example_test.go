package log_test

import (
	log2 "api/internal/app/driver/log"
	"api/internal/app/util/log"
	"context"
	"net/http"
	"net/http/httptest"
)

func ExampleDebug() {
	logger := log2.NewStdoutLogger()
	log.SetLogger(logger)

	log.SetLevel(log.LevelDebug)

	ctx := context.Background()
	log.Info(ctx, log.Value{"message": "test"})

	reqID := "107e9a86-f6e3-43d4-9906-29aefc269728"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add("Key", "value1")
	date := "2021-02-17T06:56:39+09:00"

	ctx = log.WithLogContextValue(ctx, reqID, req, date)
	log.Debug(ctx, log.Value{"message": "test"})

	// Output:
	// {"level":"INFO","message":"test"}
	// {"date":"2021-02-17T06:56:39+09:00","header":{"Key":["value1"]},"host":"example.com","level":"DEBUG","message":"test","method":"GET","path":"/test","rid":"107e9a86-f6e3-43d4-9906-29aefc269728"}
}
