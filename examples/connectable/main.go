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
	"fmt"
	"time"

	"github.com/samber/ro"
)

func sinkToDB(v map[string]map[string]string) {
	// ...
}

func sinkToStdout(ts map[string]map[string]string) {
	for k, v := range ts {
		println(k, v["4. close"])
	}
}

var source = ro.Connectable(
	ro.Pipe6(
		ro.Interval(5*time.Second),
		ro.While[int64](func() bool { return true }),
		ro.MapErr(func(v int64) (map[string]map[string]string, error) {
			fmt.Println("Fetch stock")
			return getMSFTStock()
		}),
		ro.TapOnError[map[string]map[string]string](func(err error) {
			println("Error:", err.Error())
		}),
		ro.Retry[map[string]map[string]string](),
		ro.TapOnSubscribe[map[string]map[string]string](func() {
			println("Subscribed")
		}),
		ro.TapOnFinalize[map[string]map[string]string](func() {
			println("Finalized")
		}),
	),
)

func main() {
	_ = source.Subscribe(ro.OnNext(sinkToDB))
	_ = source.Subscribe(ro.OnNext(sinkToStdout))

	subscription := source.Connect()

	// Note: using .Wait() is not recommended.
	subscription.Wait()
}
