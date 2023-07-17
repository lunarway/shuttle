package telemetry

import (
	"context"
	"fmt"

	"github.com/matishsiao/goInfo"
)

const (
	telemetryContextID   string = "shuttle.contextID"
	telemetryRunID       string = "shuttle.runID"
	TelemetryCommand     string = "shuttle.command"
	TelemetryCommandArgs string = "shuttle.command.args"
)

func WithPhase(phase string) TelemetryOption {
	return WithEntry("phase", phase)
}

func WithLabel(label string) TelemetryOption {
	return WithEntry("label", label)
}

func WithEntry(key, value string) TelemetryOption {
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

func includeContext(ctx context.Context, properties map[string]string) map[string]string {
	getFromContext(ctx, telemetryContextID, properties)
	getFromContext(ctx, TelemetryCommand, properties)
	getFromContext(ctx, TelemetryCommandArgs, properties)

	return properties
}

func getFromContext(ctx context.Context, key string, properties map[string]string) {
	if val, ok := ctx.Value(TelemetryCommandArgs).(string); ok && val != "" {
		properties[key] = val
	}
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
