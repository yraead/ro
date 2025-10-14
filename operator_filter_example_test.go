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

func ExampleFilter_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Filter(func(i int) bool {
			return i%2 == 0
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 4
	// Completed
}

func ExampleFilter_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Filter(func(i int) bool {
			return i%2 == 0
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Error: assert.AnError general error for testing
}

func ExampleDistinct_ok() {
	observable := Pipe1(
		Just(1, 1, 2, 2, 3, 3, 4, 4, 5, 5),
		Distinct[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleDistinct_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(1)
			observer.Next(2)
			observer.Next(2)
			observer.Next(3)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)
			observer.Next(4)

			return nil
		}),
		Distinct[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleIgnoreElements_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		IgnoreElements[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Completed
}

func ExampleIgnoreElements_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		IgnoreElements[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleSkip_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Skip[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleSkip_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Skip[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleSkipWhile_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		SkipWhile(func(v int) bool {
			return v > 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleSkipWhile_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		SkipWhile(func(v int) bool {
			return v > 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleSkipLast_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		SkipLast[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleSkipLast_empty() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		SkipLast[int](10),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Completed
}

func ExampleSkipLast_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		SkipLast[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Error: assert.AnError general error for testing
}

func ExampleSkipUntil_ok() {
	observable := Pipe1(
		RangeWithInterval(0, 5, 40*time.Millisecond),
		SkipUntil[int64](Interval(100*time.Millisecond)),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 2
	// Next: 3
	// Next: 4
	// Completed
}

func ExampleSkipUntil_empty() {
	observable := Pipe1(
		RangeWithInterval(0, 5, 10*time.Millisecond),
		SkipUntil[int64](Interval(100*time.Millisecond)),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Completed
}

func ExampleSkipUntil_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			go func() {
				time.Sleep(30 * time.Millisecond)
				observer.Next(1)
				time.Sleep(30 * time.Millisecond)
				observer.Next(2)
				time.Sleep(30 * time.Millisecond)
				observer.Next(3)
				observer.Error(assert.AnError)
				observer.Next(4)
			}()

			return nil
		}),
		SkipUntil[int](Interval(45*time.Millisecond)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleTake_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Take[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Completed
}

func ExampleTake_error1() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Take[int](5),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleTake_error2() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Take[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Completed
}

func ExampleTakeLast_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		TakeLast[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleTakeLast_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		TakeLast[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleTakeWhile_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		TakeWhile(func(n int) bool {
			return n < 3
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Completed
}

func ExampleTakeWhile_error1() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		TakeWhile(func(n int) bool {
			return n < 5
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleTakeWhile_error2() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		TakeWhile(func(n int) bool {
			return n < 3
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Completed
}

func ExampleTakeUntil_ok() {
	observable := Pipe1(
		RangeWithInterval(0, 5, 40*time.Millisecond),
		TakeUntil[int64](Interval(100*time.Millisecond)),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 0
	// Next: 1
	// Completed
}

func ExampleTakeUntil_empty() {
	observable := Pipe1(
		RangeWithInterval(0, 5, 50*time.Millisecond),
		TakeUntil[int64](Interval(10*time.Millisecond)),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Completed
}

func ExampleTakeUntil_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			go func() {
				time.Sleep(20 * time.Millisecond)
				observer.Next(1)
				time.Sleep(20 * time.Millisecond)
				observer.Next(2)
				time.Sleep(20 * time.Millisecond)
				observer.Next(3)
				observer.Error(assert.AnError)
				observer.Next(4)
			}()

			return nil
		}),
		TakeUntil[int](Interval(50*time.Millisecond)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 1
	// Next: 2
	// Completed
}

func ExampleHead_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Head[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Completed
}

func ExampleHead_error() {
	observable1 := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Head[int](),
	)

	subscription1 := observable1.Subscribe(PrintObserver[int]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Throw[int](assert.AnError), // no item transmitted
		Head[int](),
	)

	subscription2 := observable2.Subscribe(PrintObserver[int]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: 1
	// Completed
	// Error: assert.AnError general error for testing
}

func ExampleTail_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Tail[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 5
	// Completed
}

func ExampleTail_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Tail[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleFirst_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		First(func(n int) bool {
			return n > 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Completed
}

func ExampleFirst_error() {
	observable1 := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		First(func(n int) bool {
			return n > 2
		}),
	)

	subscription1 := observable1.Subscribe(PrintObserver[int]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Throw[int](assert.AnError), // no item transmitted
		First(func(n int) bool {
			return n > 2
		}),
	)

	subscription2 := observable2.Subscribe(PrintObserver[int]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: 3
	// Completed
	// Error: assert.AnError general error for testing
}

func ExampleLast_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Last(func(n int) bool {
			return n > 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 5
	// Completed
}

func ExampleLast_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Last(func(n int) bool {
			return n > 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleElementAt_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		ElementAt[int](2),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Completed
}

func ExampleElementAt_notFound() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		ElementAt[int](10),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: ro.ElementAt: nth element not found
}

func ExampleElementAt_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		ElementAt[int](10),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleElementAtOrDefault_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		ElementAtOrDefault(2, 100),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Completed
}

func ExampleElementAtOrDefault_notFound() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		ElementAtOrDefault(10, 100),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 100
	// Completed
}

func ExampleElementAtOrDefault_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		ElementAtOrDefault(10, 100),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}
