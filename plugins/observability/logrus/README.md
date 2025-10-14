# Logrus Plugin

The Logrus plugin provides operators for logging observables using the Logrus logging library.

## Installation

```bash
go get github.com/samber/ro/plugins/observability/logrus
```

## Operators

### Log

Logs all observable events (next, error, complete) at the specified log level.

```go
import (
    "github.com/samber/ro"
    rologrus "github.com/samber/ro/plugins/observability/logrus"
    "github.com/sirupsen/logrus"
)

// Create a logger
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)

observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    rologrus.Log[int](logger, logrus.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// level=info msg="ro.Next: 1"
// level=info msg="ro.Next: 2"
// level=info msg="ro.Next: 3"
// level=info msg="ro.Next: 4"
// level=info msg="ro.Next: 5"
// level=info msg="ro.Complete"
```

### LogWithNotification

Logs observable events with structured logging, including the value in the log entry.

```go
observable := ro.Pipe1(
    ro.Just("Hello", "World", "Golang"),
    rologrus.LogWithNotification[string](logger, logrus.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// level=info msg="ro.Next" value=Hello
// level=info msg="ro.Next" value=World
// level=info msg="ro.Next" value=Golang
// level=info msg="ro.Complete"
```

### FatalOnError

Logs fatal errors when the observable emits an error.

```go
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rologrus.FatalOnError[int](logger),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value int) {
            // Handle successful value
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

## Log Levels

The plugin supports all Logrus log levels:

- `logrus.PanicLevel`
- `logrus.FatalLevel`
- `logrus.ErrorLevel`
- `logrus.WarnLevel`
- `logrus.InfoLevel`
- `logrus.DebugLevel`
- `logrus.TraceLevel`

```go
// Log at debug level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rologrus.Log[int](logger, logrus.DebugLevel),
)
```

## Context Support

All operators support context for structured logging:

```go
ctx := context.WithValue(context.Background(), "request_id", "12345")

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rologrus.Log[int](logger, logrus.InfoLevel),
)

subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[int]())
defer subscription.Unsubscribe()
```

## Structured Data

The plugin works well with structured data types:

```go
type User struct {
    Name string
    Age  int
}

observable := ro.Pipe1(
    ro.Just(
        User{Name: "Alice", Age: 30},
        User{Name: "Bob", Age: 25},
        User{Name: "Charlie", Age: 35},
    ),
    rologrus.LogWithNotification[User](logger, logrus.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[User]())
defer subscription.Unsubscribe()

// Output:
// level=info msg="ro.Next" value="{Alice 30}"
// level=info msg="ro.Next" value="{Bob 25}"
// level=info msg="ro.Next" value="{Charlie 35}"
// level=info msg="ro.Complete"
```

## Logger Configuration

You can configure the logger with various options:

```go
// Create a logger with custom configuration
logger := logrus.New()
logger.SetFormatter(&logrus.JSONFormatter{})
logger.SetOutput(os.Stdout)
logger.SetLevel(logrus.InfoLevel)

// Add fields to the logger
logger = logger.WithField("service", "my-app")
logger = logger.WithField("version", "1.0.0")

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rologrus.Log[int](logger, logrus.InfoLevel),
)
```

## Error Handling

The plugin provides different error handling strategies:

### Log Errors

```go
// Log errors at a specific level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rologrus.Log[int](logger, logrus.ErrorLevel),
)
```

### Fatal on Error

```go
// Log fatal errors and terminate
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rologrus.FatalOnError[int](logger),
)
```

## Real-world Example

Here's a practical example that logs API requests:

```go
import (
    "context"
    "github.com/samber/ro"
    rologrus "github.com/samber/ro/plugins/observability/logrus"
    "github.com/sirupsen/logrus"
)

// Create a logger for API requests
logger := logrus.New()
logger.SetFormatter(&logrus.JSONFormatter{})
logger.SetLevel(logrus.InfoLevel)

// Process API requests with logging
pipeline := ro.Pipe2(
    // Source: API requests
    ro.Just("GET /users", "POST /users", "GET /users/123"),
    // Log each request
    rologrus.LogWithNotification[string](logger, logrus.InfoLevel),
)

subscription := pipeline.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses Logrus's efficient logging mechanisms
- Context propagation adds minimal overhead
- Structured logging with fields is optimized
- Consider log level configuration for production
- Use appropriate log formatters for your environment
- The plugin doesn't block the observable stream
- Logging is done asynchronously to avoid performance impact 