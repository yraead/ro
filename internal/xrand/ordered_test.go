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
	"reflect"
	"sort"
	"testing"

	"github.com/samber/ro/internal/xrand"
)

func TestShuffle(t *testing.T) {
	t.Parallel()
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := make([]int, len(s))
	copy(shuffled, s)

	xrand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// It's hard to test randomness deterministically, but we can check if it's not identical to original (unlikely for small slice)
	// This test might fail rarely due to pure chance, but it's good enough for basic coverage.
	notShuffled := true

	for i := range s {
		if s[i] != shuffled[i] {
			notShuffled = false
			break
		}
	}

	if notShuffled {
		t.Fatal("Warning: Slice was not shuffled (pure chance or issue).")
	}

	// Check if the shuffled slice contains the same elements as the original
	sort.Ints(shuffled)

	if !reflect.DeepEqual(s, shuffled) {
		t.Errorf("Shuffle did not preserve elements. Original: %v, Shuffled: %v", s, shuffled)
	}
}
