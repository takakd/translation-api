// Package log provides logging.
package log

import (
	"context"
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

// Logger is implemented by an application which uses Logger.
type Logger interface {
	SetLevel(level Level)
	Debug(ctx context.Context, v Value)
	Info(ctx context.Context, v Value)
	Error(ctx context.Context, v Value)
	Fatal(v Value)
}

// Use this interface for logging.
var logger Logger

// SetLogger sets the logger which is called log.Info, log.Error...
func SetLogger(l Logger) {
	logger = l
}

// SetLevel sets logging level.
func SetLevel(level Level) {
	if logger != nil {
		logger.SetLevel(level)
	}
}

// Debug outputs debug log.
func Debug(ctx context.Context, v Value) {
	defer func() {
		// don't panic
	}()
	logger.Debug(ctx, v)
}

// Info outputs info log.
func Info(ctx context.Context, v Value) {
	defer func() {
		// don't panic
	}()
	logger.Info(ctx, v)
}

// Error outputs error log.
func Error(ctx context.Context, v Value) {
	defer func() {
		// don't panic
	}()
	logger.Error(ctx, v)
}

// Fatal outputs log and exit the application.
func Fatal(v Value) {
	defer func() {
		// don't panic
	}()
	logger.Fatal(v)
}
