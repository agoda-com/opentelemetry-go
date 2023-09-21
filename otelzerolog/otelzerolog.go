package otelzerolog

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

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	otel "github.com/agoda-com/opentelemetry-logs-go/logs" // use otel so that when otel is stable, we can just change the import path
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "otelzerolog"
)

var instrumentationScope = instrumentation.Scope{
	Name:      instrumentationName,
	Version:   Version(),
	SchemaURL: semconv.SchemaURL,
}

type Hook struct {
	otel.Logger
}

var _ zerolog.Hook = (*Hook)(nil)

func NewHook(loggerProvider otel.LoggerProvider) *Hook {
	logger := loggerProvider.Logger(
		instrumentationScope.Name,
		otel.WithInstrumentationVersion(instrumentationScope.Version),
	)
	return &Hook{logger}
}

func (h Hook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if !e.Enabled() {
		return
	}

	ctx := e.GetCtx()
	span := trace.SpanFromContext(ctx).SpanContext()

	var spanID trace.SpanID
	var traceID trace.TraceID
	var traceFlags trace.TraceFlags
	if span.IsValid() {
		spanID = span.SpanID()
		traceID = span.TraceID()
		traceFlags = span.TraceFlags()
	}

	now := time.Now()
	severityText := otelLevelText(level)
	severityNumber := otelLevelNumber(level)

	logData := make(map[string]interface{})
	// create a string that appends } to the end of the buf variable you access via reflection
	ev := fmt.Sprintf("%s}", reflect.ValueOf(e).Elem().FieldByName("buf"))
	_ = json.Unmarshal([]byte(ev), &logData)

	// TODO: this is very hacky, but it works for now
	var attributes []attribute.KeyValue
	for k, v := range logData {
		attributes = append(attributes, []attribute.KeyValue{attribute.String(k, fmt.Sprintf("%v", v))}...)
	}

	lrc := otel.LogRecordConfig{
		Timestamp:            &now,
		ObservedTimestamp:    now,
		TraceId:              &traceID,
		SpanId:               &spanID,
		TraceFlags:           &traceFlags,
		SeverityText:         &severityText,
		SeverityNumber:       &severityNumber,
		Body:                 &msg,
		Resource:             nil,
		InstrumentationScope: &instrumentationScope,
		Attributes:           &attributes,
	}

	r := otel.NewLogRecord(lrc)
	h.Emit(r)
}
