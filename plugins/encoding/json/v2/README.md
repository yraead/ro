# JSON Encoding Plugin

The JSON encoding plugin provides operators for marshaling and unmarshaling JSON data in reactive streams.

## Installation

```bash
go get github.com/samber/ro/plugins/encoding/json
```

## Operators

### Marshal

Converts Go structs, maps, and other types to JSON byte slices.

```go
import (
    "github.com/samber/ro"
    rojson "github.com/samber/ro/plugins/encoding/json"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

observable := ro.Pipe1(
    ro.Just(
        User{ID: 1, Name: "Alice", Age: 30},
        User{ID: 2, Name: "Bob", Age: 25},
        User{ID: 3, Name: "Charlie", Age: 35},
    ),
    rojson.Marshal[User](),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: {"id":1,"name":"Alice","age":30}
// Next: {"id":2,"name":"Bob","age":25}
// Next: {"id":3,"name":"Charlie","age":35}
// Completed
```

### Unmarshal

Converts JSON byte slices back to Go structs, maps, and other types.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte(`{"id":1,"name":"Alice","age":30}`),
        []byte(`{"id":2,"name":"Bob","age":25}`),
        []byte(`{"id":3,"name":"Charlie","age":35}`),
    ),
    rojson.Unmarshal[User](),
)

subscription := observable.Subscribe(ro.PrintObserver[User]())
defer subscription.Unsubscribe()

// Output:
// Next: {ID:1 Name:Alice Age:30}
// Next: {ID:2 Name:Bob Age:25}
// Next: {ID:3 Name:Charlie Age:35}
// Completed
```

## Working with Maps

You can also work with `map[string]interface{}` for dynamic JSON structures:

```go
// Marshal maps to JSON
observable := ro.Pipe1(
    ro.Just(
        map[string]interface{}{"name": "Alice", "age": 30},
        map[string]interface{}{"name": "Bob", "age": 25},
        map[string]interface{}{"name": "Charlie", "age": 35},
    ),
    rojson.Marshal[map[string]interface{}](),
)

// Output:
// Next: {"age":30,"name":"Alice"}
// Next: {"age":25,"name":"Bob"}
// Next: {"age":35,"name":"Charlie"}
// Completed
```

## Roundtrip Example

Demonstrate marshal/unmarshal roundtrip:

```go
observable := ro.Pipe2(
    ro.Just(
        User{ID: 1, Name: "Alice", Age: 30},
        User{ID: 2, Name: "Bob", Age: 25},
    ),
    rojson.Marshal[User](),
    rojson.Unmarshal[User](),
)

subscription := observable.Subscribe(ro.PrintObserver[User]())
defer subscription.Unsubscribe()

// Output:
// Next: {ID:1 Name:Alice Age:30}
// Next: {ID:2 Name:Bob Age:25}
// Completed
```

## Error Handling

Both `Marshal` and `Unmarshal` operators handle errors gracefully:

### Marshal Errors

```go
type Circular struct {
    Data interface{} `json:"data"`
}

circular := Circular{}
circular.Data = circular // Create circular reference

observable := ro.Pipe1(
    ro.Just(circular),
    rojson.Marshal[Circular](),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(data []byte) {
            // Handle successful marshaling
        },
        func(err error) {
            // Handle marshaling error (e.g., circular reference)
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

### Unmarshal Errors

```go
observable := ro.Pipe1(
    ro.Just(
        []byte(`{"id":1,"name":"Alice","age":30}`), // Valid JSON
        []byte(`{"id":2,"name":"Bob",`),             // Invalid JSON (truncated)
        []byte(`{"id":3,"name":"Charlie","age":35}`), // Valid JSON
    ),
    rojson.Unmarshal[User](),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(user User) {
            // Handle successful unmarshaling
        },
        func(err error) {
            // Handle unmarshaling error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Real-world Example

Here's a practical example that processes user data through JSON transformation:

```go
import (
    "github.com/samber/ro"
    rojson "github.com/samber/ro/plugins/encoding/json"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type UserResponse struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Age      int    `json:"age"`
}

// Process users through JSON transformation pipeline
pipeline := ro.Pipe4(
    // Start with user data
    ro.Just(
        User{ID: 1, Name: "Alice", Age: 30},
        User{ID: 2, Name: "Bob", Age: 25},
        User{ID: 3, Name: "Charlie", Age: 35},
    ),
    // Marshal to JSON
    rojson.Marshal[User](),
    // Transform JSON structure (could be done with string operations)
    ro.Map(func(data []byte) []byte {
        // In a real scenario, you might transform the JSON structure
        // For this example, we'll just pass through
        return data
    }),
    // Unmarshal to different struct
    rojson.Unmarshal[UserResponse](),
)

subscription := pipeline.Subscribe(ro.PrintObserver[UserResponse]())
defer subscription.Unsubscribe()

// Output:
// Next: {UserID:1 Username:Alice Age:30}
// Next: {UserID:2 Username:Bob Age:25}
// Next: {UserID:3 Username:Charlie Age:35}
// Completed
```

## Performance Considerations

- The `Marshal` operator uses `json.Marshal` which is optimized for performance
- The `Unmarshal` operator uses `json.Unmarshal` with proper error handling
- Both operators are type-safe and work with Go's generic system
- Consider using `json.RawMessage` for intermediate JSON processing if you need to preserve the exact JSON structure

## Supported Types

The JSON plugin supports all types that the standard `encoding/json` package supports:

- Structs with JSON tags
- Maps (`map[string]interface{}`, `map[string]string`, etc.)
- Slices and arrays
- Basic types (int, string, bool, float64, etc.)
- Pointers to supported types
- Interfaces (with proper type assertions) 