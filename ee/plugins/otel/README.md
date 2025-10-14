# OpenTelemetry (Enterprise Edition)

## Configuration

The plugin uses a `CollectorConfig` struct to configure observability behavior:

```go
type CollectorConfig struct {
    // Enable or disable observability features
    EnableTracing bool
    EnableMetrics bool
    EnableLogging bool

    // Use custom providers or global defaults
    TracerProvider trace.TracerProvider
    MetricProvider metric.MeterProvider
    LoggerProvider log.LoggerProvider

    // Custom attributes for all telemetry data
    TraceAttributes   []attribute.KeyValue
    MetricAttributes  []attribute.KeyValue
    LoggingAttributes []attribute.KeyValue

    // Histogram bucket configuration for metrics
    MetricHistogramObjectivesSeconds []float64

    // Log level configuration
    LogLevelSubscription log.Severity
    LogLevelNext        log.Severity
    LogLevelError       log.Severity
}
```

## Example

### Basic Usage

```go
import (
    "github.com/samber/ro"
    rootel "github.com/samber/ro/ee/plugins/otel"
    "go.opentelemetry.io/otel/attribute"
)

var obs, collector = rootel.Pipe3(
    rootel.CollectorConfig{
        EnableTracing: true,
        EnableMetrics: true,
        EnableLogging: true,
        TraceAttributes: []attribute.KeyValue{
            attribute.String("service", "my-service"),
            attribute.String("environment", "production"),
        },
    },
    ro.Just(1, 2, 3),
    ro.Map(func(v int64) int64 {
        return v * 2
    }),
    ro.Take(2),
)

func main() {
    // Subscribe to the observable
    subscription := obs.Subscribe(
        ...
    )
    defer subscription.Unsubscribe()
}
```

### Tracing Example

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/samber/ro"
    rolicense "github.com/samber/ro/ee/pkg/license"
    rootel "github.com/samber/ro/ee/plugins/otel"
    rocsv "github.com/samber/ro/plugins/encoding/csv"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type User struct {
    ID   string
    Name string
}

func getUsers(index int64) ([]User, error) {
    // Simulate database query
    return []User{
        {ID: "1", Name: "Alice"},
        {ID: "2", Name: "Bob"},
    }, nil
}

var pipeline, collector = rootel.Pipe7(
    rootel.CollectorConfig{
        EnableTracing: true,
        EnableMetrics: false,
        EnableLogging: false,
        TraceAttributes: []attribute.KeyValue{
            attribute.String("test-mode", "tracing"),
        },
    },
    ro.Range(0, 100),
    ro.MapErr(getUsers),
    ro.RetryWithConfig[[]User](ro.RetryConfig{
        MaxRetries: 2,
        Delay:      5 * time.Second,
    }),
    ro.TakeWhile(func(users []User) bool {
        return len(users) > 0
    }),
    ro.Flatten[User](),
    ro.Map(func(user User) []string {
        return []string{user.ID, user.Name}
    }),
    ro.StartWith([]string{"ID", "Name"}),
    rocsv.NewCSVWriter(csv.NewWriter(os.Stdout)),
)

func main() {
    // Set license
    err := rolicense.SetLicense("your-license-key")
    if err != nil {
        log.Fatalf("Failed to set license: %v", err)
    }

    ctx := context.Background()

    // Initialize OpenTelemetry trace exporter
    exp, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithInsecure(),
        otlptracegrpc.WithEndpoint("localhost:4317"),
    )
    if err != nil {
        log.Fatalf("Failed to create trace exporter: %v", err)
    }

    // Create resource with service information
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName("ee-otel-tracing"),
            semconv.ServiceVersion("1.0.0"),
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create resource: %v", err)
    }

    // Create trace provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)
    defer tp.Shutdown(ctx)

    // Subscribe to the pipeline
    subscription := pipeline.Subscribe(
        ...
    )
    defer subscription.Unsubscribe()

    log.Println("Processing completed!")
}
```

### Metrics Example

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/samber/ro"
    rolicense "github.com/samber/ro/ee/pkg/license"
    rootel "github.com/samber/ro/ee/plugins/otel"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var pipeline, collector = rootel.Pipe3(
    rootel.CollectorConfig{
        EnableLogging: false,
        EnableMetrics: true,
        EnableTracing: false,
        MetricAttributes: []attribute.KeyValue{
            attribute.String("service", "ee-otel-metrics"),
            attribute.String("environment", "demo"),
        },
    },
    ro.Just(1, 2, 3, 4, 5),
    ro.Map(func(v int64) int64 {
        return v * 2
    }),
    ro.Take(3),
)

func main() {
    // Set license
    err := rolicense.SetLicense("your-license-key")
    if err != nil {
        log.Fatalf("Failed to set license: %v", err)
    }

    ctx := context.Background()

    // Create resource with service information
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName("ee-otel-metrics"),
            semconv.ServiceVersion("1.0.0"),
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create resource: %v", err)
    }

    // Create a Prometheus registry
    reg := prometheus.NewRegistry()

    // Initialize Prometheus exporter for metrics
    exporter, err := otelprometheus.New(otelprometheus.WithRegisterer(reg))
    if err != nil {
        log.Fatalf("Failed to create Prometheus exporter: %v", err)
    }

    // Create meter provider
    mp := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(exporter),
        sdkmetric.WithResource(res),
    )
    otel.SetMeterProvider(mp)
    defer mp.Shutdown(ctx)

    // Start HTTP server for metrics endpoint
    go func() {
        handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
        http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
            handler.ServeHTTP(w, r)
        })
        log.Println("Starting OpenTelemetry metrics server on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatalf("Failed to start metrics server: %v", err)
        }
    }()

    // Subscribe to the pipeline
    subscription := pipeline.Subscribe(
        ...
    )
    defer subscription.Unsubscribe()

    log.Println("Processing completed! Metrics available at http://localhost:8080/metrics")

    // Keep the main goroutine alive to serve metrics
    select {}
}
```
