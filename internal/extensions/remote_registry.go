package extensions

import "context"

type RemoteRegistry interface {
	Publish(ctx context.Context) error
}

func NewRemoteRegistry(registry string) {}
