# nsq_zap

## How to use

```go
import "github.com/timonwong/nsq-addons/logging/zap"

func main() {
    // ...
    l := nsq_zap.NewZapNsqLogger(logger)
    // For producer
    producerLogger := l.WithOptions(nsq_zap.WithLogMode(nsq_zap.TypeProducer))
    producer.SetLogger(producerLogger, nsq.LogLevelInfo)
    // For consumer
    consumerLogger := l.WithOptions(nsq_zap.WithLogMode(nsq_zap.TypeConsumer))
    consumer.SetLogger(consumerLogger, nsq.LogLevelInfo)
}
```

