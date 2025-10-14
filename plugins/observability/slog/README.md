# Slog Plugin

The Slog plugin provides operators for structured logging observables using Go's built-in `log/slog` package.

## Installation

```bash
go get github.com/samber/ro/plugins/observability/slog
```

## Operators

### Log

Logs all observable events (next, error, complete) using slog at the specified level.

```go
import (
    "log/slog"
    "os"
    "github.com/samber/ro"
    roslog "github.com/samber/ro/plugins/observability/slog"
)

// Create a slog logger
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    roslog.Log[int](*logger, slog.LevelInfo),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 1"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 2"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 3"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 4"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 5"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Complete"
```

### LogWithNotification

Logs observable events with structured attributes, including the value as a log attribute.

```go
observable := ro.Pipe1(
    ro.Just("Hello", "World", "Golang"),
    roslog.LogWithNotification[string](*logger, slog.LevelInfo),
)

subscription := observable.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next" value=Hello
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next" value=World
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next" value=Golang
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Complete"
```

## Slog Levels

The plugin supports all slog log levels:

- `slog.LevelDebug`
- `slog.LevelInfo`
- `slog.LevelWarn`
- `slog.LevelError`

```go
// Log at debug level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.Log[int](*logger, slog.LevelDebug),
)
```

## Context Support

All operators support context for structured logging:

```go
ctx := context.WithValue(context.Background(), "request_id", "12345")

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.Log[int](*logger, slog.LevelInfo),
)

subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 1"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 2"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next: 3"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Complete"
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
    roslog.LogWithNotification[User](*logger, slog.LevelInfo),
)

subscription := observable.Subscribe(ro.NoopObserver[User]())
defer subscription.Unsubscribe()

// Output:
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next" value="{Alice 30}"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next" value="{Bob 25}"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Next" value="{Charlie 35}"
// time=2024-01-01T12:00:00.000Z level=INFO msg="ro.Complete"
```

## Logger Configuration

You can configure the slog logger with various handlers:

### Text Handler

```go
// Create a text handler
handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})
logger := slog.New(handler)

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.Log[int](*logger, slog.LevelInfo),
)
```

### JSON Handler

```go
// Create a JSON handler
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})
logger := slog.New(handler)

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.Log[int](*logger, slog.LevelInfo),
)
```

### Custom Handler

```go
// Create a custom handler
handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
    AddSource: true,
})
logger := slog.New(handler)

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.Log[int](*logger, slog.LevelInfo),
)
```

## Error Handling

The plugin provides different error handling strategies:

### Log Errors

```go
// Log errors at a specific level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.Log[int](*logger, slog.LevelError),
)
```

### Structured Error Logging

```go
// Log errors with structured attributes
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    roslog.LogWithNotification[int](*logger, slog.LevelError),
)
```

## Real-world Example

Here's a practical example that logs API requests:

```go
import (
    "context"
    "log/slog"
    "os"
    "github.com/samber/ro"
    roslog "github.com/samber/ro/plugins/observability/slog"
)

// Create a logger for API requests
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})
logger := slog.New(handler)

// Process API requests with logging
pipeline := ro.Pipe2(
    // Source: API requests
    ro.Just("GET /users", "POST /users", "GET /users/123"),
    // Log each request
    roslog.LogWithNotification[string](*logger, slog.LevelInfo),
)

subscription := pipeline.Subscribe(
    ro.NewObserver(
        func(request string) {
            // Process the request
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

- The plugin uses Go's efficient slog logging mechanisms
- Context propagation adds minimal overhead
- Structured logging with attributes is optimized
- Consider log level configuration for production
- Use appropriate handlers for your environment
- The plugin doesn't block the observable stream
- Logging is done asynchronously to avoid performance impact
- Slog provides efficient JSON and text formatting 