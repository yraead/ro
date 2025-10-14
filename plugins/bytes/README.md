# Bytes Plugin

The byte plugin provides operators for manipulating byte slices and strings in reactive streams.

## Installation

```bash
go get github.com/samber/ro/plugins/bytes
```

## Operators

### CamelCase

Converts strings to camelCase format.

```go
import (
    "github.com/samber/ro"
    robytes "github.com/samber/ro/plugins/bytes"
)

observable := ro.Pipe1(
    ro.Just(
        []byte("hello world"),
        []byte("user_name"),
        []byte("API_KEY"),
    ),
    robytes.CamelCase[[]byte](),
)

// Output:
// Next: helloWorld
// Next: userName
// Next: apiKey
// Completed
```

### Capitalize

Capitalizes the first letter of each string.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("hello"),
        []byte("world"),
        []byte("golang"),
    ),
    robytes.Capitalize[[]byte](),
)

// Output:
// Next: Hello
// Next: World
// Next: Golang
// Completed
```

### KebabCase

Converts strings to kebab-case format.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("hello world"),
        []byte("userName"),
        []byte("API_KEY"),
    ),
    robytes.KebabCase[[]byte](),
)

// Output:
// Next: hello-world
// Next: user-name
// Next: api-key
// Completed
```

### PascalCase

Converts strings to PascalCase format.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("hello world"),
        []byte("user_name"),
        []byte("api_key"),
    ),
    robytes.PascalCase[[]byte](),
)

// Output:
// Next: HelloWorld
// Next: UserName
// Next: ApiKey
// Completed
```

### SnakeCase

Converts strings to snake_case format.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("hello world"),
        []byte("userName"),
        []byte("API_KEY"),
    ),
    robytes.SnakeCase[[]byte](),
)

// Output:
// Next: hello_world
// Next: user_name
// Next: api_key
// Completed
```

### Ellipsis

Truncates strings with ellipsis to a specified length.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("This is a very long string that needs to be truncated"),
        []byte("Short"),
        []byte("Another long string for demonstration"),
    ),
    robytes.Ellipsis[[]byte](20),
)

// Output:
// Next: This is a very lon...
// Next: Short
// Next: Another long string...
// Completed
```

### Words

Splits strings into words.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("hello world"),
        []byte("user_name"),
        []byte("camelCase"),
        []byte("PascalCase"),
    ),
    robytes.Words[[]byte](),
)

// Output:
// Next: [hello world]
// Next: [user name]
// Next: [camel case]
// Next: [pascal case]
// Completed
```

### Random

Generates random strings of specified size using a charset.

```go
observable := ro.Pipe1(
    ro.Just(
        []byte("prefix"),
        []byte("suffix"),
        []byte("base"),
    ),
    robytes.Random[[]byte](10, robytes.AlphanumericCharset),
)

// Output: (random strings will vary)
// Next: prefixa1b2c3d4e5
// Next: suffixf6g7h8i9j0
// Next: basek1l2m3n4o5p
// Completed
```

## Available Charsets

The Random operator provides several predefined charsets:

- `LowerCaseLettersCharset`: a-z
- `UpperCaseLettersCharset`: A-Z
- `LettersCharset`: a-z + A-Z
- `NumbersCharset`: 0-9
- `AlphanumericCharset`: a-z + A-Z + 0-9
- `SpecialCharset`: !@#$%^&*()_+-=[]{}|;':",./<>?
- `AllCharset`: All characters

## Real-world Example

Here's a practical example that processes user input and normalizes it:

```go
import (
    "github.com/samber/ro"
    robytes "github.com/samber/ro/plugins/bytes"
)

// Process user input and normalize to different formats
pipeline := ro.Pipe4(
    // Simulate user input
    ro.Just(
        []byte("user first name"),
        []byte("API_ENDPOINT_URL"),
        []byte("databaseTableName"),
    ),
    // Convert to snake_case for database
    robytes.SnakeCase[[]byte](),
    // Also generate camelCase version
    robytes.CamelCase[[]byte](),
    // And PascalCase for display
    robytes.PascalCase[[]byte](),
)

subscription := pipeline.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: user_first_name
// Next: api_endpoint_url
// Next: database_table_name
// Completed
``` 