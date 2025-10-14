# ro Iter Plugin

This plugin requires Go 1.23 or later due to the use of the new `iter` package.

This plugin provides seamless integration between Go 1.23's [iterators](https://pkg.go.dev/iter) and the [ro](https://github.com/samber/ro) reactive programming library. It allows you to:

- Convert Go iterators (`iter.Seq` and `iter.Seq2`) to Ro observables
- Convert Ro observables back to Go iterators
- Bridge the gap between traditional iteration and reactive streams

## Installation

```bash
go get github.com/samber/ro/plugins/iter
```

## Requirements

- Go 1.23 or later
- [samber/ro](https://github.com/samber/ro) reactive programming library

## Quick Start

```go
package main

import (
	"fmt"
	"iter"

	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Convert an iterator to an observable
	seq := func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			if !yield(i) {
				return
			}
		}
	}

	observable := roiter.FromSeq(seq)

	// Subscribe and process values
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(value int) {
				fmt.Printf("Received: %d\n", value)
			},
			func(err error) {
				fmt.Printf("Error: %v\n", err)
			},
			func() {
				fmt.Println("Completed")
			},
		),
	)
	defer subscription.Unsubscribe()
}
```

## API Reference

### From Iterators to Observables

#### `FromSeq[T any](iterator iter.Seq[T]) ro.Observable[T]`

Converts a single-value iterator sequence to an observable.

```go
// Create a simple sequence
seq := func(yield func(string) bool) {
	words := []string{"hello", "world", "reactive"}
	for _, word := range words {
		if !yield(word) {
			return
		}
	}
}

observable := roiter.FromSeq(seq)
```

#### `FromSeq2[K, V any](iterator iter.Seq2[K, V]) ro.Observable[lo.Tuple2[K, V]]`

Converts a key-value iterator sequence to an observable that emits tuples.

```go
// Create a key-value sequence
seq := func(yield func(string, int) bool) {
	pairs := map[string]int{"apple": 5, "banana": 3, "orange": 8}
	for k, v := range pairs {
		if !yield(k, v) {
			return
		}
	}
}

observable := roiter.FromSeq2(seq)

// Handle tuple values
subscription := observable.Subscribe(
	ro.NewObserver(func(pair lo.Tuple2[string, int]) {
		fmt.Printf("%s: %d\n", pair.A, pair.B)
	}, nil, nil),
)
```

### From Observables to Iterators

#### `ToSeq[T any](source ro.Observable[T]) iter.Seq[T]`

Converts an observable back to a single-value iterator sequence.

```go
observable := ro.Just(1, 2, 3, 4, 5)
seq := roiter.ToSeq(observable)

// Iterate using range syntax
for value := range seq {
	fmt.Printf("Value: %d\n", value)
}
```

#### `ToSeq2[T any](source ro.Observable[T]) iter.Seq2[int, T]`

Converts an observable to a key-value iterator with automatic indexing.

```go
observable := ro.Just("apple", "banana", "orange")
seq := roiter.ToSeq2(observable)

// Iterate with indices
for index, value := range seq {
	fmt.Printf("Index %d: %s\n", index, value)
}
```

## Usage Examples

### Basic Iterator to Observable

```go
package main

import (
	"fmt"
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Create a number generator
	numbers := func(yield func(int) bool) {
		for i := 1; i <= 10; i++ {
			if !yield(i) {
				return
			}
		}
	}

	// Convert to observable
	observable := roiter.FromSeq(numbers)

	// Apply reactive transformations
	result := ro.Pipe1(
		observable,
		ro.Filter(func(n int) bool {
			return n%2 == 0 // Only even numbers
		}),
		ro.Map(func(n int) int {
			return n * n // Square the numbers
		}),
	)

	// Subscribe to results
	subscription := result.Subscribe(
		ro.NewObserver(
			func(value int) {
				fmt.Printf("Result: %d\n", value)
			},
			nil,
			func() {
				fmt.Println("Stream completed")
			},
		),
	)
	defer subscription.Unsubscribe()
}
```

### Observable to Iterator for Traditional Processing

```go
package main

import (
	"fmt"
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Create a reactive stream
	observable := ro.Pipe1(
		ro.Just("hello", "world", "reactive", "programming"),
		ro.Map(func(s string) string {
			return fmt.Sprintf("🔥 %s 🔥", s)
		}),
	)

	// Convert back to iterator for traditional processing
	seq := roiter.ToSeq(observable)

	// Use standard Go iteration
	for value := range seq {
		fmt.Printf("Processing: %s\n", value)
	}
}
```

### Combining with Other Ro Operators

```go
package main

import (
	"fmt"
	"time"

	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Create a data source from iterator
	data := func(yield func(string) bool) {
		items := []string{"task1", "task2", "task3", "task4", "task5"}
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}

	// Convert to observable and apply reactive operators
	processed := ro.Pipe3(
		roiter.FromSeq(data),
		ro.Map(func(task string) string {
			// Simulate processing
			return fmt.Sprintf("processed-%s", task)
		}),
		ro.Delay[string](100*time.Millisecond),
		ro.Take[string](3), // Only take first 3 items
	)

	// Convert back to iterator for final processing
	results := roiter.ToSeq(processed)

	// Process results
	for result := range results {
		fmt.Printf("Final result: %s\n", result)
	}
}
```

### Working with Key-Value Data

```go
package main

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Create key-value iterator
	kvData := func(yield func(string, int) bool) {
		products := []struct {
			name  string
			price int
		}{
			{"laptop", 1200},
			{"mouse", 25},
			{"keyboard", 75},
			{"monitor", 300},
		}

		for _, product := range products {
			if !yield(product.name, product.price) {
				return
			}
		}
	}

	// Convert to observable of tuples
	observable := roiter.FromSeq2(kvData)

	// Process with reactive operators
	filtered := ro.Pipe1(
		observable,
		ro.Filter(func(pair lo.Tuple2[string, int]) bool {
			return pair.B > 50 // Filter by price > 50
		}),
	)

	// Convert back to iterator for final processing
	results := roiter.ToSeq(filtered)

	// Process filtered results
	for pair := range results {
		fmt.Printf("Product: %s, Price: $%d\n", pair.A, pair.B)
	}
}
```

### Error Handling

```go
package main

import (
	"errors"
	"fmt"

	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Create an iterator that might fail
	riskyData := func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			if i == 4 {
				// Simulate an error
				return
			}
			if !yield(i) {
				return
			}
		}
	}

	observable := roiter.FromSeq(riskyData)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(value int) {
				fmt.Printf("Received: %d\n", value)
			},
			func(err error) {
				fmt.Printf("Error: %v\n", err)
			},
			func() {
				fmt.Println("Completed successfully")
			},
		),
	)
	defer subscription.Unsubscribe()
}
```

### Early Termination

```go
package main

import (
	"fmt"

	"github.com/samber/ro"
	"github.com/samber/ro/plugins/iter"
)

func main() {
	// Create a large data source
	data := func(yield func(int) bool) {
		for i := 1; i <= 1000; i++ {
			if !yield(i) {
				fmt.Printf("Iterator stopped at %d\n", i)
				return
			}
		}
	}

	observable := roiter.FromSeq(data)

	// Convert to iterator and stop early
	seq := roiter.ToSeq(observable)

	count := 0
	for value := range seq {
		fmt.Printf("Processing: %d\n", value)
		count++

		// Stop after processing 5 items
		if count >= 5 {
			break
		}
	}

	fmt.Printf("Processed %d items total\n", count)
}
```

## Performance Considerations

1. **Backpressure**: The iterator-to-observable conversion respects Go iterator backpressure naturally through the `yield` function mechanism.

2. **Memory Usage**: Converting observables to iterators uses buffered channels (size 1) to handle the async nature of observables.

3. **Cancellation**: The `ToSeq` and `ToSeq2` functions handle context cancellation automatically when the iterator stops.

## Best Practices

1. **Use FromSeq** when you have existing Go iterators that want to benefit from reactive operators
2. **Use ToSeq** when you need to integrate reactive streams with traditional Go iteration patterns
3. **Combine with other Ro operators** for powerful data processing pipelines
4. **Handle early termination** properly to avoid resource leaks
5. **Consider error handling** in your iterators and reactive streams

## License

Apache 2.0 - See [LICENSE](https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md) for details.