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


package rosentry

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/samber/ro"
)

func ExampleLog() {
	// Initialize Sentry hub
	hub := createStdoutHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "observable")
	})

	// Log all notifications (Next, Error, Complete)
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](hub, sentry.LevelInfo),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Sentry event: level=info message="ro.Next: 1"
	// Sentry event: level=info message="ro.Next: 2"
	// Sentry event: level=info message="ro.Next: 3"
	// Sentry event: level=info message="ro.Next: 4"
	// Sentry event: level=info message="ro.Next: 5"
	// Sentry event: level=info message="ro.Complete"
}

func ExampleLogWithNotification() {
	// Initialize Sentry hub
	hub := createStdoutHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "observable")
	})

	// Log with structured notification data
	observable := ro.Pipe1(
		ro.Just("hello", "world", "golang"),
		LogWithNotification[string](hub, sentry.LevelDebug),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Sentry event: level=debug message="ro.Next" extra=map[value:hello]
	// Sentry event: level=debug message="ro.Next" extra=map[value:world]
	// Sentry event: level=debug message="ro.Next" extra=map[value:golang]
	// Sentry event: level=debug message="ro.Complete"
}

func ExampleLog_withError() {
	// Initialize Sentry hub
	hub := createStdoutHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "observable")
	})

	// Log including error notifications
	observable := ro.Pipe1(
		ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Error(errors.New("something went wrong"))
			observer.Next(3) // This won't be emitted due to error
			return nil
		}),
		Log[int](hub, sentry.LevelFatal),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Sentry event: level=fatal message="ro.Next: 1"
	// Sentry event: level=fatal message="ro.Next: 2"
	// Sentry event: level=fatal message="ro.Error: something went wrong" exception=error
}

func ExampleLog_inPipeline() {
	// Initialize Sentry hub
	hub := createStdoutHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "observable")
	})

	// Use logging in a complex pipeline
	observable := ro.Pipe3(
		ro.Just(1, 2, 3, 4, 5),
		ro.Filter(func(n int) bool { return n%2 == 0 }), // Keep even numbers
		Log[int](hub, sentry.LevelInfo),
		ro.Map(func(n int) string { return fmt.Sprintf("Even: %d", n) }),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Sentry event: level=info message="ro.Next: 2"
	// Sentry event: level=info message="ro.Next: 4"
	// Sentry event: level=info message="ro.Complete"
}

func ExampleLog_withContext() {
	// Initialize Sentry hub
	hub := createStdoutHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "observable")
	})

	// Log with context-aware operations
	ctx := context.Background()

	observable := ro.Pipe1(
		ro.Just("context", "aware", "logging"),
		LogWithNotification[string](hub, sentry.LevelInfo),
	)

	subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Sentry event: level=info message="ro.Next" extra=map[value:context]
	// Sentry event: level=info message="ro.Next" extra=map[value:aware]
	// Sentry event: level=info message="ro.Next" extra=map[value:logging]
	// Sentry event: level=info message="ro.Complete"
}

func ExampleLog_withCustomLevels() {
	// Initialize Sentry hub
	hub := createStdoutHub()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("component", "observable")
	})

	// Demonstrate different log levels
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](hub, sentry.LevelWarning),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Sentry event: level=warning message="ro.Next: 1"
	// Sentry event: level=warning message="ro.Next: 2"
	// Sentry event: level=warning message="ro.Next: 3"
	// Sentry event: level=warning message="ro.Next: 4"
	// Sentry event: level=warning message="ro.Next: 5"
	// Sentry event: level=warning message="ro.Complete"
}
