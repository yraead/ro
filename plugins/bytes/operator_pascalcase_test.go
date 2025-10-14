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

func TestPascalCase(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	for _, t := range allCaseTests {
		values, err := ro.Collect(
			ro.Pipe1(
				ro.Just([]byte(t.input)),
				PascalCase[[]byte](),
			),
		)
		is.Equal([]byte(t.output.PascalCase), values[0])
		is.Nil(err)

		values, err = ro.Collect(
			ro.Pipe1(
				ro.Empty[[]byte](),
				PascalCase[[]byte](),
			),
		)
		is.Empty(values)
		is.Nil(err)

		values, err = ro.Collect(
			ro.Pipe1(
				ro.Throw[[]byte](assert.AnError),
				PascalCase[[]byte](),
			),
		)
		is.Empty(values)
		is.EqualError(err, assert.AnError.Error())
	}
}
