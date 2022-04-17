package server

import "go.uber.org/fx"

var Module = fx.Invoke(StartServer)
