// Copyright 2025 samber.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ro

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleMap_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 4
	// Next: 6
	// Next: 8
	// Next: 10
	// Completed
}

func ExampleMap_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Map(func(x int) int {
			return x * 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 4
	// Next: 6
	// Error: assert.AnError general error for testing
}

func ExampleMapTo_ok() {
	observable := Pipe2(
		Just(1, 2, 3, 4, 5),
		MapTo[int, string]("Hey!"),
		Take[string](3),
	)

	subscription := observable.Subscribe(PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Hey!
	// Next: Hey!
	// Next: Hey!
	// Completed
}

func ExampleMapTo_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		MapTo[int, string]("Hey!"),
	)

	subscription := observable.Subscribe(PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Hey!
	// Next: Hey!
	// Next: Hey!
	// Error: assert.AnError general error for testing
}

func ExampleMapErr_ok() {
	observable := Pipe1(
		Just(1, 2, 3),
		MapErr(func(item int) (string, error) {
			return "Hey!", nil
		}),
	)

	subscription := observable.Subscribe(PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Hey!
	// Next: Hey!
	// Next: Hey!
	// Completed
}

func ExampleMapErr_error() {
	observable1 := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		MapErr[int, string](func(item int) (string, error) {
			return "Hey!", nil
		}),
	)

	subscription1 := observable1.Subscribe(PrintObserver[string]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Just(1, 2, 3, 4, 5),
		MapErr[int, string](func(item int) (string, error) {
			if item == 2 {
				return "Hey!", assert.AnError
			}

			return "Hey!", nil
		}),
	)

	subscription2 := observable2.Subscribe(PrintObserver[string]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: Hey!
	// Next: Hey!
	// Next: Hey!
	// Error: assert.AnError general error for testing
	// Next: Hey!
	// Error: assert.AnError general error for testing
}

func ExampleFlatMap_ok() {
	observable := Pipe1(
		Just(1, 2, 3),
		FlatMap[int](func(item int) Observable[int] {
			return Just(item, item)
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 1
	// Next: 2
	// Next: 2
	// Next: 3
	// Next: 3
	// Completed
}

func ExampleFlatMap_error() {
	observable1 := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		FlatMap[int](func(item int) Observable[int] {
			return Just(item, item)
		}),
	)

	subscription1 := observable1.Subscribe(PrintObserver[int]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Just(1, 2, 3),
		FlatMap[int](func(item int) Observable[int] {
			if item == 2 {
				return Throw[int](assert.AnError)
			}

			return Just(item, item)
		}),
	)

	subscription2 := observable2.Subscribe(PrintObserver[int]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 1
	// Next: 2
	// Next: 2
	// Next: 3
	// Next: 3
	// Error: assert.AnError general error for testing
	// Next: 1
	// Next: 1
	// Error: assert.AnError general error for testing
}

func ExampleScan_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Scan(func(agg, current int) int {
			return agg + current
		}, 42),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 43
	// Next: 45
	// Next: 48
	// Next: 52
	// Next: 57
	// Completed
}

func ExampleScan_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Scan(func(agg, current int) int {
			return agg + current
		}, 42),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 43
	// Next: 45
	// Next: 48
	// Error: assert.AnError general error for testing
}

func ExampleGroupBy_ok() {
	odd := func(v int64) bool { return v%2 == 0 }

	observable := Pipe2(
		RangeWithInterval(1, 5, 10*time.Millisecond),
		GroupBy(odd),
		MergeAll[int64](),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Completed
}

func ExampleGroupBy_error() {
	odd := func(v int) bool { return v%2 == 0 }

	observable := Pipe2(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			time.Sleep(5 * time.Millisecond)
			observer.Next(2)
			time.Sleep(5 * time.Millisecond)
			observer.Next(3)
			time.Sleep(5 * time.Millisecond)
			observer.Error(assert.AnError)
			time.Sleep(5 * time.Millisecond)
			observer.Next(4)

			return nil
		}),
		GroupBy(odd),
		MergeAll[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleBufferWhen_ok() {
	// @TODO: Implement
}

func ExampleBufferWhen_error() {
	// @TODO: Implement
}

func ExampleBufferWithTimeOrCount_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		BufferWithTimeOrCount[int](2, 100*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())

	time.Sleep(10 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: [1 2]
	// Next: [3 4]
	// Next: [5]
	// Completed
}

func ExampleBufferWithTimeOrCount_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			go func() {
				observer.Next(1)
				observer.Next(2)
				observer.Next(3)
				observer.Error(assert.AnError)
				observer.Next(4)
			}()

			return nil
		}),
		BufferWithTimeOrCount[int](2, 100*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())

	time.Sleep(10 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: [1 2]
	// Error: assert.AnError general error for testing
}

func ExampleBufferWithCount_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		BufferWithCount[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())

	time.Sleep(10 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: [1 2]
	// Next: [3 4]
	// Next: [5]
	// Completed
}

func ExampleBufferWithCount_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			go func() {
				observer.Next(1)
				observer.Next(2)
				observer.Next(3)
				observer.Error(assert.AnError)
				observer.Next(4)
			}()

			return nil
		}),
		BufferWithCount[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())

	time.Sleep(10 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: [1 2]
	// Error: assert.AnError general error for testing
}

// Commented because i get a weired conflict with other tests.
func ExampleBufferWithTime_ok() {
	observable := Pipe1(
		RangeWithInterval(1, 6, 20*time.Millisecond),
		BufferWithTime[int64](70*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[[]int64]())

	time.Sleep(200 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: [1 2 3]
	// Next: [4 5]
	// Completed
}

func ExampleBufferWithTime_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			go func() {
				observer.Next(1)
				time.Sleep(10 * time.Millisecond)
				observer.Next(2)
				time.Sleep(10 * time.Millisecond)
				observer.Next(3)

				time.Sleep(200 * time.Millisecond)
				// 1 empty buffer

				observer.Next(4)
				observer.Error(assert.AnError)
				observer.Next(5)
			}()

			return nil
		}),
		BufferWithTime[int](100*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())

	time.Sleep(300 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: [1 2 3]
	// Next: []
	// Error: assert.AnError general error for testing
}
