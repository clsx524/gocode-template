package client

import (
	"context"
	"github.com/Rippling/gocode-template/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Instrumenter interface {
	Logger() *zap.SugaredLogger
}

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
