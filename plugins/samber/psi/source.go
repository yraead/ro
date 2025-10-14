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

package ropsi

import (
	"context"
	"time"

	psinotifier "github.com/samber/go-psi"
	"github.com/samber/ro"
)

// NewPSINotifier creates an observable that emits PSI (Pressure Stall Information) statistics at regular intervals.
func NewPSINotifier(interval time.Duration) ro.Observable[psinotifier.PSIStatsResource] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[psinotifier.PSIStatsResource]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := psinotifier.AllPSIStats()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}
