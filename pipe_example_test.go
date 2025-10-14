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

func ExamplePipe() {
	observable := Pipe[int, int](
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
		Skip[int](2),
		Sum[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 24
	// Completed
}

func ExamplePipe1() {
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

func ExamplePipe2() {
	observable := Pipe2(
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
		Sum[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 30
	// Completed
}

func ExamplePipe3() {
	observable := Pipe3(
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
		Skip[int](2),
		Sum[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 24
	// Completed
}

func ExamplePipe4() {
	observable := Pipe4(
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
		Skip[int](2),
		Take[int](2),
		Sum[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 14
	// Completed
}

func ExamplePipe5() {
	observable := Pipe5(
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
		Skip[int](2),
		Take[int](2),
		Sum[int](),
		Max[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 14
	// Completed
}

func ExamplePipe6() {
	observable := Pipe6(
		Just(1, 2, 3, 4, 5),
		Map(func(x int) int {
			return x * 2
		}),
		Skip[int](2),
		Take[int](2),
		Sum[int](),
		Map(func(x int) int {
			return x / 2
		}),
		Max[int](),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 7
	// Completed
}

func ExamplePipeOp() {
	// @TODO: implement
}

func ExamplePipeOp4() {
	// @TODO: implement
}
