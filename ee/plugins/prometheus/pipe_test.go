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
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	obs, collector := Pipe4(
		CollectorConfig{
			Namespace: "test",
			Subsystem: "sub",
			ConstLabels: prometheus.Labels{
				"a": "b",
			},
		},
		ro.Just(1, 2, 3),
		ro.Map(func(i int) int { return i * 2 }),
		ro.Filter(func(i int) bool { return i%2 == 0 }), // does nothing
		ro.Take[int](10),
		ro.Sum[int](),
	)
	values, err := ro.Collect(obs)
	is.Equal([]int{12}, values)
	is.Nil(err)
	values, err = ro.Collect(obs)
	is.Equal([]int{12}, values)
	is.Nil(err)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	// @TODO: add more tests for metrics
}
