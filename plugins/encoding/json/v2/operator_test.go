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

package rojsonv2

import (
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

type nestedStruct struct {
	ID       int         `json:"id"`
	Data     testStruct  `json:"data"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("struct marshaling", func(t *testing.T) {
		input := []testStruct{
			{Name: "Alice", Age: 30, Email: "alice@example.com"},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		}
		expected := [][]byte{
			[]byte(`{"name":"Alice","age":30,"email":"alice@example.com"}`),
			[]byte(`{"name":"Bob","age":25}`),
			[]byte(`{"name":"Charlie","age":35,"email":"charlie@example.com"}`),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Marshal[testStruct](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("empty input", func(t *testing.T) {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Empty[testStruct](),
				Marshal[testStruct](),
			),
		)

		is.Equal([][]byte{}, values)
		is.Nil(err)
	})

	t.Run("primitive types", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		expected := [][]byte{
			[]byte(`"hello"`),
			[]byte(`"world"`),
			[]byte(`"test"`),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Marshal[string](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("numbers", func(t *testing.T) {
		input := []int{1, 2, 3, 42}
		expected := [][]byte{
			[]byte(`1`),
			[]byte(`2`),
			[]byte(`3`),
			[]byte(`42`),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Marshal[int](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("nested structs", func(t *testing.T) {
		input := []nestedStruct{
			{
				ID: 1,
				Data: testStruct{
					Name: "Alice",
					Age:  30,
				},
				Metadata: map[string]interface{}{
					"role": "admin",
				},
			},
			{
				ID: 2,
				Data: testStruct{
					Name:  "Bob",
					Age:   25,
					Email: "bob@example.com",
				},
			},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Marshal[nestedStruct](),
			),
		)

		is.Len(values, 2)
		is.Nil(err)
		// Verify the first result contains expected fields
		is.Contains(string(values[0]), `"id":1`)
		is.Contains(string(values[0]), `"name":"Alice"`)
		is.Contains(string(values[0]), `"role":"admin"`)
	})

	t.Run("maps", func(t *testing.T) {
		input := []map[string]interface{}{
			{"key1": "value1", "key2": 42},
			{"nested": map[string]interface{}{"inner": "value"}},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Marshal[map[string]interface{}](),
			),
		)

		is.Len(values, 2)
		is.Nil(err)
		is.Contains(string(values[0]), `"key1":"value1"`)
		is.Contains(string(values[0]), `"key2":42`)
		is.Contains(string(values[1]), `"nested"`)
	})

	t.Run("slices", func(t *testing.T) {
		input := [][]int{
			{1, 2, 3},
			{4, 5},
			{},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Marshal[[]int](),
			),
		)

		is.Len(values, 3)
		is.Nil(err)
		is.Equal([]byte(`[1,2,3]`), values[0])
		is.Equal([]byte(`[4,5]`), values[1])
		is.Equal([]byte(`[]`), values[2])
	})
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("struct unmarshaling", func(t *testing.T) {
		input := [][]byte{
			[]byte(`{"name":"Alice","age":30,"email":"alice@example.com"}`),
			[]byte(`{"name":"Bob","age":25}`),
			[]byte(`{"name":"Charlie","age":35,"email":"charlie@example.com"}`),
		}
		expected := []testStruct{
			{Name: "Alice", Age: 30, Email: "alice@example.com"},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[testStruct](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("empty input", func(t *testing.T) {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Empty[[]byte](),
				Unmarshal[testStruct](),
			),
		)

		is.Equal([]testStruct{}, values)
		is.Nil(err)
	})

	t.Run("primitive types", func(t *testing.T) {
		input := [][]byte{
			[]byte(`"hello"`),
			[]byte(`"world"`),
			[]byte(`"test"`),
		}
		expected := []string{"hello", "world", "test"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[string](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("numbers", func(t *testing.T) {
		input := [][]byte{
			[]byte(`1`),
			[]byte(`2`),
			[]byte(`3`),
			[]byte(`42`),
		}
		expected := []int{1, 2, 3, 42}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[int](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("nested structs", func(t *testing.T) {
		input := [][]byte{
			[]byte(`{"id":1,"data":{"name":"Alice","age":30},"metadata":{"role":"admin"}}`),
			[]byte(`{"id":2,"data":{"name":"Bob","age":25,"email":"bob@example.com"}}`),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[nestedStruct](),
			),
		)

		is.Len(values, 2)
		is.Nil(err)
		is.Equal(1, values[0].ID)
		is.Equal("Alice", values[0].Data.Name)
		is.Equal(30, values[0].Data.Age)
		is.Equal(2, values[1].ID)
		is.Equal("Bob", values[1].Data.Name)
		is.Equal("bob@example.com", values[1].Data.Email)
	})

	t.Run("maps", func(t *testing.T) {
		input := [][]byte{
			[]byte(`{"key1":"value1","key2":42}`),
			[]byte(`{"nested":{"inner":"value"}}`),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[map[string]interface{}](),
			),
		)

		is.Len(values, 2)
		is.Nil(err)
		is.Equal("value1", values[0]["key1"])
		is.Equal(float64(42), values[0]["key2"])
		is.NotNil(values[1]["nested"])
	})

	t.Run("slices", func(t *testing.T) {
		input := [][]byte{
			[]byte(`[1,2,3]`),
			[]byte(`[4,5]`),
			[]byte(`[]`),
		}
		expected := [][]int{
			{1, 2, 3},
			{4, 5},
			{},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[[]int](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		input := [][]byte{
			[]byte(`{"name":"Alice","age":30}`),
			[]byte(`invalid json`),
			[]byte(`{"name":"Bob","age":25}`),
		}
		expected := []testStruct{
			{Name: "Alice", Age: 30},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[testStruct](),
			),
		)

		is.Equal(expected, values)
		is.NotNil(err) // Should have an error for invalid JSON
	})

	t.Run("empty JSON object", func(t *testing.T) {
		input := [][]byte{
			[]byte(`{}`),
			[]byte(`{"name":"test"}`),
		}
		expected := []testStruct{
			{Name: "", Age: 0},
			{Name: "test", Age: 0},
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unmarshal[testStruct](),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("struct round trip", func(t *testing.T) {
		original := []testStruct{
			{Name: "Alice", Age: 30, Email: "alice@example.com"},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		}

		// Marshal
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Marshal[testStruct](),
			),
		)
		is.Nil(err)

		// Unmarshal
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Unmarshal[testStruct](),
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
					Name: "Alice",
					Age:  30,
				},
				Metadata: map[string]interface{}{
					"role": "admin",
				},
			},
			{
				ID: 2,
				Data: testStruct{
					Name:  "Bob",
					Age:   25,
					Email: "bob@example.com",
				},
			},
		}

		// Marshal
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Marshal[nestedStruct](),
			),
		)
		is.Nil(err)

		// Unmarshal
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Unmarshal[nestedStruct](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("primitive types round trip", func(t *testing.T) {
		original := []string{"hello", "world", "test"}

		// Marshal
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Marshal[string](),
			),
		)
		is.Nil(err)

		// Unmarshal
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Unmarshal[string](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("numbers round trip", func(t *testing.T) {
		original := []int{1, 2, 3, 42, -10, 0}

		// Marshal
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Marshal[int](),
			),
		)
		is.Nil(err)

		// Unmarshal
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Unmarshal[int](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("slices round trip", func(t *testing.T) {
		original := [][]int{
			{1, 2, 3},
			{4, 5},
			{},
			{10, 20, 30, 40},
		}

		// Marshal
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Marshal[[]int](),
			),
		)
		is.Nil(err)

		// Unmarshal
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Unmarshal[[]int](),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})
}
