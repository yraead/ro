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
	"strconv"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func ExampleToSlice_ok() {
	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		ToSlice[int](),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [1 2 3 4 5]
	// Completed
}

func ExampleToSlice_error() {
	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		ToSlice[int](),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleToMap_ok() {
	mapper := func(v int) (string, string) {
		return strconv.FormatInt(int64(v), 10), strconv.FormatInt(int64(v), 10)
	}

	observable := Pipe1(
		Just(1, 2, 3, 4, 5),
		ToMap(mapper),
	)

	subscription := observable.Subscribe(PrintObserver[map[string]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: map[1:1 2:2 3:3 4:4 5:5]
	// Completed
}

func ExampleToMap_error() {
	mapper := func(v int) (string, string) {
		return strconv.FormatInt(int64(v), 10), strconv.FormatInt(int64(v), 10)
	}

	observable := Pipe1(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		ToMap(mapper),
	)

	subscription := observable.Subscribe(PrintObserver[map[string]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleToChannel_ok() {
	observable := Pipe3(
		Just(1, 2, 3, 4, 5),
		ToChannel[int](42),
		Map(lo.ChannelToSlice[Notification[int]]),
		Flatten[Notification[int]](),
	)

	subscription := observable.Subscribe(PrintObserver[Notification[int]]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: Next(1)
	// Next: Next(2)
	// Next: Next(3)
	// Next: Next(4)
	// Next: Next(5)
	// Next: Complete()
	// Completed
}

func ExampleToChannel_error() {
	observable := Pipe3(
		NewObservable(func(observer Observer[int]) Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Next(3)
			observer.Error(assert.AnError)
			observer.Next(4)

			return nil
		}),
		ToChannel[int](42),
		Map(lo.ChannelToSlice[Notification[int]]),
		Flatten[Notification[int]](),
	)

	subscription := observable.Subscribe(PrintObserver[Notification[int]]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: Next(1)
	// Next: Next(2)
	// Next: Next(3)
	// Next: Error(assert.AnError general error for testing)
	// Completed
}
