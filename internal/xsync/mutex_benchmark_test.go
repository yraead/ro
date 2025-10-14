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

package xsync

import (
	"sync"
	"testing"
)

// BenchmarkMutexWithLock benchmarks the standard mutex under different contention levels.
func BenchmarkMutexWithLock_NoContention(b *testing.B) {
	mutex := NewMutexWithLock()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mutex.Lock()
		// Simulate minimal work
		_ = i + 1

		mutex.Unlock()
	}
}

func BenchmarkMutexWithLock_LowContention(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 2

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithLock_MediumContention(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithLock_HighContention(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 128

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithLock_ExtremeContention(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 1024

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

// BenchmarkMutexWithSpinlock benchmarks the spinlock mutex under different contention levels.
func BenchmarkMutexWithSpinlock_NoContention(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mutex.Lock()
		// Simulate minimal work
		_ = i + 1

		mutex.Unlock()
	}
}

func BenchmarkMutexWithSpinlock_LowContention(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 2

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithSpinlock_MediumContention(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithSpinlock_HighContention(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 128

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithSpinlock_ExtremeContention(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 1024

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

// BenchmarkMutexWithoutLock benchmarks the fake mutex under different contention levels.
func BenchmarkMutexWithoutLock_NoContention(b *testing.B) {
	mutex := NewMutexWithoutLock()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mutex.Lock()
		// Simulate minimal work
		_ = i + 1

		mutex.Unlock()
	}
}

func BenchmarkMutexWithoutLock_LowContention(b *testing.B) {
	mutex := NewMutexWithoutLock()

	var wg sync.WaitGroup

	workers := 2

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithoutLock_MediumContention(b *testing.B) {
	mutex := NewMutexWithoutLock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithoutLock_HighContention(b *testing.B) {
	mutex := NewMutexWithoutLock()

	var wg sync.WaitGroup

	workers := 128

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithoutLock_ExtremeContention(b *testing.B) {
	mutex := NewMutexWithoutLock()

	var wg sync.WaitGroup

	workers := 1024

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Simulate some work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

// BenchmarkTryLock benchmarks TryLock performance.
func BenchmarkMutexWithLock_TryLock(b *testing.B) {
	mutex := NewMutexWithLock()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mutex.TryLock()
		mutex.Unlock()
	}
}

func BenchmarkMutexWithSpinlock_TryLock(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mutex.TryLock()
		mutex.Unlock()
	}
}

func BenchmarkMutexWithoutLock_TryLock(b *testing.B) {
	mutex := NewMutexWithoutLock()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mutex.TryLock()
		mutex.Unlock()
	}
}

// BenchmarkMixedOperations benchmarks mixed Lock/TryLock operations.
func BenchmarkMutexWithLock_MixedOperations(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				if mutex.TryLock() {
					_ = k + 1

					mutex.Unlock()
				} else {
					mutex.Lock()

					_ = k + 1

					mutex.Unlock()
				}
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithSpinlock_MixedOperations(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				if mutex.TryLock() {
					_ = k + 1

					mutex.Unlock()
				} else {
					mutex.Lock()

					_ = k + 1

					mutex.Unlock()
				}
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithoutLock_MixedOperations(b *testing.B) {
	mutex := NewMutexWithoutLock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				if mutex.TryLock() {
					_ = k + 1

					mutex.Unlock()
				} else {
					mutex.Lock()

					_ = k + 1

					mutex.Unlock()
				}
			}(j)
		}

		wg.Wait()
	}
}

// BenchmarkWorkloadIntensity benchmarks with different work intensities.
func BenchmarkMutexWithLock_LightWork(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Light work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithLock_HeavyWork(b *testing.B) {
	mutex := NewMutexWithLock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Heavy work - simulate more computation
				sum := 0
				for k := 0; k < 1000; k++ {
					sum += k + 1
				}

				_ = sum

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithSpinlock_LightWork(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func(k int) {
				defer wg.Done()

				mutex.Lock()
				// Light work
				_ = k + 1

				mutex.Unlock()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkMutexWithSpinlock_HeavyWork(b *testing.B) {
	mutex := NewMutexWithSpinlock()

	var wg sync.WaitGroup

	workers := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(workers)

		for j := 0; j < workers; j++ {
			go func() {
				defer wg.Done()

				mutex.Lock()
				// Heavy work - simulate more computation
				sum := 0
				for k := 0; k < 1000; k++ {
					sum += k + 1
				}

				_ = sum

				mutex.Unlock()
			}()
		}

		wg.Wait()
	}
}

// BenchmarkContentionComparison provides a comprehensive comparison of all mutex types.
func BenchmarkContentionComparison(b *testing.B) {
	contentionLevels := []struct {
		name    string
		workers int
	}{
		{"NoContention", 1},
		{"LowContention", 2},
		{"MediumContention", 8},
		{"HighContention", 32},
	}

	mutexTypes := []struct {
		name  string
		mutex Mutex
	}{
		{"Standard", NewMutexWithLock()},
		{"Spinlock", NewMutexWithSpinlock()},
		{"Fake", NewMutexWithoutLock()},
	}

	for _, level := range contentionLevels {
		for _, mt := range mutexTypes {
			b.Run(level.name+"/"+mt.name, func(b *testing.B) {
				mutex := mt.mutex

				var wg sync.WaitGroup

				workers := level.workers

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(workers)

					for j := 0; j < workers; j++ {
						go func(k int) {
							defer wg.Done()

							mutex.Lock()
							// Simulate some work
							_ = k + 1

							mutex.Unlock()
						}(j)
					}

					wg.Wait()
				}
			})
		}
	}
}

// BenchmarkTryLockComparison compares TryLock performance across all mutex types.
func BenchmarkTryLockComparison(b *testing.B) {
	mutexTypes := []struct {
		name  string
		mutex Mutex
	}{
		{"Standard", NewMutexWithLock()},
		{"Spinlock", NewMutexWithSpinlock()},
		{"Fake", NewMutexWithoutLock()},
	}

	for _, mt := range mutexTypes {
		b.Run(mt.name, func(b *testing.B) {
			mutex := mt.mutex

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mutex.TryLock()
				mutex.Unlock()
			}
		})
	}
}
