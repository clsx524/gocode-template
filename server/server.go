package server

import (
	"context"
	"github.com/clsx524/gocode-template/config"
	pb "github.com/clsx524/gocode-template/rpc/company"
	"go.uber.org/fx"
	"net/http"
)

// Deps defines all dependencies used to start Twirp Server
type Deps struct {
	fx.In

	Lifecycle fx.Lifecycle
	config.Provider

	// APIs gathers all Twirp Server in a list
	// Each Twirp Server *needs* to have the group tag in FX
	APIs []pb.TwirpServer `group:"apis"`
}

// StartServer defines the function to start Twirp Server Mux
func StartServer(deps Deps) {
	deps.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			mux := http.NewServeMux()

			for _, api := range deps.APIs {
				mux.Handle(api.PathPrefix(), api)
			}

			return http.ListenAndServe(deps.ServiceConfig().Port, mux)
		},
	})
}
