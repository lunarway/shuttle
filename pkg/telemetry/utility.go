package telemetry

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/matishsiao/goInfo"
)

const (
	TelemetryContextID   string = "shuttle.contextID"
	TelemetryCommand     string = "shuttle.command"
	TelemetryCommandArgs string = "shuttle.command.args"
)

func WithContextID(ctx context.Context) context.Context {
	if context_id := os.Getenv("SHUTTLE_CONTEXT_ID"); context_id != "" {
		return context.WithValue(ctx, TelemetryContextID, context_id)
	}

	return context.WithValue(ctx, TelemetryContextID, uuid.New().String())
}

func WithPhase(phase string) TelemetryOption {
	return WithText("phase", phase)
}

func WithLabel(label string) TelemetryOption {
	return WithText("label", label)
}

func WithText(key, value string) TelemetryOption {
	return func(properties map[string]string) {
		properties[key] = value
	}
}

func WithGoInfo() TelemetryOption {
	return func(properties map[string]string) {
		gi, err := goInfo.GetInfo()
		if err != nil {

			properties["system.goinfo.error"] = err.Error()
			return
		}
		if gi.OS != "" {
			properties["system.os"] = gi.OS
		}
		if gi.Kernel != "" {
			properties["system.kernel"] = gi.Kernel
		}
		if gi.Core != "" {
			properties["system.core"] = gi.Core
		}
		if gi.Platform != "" {
			properties["system.platform"] = gi.Platform
		}
		if gi.Hostname != "" {
			properties["system.hostname"] = gi.Hostname
		}
		if gi.CPUs != 0 {
			properties["system.cpus"] = fmt.Sprintf("%d", gi.CPUs)
		}
		if gi.GoOS != "" {
			properties["system.goos"] = gi.GoOS
		}
	}
}

// TODO: rename
func mergeMaps(ctx context.Context, properties map[string]string) map[string]string {
	if runID, ok := ctx.Value(TelemetryContextID).(string); ok && runID != "" {
		properties[TelemetryContextID] = runID
	}

	if val, ok := ctx.Value(TelemetryCommand).(string); ok && val != "" {
		properties[TelemetryCommand] = val
	}

	if val, ok := ctx.Value(TelemetryCommandArgs).(string); ok && val != "" {
		properties[TelemetryCommandArgs] = val
	}

	return properties
}

func copyHostMap(original map[string]string, flowProperties map[string]string) map[string]string {
	properties := make(map[string]string, len(flowProperties)+len(original))
	for k, v := range original {
		properties[k] = v
	}

	for k, v := range flowProperties {
		properties[k] = v
	}

	return properties
}
