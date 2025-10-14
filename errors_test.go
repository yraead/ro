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
	"errors"
	"testing"
	"time"
)

func TestRecoverValueToError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string error",
			input:    "test error",
			expected: "unexpected error: test error",
		},
		{
			name:     "error type",
			input:    errors.New("test error"),
			expected: "test error",
		},
		{
			name:     "int value",
			input:    42,
			expected: "unexpected error: 42",
		},
		{
			name:     "nil value",
			input:    nil,
			expected: "unexpected error: <nil>",
		},
	}

	for _, tt := range tests {
		ttt := tt
		t.Run(ttt.name, func(t *testing.T) {
			t.Parallel()

			result := recoverValueToError(ttt.input)
			if result.Error() != ttt.expected {
				t.Errorf("recoverValueToError(%v) = %v, want %v", ttt.input, result.Error(), ttt.expected)
			}
		})
	}
}

func TestRecoverUnhandledError(t *testing.T) {
	t.Parallel()

	// Test that the function doesn't panic when callback panics
	t.Run("callback panics", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("recoverUnhandledError should not panic, got %v", r)
			}
		}()

		recoverUnhandledError(func() {
			panic("test panic")
		})
	})

	// Test that the function works normally when callback doesn't panic
	t.Run("callback doesn't panic", func(t *testing.T) {
		t.Parallel()
		called := false

		recoverUnhandledError(func() {
			called = true
		})

		if !called {
			t.Error("callback should have been called")
		}
	})
}

func TestErrorTypes(t *testing.T) {
	t.Parallel()

	t.Run("unsubscription error", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := newUnsubscriptionError(originalErr)

		if err.Error() != "ro.Subscription: original error" {
			t.Errorf("unsubscription error message = %v, want 'ro.Subscription: original error'", err.Error())
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("observable error", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := newObservableError(originalErr)

		if err.Error() != "ro.Observable: original error" {
			t.Errorf("observable error message = %v, want 'ro.Observable: original error'", err.Error())
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("observer error", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := newObserverError(originalErr)

		if err.Error() != "ro.Observer: original error" {
			t.Errorf("observer error message = %v, want 'ro.Observer: original error'", err.Error())
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("observer error with nil", func(t *testing.T) {
		t.Parallel()
		err := newObserverError(nil)

		if err.Error() != "ro.Observer: <nil>" {
			t.Errorf("observer error message = %v, want 'ro.Observer: <nil>'", err.Error())
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped != nil {
			t.Errorf("unwrapped error = %v, want nil", unwrapped)
		}
	})

	t.Run("timeout error", func(t *testing.T) {
		t.Parallel()
		duration := 5 * time.Second
		err := newTimeoutError(duration)

		expected := "ro.Timeout: timeout after 5s"
		if err.Error() != expected {
			t.Errorf("timeout error message = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("cast error", func(t *testing.T) {
		t.Parallel()
		err := newCastError[int, string]()

		expected := "ro.Cast: unable to cast int to string"
		if err.Error() != expected {
			t.Errorf("cast error message = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("pipe error", func(t *testing.T) {
		t.Parallel()
		err := newPipeError("test error: %s", "details")

		expected := "ro.Pipe: test error: details"
		if err.Error() != expected {
			t.Errorf("pipe error message = %v, want %v", err.Error(), expected)
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			t.Error("pipe error should have an unwrapped error")
		}
	})
}

func TestErrorUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("unsubscription error unwrap", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := &unsubscriptionError{err: originalErr}

		unwrapped := err.Unwrap()
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("observable error unwrap", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := &observableError{err: originalErr}

		unwrapped := err.Unwrap()
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("observer error unwrap", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := &observerError{err: originalErr}

		unwrapped := err.Unwrap()
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})

	t.Run("pipe error unwrap", func(t *testing.T) {
		t.Parallel()
		originalErr := errors.New("original error")
		err := &pipeError{err: originalErr}

		unwrapped := err.Unwrap()
		if unwrapped != originalErr {
			t.Errorf("unwrapped error = %v, want %v", unwrapped, originalErr)
		}
	})
}
