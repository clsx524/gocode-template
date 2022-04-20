package main

import (
	"context"
	"github.com/clsx524/gocode-template/apis"
	client "github.com/clsx524/gocode-template/clients"
	config "github.com/clsx524/gocode-template/config"
	"github.com/clsx524/gocode-template/repositories"
	"github.com/clsx524/gocode-template/server"
	"github.com/clsx524/gocode-template/services"
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
