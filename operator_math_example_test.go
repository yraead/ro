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
	"github.com/stretchr/testify/assert"
)

func ExampleAverage_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Average[int](),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Completed
}

func ExampleAverage_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Average[int](),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleCount_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Count[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 5
	// Completed
}

func ExampleCount_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Count[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleSum_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Sum[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 15
	// Completed
}

func ExampleSum_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Sum[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleRound_ok() {
	observable := Pipe1(
		Just[float64](1, 2, 3, 4, 5),
		Round(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleRound_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[float64]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Round(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleMin_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Min[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Completed
}

func ExampleMin_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Min[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleMax_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Max[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 5
	// Completed
}

func ExampleMax_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Max[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleClamp_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Clamp[int](2, 4),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 4
	// Completed
}

func ExampleClamp_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Clamp(2, 4),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleAbs_ok() {
	observable := Pipe1(
		Just[float64](-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5),
		Abs(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 5
	// Next: 4
	// Next: 3
	// Next: 2
	// Next: 1
	// Next: 0
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleAbs_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[float64]) Teardown {
			observer.Next(-3)
			observer.Next(-2)
			observer.Next(-1)
			observer.Next(0)
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Abs(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Next: 2
	// Next: 1
	// Next: 0
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleFloor_ok() {
	observable := Pipe1(
		Just(1.1, 2.4, 3.5, 4.9, 5.0),
		Floor(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleFloor_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[float64]) Teardown {
			observer.Next(1.1)
			observer.Next(2.5)
			observer.Next(3.9)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Floor(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleCeil_ok() {
	observable := Pipe1(
		Just(1.1, 2.4, 3.5, 4.9, 5.0),
		Ceil(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Next: 5
	// Completed
}

func ExampleCeil_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[float64]) Teardown {
			observer.Next(1.1)
			observer.Next(2.5)
			observer.Next(3.9)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Ceil(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next: 3
	// Next: 4
	// Error: assert.AnError general error for testing
}

func ExampleTrunc_ok() {
	observable := Pipe1(
		Just(1.1, 2.4, 3.5, 4.9, 5.0),
		Trunc(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Completed
}

func ExampleTrunc_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[float64]) Teardown {
			observer.Next(1.1)
			observer.Next(2.5)
			observer.Next(3.9)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Trunc(),
	)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleReduce_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Reduce(func(agg, current int) int {
			return agg + current
		}, 42),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 57
	// Completed
}

func ExampleReduce_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Reduce(func(agg, current int) int {
			return agg + current
		}, 42),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}
