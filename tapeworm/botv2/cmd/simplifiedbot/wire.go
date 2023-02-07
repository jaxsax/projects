//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/dimension"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/httpserver"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/services"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/telegrampoller"
)

func initialize() (*App, error) {
	wire.Build(
		internal.CommonSet,
		httpserver.New,
		telegrampoller.New,
		services.Set,
		dimension.Set,
		wire.Struct(new(App), "*"),
	)

	return nil, nil
}
