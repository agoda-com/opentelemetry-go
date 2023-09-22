# otelzerolog

## Description

This project aims to integrate OpenTelemetry, a set of APIs, libraries, agents, and collector services to capture distributed traces and metrics from your application, with Zerolog, a zero-allocation JSON logger in Go, for monitoring and managing your application's performance.

## Features

- Seamless integration of OpenTelemetry and Zerolog
- Efficient tracing and logging of application activities
- Easy debugging and monitoring

## Prerequisites

- Go (version 1.21)
- Basic understanding of OpenTelemetry and Zerolog

## Installation

```console
go get -u github.com/agoda-com/opentelemetry-go/otelzerolog
```

## Usage

```go
func main() {
  ctx := context.Background()

  // Setup opentelemetry provider
  loggerProvider := sdk.NewLoggerProvider()
  hook := otelzerolog.NewHook(loggerProvider)
  log := log.Hook(hook)

  log.Info().Ctx(ctx).Str("string", "string-value").Msg("Hello OpenTelemetry")
  time.Sleep(10 * time.Second)
}
```
