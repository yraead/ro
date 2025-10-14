
# Prometheus (Enterprise Edition)

## Configuration

The plugin uses a `CollectorConfig` struct to configure Prometheus metrics collection:

```go
type CollectorConfig struct {
    // Namespace is the namespace of the metrics.
    // Example: "myapp" -> metrics will be prefixed with "myapp_"
    Namespace string
    
    // Subsystem is the subsystem of the metrics.
    // Example: "pipeline" -> metrics will be prefixed with "myapp_pipeline_"
    Subsystem string
    
    // ConstLabels are labels that are applied to all metrics.
    // These labels help identify the source of metrics across different instances.
    ConstLabels prometheus.Labels
    
    // SummaryObjectives defines the quantile objectives for summary metrics.
    // Default values provide good coverage for most use cases.
    SummaryObjectives map[float64]float64
}
```

### Configuration Options

- **Namespace**: Sets the namespace prefix for all metrics. This helps organize metrics by application or service.
- **Subsystem**: Sets the subsystem prefix for all metrics. This provides additional categorization within the namespace.
- **ConstLabels**: Adds constant labels to all metrics. Useful for identifying environment, version, or instance information.
- **SummaryObjectives**: Configures quantile objectives for summary metrics (latency, processing time, etc.).

### Default Values

The plugin provides sensible defaults for summary objectives:

```go
var DefaultSummaryObjectives = map[float64]float64{
    0.5:  0.05,  // 50th percentile with 5% error
    0.9:  0.01,  // 90th percentile with 1% error
    0.95: 0.005, // 95th percentile with 0.5% error
    0.99: 0.001, // 99th percentile with 0.1% error
}
```

## Example

```go
import (
    "github.com/samber/ro"
    roprometheus "github.com/samber/ro/ee/plugins/prometheus"
)

var obs, collector = roprometheus.Pipe3(
    roprometheus.CollectorConfig{
        Namespace: "example",
        Subsystem: "",
        ConstLabels: prometheus.Labels{
            "foo": "bar",
        },
    },
    ro.Just(1, 2, 3),
    ro.Map(func(v int64) int64 {
        return v*2
    }),
    ro.Take(2),
)

func main() {
    prometheus.MustRegister(collector)

    obs.Subscribe(
        ...
    )
}
```
