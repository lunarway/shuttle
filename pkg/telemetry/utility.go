package telemetry

import (
	"context"

	"github.com/google/uuid"
)

const telemetryContextID = "shuttle.contextID"

func WithContextID(ctx context.Context) context.Context {
	return context.WithValue(ctx, telemetryContextID, uuid.New().String())
}

func mergeMaps(ctx context.Context, properties map[string]string) map[string]string {
	if runID, ok := ctx.Value(telemetryContextID).(string); ok && runID != "" {
		properties[telemetryContextID] = runID
	}

	return properties
}
