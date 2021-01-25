// Package log provides logging.
package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringValue(t *testing.T) {
	got := StringValue("a")
	assert.Equal(t, Value{"msg": "a"}, got)
}

func TestSetLevel(t *testing.T) {
	SetLevel(LevelDebug)
	assert.Equal(t, LevelDebug, logLevel)
}

func captureOutputLog(ctx context.Context, level Level, v Value) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputLog(ctx, level, v)

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
			level    Level
			logLevel Level
			output   bool
		}{
			{name: "debug with debug", level: LevelDebug, logLevel: LevelDebug, output: true},
			{name: "info with debug", level: LevelDebug, logLevel: LevelInfo, output: true},
			{name: "error with debug", level: LevelDebug, logLevel: LevelError, output: true},

			{name: "debug with info", level: LevelInfo, logLevel: LevelDebug, output: false},
			{name: "info with info", level: LevelInfo, logLevel: LevelInfo, output: true},
			{name: "error with info", level: LevelInfo, logLevel: LevelError, output: true},

			{name: "debug with error", level: LevelError, logLevel: LevelDebug, output: false},
			{name: "info with error", level: LevelError, logLevel: LevelInfo, output: false},
			{name: "error with error", level: LevelError, logLevel: LevelError, output: true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				SetLevel(tt.level)
				got := captureOutputLog(context.TODO(), tt.logLevel, Value{"msg": "a"})
				if tt.output {
					assert.NotEmpty(t, got)
				} else {
					assert.Empty(t, got)
				}
			})
		}
	})

	t.Run("empty value", func(t *testing.T) {
		SetLevel(LevelDebug)
		got := captureOutputLog(context.TODO(), LevelDebug, Value{})
		assert.Equal(t, "", got)
	})

	t.Run("with request id", func(t *testing.T) {
		SetLevel(LevelDebug)
		ctx := WithLogContextValue(context.Background(), "req123")
		got := captureOutputLog(ctx, LevelDebug, Value{"msg": "a"})
		assert.Contains(t, got, `"rid":"req123"`)
	})

	t.Run("levels", func(t *testing.T) {
		tests := []struct {
			name  string
			level Level
			label string
		}{
			{name: "error", level: LevelError, label: "ERROR"},
			{name: "info", level: LevelInfo, label: "INFO"},
			{name: "debug", level: LevelDebug, label: "DEBUG"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				SetLevel(tt.level)
				got := captureOutputLog(context.TODO(), tt.level, Value{"msg": "a"})
				assert.Contains(t, got, `"msg":"a"`)
				assert.Contains(t, got, `"rid":""`)
				assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, tt.label))
			})
		}
	})

	t.Run("marshal error", func(t *testing.T) {
		SetLevel(LevelDebug)
		ch := make(chan int)
		got := captureOutputLog(context.TODO(), LevelDebug, Value{"msg": ch})
		assert.Contains(t, got, "marshal error: json: unsupported type: chan int")
	})
}

func TestDebug(t *testing.T) {
	SetLevel(LevelDebug)
	got := captureOutputLog(context.TODO(), LevelDebug, Value{"msg": "a"})
	assert.Contains(t, got, `"msg":"a"`)
	assert.Contains(t, got, `"rid":""`)
	assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, "DEBUG"))
}

func TestInfo(t *testing.T) {
	SetLevel(LevelInfo)
	got := captureOutputLog(context.TODO(), LevelInfo, Value{"msg": "a"})
	assert.Contains(t, got, `"msg":"a"`)
	assert.Contains(t, got, `"rid":""`)
	assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, "INFO"))
}

func TestError(t *testing.T) {
	SetLevel(LevelError)
	got := captureOutputLog(context.TODO(), LevelError, Value{"msg": "a"})
	assert.Contains(t, got, `"msg":"a"`)
	assert.Contains(t, got, `"rid":""`)
	assert.Contains(t, got, fmt.Sprintf(`"level":"%s"`, "ERROR"))
}
