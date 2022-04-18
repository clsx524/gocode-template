package repositories

import (
	"go.uber.org/fx"
)

// Module gathers all repository components under this directory
var Module = fx.Provide(
	New,
)
