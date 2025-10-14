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

func ExampleMergeWith_ok() {
	// @TODO: implement
}

func ExampleMergeWith_error() {
	// @TODO: implement
}

func ExampleMergeWith1_ok() {
	observable := Pipe1(
		Delay[int](20*time.Millisecond)(Just(2)),
		MergeWith(Just(1)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(30 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Completed
}

func ExampleMergeWith1_error() {
	observable := Pipe1(
		Delay[int](20*time.Millisecond)(Throw[int](assert.AnError)),
		MergeWith(Just(1)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(30 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Error: assert.AnError general error for testing
}

func ExampleMergeWith2_ok() {
	// @TODO: implement
}

func ExampleMergeWith2_error() {
	// @TODO: implement
}

func ExampleMergeWith3_ok() {
	// @TODO: implement
}

func ExampleMergeWith3_error() {
	// @TODO: implement
}

func ExampleMergeWith4_ok() {
	// @TODO: implement
}

func ExampleMergeWith4_error() {
	// @TODO: implement
}

func ExampleMergeWith5_ok() {
	// @TODO: implement
}

func ExampleMergeWith5_error() {
	// @TODO: implement
}

func ExampleMergeAll_ok() {
	observable := Pipe1(
		Just(
			Just(1),
			Delay[int](20*time.Millisecond)(Just(2)),
			Delay[int](40*time.Millisecond)(Just(3)),
		),
		MergeAll[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(60 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleMergeAll_error() {
	observable := Pipe1(
		Just(
			Delay[int](10*time.Millisecond)(Just(1)),
			Delay[int](50*time.Millisecond)(Throw[int](assert.AnError)),
			Delay[int](100*time.Millisecond)(Just(3)),
		),
		MergeAll[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(300 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Error: assert.AnError general error for testing
}

func ExampleMergeMap_ok() {
	// @TODO: implement
}

func ExampleMergeMap_error() {
	// @TODO: implement
}

func ExampleCombineLatestWith_ok() {
	// @TODO: implement
}

func ExampleCombineLatestWith_error() {
	// @TODO: implement
}

func ExampleCombineLatestWith1_ok() {
	// @TODO: implement
}

func ExampleCombineLatestWith1_error() {
	// @TODO: implement
}

func ExampleCombineLatestWith2_ok() {
	// @TODO: implement
}

func ExampleCombineLatestWith2_error() {
	// @TODO: implement
}

func ExampleCombineLatestWith3_ok() {
	// @TODO: implement
}

func ExampleCombineLatestWith3_error() {
	// @TODO: implement
}

func ExampleCombineLatestWith4_ok() {
	// @TODO: implement
}

func ExampleCombineLatestWith4_error() {
	// @TODO: implement
}

func ExampleCombineLatestAll_ok() {
	// @TODO: implement
}

func ExampleCombineLatestAll_error() {
	// @TODO: implement
}

func ExampleCombineLatestAllAny_ok() {
	// @TODO: implement
}

func ExampleCombineLatestAllAny_error() {
	// @TODO: implement
}

func ExampleConcatWith_ok() {
	// @TODO: implement
}

func ExampleConcatWith_error() {
	// @TODO: implement
}

func ExampleConcatAll_ok() {
	observable := Pipe1(
		Just(
			Just(1, 2, 3),
			Just(4, 5, 6),
		),
		ConcatAll[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(30 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Next: 6
	// Completed
}

func ExampleConcatAll_error() {
	observable := Pipe1(
		Just(
			Just(1, 2, 3),
			Throw[int](assert.AnError),
			Just(4, 5, 6),
		),
		ConcatAll[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())

	time.Sleep(30 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleStartWith_ok() {
	observable := Pipe1(
		Just(4, 5, 6),
		StartWith(1, 2, 3),
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

func ExampleStartWith_error() {
	observable := Pipe1(
		Throw[int](assert.AnError),
		StartWith(1, 2, 3),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleEndWith_ok() {
	observable := Pipe1(
		Just(1, 2, 3),
		EndWith(4, 5, 6),
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

func ExampleEndWith_error() {
	observable := Pipe1(
		Throw[int](assert.AnError),
		EndWith(1, 2, 3),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExamplePairwise_ok() {
	obsercable := Pipe1(
		Just(1, 2, 3, 4, 5),
		Pairwise[int](),
	)

	subscription := obsercable.Subscribe(PrintObserver[[]int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [1 2]
	// Next: [2 3]
	// Next: [3 4]
	// Next: [4 5]
	// Completed
}

func ExamplePairwise_error() {
	observable := Pipe1(
		Throw[int](assert.AnError),
		Pairwise[int](),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleRaceWith_ok() {
	// @TODO: implement
}

func ExampleRaceWith_error() {
	// @TODO: implement
}

func ExampleZipWith_ok() {
	// @TODO: implement
}

func ExampleZipWith_error() {
	// @TODO: implement
}

func ExampleZipWith1_ok() {
	// @TODO: implement
}

func ExampleZipWith1_error() {
	// @TODO: implement
}

func ExampleZipWith2_ok() {
	// @TODO: implement
}

func ExampleZipWith2_error() {
	// @TODO: implement
}

func ExampleZipWith3_ok() {
	// @TODO: implement
}

func ExampleZipWith3_error() {
	// @TODO: implement
}

func ExampleZipWith4_ok() {
	// @TODO: implement
}

func ExampleZipWith4_error() {
	// @TODO: implement
}

func ExampleZipWith5_ok() {
	// @TODO: implement
}

func ExampleZipWith5_error() {
	// @TODO: implement
}

func ExampleZipAll_ok() {
	// @TODO: implement
}

func ExampleZipAll_error() {
	// @TODO: implement
}
