//go:build wireinject
// +build wireinject

package app

import (
	"boilerplate/internal/app/router"
	"github.com/google/wire"
)

func BuildInjector() (*Injector, error) {
	wire.Build(
		router.RouterSet,
		InitGinEngine,
		InjectorSet,
	)

	return new(Injector), nil
}
