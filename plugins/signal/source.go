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
	"os/signal"

	"github.com/samber/ro"
)

const IOReaderBufferSize = 1024

// Notify causes package signal to relay incoming signals to c.
// If no signals are provided, all incoming signals will be relayed to c.
func NewSignalCatcher(signals ...os.Signal) ro.Observable[os.Signal] {
	return ro.NewUnsafeObservableWithContext(func(ctx context.Context, destination ro.Observer[os.Signal]) ro.Teardown {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, signals...)

		go func() {
			for sig := range ch {
				destination.NextWithContext(ctx, sig)
			}
			destination.CompleteWithContext(ctx)
		}()

		return func() {
			signal.Stop(ch)
			close(ch)
		}
	})
}
