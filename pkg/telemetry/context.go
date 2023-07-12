package telemetry

import (
	"context"
	"os"

	"github.com/google/uuid"
)

const envContextID = "SHUTTLE_CONTEXT_ID"

func WithContextID(ctx context.Context) context.Context {
	if context_id := os.Getenv(envContextID); context_id != "" {
		return context.WithValue(ctx, telemetryContextID, context_id)
	}

	return context.WithValue(ctx, telemetryContextID, uuid.New().String())
}

func ContextIDFrom(ctx context.Context) string {
	if contextID, ok := ctx.Value(envContextID).(string); ok {
		return contextID
	}
	return ""
}

func WithContextValue(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
