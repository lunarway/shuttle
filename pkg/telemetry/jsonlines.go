package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type JsonLinesTelemetryClient struct {
	url         string
	labelPrefix string
	properties  map[string]string
	*http.Client
	logLocation string
	writeMutex  sync.Mutex
}

func (t *JsonLinesTelemetryClient) Trace(
	ctx context.Context,
	properties map[string]string,
) {
	copyHostMap(t.properties, properties)

	event := &UploadTraceEvent{
		App:        appKey,
		Timestamp:  time.Now().UTC(),
		Properties: includeContext(ctx, properties),
	}

	content, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal trace event: %s", err)
		return
	}

	if err = t.writeLogLine(ctx, content); err != nil {
		log.Printf("failed to write to file: %s", err)
		return
	}
}

var _ TelemetryClient = &JsonLinesTelemetryClient{}

type UploadTraceEvent struct {
	App        string            `json:"app"`
	Timestamp  time.Time         `json:"timestamp"`
	Properties map[string]string `json:"properties"`
}

// filename
const fileNameShuttleJsonLines = "shuttle-telemetry"

// extensions
const extensionShuttleJsonLines = ".jsonl"

func (t *JsonLinesTelemetryClient) writeLogLine(ctx context.Context, content []byte) error {
	// Lock the mutex so multiple writers don't write at the same time
	t.writeMutex.Lock()
	defer t.writeMutex.Unlock()

	runID := RunIDFrom(ctx)

	file, err := os.OpenFile(
		path.Join(
			t.logLocation,
			fmt.Sprintf("%s-%s%s", fileNameShuttleJsonLines, runID, extensionShuttleJsonLines),
		),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0o644,
	)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}
