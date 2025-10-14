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


package rocsv

import (
	"bytes"
	"encoding/csv"

	"github.com/samber/ro"
)

func ExampleNewCSVWriter() {
	// Write CSV data to a buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	data := ro.Just(
		[]string{"name", "age", "city"},
		[]string{"Alice", "30", "New York"},
		[]string{"Bob", "25", "Los Angeles"},
		[]string{"Charlie", "35", "Chicago"},
	)

	observable := ro.Pipe1(
		data,
		NewCSVWriter(writer),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 4
	// Completed
}

func ExampleNewCSVWriter_withCustomDelimiter() {
	// Write CSV data with custom delimiter
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';' // Use semicolon as delimiter

	data := ro.Just(
		[]string{"name", "age", "city"},
		[]string{"Alice", "30", "New York"},
		[]string{"Bob", "25", "Los Angeles"},
		[]string{"Charlie", "35", "Chicago"},
	)

	observable := ro.Pipe1(
		data,
		NewCSVWriter(writer),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 4
	// Completed
}

func ExampleNewCSVWriter_withQuotedFields() {
	// Write CSV data with quoted fields
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	data := ro.Just(
		[]string{"name", "age", "city"},
		[]string{"Alice Smith", "30", "New York, NY"},
		[]string{"Bob Johnson", "25", "Los Angeles, CA"},
		[]string{"Charlie Brown", "35", "Chicago, IL"},
	)

	observable := ro.Pipe1(
		data,
		NewCSVWriter(writer),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 4
	// Completed
}

func ExampleNewCSVWriter_withError() {
	// Write CSV data with a writer that fails on write
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = '\n'

	data := ro.Just(
		[]string{"name", "age", "city"},
		[]string{"Alice", "30", "New York"},
		[]string{"Bob", "25", "Los Angeles"},
		[]string{"Charlie", "35", "Chicago"},
	)

	observable := ro.Pipe1(
		data,
		NewCSVWriter(writer),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 0
	// Error: csv: invalid field or comment delimiter
}
