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
	otel "github.com/agoda-com/opentelemetry-logs-go/logs"
	"github.com/rs/zerolog"
)

func otelLevelNumber(level zerolog.Level) otel.SeverityNumber {
	switch level {
	case zerolog.TraceLevel:
		return otel.TRACE
	case zerolog.DebugLevel:
		return otel.DEBUG
	case zerolog.InfoLevel:
		return otel.INFO
	case zerolog.WarnLevel:
		return otel.WARN
	case zerolog.ErrorLevel:
		return otel.ERROR
	case zerolog.FatalLevel:
		return otel.FATAL
	case zerolog.PanicLevel:
		return otel.FATAL2
	default:
		return otel.INFO
	}
}

func otelLevelText(level zerolog.Level) string {
	switch level {
	case zerolog.TraceLevel:
		return "TRACE"
	case zerolog.DebugLevel:
		return "DEBUG"
	case zerolog.InfoLevel:
		return "INFO"
	case zerolog.WarnLevel:
		return "WARN"
	case zerolog.ErrorLevel:
		return "ERROR"
	case zerolog.FatalLevel:
		return "FATAL"
	case zerolog.PanicLevel:
		return "FATAL"
	default:
		return "INFO"
	}
}
