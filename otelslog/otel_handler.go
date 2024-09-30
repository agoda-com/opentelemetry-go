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
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"sync"
)

const (
	instrumentationName = "github.com/agoda-com/otelslog"
)

// OtelHandler is a Handler that writes Records to OTLP
type OtelHandler struct {
	otelHandler
}

// HandlerOptions are options for a OtelHandler.
// A zero HandlerOptions consists entirely of default values.
type HandlerOptions struct {
	Level slog.Leveler
	AddBaggage bool
}

type otelHandler struct {
	logger      otel.Logger
	opts        HandlerOptions
	groupPrefix string
	attrs       []slog.Attr
	mu          *sync.Mutex
	w           io.Writer
}

// compilation time verification handler implement interface
var _ slog.Handler = &otelHandler{}

var instrumentationScope = instrumentation.Scope{
	Name:      instrumentationName,
	Version:   Version(),
	SchemaURL: semconv.SchemaURL,
}

func (o otelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if o.opts.Level != nil {
		minLevel = o.opts.Level.Level()
	}
	return level >= minLevel
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

	if o.opts.AddBaggage {
		b := baggage.FromContext(ctx)
		// Iterate over baggage items and add them to log attributes
		for _, i := range b.Members() {
			attributes = append(attributes, attribute.String(i.Key(), i.Value()))
		}
	}
	for _, attr := range o.attrs {
		attributes = append(attributes, otelAttribute(attr)...)
	}


	record.Attrs(func(attr slog.Attr) bool {
		attributes = append(attributes, otelAttribute(withGroupPrefix(o.groupPrefix, attr))...)
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

func withGroupPrefix(groupPrefix string, attr slog.Attr) slog.Attr {
	if groupPrefix != "" {
		attr.Key = groupPrefix + attr.Key
	}
	return attr
}

func (o otelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	for i, attr := range attrs {
		attrs[i] = withGroupPrefix(o.groupPrefix, attr)
	}

	return &otelHandler{
		logger:      o.logger,
		opts:        o.opts,
		groupPrefix: o.groupPrefix,
		attrs:       append(o.attrs, attrs...),
	}
}

func (o otelHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return o
	}
	prefix := name + "."
	if o.groupPrefix != "" {
		prefix = o.groupPrefix + prefix
	}

	return &otelHandler{
		logger:      o.logger,
		opts:        o.opts,
		attrs:       o.attrs,
		groupPrefix: prefix,
	}
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
		otelHandler: otelHandler{
			logger: logger,
			opts:   *opts,
		},
	}
}
