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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type CollectorConfig struct {
	// Enable or disable the collection of traces, metrics and logs
	EnableTracing bool
	EnableMetrics bool
	EnableLogging bool

	// Use the global providers if nil
	TracerProvider trace.TracerProvider
	MetricProvider metric.MeterProvider
	LoggerProvider log.LoggerProvider

	// Attributes to add to all traces, metrics and logs
	TraceAttributes   []attribute.KeyValue
	MetricAttributes  []attribute.KeyValue
	LoggingAttributes []attribute.KeyValue

	// Histogram buckets in seconds
	MetricHistogramObjectivesSeconds []float64

	// On subscription, completion or cancellation
	LogLevelSubscription log.Severity
	// On next
	LogLevelNext log.Severity
	// On error
	LogLevelError log.Severity
}

var DefaultHistogramObjectivesSeconds = []float64{
	0.0010,
	0.0025,
	0.0050,
	0.0075,

	0.010,
	0.025,
	0.050,
	0.075,

	0.10,
	0.25,
	0.50,
	0.75,

	1,
	2.5,
	5,
	10,
	15,
	30,

	60,
	120,
	300,
	600,
	1800,
	3600,
	7200,
	10800,
	21600,
}
