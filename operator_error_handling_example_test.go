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
	"errors"

	"github.com/stretchr/testify/assert"
)

func ExampleCatch() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)
			observer.Complete()

			return nil
		}),
		Catch(func(err error) Observable[int] {
			return Of(4, 5, 6)
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Next: 6
	// Completed
}

func ExampleOnErrorResumeNextWith() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)
			observer.Complete()

			return nil
		}),
		OnErrorResumeNextWith(Of(4, 5, 6)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Next: 6
	// Completed
}

func ExampleOnErrorReturn() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)
			observer.Complete()

			return nil
		}),
		OnErrorReturn(42),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 42
	// Completed
}

func ExampleRetryWithConfig() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)
			observer.Complete()

			return nil
		}),
		RetryWithConfig[int](RetryConfig{
			MaxRetries: 1,
		}),
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
	// Error: assert.AnError general error for testing
}

func ExampleThrowIfEmpty() {
	observable := Pipe1(
		Empty[int](),
		ThrowIfEmpty[int](func() error {
			return errors.New("empty")
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: empty
}

func ExampleDoWhile() {
	i := 0

	observable := Pipe1(
		Just(1, 2, 3),
		DoWhile[int](func() bool {
			i++
			return i < 2
		}),
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
	// Completed
}

func ExampleWhile() {
	i := 0

	observable := Pipe1(
		Just(1, 2, 3),
		While[int](func() bool {
			i++
			return i < 2
		}),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}
