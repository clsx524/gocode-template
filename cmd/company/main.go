package main

import (
	"context"
	"github.com/Rippling/gocode-template/apis"
	client "github.com/Rippling/gocode-template/clients"
	config "github.com/Rippling/gocode-template/config"
	"github.com/Rippling/gocode-template/repositories"
	"github.com/Rippling/gocode-template/server"
	"github.com/Rippling/gocode-template/services"
	"go.uber.org/fx"
)

func opts() fx.Option {
	return fx.Options(
		config.Module,
		client.Module,
		repositories.Module,
		services.Module,
		apis.Module,
		server.Module,
	)
}

func main() {
	if err := fx.New(opts()).Start(context.Background()); err != nil {
		panic(err)
	}
}
