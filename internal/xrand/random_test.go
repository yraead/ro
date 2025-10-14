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

package xrand_test

import (
	"testing"

	"github.com/samber/ro/internal/xrand"
)

func TestIntN(t *testing.T) {
	t.Parallel()
	n := 100
	for i := 0; i < 1000; i++ {
		val := xrand.IntN(n)
		if val < 0 || val >= n {
			t.Errorf("IntN(%d) returned %d, which is out of range [0, %d)", n, val, n)
		}
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()
	for i := 0; i < 1000; i++ {
		val := xrand.Int64()
		// Int64 can return 0, so we just check if it doesn't panic and is a valid int64
		_ = val
	}
}
