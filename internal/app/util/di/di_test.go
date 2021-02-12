package di

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		v, err := Get("name")
		assert.Nil(t, v)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		tests := []struct {
			name string
			key  string
			args []interface{}
		}{
			{name: "no args", key: "name", args: nil},
			{name: "args", key: "name", args: []interface{}{"arg1", "arg2"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				md := NewMockDI(ctrl)
				want := struct{}{}
				md.EXPECT().Get(tt.key, tt.args...).Return(want, nil)

				SetDi(md)

				var (
					got interface{}
					err error
				)
				if tt.args != nil {
					got, err = Get(tt.key, tt.args...)
				} else {
					got, err = Get(tt.key)
				}
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			})
		}
	})
}
