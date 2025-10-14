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


package rogob

import (
	"testing"
	"time"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string
	Age   int
	Email string
}

type nestedStruct struct {
	ID       int
	Data     testStruct
	Metadata map[string]interface{}
	Created  time.Time
}

type complexStruct struct {
	ID       int
	Name     string
	Tags     []string
	Settings map[string]interface{}
	Data     []byte
	Active   bool
	Score    float64
}

func TestEncode(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("struct encoding", func(t *testing.T) {
		input := []testStruct{
			{Name: "Alice", Age: 30, Email: "alice@example.com"},
			{Name: "Bob", Age: 25, Email: "bob@example.com"},
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[testStruct](),
			),
		)

		is.Len(values, 3)
		is.Nil(err)
		// Verify that encoded values are not empty
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Empty[testStruct](),
				Encode[testStruct](),
			),
		)

		is.Equal([][]byte{}, values)
		is.Nil(err)
	})

	t.Run("primitive types", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[string](),
			),
		)

		is.Len(values, 3)
		is.Nil(err)
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("numbers", func(t *testing.T) {
		input := []int{1, 2, 3, 42, -10, 0}
		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[int](),
			),
		)

		is.Len(values, 6)
		is.Nil(err)
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("nested structs", func(t *testing.T) {
		input := []nestedStruct{
			{
				ID: 1,
				Data: testStruct{
					Name:  "Alice",
					Age:   30,
					Email: "alice@example.com",
				},
				Metadata: map[string]interface{}{
					"role":  "admin",
					"level": 5,
				},
				Created: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				ID: 2,
				Data: testStruct{
					Name:  "Bob",
					Age:   25,
					Email: "bob@example.com",
				},
				Metadata: map[string]interface{}{
					"role": "user",
				},
				Created: time.Date(2023, 2, 1, 12, 0, 0, 0, time.UTC),
			},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[nestedStruct](),
			),
		)

		is.Len(values, 2)
		is.Nil(err)
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("complex structs", func(t *testing.T) {
		input := []complexStruct{
			{
				ID:   1,
				Name: "Test User",
				Tags: []string{"admin", "developer"},
				Settings: map[string]interface{}{
					"theme":         "dark",
					"notifications": true,
				},
				Data:   []byte("test data"),
				Active: true,
				Score:  95.5,
			},
			{
				ID:   2,
				Name: "Another User",
				Tags: []string{"user"},
				Settings: map[string]interface{}{
					"theme": "light",
				},
				Data:   []byte("more data"),
				Active: false,
				Score:  87.2,
			},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[complexStruct](),
			),
		)

		is.Len(values, 2)
		is.Nil(err)
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("slices", func(t *testing.T) {
		input := [][]int{
			{1, 2, 3},
			{4, 5},
			{},
			{10, 20, 30, 40},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[[]int](),
			),
		)

		is.Len(values, 4)
		is.Nil(err)
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("maps", func(t *testing.T) {
		input := []map[string]interface{}{
			{"key1": "value1", "key2": 42, "key3": true},
			{"nested": map[string]interface{}{"inner": "value"}},
			{},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[map[string]interface{}](),
			),
		)

		is.Len(values, 1)
		is.ErrorContains(err, "gob: type not registered for interface: map[string]interface {}")
		for _, value := range values {
			is.NotEmpty(value)
		}
	})

	t.Run("pointers", func(t *testing.T) {
		ptr1 := &testStruct{Name: "Alice", Age: 30}
		ptr2 := &testStruct{Name: "Bob", Age: 25}
		input := []*testStruct{ptr1, ptr2, nil}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[*testStruct](),
			),
		)

		is.Len(values, 2)
		is.ErrorContains(err, "ro.Observer: unexpected error: gob: cannot encode nil pointer of type *rogob.testStruct")
		// First two should be non-empty, third (nil) should be empty
		is.NotEmpty(values[0])
		is.NotEmpty(values[1])
	})
}

func TestDecode(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("struct decoding", func(t *testing.T) {
		original := []testStruct{
			{Name: "Alice", Age: 30, Email: "alice@example.com"},
			{Name: "Bob", Age: 25, Email: "bob@example.com"},
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[testStruct](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[testStruct](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("empty input", func(t *testing.T) {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Empty[[]byte](),
				Decode[testStruct](),
			),
		)

		is.Equal([]testStruct{}, values)
		is.Nil(err)
	})

	t.Run("primitive types", func(t *testing.T) {
		original := []string{"hello", "world", "test"}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[string](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[string](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("numbers", func(t *testing.T) {
		original := []int{1, 2, 3, 42, -10, 0}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[int](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[int](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("nested structs", func(t *testing.T) {
		original := []nestedStruct{
			{
				ID: 1,
				Data: testStruct{
					Name:  "Alice",
					Age:   30,
					Email: "alice@example.com",
				},
				Metadata: map[string]interface{}{
					"role":  "admin",
					"level": 5,
				},
				Created: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				ID: 2,
				Data: testStruct{
					Name:  "Bob",
					Age:   25,
					Email: "bob@example.com",
				},
				Metadata: map[string]interface{}{
					"role": "user",
				},
				Created: time.Date(2023, 2, 1, 12, 0, 0, 0, time.UTC),
			},
		}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[nestedStruct](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[nestedStruct](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("complex structs", func(t *testing.T) {
		original := []complexStruct{
			{
				ID:   1,
				Name: "Test User",
				Tags: []string{"admin", "developer"},
				Settings: map[string]interface{}{
					"theme":         "dark",
					"notifications": true,
				},
				Data:   []byte("test data"),
				Active: true,
				Score:  95.5,
			},
			{
				ID:   2,
				Name: "Another User",
				Tags: []string{"user"},
				Settings: map[string]interface{}{
					"theme": "light",
				},
				Data:   []byte("more data"),
				Active: false,
				Score:  87.2,
			},
		}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[complexStruct](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[complexStruct](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("slices", func(t *testing.T) {
		original := [][]int{
			{1, 2, 3},
			{4, 5},
			nil,
			{10, 20, 30, 40},
		}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[[]int](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[[]int](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("maps", func(t *testing.T) {
		original := []map[string]interface{}{
			{"key1": "value1", "key2": 42, "key3": true},
			{},
		}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[map[string]interface{}](),
			),
		)
		is.Nil(err)

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[map[string]interface{}](),
			),
		)

		is.Equal(original, decoded)
		is.Nil(err)
	})

	t.Run("pointers", func(t *testing.T) {
		ptr1 := &testStruct{Name: "Alice", Age: 30}
		ptr2 := &testStruct{Name: "Bob", Age: 25}
		original := []*testStruct{ptr1, ptr2, nil}

		// First encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[*testStruct](),
			),
		)
		is.ErrorContains(err, "ro.Observer: unexpected error: gob: cannot encode nil pointer of type *rogob.testStruct")

		// Then decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[*testStruct](),
			),
		)

		is.Len(decoded, 2)
		is.Nil(err)
		// Check that the first two are not nil and have correct values
		is.NotNil(decoded[0])
		is.Equal("Alice", decoded[0].Name)
		is.Equal(30, decoded[0].Age)
		is.NotNil(decoded[1])
		is.Equal("Bob", decoded[1].Name)
		is.Equal(25, decoded[1].Age)
	})

	t.Run("invalid gob data", func(t *testing.T) {
		input := [][]byte{
			[]byte("invalid gob data"),
			[]byte("more invalid data"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[testStruct](),
			),
		)

		is.Empty(values)
		is.NotNil(err) // Should have an error for invalid gob data
	})
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("struct round trip", func(t *testing.T) {
		original := []testStruct{
			{Name: "Alice", Age: 30, Email: "alice@example.com"},
			{Name: "Bob", Age: 25, Email: "bob@example.com"},
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[testStruct](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[testStruct](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("nested struct round trip", func(t *testing.T) {
		original := []nestedStruct{
			{
				ID: 1,
				Data: testStruct{
					Name:  "Alice",
					Age:   30,
					Email: "alice@example.com",
				},
				Metadata: map[string]interface{}{
					"role":  "admin",
					"level": 5,
				},
				Created: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				ID: 2,
				Data: testStruct{
					Name:  "Bob",
					Age:   25,
					Email: "bob@example.com",
				},
				Metadata: map[string]interface{}{
					"role": "user",
				},
				Created: time.Date(2023, 2, 1, 12, 0, 0, 0, time.UTC),
			},
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[nestedStruct](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[nestedStruct](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("complex struct round trip", func(t *testing.T) {
		original := []complexStruct{
			{
				ID:   1,
				Name: "Test User",
				Tags: []string{"admin", "developer"},
				Settings: map[string]interface{}{
					"theme":         "dark",
					"notifications": true,
				},
				Data:   []byte("test data"),
				Active: true,
				Score:  95.5,
			},
			{
				ID:   2,
				Name: "Another User",
				Tags: []string{"user"},
				Settings: map[string]interface{}{
					"theme": "light",
				},
				Data:   []byte("more data"),
				Active: false,
				Score:  87.2,
			},
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[complexStruct](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[complexStruct](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("primitive types round trip", func(t *testing.T) {
		original := []string{"hello", "world", "test"}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[string](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[string](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("numbers round trip", func(t *testing.T) {
		original := []int{1, 2, 3, 42, -10, 0}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[int](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[int](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("slices round trip", func(t *testing.T) {
		original := [][]int{
			{1, 2, 3},
			{4, 5},
			{10, 20, 30, 40},
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[[]int](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[[]int](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("maps round trip", func(t *testing.T) {
		original := []map[string]interface{}{
			{"key1": "value1", "key2": 42, "key3": true},
			{},
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[map[string]interface{}](),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[map[string]interface{}](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})
}
