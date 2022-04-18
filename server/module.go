package server

import "go.uber.org/fx"

// Module only contains the function to start Twirp Server
var Module = fx.Invoke(StartServer)
