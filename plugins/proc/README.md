# Proc Plugin

The proc plugin provides sources for monitoring system resources and processes using the [gopsutil](https://github.com/shirou/gopsutil) library.

## Installation

```bash
go get github.com/samber/ro/plugins/proc
```

## Sources

### Memory Monitoring

#### VirtualMemoryWatcher

Monitors virtual memory usage statistics.

```go
import (
    "github.com/samber/ro"
    roproc "github.com/samber/ro/plugins/proc"
    "time"
)

observable := roproc.NewVirtualMemoryWatcher(5 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats *mem.VirtualMemoryStat) {
        fmt.Printf("Total: %d, Used: %d, Free: %d\n", 
            stats.Total, stats.Used, stats.Free)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))

// Output:
// Total: 16777216000, Used: 8589934592, Free: 8187281408
// Total: 16777216000, Used: 8590065664, Free: 8187149824
// ...
```

#### SwapMemoryWatcher

Monitors swap memory usage statistics.

```go
observable := roproc.NewSwapMemoryWatcher(10 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats *mem.SwapMemoryStat) {
        fmt.Printf("Swap Total: %d, Used: %d, Free: %d\n", 
            stats.Total, stats.Used, stats.Free)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### CPU Monitoring

#### CPUInfoWatcher

Monitors CPU information and statistics.

```go
observable := roproc.NewCPUInfoWatcher(2 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(cpuInfo cpu.InfoStat) {
        fmt.Printf("CPU: %s, Cores: %d, Model: %s\n", 
            cpuInfo.Name, cpuInfo.Cores, cpuInfo.ModelName)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Disk Monitoring

#### DiskUsageWatcher

Monitors disk usage for a specific mountpoint or device.

```go
observable := roproc.NewDiskUsageWatcher(30 * time.Second, "/")

subscription := observable.Subscribe(ro.NewObserver(
    func(usage *disk.UsageStat) {
        fmt.Printf("Path: %s, Total: %d, Used: %d, Free: %d\n", 
            usage.Path, usage.Total, usage.Used, usage.Free)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

#### DiskIOCountersWatcher

Monitors disk I/O counters for specified devices.

```go
observable := roproc.NewDiskIOCountersWatcher(5 * time.Second, "sda", "sdb")

subscription := observable.Subscribe(ro.NewObserver(
    func(counters map[string]disk.IOCountersStat) {
        for device, stats := range counters {
            fmt.Printf("Device: %s, ReadBytes: %d, WriteBytes: %d\n", 
                device, stats.ReadBytes, stats.WriteBytes)
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

### Network Monitoring

#### NetIOCountersWatcher

Monitors network I/O counters.

```go
observable := roproc.NewNetIOCountersWatcher(5 * time.Second, true)

subscription := observable.Subscribe(ro.NewObserver(
    func(counters net.IOCountersStat) {
        fmt.Printf("Interface: %s, BytesSent: %d, BytesRecv: %d\n", 
            counters.Name, counters.BytesSent, counters.BytesRecv)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Host Information

#### HostInfoWatcher

Monitors host information.

```go
observable := roproc.NewHostInfoWatcher(60 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(info *host.InfoStat) {
        fmt.Printf("Hostname: %s, OS: %s, Platform: %s\n", 
            info.Hostname, info.OS, info.Platform)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Load Average

#### LoadAverageWatcher

Monitors system load average.

```go
observable := roproc.NewLoadAverageWatcher(5 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(load *load.AvgStat) {
        fmt.Printf("Load1: %.2f, Load5: %.2f, Load15: %.2f\n", 
            load.Load1, load.Load5, load.Load15)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Sensors

#### SensorsTemperatureWatcher

Monitors temperature sensors.

```go
observable := roproc.NewSensorsTemperatureWatcher(10 * time.Second, false)

subscription := observable.Subscribe(ro.NewObserver(
    func(temp sensors.TemperatureStat) {
        fmt.Printf("Sensor: %s, Temperature: %.2f°C\n", 
            temp.Name, temp.Temperature)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

## Advanced Usage

### Combining Multiple Sources

```go
import (
    "github.com/samber/ro"
    roproc "github.com/samber/ro/plugins/proc"
    "time"
)

// Monitor both CPU and memory
cpuObservable := roproc.NewCPUInfoWatcher(2 * time.Second)
memObservable := roproc.NewVirtualMemoryWatcher(2 * time.Second)

// Combine them using Merge
combined := ro.Merge(cpuObservable, memObservable)

subscription := combined.Subscribe(ro.NewObserver(
    func(value interface{}) {
        switch v := value.(type) {
        case cpu.InfoStat:
            fmt.Printf("CPU: %s\n", v.Name)
        case *mem.VirtualMemoryStat:
            fmt.Printf("Memory Used: %d\n", v.Used)
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

### Filtering and Processing

```go
import (
    "github.com/samber/ro"
    roproc "github.com/samber/ro/plugins/proc"
    "time"
)

// Monitor memory and filter high usage
observable := ro.Pipe1(
    roproc.NewVirtualMemoryWatcher(1 * time.Second),
    ro.Filter(func(stats *mem.VirtualMemoryStat) bool {
        return float64(stats.Used)/float64(stats.Total) > 0.8 // > 80% usage
    }),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats *mem.VirtualMemoryStat) {
        usage := float64(stats.Used) / float64(stats.Total) * 100
        fmt.Printf("High memory usage: %.1f%%\n", usage)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

## Error Handling

All sources handle errors gracefully and will emit error notifications if the underlying system calls fail. The error handling follows the reactive stream pattern:

```go
observable := roproc.NewVirtualMemoryWatcher(5 * time.Second)

subscription := observable.Subscribe(ro.NewObserver(
    func(stats *mem.VirtualMemoryStat) {
        // Handle successful data
    },
    func(err error) {
        // Handle errors (e.g., permission denied, system call failed)
        log.Printf("Monitoring error: %v", err)
    },
    func() {
        // Handle completion
    },
))
```

## Dependencies

This plugin requires the [gopsutil](https://github.com/shirou/gopsutil) library:

```bash
go get github.com/shirou/gopsutil/v4
```

## Performance Considerations

- Choose appropriate intervals for your monitoring needs
- Consider using `ro.ObserveOn` to control the scheduler for CPU-intensive operations
- Use `ro.SubscribeOn` to control where the source runs
- Be mindful of system resources when monitoring frequently 