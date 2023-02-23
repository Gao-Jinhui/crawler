package mysql

import "go.uber.org/zap"

type options struct {
	logger *zap.Logger
}

type Option func(opts *options)

var defaultOptions = options{
	logger: zap.NewNop(),
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}
