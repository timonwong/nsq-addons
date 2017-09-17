package nsq_zap

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapNsqLogger_WithOptions(t *testing.T) {
	core, _ := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())

	l1 := NewZapNsqLogger(zapLogger, WithLogMode(TypeUndefined))
	l2 := l1.WithOptions(WithLogMode(TypeConsumer))
	l3 := l2.WithOptions(WithLogMode(TypeProducer))

	assert.Equal(t, TypeUndefined, l1.logType)
	assert.Equal(t, TypeConsumer, l2.logType)
	assert.Equal(t, TypeProducer, l3.logType)
}

func TestZapNsqLogger_LogTypeUndefined(t *testing.T) {
	core, logEntries := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())

	l := NewZapNsqLogger(zapLogger, WithLogMode(TypeUndefined))
	l.Output(2, "WRN Test message")

	allEntries := logEntries.TakeAll()
	assert.Len(t, allEntries, 1, "There should be only one log entry.")
	assert.Equal(t, zapcore.WarnLevel, allEntries[0].Level)
	assert.Equal(t, "Test message", allEntries[0].Message)
}

func TestZapNsqLogger_LogTypeProducer(t *testing.T) {
	core, logEntries := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())

	l := NewZapNsqLogger(zapLogger, WithLogMode(TypeProducer))
	producerID := rand.Intn(9999) + 1
	lvls := []string{"DBG", "INF", "WRN", "ERR"}
	for _, lvl := range lvls {
		logProducer(l, lvl, producerID, "(%s) connecting to nsqd", "localhost")
	}

	allEntries := logEntries.TakeAll()
	assert.Len(t, allEntries, len(lvls))

	zapLevels := []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel}
	for idx, logEntry := range allEntries {
		assert.Equal(t, zapLevels[idx], logEntry.Level, "Log level not match.")
		assert.Regexp(t,
			`.+/nsq-addons/logging/zap/nsqlogger_test.go:\d+$`,
			logEntry.Caller,
			"Expected to find package name and file name in output.")
		assert.Equal(t,
			"(localhost) connecting to nsqd",
			logEntry.Message,
			"log message should match.")

		ctxMap := logEntry.ContextMap()
		assert.Equal(t, int64(producerID), ctxMap["producerID"], "producerID not match")
	}
}

func TestZapNsqLogger_LogTypeConsumer(t *testing.T) {
	core, logEntries := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())

	l := NewZapNsqLogger(zapLogger, WithLogMode(TypeConsumer))
	consumerID := rand.Intn(9999) + 1
	lvls := []string{"DBG", "INF", "WRN", "ERR"}
	for _, lvl := range lvls {
		logConsumer(l, lvl, consumerID, "test-topic", "test-channel", "querying nsqlookupd %s", "localhost:4161")
	}

	allEntries := logEntries.TakeAll()
	assert.Len(t, allEntries, len(lvls))

	zapLevels := []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel}
	for idx, logEntry := range allEntries {
		assert.Equal(t, zapLevels[idx], logEntry.Level, "Log level should match.")
		assert.Regexp(t,
			`.+/nsq-addons/logging/zap/nsqlogger_test.go:\d+$`,
			logEntry.Caller,
			"Expected to find package name and file name in output.")
		assert.Equal(t,
			"querying nsqlookupd localhost:4161",
			logEntry.Message,
			"log message should match.")

		ctxMap := logEntry.ContextMap()
		assert.Equal(t, int64(consumerID), ctxMap["consumerID"], "consumerID should match.")
		assert.Equal(t, "test-topic", ctxMap["topic"], "topic should match.")
		assert.Equal(t, "test-channel", ctxMap["channel"], "channel should match.")
	}
}

type logger interface {
	Output(calldepth int, s string) error
}

func logProducer(logger logger, lvl string, producerID int, line string, args ...interface{}) {
	logger.Output(2, fmt.Sprintf("%-4s %3d %s", lvl, producerID, fmt.Sprintf(line, args...)))
}

func logConsumer(logger logger, lvl string, consumerID int, topic, channel, line string, args ...interface{}) {
	logger.Output(2, fmt.Sprintf("%-4s %3d [%s/%s] %s", lvl, consumerID, topic, channel, fmt.Sprintf(line, args...)))
}
