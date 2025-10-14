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


package main

import (
	"time"

	"github.com/samber/lo"
	"github.com/samber/ro"
)

type result struct {
	Orders          []string
	PaymentHistory  []string
	UserPreferences map[string]any
}

func getOrders() ([]string, error) {
	// simulate a microservice request
	time.Sleep(1 * time.Second)
	return []string{"order-1", "order-2"}, nil
}

func getPaymentHistory() ([]string, error) {
	// simulate a microservice request
	time.Sleep(100 * time.Millisecond)
	return []string{"payment-1", "payment-2"}, nil
}

func getUserPreferences() (map[string]any, error) {
	// simulate a microservice request
	time.Sleep(10 * time.Millisecond)
	return map[string]any{"role": "user"}, nil
}

// Wait for all requests to complete and combine the results.
var pipeline = ro.Pipe1(
	ro.CombineLatest3(
		// Future is a helper function to run a function in a goroutine.
		// It returns an Observable that emits a single result: a value or an error.
		ro.Future(getOrders),
		ro.Future(getPaymentHistory),
		ro.Future(getUserPreferences),
	),
	ro.Map(func(values lo.Tuple3[[]string, []string, map[string]any]) result {
		return result{
			Orders:          values.A,
			PaymentHistory:  values.B,
			UserPreferences: values.C,
		}
	}),
)

func main() {
	subscription := pipeline.Subscribe(
		// Print the result
		ro.PrintObserver[result](),
	)

	// Note: using .Wait() is not recommended.
	subscription.Wait()
}
