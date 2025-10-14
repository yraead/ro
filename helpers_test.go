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

package ro

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/samber/lo"
)

// https://github.com/stretchr/testify/issues/1101
func testWithTimeout(t *testing.T, timeout time.Duration) {
	t.Helper()

	testFinished := make(chan struct{})

	t.Cleanup(func() { close(testFinished) })

	go func() {
		select {
		case <-testFinished:
		case <-time.After(timeout):
			t.Errorf("test timed out after %s", timeout)
			os.Exit(1)
		}
	}()
}

func passThrough[T any]() func(Observable[T]) Observable[T] {
	return func(observable Observable[T]) Observable[T] {
		return observable
	}
}

func syncMapLength(m *sync.Map) int {
	size := 0

	m.Range(func(key, value any) bool {
		size++
		return true
	})

	return size
}

func t2ToSliceB[A, B any](slice []lo.Tuple2[A, B]) []B {
	return lo.Map(slice, func(t lo.Tuple2[A, B], _ int) B {
		return t.B
	})
}
