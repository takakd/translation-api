package config

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{name: "params", key: "test1", value: "value1", want: "value1"},
		{name: "no params", key: "test2", value: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mc := NewMockConfig(ctrl)
			mc.EXPECT().Get(tt.key).Return(tt.value)
			SetConfig(mc)

			got := Get(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetError(t *testing.T) {
	SetConfig(nil)
	assert.Panics(t, func() {
		Get("test1")
	})
}
