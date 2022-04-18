package services

import (
	"go.uber.org/fx"
)

// Module gathers all service components defined under services directory
var Module = fx.Provide(New)
