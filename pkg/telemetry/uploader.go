package telemetry

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type (
	TelemetryUploader struct {
		url      string
		rate     *time.Duration
		uploadmu sync.Mutex

		storageLocation string

		upload            UploadFunc
		getTelemetryFiles GetTelemetryFilesFunc
		getTelemetryFile  GetTelemetryFileFunc
	}

	UploadFunc            = func(ctx context.Context) error
	GetTelemetryFilesFunc = func(ctx context.Context, location string) ([]string, error)
	GetTelemetryFileFunc  = func(ctx context.Context, telemetryFilePath string) ([]uploadTraceEvent, error)

	UploadOptions = func(*TelemetryUploader)
)

func WithRate(rate time.Duration) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.rate = &rate
	}
}

func WithUploadFunction(uploadFunc UploadFunc) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.upload = uploadFunc
	}
}

func WithGetTelemetryFiles(getTelemetryFilesFunc GetTelemetryFilesFunc) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.getTelemetryFiles = getTelemetryFilesFunc
	}
}

func WithGetTelemetryFile(getTelemetryFileFunc GetTelemetryFileFunc) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.getTelemetryFile = getTelemetryFileFunc
	}
}

func WithRemoteLogLocation(location string) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.storageLocation = location
	}
}

func NewTelemetryUploader(url string, options ...UploadOptions) *TelemetryUploader {
	uploader := &TelemetryUploader{
		url:               url,
		upload:            upload,
		getTelemetryFiles: getTelemetryFiles,
		getTelemetryFile:  getTelemetryFile,
		storageLocation:   getRemoteLogLocation(),
	}

	for _, o := range options {
		o(uploader)
	}

	return uploader
}

func (tu *TelemetryUploader) Upload(ctx context.Context) {
	tu.getTelemetryFiles(ctx, tu.storageLocation)
}

func upload(ctx context.Context) error {
	return nil
}

func getTelemetryFiles(ctx context.Context, location string) ([]string, error) {
	if _, err := os.Stat(location); errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}

	files, err := os.ReadDir(location)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}

	shuttleTelemetryFiles := make([]string, 0)
	for _, file := range files {
		fileName := file.Name()
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		isFile := fileInfo.Mode().IsRegular()
		if strings.HasPrefix(fileName, fileNameShuttleJsonLines) &&
			strings.HasSuffix(fileName, extensionShuttleJsonLines) &&
			isFile {
			shuttleTelemetryFiles = append(shuttleTelemetryFiles, path.Join(location, fileName))
		}
	}

	return shuttleTelemetryFiles, nil
}

func getTelemetryFile(
	ctx context.Context,
	shuttleTelemetryFilePath string,
) ([]uploadTraceEvent, error) {
	return nil, nil
}
