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
	"context"

	"github.com/rs/zerolog"
	"github.com/samber/ro"
)

func Log[T any](logger *zerolog.Logger, level zerolog.Level) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapWithContext(
		func(ctx context.Context, value T) {
			logger.WithLevel(level).Msgf("ro.Next: %v", value)
		},
		func(ctx context.Context, err error) {
			logger.WithLevel(level).Msgf("ro.Error: %s", err.Error())
		},
		func(ctx context.Context) {
			logger.WithLevel(level).Msgf("ro.Complete")
		},
	)
}

func LogWithNotification[T any](logger *zerolog.Logger, level zerolog.Level) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapWithContext(
		func(ctx context.Context, value T) {
			logger.WithLevel(level).Any("value", value).Msgf("ro.Next")
		},
		func(ctx context.Context, err error) {
			logger.WithLevel(level).Err(err).Msgf("ro.Error")
		},
		func(ctx context.Context) {
			logger.WithLevel(level).Msgf("ro.Complete")
		},
	)
}

func FatalOnError[T any](logger *zerolog.Logger) func(ro.Observable[T]) ro.Observable[T] {
	return ro.TapOnErrorWithContext[T](
		func(ctx context.Context, err error) {
			logger.Fatal().Err(err).Msgf("ro.Error")
		},
	)
}
