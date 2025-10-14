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
	_ "unsafe" // required for runtime.nanotime
)

//go:linkname nanotime runtime.nanotime
func nanotime() int64

// NowNanoMonotonic returns the current time in nanoseconds.
// It is approximately 3 times faster than time.Now() for high-frequency operations.
// This function uses runtime.nanotime() for better performance.
func NowNanoMonotonic() int64 {
	return nanotime()
}

// Using go:linkname is against the Go rules. There is another way to mesure the
// duration with monotonic time: using time.Since(startTime) where startTime is
// the program start time.
// This method is 1ns slower than calling nanotime(), which is not a big deal, but
// the improvement in code quality is not worth it.
//
// If the go:linkname directive become an issue in the future, please uncomment
// the following code, open a pull-request and explain why you did it.
//
// Follow-up: https://github.com/samber/hot/issues/39

// var startTime = time.Now()
//
// func NowNano() int64 {
// 	return time.Since(startTime).Nanoseconds()
// }
