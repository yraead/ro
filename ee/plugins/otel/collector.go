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
	"github.com/samber/ro/ee/internal/introspection"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	logglobal "go.opentelemetry.io/otel/log/global"
	lognoop "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	meternoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

const (
	serviceName = "github.com/samber/ro"

	labelNamePipe             = "pipe.name"
	LabelNamePipePosition     = "pipe.position"
	labelNameOperator         = "operator.name"
	labelNameOperatorPosition = "operator.position"
	labelNameOperatorIndex    = "operator.index"
)

func newOtelCollector(opts CollectorConfig, pipeDescription *introspection.FunDesc) *otelCollector {
	if opts.TraceAttributes == nil {
		opts.TraceAttributes = []attribute.KeyValue{}
	}
	if opts.MetricAttributes == nil {
		opts.MetricAttributes = []attribute.KeyValue{}
	}
	if opts.LoggingAttributes == nil {
		opts.LoggingAttributes = []attribute.KeyValue{}
	}

	opts.TraceAttributes = concatAttributes(
		opts.TraceAttributes,
		[]attribute.KeyValue{
			attribute.String(labelNamePipe, pipeDescription.Name),
			attribute.String(LabelNamePipePosition, pipeDescription.Pos),
		},
	)
	opts.MetricAttributes = concatAttributes(
		opts.MetricAttributes,
		[]attribute.KeyValue{
			attribute.String(labelNamePipe, pipeDescription.Name),
			attribute.String(LabelNamePipePosition, pipeDescription.Pos),
		},
	)
	opts.LoggingAttributes = concatAttributes(
		opts.LoggingAttributes,
		[]attribute.KeyValue{
			attribute.String(labelNamePipe, pipeDescription.Name),
			attribute.String(LabelNamePipePosition, pipeDescription.Pos),
		},
	)

	if opts.MetricHistogramObjectivesSeconds == nil {
		opts.MetricHistogramObjectivesSeconds = DefaultHistogramObjectivesSeconds
	}

	if opts.LogLevelSubscription == 0 {
		opts.LogLevelSubscription = log.SeverityInfo
	}
	if opts.LogLevelNext == 0 {
		opts.LogLevelNext = log.SeverityInfo
	}
	if opts.LogLevelError == 0 {
		opts.LogLevelError = log.SeverityError
	}

	collector := &otelCollector{config: opts}
	collector.initTracer()
	collector.initMeter()
	collector.initLogger()
	return collector
}

type otelCollector struct {
	config CollectorConfig

	tracer trace.Tracer
	meter  metric.Meter
	logger log.Logger

	// @TODO:
	//   - SubscriptionDurationSeconds (summary): time between subscription and complete/error/unsuscribe
	//   - OperatorNotificationInflationRate (gauge): from an operator POV, the difference between the rx and tx notifications.
	SubscriptionsTotal            metric.Int64Counter
	NotificationsInTotal          metric.Int64Counter
	NotificationsOutTotal         metric.Int64Counter
	NotificationLagSeconds        metric.Float64Histogram
	OperatorProcessingTimeSeconds metric.Float64Histogram
}

func (c *otelCollector) initTracer() {
	if c.config.EnableTracing {
		if c.config.TracerProvider == nil {
			c.config.TracerProvider = otel.GetTracerProvider()
		}
		c.tracer = c.config.TracerProvider.Tracer(serviceName)
	} else {
		c.tracer = tracenoop.NewTracerProvider().Tracer(serviceName)
	}
}

func (c *otelCollector) initMeter() {
	if c.config.EnableMetrics {
		if c.config.MetricProvider == nil {
			c.config.MetricProvider = otel.GetMeterProvider()
		}
		c.meter = c.config.MetricProvider.Meter(serviceName)
		c.registerMetrics()
	} else {
		c.meter = meternoop.NewMeterProvider().Meter(serviceName)
	}
}

func (c *otelCollector) initLogger() {
	if c.config.EnableLogging {
		if c.config.LoggerProvider == nil {
			c.config.LoggerProvider = logglobal.GetLoggerProvider()
		}
		c.logger = c.config.LoggerProvider.Logger(serviceName)
	} else {
		c.logger = lognoop.NewLoggerProvider().Logger(
			serviceName,
			log.WithInstrumentationAttributes(c.config.LoggingAttributes...),
		)
	}
}

func (c *otelCollector) registerMetrics() {
	counter, err := c.meter.Int64Counter(
		"samber/ro.subscriptions_total",
		metric.WithDescription("Total number of subscriptions to a chain of operators."),
	)
	if err == nil {
		c.SubscriptionsTotal = counter
	}

	counter, err = c.meter.Int64Counter(
		"samber/ro.notifications_in_total",
		metric.WithDescription("Total number of notifications emitted by a source in a chain of operators."),
	)
	if err == nil {
		c.NotificationsInTotal = counter
	}

	counter, err = c.meter.Int64Counter(
		"samber/ro.notifications_out_total",
		metric.WithDescription("Total number of notifications emitted by a chain of operators."),
	)
	if err == nil {
		c.NotificationsOutTotal = counter
	}

	histogram, err := c.meter.Float64Histogram(
		"samber/ro.notification_lag_seconds",
		metric.WithDescription("Time for notifications to traverse from the source observable to the destination observer."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(c.config.MetricHistogramObjectivesSeconds...),
	)
	if err == nil {
		c.NotificationLagSeconds = histogram
	}

	histogram, err = c.meter.Float64Histogram(
		"samber/ro.operator_processing_time_seconds",
		metric.WithDescription("Total number of subscriptions to a chain of operators."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(c.config.MetricHistogramObjectivesSeconds...),
	)
	if err == nil {
		c.OperatorProcessingTimeSeconds = histogram
	}
}
