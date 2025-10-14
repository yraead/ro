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
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleTap_ok() {
	observable := Pipe1(
		Range(1, 4),
		Tap(
			func(value int64) {
				fmt.Printf("Next: %v\n", value)
			},
			func(err error) {
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)

	subscription := observable.Subscribe(NoopObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleTap_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Tap(
			func(value int) {
				fmt.Printf("Next: %v\n", value)
			},
			func(err error) {
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)

	subscription := observable.Subscribe(NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleTapOnNext_ok() {
	observable := Pipe1(
		Range(1, 4),
		TapOnNext(func(v int64) { fmt.Println("Next:", v) }),
	)

	subscription := observable.Subscribe(NoopObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
}

func ExampleTapOnNext_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int64]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		TapOnNext(func(v int64) { fmt.Println("Next:", v) }),
	)

	subscription := observable.Subscribe(NoopObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
}

func ExampleTapOnError_ok() {
	observable := Pipe1(
		Range(1, 4),
		TapOnError[int64](func(err error) { fmt.Printf("Error: %s\n", err.Error()) }),
	)

	subscription := observable.Subscribe(NoopObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
}

func ExampleTapOnError_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		TapOnError[int](func(err error) { fmt.Printf("Error: %s\n", err.Error()) }),
	)

	subscription := observable.Subscribe(NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleTapOnComplete_ok() {
	observable := Pipe1(
		Range(1, 4),
		TapOnComplete[int64](func() { fmt.Printf("Completed") }),
	)

	subscription := observable.Subscribe(NoopObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Completed
}

func ExampleTapOnComplete_error() {
	observable := Pipe2(
		Throw[int](assert.AnError),
		Delay[int](10*time.Millisecond),
		TapOnComplete[int](func() { fmt.Printf("Completed") }),
	)

	subscription := observable.Subscribe(NoopObserver[int]())
	subscription.Wait()

	// Output:
}

func ExampleTimeInterval() {
	observable := Pipe1(
		RangeWithInterval(0, 3, 10*time.Millisecond),
		TimeInterval[int64](),
	)

	subscription := observable.Subscribe(NoopObserver[IntervalValue[int64]]())
	defer subscription.Unsubscribe()
}

func ExampleTimestamp() {
	observable := Pipe1(
		RangeWithInterval(0, 3, 10*time.Millisecond),
		Timestamp[int64](),
	)

	subscription := observable.Subscribe(NoopObserver[TimestampValue[int64]]())
	defer subscription.Unsubscribe()
}

func ExampleDelay_ok() {
	observable := Pipe1(
		Just(1, 2, 3),
		Delay[int](10*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleDelay_cancel() {
	observable := Pipe1(
		Of(1),
		Delay[int](100*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(50 * time.Millisecond)
	subscription.Unsubscribe() // canceled before first message

	// Output:
}

func ExampleDelay_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Delay[int](10*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleRepeatWith_ok() {
	observable := Pipe1(
		Just(1, 2, 3),
		RepeatWith[int](3),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleRepeatWith_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		RepeatWith[int](3),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleTimeout_ok() {
	observable := Pipe1(
		Range(1, 4),
		Timeout[int64](20*time.Millisecond),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait()
	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleTimeout_error() {
	subscription := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			go func() {
				observer.Next(1)
				time.Sleep(100 * time.Millisecond)
				observer.Next(2)
				time.Sleep(100 * time.Millisecond)
				observer.Next(3)
				observer.Error(assert.AnError)
				observer.Next(4)
			}()
			return nil
		}),
		Timeout[int](50*time.Millisecond),
	).Subscribe(PrintObserver[int]())

	subscription.Wait()

	// Output:
	// Next: 1
	// Error: ro.Timeout: timeout after 50ms
}

func ExampleMaterialize_ok() {
	observable := Pipe1(
		Just(1, 2, 3),
		Materialize[int](),
	)

	subscription := observable.Subscribe(PrintObserver[Notification[int]]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Next(1)
	// Next: Next(2)
	// Next: Next(3)
	// Next: Complete()
	// Completed
}

func ExampleMaterialize_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Materialize[int](),
	)

	subscription := observable.Subscribe(PrintObserver[Notification[int]]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Next(1)
	// Next: Next(2)
	// Next: Next(3)
	// Next: Error(assert.AnError general error for testing)
	// Completed
}

func ExampleDematerialize_ok() {
	observable := Pipe1(
		Just(
			Notification[int]{Kind: KindNext, Value: 1, Err: nil},
			Notification[int]{Kind: KindNext, Value: 2, Err: nil},
			Notification[int]{Kind: KindNext, Value: 3, Err: nil},
			Notification[int]{Kind: KindComplete, Value: 0, Err: nil},
		),
		Dematerialize[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleDematerialize_error() {
	observable := Pipe1(
		Just(
			Notification[int]{Kind: KindNext, Value: 1, Err: nil},
			Notification[int]{Kind: KindNext, Value: 2, Err: nil},
			Notification[int]{Kind: KindNext, Value: 3, Err: nil},
			Notification[int]{Kind: KindError, Value: 0, Err: assert.AnError},
		),
		Dematerialize[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}
