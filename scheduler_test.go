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
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestNewScheduler(t *testing.T) {
	t.Parallel()
	// testWithTimeout(t, 100*time.Millisecond)
	is := assert.New(t)

	is.Panics(func() {
		NewScheduler()
	})
}

func TestOperatorSchedulerSubscribeOn(t *testing.T) { //nolint:paralleltest
	// t.Parallel()
	testWithTimeout(t, 400*time.Millisecond)
	is := assert.New(t)

	is.PanicsWithError(
		"ro.SubscribeOn: buffer size must be greater than 0",
		func() {
			_, _ = Collect(
				Pipe2(
					Just[int64](1, 2, 3),
					SubscribeOn[int64](-42),
					Map(func(x int64) int64 {
						time.Sleep(10 * time.Millisecond) // simulate slow processing
						return x
					}),
				),
			)
		},
	)

	values, err := Collect(
		Pipe2(
			Just[int64](1, 2, 3),
			SubscribeOn[int64](42),
			Map(func(x int64) int64 {
				time.Sleep(10 * time.Millisecond) // simulate slow processing
				return x
			}),
		),
	)
	is.Equal([]int64{1, 2, 3}, values)
	is.NoError(err)

	// check that either the upstream or downstream run in a goroutine
	mu := lo.Synchronize()
	order := []int64{}
	values, err = Collect(
		Pipe3(
			Range(1, 4),
			TapOnNext(func(value int64) {
				mu.Do(func() {
					order = append(order, value)
				})
			}),
			SubscribeOn[int64](42),
			TapOnNext(func(value int64) {
				time.Sleep(10 * time.Millisecond)
				mu.Do(func() {
					order = append(order, value*-1)
				})
			}),
		),
	)
	is.Equal([]int64{1, 2, 3}, values)
	is.Equal([]int64{1, 2, 3, -1, -2, -3}, order)
	is.NoError(err)

	// check that goroutine is used on downstream instead of upstream
	start := time.Now()
	obs := Pipe1(
		RangeWithInterval(0, 3, 50*time.Millisecond),
		SubscribeOn[int64](42),
	)
	sub := obs.Subscribe(NoopObserver[int64]())

	is.InDelta(150*time.Millisecond, time.Since(start), float64(15*time.Millisecond))
	is.True(sub.IsClosed())

	// @TODO: write some tests for channel buffer overflow
}

func TestOperatorSchedulerObserveOn(t *testing.T) { //nolint:paralleltest
	// t.Parallel()
	testWithTimeout(t, 700*time.Millisecond)
	is := assert.New(t)

	is.PanicsWithError(
		"ro.ObserveOn: buffer size must be greater than 0",
		func() {
			_, _ = Collect(
				Pipe2(
					Just[int64](1, 2, 3),
					ObserveOn[int64](-42),
					Map(func(x int64) int64 {
						time.Sleep(10 * time.Millisecond) // simulate slow processing
						return x
					}),
				),
			)
		},
	)

	values, err := Collect(
		Pipe2(
			Just[int64](1, 2, 3),
			ObserveOn[int64](42),
			Map(func(x int64) int64 {
				time.Sleep(10 * time.Millisecond) // simulate slow processing
				return x
			}),
		),
	)
	is.Equal([]int64{1, 2, 3}, values)
	is.NoError(err)

	// check that either the upstream or downstream run in a goroutine
	mu := lo.Synchronize()
	order := []int64{}
	values, err = Collect(
		Pipe3(
			Range(1, 4),
			TapOnNext(func(value int64) {
				mu.Do(func() {
					order = append(order, value)
				})
			}),
			ObserveOn[int64](42),
			TapOnNext(func(value int64) {
				time.Sleep(20 * time.Millisecond)
				mu.Do(func() {
					order = append(order, value*-1)
				})
			}),
		),
	)
	is.Equal([]int64{1, 2, 3}, values)
	is.Equal([]int64{1, 2, 3, -1, -2, -3}, order)
	is.NoError(err)

	// check that goroutine is used on downstream instead of upstream
	start := time.Now()
	obs := Pipe1(
		RangeWithInterval(0, 3, 50*time.Millisecond),
		ObserveOn[int64](42),
	)
	sub := obs.Subscribe(NoopObserver[int64]())

	is.InDelta(0, time.Since(start), float64(15*time.Millisecond))
	is.False(sub.IsClosed())
	sub.Wait() // Note: using .Wait() is not recommended.
	is.InDelta(150*time.Millisecond, time.Since(start), float64(15*time.Millisecond))
	is.True(sub.IsClosed())

	// @TODO: write some tests for channel buffer overflow
}
