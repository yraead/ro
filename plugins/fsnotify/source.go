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

package rofsnotify

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/ro"
)

// NewFSListener creates a file system watcher that emits file system events.
func NewFSListener(paths ...string) ro.Observable[fsnotify.Event] {
	return ro.NewUnsafeObservableWithContext(func(ctx context.Context, destination ro.Observer[fsnotify.Event]) ro.Teardown {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			destination.ErrorWithContext(ctx, err)
			return nil
		}

		// Start listening for events.
		go func() {
			for _, path := range paths {
				// Add a path.
				err = watcher.Add(path)
				if err != nil {
					destination.ErrorWithContext(ctx, err)
					return
				}
			}

			for {
				select {
				case event, ok := <-watcher.Events:
					if ok {
						destination.NextWithContext(ctx, event)
					} else {
						destination.CompleteWithContext(ctx)
						return
					}

				case err, ok := <-watcher.Errors:
					if ok {
						destination.ErrorWithContext(ctx, err)
					} else {
						destination.CompleteWithContext(ctx)
					}
					return

				case <-ctx.Done():
					if ctx.Err() != nil {
						destination.ErrorWithContext(ctx, ctx.Err())
					} else {
						destination.CompleteWithContext(ctx)
					}
					return
				}
			}
		}()

		return func() {
			_ = watcher.Close()
		}
	})
}
