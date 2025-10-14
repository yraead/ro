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
	"context"
	"fmt"
	"log"
)

var (
	// By default, the library will ignore unhandled errors and dropped notifications.
	// You can change this behavior by setting the following variables to your own
	// error handling functions.
	//
	// Example:
	//
	// 	ro.OnUnhandledError = func(ctx context.Context, err error) {
	// 		slog.Error(fmt.Sprintf("unhandled error: %s\n", err.Error()))
	// 	}
	//
	// 	ro.OnDroppedNotification = func(ctx context.Context, notification fmt.Stringer) {
	// 		slog.Warn(fmt.Sprintf("dropped notification: %s\n", notification.String()))
	// 	}
	//
	// Note: `OnUnhandledError` and `OnDroppedNotification` are called synchronously from
	// the goroutine that emits the error or the notification. A slow callback will slow
	// down the whole pipeline.

	// OnUnhandledError is called when an error is emitted by an Observable and
	// no error handler is registered.
	OnUnhandledError = IgnoreOnUnhandledError
	// OnDroppedNotification is called when a notification is emitted by an Observable and
	// no notification handler is registered.
	OnDroppedNotification = IgnoreOnDroppedNotification
)

// IgnoreOnUnhandledError is the default implementation of `OnUnhandledError`.
func IgnoreOnUnhandledError(ctx context.Context, err error) {}

// IgnoreOnDroppedNotification is the default implementation of `OnDroppedNotification`.
func IgnoreOnDroppedNotification(ctx context.Context, notification fmt.Stringer) {}

// DefaultOnUnhandledError is the default implementation of `OnUnhandledError`.
func DefaultOnUnhandledError(ctx context.Context, err error) {
	if err != nil {
		// bearer:disable go_lang_logger_leak
		log.Printf("samber/ro: unhandled error: %s\n", err.Error())
	}
}

var _ fmt.Stringer = (*Notification[int])(nil) // see below

// DefaultOnDroppedNotification is the default implementation of `OnDroppedNotification`.
//
// Since we cannot assign a generic callback to `OnDroppedNotification`,
// we had to use a `fmt.Stringer` instead a `Notification[T any]`.
func DefaultOnDroppedNotification(ctx context.Context, notification fmt.Stringer) {
	// bearer:disable go_lang_logger_leak
	log.Printf("samber/ro: dropped notification: %s\n", notification.String())
}
