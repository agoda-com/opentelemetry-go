Open-telemetry extensions for go language
---

## Open-telemetry Loggers

| Logger                     | Version | Minimal go version |
|----------------------------|---------|--------------------|
| [otelslog](otelslog)       | v0.0.1  | 1.21               |
| [otelzap](otelzap)         | v0.1.1  | 1.20               |
| [otelzerolog](otelzerolog) | v0.0.1  | 1.21               |

### Quick start with open-telemetry loggers

Before configure your logger it is required to configure open-telemetry exporter first:

[Export env variable](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/#otel_exporter_otlp_endpoint)  `OTEL_EXPORTER_OTLP_ENDPOINT=https://localhost:4318`
to your OTLP collector

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

	// configure your logger with logger provider
	...

	// to have complete example with correlation between logs and tracing start new span
	// see official trace documentation https://github.com/open-telemetry/opentelemetry-go
	tracer := otel.Tracer("my-tracer")
	spanCtx, span := tracer.Start(context.Background(), "My Span")
	defer func() {
		span.End()
	}()
	
	// now we can call function to execute logs with tracing context
	doSomething(spanCtx)
}

// call function with opentelemetry context provided
func doSomething(ctx context.Context) {
	// log your messages here
	...
}
```
See configuration details for every logger
