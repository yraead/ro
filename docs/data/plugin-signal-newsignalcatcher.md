---
name: NewSignalCatcher
slug: newsignalcatcher
sourceRef: plugins/signal/source.go#L30
type: plugin
category: signal
signatures:
  - "func NewSignalCatcher(signals ...os.Signal)"
playUrl: ""
variantHelpers:
  - plugin#signal#newsignalcatcher
similarHelpers: []
position: 0
---

Catches system signals and emits them as observable values.

```go
import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

obs := rosignal.NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM)

sub := obs.Subscribe(ro.NewObserver[os.Signal](
    func(sig os.Signal) {
        fmt.Printf("Received signal: %v\n", sig)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Signal stream completed")
    },
))
defer sub.Unsubscribe()

// When user presses Ctrl+C:
// Received signal: interrupt
// Completed
```

If no signals are provided, all incoming signals will be relayed. The observable completes when the signal channel is closed.