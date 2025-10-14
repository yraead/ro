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

// BenchmarkMutexVsRWMutex_ReadOnly benchmarks read-only workloads.
func BenchmarkMutexVsRWMutex_ReadOnly(b *testing.B) {
	scenarios := []struct {
		name       string
		readers    int
		iterations int
	}{
		{"SingleReader", 1, 1000},
		{"FewReaders", 4, 1000},
		{"ManyReaders", 16, 1000},
		{"ExtremeReaders", 64, 1000},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			// Test with regular mutex
			b.Run("Mutex", func(b *testing.B) {
				mutex := NewMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.readers)

					for j := 0; j < scenario.readers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								mutex.Lock()
								// Simulate read operation
								_ = k * k

								mutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})

			// Test with RWMutex
			b.Run("RWMutex", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.readers)

					for j := 0; j < scenario.readers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								rwmutex.RLock()
								// Simulate read operation
								_ = k * k

								rwmutex.RUnlock()
							}
						}()
					}

					wg.Wait()
				}
			})
		})
	}
}

// BenchmarkMutexVsRWMutex_WriteOnly benchmarks write-only workloads.
func BenchmarkMutexVsRWMutex_WriteOnly(b *testing.B) {
	scenarios := []struct {
		name       string
		writers    int
		iterations int
	}{
		{"SingleWriter", 1, 1000},
		{"FewWriters", 4, 1000},
		{"ManyWriters", 16, 1000},
		{"ExtremeWriters", 64, 1000},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			// Test with regular mutex
			b.Run("Mutex", func(b *testing.B) {
				mutex := NewMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.writers)

					for j := 0; j < scenario.writers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								mutex.Lock()
								// Simulate write operation
								_ = k * k

								mutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})

			// Test with RWMutex
			b.Run("RWMutex", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.writers)

					for j := 0; j < scenario.writers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								rwmutex.Lock()
								// Simulate write operation
								_ = k * k

								rwmutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})
		})
	}
}

// BenchmarkMutexVsRWMutex_ReadWriteMixed benchmarks mixed read-write workloads.
func BenchmarkMutexVsRWMutex_ReadWriteMixed(b *testing.B) {
	scenarios := []struct {
		name       string
		readers    int
		writers    int
		iterations int
	}{
		{"ReadHeavy_1W_4R", 4, 1, 1000},
		{"ReadHeavy_1W_16R", 16, 1, 1000},
		{"ReadHeavy_1W_64R", 64, 1, 1000},
		{"Balanced_4W_4R", 4, 4, 1000},
		{"Balanced_16W_16R", 16, 16, 1000},
		{"WriteHeavy_4W_1R", 1, 4, 1000},
		{"WriteHeavy_16W_1R", 1, 16, 1000},
		{"WriteHeavy_64W_1R", 1, 64, 1000},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			// Test with regular mutex
			b.Run("Mutex", func(b *testing.B) {
				mutex := NewMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.readers + scenario.writers)

					// Start readers
					for j := 0; j < scenario.readers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								mutex.Lock()
								// Simulate read operation
								_ = k * k

								mutex.Unlock()
							}
						}()
					}

					// Start writers
					for j := 0; j < scenario.writers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								mutex.Lock()
								// Simulate write operation
								_ = k * k

								mutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})

			// Test with RWMutex
			b.Run("RWMutex", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.readers + scenario.writers)

					// Start readers
					for j := 0; j < scenario.readers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								rwmutex.RLock()
								// Simulate read operation
								_ = k * k

								rwmutex.RUnlock()
							}
						}()
					}

					// Start writers
					for j := 0; j < scenario.writers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.iterations; k++ {
								rwmutex.Lock()
								// Simulate write operation
								_ = k * k

								rwmutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})
		})
	}
}

// BenchmarkMutexVsRWMutex_TryLockVsTryRLock benchmarks TryLock vs TryRLock performance.
func BenchmarkMutexVsRWMutex_TryLockVsTryRLock(b *testing.B) {
	// Test TryLock performance
	b.Run("Mutex_TryLock", func(b *testing.B) {
		mutex := NewMutexWithLock()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			mutex.TryLock()
			mutex.Unlock()
		}
	})

	// Test TryRLock performance
	b.Run("RWMutex_TryRLock", func(b *testing.B) {
		rwmutex := NewRWMutexWithLock()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			rwmutex.TryRLock()
			rwmutex.RUnlock()
		}
	})

	// Test TryLock on RWMutex
	b.Run("RWMutex_TryLock", func(b *testing.B) {
		rwmutex := NewRWMutexWithLock()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			rwmutex.TryLock()
			rwmutex.Unlock()
		}
	})
}

// BenchmarkMutexVsRWMutex_ContentionLevels benchmarks different contention levels.
func BenchmarkMutexVsRWMutex_ContentionLevels(b *testing.B) {
	contentionLevels := []struct {
		name    string
		workers int
	}{
		{"LowContention", 4},
		{"MediumContention", 16},
		{"HighContention", 64},
		{"ExtremeContention", 256},
	}

	for _, level := range contentionLevels {
		b.Run(level.name, func(b *testing.B) {
			// Test with regular mutex
			b.Run("Mutex", func(b *testing.B) {
				mutex := NewMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(level.workers)

					for j := 0; j < level.workers; j++ {
						go func(k int) {
							defer wg.Done()

							mutex.Lock()
							// Simulate work
							_ = k * k

							mutex.Unlock()
						}(j)
					}

					wg.Wait()
				}
			})

			// Test with RWMutex (all readers)
			b.Run("RWMutex_AllReaders", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(level.workers)

					for j := 0; j < level.workers; j++ {
						go func(k int) {
							defer wg.Done()

							rwmutex.RLock()
							// Simulate work
							_ = k * k

							rwmutex.RUnlock()
						}(j)
					}

					wg.Wait()
				}
			})

			// Test with RWMutex (all writers)
			b.Run("RWMutex_AllWriters", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(level.workers)

					for j := 0; j < level.workers; j++ {
						go func(k int) {
							defer wg.Done()

							rwmutex.Lock()
							// Simulate work
							_ = k * k

							rwmutex.Unlock()
						}(j)
					}

					wg.Wait()
				}
			})
		})
	}
}

// BenchmarkMutexVsRWMutex_WorkloadIntensity benchmarks with different work intensities.
func BenchmarkMutexVsRWMutex_WorkloadIntensity(b *testing.B) {
	workloads := []struct {
		name     string
		workSize int
	}{
		{"LightWork", 1},
		{"MediumWork", 100},
		{"HeavyWork", 10000},
	}

	for _, workload := range workloads {
		b.Run(workload.name, func(b *testing.B) {
			// Test with regular mutex
			b.Run("Mutex", func(b *testing.B) {
				mutex := NewMutexWithLock()

				var wg sync.WaitGroup

				workers := 8

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(workers)

					for j := 0; j < workers; j++ {
						go func() {
							defer wg.Done()

							mutex.Lock()
							// Simulate work of varying intensity
							sum := 0
							for k := 0; k < workload.workSize; k++ {
								sum += k * k
							}

							_ = sum

							mutex.Unlock()
						}()
					}

					wg.Wait()
				}
			})

			// Test with RWMutex (readers)
			b.Run("RWMutex_Readers", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				workers := 8

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(workers)

					for j := 0; j < workers; j++ {
						go func() {
							defer wg.Done()

							rwmutex.RLock()
							// Simulate work of varying intensity
							sum := 0
							for k := 0; k < workload.workSize; k++ {
								sum += k * k
							}

							_ = sum

							rwmutex.RUnlock()
						}()
					}

					wg.Wait()
				}
			})

			// Test with RWMutex (writers)
			b.Run("RWMutex_Writers", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				workers := 8

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(workers)

					for j := 0; j < workers; j++ {
						go func() {
							defer wg.Done()

							rwmutex.Lock()
							// Simulate work of varying intensity
							sum := 0
							for k := 0; k < workload.workSize; k++ {
								sum += k * k
							}

							_ = sum

							rwmutex.Unlock()
						}()
					}

					wg.Wait()
				}
			})
		})
	}
}

// BenchmarkMutexVsRWMutex_RealWorldScenarios benchmarks realistic scenarios.
func BenchmarkMutexVsRWMutex_RealWorldScenarios(b *testing.B) {
	scenarios := []struct {
		name     string
		readers  int
		writers  int
		readFreq int // How often reads happen relative to writes
	}{
		{"CacheLike_100R_1W", 100, 1, 100},
		{"DatabaseLike_50R_5W", 50, 5, 10},
		{"ConfigLike_10R_1W", 10, 1, 10},
		{"LogLike_1R_10W", 1, 10, 1},
		{"Balanced_10R_10W", 10, 10, 1},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			// Test with regular mutex
			b.Run("Mutex", func(b *testing.B) {
				mutex := NewMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.readers + scenario.writers)

					// Start readers
					for j := 0; j < scenario.readers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.readFreq; k++ {
								mutex.Lock()
								// Simulate read operation
								_ = k * k

								mutex.Unlock()
							}
						}()
					}

					// Start writers
					for j := 0; j < scenario.writers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < 1; k++ { // Writers do less work
								mutex.Lock()
								// Simulate write operation
								_ = k * k

								mutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})

			// Test with RWMutex
			b.Run("RWMutex", func(b *testing.B) {
				rwmutex := NewRWMutexWithLock()

				var wg sync.WaitGroup

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					wg.Add(scenario.readers + scenario.writers)

					// Start readers
					for j := 0; j < scenario.readers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < scenario.readFreq; k++ {
								rwmutex.RLock()
								// Simulate read operation
								_ = k * k

								rwmutex.RUnlock()
							}
						}()
					}

					// Start writers
					for j := 0; j < scenario.writers; j++ {
						go func() {
							defer wg.Done()

							for k := 0; k < 1; k++ { // Writers do less work
								rwmutex.Lock()
								// Simulate write operation
								_ = k * k

								rwmutex.Unlock()
							}
						}()
					}

					wg.Wait()
				}
			})
		})
	}
}
