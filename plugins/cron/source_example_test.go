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


package rocron

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/samber/ro"
)

// extractCounter operator that extracts just the counter from ScheduleJob
func extractCounter() func(ro.Observable[ScheduleJob]) ro.Observable[int] {
	return func(source ro.Observable[ScheduleJob]) ro.Observable[int] {
		return ro.Map(func(job ScheduleJob) int {
			return job.Counter
		})(source)
	}
}

func ExampleNewScheduler_everySecond() {
	// Create a scheduler that emits every 50ms for testing
	observable := ro.Pipe1(
		NewScheduler(
			gocron.DurationJob(50*time.Millisecond),
		),
		extractCounter(),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Wait for a few events to be emitted
	time.Sleep(175 * time.Millisecond)

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
}

func ExampleNewScheduler_everyMinute() {
	// Create a scheduler that emits every 100ms for testing
	observable := ro.Pipe1(
		NewScheduler(
			gocron.DurationJob(100*time.Millisecond),
		),
		extractCounter(),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Wait for a few events to be emitted
	time.Sleep(325 * time.Millisecond)

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
}

func ExampleNewScheduler_dailyAtSpecificTime() {
	// Create a scheduler that emits every 75ms for testing
	observable := ro.Pipe1(
		NewScheduler(
			gocron.DurationJob(75*time.Millisecond),
		),
		extractCounter(),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Wait for a few events to be emitted
	time.Sleep(250 * time.Millisecond)

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
}

func ExampleNewScheduler_weeklyOnMonday() {
	// Create a scheduler that emits every 125ms for testing
	observable := ro.Pipe1(
		NewScheduler(
			gocron.DurationJob(125*time.Millisecond),
		),
		extractCounter(),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Wait for a few events to be emitted
	time.Sleep(400 * time.Millisecond)

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
}

func ExampleNewScheduler_monthlyOnFirstDay() {
	// Create a scheduler that emits every 150ms for testing
	observable := ro.Pipe1(
		NewScheduler(
			gocron.DurationJob(150*time.Millisecond),
		),
		extractCounter(),
	)

	subscription := observable.Subscribe(ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Wait for a few events to be emitted
	time.Sleep(475 * time.Millisecond)

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
}

func ExampleNewScheduler_withContext() {
	// Create a scheduler with context for cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 225*time.Millisecond)
	defer cancel()

	observable := ro.Pipe1(
		NewScheduler(
			gocron.DurationJob(50*time.Millisecond),
		),
		extractCounter(),
	)

	subscription := observable.SubscribeWithContext(ctx, ro.PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Wait for context to timeout
	time.Sleep(300 * time.Millisecond)

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: context deadline exceeded
}

func ExampleNewScheduler_withProcessing() {
	// Create a scheduler and process the events
	observable := ro.Pipe3(
		NewScheduler(
			gocron.DurationJob(50*time.Millisecond),
		),
		extractCounter(),
		ro.Map(func(counter int) string {
			return fmt.Sprintf("Scheduled job #%d executed", counter)
		}),
		ro.Take[string](3), // Only take first 3 events
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Wait for events to be processed
	time.Sleep(200 * time.Millisecond)

	// Output:
	// Next: Scheduled job #0 executed
	// Next: Scheduled job #1 executed
	// Next: Scheduled job #2 executed
	// Completed
}
