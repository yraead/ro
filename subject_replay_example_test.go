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

import "github.com/stretchr/testify/assert"

func ExampleNewReplaySubject() {
	subject := NewReplaySubject[int](42)

	subject.Subscribe(PrintObserver[int]())

	subject.Next(123) // 123 logged by first subscriber

	subject.Subscribe(PrintObserver[int]()) // 123 logged by second subscriber

	subject.Next(456) // 456 logged by both subscriber

	subject.Complete()

	subject.Subscribe(PrintObserver[int]()) // 123 and 456 logged by third subscriber

	subject.Next(789) // nothing logged

	// Output:
	// Next: 123
	// Next: 123
	// Next: 456
	// Next: 456
	// Completed
	// Completed
	// Next: 123
	// Next: 456
	// Completed
}

func ExampleNewReplaySubject_error() {
	subject := NewReplaySubject[int](42)

	subject.Subscribe(PrintObserver[int]())

	subject.Next(123) // 123 logged by first subscriber

	subject.Subscribe(PrintObserver[int]()) // 123 logged by second subscriber

	subject.Next(456) // 456 logged by both subscriber

	subject.Error(assert.AnError) // error logged by both subscribers

	subject.Subscribe(PrintObserver[int]()) // error logged by last subscriber

	subject.Next(789)  // nothing logged
	subject.Complete() // nothing logged

	// Output:
	// Next: 123
	// Next: 123
	// Next: 456
	// Next: 456
	// Error: assert.AnError general error for testing
	// Error: assert.AnError general error for testing
	// Next: 123
	// Next: 456
	// Error: assert.AnError general error for testing
}

func ExampleNewReplaySubject_empty() {
	subject := NewReplaySubject[int](42)

	subject.Subscribe(PrintObserver[int]())

	subject.Complete() // nothing logged

	subject.Subscribe(PrintObserver[int]())

	subject.Next(123) // nothing logged

	// Output:
	// Completed
	// Completed
}

func ExampleNewReplaySubject_overflow() {
	subject := NewReplaySubject[int](2)

	subject.Next(123)  // nothing logged
	subject.Next(456)  // nothing logged
	subject.Next(789)  // nothing logged
	subject.Complete() // nothing logged

	subject.Subscribe(PrintObserver[int]()) // 456 and 789 logged

	// Output:
	// Next: 456
	// Next: 789
	// Completed
}
