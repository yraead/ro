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
	"fmt"

	"github.com/samber/ro"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log[T any](logger *zap.Logger, level zapcore.Level) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapWithContext(
		func(ctx context.Context, value T) {
			logger.Log(level, fmt.Sprintf("ro.Next: %v", value))
		},
		func(ctx context.Context, err error) {
			logger.Log(level, fmt.Sprintf("ro.Error: %s", err.Error()))
		},
		func(ctx context.Context) {
			logger.Log(level, "ro.Complete")
		},
	)
}

func LogWithNotification[T any](logger *zap.Logger, level zapcore.Level) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapWithContext(
		func(ctx context.Context, value T) {
			logger.Log(level, "ro.Next", zap.Any("value", value))
		},
		func(ctx context.Context, err error) {
			logger.Log(level, "ro.Error", zap.Error(err))
		},
		func(ctx context.Context) {
			logger.Log(level, "ro.Complete")
		},
	)
}

func FatalOnError[T any](logger *zap.Logger) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapOnErrorWithContext[T](
		func(ctx context.Context, err error) {
			logger.Fatal("ro.Error", zap.Error(err))
		},
	)
}
