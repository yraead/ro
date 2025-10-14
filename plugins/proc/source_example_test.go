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

func ExampleNewVirtualMemoryWatcher() {
	// Monitor virtual memory usage
	observable := NewVirtualMemoryWatcher(1 * time.Second)

	subscription := observable.Subscribe(ro.NoopObserver[*mem.VirtualMemoryStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(3 * time.Second)
}

func ExampleNewCPUInfoWatcher() {
	// Monitor CPU information
	observable := NewCPUInfoWatcher(2 * time.Second)

	subscription := observable.Subscribe(ro.NoopObserver[cpu.InfoStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(4 * time.Second)
}

func ExampleNewDiskUsageWatcher() {
	// Monitor disk usage for root filesystem
	observable := NewDiskUsageWatcher(5*time.Second, "/")

	subscription := observable.Subscribe(ro.NoopObserver[*disk.UsageStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(10 * time.Second)
}

func ExampleNewNetIOCountersWatcher() {
	// Monitor network I/O counters
	observable := NewNetIOCountersWatcher(3*time.Second, true)

	subscription := observable.Subscribe(ro.NoopObserver[net.IOCountersStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(6 * time.Second)
}

func ExampleNewLoadAverageWatcher() {
	// Monitor system load average
	observable := NewLoadAverageWatcher(2 * time.Second)

	subscription := observable.Subscribe(ro.NoopObserver[*load.AvgStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(4 * time.Second)
}

func ExampleNewHostInfoWatcher() {
	// Monitor host information
	observable := NewHostInfoWatcher(10 * time.Second)

	subscription := observable.Subscribe(ro.NoopObserver[*host.InfoStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(12 * time.Second)
}

func ExampleNewSensorsTemperatureWatcher() {
	// Monitor temperature sensors
	observable := NewSensorsTemperatureWatcher(5*time.Second, false)

	subscription := observable.Subscribe(ro.NoopObserver[sensors.TemperatureStat]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(10 * time.Second)
}
