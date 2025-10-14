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


package roiter

import (
	"strconv"
	"testing"

	"github.com/samber/lo"
	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestFromSeq(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	iterateItems := func(yield func(int) bool) {
		items := []int{1, 2, 3}
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}

	obs := FromSeq(iterateItems)
	values, err := ro.Collect(obs)

	is.Equal([]int{1, 2, 3}, values)
	is.Nil(err)
}

func TestFromSeq2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	iterateItems := func(yield func(int, string) bool) {
		items := []int{1, 2, 3}
		for _, v := range items {
			if !yield(v, strconv.Itoa(v)) {
				return
			}
		}
	}

	obs := FromSeq2(iterateItems)
	values, err := ro.Collect(obs)

	is.Equal([]lo.Tuple2[int, string]{lo.T2(1, "1"), lo.T2(2, "2"), lo.T2(3, "3")}, values)
	is.Nil(err)
}
