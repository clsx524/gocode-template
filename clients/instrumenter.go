package clients

import (
	"context"
	"github.com/clsx524/gocode-template/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Instrumenter defines the tools to instrument the code
// including logging, metrics, profiling etc.
type Instrumenter interface {
	Logger() *zap.SugaredLogger
}

// InstrumenterDeps defines the dependencies used in the instrumenter module
type InstrumenterDeps struct {
	fx.In

	config.Provider

	Lifecycle fx.Lifecycle
}

type instrumenter struct {
	logger *zap.SugaredLogger
}

func (i *instrumenter) Logger() *zap.SugaredLogger {
	return i.logger
}

// ProvideInstrumenter provides the Instrumenter for use in other modules
func ProvideInstrumenter(deps InstrumenterDeps) Instrumenter {
	return &instrumenter{logger: newLogger(deps)}
}

func newLogger(deps InstrumenterDeps) *zap.SugaredLogger {
	logger, err := deps.LoggerConfig().Build()
	if err != nil {
		panic(err)
	}

	deps.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = logger.Sync()
			return nil
		},
	})

	return logger.Sugar()
}
