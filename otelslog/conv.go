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
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
)

func otelAttribute(attr slog.Attr) []attribute.KeyValue {
	switch attr.Value.Kind() {
	case slog.KindBool:
		return []attribute.KeyValue{attribute.Bool(attr.Key, attr.Value.Bool())}
	//case slog.KindDuration: ???
	case slog.KindFloat64:
		return []attribute.KeyValue{attribute.Float64(attr.Key, attr.Value.Float64())}
	case slog.KindInt64:
		return []attribute.KeyValue{attribute.Int64(attr.Key, attr.Value.Int64())}
	case slog.KindString:
		return []attribute.KeyValue{attribute.String(attr.Key, attr.Value.String())}
	//case slog.KindTime: ???
	case slog.KindUint64:
		return []attribute.KeyValue{attribute.Int64(attr.Key, int64(attr.Value.Uint64()))}
	case slog.KindGroup:
		group := attr.Value.Group()
		var result []attribute.KeyValue
		for _, v := range group {
			v.Key = attr.Key + "." + v.Key
			result = append(result, otelAttribute(v)...)
		}
		return result
	}
	return []attribute.KeyValue{attribute.String(attr.Key, attr.Value.String())}
}
