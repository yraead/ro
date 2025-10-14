# Gob Encoding Plugin

The gob encoding plugin provides operators for encoding and decoding data using Go's `encoding/gob` package for binary serialization.

## Installation

```bash
go get github.com/samber/ro/plugins/encoding/gob
```

## Operators

### Encode

Encodes values to gob bytes using Go's binary serialization format.

```go
import (
    "github.com/samber/ro"
    rogob "github.com/samber/ro/plugins/encoding/gob"
)

type Person struct {
    Name string
    Age  int
}

observable := ro.Pipe1(
    ro.Just(
        Person{Name: "Alice", Age: 30},
        Person{Name: "Bob", Age: 25},
        Person{Name: "Charlie", Age: 35},
    ),
    rogob.Encode[Person](),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [15 255 129 3 1 1 6 80 101 114 115 111 110 1 255 130 0 1 2 1 4 78 97 109 101 1 12 0 1 3 65 103 101 1 4 0 0 0 34 255 130 1 5 65 108 105 99 101 1 60 0 0 0 33 255 130 1 3 66 111 98 1 50 0 0 0 34 255 130 1 7 67 104 97 114 108 105 101 1 70 0 0]
// Completed
```

### Decode

Decodes gob bytes to values using Go's binary deserialization format.

```go
encoded := []byte{15, 255, 129, 3, 1, 1, 6, 80, 101, 114, 115, 111, 110, 1, 255, 130, 0, 1, 2, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 3, 65, 103, 101, 1, 4, 0, 0, 0, 34, 255, 130, 1, 5, 65, 108, 105, 99, 101, 1, 60, 0, 0, 0, 33, 255, 130, 1, 3, 66, 111, 98, 1, 50, 0, 0, 0, 34, 255, 130, 1, 7, 67, 104, 97, 114, 108, 105, 101, 1, 70, 0, 0}

observable := ro.Pipe1(
    ro.Just(encoded),
    rogob.Decode[Person](),
)

subscription := observable.Subscribe(ro.PrintObserver[Person]())
defer subscription.Unsubscribe()

// Output:
// Next: {Alice 30}
// Completed
```

## Supported Types

The gob plugin supports encoding and decoding of all Go types that are supported by the `encoding/gob` package:

### Basic Types

```go
// Strings
observable := ro.Pipe1(
    ro.Just("hello", "world", "golang"),
    rogob.Encode[string](),
)

// Integers
observable := ro.Pipe1(
    ro.Just(1, 2, 3, 4, 5),
    rogob.Encode[int](),
)

// Floats
observable := ro.Pipe1(
    ro.Just(3.14, 2.718, 1.414),
    rogob.Encode[float64](),
)

// Booleans
observable := ro.Pipe1(
    ro.Just(true, false, true),
    rogob.Encode[bool](),
)
```

### Complex Types

```go
// Structs
type User struct {
    ID   int
    Name string
    Tags []string
}

observable := ro.Pipe1(
    ro.Just(User{ID: 1, Name: "Alice", Tags: []string{"admin", "user"}}),
    rogob.Encode[User](),
)

// Maps
observable := ro.Pipe1(
    ro.Just(map[string]int{"a": 1, "b": 2, "c": 3}),
    rogob.Encode[map[string]int](),
)

// Slices
observable := ro.Pipe1(
    ro.Just([]int{1, 2, 3, 4, 5}),
    rogob.Encode[[]int](),
)
```

## Error Handling

Both `Encode` and `Decode` operators handle errors gracefully and will emit error notifications for invalid operations:

```go
// Encoding error (e.g., unsupported type)
observable := ro.Pipe1(
    ro.Just(make(chan int)), // Channels are not supported by gob
    rogob.Encode[chan int](),
)

// Decoding error (e.g., invalid gob data)
observable := ro.Pipe1(
    ro.Just([]byte("invalid-gob-data")),
    rogob.Decode[Person](),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value interface{}) {
            // Handle successful operation
        },
        func(err error) {
            // Handle error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Roundtrip Examples

Demonstrate roundtrip encoding and decoding:

```go
// Struct roundtrip
observable := ro.Pipe2(
    ro.Just(Person{Name: "Alice", Age: 30}),
    rogob.Encode[Person](),
    rogob.Decode[Person](),
)

subscription := observable.Subscribe(ro.PrintObserver[Person]())
defer subscription.Unsubscribe()

// Output:
// Next: {Alice 30}
// Completed
```

## Real-world Example

Here's a practical example that serializes user data for storage:

```go
import (
    "github.com/samber/ro"
    rogob "github.com/samber/ro/plugins/encoding/gob"
)

type User struct {
    ID       int
    Username string
    Email    string
    Settings map[string]interface{}
}

// Serialize users for storage
pipeline := ro.Pipe3(
    // Simulate user data
    ro.Just(
        User{
            ID:       1,
            Username: "alice",
            Email:    "alice@example.com",
            Settings: map[string]interface{}{"theme": "dark", "notifications": true},
        },
        User{
            ID:       2,
            Username: "bob",
            Email:    "bob@example.com",
            Settings: map[string]interface{}{"theme": "light", "notifications": false},
        },
    ),
    // Encode to gob bytes
    rogob.Encode[User](),
    // Process encoded data (e.g., store to database)
    ro.Map(func(data []byte) string {
        return "stored: " + string(data)
    }),
)

subscription := pipeline.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

## Performance Considerations

- The plugin uses Go's standard `encoding/gob` package for all operations
- Gob encoding is efficient for Go-specific data structures
- Error handling is built into both `Encode` and `Decode` operators
- Gob format is binary and more compact than text-based formats like JSON
- Gob is Go-specific and not suitable for cross-language communication
- Consider the size of your data when encoding/decoding large structures
- Gob encoding includes type information, making it self-describing 