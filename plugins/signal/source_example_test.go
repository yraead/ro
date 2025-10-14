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


package rosignal

import (
	"context"
	"os"
	"syscall"

	"github.com/samber/ro"
)

func ExampleNewSignalCatcher() {
	// Catch all incoming signals
	observable := NewSignalCatcher()

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(signal os.Signal) {
				// Handle incoming signal
				switch signal {
				case syscall.SIGINT:
					// Handle Ctrl+C
				case syscall.SIGTERM:
					// Handle termination signal
				case syscall.SIGHUP:
					// Handle hangup signal
				}
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()
}

func ExampleNewSignalCatcher_withSpecificSignals() {
	// Catch specific signals only
	observable := NewSignalCatcher(
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Termination
		syscall.SIGHUP,  // Hangup
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(signal os.Signal) {
				// Handle specific signals
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()
}

func ExampleNewSignalCatcher_withFiltering() {
	// Catch all signals but filter for specific ones
	observable := ro.Pipe1(
		NewSignalCatcher(),
		ro.Filter(func(signal os.Signal) bool {
			// Only process SIGINT and SIGTERM
			return signal == syscall.SIGINT || signal == syscall.SIGTERM
		}),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(signal os.Signal) {
				// Handle filtered signals
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()
}

func ExampleNewSignalCatcher_withTransformation() {
	// Catch signals and transform them to string descriptions
	observable := ro.Pipe1(
		NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP),
		ro.Map(func(signal os.Signal) string {
			switch signal {
			case syscall.SIGINT:
				return "Interrupt signal received"
			case syscall.SIGTERM:
				return "Termination signal received"
			case syscall.SIGHUP:
				return "Hangup signal received"
			default:
				return "Unknown signal received"
			}
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()
}

func ExampleNewSignalCatcher_withErrorHandling() {
	// Catch signals with comprehensive error handling
	observable := NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(signal os.Signal) {
				// Handle successful signal reception
			},
			func(err error) {
				// Handle signal catching error
				// This could be due to:
				// - Insufficient permissions
				// - Signal not supported on platform
				// - Other system limitations
			},
			func() {
				// Handle completion (when signal catching stops)
			},
		),
	)
	defer subscription.Unsubscribe()
}

func ExampleNewSignalCatcher_withContext() {
	// Catch signals with context for cancellation
	observable := NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM)

	subscription := observable.SubscribeWithContext(
		context.Background(),
		ro.NewObserverWithContext(
			func(ctx context.Context, signal os.Signal) {
				// Handle signal with context
			},
			func(ctx context.Context, err error) {
				// Handle error with context
			},
			func(ctx context.Context) {
				// Handle completion with context
			},
		),
	)
	defer subscription.Unsubscribe()
}

func ExampleNewSignalCatcher_withGracefulShutdown() {
	// Catch signals for graceful shutdown
	observable := ro.Pipe1(
		NewSignalCatcher(syscall.SIGINT, syscall.SIGTERM),
		ro.Map(func(signal os.Signal) string {
			// Transform signal to shutdown action
			switch signal {
			case syscall.SIGINT:
				return "Graceful shutdown initiated by user"
			case syscall.SIGTERM:
				return "Graceful shutdown initiated by system"
			default:
				return "Unknown shutdown signal"
			}
		}),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(action string) {
				// Perform graceful shutdown
				// e.g., close connections, save state, etc.
			},
			func(err error) {
				// Handle error during shutdown
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()
}
