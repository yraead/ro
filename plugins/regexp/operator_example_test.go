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


package roregexp

import (
	"regexp"

	"github.com/samber/ro"
)

func ExampleFind() {
	// Find first match in byte slices
	pattern := regexp.MustCompile(`\d+`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("Hello 123 World"),
			[]byte("Test 456 Example"),
			[]byte("No numbers here"),
		),
		Find[[]byte](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [49 50 51]
	// Next: [52 53 54]
	// Next: []
	// Completed
}

func ExampleFindString() {
	// Find first match in strings
	pattern := regexp.MustCompile(`\d+`)
	observable := ro.Pipe2(
		ro.Just(
			"Hello 123 World",
			"Test 4567 Example",
			"No numbers here",
		),
		FindString[string](pattern),
		ro.Map(func(s string) int {
			return len(s)
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Next: 4
	// Next: 0
	// Completed
}

func ExampleFindSubmatch() {
	// Find first submatch in byte slices
	pattern := regexp.MustCompile(`(\d+)-(\w+)`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("123-abc"),
			[]byte("456-def"),
			[]byte("No match"),
		),
		FindSubmatch[[]byte](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[][]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [[49 50 51 45 97 98 99] [49 50 51] [97 98 99]]
	// Next: [[52 53 54 45 100 101 102] [52 53 54] [100 101 102]]
	// Next: []
	// Completed
}

func ExampleFindStringSubmatch() {
	// Find first submatch in strings
	pattern := regexp.MustCompile(`(\d+)-(\w+)`)
	observable := ro.Pipe1(
		ro.Just(
			"123-abc",
			"456-def",
			"No match",
		),
		FindStringSubmatch[string](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [123-abc 123 abc]
	// Next: [456-def 456 def]
	// Next: []
	// Completed
}

func ExampleFindAll() {
	// Find all matches in byte slices
	pattern := regexp.MustCompile(`\d+`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("Hello 123 World 456"),
			[]byte("Test 789 Example"),
			[]byte("No numbers here"),
		),
		FindAll[[]byte](pattern, -1),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[][]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [[49 50 51] [52 53 54]]
	// Next: [[55 56 57]]
	// Next: []
	// Completed
}

func ExampleFindAllString() {
	// Find all matches in strings
	pattern := regexp.MustCompile(`\d+`)
	observable := ro.Pipe1(
		ro.Just(
			"Hello 123 World 456",
			"Test 789 Example",
			"No numbers here",
		),
		FindAllString[string](pattern, -1),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [123 456]
	// Next: [789]
	// Next: []
	// Completed
}

func ExampleMatch() {
	// Check if byte slices match pattern
	pattern := regexp.MustCompile(`^\d+$`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("123"),
			[]byte("abc"),
			[]byte("456"),
		),
		Match[[]byte](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[bool]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: true
	// Next: false
	// Next: true
	// Completed
}

func ExampleMatchString() {
	// Check if strings match pattern
	pattern := regexp.MustCompile(`^\d+$`)
	observable := ro.Pipe1(
		ro.Just(
			"123",
			"abc",
			"456",
		),
		MatchString[string](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[bool]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: true
	// Next: false
	// Next: true
	// Completed
}

func ExampleReplaceAll() {
	// Replace matches in byte slices
	pattern := regexp.MustCompile(`\d+`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("Hello 123 World"),
			[]byte("Test 456 Example"),
		),
		ReplaceAll[[]byte](pattern, []byte("XXX")),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [72 101 108 108 111 32 88 88 88 32 87 111 114 108 100]
	// Next: [84 101 115 116 32 88 88 88 32 69 120 97 109 112 108 101]
	// Completed
}

func ExampleReplaceAllString() {
	// Replace matches in strings
	pattern := regexp.MustCompile(`\d+`)
	observable := ro.Pipe1(
		ro.Just(
			"Hello 123 World",
			"Test 456 Example",
		),
		ReplaceAllString[string](pattern, "XXX"),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Hello XXX World
	// Next: Test XXX Example
	// Completed
}

func ExampleFilterMatch() {
	// Filter byte slices that match pattern
	pattern := regexp.MustCompile(`^\d+$`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("123"),
			[]byte("abc"),
			[]byte("456"),
			[]byte("def"),
		),
		FilterMatch[[]byte](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [49 50 51]
	// Next: [52 53 54]
	// Completed
}

func ExampleFilterMatchString() {
	// Filter strings that match pattern
	pattern := regexp.MustCompile(`^\d+$`)
	observable := ro.Pipe1(
		ro.Just(
			"123",
			"abc",
			"456",
			"def",
		),
		FilterMatchString[string](pattern),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 123
	// Next: 456
	// Completed
}

func ExampleFindAllSubmatch() {
	// Find all submatches in byte slices
	pattern := regexp.MustCompile(`(\d+)-(\w+)`)
	observable := ro.Pipe1(
		ro.Just(
			[]byte("123-abc 456-def"),
			[]byte("789-ghi"),
			[]byte("No matches"),
		),
		FindAllSubmatch[[]byte](pattern, -1),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[][][]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [[[49 50 51 45 97 98 99] [49 50 51] [97 98 99]] [[52 53 54 45 100 101 102] [52 53 54] [100 101 102]]]
	// Next: [[[55 56 57 45 103 104 105] [55 56 57] [103 104 105]]]
	// Next: []
	// Completed
}

func ExampleFindAllStringSubmatch() {
	// Find all submatches in strings
	pattern := regexp.MustCompile(`(\d+)-(\w+)`)
	observable := ro.Pipe1(
		ro.Just(
			"123-abc 456-def",
			"789-ghi",
			"No matches",
		),
		FindAllStringSubmatch[string](pattern, -1),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[][]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [[123-abc 123 abc] [456-def 456 def]]
	// Next: [[789-ghi 789 ghi]]
	// Next: []
	// Completed
}
