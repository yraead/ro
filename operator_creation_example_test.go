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

package ro

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func ExampleNewObservable_ok() {
	observable := NewObservable(func(observer Observer[int]) Teardown {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		observer.Next(4)
		observer.Complete()

		return nil
	})

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Completed
}

func ExampleNewObservable_error() {
	observable := NewObservable(func(observer Observer[int]) Teardown {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		observer.Error(assert.AnError)
		observer.Next(4)

		return nil
	})

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleOf() {
	observable := Of(1, 2, 3)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleStart() {
	observable := Start(func() int {
		fmt.Println("Start!")
		return 42
	})

	subscription1 := observable.Subscribe(PrintObserver[int]())
	subscription2 := observable.Subscribe(PrintObserver[int]())

	subscription1.Wait() // Note: using .Wait() is not recommended.
	subscription2.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Start!
	// Next: 42
	// Completed
	// Start!
	// Next: 42
	// Completed
}

func ExampleJust() {
	observable := Just(1, 2, 3)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleTimer() {
	observable := Timer(10 * time.Millisecond)

	subscription := observable.Subscribe(PrintObserver[time.Duration]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 10ms
	// Completed
}

func ExampleInterval() {
	observable := Interval(100 * time.Millisecond)

	subscription := observable.Subscribe(PrintObserver[int64]())

	time.Sleep(250 * time.Millisecond)
	subscription.Unsubscribe() // "Completed" event is not transmitted

	// Output:
	// Next: 0
	// Next: 1
}

func ExampleIntervalWithInitial() {
	observable := IntervalWithInitial(50*time.Millisecond, 100*time.Millisecond)

	subscription := observable.Subscribe(PrintObserver[int64]())

	time.Sleep(300 * time.Millisecond)
	subscription.Unsubscribe() // "Completed" event is not transmitted

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
}

func ExampleRange() {
	observable := Range(0, 5)

	subscription := observable.Subscribe(PrintObserver[int64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Completed
}

func ExampleRangeWithStep() {
	observable := RangeWithStep(0, 5, 0.5)

	subscription := observable.Subscribe(PrintObserver[float64]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 0
	// Next: 0.5
	// Next: 1
	// Next: 1.5
	// Next: 2
	// Next: 2.5
	// Next: 3
	// Next: 3.5
	// Next: 4
	// Next: 4.5
	// Completed
}

func ExampleRangeWithInterval() {
	observable := RangeWithInterval(0, 5, 10*time.Millisecond)

	subscription := observable.Subscribe(PrintObserver[int64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 0
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Completed
}

func ExampleRangeWithStepAndInterval() {
	observable := RangeWithStepAndInterval(0, 5, 0.5, 10*time.Millisecond)

	subscription := observable.Subscribe(PrintObserver[float64]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 0
	// Next: 0.5
	// Next: 1
	// Next: 1.5
	// Next: 2
	// Next: 2.5
	// Next: 3
	// Next: 3.5
	// Next: 4
	// Next: 4.5
	// Completed
}

func ExampleRepeat() {
	observable := Repeat(42, 3)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 42
	// Next: 42
	// Next: 42
	// Completed
}

func ExampleRepeatWithInterval() {
	// @TODO: implment
}

func ExampleFromChannel() {
	ch := make(chan int, 10)
	observable := FromChannel(ch)

	subscription := observable.Subscribe(PrintObserver[int]())

	ch <- 1

	ch <- 2

	ch <- 3

	close(ch)

	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleFromSlice() {
	observable := FromSlice([]int{1, 2, 3})

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Completed
}

func ExampleEmpty() {
	observable := Empty[int]()

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Completed
}

func ExampleNever() {
	observable := Never()

	subscription := observable.Subscribe(PrintObserver[struct{}]())

	time.Sleep(10 * time.Millisecond)
	subscription.Unsubscribe()

	// Output:
}

func ExampleThrow() {
	observable := Throw[int](assert.AnError)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleDefer() {
	// will capture current date time
	observable1 := Of(time.Now())

	// will capture date time at the moment of subscription
	observable2 := Defer(func() Observable[time.Time] {
		return Of(time.Now())
	})

	subscription := Concat(observable1, observable2).Subscribe(NoopObserver[time.Time]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
}

func ExampleFuture_ok() {
	observable := Future(func() (int, error) {
		req, err := http.NewRequest("GET", "https://postman-echo.com/get", nil)
		if err != nil {
			return 0, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}

		defer res.Body.Close()

		// For some reason, removing the 2 following lines causes
		// the example to fail (see goleak).
		// See https://github.com/uber-go/goleak/issues/102
		_, _ = io.ReadAll(res.Body)

		defer http.DefaultClient.CloseIdleConnections()

		return 42, nil
	})

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Next: 42
	// Completed
}

func ExampleFuture_error() {
	observable := Future(func() (int, error) {
		req, err := http.NewRequest("", "", nil)
		if err != nil {
			return 0, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}

		defer res.Body.Close()

		// For some reason, removing the 2 following lines causes
		// the example to fail (see goleak).
		// See https://github.com/uber-go/goleak/issues/102
		_, _ = io.ReadAll(res.Body)

		defer http.DefaultClient.CloseIdleConnections()

		return 42, nil
	})

	subscription := observable.Subscribe(PrintObserver[int]())
	subscription.Wait() // Note: using .Wait() is not recommended.

	// Output:
	// Error: Get "": unsupported protocol scheme ""
}

func ExampleMerge_ok() {
	observable := Merge(
		RangeWithInterval(0, 2, 50*time.Millisecond),
		Pipe1(
			RangeWithInterval(10, 12, 50*time.Millisecond),
			Delay[int64](25*time.Millisecond),
		),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())

	time.Sleep(200 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: 0
	// Next: 10
	// Next: 1
	// Next: 11
	// Completed
}

func ExampleMerge_error() {
	observable := Merge(
		RangeWithInterval(0, 2, 50*time.Millisecond),
		Pipe1(
			Throw[int64](assert.AnError),
			Delay[int64](75*time.Millisecond),
		),
	)

	subscription := observable.Subscribe(PrintObserver[int64]())

	time.Sleep(100 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: 0
	// Error: assert.AnError general error for testing
}

func ExampleCombineLatest2_ok() {
	observable1 := Delay[int64](25 * time.Millisecond)(RangeWithInterval(1, 3, 50*time.Millisecond))
	observable2 := RangeWithInterval(3, 5, 50*time.Millisecond)

	observable := Pipe1(
		CombineLatest2(
			observable1,
			observable2,
		),
		Map(func(snapshot lo.Tuple2[int64, int64]) []int64 {
			return []int64{snapshot.A, snapshot.B}
		}),
	)

	subscription := observable.Subscribe(PrintObserver[[]int64]())

	time.Sleep(200 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: [1 3]
	// Next: [1 4]
	// Next: [2 4]
	// Completed
}

func ExampleCombineLatest2_error() {
	observable1 := NewObservable(func(observer Observer[int]) Teardown {
		go func() {
			time.Sleep(10 * time.Millisecond)
			observer.Next(1)
			observer.Error(assert.AnError)
		}()

		return nil
	})

	observable2 := NewObservable(func(observer Observer[int]) Teardown {
		go func() {
			observer.Next(2)
			observer.Complete()
		}()

		return nil
	})

	observable := Pipe1(
		CombineLatest2(
			observable1,
			observable2,
		),
		Map(func(snapshot lo.Tuple2[int, int]) []int {
			return []int{snapshot.A, snapshot.B}
		}),
	)

	subscription := observable.Subscribe(PrintObserver[[]int]())

	time.Sleep(30 * time.Millisecond)

	defer subscription.Unsubscribe()

	// Output:
	// Next: [1 2]
	// Error: assert.AnError general error for testing
}

func ExampleCombineLatest3_ok() {
	// @TODO: implement
}

func ExampleCombineLatest3_error() {
	// @TODO: implement
}

func ExampleCombineLatest4_ok() {
	// @TODO: implement
}

func ExampleCombineLatest4_error() {
	// @TODO: implement
}

func ExampleCombineLatest5_ok() {
	// @TODO: implement
}

func ExampleCombineLatest5_error() {
	// @TODO: implement
}

func ExampleCombineLatestAny_ok() {
	// @TODO: implement
}

func ExampleCombineLatestAny_error() {
	// @TODO: implement
}

func ExampleZip_ok() {
	// @TODO: implement
}

func ExampleZip_error() {
	// @TODO: implement
}

func ExampleZip2_ok() {
	observable := Zip2(
		Range(0, 10),
		Skip[int64](1)(Range(0, 4)),
	)

	subscription := observable.Subscribe(PrintObserver[lo.Tuple2[int64, int64]]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: {0 1}
	// Next: {1 2}
	// Next: {2 3}
	// Completed
}

func ExampleZip2_error() {
	observable := Zip2(
		Range(0, 10),
		Throw[int64](assert.AnError),
	)

	subscription := observable.Subscribe(PrintObserver[lo.Tuple2[int64, int64]]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleZip3_ok() {
	// @TODO: implement
}

func ExampleZip3_error() {
	// @TODO: implement
}

func ExampleZip4_ok() {
	// @TODO: implement
}

func ExampleZip4_error() {
	// @TODO: implement
}

func ExampleZip5_ok() {
	// @TODO: implement
}

func ExampleZip5_error() {
	// @TODO: implement
}

func ExampleZip6_ok() {
	// @TODO: implement
}

func ExampleZip6_error() {
	// @TODO: implement
}

func ExampleConcat_ok() {
	observable := Concat(
		Just(1, 2, 3),
		Just(4, 5, 6),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Next: 4
	// Next: 5
	// Next: 6
	// Completed
}

func ExampleConcat_error() {
	observable := Concat(
		Just(1, 2, 3),
		Throw[int](assert.AnError),
		Just(4, 5, 6),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: 1
	// Next: 2
	// Next: 3
	// Error: assert.AnError general error for testing
}

func ExampleRace_ok() {
	observable := Race(
		Delay[int](50*time.Millisecond)(Just(1, 2, 3)),
		Delay[int](20*time.Millisecond)(Just(4, 5, 6)),
		Delay[int](100*time.Millisecond)(Just(7, 8, 9)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: 4
	// Next: 5
	// Next: 6
	// Completed
}

func ExampleRace_error() {
	observable := Race(
		Delay[int](50*time.Millisecond)(Just(1, 2, 3)),
		Delay[int](20*time.Millisecond)(Throw[int](assert.AnError)),
		Delay[int](100*time.Millisecond)(Just(7, 8, 9)),
	)

	subscription := observable.Subscribe(PrintObserver[int]())
	defer subscription.Unsubscribe()

	time.Sleep(50 * time.Millisecond)

	// Output:
	// Error: assert.AnError general error for testing
}

func ExampleAmb_ok() {
	// @TODO: implement
}

func ExampleAmb_error() {
	// @TODO: implement
}

func ExampleRandIntN() {
	// @TODO: implement
}

func ExampleRandFloat64() {
	// @TODO: implement
}
