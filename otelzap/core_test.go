package otelzap

import (
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestLevelEnabler_ChangeLevelAfterCreation(t *testing.T) {
	atomicLevel := zap.NewAtomicLevelAt(zap.WarnLevel)
	loggerProvider := sdk.NewLoggerProvider()
	core := NewOtelCore(loggerProvider, WithLevelEnabler(atomicLevel))

	assert.False(t, core.Enabled(zap.InfoLevel))

	atomicLevel.SetLevel(zap.InfoLevel)

	assert.True(t, core.Enabled(zap.InfoLevel))
}
