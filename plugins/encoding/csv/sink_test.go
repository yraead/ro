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
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	io.Writer
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("mock error")
	// return m.Writer.Write(p)
}

func TestNewCSVWriter(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var writer strings.Builder

	sub := ro.Pipe1(
		ro.Just([]string{"a", "b", "c"}, []string{"1", "2", "3"}, []string{"4", "5", "6"}),
		NewCSVWriter(csv.NewWriter(&writer)),
	).Subscribe(ro.NewObserver(
		func(v int) {
			is.Equal(3, v)
			is.Equal("a,b,c\n1,2,3\n4,5,6\n", writer.String())
		},
		func(err error) {
			is.Fail("never")
		},
		func() {
		},
	))
	defer sub.Unsubscribe()
}

func TestNewCSVWriter_Error(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var buf bytes.Buffer
	mockWriter := &mockWriter{Writer: &buf}
	writer := csv.NewWriter(mockWriter)

	sub := ro.Pipe1(
		ro.NewObservable(func(destination ro.Observer[[]string]) ro.Teardown {
			destination.Next([]string{"a", "b", "c"})
			writer.Flush() // I need to flush manually here, because i want to test the error case
			destination.Next([]string{"1", "2", "3"})
			return nil
		}),
		NewCSVWriter(writer),
	).Subscribe(ro.NewObserver(
		func(v int) {
			is.Equal(1, v)
		},
		func(err error) {
			is.Equal("mock error", err.Error())
		},
		func() {
			is.Fail("should not complete")
		},
	))
	defer sub.Unsubscribe()
}
