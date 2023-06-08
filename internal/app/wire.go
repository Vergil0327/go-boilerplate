//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
)

func BuildInjector() (*Injector, error) {
	wire.Build(
		InitGinEngine,
		InjectorSet,
	)

	return new(Injector), nil
}
