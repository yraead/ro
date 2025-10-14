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

func ExampleNewBehaviorSubject() {
	subject := NewBehaviorSubject(42)

	subject.Subscribe(PrintObserver[int]()) // 42 logged by first subscriber

	subject.Next(123) // 123 logged by first subscriber

	subject.Subscribe(PrintObserver[int]()) // 123 logged by second subscriber

	subject.Next(456) // 123 logged by second subscriber

	subject.Complete() // 456 logged by both subscribers

	subject.Next(789)                       // nothing logged
	subject.Subscribe(PrintObserver[int]()) // nothing logged

	// Output:
	// Next: 42
	// Next: 123
	// Next: 123
	// Next: 456
	// Next: 456
	// Completed
	// Completed
	// Completed
}

func ExampleNewBehaviorSubject_error() {
	subject := NewBehaviorSubject(42)

	subject.Subscribe(PrintObserver[int]()) // 42 logged by first subscriber

	subject.Next(123) // 123 logged by first subscriber

	subject.Subscribe(PrintObserver[int]()) // 123 logged by second subscriber

	subject.Next(456) // nothing logged

	subject.Error(assert.AnError) // error logged by both subscribers

	subject.Subscribe(PrintObserver[int]()) // error logged by last subscriber

	subject.Next(789) // nothing logged

	// Output:
	// Next: 42
	// Next: 123
	// Next: 123
	// Next: 456
	// Next: 456
	// Error: assert.AnError general error for testing
	// Error: assert.AnError general error for testing
	// Error: assert.AnError general error for testing
}

func ExampleNewBehaviorSubject_empty() {
	subject := NewBehaviorSubject(42)

	subject.Complete() // nothing logged

	subject.Subscribe(PrintObserver[int]()) // nothing logged
	subject.Subscribe(PrintObserver[int]()) // nothing logged

	subject.Next(123) // nothing logged

	// Output:
	// Completed
	// Completed
}
