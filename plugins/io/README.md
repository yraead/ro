# IO Plugin

The IO plugin provides operators for reading from and writing to various input/output streams.

## Installation

```bash
go get github.com/samber/ro/plugins/io
```

## Operators

### NewIOReader

Creates an observable that reads data from an `io.Reader`.

```go
import (
    "strings"
    "github.com/samber/ro"
    roio "github.com/samber/ro/plugins/io"
)

// Read data from a string reader
data := "Hello, World! This is a test."
reader := strings.NewReader(data)
observable := roio.NewIOReader(reader)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [72 101 108 108 111 44 32 87 111 114 108 100 33 32 84 104 105 115 32 105 115 32 97 32 116 101 115 116 46]
// Completed
```

### NewIOReaderLine

Creates an observable that reads lines from an `io.Reader`.

```go
// Read lines from a string reader
data := "Line 1\nLine 2\nLine 3\n"
reader := strings.NewReader(data)
observable := roio.NewIOReaderLine(reader)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [76 105 110 101 32 49]
// Next: [76 105 110 101 32 50]
// Next: [76 105 110 101 32 51]
// Completed
```

### NewIOWriter

Creates an operator that writes data to an `io.Writer` and returns the number of bytes written.

```go
import (
    "bytes"
    "github.com/samber/ro"
    roio "github.com/samber/ro/plugins/io"
)

// Write data to a buffer
var buf bytes.Buffer
writer := &buf

data := ro.Just(
    []byte("Hello, "),
    []byte("World!"),
    []byte(" This is a test."),
)

observable := ro.Pipe1(
    data,
    roio.NewIOWriter(writer),
)

subscription := observable.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()

// Output:
// Next: 32
// Completed
```

### NewStdReader

Creates an observable that reads from standard input.

```go
observable := roio.NewStdReader()

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value []byte) {
            // Handle data from stdin
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

### NewStdReaderLine

Creates an observable that reads lines from standard input.

```go
observable := roio.NewStdReaderLine()

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value []byte) {
            // Handle line from stdin
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

### NewPrompt

Creates an observable that prompts the user for input and reads their response.

```go
observable := roio.NewPrompt("Enter your name: ")

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value []byte) {
            // Handle user input
            fmt.Printf("You entered: %s\n", string(value))
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

### NewStdWriter

Creates an operator that writes data to standard output.

```go
data := ro.Just(
    []byte("Hello, "),
    []byte("World!"),
    []byte(" This is a test."),
)

observable := ro.Pipe1(
    data,
    roio.NewStdWriter(),
)

subscription := observable.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()

// Output:
// Next: 32
// Completed
```

## Supported Reader Types

The plugin supports various `io.Reader` implementations:

### String Reader

```go
import "strings"

reader := strings.NewReader("Hello, World!")
observable := roio.NewIOReader(reader)
```

### File Reader

```go
import "os"

file, err := os.Open("data.txt")
if err != nil {
    // Handle error
}
defer file.Close()

observable := roio.NewIOReader(file)
```

### Buffer Reader

```go
import "bytes"

data := []byte("Hello, World!")
reader := bytes.NewReader(data)
observable := roio.NewIOReader(reader)
```

## Supported Writer Types

The plugin supports various `io.Writer` implementations:

### Buffer Writer

```go
import "bytes"

var buf bytes.Buffer
writer := &buf

data := ro.Just([]byte("Hello, World!"))
observable := ro.Pipe1(
    data,
    roio.NewIOWriter(writer),
)
```

### File Writer

```go
import "os"

file, err := os.Create("output.txt")
if err != nil {
    // Handle error
}
defer file.Close()

data := ro.Just([]byte("Hello, World!"))
observable := ro.Pipe1(
    data,
    roio.NewIOWriter(file),
)
```

## Error Handling

The plugin handles IO errors gracefully:

### Reading Errors

```go
reader := strings.NewReader("Hello, World!")
observable := roio.NewIOReader(reader)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value []byte) {
            // Handle successful read
        },
        func(err error) {
            // Handle read error
            // This could be due to:
            // - Network errors
            // - File system errors
            // - Permission errors
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

### Writing Errors

```go
var buf bytes.Buffer
writer := &buf

data := ro.Just([]byte("Hello, World!"))
observable := ro.Pipe1(
    data,
    roio.NewIOWriter(writer),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value int) {
            // Handle successful write count
        },
        func(err error) {
            // Handle write error
            // This could be due to:
            // - Disk full
            // - Permission errors
            // - Network errors
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Real-world Example

Here's a practical example that processes a file line by line:

```go
import (
    "bytes"
    "os"
    "strings"
    "github.com/samber/ro"
    roio "github.com/samber/ro/plugins/io"
)

// Process a file line by line
pipeline := ro.Pipe2(
    // Read lines from file
    roio.NewIOReaderLine(strings.NewReader(`Line 1: Hello
Line 2: World
Line 3: Test`)),
    // Transform lines
    ro.Map(func(line []byte) []byte {
        // Convert to uppercase
        return bytes.ToUpper(line)
    }),
)

subscription := pipeline.Subscribe(
    ro.NewObserver(
        func(line []byte) {
            // Process transformed line
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

## Performance Considerations

- The plugin uses Go's standard `io` package for all operations
- Reading is done in chunks of 1024 bytes by default
- Line reading uses buffered I/O for efficiency
- The plugin automatically handles resource cleanup
- Consider the size of your data when reading/writing large files
- Use appropriate buffer sizes for your use case
- The plugin supports streaming for large files
- Context cancellation properly stops IO operations 