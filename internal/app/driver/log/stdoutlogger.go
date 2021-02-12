// Package log provides implementation of logger.
package log

import (
	log2 "api/internal/app/util/log"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// StdoutLogger implements log.Logger.
type StdoutLogger struct {
	logger *log.Logger
	level  log2.Level
}

var _ log2.Logger = (*StdoutLogger)(nil)

// NewStdoutLogger creates new struct.
func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{
		logger: log.New(os.Stdout, "", 0),
		level:  log2.LevelInfo,
	}
}

// SetLevel sets logging level.
func (l *StdoutLogger) SetLevel(level log2.Level) {
	l.level = level
}

// Debug outputs debug log.
func (l *StdoutLogger) Debug(ctx context.Context, v log2.Value) {
	defer func() {
		// don't panic
	}()
	l.outputLog(ctx, log2.LevelDebug, v)
}

// Info outputs info log.
func (l *StdoutLogger) Info(ctx context.Context, v log2.Value) {
	defer func() {
		// don't panic
	}()
	l.outputLog(ctx, log2.LevelInfo, v)
}

// Error outputs error log.
func (l *StdoutLogger) Error(ctx context.Context, v log2.Value) {
	defer func() {
		// don't panic
	}()
	l.outputLog(ctx, log2.LevelError, v)
}

// Fatal outputs log and exit the application.
func (l *StdoutLogger) Fatal(v log2.Value) {
	defer func() {
		// don't panic
	}()
	log2.Fatal(v)
}

func (l *StdoutLogger) outputLog(ctx context.Context, level log2.Level, v log2.Value) {
	if l.level < level {
		// Ignore the log with lower priorities than the output level.
		return
	}

	if len(v) == 0 {
		return
	}

	var label string
	switch level {
	case log2.LevelError:
		label = "ERROR"
	case log2.LevelInfo:
		label = "INFO"
	case log2.LevelDebug:
		label = "DEBUG"
	}

	data := log2.Value{}
	for vk, vv := range v {
		data[vk] = vv
	}
	data["level"] = label
	data["rid"] = log2.GetRequestID(ctx)

	var msg string
	if body, err := json.Marshal(data); err == nil {
		var buf bytes.Buffer
		if json.Compact(&buf, body); err == nil {
			msg = buf.String()
		} else {
			msg = fmt.Sprintf("json compact error: %v", err)
		}
	} else {
		msg = fmt.Sprintf("marshal error: %v", err)
	}

	logger := log.New(os.Stdout, "", 0)
	logger.Println(msg)
}
