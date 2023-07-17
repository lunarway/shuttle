package telemetry

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
		availabilityCheck AvailabilityCheckFunc
		lock              LockFunc
	}

	// UploadFunc handles upload for a set of trace events
	UploadFunc = func(ctx context.Context, url string, event []UploadTraceEvent) error

	// GetTelemetryFilesFunc fetches trace event file name from a certain location. Each read is delegated to GetTelemetryFileFunc
	GetTelemetryFilesFunc = func(ctx context.Context, location string) ([]string, error)

	// GetTelemetryFileFunc reads the trace event files and returns a set of tracevents and the option to delete the file after upload is finished
	GetTelemetryFileFunc = func(ctx context.Context, telemetryFilePath string) ([]UploadTraceEvent, func(ctx context.Context) error, error)

	// AvailabilityCheckFunc gets whether the telemetry uploader is available for upload
	AvailabilityCheckFunc = func(ctx context.Context) (bool, error)

	// LockFunc  makes sure only a single upload process is run pr storage location
	LockFunc = func(ctx context.Context) (UnlockFunc, bool, error)

	// UnluckFunc clears the locks set for the storage location
	UnlockFunc = func(ctx context.Context) error

	UploadOptions = func(*TelemetryUploader)
)

// WithRate sets the current rate of which events will be uploaded
func WithRate(rate time.Duration) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.rate = &rate
	}
}

// WithUploadFunction sets the upload function to something custom
func WithUploadFunction(uploadFunc UploadFunc) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.upload = uploadFunc
	}
}

// WithGetTelemetryFiles sets the get telemetry files
func WithGetTelemetryFiles(getTelemetryFilesFunc GetTelemetryFilesFunc) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.getTelemetryFiles = getTelemetryFilesFunc
	}
}

// WithGetTelemetryFile sets the get telemetry file func
func WithGetTelemetryFile(getTelemetryFileFunc GetTelemetryFileFunc) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.getTelemetryFile = getTelemetryFileFunc
	}
}

// WithRemoteLogLocation sets where shuttle telemetry files are located
func WithRemoteLogLocation(location string) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.storageLocation = location
	}
}

// WithCleanUp determines whether to remove telemetry files after they've been read.
// There are no consistency guarantees that no duplicates will be uploaded on error if run, or dropped traces
func WithCleanUp(enabled bool) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.cleanUp = enabled
	}
}

// WithAvailabilityCheck adds a check for the upload location, this is useful if there are scenarios where the upload
// process shouldn't be run. I.e. you're not on a vpn, or on a slow internet connection
func WithAvailabilityCheck(url string) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.availabilityCheck = func(ctx context.Context) (bool, error) {
			return availabilityCheck(ctx, url)
		}
	}
}

// WithDefaultAvailabilityCheck sets a Noop availability check
func WithDefaultAvailabilityCheck() UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.availabilityCheck = func(ctx context.Context) (bool, error) {
			return true, nil
		}
	}
}

// WithFileLock this adds a file lock at a certain location,
// this is useful if there may be concurrent/parallel processes on the same storage location
func WithFileLock(storageLocation string) UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.lock = lockFunc(storageLocation)
	}
}

// WithNoLock removes the default file lock
func WithNoLock() UploadOptions {
	return func(tu *TelemetryUploader) {
		tu.lock = func(ctx context.Context) (UnlockFunc, bool, error) {
			return func(ctx context.Context) error { return nil }, false, nil
		}
	}
}

func NewTelemetryUploader(url string, options ...UploadOptions) *TelemetryUploader {
	storageLocation := getRemoteLogLocation()
	uploader := &TelemetryUploader{
		url:               url,
		upload:            upload,
		getTelemetryFiles: getTelemetryFiles,
		getTelemetryFile:  getTelemetryFile,
		availabilityCheck: func(ctx context.Context) (bool, error) {
			return true, nil
		},
		storageLocation: storageLocation,
		cleanUp:         true,
		lock:            lockFunc(storageLocation),
	}

	for _, o := range options {
		o(uploader)
	}

	return uploader
}

// Upload kicks off the file upload process. This involves reading trace event files and uploading them to a determined location.
// See Options (With*) for extra documentation
func (tu *TelemetryUploader) Upload(ctx context.Context) error {
	// Makes sure only a single upload process is run for each instance
	tu.uploadmu.Lock()
	defer tu.uploadmu.Unlock()

	unlock, locked, err := tu.lock(ctx)
	if err != nil {
		return err
	}
	if locked {
		log.Println("file is already locked returning")
		return nil
	}
	defer func(ctx context.Context, unlock UnlockFunc) {
		return
		if err := unlock(ctx); err != nil {
			log.Printf("failed to clean up lock: %s", err)
		}
	}(ctx, unlock)

	ok, err := tu.availabilityCheck(ctx)
	if err != nil {
		return fmt.Errorf("checking for endpoint failed so bad that process needs to stop: %w", err)
	}

	if !ok {
		log.Println("endpoint was not ready, try again at some other time")
		return nil
	}

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

func availabilityCheck(ctx context.Context, url string) (bool, error) {
	select {
	case <-ctx.Done():
		log.Println("availability deadline exceeded")
		return false, errors.New("met deadline returning")
	default:
		resp, err := http.DefaultClient.Get(url)
		if err != nil {
			log.Printf("checking endpoint failed with: %s", err)
			return false, nil
		}

		if resp.StatusCode != 200 {
			log.Printf("status code is not 200: status=%d", resp.StatusCode)
			return false, nil
		}

		return true, nil
	}
}

func lockFunc(storageLocation string) LockFunc {
	handleFolderExists := func() error {
		if _, err := os.Stat(storageLocation); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if err = os.MkdirAll(
					storageLocation,
					0o700,
				); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	}

	checkLockExists := func(lockFile string) (exists bool, err error) {
		file, err := os.Stat(lockFile)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}

			return false, err
		}

		if file.ModTime().Add(time.Minute * 5).Before(time.Now()) {
			return false, nil
		}

		return true, nil

	}

	unlockFunc := func(lockFile string) func(ctx context.Context) error {
		return func(ctx context.Context) error {
			if err := os.Remove(lockFile); err != nil {
				return err
			}

			return nil
		}
	}

	return func(ctx context.Context) (UnlockFunc, bool, error) {
		if err := handleFolderExists(); err != nil {
			return nil, false, err
		}
		lockFile := path.Join(storageLocation, ".shuttle-telemetry-lock")
		exists, err := checkLockExists(lockFile)
		if err != nil {
			return nil, false, err
		}
		if exists {
			return nil, true, nil
		}

		if _, err := os.Create(lockFile); err != nil {
			return nil, false, err
		}

		return unlockFunc(lockFile), false, nil
	}
}
