package nsq_zap

import (
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	producerRE = regexp.MustCompile(`^(?P<lvl>[A-Z]{3})[ ]+(?P<id>\d+) (?P<msg>.+)$`)
	consumerRE = regexp.MustCompile(`^(?P<lvl>[A-Z]{3})[ ]+(?P<id>\d+) \[(?P<topic>[.a-zA-Z0-9_-]+(#ephemeral)?)/(?P<channel>[.a-zA-Z0-9_-]+(#ephemeral)?)\] (?P<msg>.+)$`)
)

// ZapNsqLogger wraps zap logger for nsq producer.
type ZapNsqLogger struct {
	logger  *zap.Logger
	logType LogType
}

// NewZapNsqLogger returns a new ZapNSQLogger.
func NewZapNsqLogger(logger *zap.Logger, opts ...Option) *ZapNsqLogger {
	l := &ZapNsqLogger{logger: logger}
	for _, opt := range opts {
		opt.apply(l)
	}
	return l
}

func (l *ZapNsqLogger) clone() *ZapNsqLogger {
	copy := *l
	return &copy
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (l *ZapNsqLogger) WithOptions(opts ...Option) *ZapNsqLogger {
	c := l.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// Output implements the nsq.logger interface.
func (l *ZapNsqLogger) Output(calldepth int, s string) error {
	var lvl, msg string
	var fields []zapcore.Field

	switch l.logType {
	case TypeProducer:
		matches := producerRE.FindStringSubmatch(s)
		if len(matches) != 4 {
			return nil
		}

		lvl = matches[1]
		msg = matches[3]
		producerID, _ := strconv.Atoi(matches[2])
		fields = []zapcore.Field{
			zap.Int("producerID", producerID),
		}
	case TypeConsumer:
		matches := consumerRE.FindStringSubmatch(s)
		if len(matches) != 8 {
			return nil
		}

		lvl = matches[1]
		msg = matches[7]
		consumerID, _ := strconv.Atoi(matches[2])
		fields = []zapcore.Field{
			zap.Int("consumerID", consumerID),
			zap.String("topic", matches[3]),
			zap.String("channel", matches[5]),
		}
	case TypeUndefined:
		xs := strings.SplitN(s, " ", 2)
		if len(xs) != 2 {
			return nil
		}
		lvl, msg = xs[0], xs[1]
	default:
		return nil
	}

	logger := l.logger.WithOptions(zap.AddCallerSkip(calldepth)).With(fields...)
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
