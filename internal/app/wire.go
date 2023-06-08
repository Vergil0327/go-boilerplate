//go:build wireinject
// +build wireinject

package app

import (
	"boilerplate/internal/app/router"
	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		router.RouterSet,
		InitGinEngine,
		InjectorSet,
	)

	return new(Injector), nil, nil
}
