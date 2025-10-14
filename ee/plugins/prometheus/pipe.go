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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/ro"
	"github.com/samber/ro/ee/internal/introspection"
)

// wrapPipeWithObservability is a wrapper around PipeOpX that adds multiple observability operators to the wrapPipeWithObservability.
func wrapPipeWithObservability[First any, Last any](collector *prometheusCollector, operators func(ro.Observable[First]) ro.Observable[Last]) func(ro.Observable[First]) ro.Observable[Last] {
	return ro.PipeOp3(
		// // Add input notification counter between source and first operator.
		// IncCounterOnNext[First](collector.NotificationsInTotal.With(prometheus.Labels{})),
		// // Track the time it takes for a notification to traverse from the source observable to the destination observer.
		// ObserveNextLag[First](collector.NotificationLagSeconds.With(prometheus.Labels{})),
		// // Track the time it takes for an operator to process a notification.
		// InitOperatorProcessingTimeObserver[First](),

		// observeBeforePipe is the aggregation of the following operators:
		//   - IncCounterOnNext
		//   - InitOperatorProcessingTimeObserver
		//   - ObserveNextLag
		observeBeforePipe[First](
			collector.NotificationsInTotal.With(prometheus.Labels{}),
			collector.NotificationLagSeconds.With(prometheus.Labels{}),
		),

		// PipeX(...)
		operators,

		// // Add output notification counter between last operator and final subscriber.
		// IncCounterOnNext[Last](collector.NotificationsOutTotal.With(prometheus.Labels{})),
		// // Add subscriptions counter at the end: it will be the first called operator on subscription.
		// IncCounterOnSubscription[Last](collector.SubscriptionsTotal),

		// observeAfterPipe is the aggregation of the following operators:
		//   - IncCounterOnNext
		//   - IncCounterOnSubscription
		observeAfterPipe[Last](
			collector.NotificationsOutTotal.With(prometheus.Labels{}),
			collector.SubscriptionsTotal,
		),
	)
}

// Pipe1 is a typesafe 🎉 implementation of Pipe, that takes a source and 1 operator.
func Pipe1[A any, B any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
) (ro.Observable[B], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[B](err), nil
	}

	arg0 := pipeDescription.Arguments[0]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp1(
			operator1,
		),
		// with license
		ro.PipeOp2(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
		),
	), collector
}

// Pipe2 is a typesafe 🎉 implementation of Pipe, that takes a source and 2 operators.
func Pipe2[A any, B any, C any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
) (ro.Observable[C], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[C](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp2(
			operator1,
			operator2,
		),
		// with license
		ro.PipeOp4(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
		),
	), collector
}

// Pipe3 is a typesafe 🎉 implementation of Pipe, that takes a source and 3 operators.
func Pipe3[A any, B any, C any, D any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
) (ro.Observable[D], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[D](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp3(
			operator1,
			operator2,
			operator3,
		),
		// with license
		ro.PipeOp6(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
		),
	), collector
}

// Pipe4 is a typesafe 🎉 implementation of Pipe, that takes a source and 4 operators.
func Pipe4[A any, B any, C any, D any, E any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
) (ro.Observable[E], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[E](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp4(
			operator1,
			operator2,
			operator3,
			operator4,
		),
		// with license
		ro.PipeOp8(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
		),
	), collector
}

// Pipe5 is a typesafe 🎉 implementation of Pipe, that takes a source and 5 operators.
func Pipe5[A any, B any, C any, D any, E any, F any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
) (ro.Observable[F], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[F](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp5(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
		),
		// with license
		ro.PipeOp10(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
		),
	), collector
}

// Pipe6 is a typesafe 🎉 implementation of Pipe, that takes a source and 6 operators.
func Pipe6[A any, B any, C any, D any, E any, F any, G any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
) (ro.Observable[G], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[G](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp6(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
		),
		// with license
		ro.PipeOp12(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
		),
	), collector
}

// Pipe7 is a typesafe 🎉 implementation of Pipe, that takes a source and 7 operators.
func Pipe7[A any, B any, C any, D any, E any, F any, G any, H any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
) (ro.Observable[H], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[H](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp7(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
		),
		// with license
		ro.PipeOp14(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
			operator7,
			observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
		),
	), collector
}

// Pipe8 is a typesafe 🎉 implementation of Pipe, that takes a source and 8 operators.
func Pipe8[A any, B any, C any, D any, E any, F any, G any, H any, I any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
) (ro.Observable[I], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[I](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp8(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
		),
		// with license
		ro.PipeOp16(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
			operator7,
			observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
			operator8,
			observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
		),
	), collector
}

func Pipe9[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
) (ro.Observable[J], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[J](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp9(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
		),
		// with license
		ro.PipeOp18(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
			operator7,
			observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
			operator8,
			observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
			operator9,
			observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
		),
	), collector
}

func Pipe10[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
) (ro.Observable[K], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[K](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp10(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
		),
		// with license
		ro.PipeOp20(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
			operator7,
			observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
			operator8,
			observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
			operator9,
			observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
			operator10,
			observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
		),
	), collector
}

// Pipe11 is a typesafe 🎉 implementation of Pipe, that takes a source and 11 operators.
func Pipe11[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
) (ro.Observable[L], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[L](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp11(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
		),
		// with license
		ro.PipeOp22(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
			operator7,
			observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
			operator8,
			observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
			operator9,
			observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
			operator10,
			observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
			operator11,
			observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
		),
	), collector
}

// Pipe12 is a typesafe 🎉 implementation of Pipe, that takes a source and 12 operators.
func Pipe12[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
) (ro.Observable[M], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[M](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp12(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
		),
		// with license
		ro.PipeOp24(
			operator1,
			observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
			operator2,
			observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
			operator3,
			observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
			operator4,
			observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
			operator5,
			observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
			operator6,
			observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
			operator7,
			observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
			operator8,
			observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
			operator9,
			observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
			operator10,
			observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
			operator11,
			observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
			operator12,
			observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
		),
	), collector
}

// Pipe13 is a typesafe 🎉 implementation of Pipe, that takes a source and 13 operators.
func Pipe13[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
) (ro.Observable[N], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[N](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp13(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp2(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
			),
		),
	), collector
}

// Pipe14 is a typesafe 🎉 implementation of Pipe, that takes a source and 14 operators.
func Pipe14[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
) (ro.Observable[O], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[O](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp14(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp4(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
			),
		),
	), collector
}

// Pipe15 is a typesafe 🎉 implementation of Pipe, that takes a source and 15 operators.
func Pipe15[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
) (ro.Observable[P], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[P](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp15(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp6(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
			),
		),
	), collector
}

// Pipe16 is a typesafe 🎉 implementation of Pipe, that takes a source and 16 operators.
func Pipe16[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
) (ro.Observable[Q], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[Q](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp16(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp8(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
			),
		),
	), collector
}

// Pipe17 is a typesafe 🎉 implementation of Pipe, that takes a source and 17 operators.
func Pipe17[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
) (ro.Observable[R], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[R](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp17(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp10(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
			),
		),
	), collector
}

// Pipe18 is a typesafe 🎉 implementation of Pipe, that takes a source and 18 operators.
func Pipe18[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
) (ro.Observable[S], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[S](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp18(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp12(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
			),
		),
	), collector
}

// Pipe19 is a typesafe 🎉 implementation of Pipe, that takes a source and 19 operators.
func Pipe19[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any, T any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
	operator19 func(ro.Observable[S]) ro.Observable[T],
) (ro.Observable[T], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[T](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]
	arg18 := pipeDescription.Arguments[18]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp19(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
			operator19,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp14(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
				operator19,
				observeOperatorProcessingTime[T](collector.OperatorProcessingTimeSeconds, arg18.Name, arg18.Pos, 18),
			),
		),
	), collector
}

// Pipe20 is a typesafe 🎉 implementation of Pipe, that takes a source and 20 operators.
func Pipe20[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any, T any, U any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
	operator19 func(ro.Observable[S]) ro.Observable[T],
	operator20 func(ro.Observable[T]) ro.Observable[U],
) (ro.Observable[U], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[U](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]
	arg18 := pipeDescription.Arguments[18]
	arg19 := pipeDescription.Arguments[19]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp20(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
			operator19,
			operator20,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp16(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
				operator19,
				observeOperatorProcessingTime[T](collector.OperatorProcessingTimeSeconds, arg18.Name, arg18.Pos, 18),
				operator20,
				observeOperatorProcessingTime[U](collector.OperatorProcessingTimeSeconds, arg19.Name, arg19.Pos, 19),
			),
		),
	), collector
}

// Pipe21 is a typesafe 🎉 implementation of Pipe, that takes a source and 21 operators.
func Pipe21[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any, T any, U any, V any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
	operator19 func(ro.Observable[S]) ro.Observable[T],
	operator20 func(ro.Observable[T]) ro.Observable[U],
	operator21 func(ro.Observable[U]) ro.Observable[V],
) (ro.Observable[V], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[V](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]
	arg18 := pipeDescription.Arguments[18]
	arg19 := pipeDescription.Arguments[19]
	arg20 := pipeDescription.Arguments[20]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp21(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
			operator19,
			operator20,
			operator21,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp18(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
				operator19,
				observeOperatorProcessingTime[T](collector.OperatorProcessingTimeSeconds, arg18.Name, arg18.Pos, 18),
				operator20,
				observeOperatorProcessingTime[U](collector.OperatorProcessingTimeSeconds, arg19.Name, arg19.Pos, 19),
				operator21,
				observeOperatorProcessingTime[V](collector.OperatorProcessingTimeSeconds, arg20.Name, arg20.Pos, 20),
			),
		),
	), collector
}

// Pipe22 is a typesafe 🎉 implementation of Pipe, that takes a source and 22 operators.
func Pipe22[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any, T any, U any, V any, W any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
	operator19 func(ro.Observable[S]) ro.Observable[T],
	operator20 func(ro.Observable[T]) ro.Observable[U],
	operator21 func(ro.Observable[U]) ro.Observable[V],
	operator22 func(ro.Observable[V]) ro.Observable[W],
) (ro.Observable[W], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[W](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]
	arg18 := pipeDescription.Arguments[18]
	arg19 := pipeDescription.Arguments[19]
	arg20 := pipeDescription.Arguments[20]
	arg21 := pipeDescription.Arguments[21]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp22(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
			operator19,
			operator20,
			operator21,
			operator22,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp20(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
				operator19,
				observeOperatorProcessingTime[T](collector.OperatorProcessingTimeSeconds, arg18.Name, arg18.Pos, 18),
				operator20,
				observeOperatorProcessingTime[U](collector.OperatorProcessingTimeSeconds, arg19.Name, arg19.Pos, 19),
				operator21,
				observeOperatorProcessingTime[V](collector.OperatorProcessingTimeSeconds, arg20.Name, arg20.Pos, 20),
				operator22,
				observeOperatorProcessingTime[W](collector.OperatorProcessingTimeSeconds, arg21.Name, arg21.Pos, 21),
			),
		),
	), collector
}

// Pipe23 is a typesafe 🎉 implementation of Pipe, that takes a source and 23 operators.
func Pipe23[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any, T any, U any, V any, W any, X any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
	operator19 func(ro.Observable[S]) ro.Observable[T],
	operator20 func(ro.Observable[T]) ro.Observable[U],
	operator21 func(ro.Observable[U]) ro.Observable[V],
	operator22 func(ro.Observable[V]) ro.Observable[W],
	operator23 func(ro.Observable[W]) ro.Observable[X],
) (ro.Observable[X], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[X](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]
	arg18 := pipeDescription.Arguments[18]
	arg19 := pipeDescription.Arguments[19]
	arg20 := pipeDescription.Arguments[20]
	arg21 := pipeDescription.Arguments[21]
	arg22 := pipeDescription.Arguments[22]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp23(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
			operator19,
			operator20,
			operator21,
			operator22,
			operator23,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp22(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
				operator19,
				observeOperatorProcessingTime[T](collector.OperatorProcessingTimeSeconds, arg18.Name, arg18.Pos, 18),
				operator20,
				observeOperatorProcessingTime[U](collector.OperatorProcessingTimeSeconds, arg19.Name, arg19.Pos, 19),
				operator21,
				observeOperatorProcessingTime[V](collector.OperatorProcessingTimeSeconds, arg20.Name, arg20.Pos, 20),
				operator22,
				observeOperatorProcessingTime[W](collector.OperatorProcessingTimeSeconds, arg21.Name, arg21.Pos, 21),
				operator23,
				observeOperatorProcessingTime[X](collector.OperatorProcessingTimeSeconds, arg22.Name, arg22.Pos, 22),
			),
		),
	), collector
}

// Pipe24 is a typesafe 🎉 implementation of Pipe, that takes a source and 24 operators.
func Pipe24[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any, M any, N any, O any, P any, Q any, R any, S any, T any, U any, V any, W any, X any, Y any](
	collectorConfig CollectorConfig,
	source ro.Observable[A],
	operator1 func(ro.Observable[A]) ro.Observable[B],
	operator2 func(ro.Observable[B]) ro.Observable[C],
	operator3 func(ro.Observable[C]) ro.Observable[D],
	operator4 func(ro.Observable[D]) ro.Observable[E],
	operator5 func(ro.Observable[E]) ro.Observable[F],
	operator6 func(ro.Observable[F]) ro.Observable[G],
	operator7 func(ro.Observable[G]) ro.Observable[H],
	operator8 func(ro.Observable[H]) ro.Observable[I],
	operator9 func(ro.Observable[I]) ro.Observable[J],
	operator10 func(ro.Observable[J]) ro.Observable[K],
	operator11 func(ro.Observable[K]) ro.Observable[L],
	operator12 func(ro.Observable[L]) ro.Observable[M],
	operator13 func(ro.Observable[M]) ro.Observable[N],
	operator14 func(ro.Observable[N]) ro.Observable[O],
	operator15 func(ro.Observable[O]) ro.Observable[P],
	operator16 func(ro.Observable[P]) ro.Observable[Q],
	operator17 func(ro.Observable[Q]) ro.Observable[R],
	operator18 func(ro.Observable[R]) ro.Observable[S],
	operator19 func(ro.Observable[S]) ro.Observable[T],
	operator20 func(ro.Observable[T]) ro.Observable[U],
	operator21 func(ro.Observable[U]) ro.Observable[V],
	operator22 func(ro.Observable[V]) ro.Observable[W],
	operator23 func(ro.Observable[W]) ro.Observable[X],
	operator24 func(ro.Observable[X]) ro.Observable[Y],
) (ro.Observable[Y], prometheus.Collector) {
	pipeDescription, err := introspection.GetFunctionDescription(0, 2)
	if err != nil {
		return ro.Throw[Y](err), nil
	}

	arg0 := pipeDescription.Arguments[0]
	arg1 := pipeDescription.Arguments[1]
	arg2 := pipeDescription.Arguments[2]
	arg3 := pipeDescription.Arguments[3]
	arg4 := pipeDescription.Arguments[4]
	arg5 := pipeDescription.Arguments[5]
	arg6 := pipeDescription.Arguments[6]
	arg7 := pipeDescription.Arguments[7]
	arg8 := pipeDescription.Arguments[8]
	arg9 := pipeDescription.Arguments[9]
	arg10 := pipeDescription.Arguments[10]
	arg11 := pipeDescription.Arguments[11]
	arg12 := pipeDescription.Arguments[12]
	arg13 := pipeDescription.Arguments[13]
	arg14 := pipeDescription.Arguments[14]
	arg15 := pipeDescription.Arguments[15]
	arg16 := pipeDescription.Arguments[16]
	arg17 := pipeDescription.Arguments[17]
	arg18 := pipeDescription.Arguments[18]
	arg19 := pipeDescription.Arguments[19]
	arg20 := pipeDescription.Arguments[20]
	arg21 := pipeDescription.Arguments[21]
	arg22 := pipeDescription.Arguments[22]
	arg23 := pipeDescription.Arguments[23]

	collector := newPrometheusCollector(collectorConfig, pipeDescription)

	return checkLicenseAndPipe(
		collector,
		source,
		// no license
		ro.PipeOp24(
			operator1,
			operator2,
			operator3,
			operator4,
			operator5,
			operator6,
			operator7,
			operator8,
			operator9,
			operator10,
			operator11,
			operator12,
			operator13,
			operator14,
			operator15,
			operator16,
			operator17,
			operator18,
			operator19,
			operator20,
			operator21,
			operator22,
			operator23,
			operator24,
		),
		// with license
		ro.PipeOp2(
			ro.PipeOp24(
				operator1,
				observeOperatorProcessingTime[B](collector.OperatorProcessingTimeSeconds, arg0.Name, arg0.Pos, 0),
				operator2,
				observeOperatorProcessingTime[C](collector.OperatorProcessingTimeSeconds, arg1.Name, arg1.Pos, 1),
				operator3,
				observeOperatorProcessingTime[D](collector.OperatorProcessingTimeSeconds, arg2.Name, arg2.Pos, 2),
				operator4,
				observeOperatorProcessingTime[E](collector.OperatorProcessingTimeSeconds, arg3.Name, arg3.Pos, 3),
				operator5,
				observeOperatorProcessingTime[F](collector.OperatorProcessingTimeSeconds, arg4.Name, arg4.Pos, 4),
				operator6,
				observeOperatorProcessingTime[G](collector.OperatorProcessingTimeSeconds, arg5.Name, arg5.Pos, 5),
				operator7,
				observeOperatorProcessingTime[H](collector.OperatorProcessingTimeSeconds, arg6.Name, arg6.Pos, 6),
				operator8,
				observeOperatorProcessingTime[I](collector.OperatorProcessingTimeSeconds, arg7.Name, arg7.Pos, 7),
				operator9,
				observeOperatorProcessingTime[J](collector.OperatorProcessingTimeSeconds, arg8.Name, arg8.Pos, 8),
				operator10,
				observeOperatorProcessingTime[K](collector.OperatorProcessingTimeSeconds, arg9.Name, arg9.Pos, 9),
				operator11,
				observeOperatorProcessingTime[L](collector.OperatorProcessingTimeSeconds, arg10.Name, arg10.Pos, 10),
				operator12,
				observeOperatorProcessingTime[M](collector.OperatorProcessingTimeSeconds, arg11.Name, arg11.Pos, 11),
			),
			ro.PipeOp24(
				operator13,
				observeOperatorProcessingTime[N](collector.OperatorProcessingTimeSeconds, arg12.Name, arg12.Pos, 12),
				operator14,
				observeOperatorProcessingTime[O](collector.OperatorProcessingTimeSeconds, arg13.Name, arg13.Pos, 13),
				operator15,
				observeOperatorProcessingTime[P](collector.OperatorProcessingTimeSeconds, arg14.Name, arg14.Pos, 14),
				operator16,
				observeOperatorProcessingTime[Q](collector.OperatorProcessingTimeSeconds, arg15.Name, arg15.Pos, 15),
				operator17,
				observeOperatorProcessingTime[R](collector.OperatorProcessingTimeSeconds, arg16.Name, arg16.Pos, 16),
				operator18,
				observeOperatorProcessingTime[S](collector.OperatorProcessingTimeSeconds, arg17.Name, arg17.Pos, 17),
				operator19,
				observeOperatorProcessingTime[T](collector.OperatorProcessingTimeSeconds, arg18.Name, arg18.Pos, 18),
				operator20,
				observeOperatorProcessingTime[U](collector.OperatorProcessingTimeSeconds, arg19.Name, arg19.Pos, 19),
				operator21,
				observeOperatorProcessingTime[V](collector.OperatorProcessingTimeSeconds, arg20.Name, arg20.Pos, 20),
				operator22,
				observeOperatorProcessingTime[W](collector.OperatorProcessingTimeSeconds, arg21.Name, arg21.Pos, 21),
				operator23,
				observeOperatorProcessingTime[X](collector.OperatorProcessingTimeSeconds, arg22.Name, arg22.Pos, 22),
				operator24,
				observeOperatorProcessingTime[Y](collector.OperatorProcessingTimeSeconds, arg23.Name, arg23.Pos, 23),
			),
		),
	), collector
}
