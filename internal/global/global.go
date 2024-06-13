package global

import "os"

// GlobalStore represents the ~/.shuttle folder it acts as an abstraction for said folder and ensures operations against it are consistent and controlled
type GlobalStore struct {
	options *GlobalStoreOptions
}

type GlobalStoreOption func(options *GlobalStoreOptions)

func WithShuttleConfig(shuttleConfig string) GlobalStoreOption {
	return func(options *GlobalStoreOptions) {
		options.ShuttleConfig = shuttleConfig
	}
}

type GlobalStoreOptions struct {
	ShuttleConfig string
}

func newDefaultGlobalStoreOptions() *GlobalStoreOptions {
	return &GlobalStoreOptions{
		ShuttleConfig: "$HOME/.shuttle",
	}
}

func NewGlobalStore(options ...GlobalStoreOption) *GlobalStore {
	defaultOptions := newDefaultGlobalStoreOptions()
	for _, opt := range options {
		opt(defaultOptions)
	}

	return &GlobalStore{
		options: defaultOptions,
	}
}

func (gs *GlobalStore) Root() string {
	return os.ExpandEnv(gs.options.ShuttleConfig)
}
