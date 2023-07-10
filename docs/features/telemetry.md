# Shuttle Telemetry

This feature introduces telemetry to shuttle, it is a bit different than what
you might be used to. This is not for shuttle to send telemetry to us (Lunar),
but rather an option for you (the user/org) to collect telemetry, for your own
needs. This is an opt-in feature, and you will need to build some tooling around
it.

The goal of the telemetry feature is to collect certain information from shuttle
runs which may be useful for analytics purposes. The collected data has
intentionally been anonymized.

To enable telemetry you can either enable local or remote tracing tracing.

- Local tracing just outputs what we would trace to standard out. You can enable
  it so: `export SHUTTLE_LOG_TRACING=true`, and then run a shuttle command
  `shuttle run build`. You should now see extra log statements in your console
  output.
- Remote does a few more things, by itself it doesn't actually upload files, but
  it puts them in a ready to use format (json lines, .jsonl) in your
  `~/.local/share/shuttle/telemetry` folder. Each process will have a unique log
  file in there. To enable this feature simply
  `export SHUTTLE_REMOTE_TRACING=true`, you can set a custom telemetry folder
  directory using `SHUTTLE_REMOTE_LOG_LOCATION` as well.

Finally you can choose to upload the telemetry: simply
`shuttle upload --url https://<your-custom-backend>/publish`. The backend
implementation is intentionally left blank, and for now no reference
implementation is provided.

However, the schema is super straightforward, you just need an endpoint that
accepts the following body:

```json
{
  "app": "string",
  "timestamp": "string", // follows RFC3339
  "properties": {
    "keys": "values"
    // ...
  }
}
```

You can now choose to ingest the logs however, you want, maybe you'd like to
trace it with opentelemetry, or log it to your preferred logging solution, or
put it in a datalake, the choice is yours.

To see what values you can expect, simply
`SHUTTLE_LOG_TRACING=true shuttle run build` and capture some of the output, to
see some of the values you may expect.

Each run will have a `shuttle.contextID` field, this is used to tie a run
together, so that if a build.sh file calls another shuttle command internally
that will be logged under the same contextID as well.
