package client

import "go.uber.org/fx"

var Module = fx.Provide(
	ProvideInstrumenter,
	ProvideMongoClient,
)
