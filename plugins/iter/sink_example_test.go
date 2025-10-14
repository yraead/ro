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


package roiter

import (
	"fmt"

	"github.com/samber/ro"
)

func ExampleToSeq() {
	// Convert observable to sequence
	observable := ro.Just(1, 2, 3, 4, 5)
	seq := ToSeq[int](observable)

	// Iterate over the sequence
	for value := range seq {
		// Process each value
		_ = value
	}

	// Output: Will iterate over 1, 2, 3, 4, 5
}

func ExampleToSeq2() {
	// Convert observable to key-value sequence
	observable := ro.Just("Alice", "Bob", "Charlie")
	seq := ToSeq2[string](observable)

	// Iterate over the key-value sequence
	for index, value := range seq {
		// Process each key-value pair
		_ = index
		_ = value
	}

	// Output: Will iterate over (0, "Alice"), (1, "Bob"), (2, "Charlie")
}

func ExampleToSeq_withFiltering() {
	// Convert filtered observable to sequence
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		ro.Filter(func(n int) bool {
			return n%2 == 0 // Only even numbers
		}),
	)
	seq := ToSeq[int](observable)

	// Iterate over the filtered sequence
	for value := range seq {
		// Process each even value
		_ = value
	}

	// Output: Will iterate over 2, 4, 6, 8, 10
}

func ExampleToSeq_withTransformation() {
	// Convert transformed observable to sequence
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		ro.Map(func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		}),
	)
	seq := ToSeq[string](observable)

	// Iterate over the transformed sequence
	for value := range seq {
		// Process each transformed value
		_ = value
	}

	// Output: Will iterate over "Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"
}

func ExampleToSeq_withEarlyTermination() {
	// Convert observable to sequence with early termination
	observable := ro.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	seq := ToSeq[int](observable)

	// Iterate over the sequence with early termination
	for value := range seq {
		// Process each value
		if value > 5 {
			// Terminate early
			break
		}
	}

	// Output: Will iterate over 1, 2, 3, 4, 5 then terminate
}
