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

package roproc

import (
	"context"
	"time"

	"github.com/samber/ro"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

// NewVirtualMemoryWatcher creates an observable that emits virtual memory statistics at regular intervals.
func NewVirtualMemoryWatcher(interval time.Duration) ro.Observable[*mem.VirtualMemoryStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*mem.VirtualMemoryStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := mem.VirtualMemory()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewSwapMemoryWatcher creates an observable that emits swap memory statistics at regular intervals.
func NewSwapMemoryWatcher(interval time.Duration) ro.Observable[*mem.SwapMemoryStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*mem.SwapMemoryStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := mem.SwapMemory()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewSwapDeviceWatcher creates an observable that emits swap device information at regular intervals.
func NewSwapDeviceWatcher(interval time.Duration) ro.Observable[*mem.SwapDevice] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*mem.SwapDevice]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := mem.SwapDevices()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewCPUInfoWatcher creates an observable that emits CPU information statistics at regular intervals.
func NewCPUInfoWatcher(interval time.Duration) ro.Observable[cpu.InfoStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[cpu.InfoStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := cpu.Info()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewDiskUsageWatcher creates an observable that emits disk usage statistics at regular intervals.
func NewDiskUsageWatcher(interval time.Duration, mountpointOrDevicePath string) ro.Observable[*disk.UsageStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*disk.UsageStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := disk.Usage(mountpointOrDevicePath)
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewDiskIOCountersWatcher creates an observable that emits disk I/O counters at regular intervals.
func NewDiskIOCountersWatcher(interval time.Duration, names ...string) ro.Observable[map[string]disk.IOCountersStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[map[string]disk.IOCountersStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := disk.IOCounters(names...)
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewDiskPartitionWatcher creates an observable that emits disk partition information at regular intervals.
func NewDiskPartitionWatcher(interval time.Duration) ro.Observable[disk.PartitionStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[disk.PartitionStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := disk.Partitions(true)
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewHostInfoWatcher creates an observable that emits host information at regular intervals.
func NewHostInfoWatcher(interval time.Duration) ro.Observable[*host.InfoStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*host.InfoStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := host.Info()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewHostUserWatcher creates an observable that emits host user information at regular intervals.
func NewHostUserWatcher(interval time.Duration) ro.Observable[host.UserStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[host.UserStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := host.Users()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewLoadAverageWatcher creates an observable that emits load average statistics at regular intervals.
func NewLoadAverageWatcher(interval time.Duration) ro.Observable[*load.AvgStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*load.AvgStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := load.Avg()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewLoadMiscWatcher creates an observable that emits miscellaneous load statistics at regular intervals.
func NewLoadMiscWatcher(interval time.Duration) ro.Observable[*load.MiscStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[*load.MiscStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := load.Misc()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						destination.NextWithContext(ctx, stats)
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewNetConnectionsWatcher creates an observable that emits network connection statistics at regular intervals.
func NewNetConnectionsWatcher(interval time.Duration) ro.Observable[net.ConnectionStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[net.ConnectionStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := net.Connections("all")
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewNetConntrackWatcher creates an observable that emits network conntrack statistics at regular intervals.
func NewNetConntrackWatcher(interval time.Duration, perCPU bool) ro.Observable[net.ConntrackStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[net.ConntrackStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := net.ConntrackStats(perCPU)
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewNetFilterCountersWatcher creates an observable that emits network filter counters at regular intervals.
func NewNetFilterCountersWatcher(interval time.Duration) ro.Observable[net.FilterStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[net.FilterStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := net.FilterCounters()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewNetIOCountersWatcher creates an observable that emits network I/O counters at regular intervals.
func NewNetIOCountersWatcher(interval time.Duration, perNIC bool) ro.Observable[net.IOCountersStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[net.IOCountersStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := net.IOCounters(perNIC)
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}

// NewSensorsTemperatureWatcher creates an observable that emits sensor temperature statistics at regular intervals.
func NewSensorsTemperatureWatcher(interval time.Duration, perNIC bool) ro.Observable[sensors.TemperatureStat] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[sensors.TemperatureStat]) ro.Teardown {
		sub := ro.Interval(interval).
			SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value int64) {
						stats, err := sensors.SensorsTemperatures()
						if err != nil {
							destination.ErrorWithContext(ctx, err)
							return
						}

						for _, stat := range stats {
							destination.NextWithContext(ctx, stat)
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

		return sub.Unsubscribe
	})
}
