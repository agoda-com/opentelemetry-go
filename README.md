# otelslog
log/slog handler for OTel

## Quick start

[Export env variable](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/#otel_exporter_otlp_endpoint)  `OTEL_EXPORTER_OTLP_ENDPOINT=https://localhost:4318`
to your OTLP collector

To start with otelslog first you need to configure opentelemetry-logs-go exporters. See full example bellow:

```go
package main

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/agoda-com/otelslog"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log/slog"
	"os"
)

// configure common attributes for all logs
func newResource() *resource.Resource {
	hostName, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otelslog-example"),
		semconv.ServiceVersion("1.0.0"),
		semconv.HostName(hostName),
	)
}

func main() {

	ctx := context.Background()

	// configure opentelemetry logger provider
	logExporter, _ := otlplogs.NewExporter(ctx)
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(logExporter),
		sdk.WithResource(newResource()),
	)
	// gracefully shutdown logger to flush accumulated signals before program finish
	defer loggerProvider.Shutdown(ctx)

	otelLogger := slog.New(otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{}))
	
	//configure default logger
	slog.SetDefault(otelLogger)

	doSomething(ctx)
}

func doSomething(ctx context.Context) {
	slog.InfoContext(ctx, "hello", slog.String("myKey", "myValue"))
}
```