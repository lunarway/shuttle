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
	properties map[string]string,
) {
	properties = copyHostMap(t.properties, properties)

	content, err := json.Marshal(properties)
	if err != nil {
		log.Printf("failed to serialize properties")
		return
	}

	log.Printf("%s: %s\n", t.labelPrefix, string(content))
}

var _ TelemetryClient = &LoggingTelemetryClient{}
