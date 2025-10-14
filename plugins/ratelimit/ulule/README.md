# Ulule Rate Limiter Plugin

The Ulule rate limiter plugin provides operators for rate limiting using the `ulule/limiter` package.

## Installation

```bash
go get github.com/samber/ro/plugins/ratelimit/ulule
```

## Operators

### NewRateLimiter

Creates a rate limiter operator that filters values based on rate limits.

```go
import (
    "time"
    "github.com/samber/ro"
    roratelimit "github.com/samber/ro/plugins/ratelimit/ulule"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

// Create a rate limiter with 5 requests per second
store := memory.NewStore()
rate := limiter.Rate{
    Period: time.Second,
    Limit:  5,
}
limiter := limiter.New(store, rate)

// Rate limit by user ID
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
    roratelimit.NewRateLimiter[string](limiter, func(userID string) string {
        return userID
    }),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

## Rate Limiter Configuration

### Rate Definition

```go
// 5 requests per second
rate := limiter.Rate{
    Period: time.Second,
    Limit:  5,
}

// 100 requests per minute
rate := limiter.Rate{
    Period: time.Minute,
    Limit:  100,
}

// 1000 requests per hour
rate := limiter.Rate{
    Period: time.Hour,
    Limit:  1000,
}
```

### Store Configuration

```go
// Memory store (in-process)
store := memory.NewStore()

// Redis store (distributed)
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
store := redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
    Prefix: "rate_limiter",
})
```

## Key Generation Strategies

### User-based Rate Limiting

```go
type Request struct {
    UserID string
    Action string
    Data   string
}

observable := ro.Pipe1(
    ro.Just(
        Request{UserID: "user1", Action: "login", Data: "data1"},
        Request{UserID: "user2", Action: "login", Data: "data2"},
        Request{UserID: "user1", Action: "logout", Data: "data3"},
    ),
    roratelimit.NewRateLimiter[Request](limiter, func(req Request) string {
        return req.UserID
    }),
)
```

### IP-based Rate Limiting

```go
type APIRequest struct {
    IPAddress string
    Endpoint  string
    Method    string
}

observable := ro.Pipe1(
    ro.Just(
        APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "GET"},
        APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/users", Method: "GET"},
        APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/posts", Method: "POST"},
    ),
    roratelimit.NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
        return req.IPAddress
    }),
)
```

### Endpoint-based Rate Limiting

```go
observable := ro.Pipe1(
    ro.Just(
        APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "GET"},
        APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/posts", Method: "GET"},
        APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "POST"},
    ),
    roratelimit.NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
        return req.Endpoint
    }),
)
```

### Composite Key Rate Limiting

```go
observable := ro.Pipe1(
    ro.Just(
        APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "GET"},
        APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/users", Method: "GET"},
        APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/posts", Method: "GET"},
    ),
    roratelimit.NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
        return req.IPAddress + ":" + req.Endpoint
    }),
)
```

## Error Handling

The plugin handles rate limiting errors gracefully:

```go
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
    roratelimit.NewRateLimiter[string](limiter, func(userID string) string {
        return userID
    }),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value string) {
            // Handle successful rate-limited value
        },
        func(err error) {
            // Handle rate limiting error
            // This could be due to:
            // - Store errors
            // - Context cancellation
            // - Other limiter errors
        },
        func() {
            // Handle completion
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
    roratelimit "github.com/samber/ro/plugins/ratelimit/ulule"
)

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

observable := ro.Pipe1(
    ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
    roratelimit.NewRateLimiter[string](limiter, func(userID string) string {
        return userID
    }),
)

subscription := observable.SubscribeWithContext(
    ctx,
    ro.NewObserverWithContext(
        func(ctx context.Context, value string) {
            // Handle rate-limited value with context
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

## Real-world Example

Here's a practical example that rate limits API requests:

```go
import (
    "time"
    "github.com/samber/ro"
    roratelimit "github.com/samber/ro/plugins/ratelimit/ulule"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

type APIRequest struct {
    UserID   string
    Endpoint string
    Method   string
    Data     string
}

// Create rate limiter: 10 requests per minute per user
store := memory.NewStore()
rate := limiter.Rate{
    Period: time.Minute,
    Limit:  10,
}
limiter := limiter.New(store, rate)

// Process API requests with rate limiting
pipeline := ro.Pipe2(
    // Simulate API requests
    ro.Just(
        APIRequest{UserID: "user1", Endpoint: "/api/users", Method: "GET", Data: "data1"},
        APIRequest{UserID: "user2", Endpoint: "/api/posts", Method: "GET", Data: "data2"},
        APIRequest{UserID: "user1", Endpoint: "/api/users", Method: "POST", Data: "data3"},
        APIRequest{UserID: "user3", Endpoint: "/api/comments", Method: "GET", Data: "data4"},
        APIRequest{UserID: "user1", Endpoint: "/api/users", Method: "PUT", Data: "data5"},
    ),
    // Apply rate limiting per user
    roratelimit.NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
        return req.UserID
    }),
)

subscription := pipeline.Subscribe(
    ro.NewObserver(
        func(req APIRequest) {
            // Process rate-limited request
            // Only requests within rate limit will be processed
        },
        func(err error) {
            // Handle rate limiting errors
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses the `ulule/limiter` package for rate limiting
- Rate limiting is applied per key (user, IP, endpoint, etc.)
- Use appropriate stores for your use case:
  - Memory store for single-instance applications
  - Redis store for distributed applications
- Consider the rate limit period and limit for your use case
- The plugin automatically handles rate limit checking and filtering
- Context cancellation properly stops rate limiting operations
- Choose appropriate key generation strategies for your application 