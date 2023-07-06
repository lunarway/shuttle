package telemetry

import (
	"context"
	"encoding/json"
	"log"
)

type LoggingTelemetryClient struct {
	labelPrefix string
	properties  map[string]string
}

func (t *LoggingTelemetryClient) Trace(
	ctx context.Context,
	label string,
	options ...TelemetryOption,
) {
	var properties map[string]string
	for _, o := range options {
		o(properties)
	}

	content, err := json.Marshal(mergeMaps(ctx, properties))
	if err != nil {
		log.Printf("failed to serialize properties")
		return
	}

	log.Printf("%s.%s: %s\n", t.labelPrefix, label, string(content))
}

func (t *LoggingTelemetryClient) TraceError(
	ctx context.Context,
	label string,
	err error,
	options ...TelemetryOption,
) {
	var properties map[string]string
	for _, o := range options {
		o(properties)
	}

	content, marshalErr := json.Marshal(mergeMaps(ctx, properties))
	if marshalErr != nil {
		log.Printf("failed to serialize properties")
		return
	}

	log.Printf("%s.%s: (error=%s) %s\n", t.labelPrefix, label, err, string(content))
}

var _ TelemetryClient = &LoggingTelemetryClient{}
