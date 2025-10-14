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


package rostrconv

import (
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestAtoi(t *testing.T) {
	t.Run("successful conversions", func(t *testing.T) {
		input := []string{"123", "456", "789"}
		expected := []int{123, 456, 789}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Atoi[string](),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("invalid conversions", func(t *testing.T) {
		input := []string{"123", "abc", "456"}
		expected := []int{123} // Stream stops at first error

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Atoi[string](),
			),
		)

		assert.Equal(t, expected, values)
		assert.NotNil(t, err) // Should have an error for "abc"
	})
}

func TestParseInt(t *testing.T) {
	t.Run("decimal conversion", func(t *testing.T) {
		input := []string{"123", "456", "789"}
		expected := []int64{123, 456, 789}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseInt[string](10, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("hexadecimal conversion", func(t *testing.T) {
		input := []string{"FF", "1A", "2B"}
		expected := []int64{255, 26, 43}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseInt[string](16, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("invalid conversion", func(t *testing.T) {
		input := []string{"123", "invalid", "456"}
		expected := []int64{123} // Stream stops at first error

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseInt[string](10, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.NotNil(t, err) // Should have an error for "invalid"
	})
}

func TestParseFloat(t *testing.T) {
	t.Run("successful conversions", func(t *testing.T) {
		input := []string{"3.14", "2.718", "1.414"}
		expected := []float64{3.14, 2.718, 1.414}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseFloat[string](64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("scientific notation", func(t *testing.T) {
		input := []string{"1.23e+2", "4.56e-1", "7.89e0"}
		expected := []float64{123.0, 0.456, 7.89}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseFloat[string](64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("invalid conversion", func(t *testing.T) {
		input := []string{"3.14", "not_a_number", "2.718"}
		expected := []float64{3.14} // Stream stops at first error

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseFloat[string](64),
			),
		)

		assert.Equal(t, expected, values)
		assert.NotNil(t, err) // Should have an error for "not_a_number"
	})
}

func TestParseBool(t *testing.T) {
	t.Run("true values", func(t *testing.T) {
		input := []string{"true", "TRUE", "True", "1", "t", "T"}
		expected := []bool{true, true, true, true, true, true}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseBool[string](),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("false values", func(t *testing.T) {
		input := []string{"false", "FALSE", "False", "0", "f", "F"}
		expected := []bool{false, false, false, false, false, false}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseBool[string](),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		input := []string{"true", "maybe", "false"}
		expected := []bool{true} // Stream stops at first error

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseBool[string](),
			),
		)

		assert.Equal(t, expected, values)
		assert.NotNil(t, err) // Should have an error for "maybe"
	})
}

func TestParseUint(t *testing.T) {
	t.Run("decimal conversion", func(t *testing.T) {
		input := []string{"123", "456", "789"}
		expected := []uint64{123, 456, 789}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseUint[string](10, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("hexadecimal conversion", func(t *testing.T) {
		input := []string{"FF", "1A", "2B"}
		expected := []uint64{255, 26, 43}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseUint[string](16, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("invalid conversion", func(t *testing.T) {
		input := []string{"123", "invalid", "456"}
		expected := []uint64{123} // Stream stops at first error

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseUint[string](10, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.NotNil(t, err) // Should have an error for "invalid"
	})
}

func TestParseUint64(t *testing.T) {
	t.Run("successful conversion", func(t *testing.T) {
		input := []string{"123", "456", "789"}
		expected := []uint64{123, 456, 789}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				ParseUint64[string](10, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestFormatBool(t *testing.T) {
	t.Run("format boolean values", func(t *testing.T) {
		input := []bool{true, false, true}
		expected := []string{"true", "false", "true"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatBool(),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestFormatFloat(t *testing.T) {
	t.Run("fixed point format", func(t *testing.T) {
		input := []float64{3.14159, 2.71828}
		expected := []string{"3.14", "2.72"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatFloat('f', 2, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("scientific format", func(t *testing.T) {
		input := []float64{1234.5678, 0.001234}
		expected := []string{"1.23e+03", "1.23e-03"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatFloat('e', 2, 64),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestFormatComplex(t *testing.T) {
	t.Run("format complex numbers", func(t *testing.T) {
		input := []complex128{3 + 4i, 1 + 2i}
		expected := []string{"(3.00+4.00i)", "(1.00+2.00i)"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatComplex('f', 2, 128),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestFormatInt(t *testing.T) {
	t.Run("decimal format", func(t *testing.T) {
		input := []int64{123, 456, 789}
		expected := []string{"123", "456", "789"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatInt[string](10),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("hexadecimal format", func(t *testing.T) {
		input := []int64{255, 26, 43}
		expected := []string{"ff", "1a", "2b"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatInt[string](16),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestFormatUint(t *testing.T) {
	t.Run("decimal format", func(t *testing.T) {
		input := []uint64{123, 456, 789}
		expected := []string{"123", "456", "789"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatUint[string](10),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("hexadecimal format", func(t *testing.T) {
		input := []uint64{255, 26, 43}
		expected := []string{"ff", "1a", "2b"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				FormatUint[string](16),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestItoa(t *testing.T) {
	t.Run("convert integers to strings", func(t *testing.T) {
		input := []int{123, 456, 789}
		expected := []string{"123", "456", "789"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Itoa(),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestQuote(t *testing.T) {
	t.Run("quote strings", func(t *testing.T) {
		input := []string{"hello", "world\n", "test"}
		expected := []string{`"hello"`, `"world\n"`, `"test"`}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Quote(),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestQuoteRune(t *testing.T) {
	t.Run("quote runes", func(t *testing.T) {
		input := []rune{'a', 'b', 'c'}
		expected := []string{`'a'`, `'b'`, `'c'`}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				QuoteRune(),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})
}

func TestUnquote(t *testing.T) {
	t.Run("successful unquote", func(t *testing.T) {
		input := []string{`"hello"`, `"world\n"`, `"test"`}
		expected := []string{"hello", "world\n", "test"}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unquote(),
			),
		)

		assert.Equal(t, expected, values)
		assert.Nil(t, err)
	})

	t.Run("invalid quoted strings", func(t *testing.T) {
		input := []string{`"hello"`, `invalid`, `"test"`}
		expected := []string{"hello"} // Stream stops at first error

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Unquote(),
			),
		)

		assert.Equal(t, expected, values)
		assert.NotNil(t, err) // Should have an error for "invalid"
	})
}
