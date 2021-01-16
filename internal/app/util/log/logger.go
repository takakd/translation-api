// Package log provides logging feature.
package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
)

// Level represents the type of logging level.
type Level int

const (
	// LevelError outputs Error logs.
	LevelError Level = iota
	// LevelInfo outputs Error and Info logs.
	LevelInfo
	// LevelDebug outputs all level logs.
	LevelDebug
)

// Value is contents of log.
type Value map[string]interface{}

// StringValue creates a simple string log.
func StringValue(s string) Value {
	return map[string]interface{}{"msg": s}
}

// Output log above this level.
var logLevel Level

// SetLevel sets logging level.
func SetLevel(level Level) {
	logLevel = level
}

// Debug outputs debug log.
func Debug(ctx context.Context, v Value) {
	defer func() {
		// don't panic
	}()
	outputLog(ctx, LevelDebug, v)
}

// Info outputs info log.
func Info(ctx context.Context, v Value) {
	defer func() {
		// don't panic
	}()
	outputLog(ctx, LevelInfo, v)
}

// Error outputs info log.
func Error(ctx context.Context, v Value) {
	defer func() {
		// don't panic
	}()
	outputLog(ctx, LevelError, v)
}

// Fatal calls log.Fatal.
func Fatal(ctx context.Context, v ...interface{}) {
	log.Fatal(v...)
}

func outputLog(ctx context.Context, level Level, v Value) {
	if logLevel < level {
		// Ignore the log with lower priorities than the output level.
		return
	}

	if len(v) == 0 {
		return
	}

	var label string
	switch level {
	case LevelError:
		label = "ERROR"
	case LevelInfo:
		label = "INFO"
	case LevelDebug:
		label = "DEBUG"
	}

	data := Value{}
	for vk, vv := range v {
		data[vk] = vv
	}
	data["level"] = label
	data["rid"] = getRequestID(ctx)

	var msg string
	if body, err := json.Marshal(data); err == nil {
		var buf bytes.Buffer
		if json.Compact(&buf, body); err == nil {
			msg = buf.String()
		} else {
			msg = "marshal error in logging"
		}
	} else {
		msg = "marshal error in logging"
	}

	logger := log.New(os.Stdout, "", 0)
	logger.Println(msg)
}
