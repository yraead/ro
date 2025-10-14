# Samber PSI Plugin

The samber/psi plugin provides a source for monitoring Pressure Stall Information (PSI) using the [samber/go-psi](https://github.com/samber/go-psi) library.

## Installation

```bash
go get github.com/samber/ro/plugins/samber/psi
```

## Sources

### PSINotifier

Monitors system pressure stall information (PSI) at regular intervals.

```go
import (
    "github.com/samber/ro"
    ropsi "github.com/samber/ro/plugins/samber/psi"
    psinotifier "github.com/samber/go-psi"
    "time"
)

observable := ropsi.NewPSINotifier(5 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats psinotifier.PSIStatsResource) {
        fmt.Printf("CPU Pressure: %.2f%%, Memory Pressure: %.2f%%, IO Pressure: %.2f%%\n",
            stats.CPU.Some.Avg10, stats.Memory.Some.Avg10, stats.IO.Some.Avg10)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))

// Output:
// CPU Pressure: 0.50%, Memory Pressure: 2.30%, IO Pressure: 1.20%
// CPU Pressure: 0.45%, Memory Pressure: 2.25%, IO Pressure: 1.15%
// ...
```

## Advanced Usage

### Filtering High Pressure

```go
import (
    "github.com/samber/ro"
    ropsi "github.com/samber/ro/plugins/samber/psi"
    psinotifier "github.com/samber/go-psi"
    "time"
)

// Monitor PSI and filter for high pressure events
observable := ro.Pipe1(
    ropsi.NewPSINotifier(1 * time.Second),
    ro.Filter(func(stats psinotifier.PSIStatsResource) bool {
        return stats.CPU.Some.Avg10 > 5.0 || // > 5% CPU pressure
               stats.Memory.Some.Avg10 > 10.0 || // > 10% memory pressure
               stats.IO.Some.Avg10 > 5.0 // > 5% IO pressure
    }),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats psinotifier.PSIStatsResource) {
        fmt.Printf("High pressure detected!\n")
        fmt.Printf("CPU: %.2f%%, Memory: %.2f%%, IO: %.2f%%\n",
            stats.CPU.Some.Avg10, stats.Memory.Some.Avg10, stats.IO.Some.Avg10)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Monitoring Specific Pressure Types

```go
import (
    "github.com/samber/ro"
    ropsi "github.com/samber/ro/plugins/samber/psi"
    psinotifier "github.com/samber/go-psi"
    "time"
)

// Monitor only memory pressure
observable := ro.Pipe1(
    ropsi.NewPSINotifier(2 * time.Second),
    ro.Map(func(stats psinotifier.PSIStatsResource) float64 {
        return stats.Memory.Some.Avg10
    }),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(memoryPressure float64) {
        fmt.Printf("Memory pressure: %.2f%%\n", memoryPressure)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Combining with Other Monitoring

```go
import (
    "github.com/samber/ro"
    ropsi "github.com/samber/ro/plugins/samber/psi"
    roproc "github.com/samber/ro/plugins/proc"
    psinotifier "github.com/samber/go-psi"
    "time"
)

// Monitor both PSI and system load
psiObservable := ropsi.NewPSINotifier(5 * time.Second)
loadObservable := roproc.NewLoadAverageWatcher(5 * time.Second)

// Combine them using Merge
combined := ro.Merge(psiObservable, loadObservable)

subscription := combined.Subscribe(ro.NewObserver(
    func(value interface{}) {
        switch v := value.(type) {
        case psinotifier.PSIStatsResource:
            fmt.Printf("PSI - CPU: %.2f%%, Memory: %.2f%%\n",
                v.CPU.Some.Avg10, v.Memory.Some.Avg10)
        case *load.AvgStat:
            fmt.Printf("Load - 1min: %.2f, 5min: %.2f\n",
                v.Load1, v.Load5)
        }
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

## Understanding PSI

Pressure Stall Information (PSI) provides insights into system resource contention:

- **CPU Pressure**: Indicates when processes are waiting for CPU time
- **Memory Pressure**: Indicates when processes are waiting for memory allocation
- **IO Pressure**: Indicates when processes are waiting for IO operations

### PSI Metrics

Each pressure type provides several metrics:

- **Some**: Percentage of time that at least some tasks were stalled
- **Full**: Percentage of time that all non-idle tasks were stalled
- **Avg10**: 10-second average
- **Avg60**: 1-minute average
- **Avg300**: 5-minute average

### Thresholds

Common PSI thresholds for monitoring:

- **Low pressure**: < 5%
- **Medium pressure**: 5-10%
- **High pressure**: > 10%

## Error Handling

The PSI source handles errors gracefully and will emit error notifications if the underlying system calls fail:

```go
observable := ropsi.NewPSINotifier(5 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats psinotifier.PSIStatsResource) {
        // Handle successful PSI data
    },
    func(err error) {
        // Handle errors (e.g., permission denied, PSI not available)
        log.Printf("PSI monitoring error: %v", err)
    },
    func() {
        // Handle completion
    },
))
```

## Dependencies

This plugin requires the [samber/go-psi](https://github.com/samber/go-psi) library:

```bash
go get github.com/samber/go-psi
```

## Performance Considerations

- Choose appropriate intervals for your monitoring needs
- PSI monitoring is relatively lightweight but consider system impact
- Use `ro.ObserveOn` to control the scheduler for PSI operations
- Be mindful of system resources when monitoring frequently

## System Requirements

PSI monitoring requires:

- Linux kernel 4.20+ (for full PSI support)
- Access to `/proc/pressure/` filesystem
- Appropriate permissions to read PSI data 