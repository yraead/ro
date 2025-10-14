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


package rotestify

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func ExampleTestify() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that emits values
	observable := ro.Just(1, 2, 3, 4, 5)

	// Test the observable behavior
	Testify[int](is).
		Source(observable).
		ExpectNext(1).
		ExpectNext(2).
		ExpectNext(3).
		ExpectNext(4).
		ExpectNext(5).
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if observable emits exactly 1, 2, 3, 4, 5 and completes")

	// Output: Test passes if observable emits exactly 1, 2, 3, 4, 5 and completes
}

func ExampleTestify_error() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that emits an error
	observable := ro.Pipe1(
		ro.Just(1, 2, 3),
		ro.MapErr(func(n int) (int, error) {
			if n == 2 {
				return n, errors.New("error on 2")
			}
			return n, nil
		}),
	)

	// Test the observable behavior with error handling
	Testify[int](is).
		Source(observable).
		ExpectNext(1).
		ExpectError(errors.New("error on 2")).
		Verify()

	fmt.Println("Test passes if observable emits 1, then error \"error on 2\"")

	// Output: Test passes if observable emits 1, then error "error on 2"
}

func ExampleTestify_sequence() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that emits a sequence
	observable := ro.Just("a", "b", "c")

	// Test the observable behavior using sequence expectation
	Testify[string](is).
		Source(observable).
		ExpectNextSeq("a", "b", "c").
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if observable emits exactly \"a\", \"b\", \"c\" in sequence and completes")

	// Output: Test passes if observable emits exactly "a", "b", "c" in sequence and completes
}

func ExampleTestify_customMessages() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that emits a single value
	observable := ro.Just(42)

	// Test the observable behavior with custom error messages
	Testify[int](is).
		Source(observable).
		ExpectNext(42, "expected the answer to life").
		ExpectComplete("should complete after one value").
		Verify()

	fmt.Println("Test passes if observable emits 42 and completes, with custom messages on failure")

	// Output: Test passes if observable emits 42 and completes, with custom messages on failure
}

func ExampleTestify_filtered() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that filters even numbers
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5, 6),
		ro.Filter(func(n int) bool {
			return n%2 == 0 // Keep only even numbers
		}),
	)

	// Test the filtered observable behavior
	Testify[int](is).
		Source(observable).
		ExpectNextSeq(2, 4, 6).
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if observable emits only even numbers: 2, 4, 6")

	// Output: Test passes if observable emits only even numbers: 2, 4, 6
}

func ExampleTestify_mapped() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that transforms strings to uppercase
	observable := ro.Pipe1(
		ro.Just("hello", "world"),
		ro.Map(func(s string) string {
			return strings.ToUpper(s)
		}),
	)

	// Test the mapped observable behavior
	Testify[string](is).
		Source(observable).
		ExpectNextSeq("HELLO", "WORLD").
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if observable transforms strings to uppercase: \"HELLO\", \"WORLD\"")

	// Output: Test passes if observable transforms strings to uppercase: "HELLO", "WORLD"
}

func ExampleTestify_errorScenarios() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Test immediate error scenario
	immediateErrorObservable := ro.Throw[int](errors.New("test error"))

	Testify[int](is).
		Source(immediateErrorObservable).
		ExpectError(errors.New("test error")).
		Verify()

	// Test error after values scenario
	errorAfterValuesObservable := ro.Pipe1(
		ro.Just(1, 2, 3),
		ro.MapErr(func(n int) (int, error) {
			if n == 3 {
				return n, errors.New("error on 3")
			}
			return n, nil
		}),
	)

	Testify[int](is).
		Source(errorAfterValuesObservable).
		ExpectNext(1).
		ExpectNext(2).
		ExpectError(errors.New("error on 3")).
		Verify()

	fmt.Println("First test passes if observable immediately emits error \"test error\"")
	fmt.Println("Second test passes if observable emits 1, 2, then error \"error on 3\"")

	// Output:
	// First test passes if observable immediately emits error "test error"
	// Second test passes if observable emits 1, 2, then error "error on 3"
}

func ExampleTestify_context() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)
	ctx := context.Background()

	// Create an observable that emits values
	observable := ro.Just(1, 2, 3)

	// Test the observable behavior with context
	Testify[int](is).
		Source(observable).
		ExpectNextSeq(1, 2, 3).
		ExpectComplete().
		VerifyWithContext(ctx)

	fmt.Println("Test passes if observable emits 1, 2, 3 and completes with context support")

	// Output: Test passes if observable emits 1, 2, 3 and completes with context support
}

func ExampleTestify_complexPipeline() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create a complex pipeline: filter -> map -> take
	observable := ro.Pipe3(
		ro.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		ro.Filter(func(n int) bool {
			return n%2 == 0 // Keep only even numbers
		}),
		ro.Map(func(n int) int {
			return n * 2 // Double the values
		}),
		ro.Take[int](3), // Take only first 3 values
	)

	// Test the complex pipeline behavior
	Testify[int](is).
		Source(observable).
		ExpectNextSeq(4, 8, 12). // 2*2, 4*2, 6*2
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if pipeline filters even numbers, doubles them, and takes first 3: 4, 8, 12")

	// Output: Test passes if pipeline filters even numbers, doubles them, and takes first 3: 4, 8, 12
}

func ExampleTestify_empty() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an empty observable
	observable := ro.Empty[int]()

	// Test that empty observable completes without emitting values
	Testify[int](is).
		Source(observable).
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if observable completes without emitting any values")

	// Output: Test passes if observable completes without emitting any values
}

func ExampleTestify_singleValue() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that emits a single value
	observable := ro.Just("single value")

	// Test single value observable
	Testify[string](is).
		Source(observable).
		ExpectNext("single value").
		ExpectComplete().
		Verify()

	fmt.Println("Test passes if observable emits exactly \"single value\" and completes")

	// Output: Test passes if observable emits exactly "single value" and completes
}

func ExampleTestify_errorOnly() {
	// Create a test instance
	t := &testing.T{}
	is := assert.New(t)

	// Create an observable that only emits an error
	observable := ro.Throw[string](errors.New("something went wrong"))

	// Test error-only observable
	Testify[string](is).
		Source(observable).
		ExpectError(errors.New("something went wrong")).
		Verify()

	fmt.Println("Test passes if observable emits exactly the expected error")

	// Output: Test passes if observable emits exactly the expected error
}
