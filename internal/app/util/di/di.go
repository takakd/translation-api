// Package di provides dependency injection container.
package di

import "fmt"

var (
	di DI
)

// DI returns interfaces implemented per environment.
type DI interface {
	Get(name string) (interface{}, error)
}

// Get is helper function of DI.Get
func Get(name string) (interface{}, error) {
	if di == nil {
		return nil, fmt.Errorf("not set container error")
	}
	return di.Get(name)
}

// SetDi sets DI used throughout the application.
func SetDi(d DI) {
	di = d
}
