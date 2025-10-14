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

import "github.com/prometheus/client_golang/prometheus"

type CollectorConfig struct {
	// Namespace is the namespace of the metrics.
	Namespace string
	// Subsystem is the subsystem of the metrics.
	Subsystem string
	// ConstLabels are labels that are applied to all metrics.
	ConstLabels prometheus.Labels

	// Custom options
	SummaryObjectives map[float64]float64
}

var DefaultSummaryObjectives = map[float64]float64{
	0.5:  0.05,
	0.9:  0.01,
	0.95: 0.005,
	0.99: 0.001,
}
