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

package xatomic

import (
	"testing"
)

func TestPointer_NewPointer(t *testing.T) {
	t.Parallel()
	x := 42
	ptr := NewPointer(&x)

	if ptr.Load() != &x {
		t.Errorf("NewPointer failed, expected %p, got %p", &x, ptr.Load())
	}
}

func TestPointer_StoreAndLoad(t *testing.T) {
	t.Parallel()
	var p Pointer[int]

	// Test storing nil
	p.Store(nil)
	if p.Load() != nil {
		t.Error("Expected nil after storing nil")
	}

	// Test storing a value
	x := 42
	p.Store(&x)

	if p.Load() != &x {
		t.Errorf("Load failed, expected %p, got %p", &x, p.Load())
	}

	if *p.Load() != 42 {
		t.Errorf("Load returned wrong value, expected 42, got %d", *p.Load())
	}
}

func TestPointer_Swap(t *testing.T) {
	t.Parallel()
	var p Pointer[string]

	// Initial state
	old := p.Swap(nil)
	if old != nil {
		t.Error("Expected nil from initial Swap")
	}

	// Swap with a value
	s1 := "hello"
	old = p.Swap(&s1)
	if old != nil {
		t.Error("Expected nil from Swap when swapping to first value")
	}

	// Swap to another value
	s2 := "world"
	old = p.Swap(&s2)
	if old != &s1 {
		t.Error("Swap didn't return the old value")
	}

	if *p.Load() != "world" {
		t.Error("Swap didn't store the new value")
	}
}

func TestPointer_CompareAndSwap(t *testing.T) {
	t.Parallel()
	var p Pointer[int]

	x := 1
	y := 2
	z := 3

	// Store initial value
	p.Store(&x)

	// Successful compare-and-swap
	swapped := p.CompareAndSwap(&x, &y)
	if !swapped {
		t.Error("CompareAndSwap should have succeeded")
	}

	if p.Load() != &y {
		t.Error("CompareAndSwap didn't store the new value")
	}

	// Failed compare-and-swap
	swapped = p.CompareAndSwap(&x, &z)
	if swapped {
		t.Error("CompareAndSwap should have failed")
	}

	if p.Load() != &y {
		t.Error("CompareAndSwap shouldn't have changed the value on failure")
	}

	// Successful compare-and-swap with correct old value
	swapped = p.CompareAndSwap(&y, &z)
	if !swapped {
		t.Error("CompareAndSwap should have succeeded with correct old value")
	}

	if p.Load() != &z {
		t.Error("CompareAndSwap didn't store the new value on successful swap")
	}
}

func TestPointer_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	var p Pointer[int]

	done := make(chan bool, 2)

	// Goroutine 1: store values
	go func() {
		for i := 0; i < 1000; i++ {
			val := i
			p.Store(&val)
		}
		done <- true
	}()

	// Goroutine 2: load values
	go func() {
		for i := 0; i < 1000; i++ {
			_ = p.Load()
		}
		done <- true
	}()

	<-done
	<-done
}

func TestPointer_UnsafePointer(t *testing.T) {
	t.Parallel()
	var p Pointer[int]

	x := 42
	p.Store(&x)

	if p.Load() != &x {
		t.Error("Store didn't work")
	}

	if *p.Load() != 42 {
		t.Error("Store returned wrong value")
	}
}
