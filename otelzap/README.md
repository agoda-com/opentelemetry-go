# otelzap

Zap logger with OpenTelemetry support. This logger will export LogRecord's in OTLP format.

## Quick start

[Export env variable](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/#otel_exporter_otlp_endpoint)  `OTEL_EXPORTER_OTLP_ENDPOINT=https://localhost:4318`
to your OTLP collector

```go
package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	semconv2 "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"github.com/agoda-com/otelzap"
	otellogs "github.com/agoda-com/opentelemetry-logs-go"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	"os"
)

// configure common attributes for all logs 
func newResource() *resource.Resource {
	hostName, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otelzap-example"),
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

	// set opentelemetry logger provider globally 
	otellogs.SetLoggerProvider(loggerProvider)

	// create new  logger with opentelemetry zap core and set it globally
	logger := zap.New(otelzap.NewOtelCore(loggerProvider))
	zap.ReplaceGlobals(logger)

	// now your application ready to produce logs to opentelemetry collector
	doSomething()

}

func doSomething() {
	// start new span
	// see official trace documentation https://github.com/open-telemetry/opentelemetry-go
	tracer := otel.Tracer("my-tracer")
	spanCtx, span := tracer.Start(context.Background(), "My Span")
	defer func() {
		span.End()
	}()

	// send log with opentelemetry context
	otelzap.Ctx(spanCtx).Info("My message with trace context")
}

```
