# otelslog

log/slog handler for OTel

## Quick start

Configure open-telemetry provider first. See [example here](../README.md)

Then configure slog logger with otelslog handler:

```go
package main

import (
	"context"
	"github.com/agoda-com/opentelemetry-go/otelslog"
	"log/slog"
)

func main() {
	// configure logger provider
	loggerProvider :=  ...

	otelLogger := slog.New(otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{}))

	//configure default logger
	slog.SetDefault(otelLogger)

	doSomething(ctx)
}

// call function with opentelemetry context provided
func doSomething(ctx context.Context) {
	slog.InfoContext(ctx, "hello", slog.String("myKey", "myValue"))
}
```