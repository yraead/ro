# File System Notify Plugin

The file system notify plugin provides operators for monitoring file system events using the `fsnotify` package.

## Installation

```bash
go get github.com/samber/ro/plugins/fsnotify
```

## Operators

### NewFSListener

Creates an observable that monitors file system events for specified paths.

```go
import (
    "os"
    "github.com/fsnotify/fsnotify"
    "github.com/samber/ro"
    rofsnotify "github.com/samber/ro/plugins/fsnotify"
)

// Monitor a single directory
tempDir := os.TempDir()
observable := rofsnotify.NewFSListener(tempDir)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(event fsnotify.Event) {
            // Handle file system event
            switch event.Op {
            case fsnotify.Create:
                // File was created
            case fsnotify.Write:
                // File was written to
            case fsnotify.Remove:
                // File was removed
            case fsnotify.Rename:
                // File was renamed
            case fsnotify.Chmod:
                // File permissions changed
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

## Event Types

The plugin monitors the following file system events:

- **Create**: A new file or directory was created
- **Write**: A file was written to
- **Remove**: A file or directory was removed
- **Rename**: A file or directory was renamed
- **Chmod**: File permissions were changed

## Multiple Path Monitoring

You can monitor multiple directories simultaneously:

```go
paths := []string{
    "/path/to/dir1",
    "/path/to/dir2",
    "/path/to/dir3",
}

observable := rofsnotify.NewFSListener(paths...)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(event fsnotify.Event) {
            // Handle events from any of the monitored paths
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

## Event Filtering

You can filter events based on file extensions or event types:

```go
import (
    "path/filepath"
    "github.com/fsnotify/fsnotify"
    "github.com/samber/ro"
    rofsnotify "github.com/samber/ro/plugins/fsnotify"
)

// Filter by file extension
observable := ro.Pipe1(
    rofsnotify.NewFSListener("/path/to/monitor"),
    ro.Filter(func(event fsnotify.Event) bool {
        // Only process .txt files
        return filepath.Ext(event.Name) == ".txt"
    }),
)

// Filter by event type
observable := ro.Pipe1(
    rofsnotify.NewFSListener("/path/to/monitor"),
    ro.Filter(func(event fsnotify.Event) bool {
        // Only process create and write events
        return event.Op&(fsnotify.Create|fsnotify.Write) != 0
    }),
)
```

## Event Throttling

To avoid processing too many rapid successive events, you can throttle the events:

```go
import (
    "time"
    "github.com/samber/ro"
    rofsnotify "github.com/samber/ro/plugins/fsnotify"
)

observable := ro.Pipe1(
    rofsnotify.NewFSListener("/path/to/monitor"),
    ro.ThrottleTime[fsnotify.Event](100 * time.Millisecond),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(event fsnotify.Event) {
            // Handle throttled file system event
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

## Error Handling

The plugin handles various file system errors gracefully:

```go
observable := rofsnotify.NewFSListener("/path/to/monitor")

subscription := observable.Subscribe(
    ro.NewObserver(
        func(event fsnotify.Event) {
            // Handle successful file system event
        },
        func(err error) {
            // Handle file system monitoring error
            // This could be due to:
            // - Insufficient permissions
            // - Directory not existing
            // - File system limitations
            // - Other file system issues
        },
        func() {
            // Handle completion (when monitoring stops)
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
    "github.com/samber/ro"
    rofsnotify "github.com/samber/ro/plugins/fsnotify"
)

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

observable := rofsnotify.NewFSListener("/path/to/monitor")

subscription := observable.SubscribeWithContext(
    ctx,
    ro.NewObserverWithContext(
        func(ctx context.Context, event fsnotify.Event) {
            // Handle file system event with context
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

## Event Transformation

You can transform file system events into other formats:

```go
observable := ro.Pipe1(
    rofsnotify.NewFSListener("/path/to/monitor"),
    ro.Map(func(event fsnotify.Event) string {
        // Transform event to string representation
        return event.Name + " - " + event.Op.String()
    }),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

## Real-world Example

Here's a practical example that monitors a log directory and processes new log files:

```go
import (
    "path/filepath"
    "strings"
    "github.com/fsnotify/fsnotify"
    "github.com/samber/ro"
    rofsnotify "github.com/samber/ro/plugins/fsnotify"
)

// Monitor log directory for new log files
pipeline := ro.Pipe3(
    // Monitor the logs directory
    rofsnotify.NewFSListener("/var/log"),
    // Filter for create events on .log files
    ro.Filter(func(event fsnotify.Event) bool {
        return event.Op == fsnotify.Create && 
               filepath.Ext(event.Name) == ".log"
    }),
    // Transform to log file path
    ro.Map(func(event fsnotify.Event) string {
        return event.Name
    }),
)

subscription := pipeline.Subscribe(
    ro.NewObserver(
        func(logFile string) {
            // Process new log file
            // e.g., read and analyze the log file
        },
        func(err error) {
            // Handle monitoring error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses the `fsnotify` package for efficient file system monitoring
- Events are emitted asynchronously to avoid blocking
- Monitor only necessary directories to reduce system load
- Use filtering to process only relevant events
- Consider throttling for high-frequency event sources
- The plugin automatically handles file system limitations and errors
- Context cancellation properly cleans up monitoring resources 