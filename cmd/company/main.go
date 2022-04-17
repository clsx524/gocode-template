package main

import (
	"context"
	client "github.com/Rippling/gocode-template/client"
	config "github.com/Rippling/gocode-template/config"
	"github.com/Rippling/gocode-template/handler"
	"github.com/Rippling/gocode-template/repository"
	"github.com/Rippling/gocode-template/server"
	"github.com/Rippling/gocode-template/service"
	"go.uber.org/fx"
)

func opts() fx.Option {
	return fx.Options(
		config.Module,
		client.Module,
		repository.Module,
		service.Module,
		handler.Module,
		server.Module,
	)
}

func main() {
	if err := fx.New(opts()).Start(context.Background()); err != nil {
		panic(err)
	}
}
