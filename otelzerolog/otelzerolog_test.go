package otelzerolog

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// configure common attributes for all logs
func newResource() *resource.Resource {
	hostName, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otelzerolog-example"),
		semconv.ServiceVersion("1.0.0"),
		semconv.HostName(hostName),
	)
}

func TestZerologHook(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer

	// configure opentelemetry logger provider
	logExporter, _ := stdoutlogs.NewExporter(stdoutlogs.WithWriter(&buf))
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithSyncer(logExporter), // use syncer to make sure all logs are flushed before test ends
		sdk.WithResource(newResource()),
	)

	hook := NewHook(loggerProvider)
	log := log.Hook(hook)
	log.Info().Ctx(ctx).Str("key", "value").Msg("hello zerolog")

	_ = loggerProvider.Shutdown(ctx)

	actual := buf.String()
	assert.Contains(t, actual, "INFO")                                                    // ensure th log level
	assert.Contains(t, actual, "hello zerolog")                                           // ensure the message
	assert.Contains(t, actual, "[scopeInfo: otelzerolog:0.0.1]")                          // ensure the scope info
	assert.Contains(t, actual, "service.name=otelzerolog-example, service.version=1.0.0") // ensure the resource attributes
	assert.Contains(t, actual, "level=info")                                              // ensure the severity attributes
	assert.Contains(t, actual, "key=value")                                               // ensure the log fields
}

func TestZerologHook_ValidSpan(t *testing.T) {
	var buf bytes.Buffer

	// configure opentelemetry logger provider
	logExporter, _ := stdoutlogs.NewExporter(stdoutlogs.WithWriter(&buf))
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithSyncer(logExporter), // use syncer to make sure all logs are flushed before test ends
		sdk.WithResource(newResource()),
	)

	// create a span
	tracerProvider := oteltrace.NewTracerProvider(oteltrace.WithResource(newResource()))
	tracer := tracerProvider.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test")
	defer span.End()

	hook := NewHook(loggerProvider)
	log := log.Hook(hook)
	log.Warn().Ctx(ctx).Str("key", "value").Msg("hello zerolog")

	actual := buf.String()
	assert.Contains(t, actual, span.SpanContext().SpanID().String())  // ensure the spanID is logged
	assert.Contains(t, actual, span.SpanContext().TraceID().String()) // ensure the traceID is logged

	log.Error().Ctx(ctx).Str("key", "value").Discard().Msg("this should not be logged")
	_ = loggerProvider.Shutdown(ctx)
}
