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
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("standard encoding", func(t *testing.T) {
		input := [][]byte{
			[]byte("hello"),
			[]byte("world"),
			[]byte("test"),
		}
		expected := []string{
			"aGVsbG8=",
			"d29ybGQ=",
			"dGVzdA==",
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[[]byte](base64.StdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("empty input", func(t *testing.T) {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Empty[[]byte](),
				Encode[[]byte](base64.StdEncoding),
			),
		)

		is.Equal([]string{}, values)
		is.Nil(err)
	})

	t.Run("empty bytes", func(t *testing.T) {
		input := [][]byte{
			[]byte{},
			[]byte("non-empty"),
		}
		expected := []string{
			"",
			"bm9uLWVtcHR5",
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[[]byte](base64.StdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("URL safe encoding", func(t *testing.T) {
		input := [][]byte{
			[]byte("hello world"),
			[]byte("test data"),
		}
		expected := []string{
			"aGVsbG8gd29ybGQ=",
			"dGVzdCBkYXRh",
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[[]byte](base64.URLEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("raw standard encoding", func(t *testing.T) {
		input := [][]byte{
			[]byte("hello"),
			[]byte("world"),
		}
		expected := []string{
			"aGVsbG8",
			"d29ybGQ",
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[[]byte](base64.RawStdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("raw URL encoding", func(t *testing.T) {
		input := [][]byte{
			[]byte("hello"),
			[]byte("world"),
		}
		expected := []string{
			"aGVsbG8",
			"d29ybGQ",
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Encode[[]byte](base64.RawURLEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})
}

func TestDecode(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("standard decoding", func(t *testing.T) {
		input := []string{
			"aGVsbG8=",
			"d29ybGQ=",
			"dGVzdA==",
		}
		expected := [][]byte{
			[]byte("hello"),
			[]byte("world"),
			[]byte("test"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.StdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("empty input", func(t *testing.T) {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Empty[string](),
				Decode[string](base64.StdEncoding),
			),
		)

		is.Equal([][]byte{}, values)
		is.Nil(err)
	})

	t.Run("empty string", func(t *testing.T) {
		input := []string{
			"",
			"bm9uLWVtcHR5",
		}
		expected := [][]byte{
			[]byte{},
			[]byte("non-empty"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.StdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("URL safe decoding", func(t *testing.T) {
		input := []string{
			"aGVsbG8gd29ybGQ=",
			"dGVzdCBkYXRh",
		}
		expected := [][]byte{
			[]byte("hello world"),
			[]byte("test data"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.URLEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("raw standard decoding", func(t *testing.T) {
		input := []string{
			"aGVsbG8",
			"d29ybGQ",
		}
		expected := [][]byte{
			[]byte("hello"),
			[]byte("world"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.RawStdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("raw URL decoding", func(t *testing.T) {
		input := []string{
			"aGVsbG8",
			"d29ybGQ",
		}
		expected := [][]byte{
			[]byte("hello"),
			[]byte("world"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.RawURLEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})

	t.Run("invalid base64 string", func(t *testing.T) {
		input := []string{
			"aGVsbG8=",
			"invalid-base64!",
			"d29ybGQ=",
		}
		expected := [][]byte{
			[]byte("hello"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.StdEncoding),
			),
		)

		is.Equal(expected, values)
		is.NotNil(err) // Should have an error for invalid base64
	})

	t.Run("incomplete padding", func(t *testing.T) {
		input := []string{
			"aGVsbG8",
			"d29ybGQ",
		}
		expected := [][]byte{
			[]byte("hello"),
			[]byte("world"),
		}

		values, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(input),
				Decode[string](base64.RawStdEncoding),
			),
		)

		is.Equal(expected, values)
		is.Nil(err)
	})
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("standard encoding round trip", func(t *testing.T) {
		original := [][]byte{
			[]byte("hello world"),
			[]byte("test data"),
			[]byte(""),
			[]byte("special chars: !@#$%^&*()"),
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[[]byte](base64.StdEncoding),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[string](base64.StdEncoding),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("URL safe encoding round trip", func(t *testing.T) {
		original := [][]byte{
			[]byte("hello world"),
			[]byte("test data"),
			[]byte(""),
			[]byte("special chars: !@#$%^&*()"),
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[[]byte](base64.URLEncoding),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[string](base64.URLEncoding),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})

	t.Run("raw encoding round trip", func(t *testing.T) {
		original := [][]byte{
			[]byte("hello world"),
			[]byte("test data"),
			[]byte(""),
			[]byte("special chars: !@#$%^&*()"),
		}

		// Encode
		encoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(original),
				Encode[[]byte](base64.RawStdEncoding),
			),
		)
		is.Nil(err)

		// Decode
		decoded, err := ro.Collect(
			ro.Pipe1(
				ro.FromSlice(encoded),
				Decode[string](base64.RawStdEncoding),
			),
		)
		is.Nil(err)

		is.Equal(original, decoded)
	})
}
