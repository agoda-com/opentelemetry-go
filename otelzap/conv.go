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

package otelzap

import (
	otel "github.com/agoda-com/opentelemetry-logs-go/logs"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/zap/zapcore"
	"math"
)

// otelLevel zap level to otlp level converter
func otelLevel(level zapcore.Level) otel.SeverityNumber {
	switch level {
	case zapcore.DebugLevel:
		return otel.DEBUG
	case zapcore.InfoLevel:
		return otel.INFO
	case zapcore.WarnLevel:
		return otel.WARN
	case zapcore.ErrorLevel:
		return otel.ERROR
	case zapcore.DPanicLevel:
		return otel.ERROR
	case zapcore.PanicLevel:
		return otel.ERROR
	case zapcore.FatalLevel:
		return otel.FATAL
	}
	return otel.TRACE
}

// otelAttribute convert zap Field into OpenTelemetry Attribute
func otelAttribute(f zapcore.Field) []attribute.KeyValue {
	switch f.Type {
	case zapcore.UnknownType:
		return []attribute.KeyValue{attribute.String(f.Key, f.String)}
	case zapcore.BoolType:
		return []attribute.KeyValue{attribute.Bool(f.Key, f.Integer == 1)}
	case zapcore.Float64Type:
		return []attribute.KeyValue{attribute.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))}
	case zapcore.Float32Type:
		return []attribute.KeyValue{attribute.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))}
	case zapcore.Int64Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, f.Integer)}
	case zapcore.Int32Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, f.Integer)}
	case zapcore.Int16Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, f.Integer)}
	case zapcore.Int8Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, f.Integer)}
	case zapcore.StringType:
		return []attribute.KeyValue{attribute.String(f.Key, f.String)}
	case zapcore.Uint64Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, int64(uint64(f.Integer)))}
	case zapcore.Uint32Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, int64(uint64(f.Integer)))}
	case zapcore.Uint16Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, int64(uint64(f.Integer)))}
	case zapcore.Uint8Type:
		return []attribute.KeyValue{attribute.Int64(f.Key, int64(uint64(f.Integer)))}
	case zapcore.ErrorType:
		err := f.Interface.(error)
		if err != nil {
			return []attribute.KeyValue{semconv.ExceptionMessage(err.Error())}
		}
		return []attribute.KeyValue{}
	case zapcore.SkipType:
		return []attribute.KeyValue{}
	}
	// unhandled types will be treated as string
	return []attribute.KeyValue{attribute.String(f.Key, f.String)}
}
