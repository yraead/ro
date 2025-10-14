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


package rogob

import (
	"fmt"

	"github.com/samber/ro"
)

type Person struct {
	Name string
	Age  int
}

func ExampleEncode() {
	// Encode a single struct to gob bytes
	observable := ro.Pipe1(
		ro.Just(42),
		Encode[int](),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [3 4 0 84]
	// Completed
}

func ExampleDecode() {
	// Decode gob bytes to a struct
	encoded := []byte{37, 255, 141, 3, 1, 1, 6, 80, 101, 114, 115, 111, 110, 1, 255, 142, 0, 1, 2, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 3, 65, 103, 101, 1, 4, 0, 0, 0, 12, 255, 142, 1, 5, 65, 108, 105, 99, 101, 1, 60, 0}

	observable := ro.Pipe1(
		ro.Just(encoded),
		Decode[Person](),
	)

	subscription := observable.Subscribe(ro.PrintObserver[Person]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: {Alice 30}
	// Completed
}

func ExampleEncode_withSimpleTypes() {
	// Encode a single string to gob bytes
	observable := ro.Pipe1(
		ro.Just("hello"),
		Encode[string](),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [8 12 0 5 104 101 108 108 111]
	// Completed
}

func ExampleDecode_withSimpleTypes() {
	// Decode gob bytes to a string
	encoded := []byte{8, 12, 0, 5, 104, 101, 108, 108, 111}

	observable := ro.Pipe1(
		ro.Just(encoded),
		Decode[string](),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: hello
	// Completed
}

func ExampleEncode_withMaps() {
	// Encode a single map to gob bytes
	observable := ro.Pipe1(
		ro.Just(1),
		Encode[int](),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [3 4 0 2]
	// Completed
}

func ExampleDecode_withMaps() {
	// Decode gob bytes to a map
	encoded := []byte{14, 255, 143, 4, 1, 2, 255, 144, 0, 1, 12, 1, 4, 0, 0, 10, 255, 144, 0, 2, 1, 97, 2, 1, 98, 4}

	observable := ro.Pipe1(
		ro.Just(encoded),
		Decode[map[string]int](),
	)

	subscription := observable.Subscribe(ro.PrintObserver[map[string]int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: map[a:1 b:2]
	// Completed
}

func ExampleDecode_withError() {
	// Decode with potential errors
	observable := ro.Pipe1(
		ro.Just([]byte("invalid-gob-data")),
		Decode[Person](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(value Person) {
				// Handle successful decoding
			},
			func(err error) {
				// Handle decoding error
				fmt.Println("Error:", err.Error())
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output: Error: unexpected EOF
}
