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


package robase64

import (
	"encoding/base64"

	"github.com/samber/ro"
)

func ExampleEncode() {
	// Encode byte slices to base64 strings
	observable := ro.Pipe1(
		ro.Just([]byte("hello"), []byte("world"), []byte("golang")),
		Encode[[]byte](base64.StdEncoding),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: aGVsbG8=
	// Next: d29ybGQ=
	// Next: Z29sYW5n
	// Completed
}

func ExampleDecode() {
	// Decode base64 strings to byte slices
	observable := ro.Pipe1(
		ro.Just("aGVsbG8=", "d29ybGQ=", "Z29sYW5n"),
		Decode[string](base64.StdEncoding),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [104 101 108 108 111]
	// Next: [119 111 114 108 100]
	// Next: [103 111 108 97 110 103]
	// Completed
}

func ExampleEncode_withURLEncoding() {
	// Encode using URL-safe base64 encoding
	observable := ro.Pipe1(
		ro.Just([]byte("hello world"), []byte("golang programming")),
		Encode[[]byte](base64.URLEncoding),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: aGVsbG8gd29ybGQ=
	// Next: Z29sYW5nIHByb2dyYW1taW5n
	// Completed
}

func ExampleDecode_withURLEncoding() {
	// Decode using URL-safe base64 encoding
	observable := ro.Pipe1(
		ro.Just("aGVsbG8gd29ybGQ=", "Z29sYW5nIHByb2dyYW1taW5n"),
		Decode[string](base64.URLEncoding),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [104 101 108 108 111 32 119 111 114 108 100]
	// Next: [103 111 108 97 110 103 32 112 114 111 103 114 97 109 109 105 110 103]
	// Completed
}

func ExampleDecode_withError() {
	// Decode with potential errors
	observable := ro.Pipe1(
		ro.Just("aGVsbG8=", "invalid-base64", "d29ybGQ="),
		Decode[string](base64.StdEncoding),
	)

	subscription := observable.Subscribe(ro.PrintObserver[[]byte]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [104 101 108 108 111]
	// Error: illegal base64 data at input byte 7
}
