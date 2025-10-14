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


package ropsi

import (
	"time"

	psinotifier "github.com/samber/go-psi"
	"github.com/samber/ro"
)

func ExampleNewPSINotifier() {
	// Monitor system pressure stall information (PSI)
	observable := NewPSINotifier(5 * time.Second)

	subscription := observable.Subscribe(ro.NoopObserver[psinotifier.PSIStatsResource]())
	defer subscription.Unsubscribe()

	// Let it run for a few seconds
	time.Sleep(10 * time.Second)
}
