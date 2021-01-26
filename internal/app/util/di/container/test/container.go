package test

import (
	"api/internal/app/util/di"
)

// Container implements DI.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns a concrete struct identified by name.
func (d *Container) Get(name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
