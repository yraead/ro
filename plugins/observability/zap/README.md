# Zap Plugin

The Zap plugin provides operators for structured logging observables using the Uber Zap logging library.

## Installation

```bash
go get github.com/samber/ro/plugins/observability/zap
```

## Operators

### Log

Logs all observable events (next, error, complete) using Zap at the specified level.

```go
import (
    "github.com/samber/ro"
    rozap "github.com/samber/ro/plugins/observability/zap"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Create a Zap logger
logger, _ := zap.NewProduction()

observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    rozap.Log[int](logger, zapcore.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

logger.Sync()

// Output:
// {"level":"info","msg":"ro.Next: 1"}
// {"level":"info","msg":"ro.Next: 2"}
// {"level":"info","msg":"ro.Next: 3"}
// {"level":"info","msg":"ro.Next: 4"}
// {"level":"info","msg":"ro.Next: 5"}
// {"level":"info","msg":"ro.Complete"}
```

### LogWithNotification

Logs observable events with structured fields, including the value as a log field.

```go
observable := ro.Pipe1(
    ro.Just("Hello", "World", "Golang"),
    rozap.LogWithNotification[string](logger, zapcore.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()

logger.Sync()

// Output:
// {"level":"info","msg":"ro.Next","value":"Hello"}
// {"level":"info","msg":"ro.Next","value":"World"}
// {"level":"info","msg":"ro.Next","value":"Golang"}
// {"level":"info","msg":"ro.Complete"}
```

### FatalOnError

Logs fatal errors when the observable emits an error.

```go
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.FatalOnError[int](logger),
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

logger.Sync()
```

## Zap Levels

The plugin supports all Zap log levels:

- `zapcore.DebugLevel`
- `zapcore.InfoLevel`
- `zapcore.WarnLevel`
- `zapcore.ErrorLevel`
- `zapcore.DPanicLevel`
- `zapcore.PanicLevel`
- `zapcore.FatalLevel`

```go
// Log at debug level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.DebugLevel),
)
```

## Context Support

All operators support context for structured logging:

```go
ctx := context.WithValue(context.Background(), "request_id", "12345")

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.InfoLevel),
)

subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[int]())
defer subscription.Unsubscribe()

logger.Sync()

// Output:
// {"level":"info","msg":"ro.Next: 1"}
// {"level":"info","msg":"ro.Next: 2"}
// {"level":"info","msg":"ro.Next: 3"}
// {"level":"info","msg":"ro.Complete"}
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
    rozap.LogWithNotification[User](logger, zapcore.InfoLevel),
)

subscription := observable.Subscribe(ro.NoopObserver[User]())
defer subscription.Unsubscribe()

logger.Sync()

// Output:
// {"level":"info","msg":"ro.Next","value":"{Alice 30}"}
// {"level":"info","msg":"ro.Next","value":"{Bob 25}"}
// {"level":"info","msg":"ro.Next","value":"{Charlie 35}"}
// {"level":"info","msg":"ro.Complete"}
```

## Logger Configuration

You can configure the Zap logger with various options:

### Production Logger

```go
// Create a production logger
logger, _ := zap.NewProduction()
defer logger.Sync()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.InfoLevel),
)
```

### Development Logger

```go
// Create a development logger
logger, _ := zap.NewDevelopment()
defer logger.Sync()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.InfoLevel),
)
```

### Custom Logger

```go
// Create a custom logger
config := zap.NewProductionConfig()
config.OutputPaths = []string{"stdout", "logs/app.log"}
logger, _ := config.Build()
defer logger.Sync()

observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.InfoLevel),
)
```

### Sugared Logger

```go
// Create a sugared logger
logger, _ := zap.NewProduction()
sugar := logger.Sugar()

// Note: The plugin works with the core logger, not the sugared logger
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.InfoLevel),
)
```

## Error Handling

The plugin provides different error handling strategies:

### Log Errors

```go
// Log errors at a specific level
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.Log[int](logger, zapcore.ErrorLevel),
)
```

### Fatal on Error

```go
// Log fatal errors and terminate
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.FatalOnError[int](logger),
)
```

### Structured Error Logging

```go
// Log errors with structured fields
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rozap.LogWithNotification[int](logger, zapcore.ErrorLevel),
)
```

## Real-world Example

Here's a practical example that logs API requests:

```go
import (
    "context"
    "github.com/samber/ro"
    rozap "github.com/samber/ro/plugins/observability/zap"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Create a logger for API requests
logger, _ := zap.NewProduction()
defer logger.Sync()

// Process API requests with logging
pipeline := ro.Pipe2(
    // Source: API requests
    ro.Just("GET /users", "POST /users", "GET /users/123"),
    // Log each request
    rozap.LogWithNotification[string](logger, zapcore.InfoLevel),
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

- The plugin uses Zap's efficient logging mechanisms
- Context propagation adds minimal overhead
- Structured logging with fields is optimized
- Consider log level configuration for production
- Use appropriate logger configuration for your environment
- The plugin doesn't block the observable stream
- Logging is done asynchronously to avoid performance impact
- Zap provides efficient JSON encoding
- Remember to call `logger.Sync()` in production 