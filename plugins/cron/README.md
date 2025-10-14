# Cron Plugin

The cron plugin provides scheduling capabilities for reactive streams using the [gocron](https://github.com/go-co-op/gocron) library.

## Installation

```bash
go get github.com/samber/ro/plugins/cron
```

## Features

- Schedule jobs using cron expressions or duration intervals
- Automatic job execution with reactive stream notifications
- Context-aware cancellation
- Thread-safe job execution

## Usage

### Basic Scheduling

```go
import (
    "time"
    "github.com/samber/ro"
    rocron "github.com/samber/ro/plugins/cron"
    "github.com/go-co-op/gocron/v2"
)

// Schedule a job every 5 seconds
observable := rocron.NewScheduler(
    gocron.DurationJob(5 * time.Second),
)

subscription := observable.Subscribe(ro.PrintObserver[rocron.ScheduleJob]())
defer subscription.Unsubscribe()

// Output: (will emit every 5 seconds)
// Next: {Counter: 0, Time: 2024-01-01 12:00:05 +0000 UTC}
// Next: {Counter: 1, Time: 2024-01-01 12:00:10 +0000 UTC}
// Next: {Counter: 2, Time: 2024-01-01 12:00:15 +0000 UTC}
// ... (continues)
```

### Cron Expressions

```go
// Schedule a job daily at 23:42
observable := rocron.NewScheduler(
    gocron.CronJob("42 23 * * *", false), // Daily at 23:42
)

subscription := observable.Subscribe(ro.PrintObserver[rocron.ScheduleJob]())
defer subscription.Unsubscribe()

// Output: (will emit daily at 23:42)
// Next: {Counter: 0, Time: 2024-01-01 23:42:00 +0000 UTC}
// Next: {Counter: 1, Time: 2024-01-02 23:42:00 +0000 UTC}
// Next: {Counter: 2, Time: 2024-01-03 23:42:00 +0000 UTC}
// ... (continues daily)
```

### Context Cancellation

```go
import (
    "context"
    "time"
    "github.com/samber/ro"
    rocron "github.com/samber/ro/plugins/cron"
    "github.com/go-co-op/gocron/v2"
)

// Create a scheduler with context for cancellation
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

observable := rocron.NewScheduler(
    gocron.DurationJob(1 * time.Second),
)

subscription := observable.SubscribeWithContext(ctx, ro.PrintObserver[rocron.ScheduleJob]())
defer subscription.Unsubscribe()

// Output: (will emit every second for 10 seconds, then complete)
// Next: {Counter: 0, Time: 2024-01-01 12:00:01 +0000 UTC}
// Next: {Counter: 1, Time: 2024-01-01 12:00:02 +0000 UTC}
// ...
// Next: {Counter: 9, Time: 2024-01-01 12:00:10 +0000 UTC}
// Completed
```

### Processing Scheduled Events

```go
import (
    "time"
    "github.com/samber/ro"
    rocron "github.com/samber/ro/plugins/cron"
    "github.com/go-co-op/gocron/v2"
)

// Create a scheduler and process the events
observable := ro.Pipe2(
    rocron.NewScheduler(
        gocron.DurationJob(1 * time.Second),
    ),
    ro.Map(func(job rocron.ScheduleJob) string {
        return "Scheduled job executed at " + job.Time.Format("15:04:05")
    }),
    ro.Take[string](3), // Only take first 3 events
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: Scheduled job executed at 12:00:01
// Next: Scheduled job executed at 12:00:02
// Next: Scheduled job executed at 12:00:03
// Completed
```

## ScheduleJob Structure

Each scheduled job emits a `ScheduleJob` struct:

```go
type ScheduleJob struct {
    Counter int       // Incremental counter starting from 0
    Time    time.Time // Timestamp when the job was executed
}
```

## Common Cron Patterns

| Pattern     | Description              | Example                              |
| ----------- | ------------------------ | ------------------------------------ |
| `* * * * *` | Every minute             | `gocron.CronJob("* * * * *", false)` |
| `0 * * * *` | Every hour               | `gocron.CronJob("0 * * * *", false)` |
| `0 0 * * *` | Every day at midnight    | `gocron.CronJob("0 0 * * *", false)` |
| `0 0 * * 0` | Every Sunday at midnight | `gocron.CronJob("0 0 * * 0", false)` |
| `0 0 1 * *` | First day of each month  | `gocron.CronJob("0 0 1 * *", false)` |
| `0 0 1 1 *` | January 1st at midnight  | `gocron.CronJob("0 0 1 1 *", false)` |

## Duration Scheduling

For simple interval-based scheduling, use `DurationJob`:

```go
// Every 30 seconds
gocron.DurationJob(30 * time.Second)

// Every 5 minutes
gocron.DurationJob(5 * time.Minute)

// Every hour
gocron.DurationJob(1 * time.Hour)

// Every day
gocron.DurationJob(24 * time.Hour)
```

## Real-world Example

Here's a practical example that processes scheduled events and logs them:

```go
import (
    "context"
    "log"
    "time"
    "github.com/samber/ro"
    rocron "github.com/samber/ro/plugins/cron"
    "github.com/go-co-op/gocron/v2"
)

// Create a pipeline that processes scheduled events
pipeline := ro.Pipe3(
    // Schedule a job every 5 seconds
    rocron.NewScheduler(
        gocron.DurationJob(5 * time.Second),
    ),
    // Process the job and create a log message
    ro.Map(func(job rocron.ScheduleJob) string {
        return "Job executed at " + job.Time.Format("2006-01-02 15:04:05") + 
               " (execution #" + fmt.Sprintf("%d", job.Counter+1) + ")"
    }),
    // Log the message
    ro.Tap(func(msg string) {
        log.Println(msg)
    }),
)

// Run for 20 seconds
ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
defer cancel()

subscription := pipeline.SubscribeWithContext(ctx, ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// 2024/01/01 12:00:05 Job executed at 2024-01-01 12:00:05 (execution #1)
// Next: Job executed at 2024-01-01 12:00:05 (execution #1)
// 2024/01/01 12:00:10 Job executed at 2024-01-01 12:00:10 (execution #2)
// Next: Job executed at 2024-01-01 12:00:10 (execution #2)
// 2024/01/01 12:00:15 Job executed at 2024-01-01 12:00:15 (execution #3)
// Next: Job executed at 2024-01-01 12:00:15 (execution #3)
// 2024/01/01 12:00:20 Job executed at 2024-01-01 12:00:20 (execution #4)
// Next: Job executed at 2024-01-01 12:00:20 (execution #4)
// Completed
```

## Error Handling

The scheduler will emit errors if there are issues with job creation or execution:

```go
observable := rocron.NewScheduler(
    gocron.CronJob("invalid cron", false), // Invalid cron expression
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(job rocron.ScheduleJob) {
            // Handle successful job execution
        },
        func(err error) {
            log.Printf("Scheduler error: %v", err)
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
``` 