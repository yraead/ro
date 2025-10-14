# Testify Plugin

The testify plugin provides testing utilities for reactive streams using the [testify](https://github.com/stretchr/testify) assertion library.

## Installation

```bash
go get github.com/samber/ro/plugins/testify
```

## Usage

### Basic Testing

```go
import (
    "testing"
    "github.com/samber/ro"
    rotestify "github.com/samber/ro/plugins/testify"
    "github.com/stretchr/testify/assert"
)

func TestObservable(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Just(1, 2, 3, 4, 5)
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectNext(1).
        ExpectNext(2).
        ExpectNext(3).
        ExpectNext(4).
        ExpectNext(5).
        ExpectComplete().
        Verify()
}
```

### Testing with Error Handling

```go
func TestObservableWithError(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Pipe1(
        ro.Just(1, 2, 3),
        ro.MapErr(func(n int) (int, error) {
            if n == 2 {
                return n, errors.New("error on 2")
            }
            return n, nil
        }),
    )
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectNext(1).
        ExpectError(errors.New("error on 2")).
        Verify()
}
```

### Testing Sequences

```go
func TestObservableSequence(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Just("a", "b", "c")
    
    rotestify.Testify[string](is).
        Source(observable).
        ExpectNextSeq("a", "b", "c").
        ExpectComplete().
        Verify()
}
```

### Testing with Custom Messages

```go
func TestObservableWithMessages(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Just(42)
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectNext(42, "expected the answer to life").
        ExpectComplete("should complete after one value").
        Verify()
}
```

## Advanced Usage

### Testing Filtered Streams

```go
func TestFilteredObservable(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Pipe1(
        ro.Just(1, 2, 3, 4, 5, 6),
        ro.Filter(func(n int) bool {
            return n%2 == 0 // Keep only even numbers
        }),
    )
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectNextSeq(2, 4, 6).
        ExpectComplete().
        Verify()
}
```

### Testing Mapped Streams

```go
func TestMappedObservable(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Pipe1(
        ro.Just("hello", "world"),
        ro.Map(func(s string) string {
            return strings.ToUpper(s)
        }),
    )
    
    rotestify.Testify[string](is).
        Source(observable).
        ExpectNextSeq("HELLO", "WORLD").
        ExpectComplete().
        Verify()
}
```

### Testing Error Scenarios

```go
func TestErrorScenarios(t *testing.T) {
    is := assert.New(t)
    
    t.Run("immediate error", func(t *testing.T) {
        observable := ro.Throw[int](errors.New("test error"))
        
        rotestify.Testify[int](is).
            Source(observable).
            ExpectError(errors.New("test error")).
            Verify()
    })
    
    t.Run("error after values", func(t *testing.T) {
        observable := ro.Pipe1(
            ro.Just(1, 2, 3),
            ro.MapErr(func(n int) (int, error) {
                if n == 3 {
                    return n, errors.New("error on 3")
                }
                return n, nil
            }),
        )
        
        rotestify.Testify[int](is).
            Source(observable).
            ExpectNext(1).
            ExpectNext(2).
            ExpectError(errors.New("error on 3")).
            Verify()
    })
}
```

### Testing with Context

```go
func TestObservableWithContext(t *testing.T) {
    is := assert.New(t)
    ctx := context.Background()
    
    observable := ro.Just(1, 2, 3)
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectNextSeq(1, 2, 3).
        ExpectComplete().
        VerifyWithContext(ctx)
}
```

## Testing Patterns

### Testing Empty Streams

```go
func TestEmptyObservable(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Empty[int]()
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectComplete().
        Verify()
}
```

### Testing Single Value Streams

```go
func TestSingleValueObservable(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Just(42)
    
    rotestify.Testify[int](is).
        Source(observable).
        ExpectNext(42).
        ExpectComplete().
        Verify()
}
```

### Testing Infinite Streams

```go
func TestInfiniteObservable(t *testing.T) {
    is := assert.New(t)
    
    observable := ro.Interval(100 * time.Millisecond)
    
    rotestify.Testify[ro.IntervalValue](is).
        Source(observable).
        ExpectNext(ro.IntervalValue(0)).
        ExpectNext(ro.IntervalValue(1)).
        ExpectNext(ro.IntervalValue(2)).
        // Note: Don't call ExpectComplete() for infinite streams
        Verify()
}
```

## Best Practices

### Use Descriptive Test Names

```go
func TestUserService_GetActiveUsers_ReturnsOnlyActiveUsers(t *testing.T) {
    is := assert.New(t)
    
    observable := userService.GetActiveUsers()
    
    rotestify.Testify[User](is).
        Source(observable).
        ExpectNextSeq(
            User{ID: "1", Name: "Alice", Active: true},
            User{ID: "2", Name: "Bob", Active: true},
        ).
        ExpectComplete().
        Verify()
}
```

### Test Error Conditions

```go
func TestUserService_GetUser_ReturnsErrorForInvalidID(t *testing.T) {
    is := assert.New(t)
    
    observable := userService.GetUser("invalid-id")
    
    rotestify.Testify[User](is).
        Source(observable).
        ExpectError(errors.New("user not found")).
        Verify()
}
```

### Test Edge Cases

```go
func TestDataProcessor_ProcessData_HandlesEmptyInput(t *testing.T) {
    is := assert.New(t)
    
    observable := dataProcessor.ProcessData([]string{})
    
    rotestify.Testify[ProcessedData](is).
        Source(observable).
        ExpectComplete().
        Verify()
}
```

## Integration with Other Testing Libraries

### Using with Testify Suite

```go
import (
    "github.com/stretchr/testify/suite"
)

type ObservableTestSuite struct {
    suite.Suite
}

func (suite *ObservableTestSuite) TestBasicObservable() {
    observable := ro.Just(1, 2, 3)
    
    rotestify.Testify[int](suite.Assert()).
        Source(observable).
        ExpectNextSeq(1, 2, 3).
        ExpectComplete().
        Verify()
}

func TestObservableTestSuite(t *testing.T) {
    suite.Run(t, new(ObservableTestSuite))
}
```

## Dependencies

This plugin requires the [testify](https://github.com/stretchr/testify) library:

```bash
go get github.com/stretchr/testify
```

## Limitations

- The testing framework is designed for synchronous testing
- Infinite streams should not call `ExpectComplete()`
- Error testing requires exact error matching
- Context-aware testing is supported but not required 