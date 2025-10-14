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


package rosignal

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestNewSignalCatcher_Basic(t *testing.T) {
	observable := NewSignalCatcher(syscall.SIGUSR1)

	var (
		mu              sync.Mutex
		receivedSignals []os.Signal
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu.Lock()
			receivedSignals = append(receivedSignals, sig)
			mu.Unlock()
		},
		nil,
		nil,
	))

	time.Sleep(10 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(50 * time.Millisecond)

	subscription.Unsubscribe()

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, receivedSignals, 1)
	assert.Equal(t, syscall.SIGUSR1, receivedSignals[0])
}

func TestNewSignalCatcher_MultipleSignals(t *testing.T) {
	observable := NewSignalCatcher(syscall.SIGUSR1, syscall.SIGUSR2)

	var (
		mu              sync.Mutex
		receivedSignals []os.Signal
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu.Lock()
			receivedSignals = append(receivedSignals, sig)
			mu.Unlock()
		},
		nil,
		nil,
	))

	time.Sleep(10 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(20 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
	time.Sleep(50 * time.Millisecond)

	subscription.Unsubscribe()

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, receivedSignals, 2)
	assert.Contains(t, receivedSignals, syscall.SIGUSR1)
	assert.Contains(t, receivedSignals, syscall.SIGUSR2)
}

func TestNewSignalCatcher_NoSignals(t *testing.T) {
	observable := NewSignalCatcher()

	var (
		mu              sync.Mutex
		receivedSignals []os.Signal
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu.Lock()
			receivedSignals = append(receivedSignals, sig)
			mu.Unlock()
		},
		nil,
		nil,
	))

	time.Sleep(10 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(50 * time.Millisecond)

	subscription.Unsubscribe()

	mu.Lock()
	defer mu.Unlock()

	// Filter for the signal we care about
	var gotUSR1 bool
	for _, sig := range receivedSignals {
		if sig == syscall.SIGUSR1 {
			gotUSR1 = true
		}
	}
	assert.True(t, gotUSR1, "should receive SIGUSR1")
}

func TestNewSignalCatcher_Unsubscribe(t *testing.T) {
	observable := NewSignalCatcher(syscall.SIGUSR1)

	var (
		mu              sync.Mutex
		receivedSignals []os.Signal
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu.Lock()
			receivedSignals = append(receivedSignals, sig)
			mu.Unlock()
		},
		nil,
		nil,
	))

	time.Sleep(10 * time.Millisecond)
	subscription.Unsubscribe()
	time.Sleep(10 * time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, receivedSignals, 0)
}

func TestNewSignalCatcher_TwoSubscribers(t *testing.T) {
	observable := NewSignalCatcher(syscall.SIGUSR1)

	var (
		mu1, mu2           sync.Mutex
		signals1, signals2 []os.Signal
	)

	sub1 := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu1.Lock()
			signals1 = append(signals1, sig)
			mu1.Unlock()
		},
		nil,
		nil,
	))

	sub2 := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu2.Lock()
			signals2 = append(signals2, sig)
			mu2.Unlock()
		},
		nil,
		nil,
	))

	time.Sleep(10 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(50 * time.Millisecond)

	sub1.Unsubscribe()
	sub2.Unsubscribe()

	mu1.Lock()
	mu2.Lock()
	defer mu1.Unlock()
	defer mu2.Unlock()

	assert.Len(t, signals1, 1)
	assert.Len(t, signals2, 1)
	assert.Equal(t, syscall.SIGUSR1, signals1[0])
	assert.Equal(t, syscall.SIGUSR1, signals2[0])
}

func TestNewSignalCatcher_ErrorCallback(t *testing.T) {
	observable := NewSignalCatcher(syscall.SIGUSR1)

	var (
		mu              sync.Mutex
		receivedSignals []os.Signal
		errors          []error
	)

	subscription := observable.Subscribe(ro.NewObserver(
		func(sig os.Signal) {
			mu.Lock()
			receivedSignals = append(receivedSignals, sig)
			mu.Unlock()
		},
		func(err error) {
			mu.Lock()
			errors = append(errors, err)
			mu.Unlock()
		},
		nil,
	))

	time.Sleep(10 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(50 * time.Millisecond)

	subscription.Unsubscribe()

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, receivedSignals, 1)
	assert.Len(t, errors, 0)
}
