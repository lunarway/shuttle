package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type UploadTelemetryClient struct {
	url         string
	labelPrefix string
	properties  map[string]string
	*http.Client
}

func (t *UploadTelemetryClient) Trace(
	ctx context.Context,
	label string,
	options ...TelemetryOption,
) {
	properties := make(map[string]string)
	for _, o := range options {
		o(properties)
	}

	properties["shuttle.label"] = label

	err := t.uploadTrace(ctx, properties)
	if err != nil {
		fmt.Printf("failed to publish trace: %s", err)
	}
}

func (t *UploadTelemetryClient) TraceError(
	ctx context.Context,
	label string,
	err error,
	options ...TelemetryOption,
) {
	properties := make(map[string]string)
	for _, o := range options {
		o(properties)
	}

	properties["shuttle.label"] = label
	properties["phase"] = "error"
	properties["error"] = err.Error()
	err = t.uploadTrace(ctx, properties)
	if err != nil {
		fmt.Printf("failed to publish trace: %s", err)
	}
}

var _ TelemetryClient = &UploadTelemetryClient{}

type uploadTraceEvent struct {
	App        string            `json:"app"`
	Type       string            `json:"type"`
	Timestamp  time.Time         `json:"timestamp"`
	Properties map[string]string `json:"properties"`
}

func (t *UploadTelemetryClient) uploadTrace(
	ctx context.Context,
	properties map[string]string,
) error {
	copyHostMap(t.properties, properties)

	event := &uploadTraceEvent{
		App:        appKey,
		Type:       "",
		Timestamp:  time.Now().UTC(),
		Properties: mergeMaps(ctx, properties),
	}

	content, err := json.Marshal(event)
	if err != nil {
		return err
	}

	resp, err := t.Client.Post(t.url, "application/json", bytes.NewReader(content))
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
