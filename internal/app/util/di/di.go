// Package di provides a dependency injection.
package di

import "fmt"

var (
	di DI
)

// DI is implemented by an application which uses DI container.
type DI interface {
	// Get returns a concrete struct identified by name.
	Get(name string) (interface{}, error)
}

// Get is a helper function.
func Get(name string) (interface{}, error) {
	if di == nil {
		return nil, fmt.Errorf("not set container error")
	}
	return di.Get(name)
}

// SetDi sets DI used by an application.
func SetDi(d DI) {
	di = d
}
