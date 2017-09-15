package nsq_zap

import (
	"regexp"

	"go.uber.org/zap"
)

var (
	spaces = regexp.MustCompile(`\s+`)
)

// ZapNSQLogger wraps zap logger for nsq producer.
type ZapNSQLogger struct {
	logger *zap.Logger
}

// NewZapNsqLogger returns a new ZapNSQLogger.
func NewZapNsqLogger(logger *zap.Logger) *ZapNSQLogger {
	return &ZapNsqLogger{logger: logger}
}

// Output implements the nsq.logger interface.
func (l *ZapNsqLogger) Output(calldepth int, s string) error {
	logger := l.logger

	fields := spaces.Split(s, 2)
	if len(fields) != 2 {
		logger.Warn(s)
	}

	lvl, msg := fields[0], fields[1]
	switch lvl {
	case "INF":
		logger.Info(msg)
	case "WRN":
		logger.Warn(msg)
	case "ERR":
		logger.Error(msg)
	default:
		logger.Debug(msg)
	}

	return nil
}
