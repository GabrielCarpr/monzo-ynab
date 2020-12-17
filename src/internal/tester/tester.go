package tester

import (
	"github.com/sarulabs/di/v2"
)

// SetDep is a shorthand for setting a dependency on a DI builder
func SetDep(builder *di.Builder, name string, object interface{}) {
	builder.Add(di.Def{
		Name: name,
		Build: func(ctn di.Container) (interface{}, error) {
			return object, nil
		},
	})
}
