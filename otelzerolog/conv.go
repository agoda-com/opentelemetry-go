package otelzerolog

import (
	"fmt"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

func otelAttribute(key string, value interface{}) []attribute.KeyValue {
	switch value := value.(type) {
	case bool:
		return []attribute.KeyValue{attribute.Bool(key, value)}
		// Number information is lost when we're converting to byte to interface{}, let's recover it
	case float64:
		if _, frac := math.Modf(value); frac == 0.0 {
			return []attribute.KeyValue{attribute.Int64(key, int64(value))}
		} else {
			return []attribute.KeyValue{attribute.Float64(key, value)}
		}
	case string:
		return []attribute.KeyValue{attribute.String(key, value)}
	case []interface{}:
		var result []attribute.KeyValue
		for _, v := range value {
			// recursively call otelAttribute to handle nested arrays
			result = append(result, otelAttribute(key, v)...)
		}
		return result
	}
	// Default case
	return []attribute.KeyValue{attribute.String(key, fmt.Sprintf("%v", value))}
}
