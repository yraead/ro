# String Conversion Plugin

The string conversion plugin provides operators for converting between strings and various data types using Go's `strconv` package.

## Installation

```bash
go get github.com/samber/ro/plugins/strconv
```

## Operators

### Atoi

Converts strings to integers using `strconv.Atoi`.

```go
import (
    "github.com/samber/ro"
    rostrconv "github.com/samber/ro/plugins/strconv"
)

observable := ro.Pipe1(
    ro.Just("123", "456", "789"),
    rostrconv.Atoi[string](),
)

subscription := observable.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()

// Output:
// Next: 123
// Next: 456
// Next: 789
// Completed
```

### ParseInt

Converts strings to int64 values with specified base and bit size.

```go
observable := ro.Pipe1(
    ro.Just("123", "FF", "1010"),
    rostrconv.ParseInt[string](16, 64), // Parse as hex, 64-bit
)

subscription := observable.Subscribe(ro.PrintObserver[int64]())
defer subscription.Unsubscribe()

// Output:
// Next: 291
// Next: 255
// Next: 4112
// Completed
```

### ParseFloat

Converts strings to float64 values with specified bit size.

```go
observable := ro.Pipe1(
    ro.Just("3.14", "2.718", "1.414"),
    rostrconv.ParseFloat[string](64), // Parse as 64-bit float
)

subscription := observable.Subscribe(ro.PrintObserver[float64]())
defer subscription.Unsubscribe()

// Output:
// Next: 3.14
// Next: 2.718
// Next: 1.414
// Completed
```

### ParseBool

Converts strings to boolean values using `strconv.ParseBool`.

```go
observable := ro.Pipe1(
    ro.Just("true", "false", "1", "0"),
    rostrconv.ParseBool[string](),
)

subscription := observable.Subscribe(ro.PrintObserver[bool]())
defer subscription.Unsubscribe()

// Output:
// Next: true
// Next: false
// Next: true
// Next: false
// Completed
```

### ParseUint

Converts strings to uint64 values with specified base and bit size.

```go
observable := ro.Pipe1(
    ro.Just("123", "456", "789"),
    rostrconv.ParseUint[string](10, 64), // Parse as decimal, 64-bit unsigned
)

subscription := observable.Subscribe(ro.PrintObserver[uint64]())
defer subscription.Unsubscribe()

// Output:
// Next: 123
// Next: 456
// Next: 789
// Completed
```

### FormatBool

Converts boolean values to strings using `strconv.FormatBool`.

```go
observable := ro.Pipe1(
    ro.Just(true, false, true),
    rostrconv.FormatBool(),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: true
// Next: false
// Next: true
// Completed
```

### FormatFloat

Converts float64 values to strings with specified format, precision, and bit size.

```go
observable := ro.Pipe1(
    ro.Just(3.14159, 2.71828, 1.41421),
    rostrconv.FormatFloat('f', 3, 64), // Format with 3 decimal places
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: 3.142
// Next: 2.718
// Next: 1.414
// Completed
```

### FormatInt

Converts int64 values to strings with specified base.

```go
observable := ro.Pipe1(
    ro.Just(int64(255), int64(123), int64(456)),
    rostrconv.FormatInt[string](16), // Format as hex
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: ff
// Next: 7b
// Next: 1c8
// Completed
```

### FormatUint

Converts uint64 values to strings with specified base.

```go
observable := ro.Pipe1(
    ro.Just(uint64(255), uint64(123), uint64(456)),
    rostrconv.FormatUint[string](10), // Format as decimal
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: 255
// Next: 123
// Next: 456
// Completed
```

### Itoa

Converts integers to strings using `strconv.Itoa`.

```go
observable := ro.Pipe1(
    ro.Just(123, 456, 789),
    rostrconv.Itoa(),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: 123
// Next: 456
// Next: 789
// Completed
```

### Quote

Quotes strings using `strconv.Quote`.

```go
observable := ro.Pipe1(
    ro.Just("hello", "world", "golang"),
    rostrconv.Quote(),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: "hello"
// Next: "world"
// Next: "golang"
// Completed
```

### QuoteRune

Quotes runes using `strconv.QuoteRune`.

```go
observable := ro.Pipe1(
    ro.Just('a', 'b', 'c'),
    rostrconv.QuoteRune(),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: 'a'
// Next: 'b'
// Next: 'c'
// Completed
```

### Unquote

Unquotes strings using `strconv.Unquote`.

```go
observable := ro.Pipe1(
    ro.Just(`"hello"`, `"world"`, `"golang"`),
    rostrconv.Unquote(),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: hello
// Next: world
// Next: golang
// Completed
```

## Error Handling

All parsing operators handle errors gracefully and will emit error notifications for invalid input:

```go
observable := ro.Pipe1(
    ro.Just("123", "abc", "456"), // "abc" will cause an error
    rostrconv.ParseInt[string](10, 64),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value int64) {
            // Handle successful parsing
        },
        func(err error) {
            // Handle parsing error
        },
        func() {
            // Handle completion
        },
    ),
)
defer subscription.Unsubscribe()
```

## Roundtrip Examples

Demonstrate roundtrip conversions:

```go
// Integer roundtrip
observable := ro.Pipe2(
    ro.Just(123, 456, 789),
    rostrconv.Itoa(),
    rostrconv.Atoi[string](),
)

subscription := observable.Subscribe(ro.PrintObserver[int]())
defer subscription.Unsubscribe()

// Output:
// Next: 123
// Next: 456
// Next: 789
// Completed
```

## Real-world Example

Here's a practical example that processes CSV data with type conversions:

```go
import (
    "github.com/samber/ro"
    rostrconv "github.com/samber/ro/plugins/strconv"
)

type User struct {
    ID   int
    Name string
    Age  int
}

// Process CSV data with type conversions
pipeline := ro.Pipe4(
    // Simulate CSV data as strings
    ro.Just(
        []string{"1", "Alice", "30"},
        []string{"2", "Bob", "25"},
        []string{"3", "Charlie", "35"},
    ),
    // Convert ID and Age to integers
    ro.Map(func(row []string) []string {
        // In a real scenario, you'd use rostrconv.Atoi for each field
        return row
    }),
    // Extract and convert ID
    ro.Map(func(row []string) int {
        // This would be done with rostrconv.Atoi in practice
        return 1 // Simplified for example
    }),
    // Convert back to string for output
    rostrconv.Itoa(),
)

subscription := pipeline.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()
```

## Base and Bit Size Parameters

### ParseInt and ParseUint

- **base**: Number system (10 for decimal, 16 for hex, 8 for octal, 2 for binary)
- **bitSize**: Integer type size (32 for int32/uint32, 64 for int64/uint64)

### ParseFloat

- **bitSize**: Float type size (32 for float32, 64 for float64)

### FormatFloat

- **fmt**: Format byte ('f', 'e', 'E', 'g', 'G')
- **prec**: Precision (number of decimal places)
- **bitSize**: Float type size (32 or 64)

## Performance Considerations

- The plugin uses Go's standard `strconv` package for all conversions
- Error handling is built into all parsing operators
- Consider using appropriate bit sizes to avoid overflow
- Use base 10 for decimal numbers, base 16 for hexadecimal, etc. 