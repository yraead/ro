# Signal Plugin

The signal plugin provides operators for handling operating system signals using Go's `os/signal` package.

## Installation

```bash
go get github.com/samber/ro/plugins/signal
```

## Operators

### NewSignalCatcher

Creates an observable that catches and emits operating system signals.

```go
import (
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

// Catch all incoming signals
observable := rosignal.NewSignalCatcher()

subscription := observable.Subscribe(
    ro.NewObserver(
        func(signal os.Signal) {
            // Handle incoming signal
            switch signal {
            case syscall.SIGINT:
                // Handle Ctrl+C
            case syscall.SIGTERM:
                // Handle termination signal
            case syscall.SIGHUP:
                // Handle hangup signal
            }
        },
        func(err error) {
            // Handle error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Signal Types

The plugin can catch various operating system signals:

- **SIGINT**: Interrupt signal (Ctrl+C)
- **SIGTERM**: Termination signal
- **SIGHUP**: Hangup signal
- **SIGUSR1**: User-defined signal 1
- **SIGUSR2**: User-defined signal 2
- **SIGQUIT**: Quit signal
- **SIGKILL**: Kill signal (cannot be caught)

## Specific Signal Catching

You can catch specific signals only:

```go
import (
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

// Catch only specific signals
observable := rosignal.NewSignalCatcher(
    syscall.SIGINT,  // Ctrl+C
    syscall.SIGTERM, // Termination
    syscall.SIGHUP,  // Hangup
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(signal os.Signal) {
            // Handle specific signals
        },
        func(err error) {
            // Handle error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Signal Filtering

You can filter signals based on your requirements:

```go
import (
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

// Catch all signals but filter for specific ones
observable := ro.Pipe1(
    rosignal.NewSignalCatcher(),
    ro.Filter(func(signal os.Signal) bool {
        // Only process SIGINT and SIGTERM
        return signal == syscall.SIGINT || signal == syscall.SIGTERM
    }),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(signal os.Signal) {
            // Handle filtered signals
        },
        func(err error) {
            // Handle error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Signal Transformation

You can transform signals into other formats:

```go
import (
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

// Transform signals to string descriptions
observable := ro.Pipe1(
    rosignal.NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP),
    ro.Map(func(signal os.Signal) string {
        switch signal {
        case syscall.SIGINT:
            return "Interrupt signal received"
        case syscall.SIGTERM:
            return "Termination signal received"
        case syscall.SIGHUP:
            return "Hangup signal received"
        default:
            return "Unknown signal received"
        }
    }),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

## Error Handling

The plugin handles signal catching errors gracefully:

```go
observable := rosignal.NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(signal os.Signal) {
            // Handle successful signal reception
        },
        func(err error) {
            // Handle signal catching error
            // This could be due to:
            // - Insufficient permissions
            // - Signal not supported on platform
            // - Other system limitations
        },
        func() {
            // Handle completion (when signal catching stops)
        },
    ),
)
defer subscription.Unsubscribe()
```

## Context Support

You can use context for cancellation and timeout:

```go
import (
    "context"
    "time"
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

observable := rosignal.NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM)

subscription := observable.SubscribeWithContext(
    ctx,
    ro.NewObserverWithContext(
        func(ctx context.Context, signal os.Signal) {
            // Handle signal with context
        },
        func(ctx context.Context, err error) {
            // Handle error with context
        },
        func(ctx context.Context) {
            // Handle completion with context
        },
    ),
)
defer subscription.Unsubscribe()
```

## Graceful Shutdown Example

Here's a practical example for graceful shutdown:

```go
import (
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

// Catch signals for graceful shutdown
observable := ro.Pipe1(
    rosignal.NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM),
    ro.Map(func(signal os.Signal) string {
        // Transform signal to shutdown action
        switch signal {
        case syscall.SIGINT:
            return "Graceful shutdown initiated by user"
        case syscall.SIGTERM:
            return "Graceful shutdown initiated by system"
        default:
            return "Unknown shutdown signal"
        }
    }),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(action string) {
            // Perform graceful shutdown
            // e.g., close connections, save state, etc.
        },
        func(err error) {
            // Handle error during shutdown
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Real-world Example

Here's a practical example that handles different types of signals:

```go
import (
    "os"
    "syscall"
    "github.com/samber/ro"
    rosignal "github.com/samber/ro/plugins/signal"
)

// Handle different signal types
pipeline := ro.Pipe2(
    // Catch common signals
    rosignal.NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP),
    // Transform to actions
    ro.Map(func(signal os.Signal) string {
        switch signal {
        case syscall.SIGINT:
            return "shutdown"
        case syscall.SIGTERM:
            return "shutdown"
        case syscall.SIGHUP:
            return "reload"
        default:
            return "unknown"
        }
    }),
)

subscription := pipeline.Subscribe(
    ro.NewObserver(
        func(action string) {
            switch action {
            case "shutdown":
                // Perform shutdown operations
            case "reload":
                // Perform reload operations
            }
        },
        func(err error) {
            // Handle error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses Go's standard `os/signal` package for signal handling
- Signal catching is asynchronous and non-blocking
- Only catch signals that your application needs to handle
- Use context cancellation to properly clean up signal handlers
- Consider platform-specific signal behavior differences
- The plugin automatically handles signal registration and cleanup
- Context cancellation properly stops signal monitoring 