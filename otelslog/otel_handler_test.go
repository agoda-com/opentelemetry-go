package otelslog

import (
	"bytes"
	"context"
	"go.opentelemetry.io/otel/baggage"
	"log/slog"
	"os"
	"testing"

	"github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	"github.com/stretchr/testify/assert"

	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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

func doSomething(ctx context.Context) {
	slog.InfoContext(ctx, "hello slog", slog.String("myKey", "myValue"),
		slog.Group("myGroup", slog.String("groupKey", "groupValue")))
}

func TestNewOtelHandler(t *testing.T) {
	ctx := context.Background()

	var buf bytes.Buffer

	// configure opentelemetry logger provider
	logExporter, _ := stdoutlogs.NewExporter(stdoutlogs.WithWriter(&buf))
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(logExporter),
		sdk.WithResource(newResource()),
	)

	handler := NewOtelHandler(loggerProvider, &HandlerOptions{
		Level:      slog.LevelInfo,
		AddBaggage: true,
	}).
		WithAttrs([]slog.Attr{slog.String("first", "value1")}).
		WithGroup("group1").
		WithAttrs([]slog.Attr{slog.String("second", "value2")}).
		WithGroup("group2")

	otelLogger := slog.New(handler)
	slog.SetDefault(otelLogger)

	member, _ := baggage.NewMember("baggage.key", "true")
	bag, _ := baggage.New(member)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	doSomething(ctx)

	loggerProvider.Shutdown(ctx)

	actual := buf.String()

	assert.Contains(t, actual, "INFO hello slog [scopeInfo: github.com/agoda-com/otelslog:0.2.0] {host.name=")
	assert.Contains(t, actual, "service.name=otelslog-example, service.version=1.0.0, baggage.key=true, first=value1, group1.second=value2, group1.group2.myKey=myValue, group1.group2.myGroup.groupKey=groupValue}")
}
