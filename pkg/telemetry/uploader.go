package telemetry

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type (
	TelemetryUploader struct {
		url      string
		rate     *time.Duration
		uploadmu sync.Mutex

		storageLocation string
		cleanUp         bool

		upload            UploadFunc
		getTelemetryFiles GetTelemetryFilesFunc
		getTelemetryFile  GetTelemetryFileFunc
	}

	UploadFunc            = func(ctx context.Context, url string, event []UploadTraceEvent) error
	GetTelemetryFilesFunc = func(ctx context.Context, location string) ([]string, error)
	GetTelemetryFileFunc  = func(ctx context.Context, telemetryFilePath string) ([]UploadTraceEvent, func(ctx context.Context) error, error)

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

func WithCleanUp(enabled bool) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.cleanUp = enabled
	}
}

func NewTelemetryUploader(url string, options ...UploadOptions) *TelemetryUploader {
	uploader := &TelemetryUploader{
		url:               url,
		upload:            upload,
		getTelemetryFiles: getTelemetryFiles,
		getTelemetryFile:  getTelemetryFile,
		storageLocation:   getRemoteLogLocation(),
		cleanUp:           true,
	}

	for _, o := range options {
		o(uploader)
	}

	return uploader
}

func (tu *TelemetryUploader) Upload(ctx context.Context) error {
	tu.uploadmu.Lock()
	defer tu.uploadmu.Unlock()

	files, err := tu.getTelemetryFiles(ctx, tu.storageLocation)
	if err != nil {
		return fmt.Errorf("failed to get telemetry files: %w", err)
	}

	egrp, ctx := errgroup.WithContext(ctx)
	for _, file := range files {
		file := file
		egrp.Go(func() error {
			events, cleanUpFunc, err := tu.getTelemetryFile(ctx, file)
			if err != nil {
				return fmt.Errorf(
					"failed to read events from shuttle telemetry file: %s, err: %w",
					file,
					err,
				)
			}

			if err := tu.upload(ctx, tu.url, events); err != nil {
				return fmt.Errorf("failed to upload events: %w", err)
			}

			if tu.cleanUp {
				if err := cleanUpFunc(ctx); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := egrp.Wait(); err != nil {
		return err
	}

	return nil
}

func upload(ctx context.Context, url string, events []UploadTraceEvent) error {
	content, err := json.Marshal(events)
	if err != nil {
		return err
	}

	client := http.DefaultClient

	resp, err := client.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(
			"failed to push trace event with status code: %d, reason: %s",
			resp.StatusCode,
			string(body),
		)
	}

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
			// TODO: only read files older than a certain threshold

			shuttleTelemetryFiles = append(shuttleTelemetryFiles, path.Join(location, fileName))
		}
	}

	return shuttleTelemetryFiles, nil
}

func getTelemetryFile(
	ctx context.Context,
	shuttleTelemetryFilePath string,
) ([]UploadTraceEvent, func(ctx context.Context) error, error) {
	file, err := os.Open(shuttleTelemetryFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []UploadTraceEvent{}, nil, nil
		}
		return nil, nil, err
	}

	events := make([]UploadTraceEvent, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var event UploadTraceEvent
		err := json.Unmarshal(line, &event)
		if err != nil {
			return nil, nil, err
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	cleanUp := func(ctx context.Context) error {
		return os.Remove(shuttleTelemetryFilePath)
	}

	return events, cleanUp, nil
}
