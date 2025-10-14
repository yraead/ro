# Base64 Encoding Plugin

The base64 encoding plugin provides operators for encoding and decoding data using Go's `encoding/base64` package.

## Installation

```bash
go get github.com/samber/ro/plugins/encoding/base64
```

## Operators

### Encode

Encodes byte slices to base64 strings using the specified encoding.

```go
import (
    "encoding/base64"
    "github.com/samber/ro"
    robase64 "github.com/samber/ro/plugins/encoding/base64"
)

observable := ro.Pipe1(
    ro.Just([]byte("hello"), []byte("world"), []byte("golang")),
    robase64.Encode[[]byte](base64.StdEncoding),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: aGVsbG8=
// Next: d29ybGQ=
// Next: Z29sYW5n
// Completed
```

### Decode

Decodes base64 strings to byte slices using the specified encoding.

```go
observable := ro.Pipe1(
    ro.Just("aGVsbG8=", "d29ybGQ=", "Z29sYW5n"),
    robase64.Decode[string](base64.StdEncoding),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [104 101 108 108 111]
// Next: [119 111 114 108 100]
// Next: [103 111 108 97 110 103]
// Completed
```

## Encoding Types

### Standard Encoding

Uses `base64.StdEncoding` for standard base64 encoding:

```go
observable := ro.Pipe1(
    ro.Just([]byte("hello world")),
    robase64.Encode[[]byte](base64.StdEncoding),
)
```

### URL-Safe Encoding

Uses `base64.URLEncoding` for URL-safe base64 encoding:

```go
observable := ro.Pipe1(
    ro.Just([]byte("hello world")),
    robase64.Encode[[]byte](base64.URLEncoding),
)
```

### Raw Standard Encoding

Uses `base64.RawStdEncoding` for standard base64 encoding without padding:

```go
observable := ro.Pipe1(
    ro.Just([]byte("hello world")),
    robase64.Encode[[]byte](base64.RawStdEncoding),
)
```

### Raw URL-Safe Encoding

Uses `base64.RawURLEncoding` for URL-safe base64 encoding without padding:

```go
observable := ro.Pipe1(
    ro.Just([]byte("hello world")),
    robase64.Encode[[]byte](base64.RawURLEncoding),
)
```

## Error Handling

The `Decode` operator handles errors gracefully and will emit error notifications for invalid base64 input:

```go
observable := ro.Pipe1(
    ro.Just("aGVsbG8=", "invalid-base64", "d29ybGQ="),
    robase64.Decode[string](base64.StdEncoding),
)

subscription := observable.Subscribe(
    ro.NewObserver(
        func(value []byte) {
            // Handle successful decoding
        },
        func(err error) {
            // Handle decoding error
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
// Standard encoding roundtrip
observable := ro.Pipe2(
    ro.Just([]byte("hello world")),
    robase64.Encode[[]byte](base64.StdEncoding),
    robase64.Decode[string](base64.StdEncoding),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [104 101 108 108 111 32 119 111 114 108 100]
// Completed
```

## Real-world Example

Here's a practical example that processes binary data with base64 encoding:

```go
import (
    "encoding/base64"
    "github.com/samber/ro"
    robase64 "github.com/samber/ro/plugins/encoding/base64"
)

// Process binary data with base64 encoding
pipeline := ro.Pipe3(
    // Simulate binary data
    ro.Just(
        []byte("user:password"),
        []byte("api:key123"),
        []byte("token:secret"),
    ),
    // Encode as base64
    robase64.Encode[[]byte](base64.StdEncoding),
    // Process encoded strings
    ro.Map(func(encoded string) string {
        return "Basic " + encoded
    }),
)

subscription := pipeline.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: Basic dXNlcjpwYXNzd29yZA==
// Next: Basic YXBpOmtleTEyMw==
// Next: Basic dG9rZW46c2VjcmV0
// Completed
```

## Performance Considerations

- The plugin uses Go's standard `encoding/base64` package for all operations
- Error handling is built into the `Decode` operator
- Choose the appropriate encoding type for your use case:
  - `StdEncoding` for general base64 encoding
  - `URLEncoding` for URL-safe encoding
  - `RawStdEncoding` or `RawURLEncoding` to avoid padding characters
- Consider the size of your data when encoding/decoding large byte slices 