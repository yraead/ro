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
	"context"
	"fmt"
)

func ExampleContextWithValue() {
	type contextValue struct{}

	observable := Pipe2(
		Just(1, 2, 3, 4, 5),
		ContextWithValue[int](contextValue{}, 42),
		Filter(func(i int) bool {
			return i%2 == 0
		}),
	)

	subscription := observable.Subscribe(
		OnNextWithContext(func(ctx context.Context, value int) {
			fmt.Printf("Next: %v\n", value)
			fmt.Printf("Next context value: %v\n", ctx.Value(contextValue{}))
		}),
	)
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Next context value: 42
	// Next: 4
	// Next context value: 42
}
