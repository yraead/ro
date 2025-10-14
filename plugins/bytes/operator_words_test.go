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

package robytes

import (
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestWords(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		args []byte
		want [][]byte
	}{
		{[]byte("PascalCase"), [][]byte{[]byte("Pascal"), []byte("Case")}},
		{[]byte("camelCase"), [][]byte{[]byte("camel"), []byte("Case")}},
		{[]byte("snake_case"), [][]byte{[]byte("snake"), []byte("case")}},
		{[]byte("kebab_case"), [][]byte{[]byte("kebab"), []byte("case")}},
		{[]byte("_test text_"), [][]byte{[]byte("test"), []byte("text")}},
		{[]byte("UPPERCASE"), [][]byte{[]byte("UPPERCASE")}},
		{[]byte("HTTPCode"), [][]byte{[]byte("HTTP"), []byte("Code")}},
		{[]byte("Int8Value"), [][]byte{[]byte("Int"), []byte("8"), []byte("Value")}},
	}

	for _, t := range tests {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Just([]byte(t.args)),
				Words[[]byte](),
			),
		)
		is.Equal([][][]byte{t.want}, values)
		is.Nil(err)

		values, err = ro.Collect(
			ro.Pipe1(
				ro.Empty[[]byte](),
				Words[[]byte](),
			),
		)
		is.Equal([][][]byte{}, values)
		is.Nil(err)

		values, err = ro.Collect(
			ro.Pipe1(
				ro.Throw[[]byte](assert.AnError),
				Words[[]byte](),
			),
		)
		is.Equal([][][]byte{}, values)
		is.EqualError(err, assert.AnError.Error())
	}
}
