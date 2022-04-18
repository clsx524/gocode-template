package clients

import "go.uber.org/fx"

// Module gathers all modules defined under clients directory
var Module = fx.Provide(
	ProvideInstrumenter,
	ProvideMongoClient,
)
