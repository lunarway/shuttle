package telemetry

import (
	"context"
	"sync"
	"time"
)

type (
	TelemetryUploader struct {
		url      string
		rate     *time.Duration
		uploadmu sync.Mutex

		upload            UploadFunc
		getTelemetryFiles GetTelemetryFilesFunc
	}

	UploadFunc            = func(ctx context.Context) error
	GetTelemetryFilesFunc = func(ctx context.Context) error

	UploadOptions = func(*TelemetryUploader)
)

func WithRate(rate time.Duration) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.rate = &rate
	}
}

func WithUploadFunction(uploadFunc func(ctx context.Context) error) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.upload = uploadFunc
	}
}

func WithGetTelemetryFiles(getTelemetryFilesFunc func(ctx context.Context) error) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.getTelemetryFiles = getTelemetryFilesFunc
	}
}

func NewTelemetryUploader(url string, options ...UploadOptions) *TelemetryUploader {
	uploader := &TelemetryUploader{
		url:               url,
		upload:            upload,
		getTelemetryFiles: getTelemetryFiles,
	}

	for _, o := range options {
		o(uploader)
	}

	return uploader
}

func (tu *TelemetryUploader) Upload(ctx context.Context) {
}

func upload(ctx context.Context) error            {}
func getTelemetryFiles(ctx context.Context) error {}
