# HTTP Plugin

The HTTP plugin provides operators for making HTTP requests in reactive streams.

## Installation

```bash
go get github.com/samber/ro/plugins/http
```

## Operators

### HTTPRequest

Sends HTTP requests and returns responses as an observable stream.

```go
import (
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

// Create HTTP request
req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)

// Send HTTP request
observable := rohttp.HTTPRequest(req, nil)

subscription := observable.Subscribe(ro.PrintObserver[*http.Response]())
defer subscription.Unsubscribe()

// Output:
// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
// Completed
```

### HTTPRequestJSON

Sends HTTP requests and automatically parses JSON responses into the specified type. This is a convenience operator that combines HTTPRequest with JSON parsing.

```go
import (
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

// Create HTTP request
req, _ := http.NewRequest("GET", "https://api.example.com/users/1", nil)

// Send HTTP request and parse JSON response
observable := rohttp.HTTPRequestJSON[User](req, nil)

subscription := observable.Subscribe(ro.PrintObserver[User]())
defer subscription.Unsubscribe()

// Output:
// Next: {ID:1 Name:John Doe Email:john@example.com}
// Completed
```

## Basic Usage

### Simple GET Request

```go
import (
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

// Create a simple GET request
req, _ := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)

// Send the request
observable := rohttp.HTTPRequest(req, nil)

subscription := observable.Subscribe(ro.PrintObserver[*http.Response]())
defer subscription.Unsubscribe()
```

### Request with Headers

```go
// Create request with custom headers
req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
req.Header.Set("User-Agent", "ro-http-client/1.0")
req.Header.Set("Accept", "application/json")
req.Header.Set("Authorization", "Bearer your-token")

observable := rohttp.HTTPRequest(req, nil)

subscription := observable.Subscribe(ro.PrintObserver[*http.Response]())
defer subscription.Unsubscribe()
```

### Custom HTTP Client

```go
import (
    "net/http"
    "time"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

// Create custom HTTP client with timeout
client := &http.Client{
    Timeout: 10 * time.Second,
}

req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
observable := rohttp.HTTPRequest(req, client)

subscription := observable.Subscribe(ro.PrintObserver[*http.Response]())
defer subscription.Unsubscribe()
```

## Advanced Usage

### Processing Responses

```go
import (
    "fmt"
    "io"
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

// Send request and process response
observable := ro.Pipe1(
    rohttp.HTTPRequest(req, nil),
    ro.Map(func(resp *http.Response) string {
        // Read response body
        body, _ := io.ReadAll(resp.Body)
        defer resp.Body.Close()
        
        return fmt.Sprintf("Status: %s, Body: %s", resp.Status, string(body))
    }),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: Status: 200 OK, Body: {"id": 1, "title": "..."}
// Completed
```

### Multiple Concurrent Requests

```go
import (
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

// Create multiple requests
urls := []string{
    "https://api.example.com/users/1",
    "https://api.example.com/users/2",
    "https://api.example.com/users/3",
}

// Convert URLs to requests
requests := make([]*http.Request, len(urls))
for i, url := range urls {
    req, _ := http.NewRequest("GET", url, nil)
    requests[i] = req
}

// Process multiple requests concurrently
observable := ro.Pipe1(
    ro.FromSlice(requests),
    ro.MergeMap(func(req *http.Request) ro.Observable[*http.Response] {
        return rohttp.HTTPRequest(req, nil)
    }),
)

subscription := observable.Subscribe(ro.PrintObserver[*http.Response]())
defer subscription.Unsubscribe()

// Output:
// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
// Completed
```

### Error Handling

```go
import (
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

req, _ := http.NewRequest("GET", "https://invalid-url-that-will-fail.com", nil)

observable := rohttp.HTTPRequest(req, nil)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(resp *http.Response) {
            // Handle successful response
            // Note: HTTP status >= 400 is not considered an error by this operator
            if resp.StatusCode >= 400 {
                // Handle HTTP error status manually
            }
        },
        func(err error) {
            // Handle network errors, timeouts, etc.
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

### Context and Timeout

```go
import (
    "context"
    "net/http"
    "time"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

// Create request with timeout context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, _ := http.NewRequest("GET", "https://api.example.com/slow-endpoint", nil)
req = req.WithContext(ctx)

observable := rohttp.HTTPRequest(req, nil)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(resp *http.Response) {
            // Handle successful response
        },
        func(err error) {
            // Handle timeout or other errors
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Real-world Example

Here's a practical example that fetches user data from an API and processes it:

```go
import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
    rojson "github.com/samber/ro/plugins/encoding/json"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

// Create a pipeline that fetches and processes user data
pipeline := ro.Pipe4(
    // Create requests for multiple users
    ro.Just(1, 2, 3),
    ro.Map(func(id int) *http.Request {
        req, _ := http.NewRequest("GET", fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id), nil)
        return req
    }),
    // Send HTTP requests
    ro.MergeMap(func(req *http.Request) ro.Observable[*http.Response] {
        return rohttp.HTTPRequest(req, nil)
    }),
    // Process responses
    ro.Map(func(resp *http.Response) User {
        body, _ := io.ReadAll(resp.Body)
        defer resp.Body.Close()
        
        var user User
        json.Unmarshal(body, &user)
        return user
    }),
)

// Add timeout context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

subscription := pipeline.SubscribeWithContext(ctx, ro.PrintObserver[User]())
defer subscription.Unsubscribe()

// Output:
// Next: {ID:1 Name:Leanne Graham Email:Sincere@april.biz}
// Next: {ID:2 Name:Ervin Howell Email:Shanna@melissa.tv}
// Next: {ID:3 Name:Clementine Bauch Email:Nathan@yesenia.net}
// Completed
```

## HTTPRequestJSON Usage Examples

### Simple JSON String

```go
import (
    "net/http"
    "github.com/samber/ro"
    rohttp "github.com/samber/ro/plugins/http"
)

req, _ := http.NewRequest("GET", "https://api.example.com/message", nil)

observable := rohttp.HTTPRequestJSON[string](req, nil)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

### JSON Array

```go
req, _ := http.NewRequest("GET", "https://api.example.com/tags", nil)

observable := rohttp.HTTPRequestJSON[[]string](req, nil)

subscription := observable.Subscribe(ro.PrintObserver[[]string]())
defer subscription.Unsubscribe()
```

### Complex JSON with Error Handling

```go
type APIResponse struct {
    Success bool   `json:"success"`
    Data    []User `json:"data"`
    Error   string `json:"error,omitempty"`
}

req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)

observable := rohttp.HTTPRequestJSON[APIResponse](req, nil)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(response APIResponse) {
            if response.Success {
                fmt.Printf("Got %d users\n", len(response.Data))
            } else {
                fmt.Printf("API Error: %s\n", response.Error)
            }
        },
        func(err error) {
            fmt.Printf("Network or JSON error: %s\n", err)
        },
        func() {
            fmt.Println("Request completed")
        },
    ),
)
defer subscription.Unsubscribe()
```

### HTTPRequestJSON with Processing Pipeline

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

pipeline := ro.Pipe2(
    // Fetch multiple users concurrently
    ro.FromSlice([]int{1, 2, 3}),
    ro.MergeMap(func(userID int) ro.Observable[User] {
        req, _ := http.NewRequest("GET", fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", userID), nil)
        return rohttp.HTTPRequestJSON[User](req, nil)
    }),
)

subscription := pipeline.Subscribe(ro.PrintObserver[User]())
defer subscription.Unsubscribe()
```

## Important Notes

1. **HTTP Status Codes**: HTTP status codes >= 400 are **not** considered errors by these operators. You need to handle HTTP errors manually if needed.

2. **Response Body**: For `HTTPRequest`, remember to call `resp.Body.Close()` when you're done with the response to avoid resource leaks. `HTTPRequestJSON` handles this automatically.

3. **JSON Parsing**: `HTTPRequestJSON` automatically closes the response body and will emit an error if the JSON cannot be parsed.

4. **Context**: Use request context for timeouts and cancellation.

5. **Concurrency**: The operators are designed for concurrent use and are thread-safe.

6. **Error Handling**: Network errors, timeouts, and other transport errors will be emitted as error notifications.

## Performance Considerations

- The operator uses Go's standard `http.Client` for requests
- Requests are executed asynchronously
- Use `MergeMap` for concurrent requests
- Consider connection pooling for high-throughput scenarios
- Use appropriate timeouts to prevent hanging requests 