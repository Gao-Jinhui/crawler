package engine

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/spider"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	WorkCount   int
	Fetcher     spider.Fetcher
	Storage     collector.Storage
	Logger      *zap.Logger
	TaskSeeds   []*spider.Task
	registryURL string
	scheduler   Scheduler
}

var defaultOptions = options{
	Logger: zap.NewNop(),
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}
func WithFetcher(fetcher spider.Fetcher) Option {
	return func(opts *options) {
		opts.Fetcher = fetcher
	}
}

func WithWorkCount(workCount int) Option {
	return func(opts *options) {
		opts.WorkCount = workCount
	}
}

func WithSeeds(seed []*spider.Task) Option {
	return func(opts *options) {
		opts.TaskSeeds = seed
	}
}

func WithScheduler(scheduler Scheduler) Option {
	return func(opts *options) {
		opts.scheduler = scheduler
	}
}

func WithStorage(s collector.Storage) Option {
	return func(opts *options) {
		opts.Storage = s
	}
}

func WithregistryURL(registryURL string) Option {
	return func(opts *options) {
		opts.registryURL = registryURL
	}
}
