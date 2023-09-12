package otelslog

import (
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
)

func otelAttribute(attr slog.Attr) attribute.KeyValue {
	switch attr.Value.Kind() {
	case slog.KindBool:
		return attribute.Bool(attr.Key, attr.Value.Bool())
	//case slog.KindDuration: ???
	case slog.KindFloat64:
		return attribute.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return attribute.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return attribute.String(attr.Key, attr.Value.String())
	//case slog.KindTime: ???
	case slog.KindUint64:
		return attribute.Int64(attr.Key, int64(attr.Value.Uint64()))
		//case slog.KindGroup: ???
		//case slog.KindLogValuer: ???
	}
	return attribute.String(attr.Key, attr.Value.String())
}
