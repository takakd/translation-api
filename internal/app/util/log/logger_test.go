// Package log provides logging.
package log

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStringValue(t *testing.T) {
	got := StringValue("a")
	assert.Equal(t, Value{"msg": "a"}, got)
}

func TestSetLogger(t *testing.T) {
	var i Logger
	SetLogger(i)
	assert.Equal(t, i, logger)
}

func TestSetLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)

	level := LevelInfo
	ml.EXPECT().SetLevel(level)
	SetLogger(ml)

	SetLevel(level)
}

func TestDebug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)

	ctx := context.TODO()
	v := Value{}
	ml.EXPECT().Debug(ctx, v)
	SetLogger(ml)

	Debug(ctx, v)
}

func TestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)

	ctx := context.TODO()
	v := Value{}
	ml.EXPECT().Info(ctx, v)
	SetLogger(ml)

	Info(ctx, v)
}

func TestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)

	ctx := context.TODO()
	v := Value{}
	ml.EXPECT().Error(ctx, v)
	SetLogger(ml)

	Error(ctx, v)
}
