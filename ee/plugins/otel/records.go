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
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func concatAttributes[T any](kv ...[]T) []T {
	length := 0
	for _, v := range kv {
		length += len(v)
	}

	output := make([]T, length)
	for _, v := range kv {
		output = append(output, v...)
	}
	return output
}

func traceWithAttributes(kv ...[]attribute.KeyValue) trace.SpanStartEventOption {
	return trace.WithAttributes(concatAttributes(kv...)...)
}

func metricWithAttributes(kv ...[]attribute.KeyValue) metric.MeasurementOption {
	return metric.WithAttributes(concatAttributes(kv...)...)
}

func newRecord(msg string, severity log.Severity, kv ...log.KeyValue) log.Record {
	var r log.Record
	r.SetEventName(msg)
	r.SetObservedTimestamp(time.Now())
	r.SetSeverity(severity)
	r.AddAttributes(kv...)
	return r
}
