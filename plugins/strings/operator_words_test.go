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

package rostrings

import (
	"testing"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestWords(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		args string
		want []string
	}{
		{"PascalCase", []string{"Pascal", "Case"}},
		{"camelCase", []string{"camel", "Case"}},
		{"snake_case", []string{"snake", "case"}},
		{"kebab_case", []string{"kebab", "case"}},
		{"_test text_", []string{"test", "text"}},
		{"UPPERCASE", []string{"UPPERCASE"}},
		{"HTTPCode", []string{"HTTP", "Code"}},
		{"Int8Value", []string{"Int", "8", "Value"}},
	}

	for _, t := range tests {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Just(t.args),
				Words[string](),
			),
		)
		is.Equal([][]string{t.want}, values)
		is.Nil(err)

		values, err = ro.Collect(
			ro.Pipe1(
				ro.Empty[string](),
				Words[string](),
			),
		)
		is.Equal([][]string{}, values)
		is.Nil(err)

		values, err = ro.Collect(
			ro.Pipe1(
				ro.Throw[string](assert.AnError),
				Words[string](),
			),
		)
		is.Equal([][]string{}, values)
		is.EqualError(err, assert.AnError.Error())
	}
}
