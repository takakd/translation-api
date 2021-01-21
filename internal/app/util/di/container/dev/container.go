package dev

import (
	"api/internal/app/util/di"
)

// Container implements DI on developer env.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns interfaces corresponding name.
func (d *Container) Get(name string) (interface{}, error) {
	return nil, nil
}
