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

	"github.com/stretchr/testify/assert"
)

func ExampleNewObserver() {
	observer := NewObserver(
		func(value int) {
			fmt.Printf("Next: %d\n", value)
		},
		func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
		},
		func() {
			fmt.Printf("Completed\n")
		},
	)

	observer.Next(123)  // 123 logged
	observer.Next(456)  // 456 logged
	observer.Complete() // Completed logged

	observer.Next(789) // nothing logged

	// Output:
	// Next: 123
	// Next: 456
	// Completed
}

func ExampleNewObserver_error() {
	observer := NewObserver(
		func(value int) {
			fmt.Printf("Next: %d\n", value)
		},
		func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
		},
		func() {
			fmt.Printf("Completed\n")
		},
	)

	observer.Next(123)             // 123 logged
	observer.Next(456)             // 456 logged
	observer.Error(assert.AnError) // Completed logged

	observer.Next(789) // nothing logged

	// Output:
	// Next: 123
	// Next: 456
	// Error: assert.AnError general error for testing
}

func ExampleNewObserver_empty() {
	observer := NewObserver(
		func(value int) {
			fmt.Printf("Next: %d\n", value)
		},
		func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
		},
		func() {
			fmt.Printf("Completed\n")
		},
	)

	observer.Complete() // Completed logged

	observer.Next(123) // nothing logged

	// Output:
	// Completed
}
