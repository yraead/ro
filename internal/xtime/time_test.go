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

package xtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNowNanoMonotonic(t *testing.T) { //nolint:paralleltest
	// t.Parallel()
	is := assert.New(t)

	got1 := NowNanoMonotonic()
	time.Sleep(100 * time.Millisecond)
	got2 := NowNanoMonotonic()
	is.InDelta(100*time.Millisecond, got2-got1, float64(10*time.Millisecond))

	got3 := []int64{}
	for i := 0; i < 1000; i++ {
		got3 = append(got3, NowNanoMonotonic())

		time.Sleep(10 * time.Nanosecond)
	}

	is.IsIncreasing(got3)
}
