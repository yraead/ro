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
	"encoding/csv"
	"strings"

	"github.com/samber/ro"
)

func ExampleNewCSVReader() {
	// Read CSV data from a string
	csvData := `name,age,city
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago`

	reader := csv.NewReader(strings.NewReader(csvData))
	observable := NewCSVReader(reader)

	subscription := observable.Subscribe(ro.PrintObserver[[]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [name age city]
	// Next: [Alice 30 New York]
	// Next: [Bob 25 Los Angeles]
	// Next: [Charlie 35 Chicago]
	// Completed
}

func ExampleNewCSVReader_withCustomDelimiter() {
	// Read CSV data with custom delimiter
	csvData := `name;age;city
Alice;30;New York
Bob;25;Los Angeles
Charlie;35;Chicago`

	reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ';' // Use semicolon as delimiter
	observable := NewCSVReader(reader)

	subscription := observable.Subscribe(ro.PrintObserver[[]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [name age city]
	// Next: [Alice 30 New York]
	// Next: [Bob 25 Los Angeles]
	// Next: [Charlie 35 Chicago]
	// Completed
}

func ExampleNewCSVReader_withQuotedFields() {
	// Read CSV data with quoted fields
	csvData := `name,age,city
"Alice Smith",30,"New York, NY"
"Bob Johnson",25,"Los Angeles, CA"
"Charlie Brown",35,"Chicago, IL"`

	reader := csv.NewReader(strings.NewReader(csvData))
	observable := NewCSVReader(reader)

	subscription := observable.Subscribe(ro.PrintObserver[[]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [name age city]
	// Next: [Alice Smith 30 New York, NY]
	// Next: [Bob Johnson 25 Los Angeles, CA]
	// Next: [Charlie Brown 35 Chicago, IL]
	// Completed
}

func ExampleNewCSVReader_withError() {
	// Read CSV data with potential errors
	csvData := `name,age
Alice,30,New York
Bob,25,"Los Angeles, CA"
Charlie,35,Chicago`

	reader := csv.NewReader(strings.NewReader(csvData))
	observable := NewCSVReader(reader)

	subscription := observable.Subscribe(ro.PrintObserver[[]string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: [name age]
	// Error: record on line 2: wrong number of fields
}
