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

package rootel

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/ro"
	"github.com/samber/ro/internal/xtime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// IncCounterOnNext is a pipe operator that increments a counter
// when a new Next() notification is sent to the destination observer.
func IncCounterOnNext[T any](counter metric.Int64Counter, attributes []attribute.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		metricAttributes := metricWithAttributes(attributes)

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						counter.Add(ctx, 1, metricAttributes)
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

// LogOnNext is a pipe operator that logs a message
// when a new Next() notification is sent to the destination observer.
func LogOnNext[T any](logger log.Logger, severity log.Severity, attributes []log.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						logAttributes := concatAttributes(attributes, []log.KeyValue{
							log.String("notification.value", fmt.Sprintf("%v", value)),
						})
						logger.Emit(ctx, newRecord("ro.Next(...)", severity, logAttributes...))

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
func IncCounterOnError[T any](counter metric.Int64Counter, attributes []attribute.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		metricAttributes := metricWithAttributes(attributes)

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					func(ctx context.Context, err error) {
						counter.Add(ctx, 1, metricAttributes)

						destination.ErrorWithContext(ctx, err)
					},
					destination.CompleteWithContext,
				),
			)

			return sub.Unsubscribe
		})
	}
}

// LogOnError is a pipe operator that logs a message
// when a new Error() notification is sent to the destination observer.
func LogOnError[T any](logger log.Logger, severity log.Severity, attributes []log.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					func(ctx context.Context, err error) {
						logAttributes := concatAttributes(attributes, []log.KeyValue{
							log.String("notification.error", err.Error()),
						})
						logger.Emit(ctx, newRecord(fmt.Sprintf("ro.Error(%v)", err), severity, logAttributes...))

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
func IncCounterOnComplete[T any](counter metric.Int64Counter, attributes []attribute.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		metricAttributes := metricWithAttributes(attributes)

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					destination.ErrorWithContext,
					func(ctx context.Context) {
						counter.Add(ctx, 1, metricAttributes)

						destination.CompleteWithContext(ctx)
					},
				),
			)

			return sub.Unsubscribe
		})
	}
}

// LogOnComplete is a pipe operator that logs a message
// when a new Complete() notification is sent to the destination observer.
func LogOnComplete[T any](logger log.Logger, severity log.Severity, attributes []log.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					destination.ErrorWithContext,
					func(ctx context.Context) {
						logger.Emit(ctx, newRecord("ro.Complete()", severity, attributes...))

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
func IncCounterOnSubscription[T any](counter metric.Int64Counter, attributes []attribute.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		metricAttributes := metricWithAttributes(attributes)

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			counter.Add(subscriberCtx, 1, metricAttributes)

			sub := source.SubscribeWithContext(
				subscriberCtx,
				destination,
			)

			return sub.Unsubscribe
		})
	}
}

// LogOnSubscription is a pipe operator that logs a message
// when a new subscription is created.
func LogOnSubscription[T any](logger log.Logger, severity log.Severity, attributes []log.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			logger.Emit(subscriberCtx, newRecord("ro.Subscribe(...)", severity, attributes...))

			sub := source.SubscribeWithContext(
				subscriberCtx,
				destination,
			)

			return func() {
				sub.Unsubscribe()
				logger.Emit(subscriberCtx, newRecord("ro.Unsubscribe()", severity, attributes...))
			}
		})
	}
}

// StartTraceOnSubscription is a pipe operator that create a new OTEL trace
// when a new subscription is created.
func StartTraceOnSubscription[T any](collector *otelCollector) func(ro.Observable[T]) ro.Observable[T] {
	traceAttributes := traceWithAttributes(collector.config.TraceAttributes)

	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			// Create a new OTEL trace for each subscription.
			ctx, span := collector.tracer.Start(
				subscriberCtx,
				"ro.Subscribe(...)",
				traceAttributes,
			)

			sub := source.SubscribeWithContext(
				ctx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

			return func() {
				sub.Unsubscribe()
				span.End()
			}
		})
	}
}

// TraceOnError is a pipe operator that records an error in the current OTEL trace.
func TraceOnError[T any](collector *otelCollector) func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					destination.NextWithContext,
					func(ctx context.Context, err error) {
						span := trace.SpanFromContext(ctx)
						span.SetStatus(codes.Error, "ro.Error(...)")
						span.SetAttributes(attribute.String("notification.error", err.Error()))
						span.RecordError(err)
						destination.ErrorWithContext(ctx, err)
					},
					destination.CompleteWithContext,
				),
			)

			return sub.Unsubscribe
		})
	}
}

// ObserveNextLag is a pipe operator that tracks the time it takes for a notification
// to traverse from the source observable to the destination observer.
// It mesures the time the source pauses while waiting for the destination
// to process the notification.
func ObserveNextLag[T any](tracer trace.Tracer, operatorName string, histogram metric.Float64Histogram, attributes []attribute.KeyValue) func(ro.Observable[T]) ro.Observable[T] {
	metricAttributes := metricWithAttributes(attributes)

	return func(source ro.Observable[T]) ro.Observable[T] {
		if !isOtelEnabled() {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						ctx, span := tracer.Start(ctx, operatorName+"/ro.Next(...)")
						defer span.End()

						start := xtime.NowNanoMonotonic()
						destination.NextWithContext(ctx, value)
						end := xtime.NowNanoMonotonic()

						histogram.Record(
							ctx,
							float64(end-start)/1e9,
							metricAttributes,
						)
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
//   - TraceOnError
//   - StartTraceOnSubscription
//
// Aggregating avoid the need to add a short lock for each notification.
func observeBeforePipe[T any](collector *otelCollector) func(ro.Observable[T]) ro.Observable[T] {
	metricAttributes := metricWithAttributes(collector.config.MetricAttributes)

	return func(source ro.Observable[T]) ro.Observable[T] {
		if !collector.config.EnableMetrics {
			return source
		}

		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						// Count the number of messages entering the instrumented ro.PipeX
						if collector.NotificationsInTotal != nil {
							collector.NotificationsInTotal.Add(ctx, 1, metricAttributes)
						}

						// Start the checkpoint for the next operator
						start := xtime.NowNanoMonotonic()
						if collector.OperatorProcessingTimeSeconds != nil {
							ctx = context.WithValue(ctx, checkpointCtx{}, start)
						}

						// Forward the event to the next operator
						destination.NextWithContext(ctx, value)

						// Measure the processing time of the operator, using the previous checkpoint
						end := xtime.NowNanoMonotonic()
						if collector.NotificationLagSeconds != nil {
							collector.NotificationLagSeconds.Record(
								ctx,
								float64(end-start)/1e9,
								metricAttributes,
							)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)
			return sub.Unsubscribe
		})
	}
}

// observeOnNotification creates an operator that measures the processing duration of operators using OTEL metrics
func observeOnNotification[T any](collector *otelCollector, operatorName string, operatorPosition string, operatorIndex int) func(ro.Observable[T]) ro.Observable[T] {
	metricAttributes := metricWithAttributes(
		collector.config.MetricAttributes,
		[]attribute.KeyValue{
			attribute.String(labelNameOperator, operatorName),
			attribute.String(labelNameOperatorPosition, operatorPosition),
			attribute.String(labelNameOperatorIndex, strconv.Itoa(operatorIndex)),
		},
	)
	logAttributes := []log.KeyValue{
		log.String(labelNameOperator, operatorName),
		log.String(labelNameOperatorPosition, operatorPosition),
		log.String(labelNameOperatorIndex, strconv.Itoa(operatorIndex)),
	}

	return func(source ro.Observable[T]) ro.Observable[T] {
		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						// Log the call to the operator
						if collector.config.EnableLogging {
							collector.logger.Emit(
								ctx,
								newRecord(
									operatorName+"/ro.Next(...)",
									collector.config.LogLevelNext,
									concatAttributes(
										logAttributes,
										[]log.KeyValue{
											log.String("notification.value", fmt.Sprintf("%v", value)),
										},
									)...,
								),
							)
						}

						// Trace the call to the operator
						ctx, span := collector.tracer.Start(ctx, operatorName+"/ro.Next(...)")
						span.SetAttributes(collector.config.TraceAttributes...)
						span.SetAttributes(
							attribute.String(labelNameOperator, operatorName),
							attribute.String(labelNameOperatorPosition, operatorPosition),
							attribute.String(labelNameOperatorIndex, strconv.Itoa(operatorIndex)),
						)
						defer span.End()

						if collector.config.EnableMetrics && collector.OperatorProcessingTimeSeconds != nil {
							// Measure the processing time of the operator, using the previous checkpoint
							start, ok := ctx.Value(checkpointCtx{}).(int64)
							end := xtime.NowNanoMonotonic()
							if ok {
								collector.OperatorProcessingTimeSeconds.Record(
									ctx,
									float64(end-start)/1e9,
									metricAttributes,
								)
							}

							// Update the checkpoint for the next operator
							ctx = context.WithValue(ctx, checkpointCtx{}, end)
						}

						// Forward the event to the next operator
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
func observeAfterPipe[T any](collector *otelCollector) func(ro.Observable[T]) ro.Observable[T] {
	traceAttributes := trace.WithAttributes(collector.config.TraceAttributes...)
	metricAttributes := metric.WithAttributes(collector.config.MetricAttributes...)

	return func(source ro.Observable[T]) ro.Observable[T] {
		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			// Instrument subscription
			if collector.config.EnableMetrics && collector.SubscriptionsTotal != nil {
				collector.SubscriptionsTotal.Add(subscriberCtx, 1, metricAttributes)
			}

			var span trace.Span
			if collector.config.EnableTracing {
				subscriberCtx, span = collector.tracer.Start(
					subscriberCtx,
					"ro.Subscribe(...)",
					traceAttributes,
				)
			}

			if collector.config.EnableLogging {
				collector.logger.Emit(
					subscriberCtx,
					newRecord("ro.Subscribe(...)", collector.config.LogLevelSubscription),
				)
			}

			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						// Count the number of messages leaving the instrumented ro.PipeX
						if collector.config.EnableMetrics && collector.NotificationsOutTotal != nil {
							collector.NotificationsOutTotal.Add(subscriberCtx, 1, metricAttributes)
						}

						destination.NextWithContext(ctx, value)
					},
					func(ctx context.Context, err error) {
						// Log the error in the last operator of the pipeline
						if collector.config.EnableLogging {
							collector.logger.Emit(
								subscriberCtx,
								newRecord(fmt.Sprintf("ro.Error(%v)", err), collector.config.LogLevelError, log.String("notification.error", err.Error())),
							)
						}

						// Record the error in the current OTEL trace
						if collector.config.EnableTracing {
							s := trace.SpanFromContext(ctx)
							s.SetStatus(codes.Error, "ro.Error(...)")
							s.SetAttributes(attribute.String("notification.error", err.Error()))
							s.RecordError(err)
						}

						destination.ErrorWithContext(ctx, err)
					},
					func(ctx context.Context) {
						// Log completion (different from unsubscription)
						if collector.config.EnableLogging {
							collector.logger.Emit(
								subscriberCtx,
								newRecord("ro.Complete()", collector.config.LogLevelSubscription),
							)
						}

						destination.CompleteWithContext(ctx)
					},
				),
			)

			return func() {
				// Log unsubscription (different from completion)
				if collector.config.EnableLogging {
					collector.logger.Emit(
						subscriberCtx,
						newRecord("ro.Unsubscribe()", collector.config.LogLevelSubscription),
					)
				}

				sub.Unsubscribe()

				if collector.config.EnableTracing {
					span.End()
				}
			}
		})
	}
}
