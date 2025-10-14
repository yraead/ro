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

func ExampleAll_ok() {
	observable1 := Pipe1(
		Just(1, 2, 3, 4, 5),
		All(func(i int) bool { return i > 0 }),
	)

	subscription1 := observable1.Subscribe(PrintObserver[bool]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Just(1, 2, 3, 4, 5),
		All(func(i int) bool { return i%2 == 0 }),
	)

	subscription2 := observable2.Subscribe(PrintObserver[bool]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: true
	// Completed
	// Next: false
	// Completed
}

func ExampleAll_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		All(func(i int) bool { return i > 0 }),
	)

	subscription := observable.Subscribe(PrintObserver[bool]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleContains_ok() {
	observable1 := Pipe1(
		Just(1, 2, 3, 4, 5),
		Contains(func(i int) bool { return i < 0 }),
	)

	subscription1 := observable1.Subscribe(PrintObserver[bool]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Just(1, 2, 3, 4, 5),
		Contains(func(i int) bool { return i%2 == 0 }),
	)

	subscription2 := observable2.Subscribe(PrintObserver[bool]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: false
	// Completed
	// Next: true
	// Completed
}

func ExampleContains_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Contains(func(i int) bool { return i == 4 }),
	)

	subscription := observable.Subscribe(PrintObserver[bool]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleFind_ok() {
	observable1 := Pipe1(
		Just(1, 2, 3, 4, 5),
		Find(func(i int) bool { return i < 0 }),
	)

	subscription1 := observable1.Subscribe(PrintObserver[int]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Just(1, 2, 3, 4, 5),
		Find(func(i int) bool { return i%2 == 0 }),
	)

	subscription2 := observable2.Subscribe(PrintObserver[int]())
	defer subscription2.Unsubscribe()

	// Output:
	// Completed
	// Next: 2
	// Completed
}

func ExampleFind_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		Find(func(i int) bool { return i == 4 }),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleIif_ok() {
	observable := Iif(
		func() bool {
			return true
		},
		Just(1, 2, 3),
		Just(4, 5, 6),
	)()

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleIif_error() {
	observable := Iif(
		func() bool {
			return false
		},
		Just(1, 2, 3),
		Throw[int](assert.AnError),
	)()

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleDefaultIfEmpty_ok() {
	observable1 := Pipe1(
		Just(1, 2, 3),
		DefaultIfEmpty(42),
	)

	subscription1 := observable1.Subscribe(PrintObserver[int]())
	defer subscription1.Unsubscribe()

	observable2 := Pipe1(
		Empty[int](),
		DefaultIfEmpty(42),
	)

	subscription2 := observable2.Subscribe(PrintObserver[int]())
	defer subscription2.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
	// Next: 42
	// Completed
}

func ExampleDefaultIfEmpty_error() {
	observable := Pipe1(
		Throw[int](assert.AnError),
		DefaultIfEmpty(42),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}
