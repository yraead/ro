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
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/samber/ro"
)

func Log[T any](logger *sentry.Hub, level sentry.Level) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapWithContext(
		func(ctx context.Context, value T) {
			event := sentry.NewEvent()
			event.Level = level
			event.Message = fmt.Sprintf("ro.Next: %v", value)
			logger.CaptureEvent(event)
		},
		func(ctx context.Context, err error) {
			event := sentry.NewEvent()
			event.Level = level
			event.Message = "ro.Error: " + err.Error()
			event.SetException(err, 10)
			logger.CaptureEvent(event)
		},
		func(ctx context.Context) {
			event := sentry.NewEvent()
			event.Level = level
			event.Message = "ro.Complete"
			logger.CaptureEvent(event)
		},
	)
}

func LogWithNotification[T any](logger *sentry.Hub, level sentry.Level) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapWithContext(
		func(ctx context.Context, value T) {
			event := sentry.NewEvent()
			event.Level = level
			event.Message = "ro.Next"
			event.Extra["value"] = value
			logger.CaptureEvent(event)
		},
		func(ctx context.Context, err error) {
			event := sentry.NewEvent()
			event.Level = level
			event.Message = "ro.Error"
			event.Extra["error"] = err
			event.SetException(err, 10)
			logger.CaptureEvent(event)
		},
		func(ctx context.Context) {
			event := sentry.NewEvent()
			event.Level = level
			event.Message = "ro.Complete"
			logger.CaptureEvent(event)
		},
	)
}
