# Zerolog Plugin

The Zerolog plugin provides operators for structured logging observables using the Zerolog logging library.

## Installation

```bash
go get github.com/samber/ro/plugins/observability/zerolog
```

## Operators

### Log

Logs all observable events (next, error, complete) using Zerolog at the specified level.

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/samber/ro"
    rozerolog "github.com/samber/ro/plugins/observability/zerolog"
)

// Create a Zerolog logger
logger := log.With().Str("service", "my-app").Logger()

observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    rozerolog.Log[int](&logger, zerolog.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// {"level":"info","message":"ro.Next: 1"}
// {"level":"info","message":"ro.Next: 2"}
// {"level":"info","message":"ro.Next: 3"}
// {"level":"info","message":"ro.Next: 4"}
// {"level":"info","message":"ro.Next: 5"}
// {"level":"info","message":"ro.Complete"}
```

### LogWithNotification

Logs observable events with structured fields, including the value as a log field.

```go
observable := ro.Pipe1(
    ro.Just("Hello", "World", "Golang"),
    rozerolog.LogWithNotification[string](&logger, zerolog.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// {"level":"info","value":"Hello","message":"ro.Next"}
// {"level":"info","value":"World","message":"ro.Next"}
// {"level":"info","value":"Golang","message":"ro.Next"}
// {"level":"info","message":"ro.Complete"}
```

### FatalOnError

Logs fatal errors when the observable emits an error.

```go
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.FatalOnError[int](&logger),
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

## Zerolog Levels

The plugin supports all Zerolog log levels:

- `zerolog.DebugLevel`
- `zerolog.InfoLevel`
- `zerolog.WarnLevel`
- `zerolog.ErrorLevel`
- `zerolog.FatalLevel`
- `zerolog.PanicLevel`
- `zerolog.Disabled`

```go
// Log at debug level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.DebugLevel),
)
```

## Context Support

All operators support context for structured logging:

```go
ctx := context.WithValue(context.Background(), "request_id", "12345")

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.InfoLevel),
)

subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// {"level":"info","message":"ro.Next: 1"}
// {"level":"info","message":"ro.Next: 2"}
// {"level":"info","message":"ro.Next: 3"}
// {"level":"info","message":"ro.Complete"}
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
    rozerolog.LogWithNotification[User](&logger, zerolog.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[User]())
defer subscription.Unsubscribe()

// Output:
// {"level":"info","value":"{Alice 30}","message":"ro.Next"}
// {"level":"info","value":"{Bob 25}","message":"ro.Next"}
// {"level":"info","value":"{Charlie 35}","message":"ro.Next"}
// {"level":"info","message":"ro.Complete"}
```

## Logger Configuration

You can configure the Zerolog logger with various options:

### Global Logger

```go
// Use the global logger
logger := log.With().Str("service", "my-app").Logger()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.InfoLevel),
)
```

### Console Logger

```go
// Create a console logger
logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.InfoLevel),
)
```

### JSON Logger

```go
// Create a JSON logger
logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.InfoLevel),
)
```

### Custom Logger

```go
// Create a custom logger with fields
logger := zerolog.New(os.Stdout).
    With().
    Str("service", "my-app").
    Str("version", "1.0.0").
    Timestamp().
    Logger()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.InfoLevel),
)
```

## Error Handling

The plugin provides different error handling strategies:

### Log Errors

```go
// Log errors at a specific level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.Log[int](&logger, zerolog.ErrorLevel),
)
```

### Fatal on Error

```go
// Log fatal errors and terminate
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.FatalOnError[int](&logger),
)
```

### Structured Error Logging

```go
// Log errors with structured fields
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozerolog.LogWithNotification[int](&logger, zerolog.ErrorLevel),
)
```

## Real-world Example

Here's a practical example that logs API requests:

```go
import (
    "context"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/samber/ro"
    rozerolog "github.com/samber/ro/plugins/observability/zerolog"
)

// Create a logger for API requests
logger := log.With().
    Str("service", "api-gateway").
    Str("environment", "production").
    Logger()

// Process API requests with logging
pipeline := ro.Pipe2(
    // Source: API requests
    ro.Just("GET /users", "POST /users", "GET /users/123"),
    // Log each request
    rozerolog.LogWithNotification[string](&logger, zerolog.InfoLevel),
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

- The plugin uses Zerolog's efficient logging mechanisms
- Context propagation adds minimal overhead
- Structured logging with fields is optimized
- Consider log level configuration for production
- Use appropriate logger configuration for your environment
- The plugin doesn't block the observable stream
- Logging is done asynchronously to avoid performance impact
- Zerolog provides efficient JSON encoding
- Zero-allocation logging for high-performance applications 