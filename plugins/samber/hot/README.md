# Samber Hot Plugin

The samber/hot plugin provides operators for integrating with [samber/hot](https://github.com/samber/hot) caching library in reactive streams.

## Installation

```bash
go get github.com/samber/ro/plugins/samber/hot
```

## Operators

### GetOrFetch

Retrieves values from a hot cache and returns a tuple with the value and a boolean indicating if the value was found.

```go
import (
    "github.com/samber/ro"
    rohot "github.com/samber/ro/plugins/samber/hot"
    "github.com/samber/hot"
    "github.com/samber/lo"
)

// Create a hot cache
cache := hot.NewHotCache[string, User](hot.Config{
    Capacity: 1000,
    TTL:      5 * time.Minute,
})

observable := ro.Pipe1(
    ro.Just("user1", "user2", "user3"),
    rohot.GetOrFetch(cache),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(tuple lo.Tuple2[User, bool]) {
        user, found := tuple.Unpack()
        if found {
            fmt.Printf("Found user: %s\n", user.Name)
        } else {
            fmt.Printf("User not found in cache\n")
        }
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### GetOrFetchOrSkip

Retrieves values from a hot cache and only emits values that were found in the cache.

```go
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user3"),
    rohot.GetOrFetchOrSkip(cache),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(user User) {
        fmt.Printf("Found user: %s\n", user.Name)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))

// Output:
// Found user: John
// Found user: Jane
// Completed
```

### GetOrFetchOrError

Retrieves values from a hot cache and emits an error if the value is not found.

```go
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user3"),
    rohot.GetOrFetchOrError(cache),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(user User) {
        fmt.Printf("Found user: %s\n", user.Name)
    },
    func(err error) {
        if errors.Is(err, rohot.NotFound) {
            fmt.Printf("User not found in cache\n")
        } else {
            fmt.Printf("Cache error: %v\n", err)
        }
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### GetOrFetchMany

Retrieves multiple values from a hot cache at once, useful for batch operations.

```go
observable := ro.Pipe1(
    ro.Just([]string{"user1", "user2", "user3"}),
    rohot.GetOrFetchMany(cache),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(users map[string]User) {
        for id, user := range users {
            fmt.Printf("User %s: %s\n", id, user.Name)
        }
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

## Advanced Usage

### Cache Population

You can populate the cache and then use it in reactive streams:

```go
import (
    "github.com/samber/ro"
    rohot "github.com/samber/ro/plugins/samber/hot"
    "github.com/samber/hot"
)

// Create and populate cache
cache := hot.NewHotCache[string, User](hot.Config{
    Capacity: 1000,
    TTL:      5 * time.Minute,
})

// Populate cache with some users
cache.Set("user1", User{ID: "user1", Name: "John"})
cache.Set("user2", User{ID: "user2", Name: "Jane"})

// Use cache in reactive stream
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user3"),
    rohot.GetOrFetchOrSkip(cache),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(user User) {
        fmt.Printf("Found user: %s\n", user.Name)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Combining with Other Operators

```go
import (
    "github.com/samber/ro"
    rohot "github.com/samber/ro/plugins/samber/hot"
    "github.com/samber/hot"
)

// Create cache
cache := hot.NewHotCache[string, User](hot.Config{
    Capacity: 1000,
    TTL:      5 * time.Minute,
})

// Process user IDs, fetch from cache, and filter active users
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user3", "user4"),
    rohot.GetOrFetchOrSkip(cache),
    ro.Filter(func(user User) bool {
        return user.Active
    }),
    ro.Map(func(user User) string {
        return user.Name
    }),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(name string) {
        fmt.Printf("Active user: %s\n", name)
    },
    func(err error) {
        fmt.Printf("Error: %v\n", err)
    },
    func() {
        fmt.Println("Completed")
    },
))
```

### Error Handling

The plugin provides different error handling strategies:

1. **Skip missing**: Use `GetOrFetchOrSkip` to filter out missing values
2. **Error on missing**: Use `GetOrFetchOrError` to emit errors for missing values
3. **Tuple with status**: Use `GetOrFetch` to get both value and found status

```go
// Handle different error scenarios
observable := ro.Pipe1(
    ro.Just("user1", "user2", "user3"),
    rohot.GetOrFetch(cache),
)

subscription := observable.Subscribe(ro.NewObserver(
    func(tuple lo.Tuple2[User, bool]) {
        user, found := tuple.Unpack()
        if found {
            fmt.Printf("Found user: %s\n", user.Name)
        } else {
            fmt.Printf("User not found in cache\n")
        }
    },
    func(err error) {
        if errors.Is(err, rohot.NotFound) {
            fmt.Printf("User not found in cache\n")
        } else {
            fmt.Printf("Cache error: %v\n", err)
        }
    },
    func() {
        fmt.Println("Completed")
    },
))
```

## Performance Considerations

- Use `GetOrFetchMany` for batch operations to reduce cache lookups
- Configure appropriate TTL and capacity for your use case
- Consider using `ro.ObserveOn` to control the scheduler for cache operations
- Monitor cache hit rates and adjust cache size accordingly

## Dependencies

This plugin requires the [samber/hot](https://github.com/samber/hot) and [samber/lo](https://github.com/samber/lo) libraries:

```bash
go get github.com/samber/hot
go get github.com/samber/lo
``` 