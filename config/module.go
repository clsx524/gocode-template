package config

import "go.uber.org/fx"

// Module provides the configuration Provider
var Module = fx.Provide(New)
