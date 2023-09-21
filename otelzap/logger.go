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
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a thin wrapper for zap.Logger that adds Ctx method.
type Logger struct {
	*zap.Logger
}

const contextKey = "context"

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: l.Logger.Sugar(),
	}
}

func (l *Logger) Ctx(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return l.With(zap.Reflect(contextKey, span.SpanContext()))
	}
	return l
}

func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

type SugaredLogger struct {
	*zap.SugaredLogger
}

func (l *SugaredLogger) Ctx(ctx context.Context) *SugaredLogger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return &SugaredLogger{
			SugaredLogger: l.With(zap.Reflect(contextKey, span.SpanContext())),
		}
	}
	return &SugaredLogger{
		SugaredLogger: l.SugaredLogger,
	}
}
