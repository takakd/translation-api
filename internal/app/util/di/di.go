// Package di provides a dependency injection.
package di

import "fmt"

var (
	di DI
)

// DI is implemented by an application which uses DI container.
// Pass args if the concrete struct needs arguments.
type DI interface {
	// Get returns a concrete struct identified by name.
	Get(name string, args ...interface{}) (interface{}, error)
}

// Get is a helper function.
func Get(name string, args ...interface{}) (interface{}, error) {
	if di == nil {
		return nil, fmt.Errorf("not set container error")
	}
	return di.Get(name, args...)
}

// SetDi sets DI used by an application.
func SetDi(d DI) {
	di = d
}
