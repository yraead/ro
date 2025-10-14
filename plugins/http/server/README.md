# Ro HTTP Server Plugin

This plugin provides reactive HTTP server functionality for the [Ro](https://github.com/samber/ro) reactive programming library. It allows you to:

- Create HTTP servers that handle incoming requests as reactive streams
- Apply reactive operators to HTTP request processing
- Build scalable and composable HTTP APIs using reactive patterns
- Handle WebSocket connections within HTTP servers

## Installation

```bash
go get github.com/samber/ro/plugins/http/server
```

## Requirements

- [Ro](https://github.com/samber/ro) reactive programming library
- Go 1.18 or later

## Quick Start

```go
package main

import (
	"fmt"
	"net/http"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

func main() {
	// Create a simple HTTP server
	server := rohttp.NewServer(":8080")
	
	// Handle GET requests to /hello
	server.HandleRoute("GET", "/hello", func(req *http.Request) ro.Observable[string] {
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		
		return ro.Just(fmt.Sprintf("Hello, %s!", name))
	})
	
	// Start the server
	if err := server.Start(); err != nil {
		panic(err)
	}
}
```

## API Reference

### Server Configuration

#### `NewServer(addr string) *Server`

Creates a new HTTP server instance that listens on the specified address.

```go
server := rohttp.NewServer(":8080")
```

#### `NewServerWithOptions(options ServerOptions) *Server`

Creates a new HTTP server with custom configuration options.

```go
server := rohttp.NewServerWithOptions(rohttp.ServerOptions{
	Addr:         ":8080",
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
	IdleTimeout:  120 * time.Second,
})
```

### Route Handling

#### `HandleRoute(method, pattern string, handler func(*http.Request) ro.Observable[T])`

Registers a route handler that returns an observable. The observable's values will be serialized and sent as HTTP responses.

```go
server.HandleRoute("GET", "/api/users", func(req *http.Request) ro.Observable[User] {
    return fetchUsersFromDatabase()
})
```

#### `HandleRouteFunc(method, pattern string, handler func(*http.Request) (interface{}, error))`

Registers a route handler that returns a value and error, useful for traditional HTTP handlers.

```go
server.HandleRouteFunc("POST", "/api/users", func(req *http.Request) (interface{}, error) {
    var user User
    if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
        return nil, err
    }
    
    createdUser, err := createUserInDatabase(user)
    return createdUser, err
})
```

### Middleware

#### `Use(middleware func(http.Handler) http.Handler)`

Adds middleware to the server.

```go
// CORS middleware
server.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
})

// Logging middleware
server.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start)
        fmt.Printf("%s %s %d %v\n", r.Method, r.URL.Path, wrapped.statusCode, duration)
    })
})
```

### Server Control

#### `Start() error`

Starts the HTTP server and begins listening for incoming requests.

#### `Stop() error`

Gracefully stops the HTTP server.

#### `IsRunning() bool`

Returns true if the server is currently running.

## Usage Examples

### Basic REST API

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

func main() {
	server := rohttp.NewServer(":8080")
	
	// Get all users
	server.HandleRoute("GET", "/api/users", func(req *http.Request) ro.Observable[User] {
		return ro.FromSlice(users)
	})
	
	// Get user by ID
	server.HandleRoute("GET", "/api/users/{id}", func(req *http.Request) ro.Observable[User] {
		id := extractIDFromPath(req.URL.Path, "/api/users/")
		
		return ro.Pipe1(
			ro.FromSlice(users),
			ro.Filter(func(user User) bool {
				return user.ID == id
			}),
			ro.Take[User](1),
		)
	})
	
	// Create new user
	server.HandleRouteFunc("POST", "/api/users", func(req *http.Request) (interface{}, error) {
		var newUser User
		if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		
		newUser.ID = len(users) + 1
		users = append(users, newUser)
		
		return newUser, nil
	})
	
	fmt.Println("Server starting on :8080")
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func extractIDFromPath(path, prefix string) int {
	// Extract ID from path like "/api/users/123"
	// In a real implementation, you'd use a proper router
	return 123 // Placeholder
}
```

### Reactive Data Streaming

```go
package main

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

func main() {
	server := rohttp.NewServer(":8080")
	
	// Server-sent events endpoint
	server.HandleRoute("GET", "/api/events", func(req *http.Request) ro.Observable[string] {
		return ro.Pipe1(
			ro.Interval(1*time.Second),
			ro.Map(func(tick int64) string {
				return fmt.Sprintf("data: Server time: %s\n\n", time.Now().Format(time.RFC3339))
			}),
		)
	})
	
	// Real-time data processing
	server.HandleRoute("GET", "/api/processed-data", func(req *http.Request) ro.Observable[float64] {
		// Simulate incoming data stream
		dataStream := ro.Pipe1(
			ro.Interval(500*time.Millisecond),
			ro.Map(func(tick int64) float64 {
				// Simulate sensor data
				return float64(tick%100) * 1.5
			}),
		)
		
		// Apply reactive transformations
		return ro.Pipe3(
			dataStream,
			ro.Filter(func(value float64) bool {
				return value > 50.0 // Filter out low values
			}),
			ro.Map(func(value float64) float64 {
				return value * 2.0 // Double the values
			}),
			ro.BufferCount[[]float64](5), // Buffer in groups of 5
			ro.Map(func(values []float64) float64 {
				// Calculate average
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				return sum / float64(len(values))
			}),
		)
	})
	
	fmt.Println("Server starting on :8080")
	if err := server.Start(); err != nil {
		panic(err)
	}
}
```

### Authentication and Authorization

```go
package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

func main() {
	server := rohttp.NewServer(":8080")
	
	// Authentication middleware
	server.Use(authMiddleware)
	
	// Public endpoint
	server.HandleRoute("GET", "/public", func(req *http.Request) ro.Observable[string] {
		return ro.Just("This is a public endpoint")
	})
	
	// Protected endpoint
	server.HandleRoute("GET", "/protected", func(req *http.Request) ro.Observable[string] {
		user := req.Context().Value("user").(string)
		return ro.Just(fmt.Sprintf("Hello, %s! This is protected data.", user))
	})
	
	// Admin endpoint
	server.HandleRoute("GET", "/admin", func(req *http.Request) ro.Observable[string] {
		user := req.Context().Value("user").(string)
		role := req.Context().Value("role").(string)
		
		if role != "admin" {
			return ro.Throw[string](fmt.Errorf("access denied"))
		}
		
		return ro.Just(fmt.Sprintf("Admin panel for %s", user))
	})
	
	fmt.Println("Server starting on :8080")
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		if r.URL.Path == "/public" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		// Parse Basic Auth
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		
		// Validate credentials (in production, use proper authentication)
		username, password := credentials[0], credentials[1]
		if !isValidUser(username, password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		
		// Add user info to context
		ctx := context.WithValue(r.Context(), "user", username)
		ctx = context.WithValue(ctx, "role", getUserRole(username))
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isValidUser(username, password string) bool {
	// In production, validate against database
	return username == "admin" && password == "secret" ||
		username == "user" && password == "password"
}

func getUserRole(username string) string {
	if username == "admin" {
		return "admin"
	}
	return "user"
}
```

### File Upload Handler

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

type UploadResponse struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
}

func main() {
	server := rohttp.NewServer(":8080")
	
	// File upload endpoint
	server.HandleRouteFunc("POST", "/upload", func(req *http.Request) (interface{}, error) {
		// Parse multipart form (max 32MB)
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			return nil, fmt.Errorf("failed to parse form: %w", err)
		}
		
		file, handler, err := req.FormFile("file")
		if err != nil {
			return nil, fmt.Errorf("failed to get file: %w", err)
		}
		defer file.Close()
		
		// Create upload directory if it doesn't exist
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}
		
		// Create destination file
		dstPath := filepath.Join(uploadDir, handler.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dst.Close()
		
		// Copy file content
		size, err := io.Copy(dst, file)
		if err != nil {
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
		
		return UploadResponse{
			Filename: handler.Filename,
			Size:     size,
			Path:     dstPath,
		}, nil
	})
	
	// Multiple file upload with reactive processing
	server.HandleRoute("POST", "/upload-batch", func(req *http.Request) ro.Observable[UploadResponse] {
		return ro.Pipe2(
			// Create observable from uploaded files
			ro.Create(func(observer ro.Observer[UploadResponse]) {
				if err := req.ParseMultipartForm(32 << 20); err != nil {
					observer.Error(fmt.Errorf("failed to parse form: %w", err))
					return
				}
				
				files := req.MultipartForm.File["files"]
				for _, fileHeader := range files {
					file, err := fileHeader.Open()
					if err != nil {
						observer.Error(fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err))
						return
					}
					
					// Process file
					dstPath := filepath.Join("./uploads", fileHeader.Filename)
					dst, err := os.Create(dstPath)
					if err != nil {
						file.Close()
						observer.Error(fmt.Errorf("failed to create file %s: %w", fileHeader.Filename, err))
						return
					}
					
					size, err := io.Copy(dst, file)
					file.Close()
					dst.Close()
					
					if err != nil {
						observer.Error(fmt.Errorf("failed to save file %s: %w", fileHeader.Filename, err))
						return
					}
					
					observer.Next(UploadResponse{
						Filename: fileHeader.Filename,
						Size:     size,
						Path:     dstPath,
					})
				}
				observer.Complete()
			}),
			// Apply reactive transformations
			ro.Map(func(response UploadResponse) UploadResponse {
				// Add processing metadata
				return response
			}),
			// Filter successful uploads
			ro.Filter(func(response UploadResponse) bool {
				return response.Size > 0
			}),
		)
	})
	
	fmt.Println("Server starting on :8080")
	if err := server.Start(); err != nil {
		panic(err)
	}
}
```

### Error Handling and Recovery

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func main() {
	server := rohttp.NewServer(":8080")
	
	// Global error recovery middleware
	server.Use(recoveryMiddleware)
	
	// Route that might panic
	server.HandleRoute("GET", "/panic", func(req *http.Request) ro.Observable[string] {
		panic("This is a controlled panic!")
	})
	
	// Route with validation error
	server.HandleRoute("GET", "/validate", func(req *http.Request) ro.Observable[string] {
		id := req.URL.Query().Get("id")
		if id == "" {
			return ro.Throw[string](fmt.Errorf("id parameter is required"))
		}
		
		return ro.Just(fmt.Sprintf("Valid ID: %s", id))
	})
	
	// Route with async error
	server.HandleRoute("GET", "/async-error", func(req *http.Request) ro.Observable[string] {
		return ro.Pipe1(
			ro.Throw[string](fmt.Errorf("async operation failed")),
			ro.CatchError(func(err error) ro.Observable[string] {
				return ro.Just(fmt.Sprintf("Caught error: %v", err))
			}),
		)
	})
	
	fmt.Println("Server starting on :8080")
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic recovered: %v\n", err)
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				
				errorResponse := ErrorResponse{
					Error:   "internal_server_error",
					Message: "An unexpected error occurred",
					Code:    http.StatusInternalServerError,
				}
				
				json.NewEncoder(w).Encode(errorResponse)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}
```

### Custom Response Serialization

```go
package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/http/server"
)

type User struct {
	ID    int    `json:"id" xml:"id"`
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

type UserList struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `json:"users" xml:"user"`
}

func main() {
	server := rohttp.NewServer(":8080")
	
	// JSON response (default)
	server.HandleRoute("GET", "/users.json", func(req *http.Request) ro.Observable[UserList] {
		users := []User{
			{ID: 1, Name: "John Doe", Email: "john@example.com"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
		}
		
		return ro.Just(UserList{Users: users})
	})
	
	// XML response
	server.HandleRouteWithSerializer("GET", "/users.xml", func(req *http.Request) ro.Observable[UserList] {
		users := []User{
			{ID: 1, Name: "John Doe", Email: "john@example.com"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
		}
		
		return ro.Just(UserList{Users: users})
	}, func(data UserList) ([]byte, string, error) {
		xmlData, err := xml.Marshal(data)
		if err != nil {
			return nil, "", err
		}
		return xmlData, "application/xml", nil
	})
	
	// Plain text response
	server.HandleRouteWithSerializer("GET", "/users.txt", func(req *http.Request) ro.Observable[UserList] {
		users := []User{
			{ID: 1, Name: "John Doe", Email: "john@example.com"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
		}
		
		return ro.Just(UserList{Users: users})
	}, func(data UserList) ([]byte, string, error) {
		text := fmt.Sprintf("Users:\n")
		for _, user := range data.Users {
			text += fmt.Sprintf("- %s (%s)\n", user.Name, user.Email)
		}
		return []byte(text), "text/plain", nil
	})
	
	fmt.Println("Server starting on :8080")
	if err := server.Start(); err != nil {
		panic(err)
	}
}
```

## Configuration Options

### ServerOptions

```go
type ServerOptions struct {
	Addr         string        // Server address (e.g., ":8080")
	ReadTimeout  time.Duration // Maximum duration for reading requests
	WriteTimeout time.Duration // Maximum duration for writing responses
	IdleTimeout  time.Duration // Maximum time to wait for next request
	TLSConfig    *tls.Config   // TLS configuration for HTTPS
}

// Example usage
server := rohttp.NewServerWithOptions(rohttp.ServerOptions{
	Addr:         ":8080",
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
	IdleTimeout:  120 * time.Second,
})
```

## Best Practices

1. **Use reactive operators** for data transformation and filtering
2. **Implement proper error handling** with catch operators and middleware
3. **Add authentication and authorization** through middleware
4. **Use context** for request-scoped values and cancellation
5. **Implement graceful shutdown** for production deployments
6. **Add logging and monitoring** for observability
7. **Validate input** before processing requests
8. **Use appropriate content types** for different response formats

## Performance Considerations

1. **Connection pooling**: The server reuses HTTP connections efficiently
2. **Backpressure**: Reactive streams naturally handle backpressure
3. **Memory management**: Avoid holding large observables in memory
4. **Timeouts**: Set appropriate timeouts for production use
5. **Graceful shutdown**: Allow in-flight requests to complete

## License

Apache 2.0 - See [LICENSE](https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md) for details.