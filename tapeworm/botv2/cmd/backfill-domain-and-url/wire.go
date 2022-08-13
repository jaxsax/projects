//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jaxsax/projects/tapeworm/botv2/internal"
)

func initialize() (*App, error) {
	wire.Build(
		internal.CommonSet,
		wire.Struct(new(App), "*"),
	)

	return nil, nil
}
