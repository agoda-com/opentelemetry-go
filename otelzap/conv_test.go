package otelzap

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	testFieldKey = "test-123"
	testNow      = time.Now()
)

func TestOTeLAttributeMapping(t *testing.T) {
	tests := []struct {
		Input    zapcore.Field
		Expected []attribute.KeyValue
	}{
		{Input: zap.Bool(testFieldKey, true), Expected: []attribute.KeyValue{attribute.Bool(testFieldKey, true)}},
		{Input: zap.Float64(testFieldKey, 123.123), Expected: []attribute.KeyValue{attribute.Float64(testFieldKey, 123.123)}},
		{Input: zap.Int(testFieldKey, 123), Expected: []attribute.KeyValue{attribute.Int64(testFieldKey, 123)}},
		{Input: zap.String(testFieldKey, "hello"), Expected: []attribute.KeyValue{attribute.String(testFieldKey, "hello")}},
		{Input: zap.ByteString(testFieldKey, []byte("hello")), Expected: []attribute.KeyValue{attribute.String(testFieldKey, "hello")}},
		{Input: zap.Binary(testFieldKey, []byte{1, 0, 0, 1}), Expected: []attribute.KeyValue{attribute.String(testFieldKey, "AQAAAQ==")}},
		{Input: zap.Duration(testFieldKey, time.Minute), Expected: []attribute.KeyValue{attribute.Float64(testFieldKey, time.Minute.Seconds())}},
		{Input: zap.Time(testFieldKey, testNow), Expected: []attribute.KeyValue{attribute.Int64(testFieldKey, testNow.Unix())}},
		{Input: zap.Stringer(testFieldKey, bytes.NewBuffer([]byte("hello"))), Expected: []attribute.KeyValue{attribute.String(testFieldKey, "hello")}},
		{Input: zap.Error(errors.New("world")), Expected: []attribute.KeyValue{semconv.ExceptionMessage("world")}},
		{Input: zap.Skip(), Expected: []attribute.KeyValue{}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test.Input), func(t *testing.T) {
			output := otelAttribute(test.Input)
			assert.ElementsMatch(t, test.Expected, output)
		})
	}
}
