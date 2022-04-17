package server

import (
	"context"
	"github.com/Rippling/gocode-template/config"
	pb "github.com/Rippling/gocode-template/rpc/company"
	"go.uber.org/fx"
	"net/http"
)

type Deps struct {
	fx.In

	Lifecycle fx.Lifecycle
	config.Provider

	Handlers []pb.TwirpServer `group:"handlers"`
}

func StartServer(deps Deps) {
	deps.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			mux := http.NewServeMux()

			for _, handler := range deps.Handlers {
				mux.Handle(handler.PathPrefix(), handler)
			}

			return http.ListenAndServe(deps.ServiceConfig().Port, mux)
		},
	})
}
