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

func ExampleNewAsyncSubject() {
	subject := NewAsyncSubject[int]()

	subject.Subscribe(PrintObserver[int]())

	subject.Next(123) // nothing logged

	sub := Pipe1(
		subject.AsObservable(),
		Delay[int](10*time.Millisecond),
	).Subscribe(PrintObserver[int]())
	defer sub.Unsubscribe()

	subject.Next(456) // nothing logged

	subject.Complete() // 456 logged by both subscribers

	time.Sleep(30 * time.Millisecond)

	subject.Next(789)                       // nothing logged
	subject.Subscribe(PrintObserver[int]()) // 456 logged by both subscribers

	// Output:
	// Next: 456
	// Completed
	// Next: 456
	// Completed
	// Next: 456
	// Completed
}

func ExampleNewAsyncSubject_error() {
	subject := NewAsyncSubject[int]()

	subject.Subscribe(PrintObserver[int]())

	subject.Next(123) // nothing logged

	subject.Subscribe(PrintObserver[int]())

	subject.Next(456) // nothing logged

	subject.Error(assert.AnError) // error logged by both subscribers

	subject.Subscribe(PrintObserver[int]()) // error logged by last subscriber

	subject.Next(789)  // nothing logged
	subject.Complete() // nothing logged

	// Output:
	// Error: assert.AnError general error for testing
	// Error: assert.AnError general error for testing
	// Error: assert.AnError general error for testing
}

func ExampleNewAsyncSubject_empty() {
	subject := NewAsyncSubject[int]()

	subject.Subscribe(PrintObserver[int]())

	subject.Complete() // nothing logged

	subject.Subscribe(PrintObserver[int]())

	subject.Next(123) // nothing logged

	// Output:
	// Completed
	// Completed
}
