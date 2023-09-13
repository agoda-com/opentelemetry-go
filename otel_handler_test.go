package otelslog

import (
	"bytes"
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	"log/slog"

	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"os"
	"testing"
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
	slog.InfoContext(ctx, "hello slog", slog.String("myKey", "myValue"))
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

	otelLogger := slog.New(NewOtelHandler(loggerProvider, &HandlerOptions{}))
	slog.SetDefault(otelLogger)

	doSomething(ctx)

	loggerProvider.Shutdown(ctx)

	_ = buf.String()

	// TODO: require opentelemetry-logs-go updated with fixed writer
	//assert.Equal(t, "{host.name=CLX4NV72V6, service.name=otelslog-example, service.version=1.0.0, myKey=myValue}", actual)
}
