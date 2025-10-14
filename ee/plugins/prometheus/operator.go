// Copyright 2025 samber.
//
// Licensed as an Enterprise License (the "License"); you may not use
// this file except in compliance with the License. You may obtain
// a copy of the License at:
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.ee.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package roprometheus

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/ro"
	"github.com/samber/ro/internal/xtime"
)

// IncCounterOnNext is a pipe operator that increments a counter
// when a new Next() notification is sent to the destination observer.
//
// It adds a short lock for each Next() notification: 30ns for the
// prometheus/client_golang locks.
func IncCounterOnNext[T any](counter prometheus.Counter) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isPrometheusEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						counter.Inc()
						destination.NextWithContext(ctx, value)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

			return sub.Unsubscribe
		})
	}
}

// IncCounterOnError is a pipe operator that increments a counter
// when a new Error() notification is sent to the destination observer.
//
// It adds a short lock for each Error() notification: 30ns for the
// prometheus/client_golang locks.
func IncCounterOnError[T any](counter prometheus.Counter) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isPrometheusEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					func(ctx context.Context, err error) {
						counter.Inc()
						destination.ErrorWithContext(ctx, err)
					},
					destination.CompleteWithContext,
				),
			)

			return sub.Unsubscribe
		})
	}
}

// IncCounterOnComplete is a pipe operator that increments a counter
// when a new Complete() notification is sent to the destination observer.
//
// It adds a short lock for each Complete() notification: 30ns for the
// prometheus/client_golang locks.
func IncCounterOnComplete[T any](counter prometheus.Counter) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isPrometheusEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					destination.ErrorWithContext,
					func(ctx context.Context) {
						counter.Inc()
						destination.CompleteWithContext(ctx)
					},
				),
			)

			return sub.Unsubscribe
		})
	}
}

// IncCounterOnSubscription is a pipe operator that increments a counter
// when a new subscription is created.
//
// It adds a short lock for each subscription: 30ns for the
// prometheus/client_golang locks.
func IncCounterOnSubscription[T any](counter prometheus.Counter) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isPrometheusEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			counter.Inc()

			sub := source.SubscribeWithContext(subscriberCtx, destination)
			return sub.Unsubscribe
		})
	}
}

// ObserveNextLag is a pipe operator that tracks the time it takes for a notification
// to traverse from the source observable to the destination observer.
// It mesures the time the source pauses while waiting for the destination
// to process the notification.
//
// It adds a short lock for each Next() notification:
//   - 2x 15ns for the time tracking
//   - 30ns for the prometheus/client_golang locks
func ObserveNextLag[T any](summaryOrHistogram prometheus.Observer) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isPrometheusEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						start := xtime.NowNanoMonotonic()
						destination.NextWithContext(ctx, value)
						end := xtime.NowNanoMonotonic()

						summaryOrHistogram.Observe(float64(end-start) / 1e9)
					},
					// @TODO: track error and completion processing time?
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)
			return sub.Unsubscribe
		})
	}
}

type checkpointCtx struct{}

// observeBeforePipe is the aggregation of the following operators:
//   - IncCounterOnNext
//   - ObserveNextLag
//
// Aggregating avoid the need to add a short lock for each notification.
func observeBeforePipe[T any](counterOnNext prometheus.Counter, summaryOrHistogram prometheus.Observer) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						// Operator: Operator: IncCounterOnNext
						counterOnNext.Inc()

						/////////////////////////////

						// Operator: observeOnNotification
						start := xtime.NowNanoMonotonic()
						ctx = context.WithValue(ctx, checkpointCtx{}, start)

						/////////////////////////////

						// Operator: ObserveNextLag
						destination.NextWithContext(ctx, value)

						end := xtime.NowNanoMonotonic()
						summaryOrHistogram.Observe(float64(end-start) / 1e9)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)
			return sub.Unsubscribe
		})
	}
}

// observeOperatorProcessingTime is a pipe operator that observes the processing time of an operator.
func observeOperatorProcessingTime[T any](summaryOrHistogram prometheus.ObserverVec, operatorName string, operatorPosition string, operatorIndex int) func(ro.Observable[T]) ro.Observable[T] {
	prometheusObserver := summaryOrHistogram.With(prometheus.Labels{
		labelNameOperator:         operatorName,
		labelNameOperatorPosition: operatorPosition,
		labelNameOperatorIndex:    strconv.Itoa(operatorIndex),
	})

	return func(source ro.Observable[T]) ro.Observable[T] {
		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						start, ok := ctx.Value(checkpointCtx{}).(int64)
						end := xtime.NowNanoMonotonic()

						if ok {
							prometheusObserver.Observe(float64(end-start) / 1e9)
						}

						ctx = context.WithValue(ctx, checkpointCtx{}, end)
						destination.NextWithContext(ctx, value)
					},
					// @TODO: track error and completion processing time?
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)
			return sub.Unsubscribe
		})
	}
}

// observeAfterPipe is the aggregation of the following operators:
//   - IncCounterOnNext
//   - IncCounterOnSubscription
//
// Aggregating avoid the need to add a short lock for each notification.
func observeAfterPipe[T any](counterOnNext prometheus.Counter, counterOnSubscription prometheus.Counter) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			// Operator: IncCounterOnSubscription
			counterOnSubscription.Inc()

			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						// Operator: IncCounterOnNext
						counterOnNext.Inc()

						destination.NextWithContext(ctx, value)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)
			return sub.Unsubscribe
		})
	}
}
