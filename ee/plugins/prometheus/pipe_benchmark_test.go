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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/ro"
)

var operators = []func(ro.Observable[int]) ro.Observable[int]{
	ro.Map(func(i int) int { return i * 2 }),
	ro.Filter(func(i int) bool { return i%2 == 0 }),
	ro.Take[int](100),
	ro.Sum[int](),
}

// BenchmarkStdPipe4 benchmarks the non-instrumented Pipe4 implementation
// from the base ro package
func BenchmarkStdPipe4(b *testing.B) {
	// Create a large dataset to make the overhead more measurable
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	// Create the non-instrumented pipe using base ro.Pipe4
	obs := ro.Pipe4(
		ro.Just(data...),
		operators[0],
		operators[1],
		operators[2],
		operators[3],
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := ro.Collect(obs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPrometheusPipe4 benchmarks the instrumented Pipe4 implementation
// that includes Prometheus metrics collection
func BenchmarkPrometheusPipe4_goodLicense(b *testing.B) {
	bypassLicenseCheck = true

	// Create a large dataset to make the overhead more measurable
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	// Create the instrumented pipe
	obs, _ := Pipe4(
		CollectorConfig{
			Namespace: "benchmark",
			Subsystem: "test",
			ConstLabels: prometheus.Labels{
				"benchmark": "instrumented",
			},
		},
		ro.Just(data...),
		operators[0],
		operators[1],
		operators[2],
		operators[3],
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := ro.Collect(obs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPrometheusPipe4 benchmarks the instrumented Pipe4 implementation
// that includes Prometheus metrics collection
func BenchmarkPrometheusPipe4_badLicense(b *testing.B) {
	bypassLicenseCheck = false

	// Create a large dataset to make the overhead more measurable
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	// Create the instrumented pipe
	obs, _ := Pipe4(
		CollectorConfig{
			Namespace: "benchmark",
			Subsystem: "test",
			ConstLabels: prometheus.Labels{
				"benchmark": "instrumented",
			},
		},
		ro.Just(data...),
		operators[0],
		operators[1],
		operators[2],
		operators[3],
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := ro.Collect(obs)
		if err != nil {
			b.Fatal(err)
		}
	}
}
