package log

import (
	log2 "api/internal/app/util/log"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"context"

	"net/http"

	"encoding/json"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStdoutLogger(t *testing.T) {
	l := NewStdoutLogger()
	assert.Equal(t, log2.LevelInfo, l.level)
	assert.NotNil(t, l.logger)
}

func TestStdoutLogger_SetLevel(t *testing.T) {
	l := NewStdoutLogger()
	l.SetLevel(log2.LevelInfo)
	assert.Equal(t, log2.LevelInfo, l.level)
}

func captureOutputLog(ctx context.Context, l *StdoutLogger, level log2.Level, v log2.Value) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l.outputLog(ctx, level, v)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestOutputLog(t *testing.T) {
	t.Run("un match level", func(t *testing.T) {
		tests := []struct {
			name     string
			level    log2.Level
			logLevel log2.Level
			output   bool
		}{
			{name: "debug with debug", level: log2.LevelDebug, logLevel: log2.LevelDebug, output: true},
			{name: "info with debug", level: log2.LevelDebug, logLevel: log2.LevelInfo, output: true},
			{name: "error with debug", level: log2.LevelDebug, logLevel: log2.LevelError, output: true},

			{name: "debug with info", level: log2.LevelInfo, logLevel: log2.LevelDebug, output: false},
			{name: "info with info", level: log2.LevelInfo, logLevel: log2.LevelInfo, output: true},
			{name: "error with info", level: log2.LevelInfo, logLevel: log2.LevelError, output: true},

			{name: "debug with error", level: log2.LevelError, logLevel: log2.LevelDebug, output: false},
			{name: "info with error", level: log2.LevelError, logLevel: log2.LevelInfo, output: false},
			{name: "error with error", level: log2.LevelError, logLevel: log2.LevelError, output: true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				l := NewStdoutLogger()
				l.SetLevel(tt.level)
				got := captureOutputLog(context.TODO(), l, tt.logLevel, log2.Value{"msg": "a"})
				if tt.output {
					assert.NotEmpty(t, got)
				} else {
					assert.Empty(t, got)
				}
			})
		}
	})

	t.Run("empty value", func(t *testing.T) {
		l := NewStdoutLogger()
		l.SetLevel(log2.LevelDebug)
		got := captureOutputLog(context.TODO(), l, log2.LevelDebug, log2.Value{})
		assert.Equal(t, "", got)
	})

	t.Run("with request context", func(t *testing.T) {
		l := NewStdoutLogger()
		l.SetLevel(log2.LevelDebug)

		tReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		tReq.Header.Add("Key", "value1")
		tReqID := "req123"
		tDate := "2021-02-16T20:56:21Z"

		ctx := log2.WithLogContextValue(context.Background(), tReqID, tReq, tDate)
		got := captureOutputLog(ctx, l, log2.LevelDebug, log2.Value{"msg": "a"})

		assert.Contains(t, got, fmt.Sprintf(`"rid":"%s"`, tReqID))
		assert.Contains(t, got, fmt.Sprintf(`"date":"%s"`, tDate))
		assert.Contains(t, got, fmt.Sprintf(`"path":"%s"`, tReq.URL.Path))
		assert.Contains(t, got, fmt.Sprintf(`"method":"%s"`, tReq.Method))

		buf, err := json.Marshal(tReq.Header)
		require.NoError(t, err)
		assert.Contains(t, got, string(buf))
	})

	t.Run("levels", func(t *testing.T) {
		tests := []struct {
			name  string
			level log2.Level
			label string
		}{
			{name: "error", level: log2.LevelError, label: "ERROR"},
			{name: "info", level: log2.LevelInfo, label: "INFO"},
			{name: "debug", level: log2.LevelDebug, label: "DEBUG"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				l := NewStdoutLogger()
				l.SetLevel(tt.level)
				got := captureOutputLog(context.TODO(), l, tt.level, log2.Value{"msg": "a"})
				assert.Contains(t, got, `"msg":"a"`)
				assert.NotContains(t, got, `"rid":""`)
				assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, tt.label))
			})
		}
	})

	t.Run("marshal error", func(t *testing.T) {
		l := NewStdoutLogger()
		l.SetLevel(log2.LevelDebug)
		ch := make(chan int)
		got := captureOutputLog(context.TODO(), l, log2.LevelDebug, log2.Value{"msg": ch})
		assert.Contains(t, got, "marshal error: json: unsupported type: chan int")
	})
}

func TestDebug(t *testing.T) {
	l := NewStdoutLogger()
	l.SetLevel(log2.LevelDebug)
	got := captureOutputLog(context.TODO(), l, log2.LevelDebug, log2.Value{"msg": "a"})
	assert.Contains(t, got, `"msg":"a"`)
	assert.NotContains(t, got, `"rid":""`)
	assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, "DEBUG"))
}

func TestInfo(t *testing.T) {
	l := NewStdoutLogger()
	l.SetLevel(log2.LevelInfo)
	got := captureOutputLog(context.TODO(), l, log2.LevelInfo, log2.Value{"msg": "a"})
	assert.Contains(t, got, `"msg":"a"`)
	assert.NotContains(t, got, `"rid":""`)
	assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, "INFO"))
}

func TestError(t *testing.T) {
	l := NewStdoutLogger()
	l.SetLevel(log2.LevelError)
	got := captureOutputLog(context.TODO(), l, log2.LevelError, log2.Value{"msg": "a"})
	assert.Contains(t, got, `"msg":"a"`)
	assert.NotContains(t, got, `"rid":""`)
	assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, "ERROR"))
}
