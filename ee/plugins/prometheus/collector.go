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
	"github.com/samber/lo"
	"github.com/samber/ro/ee/internal/introspection"
)

const (
	labelNamePipe         = "pipe"
	LabelNamePipePosition = "pipe_position"
	// LabelNameSubscriptionID   = "subscription_id"
	labelNameOperator         = "operator"
	labelNameOperatorPosition = "operator_position"
	labelNameOperatorIndex    = "operator_index"
)

var _ prometheus.Collector = (*prometheusCollector)(nil)

func newPrometheusCollector(opts CollectorConfig, pipeDescription *introspection.FunDesc) *prometheusCollector {
	if opts.SummaryObjectives == nil {
		opts.SummaryObjectives = DefaultSummaryObjectives
	}
	if opts.ConstLabels == nil {
		opts.ConstLabels = prometheus.Labels{}
	}

	constLabels := lo.Assign(
		map[string]string{},
		opts.ConstLabels,
		map[string]string{
			labelNamePipe:         pipeDescription.Name,
			LabelNamePipePosition: pipeDescription.Pos,
		},
	)

	return &prometheusCollector{
		SubscriptionsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace:   opts.Namespace,
				Subsystem:   opts.Subsystem,
				Name:        "ro_subscriptions_total",
				Help:        "Total number of subscriptions to a chain of operators.",
				ConstLabels: constLabels,
			},
		),

		NotificationsInTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   opts.Namespace,
				Subsystem:   opts.Subsystem,
				Name:        "ro_notification_in_total",
				Help:        "Total number of notifications emitted by a source in a chain of operators.",
				ConstLabels: constLabels,
			},
			[]string{},
		),
		NotificationsOutTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   opts.Namespace,
				Subsystem:   opts.Subsystem,
				Name:        "ro_notification_out_total",
				Help:        "Total number of notifications emitted by a chain of operators.",
				ConstLabels: constLabels,
			},
			[]string{},
		),
		NotificationLagSeconds: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:   opts.Namespace,
				Subsystem:   opts.Subsystem,
				Name:        "ro_notification_lag_seconds",
				Help:        "Time for notifications to traverse from the source observable to the destination observer.",
				ConstLabels: constLabels,
				Objectives:  opts.SummaryObjectives,
			},
			[]string{},
		),

		OperatorProcessingTimeSeconds: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:   opts.Namespace,
				Subsystem:   opts.Subsystem,
				Name:        "ro_operator_processing_time_seconds_total",
				Help:        "Total number of subscriptions to a chain of operators.",
				ConstLabels: constLabels,
				Objectives:  opts.SummaryObjectives,
			},
			[]string{labelNameOperator, labelNameOperatorPosition, labelNameOperatorIndex},
		),
	}
}

type prometheusCollector struct {
	// @TODO:
	//   - SubscriptionDurationSeconds (summary): time between subscription and complete/error/unsuscribe
	//   - OperatorNotificationInflationRate (gauge): from an operator POV, the difference between the rx and tx notifications.
	SubscriptionsTotal            prometheus.Counter
	NotificationsInTotal          *prometheus.CounterVec
	NotificationsOutTotal         *prometheus.CounterVec
	NotificationLagSeconds        *prometheus.SummaryVec
	OperatorProcessingTimeSeconds *prometheus.SummaryVec
}

func (c *prometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	if !isPrometheusEnabled() {
		return
	}

	c.SubscriptionsTotal.Describe(ch)
	c.NotificationsInTotal.Describe(ch)
	c.NotificationsOutTotal.Describe(ch)
	c.NotificationLagSeconds.Describe(ch)
	c.OperatorProcessingTimeSeconds.Describe(ch)
}

func (c *prometheusCollector) Collect(ch chan<- prometheus.Metric) {
	if !isPrometheusEnabled() {
		return
	}

	c.SubscriptionsTotal.Collect(ch)
	c.NotificationsInTotal.Collect(ch)
	c.NotificationsOutTotal.Collect(ch)
	c.NotificationLagSeconds.Collect(ch)
	c.OperatorProcessingTimeSeconds.Collect(ch)
}
