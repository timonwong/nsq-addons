package nsq_zap

// LogType configures the types for nsq logger.
type LogType int

const (
	TypeUndefined LogType = iota
	TypeProducer
	TypeConsumer
)

// An Option configures a ZapNsqLogger.
type Option interface {
	apply(*ZapNsqLogger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*ZapNsqLogger)

func (f optionFunc) apply(log *ZapNsqLogger) {
	f(log)
}

// WithLogMode configures LogType for a ZapNsqLogger.
func WithLogMode(logType LogType) Option {
	return optionFunc(func(logger *ZapNsqLogger) {
		logger.logType = logType
	})
}
