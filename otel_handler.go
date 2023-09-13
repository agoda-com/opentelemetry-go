/*
Copyright Agoda Services Co.,Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package otelslog

import (
	"context"
	otel "github.com/agoda-com/opentelemetry-logs-go/logs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

const (
	instrumentationName = "github.com/agoda-com/otelslog"
)

// OtelHandler is a Handler that writes Records to OTLP
type OtelHandler struct {
	otelHandler
}

type otelHandler struct {
	logger otel.Logger
}

var instrumentationScope = instrumentation.Scope{
	Name:      instrumentationName,
	Version:   Version(),
	SchemaURL: semconv.SchemaURL,
}

func (o otelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (o otelHandler) Handle(ctx context.Context, record slog.Record) error {

	spanContext := trace.SpanFromContext(ctx).SpanContext()
	var traceID *trace.TraceID = nil
	var spanID *trace.SpanID = nil
	var traceFlags *trace.TraceFlags = nil
	if spanContext.IsValid() {
		tid := spanContext.TraceID()
		sid := spanContext.SpanID()
		tf := spanContext.TraceFlags()
		traceID = &tid
		spanID = &sid
		traceFlags = &tf
	}
	levelString := record.Level.String()
	severity := otel.SeverityNumber(int(record.Level.Level()) + 9)

	var attributes []attribute.KeyValue

	record.Attrs(func(attr slog.Attr) bool {
		attributes = append(attributes, otelAttribute(attr)...)
		return true
	})

	lrc := otel.LogRecordConfig{
		Timestamp:            &record.Time,
		ObservedTimestamp:    record.Time,
		TraceId:              traceID,
		SpanId:               spanID,
		TraceFlags:           traceFlags,
		SeverityText:         &levelString,
		SeverityNumber:       &severity,
		Body:                 &record.Message,
		Resource:             nil,
		InstrumentationScope: &instrumentationScope,
		Attributes:           &attributes,
	}

	r := otel.NewLogRecord(lrc)
	o.logger.Emit(r)
	return nil
}

func (o otelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//TODO implement me
	panic("implement me")
}

func (o otelHandler) WithGroup(name string) slog.Handler {
	//TODO implement me
	panic("implement me")
}

// compilation time verification
var _ slog.Handler = &otelHandler{}

// HandlerOptions are options for a OtelHandler.
// A zero HandlerOptions consists entirely of default values.
type HandlerOptions struct {
}

// NewOtelHandler creates a OtelHandler that writes to otlp,
// using the given options.
// If opts is nil, the default options are used.
func NewOtelHandler(loggerProvider otel.LoggerProvider, opts *HandlerOptions) *OtelHandler {
	logger := loggerProvider.Logger(
		instrumentationScope.Name,
		otel.WithInstrumentationVersion(instrumentationScope.Version),
	)
	return &OtelHandler{
		otelHandler{
			logger: logger,
		},
	}
}
