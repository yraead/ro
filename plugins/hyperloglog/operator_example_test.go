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

package rohyperloglog

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"

	"github.com/samber/ro"
)

func ExampleCountDistinct() {
	// Count distinct strings using hyperloglog
	observable := ro.Pipe1(
		ro.Just("alice", "bob", "charlie", "alice", "bob", "david"),
		CountDistinct[string](8, false, func(input string) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 4
	// Completed
}

func ExampleCountDistinct_withStructs() {
	// Count distinct structs using hyperloglog
	type User struct {
		ID   int
		Name string
	}

	observable := ro.Pipe1(
		ro.Just(
			User{ID: 1, Name: "Alice"},
			User{ID: 2, Name: "Bob"},
			User{ID: 1, Name: "Alice"}, // Duplicate
			User{ID: 3, Name: "Charlie"},
			User{ID: 2, Name: "Bob"}, // Duplicate
		),
		CountDistinct[User](10, false, func(input User) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input.Name))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 3
	// Completed
}

func ExampleCountDistinct_withPrecision() {
	// Count distinct with different precision levels
	observable := ro.Pipe1(
		ro.Just("a", "b", "c", "d", "e", "f", "g", "h", "i", "j"),
		CountDistinct[string](4, false, func(input string) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Completed
}

func ExampleCountDistinct_withSparse() {
	// Count distinct with sparse representation
	observable := ro.Pipe1(
		ro.Just("alice", "bob", "charlie", "david", "eve"),
		CountDistinct[string](8, true, func(input string) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 5
	// Completed
}

func ExampleCountDistinct_withSHA256() {
	// Count distinct using SHA256 hash
	observable := ro.Pipe1(
		ro.Just("alice", "bob", "charlie", "alice", "bob", "david"),
		CountDistinct[string](12, false, func(input string) uint64 {
			hash := sha256.Sum256([]byte(input))
			return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
				uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 4
	// Completed
}

func ExampleCountDistinctReduce() {
	// Count distinct with incremental updates
	observable := ro.Pipe1(
		ro.Just("alice", "bob", "charlie", "alice", "bob", "david"),
		CountDistinctReduce[string](8, false, func(input string) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 3
	// Next: 3
	// Next: 4
	// Completed
}

func ExampleCountDistinctReduce_withStructs() {
	// Count distinct structs with incremental updates
	type User struct {
		ID   int
		Name string
	}

	observable := ro.Pipe1(
		ro.Just(
			User{ID: 1, Name: "Alice"},
			User{ID: 2, Name: "Bob"},
			User{ID: 1, Name: "Alice"}, // Duplicate
			User{ID: 3, Name: "Charlie"},
			User{ID: 2, Name: "Bob"}, // Duplicate
		),
		CountDistinctReduce[User](10, false, func(input User) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input.Name))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 2
	// Next: 3
	// Next: 3
	// Completed
}

func ExampleCountDistinct_withError() {
	// Count distinct with potential errors
	// Note: This demonstrates what happens with invalid precision

	// This will panic when the operator is created
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic caught:", r)
		}
	}()

	_ = CountDistinct[string](20, false, func(input string) uint64 {
		// Invalid precision (should be 4-18)
		h := fnv.New64a()
		h.Write([]byte(input))
		return h.Sum64()
	})

	// Output: Panic caught: rohyperloglog.CountDistinct: precision has to be >= 4 and <= 18
}

func ExampleCountDistinct_withLargeDataset() {
	// Count distinct in a large dataset
	observable := ro.Pipe1(
		ro.Just("user1", "user2", "user3", "user4", "user5", "user6", "user7", "user8", "user9", "user10"),
		CountDistinct[string](16, false, func(input string) uint64 {
			h := fnv.New64a()
			h.Write([]byte(input))
			return h.Sum64()
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[uint64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 2
	// Completed
}
