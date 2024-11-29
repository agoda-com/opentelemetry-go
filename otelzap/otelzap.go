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
	"go.uber.org/zap/zapcore"
)

// NewOtelCore creates new OpenTelemetry Core to export logs in OTLP format
func NewOtelCore(loggerProvider otel.LoggerProvider, opts ...Option) zapcore.Core {
	logger := loggerProvider.Logger(
		instrumentationScope.Name,
		otel.WithInstrumentationVersion(instrumentationScope.Version),
	)

	c := &otlpCore{
		logger:       logger,
		levelEnabler: zapcore.InfoLevel,
	}
	for _, apply := range opts {
		apply(c)
	}

	return c
}

// Option is a function that applies an option to an OpenTelemetry Core
type Option func(c *otlpCore)

// WithLevel sets the minimum level for the OpenTelemetry Core log to be exported
func WithLevel(levelEnabler zapcore.LevelEnabler) Option {
	return Option(func(c *otlpCore) {
		c.levelEnabler = levelEnabler
	})
}
