# Sentry Plugin

This plugin provides Sentry integration for the ro library, allowing you to log observable events to Sentry for monitoring and debugging.

## Features

- Log observable events (Next, Error, Complete) to Sentry
- Support for different log levels (Debug, Info, Warning, Error)
- Structured logging with extra fields
- Error tracking with stack traces

## Usage

### Basic Logging

```go
import (
    "github.com/getsentry/sentry-go"
    "github.com/samber/ro"
    "github.com/samber/ro/plugins/observability/sentry"
)

// Initialize Sentry hub
hub := sentry.CurrentHub().Clone()
hub.ConfigureScope(func(scope *sentry.Scope) {
    scope.SetTag("component", "observable")
})

// Log all notifications
observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    rosentry.Log[int](hub, sentry.LevelInfo),
)

subscription := observable.Subscribe(ro.NoopObserver[int]())
defer subscription.Unsubscribe()
```

### Structured Logging

```go
// Initialize Sentry hub
hub := sentry.CurrentHub().Clone()

// Log with structured notification data
observable := ro.Pipe1(
    ro.Just("hello", "world", "golang"),
    rosentry.LogWithNotification[string](hub, sentry.LevelDebug),
)
```

### Error Handling

```go
// Initialize Sentry hub
hub := sentry.CurrentHub().Clone()

// Log including error notifications
observable := ro.Pipe1(
    ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
        observer.Next(1)
        observer.Next(2)
        observer.Error(errors.New("something went wrong"))
        return nil
    }),
    rosentry.Log[int](hub, sentry.LevelError),
)
```

## Available Functions

### Operators

- `Log[T](hub *sentry.Hub, level sentry.Level)` - Logs events with simple message formatting
- `LogWithNotification[T](hub *sentry.Hub, level sentry.Level)` - Logs events with structured data

## Log Levels

The plugin supports all Sentry log levels:

- `sentry.LevelDebug` - Debug information
- `sentry.LevelInfo` - General information
- `sentry.LevelWarning` - Warning messages
- `sentry.LevelError` - Error messages
- `sentry.LevelFatal` - Fatal errors

## Dependencies

- `github.com/getsentry/sentry-go` - Sentry Go SDK
- `github.com/samber/ro` - Reactive Observables library 