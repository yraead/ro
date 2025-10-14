# Regexp Plugin

The Regexp plugin provides operators for regular expression operations on observables.

## Installation

```bash
go get github.com/samber/ro/plugins/regexp
```

## Operators

### Find

Finds the first match of the pattern in byte slices.

```go
import (
    "regexp"
    "github.com/samber/ro"
    roregexp "github.com/samber/ro/plugins/regexp"
)

pattern := regexp.MustCompile(`\d+`)
observable := ro.Pipe1(
    ro.Just(
        []byte("Hello 123 World"),
        []byte("Test 456 Example"),
        []byte("No numbers here"),
    ),
    roregexp.Find[[]byte](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [49 50 51]
// Next: [52 53 54]
// Next: []
// Completed
```

### FindString

Finds the first match of the pattern in strings.

```go
pattern := regexp.MustCompile(`\d+`)
observable := ro.Pipe1(
    ro.Just(
        "Hello 123 World",
        "Test 456 Example",
        "No numbers here",
    ),
    roregexp.FindString[string](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: 123
// Next: 456
// Next:
// Completed
```

### FindSubmatch

Finds the first submatch of the pattern in byte slices.

```go
pattern := regexp.MustCompile(`(\d+)-(\w+)`)
observable := ro.Pipe1(
    ro.Just(
        []byte("123-abc"),
        []byte("456-def"),
        []byte("No match"),
    ),
    roregexp.FindSubmatch[[]byte](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[[][]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [[49 50 51 45 97 98 99] [49 50 51] [97 98 99]]
// Next: [[52 53 54 45 100 101 102] [52 53 54] [100 101 102]]
// Next: []
// Completed
```

### FindStringSubmatch

Finds the first submatch of the pattern in strings.

```go
pattern := regexp.MustCompile(`(\d+)-(\w+)`)
observable := ro.Pipe1(
    ro.Just(
        "123-abc",
        "456-def",
        "No match",
    ),
    roregexp.FindStringSubmatch[string](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[[]string]())
defer subscription.Unsubscribe()

// Output:
// Next: [123-abc 123 abc]
// Next: [456-def 456 def]
// Next: []
// Completed
```

### FindAll

Finds all matches of the pattern in byte slices.

```go
pattern := regexp.MustCompile(`\d+`)
observable := ro.Pipe1(
    ro.Just(
        []byte("Hello 123 World 456"),
        []byte("Test 789 Example"),
        []byte("No numbers here"),
    ),
    roregexp.FindAll[[]byte](pattern, -1),
)

subscription := observable.Subscribe(ro.PrintObserver[[][]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [[49 50 51] [52 53 54]]
// Next: [[55 56 57]]
// Next: []
// Completed
```

### FindAllString

Finds all matches of the pattern in strings.

```go
pattern := regexp.MustCompile(`\d+`)
observable := ro.Pipe1(
    ro.Just(
        "Hello 123 World 456",
        "Test 789 Example",
        "No numbers here",
    ),
    roregexp.FindAllString[string](pattern, -1),
)

subscription := observable.Subscribe(ro.PrintObserver[[]string]())
defer subscription.Unsubscribe()

// Output:
// Next: [123 456]
// Next: [789]
// Next: []
// Completed
```

### FindAllSubmatch

Finds all submatches of the pattern in byte slices.

```go
pattern := regexp.MustCompile(`(\d+)-(\w+)`)
observable := ro.Pipe1(
    ro.Just(
        []byte("123-abc 456-def"),
        []byte("789-ghi"),
        []byte("No matches"),
    ),
    roregexp.FindAllSubmatch[[]byte](pattern, -1),
)

subscription := observable.Subscribe(ro.PrintObserver[[][][]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [[[49 50 51 45 97 98 99] [49 50 51] [97 98 99]] [[52 53 54 45 100 101 102] [52 53 54] [100 101 102]]]
// Next: [[[55 56 57 45 103 104 105] [55 56 57] [103 104 105]]]
// Next: []
// Completed
```

### FindAllStringSubmatch

Finds all submatches of the pattern in strings.

```go
pattern := regexp.MustCompile(`(\d+)-(\w+)`)
observable := ro.Pipe1(
    ro.Just(
        "123-abc 456-def",
        "789-ghi",
        "No matches",
    ),
    roregexp.FindAllStringSubmatch[string](pattern, -1),
)

subscription := observable.Subscribe(ro.PrintObserver[[][]string]())
defer subscription.Unsubscribe()

// Output:
// Next: [[123-abc 123 abc] [456-def 456 def]]
// Next: [[789-ghi 789 ghi]]
// Next: []
// Completed
```

### Match

Checks if byte slices match the pattern.

```go
pattern := regexp.MustCompile(`^\d+$`)
observable := ro.Pipe1(
    ro.Just(
        []byte("123"),
        []byte("abc"),
        []byte("456"),
    ),
    roregexp.Match[[]byte](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[bool]())
defer subscription.Unsubscribe()

// Output:
// Next: true
// Next: false
// Next: true
// Completed
```

### MatchString

Checks if strings match the pattern.

```go
pattern := regexp.MustCompile(`^\d+$`)
observable := ro.Pipe1(
    ro.Just(
        "123",
        "abc",
        "456",
    ),
    roregexp.MatchString[string](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[bool]())
defer subscription.Unsubscribe()

// Output:
// Next: true
// Next: false
// Next: true
// Completed
```

### ReplaceAll

Replaces all matches of the pattern in byte slices with the replacement.

```go
pattern := regexp.MustCompile(`\d+`)
observable := ro.Pipe1(
    ro.Just(
        []byte("Hello 123 World"),
        []byte("Test 456 Example"),
    ),
    roregexp.ReplaceAll[[]byte](pattern, []byte("XXX")),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [72 101 108 108 111 32 88 88 88 32 87 111 114 108 100]
// Next: [84 101 115 116 32 88 88 88 32 69 120 97 109 112 108 101]
// Completed
```

### ReplaceAllString

Replaces all matches of the pattern in strings with the replacement.

```go
pattern := regexp.MustCompile(`\d+`)
observable := ro.Pipe1(
    ro.Just(
        "Hello 123 World",
        "Test 456 Example",
    ),
    roregexp.ReplaceAllString[string](pattern, "XXX"),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: Hello XXX World
// Next: Test XXX Example
// Completed
```

### FilterMatch

Filters byte slices that match the pattern.

```go
pattern := regexp.MustCompile(`^\d+$`)
observable := ro.Pipe1(
    ro.Just(
        []byte("123"),
        []byte("abc"),
        []byte("456"),
        []byte("def"),
    ),
    roregexp.FilterMatch[[]byte](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
defer subscription.Unsubscribe()

// Output:
// Next: [49 50 51]
// Next: [52 53 54]
// Completed
```

### FilterMatchString

Filters strings that match the pattern.

```go
pattern := regexp.MustCompile(`^\d+$`)
observable := ro.Pipe1(
    ro.Just(
        "123",
        "abc",
        "456",
        "def",
    ),
    roregexp.FilterMatchString[string](pattern),
)

subscription := observable.Subscribe(ro.PrintObserver[string]())
defer subscription.Unsubscribe()

// Output:
// Next: 123
// Next: 456
// Completed
```

## Pattern Compilation

You can compile patterns once and reuse them:

```go
// Compile pattern once
pattern := regexp.MustCompile(`\d+`)

// Use the same pattern for multiple operations
observable1 := ro.Pipe1(
    ro.Just("Hello 123 World"),
    roregexp.FindString[string](pattern),
)

observable2 := ro.Pipe1(
    ro.Just("Test 456 Example"),
    roregexp.ReplaceAllString[string](pattern, "XXX"),
)
```

## Error Handling

The plugin uses `regexp.MustCompile` which panics on invalid patterns. For production code, use `regexp.Compile`:

```go
pattern, err := regexp.Compile(`\d+`)
if err != nil {
    // Handle compilation error
    return
}

observable := ro.Pipe1(
    ro.Just("Hello 123 World"),
    roregexp.FindString[string](pattern),
)
```

## Real-world Example

Here's a practical example that validates and extracts email addresses:

```go
import (
    "regexp"
    "github.com/samber/ro"
    roregexp "github.com/samber/ro/plugins/regexp"
)

// Email validation pattern
emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Process user data
pipeline := ro.Pipe2(
    // Source: User data
    ro.Just(
        "user1@example.com",
        "invalid-email",
        "user2@test.org",
        "not-an-email",
    ),
    // Filter valid emails
    roregexp.FilterMatchString[string](emailPattern),
)

subscription := pipeline.Subscribe(
    ro.NewObserver(
        func(email string) {
            // Process valid email
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

- Compile patterns once and reuse them
- Use appropriate pattern complexity for your use case
- Consider using `regexp.Compile` instead of `regexp.MustCompile` for production
- The plugin doesn't block the observable stream
- Regular expression operations are done synchronously
- Complex patterns may impact performance on large datasets 