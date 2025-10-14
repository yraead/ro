// Copyright 2025 samber.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rocron

import (
	"context"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestNewScheduler(t *testing.T) {
	obs := NewScheduler(
		gocron.DurationJob(
			100 * time.Millisecond,
		),
	)
	assert.NotNil(t, obs)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(250 * time.Millisecond)
		cancel()
	}()

	items, _, err := ro.CollectWithContext(ctx, obs)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Len(t, items, 2)
	assert.Equal(t, items[0].Counter, 0)
	assert.Equal(t, items[1].Counter, 1)

	// 100ms between the first and second item
	assert.WithinDuration(t, items[0].Time.Add(100*time.Millisecond), items[1].Time, 40*time.Millisecond)
}

func TestNewScheduler_Shutdown(t *testing.T) {
	obs := NewScheduler(
		gocron.DurationJob(
			100 * time.Millisecond,
		),
	)
	assert.NotNil(t, obs)

	var items []ScheduleJob

	sub := obs.Subscribe(
		ro.NewObserver(
			func(item ScheduleJob) {
				items = append(items, item)
			},
			func(err error) {
				assert.Fail(t, "should not error")
			},
			func() {
				assert.Fail(t, "should not complete")
			},
		),
	)

	time.Sleep(250 * time.Millisecond)
	sub.Unsubscribe()

	assert.Len(t, items, 2)
	assert.Equal(t, items[0].Counter, 0)
	assert.Equal(t, items[1].Counter, 1)

	// 100ms between the first and second item
	assert.WithinDuration(t, items[0].Time.Add(100*time.Millisecond), items[1].Time, 40*time.Millisecond)
}
