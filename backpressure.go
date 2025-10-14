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

// Backpressure is a type that represents the backpressure strategy to use.
type Backpressure int8

const (
	// BackpressureBlock blocks the source observable when the destination is not ready to receive more values.
	BackpressureBlock Backpressure = iota
	// BackpressureDrop drops the source observable when the destination is not ready to receive more values.
	BackpressureDrop
)

// ConcurrencyMode is a type that represents the concurrency mode to use.
type ConcurrencyMode int8

// ConcurrencyModeSafe is a concurrency mode that is safe to use.
// Spinlock is ignored because it is too slow when chaining operators. Spinlock should be used
// only for short-lived local locks.
const (
	ConcurrencyModeSafe ConcurrencyMode = iota
	ConcurrencyModeUnsafe
	ConcurrencyModeEventualySafe
)
