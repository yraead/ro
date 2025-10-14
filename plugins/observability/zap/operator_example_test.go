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


package rozap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/ro"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// createTestLogger creates a zap logger configured for testing with consistent output
func createTestLogger(level zapcore.Level) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("2024-01-01T12:00:00.000Z")
	}
	config.EncoderConfig.CallerKey = ""
	config.EncoderConfig.FunctionKey = ""
	config.EncoderConfig.StacktraceKey = ""
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stdout"}
	logger, _ := config.Build()
	return logger
}

func ExampleLog() {
	// Initialize zap logger with custom config to match expected output
	logger := createTestLogger(zapcore.InfoLevel)

	// Log all notifications (Next, Error, Complete)
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](logger, zapcore.InfoLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	logger.Sync()

	// Output:
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 1
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 2
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 3
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 4
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 5
	// 2024-01-01T12:00:00.000Z	INFO	ro.Complete
}

func ExampleLogWithNotification() {
	// Initialize zap logger with custom config to match expected output
	logger := createTestLogger(zapcore.DebugLevel)

	// Log with structured notification data
	observable := ro.Pipe1(
		ro.Just("hello", "world", "golang"),
		LogWithNotification[string](logger, zapcore.DebugLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[string]())
	defer subscription.Unsubscribe()

	logger.Sync()

	// Output:
	// 2024-01-01T12:00:00.000Z	DEBUG	ro.Next	{"value": "hello"}
	// 2024-01-01T12:00:00.000Z	DEBUG	ro.Next	{"value": "world"}
	// 2024-01-01T12:00:00.000Z	DEBUG	ro.Next	{"value": "golang"}
	// 2024-01-01T12:00:00.000Z	DEBUG	ro.Complete
}

func ExampleLog_withError() {
	// Initialize zap logger with custom config to match expected output
	logger := createTestLogger(zapcore.DebugLevel)

	// Log including error notifications
	observable := ro.Pipe1(
		ro.NewObservable(func(observer ro.Observer[int]) ro.Teardown {
			observer.Next(1)
			observer.Next(2)
			observer.Error(errors.New("something went wrong"))
			observer.Next(3) // This won't be emitted due to error
			return nil
		}),
		Log[int](logger, zapcore.ErrorLevel),
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(value int) {
			// Consume values to trigger logging
		},
		func(err error) {
			// Handle errors
		},
		func() {
			// Handle completion
		},
	))
	defer subscription.Unsubscribe()

	logger.Sync()

	// Output:
	// 2024-01-01T12:00:00.000Z	ERROR	ro.Next: 1
	// 2024-01-01T12:00:00.000Z	ERROR	ro.Next: 2
	// 2024-01-01T12:00:00.000Z	ERROR	ro.Error: something went wrong
}

func ExampleLog_inPipeline() {
	// Initialize zap logger with custom config to match expected output
	logger := createTestLogger(zapcore.DebugLevel)

	// Use logging in a complex pipeline
	observable := ro.Pipe3(
		ro.Just(1, 2, 3, 4, 5),
		ro.Filter(func(n int) bool { return n%2 == 0 }), // Keep even numbers
		Log[int](logger, zapcore.InfoLevel),
		ro.Map(func(n int) string { return fmt.Sprintf("Even: %d", n) }),
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(value string) {
			// Consume values to trigger logging
		},
		func(err error) {
			// Handle errors
		},
		func() {
			// Handle completion
		},
	))
	defer subscription.Unsubscribe()

	logger.Sync()

	// Output:
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 2
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next: 4
	// 2024-01-01T12:00:00.000Z	INFO	ro.Complete
}

func ExampleLog_withContext() {
	// Initialize zap logger with custom config to match expected output
	logger := createTestLogger(zapcore.DebugLevel)

	// Log with context-aware operations
	ctx := context.Background()

	observable := ro.Pipe1(
		ro.Just("context", "aware", "logging"),
		LogWithNotification[string](logger, zapcore.InfoLevel),
	)

	subscription := observable.SubscribeWithContext(ctx, ro.NewObserverWithContext(
		func(ctx context.Context, value string) {
			// Consume values to trigger logging
		},
		func(ctx context.Context, err error) {
			// Handle errors
		},
		func(ctx context.Context) {
			// Handle completion
		},
	))
	defer subscription.Unsubscribe()

	logger.Sync()

	// Output:
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next	{"value": "context"}
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next	{"value": "aware"}
	// 2024-01-01T12:00:00.000Z	INFO	ro.Next	{"value": "logging"}
	// 2024-01-01T12:00:00.000Z	INFO	ro.Complete
}

func ExampleLog_withCustomLevels() {
	// Initialize zap logger with custom config to match expected output
	logger := createTestLogger(zapcore.DebugLevel)

	// Demonstrate different log levels
	observable := ro.Pipe1(
		ro.Just(1, 2, 3, 4, 5),
		Log[int](logger, zapcore.WarnLevel),
	)

	subscription := observable.Subscribe(ro.NoopObserver[int]())
	defer subscription.Unsubscribe()

	logger.Sync()

	// Output:
	// 2024-01-01T12:00:00.000Z	WARN	ro.Next: 1
	// 2024-01-01T12:00:00.000Z	WARN	ro.Next: 2
	// 2024-01-01T12:00:00.000Z	WARN	ro.Next: 3
	// 2024-01-01T12:00:00.000Z	WARN	ro.Next: 4
	// 2024-01-01T12:00:00.000Z	WARN	ro.Next: 5
	// 2024-01-01T12:00:00.000Z	WARN	ro.Complete
}
