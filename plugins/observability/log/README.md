# Log Observability Plugin

The log observability plugin provides operators for logging reactive stream notifications using Go's standard `log` package.

## Installation

```bash
go get github.com/samber/ro/plugins/observability/log
```

## Operators

### Log

Logs all notifications (Next, Error, Complete) from an observable stream.

```go
import (
    "github.com/samber/ro"
    rolog "github.com/samber/ro/plugins/observability/log"
)

observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    rolog.Log[int](),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 ro.Next: 1
// 2024/01/01 12:00:00 ro.Next: 2
// 2024/01/01 12:00:00 ro.Next: 3
// 2024/01/01 12:00:00 ro.Next: 4
// 2024/01/01 12:00:00 ro.Next: 5
// 2024/01/01 12:00:00 ro.Complete
```

### LogWithPrefix

Logs all notifications with a custom prefix for better identification.

```go
observable := ro.Pipe1(
    ro.Just("hello", "world", "golang"),
    rolog.LogWithPrefix[string]("[MyApp]"),
)

subscription := observable.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 [MyApp] ro.Next: hello
// 2024/01/01 12:00:00 [MyApp] ro.Next: world
// 2024/01/01 12:00:00 [MyApp] ro.Next: golang
// 2024/01/01 12:00:00 [MyApp] ro.Complete
```

### FatalOnError

Logs errors and calls `log.Fatal` when an error occurs, terminating the program.

```go
observable := ro.Pipe1(
    ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
        observer.Next(1)
        observer.Next(2)
        observer.Error(errors.New("critical error"))
        return nil
    }),
    rolog.FatalOnError[int](),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 ro.Error: critical error
```

### FatalOnErrorWithPrefix

Logs errors with a custom prefix and calls `log.Fatal` when an error occurs.

```go
observable := ro.Pipe1(
    ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
        observer.Next(1)
        observer.Error(errors.New("database connection failed"))
        return nil
    }),
    rolog.FatalOnErrorWithPrefix[int]("[Database]"),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 [Database] ro.Error: database connection failed
```

## Advanced Usage

### Logging in Complex Pipelines

```go
import (
    "fmt"
    "github.com/samber/ro"
    rolog "github.com/samber/ro/plugins/observability/log"
)

// Use logging in a complex pipeline
observable := ro.Pipe3(
    ro.Just(1, 2, 3, 4, 5),
    ro.Filter(func(n int) bool { return n%2 == 0 }), // Keep even numbers
    rolog.LogWithPrefix[int]("[Filter]"),
    ro.Map(func(n int) string { return fmt.Sprintf("Even: %d", n) }),
)

subscription := observable.Subscribe(ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 [Filter] ro.Next: 2
// 2024/01/01 12:00:00 [Filter] ro.Next: 4
// 2024/01/01 12:00:00 [Filter] ro.Complete
```

### Context-Aware Logging

```go
import (
    "context"
    "github.com/samber/ro"
    rolog "github.com/samber/ro/plugins/observability/log"
)

// Log with context-aware operations
ctx := context.Background()

observable := ro.Pipe1(
    ro.Just("context", "aware", "logging"),
    rolog.LogWithPrefix[string]("[Context]"),
)

subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 [Context] ro.Next: context
// 2024/01/01 12:00:00 [Context] ro.Next: aware
// 2024/01/01 12:00:00 [Context] ro.Next: logging
// 2024/01/01 12:00:00 [Context] ro.Complete
```

### Error Logging

```go
import (
    "errors"
    "github.com/samber/ro"
    rolog "github.com/samber/ro/plugins/observability/log"
)

// Log including error notifications
observable := ro.Pipe1(
    ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
        observer.Next(1)
        observer.Next(2)
        observer.Error(errors.New("something went wrong"))
        observer.Next(3) // This won't be emitted due to error
        return nil
    }),
    rolog.Log[int](),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 ro.Next: 1
// 2024/01/01 12:00:00 ro.Next: 2
// 2024/01/01 12:00:00 ro.Error: something went wrong
```

## Real-world Example

Here's a practical example that demonstrates logging in a data processing pipeline:

```go
import (
    "context"
    "errors"
    "fmt"
    "github.com/samber/ro"
    rolog "github.com/samber/ro/plugins/observability/log"
)

// Simulate a data processing pipeline with logging
pipeline := ro.Pipe4(
    // Generate data
    ro.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
    rolog.LogWithPrefix[int]("[Input]"),
    
    // Filter even numbers
    ro.Filter(func(n int) bool { return n%2 == 0 }),
    rolog.LogWithPrefix[int]("[Filter]"),
    
    // Transform data
    ro.Map(func(n int) string { return fmt.Sprintf("Processed: %d", n) }),
    rolog.LogWithPrefix[string]("[Transform]"),
    
    // Simulate some errors
    ro.Map(func(s string) string {
        if s == "Processed: 6" {
            panic("Simulated error for 6")
        }
        return s
    }),
    rolog.LogWithPrefix[string]("[Process]"),
)

// Add error handling
ctx := context.Background()

subscription := pipeline.SubscribeWithContext(ctx, ro.NoopObserver[string]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:00 [Input] ro.Next: 1
// 2024/01/01 12:00:00 [Input] ro.Next: 2
// 2024/01/01 12:00:00 [Filter] ro.Next: 2
// 2024/01/01 12:00:00 [Transform] ro.Next: Processed: 2
// 2024/01/01 12:00:00 [Process] ro.Next: Processed: 2
// 2024/01/01 12:00:00 [Input] ro.Next: 3
// 2024/01/01 12:00:00 [Input] ro.Next: 4
// 2024/01/01 12:00:00 [Filter] ro.Next: 4
// 2024/01/01 12:00:00 [Transform] ro.Next: Processed: 4
// 2024/01/01 12:00:00 [Process] ro.Next: Processed: 4
// ... (continues with logging for each step)
```

## Best Practices

1. **Use Prefixes**: Always use `LogWithPrefix` to identify which part of your pipeline is logging.

2. **Error Handling**: Use `FatalOnError` sparingly, only for truly critical errors that should terminate the application.

3. **Performance**: Logging operators add overhead, so use them judiciously in high-performance scenarios.

4. **Context**: Use context-aware logging when working with cancellable operations.

5. **Structured Logging**: Consider using other observability plugins (zap, logrus, etc.) for more structured logging capabilities.

## Integration with Other Observability Plugins

The log plugin is part of the observability suite. You can combine it with other observability plugins:

- **Zap**: For structured logging with high performance
- **Logrus**: For structured logging with hooks
- **Slog**: For Go's standard structured logging
- **Sentry**: For error tracking and monitoring

## Configuration

The log plugin uses Go's standard `log` package, so you can configure it using standard log functions:

```go
import (
    "log"
    "os"
)

// Set log output to a file
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
log.SetOutput(file)

// Set log flags
log.SetFlags(log.LstdFlags | log.Lshortfile)

// Use the logging operators
observable := ro.Pipe1(
    ro.Just(1, 2, 3),
    rolog.Log[int](),
)
```
