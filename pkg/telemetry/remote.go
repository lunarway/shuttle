package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type JsonLinesTelemetryClient struct {
	url         string
	labelPrefix string
	properties  map[string]string
	*http.Client
	logLocation string
}

func (t *JsonLinesTelemetryClient) Trace(
	ctx context.Context,
	properties map[string]string,
) {
	copyHostMap(t.properties, properties)

	event := &uploadTraceEvent{
		App:        appKey,
		Timestamp:  time.Now().UTC(),
		Properties: mergeMaps(ctx, properties),
	}

	content, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to write to file: %s", err)
	}

	if err = t.writeLogLine(ctx, content); err != nil {
		log.Printf("failed to write to file: %s", err)
	}
}

var _ TelemetryClient = &JsonLinesTelemetryClient{}

type uploadTraceEvent struct {
	App        string            `json:"app"`
	Timestamp  time.Time         `json:"timestamp"`
	Properties map[string]string `json:"properties"`
}

//func (t *JsonLinesTelemetryClient) upload(ctx context.Context, content []byte) error {
//	resp, err := t.Client.Post(t.url, "application/json", bytes.NewReader(content))
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode > 299 {
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			return err
//		}
//		return fmt.Errorf(
//			"failed to push trace event with status code: %d, reason: %s",
//			resp.StatusCode,
//			string(body),
//		)
//	}
//
//	return nil
//}

// filename
const fileNameShuttleJsonLines = "shuttle-telemetry"

// extensions
const extensionShuttleJsonLines = ".jsonl"

func (t *JsonLinesTelemetryClient) writeLogLine(ctx context.Context, content []byte) error {
	file, err := os.OpenFile(
		path.Join(
			t.logLocation,
			fmt.Sprintf("%s%s", fileNameShuttleJsonLines, extensionShuttleJsonLines),
		),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0644,
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
