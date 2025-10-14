// Copyright 2025 samber.
//
// Licensed as an Enterprise License (the "License"); you may not use
// this file except in compliance with the License. You may obtain
// a copy of the License at:
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.ee.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package introspection

import (
	"context"

	"github.com/samber/ro"
)

//nolint:unusedparams
func pipe2[A any, B any, C any](source ro.Observable[A], operator1 func(ro.Observable[A]) ro.Observable[B], operator2 func(ro.Observable[B]) ro.Observable[C]) (*FunDesc, error) {
	return GetFunctionDescription(0, 1)
}

func distinct[T comparable]() func(ro.Observable[T]) ro.Observable[T] {
	return func(source ro.Observable[T]) ro.Observable[T] {
		return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[T]) ro.Teardown {
			seen := map[T]struct{}{}

			sub := source.SubscribeWithContext(
				subscriberCtx,
				ro.NewObserverWithContext(
					func(ctx context.Context, value T) {
						if _, ok := seen[value]; !ok {
							destination.NextWithContext(ctx, value)
							seen[value] = struct{}{}
						}
					},
					destination.ErrorWithContext,
					destination.CompleteWithContext,
				),
			)

			return sub.Unsubscribe
		})
	}
}
