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


package rozerolog

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/ro"
)

func ExampleLog() {
	// Initialize zerolog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := zerolog.New(buff).With().Logger()
	defer buff.Flush()

	// Log all notifications (Next, Error, Complete)
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](&logger, zerolog.InfoLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// {"level":"info","message":"ro.Next: 1"}
	// {"level":"info","message":"ro.Next: 2"}
	// {"level":"info","message":"ro.Next: 3"}
	// {"level":"info","message":"ro.Next: 4"}
	// {"level":"info","message":"ro.Next: 5"}
	// {"level":"info","message":"ro.Complete"}
}

func ExampleLogWithNotification() {
	// Initialize zerolog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := zerolog.New(buff).With().Logger()
	defer buff.Flush()

	// Log with structured notification data
	observable := ro.Pipe1(
		ro.Just("hello", "world", "golang"),
		LogWithNotification[string](&logger, zerolog.DebugLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// {"level":"debug","value":"hello","message":"ro.Next"}
	// {"level":"debug","value":"world","message":"ro.Next"}
	// {"level":"debug","value":"golang","message":"ro.Next"}
	// {"level":"debug","message":"ro.Complete"}
}

func ExampleLog_withError() {
	// Initialize zerolog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := zerolog.New(buff).With().Logger()
	defer buff.Flush()

	// Log including error notifications
	observable := ro.Pipe1(
		ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Error(errors.New("something went wrong"))
			observer.Next(3) // This won't be emitted due to error
			return nil
		}),
		Log[int](&logger, zerolog.ErrorLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// {"level":"error","message":"ro.Next: 1"}
	// {"level":"error","message":"ro.Next: 2"}
	// {"level":"error","message":"ro.Error: something went wrong"}
}

func ExampleLog_inPipeline() {
	// Initialize zerolog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := zerolog.New(buff).With().Logger()
	defer buff.Flush()

	// Use logging in a complex pipeline
	observable := ro.Pipe3(
		ro.Just(1, 2, 3, 4, 5),
		ro.Filter(func(n int) bool { return n%2 == 0 }), // Keep even numbers
		Log[int](&logger, zerolog.InfoLevel),
		ro.Map(func(n int) string { return fmt.Sprintf("Even: %d", n) }),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// {"level":"info","message":"ro.Next: 2"}
	// {"level":"info","message":"ro.Next: 4"}
	// {"level":"info","message":"ro.Complete"}
}

func ExampleLog_withContext() {
	// Initialize zerolog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := zerolog.New(buff).With().Logger()
	defer buff.Flush()

	// Log with context-aware operations
	ctx := context.Background()

	observable := ro.Pipe1(
		ro.Just("context", "aware", "logging"),
		LogWithNotification[string](&logger, zerolog.InfoLevel),
	)

	subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// {"level":"info","value":"context","message":"ro.Next"}
	// {"level":"info","value":"aware","message":"ro.Next"}
	// {"level":"info","value":"logging","message":"ro.Next"}
	// {"level":"info","message":"ro.Complete"}
}

func ExampleLog_withCustomLevels() {
	// Initialize zerolog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := zerolog.New(buff).With().Logger()
	defer buff.Flush()

	// Demonstrate different log levels
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](&logger, zerolog.WarnLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// {"level":"warn","message":"ro.Next: 1"}
	// {"level":"warn","message":"ro.Next: 2"}
	// {"level":"warn","message":"ro.Next: 3"}
	// {"level":"warn","message":"ro.Next: 4"}
	// {"level":"warn","message":"ro.Next: 5"}
	// {"level":"warn","message":"ro.Complete"}
}

// func ExampleFatalOnError() {
// 	// Initialize zerolog logger
// 	buff := bufio.NewWriter(os.Stdout)
// 	logger := zerolog.New(buff).With().Logger()
// 	defer buff.Flush()

// 	// Fatal on error (this would terminate the program in real usage)
// 	observable := ro.Pipe1(
// 		ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
// 			observer.Next(1)
// 			observer.Next(2)
// 			observer.Error(errors.New("critical error"))
// 			return nil
// 		}),
// 		FatalOnError[int](&logger),
// 	)

// 	subscription := observable.Subscribe(ro.NoopObserver[int]())
// 	defer subscription.Unsubscribe()
// }
