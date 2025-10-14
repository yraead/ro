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
	"io"
	"strings"
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

// failingReader is a reader that fails after reading a certain amount of data
type failingReader struct {
	io.Reader
}

func (fr *failingReader) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}

func TestNewCSVReader(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	reader := strings.NewReader("a,b,c\n1,2,3\n4,5,6\n")

	expected := [][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
		{"4", "5", "6"},
	}
	step := 0

	sub := NewCSVReader(csv.NewReader(reader)).Subscribe(
		ro.NewObserver(
			func(values []string) {
				is.Equal(expected[step], values)
				step++
			},
			func(err error) {
				is.Fail("never")
			},
			func() {},
		),
	)
	defer sub.Unsubscribe()
}

func TestNewCSVReader_Error(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	sub := NewCSVReader(csv.NewReader(&failingReader{})).Subscribe(
		ro.NewObserver(
			func(values []string) {
				is.Fail("should not be called")
			},
			func(err error) {
				is.Equal(assert.AnError, err)
			},
			func() {
				is.Fail("should not complete")
			},
		),
	)
	defer sub.Unsubscribe()
}

func TestNewCSVReader_EmptyInput(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	reader := strings.NewReader("")

	completed := false
	sub := NewCSVReader(csv.NewReader(reader)).Subscribe(
		ro.NewObserver(
			func(values []string) {
				is.Fail("should not emit any values")
			},
			func(err error) {
				is.Fail("should not error")
			},
			func() {
				completed = true
			},
		),
	)
	defer sub.Unsubscribe()

	is.True(completed)
}
