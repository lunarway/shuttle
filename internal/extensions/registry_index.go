package extensions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

type (
	registryIndex struct {
		registryPath string
	}

	registryExtensionDownloadLink struct {
		Architecture string `json:"architecture"`
		Os           string `json:"os"`
		Url          string `json:"url"`
		Checksum     string `json:"checksum"`
		Provider     string `json:"provider"`
	}

	registryExtension struct {
		Name         string                          `json:"name"`
		Description  string                          `json:"description"`
		Version      string                          `json:"version"`
		DownloadUrls []registryExtensionDownloadLink `json:"downloadUrls"`
	}
)

func newRegistryIndex(registryPath string) *registryIndex {
	return &registryIndex{
		registryPath: registryPath,
	}
}

func (r *registryIndex) getExtensions(_ context.Context) ([]registryExtension, error) {
	contents, err := os.ReadDir(r.getIndexPath())
	if err != nil {
		return nil, fmt.Errorf("failed to list index in registry: %s, %w", r.getIndexPath(), err)
	}

	extensions := make([]registryExtension, 0)
	for _, dir := range contents {
		if !dir.IsDir() {
			continue
		}

		extensionPath := path.Join(r.getIndexPath(), dir.Name(), "shuttle-extension.json")

		extensionContent, err := os.ReadFile(extensionPath)
		if err != nil {
			log.Printf("failed to get extension: %s, skipping extension, the extension might be invalid at: %s, please contact your admin", err.Error(), extensionPath)
			continue
		}

		var extension registryExtension
		if err := json.Unmarshal(extensionContent, &extension); err != nil {
			return nil, fmt.Errorf("failed unmarshal extension at path: %s, err: %w", extensionPath, err)
		}

		extensions = append(extensions, extension)
	}

	return extensions, nil
}

func (r *registryIndex) getIndexPath() string {
	return path.Join(r.registryPath, "index")
}
