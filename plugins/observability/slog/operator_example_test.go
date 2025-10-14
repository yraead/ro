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


package rolog

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/samber/ro"
)

// timeFilterWriter filters out time from slog output
type timeFilterWriter struct {
	w io.Writer
}

func (w *timeFilterWriter) Write(p []byte) (n int, err error) {
	// Remove time=... from the output
	line := string(p)
	if strings.Contains(line, "time=") {
		// Find the position after the time field
		timeIndex := strings.Index(line, "time=")
		if timeIndex != -1 {
			// Find the end of the time field (next space or end of line)
			rest := line[timeIndex:]
			spaceIndex := strings.Index(rest, " ")
			if spaceIndex != -1 {
				// Remove time field and the following space
				line = line[:timeIndex] + rest[spaceIndex+1:]
			} else {
				// Time field is at the end
				line = line[:timeIndex]
			}
		}
	}
	return w.w.Write([]byte(line))
}

func ExampleLog() {
	// Initialize slog logger
	buff := bufio.NewWriter(os.Stdout)
	logger := slog.New(slog.NewTextHandler(&timeFilterWriter{w: buff}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	defer buff.Flush()

	// Log all notifications (Next, Error, Complete)
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](*logger, slog.LevelInfo),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// level=INFO msg="ro.Next: 1"
	// level=INFO msg="ro.Next: 2"
	// level=INFO msg="ro.Next: 3"
	// level=INFO msg="ro.Next: 4"
	// level=INFO msg="ro.Next: 5"
	// level=INFO msg=ro.Complete
}

func ExampleLogWithNotification() {
	// Initialize slog logger with mock handler that removes time
	buff := bufio.NewWriter(os.Stdout)
	logger := slog.New(slog.NewTextHandler(&timeFilterWriter{w: buff}, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	defer buff.Flush()

	// Log with structured notification data
	observable := ro.Pipe1(
		ro.Just("hello", "world", "golang"),
		LogWithNotification[string](*logger, slog.LevelDebug),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// level=DEBUG msg=ro.Next value=hello
	// level=DEBUG msg=ro.Next value=world
	// level=DEBUG msg=ro.Next value=golang
	// level=DEBUG msg=ro.Complete
}

func ExampleLog_withError() {
	// Initialize slog logger with mock handler that removes time
	buff := bufio.NewWriter(os.Stdout)
	logger := slog.New(slog.NewTextHandler(&timeFilterWriter{w: buff}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
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
		Log[int](*logger, slog.LevelError),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// level=ERROR msg="ro.Next: 1"
	// level=ERROR msg="ro.Next: 2"
	// level=ERROR msg="ro.Error: something went wrong"
}

func ExampleLog_inPipeline() {
	// Initialize slog logger with mock handler that removes time
	buff := bufio.NewWriter(os.Stdout)
	logger := slog.New(slog.NewTextHandler(&timeFilterWriter{w: buff}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	defer buff.Flush()

	// Use logging in a complex pipeline
	observable := ro.Pipe3(
		ro.Just(1, 2, 3, 4, 5),
		ro.Filter(func(n int) bool { return n%2 == 0 }), // Keep even numbers
		Log[int](*logger, slog.LevelInfo),
		ro.Map(func(n int) string { return fmt.Sprintf("Even: %d", n) }),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// level=INFO msg="ro.Next: 2"
	// level=INFO msg="ro.Next: 4"
	// level=INFO msg=ro.Complete
}

func ExampleLog_withContext() {
	// Initialize slog logger with mock handler that removes time
	buff := bufio.NewWriter(os.Stdout)
	logger := slog.New(slog.NewTextHandler(&timeFilterWriter{w: buff}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	defer buff.Flush()

	// Log with context-aware operations
	ctx := context.Background()

	observable := ro.Pipe1(
		ro.Just("context", "aware", "logging"),
		LogWithNotification[string](*logger, slog.LevelInfo),
	)

	subscription := observable.SubscribeWithContext(ctx, ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// level=INFO msg=ro.Next value=context
	// level=INFO msg=ro.Next value=aware
	// level=INFO msg=ro.Next value=logging
	// level=INFO msg=ro.Complete
}

func ExampleLog_withCustomLevels() {
	// Initialize slog logger with mock handler that removes time
	buff := bufio.NewWriter(os.Stdout)
	logger := slog.New(slog.NewTextHandler(&timeFilterWriter{w: buff}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	defer buff.Flush()

	// Demonstrate different log levels
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](*logger, slog.LevelWarn),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// level=WARN msg="ro.Next: 1"
	// level=WARN msg="ro.Next: 2"
	// level=WARN msg="ro.Next: 3"
	// level=WARN msg="ro.Next: 4"
	// level=WARN msg="ro.Next: 5"
	// level=WARN msg=ro.Complete
}
