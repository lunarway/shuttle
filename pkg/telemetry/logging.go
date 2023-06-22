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
	properties map[string]string,
) {
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
	properties map[string]string,
) {
	content, marshalErr := json.Marshal(mergeMaps(ctx, properties))
	if marshalErr != nil {
		log.Printf("failed to serialize properties")
		return
	}

	log.Printf("%s.%s: (error=%s) %s\n", t.labelPrefix, label, err, string(content))
}

var _ TelemetryClient = &LoggingTelemetryClient{}
