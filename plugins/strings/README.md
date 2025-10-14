# Strings Plugin

The string plugin provides operators for manipulating strings in reactive streams.

## Installation

```bash
go get github.com/samber/ro/plugins/strings
```

## Operators

### CamelCase

Converts strings to camelCase format.

```go
import (
    "github.com/samber/ro"
    rostrings "github.com/samber/ro/plugins/strings"
)

observable := ro.Pipe1(
    ro.Just(
        "hello world",
        "user_name",
        "API_KEY",
    ),
    rostrings.CamelCase[string](),
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
        "hello",
        "world",
        "golang",
    ),
    rostrings.Capitalize[string](),
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
        "hello world",
        "userName",
        "API_KEY",
    ),
    rostrings.KebabCase[string](),
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
        "hello world",
        "user_name",
        "api_key",
    ),
    rostrings.PascalCase[string](),
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
        "hello world",
        "userName",
        "API_KEY",
    ),
    rostrings.SnakeCase[string](),
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
        "This is a very long string that needs to be truncated",
        "Short",
        "Another long string for demonstration",
    ),
    rostrings.Ellipsis[string](20),
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
        "hello world",
        "user_name",
        "camelCase",
        "PascalCase",
    ),
    rostrings.Words[string](),
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
        "prefix",
        "suffix",
        "base",
    ),
    rostrings.Random[string](10, rostrings.AlphanumericCharset),
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
    rostrings "github.com/samber/ro/plugins/strings"
)

// Process user input and normalize to different formats
pipeline := ro.Pipe4(
    // Simulate user input
    ro.Just(
        "user first name",
        "API_ENDPOINT_URL",
        "databaseTableName",
    ),
    // Convert to snake_case for database
    rostrings.SnakeCase[string](),
    // Also generate camelCase version
    rostrings.CamelCase[string](),
    // And PascalCase for display
    rostrings.PascalCase[string](),
)

subscription := pipeline.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: user_first_name
// Next: api_endpoint_url
// Next: database_table_name
// Completed
```

## Comparison with Byte Plugin

The string plugin provides similar functionality to the byte plugin but works with `string` types instead of `[]byte`. Choose the appropriate plugin based on your data type:

- Use **string plugin** when working with `string` types
- Use **byte plugin** when working with `[]byte` types

## Performance Considerations

- String operations are generally faster than byte operations for text processing
- The string plugin uses Go's standard `strings` package for operations
- Consider using the byte plugin for binary data or when you need byte-level control 