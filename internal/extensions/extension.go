package extensions

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/lunarway/shuttle/internal/global"
)

// Extension is the descriptor of a single extension, it is used to add description to the cli, as well as calling the specific extension in question
type Extension struct {
	os          string
	arch        string
	globalStore *global.GlobalStore
	remote      *registryExtension
}

func newExtensionFromRegistry(globalStore *global.GlobalStore, registryExtension *registryExtension) (*Extension, error) {
	return &Extension{
		os:          runtime.GOOS,
		arch:        runtime.GOARCH,
		globalStore: globalStore,
		remote:      registryExtension,
	}, nil
}

func (e *Extension) Ensure(ctx context.Context) error {
	extensionsCachePath := getExtensionsCachePath(e.globalStore)
	binaryName := e.getExtensionBinaryName()
	if err := ensureExists(extensionsCachePath); err != nil {
		return fmt.Errorf("failed to create cache path: %w", err)
	}

	binaryPath := path.Join(extensionsCachePath, binaryName)
	if exists(binaryPath) {
		// TODO: do a checksum chck
		//return nil
	}

	downloadLink := e.getRemoteBinaryDownloadLink()
	if downloadLink == nil {
		return fmt.Errorf("failed to find a valid extension matching your os and architecture")
	}

	downloader, err := NewDownloader(downloadLink)
	if err != nil {
		return err
	}

	if err := downloader.Download(ctx, binaryPath); err != nil {
		return err
	}

	return nil
}

func (e *Extension) Name() string {
	return e.remote.Name
}

func (e *Extension) Version() string {
	return e.remote.Version
}

func (e *Extension) Description() string {
	return e.remote.Description
}

func (e *Extension) getExtensionBinaryName() string {
	return e.remote.Name
}

func (e *Extension) FullPath() string {
	return path.Join(getExtensionsCachePath(e.globalStore), e.Name())
}

func (e *Extension) getRemoteBinaryDownloadLink() *registryExtensionDownloadLink {
	for _, download := range e.remote.DownloadUrls {
		if download.Os == e.os &&
			download.Architecture == e.arch {
			return &download
		}
	}

	return nil
}
