package otelzerolog

import (
	"bytes"
	"context"
	"fmt"
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

type User struct{}

func (u *User) String() string {
	return "I am a user"
}

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
	assert.Contains(t, actual, "INFO")                                                               // ensure th log level
	assert.Contains(t, actual, "hello zerolog")                                                      // ensure the message
	assert.Contains(t, actual, "scopeInfo: github.com/agoda-com/opentelemetry-go/otelzerolog:0.0.1") // ensure the scope info
	assert.Contains(t, actual, "service.name=otelzerolog-example")                                   // ensure the resource attributes
	assert.Contains(t, actual, "service.version=1.0.0")                                              // ensure the resource attributes
	assert.Contains(t, actual, "level=info")                                                         // ensure the severity attributes
	assert.Contains(t, actual, "key=value")                                                          // ensure the log fields
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
	log.Warn().Ctx(ctx).
		Str("key", "value").
		Strs("strs", []string{"1", "2", "3"}).
		Stringer("stringer", &User{}).
		Int("int", 0).
		Int16("i16", 16).
		Int32("i32", 32).
		Int64("i64", 64).
		Dur("dur", 1).
		Uint("u", 0).
		Uint8("u", 0).
		Uint16("u", 0).
		Uint32("u", 0).
		Uint64("u", 0).
		Float32("float32", 32.32).
		Float64("float64", 64.64).
		Bool("bool", true).
		Interface("interface", &User{}).
		Interface("array", []interface{}{"1", 1, "2", 2, "3", 3}).
		Err(fmt.Errorf("new error")).
		Msg("hello zerolog")

	actual := buf.String()
	assert.Contains(t, actual, span.SpanContext().SpanID().String())  // ensure the spanID is logged
	assert.Contains(t, actual, span.SpanContext().TraceID().String()) // ensure the traceID is logged

	log.Error().Ctx(ctx).Str("key", "value").Discard().Msg("this should not be logged")
	_ = loggerProvider.Shutdown(ctx)
}
