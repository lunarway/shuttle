package telemetry

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/matishsiao/goInfo"
)

const telemetryContextID = "shuttle.contextID"

func WithContextID(ctx context.Context) context.Context {
	return context.WithValue(ctx, telemetryContextID, uuid.New().String())
}

func WithPhase(phase string) TelemetryOption {
	return WithText("phase", phase)
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

func mergeMaps(ctx context.Context, properties map[string]string) map[string]string {
	if runID, ok := ctx.Value(telemetryContextID).(string); ok && runID != "" {
		properties[telemetryContextID] = runID
	}

	return properties
}

func copyHostMap(original map[string]string, target map[string]string) {
	for k, v := range original {
		target[k] = v
	}
}
