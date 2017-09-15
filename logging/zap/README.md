# nsq_zap

## How to use

```go
import "github.com/timonwong/nsq-addons/logging/zap"

func main() {
    // ...
    l := nsq_zap.NewZapNsqLogger(logger)
    // For consumer
    consumer.SetLogger(l, nsq.LogLevelInfo)
    // For producer
    producer.SetLogger(l, nsq.LogLevelInfo)
}
```

