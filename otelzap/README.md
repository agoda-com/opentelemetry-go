# otelzap

Zap logger with OpenTelemetry support. This logger will export LogRecord's in OTLP format.

## Quick start

Configure open-telemetry provider. See [example here](../README.md)

Then configure zap logger with otel core:

```go
package main

import (
	"context"
	"github.com/agoda-com/opentelemetry-go/otelzap"
	"go.uber.org/zap"
)

func main() {

	// configure logger provider
	loggerProvider :=  ...

	// create new  logger with opentelemetry zap core and set it globally
	logger := zap.New(otelzap.NewOtelCore(loggerProvider))
	zap.ReplaceGlobals(logger)
}

// call function with opentelemetry context provided
func doSomething(ctx context.Context) {
	// send log with opentelemetry context
	otelzap.Ctx(ctx).Info("My message with trace context")
}

```
