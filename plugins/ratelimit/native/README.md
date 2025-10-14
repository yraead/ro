# Native Rate Limiter Plugin

The native rate limiter plugin provides operators for rate limiting using Go's built-in time-based windowing.

## Installation

```bash
go get github.com/samber/ro/plugins/ratelimit/native
```

## Operators

### NewRateLimiter

Creates a rate limiter operator that limits values based on a count and time interval.

```go
import (
    "time"
    "github.com/samber/ro"
    roratelimit "github.com/samber/ro/plugins/ratelimit/native"
)

// Create a rate limiter: 5 items per second per key
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
    roratelimit.NewRateLimiter[string](5, time.Second, func(userID string) string {
        return userID
    }),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

## Parameters

### Count

The maximum number of items allowed in the time window:

```go
// Allow 10 items per window
observable := ro.Pipe1(
    ro.Just("item1", "item2", "item3", "item4", "item5"),
    roratelimit.NewRateLimiter[string](10, time.Second, func(item string) string {
        return "default"
    }),
)
```

### Interval

The time window for rate limiting:

```go
// 5 items per second
observable := ro.Pipe1(
    ro.Just("item1", "item2", "item3", "item4", "item5"),
    roratelimit.NewRateLimiter[string](5, time.Second, func(item string) string {
        return "default"
    }),
)

// 100 items per minute
observable := ro.Pipe1(
    ro.Just("item1", "item2", "item3", "item4", "item5"),
    roratelimit.NewRateLimiter[string](100, time.Minute, func(item string) string {
        return "default"
    }),
)

// 1000 items per hour
observable := ro.Pipe1(
    ro.Just("item1", "item2", "item3", "item4", "item5"),
    roratelimit.NewRateLimiter[string](1000, time.Hour, func(item string) string {
        return "default"
    }),
)
```

### Key Function

A function that extracts the key for rate limiting:

```go
type Request struct {
    UserID string
    Action string
    Data   string
}

// Rate limit by user ID
observable := ro.Pipe1(
    ro.Just(
        Request{UserID: "user1", Action: "login", Data: "data1"},
        Request{UserID: "user2", Action: "login", Data: "data2"},
        Request{UserID: "user1", Action: "logout", Data: "data3"},
    ),
    roratelimit.NewRateLimiter[Request](3, time.Minute, func(req Request) string {
        return req.UserID
    }),
)
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
    roratelimit.NewRateLimiter[Request](5, time.Minute, func(req Request) string {
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
    roratelimit.NewRateLimiter[APIRequest](10, time.Minute, func(req APIRequest) string {
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
    roratelimit.NewRateLimiter[APIRequest](2, time.Second, func(req APIRequest) string {
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
    roratelimit.NewRateLimiter[APIRequest](5, time.Minute, func(req APIRequest) string {
        return req.IPAddress + ":" + req.Endpoint
    }),
)
```

## Real-world Example

Here's a practical example that rate limits API requests:

```go
import (
    "time"
    "github.com/samber/ro"
    roratelimit "github.com/samber/ro/plugins/ratelimit/native"
)

type APIRequest struct {
    UserID   string
    Endpoint string
    Method   string
    Data     string
}

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
    // Apply rate limiting: 10 requests per minute per user
    roratelimit.NewRateLimiter[APIRequest](10, time.Minute, func(req APIRequest) string {
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
            // Handle errors
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses Go's built-in time-based windowing for rate limiting
- Rate limiting is applied per key (user, IP, endpoint, etc.)
- Memory usage scales with the number of unique keys
- The algorithm uses sliding windows for accurate rate limiting
- Consider the count and interval for your use case
- The plugin automatically handles rate limit checking and filtering
- Choose appropriate key generation strategies for your application
- This implementation is suitable for single-instance applications 