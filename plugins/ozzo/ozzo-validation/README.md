# Ozzo Validation Plugin

The ozzo-validation plugin provides operators for validating data in reactive streams using the [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) library.

## Installation

```bash
go get github.com/samber/ro/plugins/ozzo/ozzo-validation
```

## Operators

### Validate

Validates values using ozzo validation rules and returns a `Result` monad.

```go
import (
    "github.com/samber/ro"
    roozzo "github.com/samber/ro/plugins/ozzo/ozzo-validation"
    ozzo "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

observable := ro.Pipe1(
    ro.Just(
        User{Name: "John", Email: "john@example.com", Age: 25},
        User{Name: "", Email: "invalid-email", Age: -5},
    ),
    roozzo.Validate(
        ozzo.Field(&User.Name, ozzo.Required, ozzo.Length(1, 50)),
        ozzo.Field(&User.Email, ozzo.Required, is.Email),
        ozzo.Field(&User.Age, ozzo.Required, ozzo.Min(0), ozzo.Max(150)),
    ),
)

// Output:
// Next: Ok(User{Name: "John", Email: "john@example.com", Age: 25})
// Next: Err(validation errors)
// Completed
```

### ValidateStruct

Validates structs that implement the `Validatable` interface.

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (u User) Validate() error {
    return ozzo.ValidateStruct(&u,
        ozzo.Field(&u.Name, ozzo.Required, ozzo.Length(1, 50)),
        ozzo.Field(&u.Email, ozzo.Required, is.Email),
    )
}

observable := ro.Pipe1(
    ro.Just(
        User{Name: "John", Email: "john@example.com"},
        User{Name: "", Email: "invalid-email"},
    ),
    roozzo.ValidateStruct[User](),
)

// Output:
// Next: Ok(User{Name: "John", Email: "john@example.com"})
// Next: Err(validation errors)
// Completed
```

### ValidateOrError

Validates values and propagates errors through the stream instead of wrapping them in a Result.

```go
observable := ro.Pipe1(
    ro.Just(
        User{Name: "John", Email: "john@example.com", Age: 25},
        User{Name: "", Email: "invalid-email", Age: -5},
    ),
    roozzo.ValidateOrError(
        ozzo.Field(&User.Name, ozzo.Required, ozzo.Length(1, 50)),
        ozzo.Field(&User.Email, ozzo.Required, is.Email),
        ozzo.Field(&User.Age, ozzo.Required, ozzo.Min(0), ozzo.Max(150)),
    ),
)

// Output:
// Next: User{Name: "John", Email: "john@example.com", Age: 25}
// Error: validation errors
```

### ValidateOrSkip

Validates values and skips invalid ones, only emitting valid values.

```go
observable := ro.Pipe1(
    ro.Just(
        User{Name: "John", Email: "john@example.com", Age: 25},
        User{Name: "", Email: "invalid-email", Age: -5},
        User{Name: "Jane", Email: "jane@example.com", Age: 30},
    ),
    roozzo.ValidateOrSkip(
        ozzo.Field(&User.Name, ozzo.Required, ozzo.Length(1, 50)),
        ozzo.Field(&User.Email, ozzo.Required, is.Email),
        ozzo.Field(&User.Age, ozzo.Required, ozzo.Min(0), ozzo.Max(150)),
    ),
)

// Output:
// Next: User{Name: "John", Email: "john@example.com", Age: 25}
// Next: User{Name: "Jane", Email: "jane@example.com", Age: 30}
// Completed
```

### Context-Aware Validation

All operators have context-aware variants that pass the context to the validation rules:

- `ValidateWithContext`
- `ValidateStructWithContext`
- `ValidateOrErrorWithContext`
- `ValidateOrSkipWithContext`

```go
observable := ro.Pipe1(
    ro.Just(
        User{Name: "John", Email: "john@example.com", Age: 25},
    ),
    roozzo.ValidateWithContext(
        ozzo.Field(&User.Name, ozzo.Required, ozzo.Length(1, 50)),
        ozzo.Field(&User.Email, ozzo.Required, is.Email),
        ozzo.Field(&User.Age, ozzo.Required, ozzo.Min(0), ozzo.Max(150)),
    ),
)
```

## Result Monad

The `Result` type provides a monadic interface for handling validation results:

```go
result := roozzo.Ok("valid value")
if result.IsOk() {
    value := result.Unwrap()
    // Use the validated value
}

errResult := roozzo.Err[string](errors.New("validation failed"))
if errResult.IsError() {
    err := errResult.Error()
    // Handle the validation error
}

// Get both value and error
value, err := result.Get()
```

## Error Handling

The plugin provides different error handling strategies:

1. **Result-based**: Use `Validate` or `ValidateStruct` to get a `Result` monad
2. **Error propagation**: Use `ValidateOrError` to propagate errors through the stream
3. **Skip invalid**: Use `ValidateOrSkip` to filter out invalid values
4. **Context-aware**: Use the `WithContext` variants for context-aware validation

## Dependencies

This plugin requires the [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) library:

```bash
go get github.com/go-ozzo/ozzo-validation/v4
``` 