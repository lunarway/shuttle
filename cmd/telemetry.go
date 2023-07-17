package cmd

import (
	stdcontext "context"
	"strings"

	"github.com/lunarway/shuttle/pkg/telemetry"
)

func trace(
	ctx stdcontext.Context,
	name string,
	args []string,
) (stdcontext.Context, func(options ...telemetry.TelemetryOption), func(err error, options ...telemetry.TelemetryOption), func()) {
	ctx = telemetry.WithContextID(ctx)
	ctx = telemetry.WithRunID(ctx)
	ctx = WithRunTelemetry(ctx, name, args)

	traceInfo := func(options ...telemetry.TelemetryOption) {
		telemetry.Trace(ctx, name, telemetry.WithPhase("start"))
	}
	traceInfo(telemetry.WithPhase("start"))
	traceErr := func(err error, options ...telemetry.TelemetryOption) {
		telemetry.TraceError(ctx, name, err)
	}

	return ctx,
		traceInfo,
		traceErr,
		func() {
			telemetry.Trace(ctx, name, telemetry.WithPhase("end"))
		}
}

func WithRunTelemetry(
	ctx stdcontext.Context,
	commandName string,
	args []string,
) stdcontext.Context {
	ctx = stdcontext.WithValue(ctx, telemetry.TelemetryCommand, commandName)
	if len(args) != 0 {
		// TODO: Make sure we sanitize secrets, somehow
		ctx = stdcontext.WithValue(ctx, telemetry.TelemetryCommandArgs, strings.Join(args[1:], " "))
	}
	return ctx
}
