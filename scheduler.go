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
	"context"
	"sync"

	"github.com/samber/lo"
)

// NewScheduler just trolls other languages. 😈
// https://reactivex.io/documentation/scheduler.html
func NewScheduler() {
	panic(`Just kidding. 😇

Go is a modern programming language and doesn't need a scheduler.
It has multithreading built-in, so you don't need to worry about it.

If you come from PHP, Python, Javascript... welcome in the future! 🎉

However, if you really want to Schedule() something such as a conference,
a meetup, a business lunch or ...a date 😘, you can reach me here:
	-> https://twitter.com/samuelberthe
	-> https://bsky.app/samber

More seriously, if you're looking for building a highly-parallel, concurrent,
scalable and reliable stream processing app, use "samber/ro" package instead
of "samber/so".
`)
}

// SubscribeOn schedule the upstream flow to a different goroutine. Next, Error and Complete notifications
// are sent to a queue first, then the consumer consume this queue.
// SubscribeOn converts a push-based Observable into a pullable stream with backpressure capabilities.
//
// To schedule the downstream flow to a different goroutine, refer to SubscribeOn.
//
// When an Observable emits values faster than they can be consumed, SubscribeOn buffers these values
// in a queue of specified capacity. This allows downstream consumers to pull values at their own pace
// while managing backpressure from upstream emissions.
//
// Note: Once the buffer reaches its capacity, upstream emissions will block until space becomes
// available, effectively implementing backpressure control.
//
// @TODO: add a backpressure policy ? drop vs block.
func SubscribeOn[T any](bufferSize int) func(Observable[T]) Observable[T] {
	if bufferSize <= 0 {
		panic(ErrSubscribeOnWrongBufferSize)
	}

	return detachOn[T](bufferSize, true, false)
}

// ObserveOn schedule the downstream flow to a different goroutine. Next, Error and Complete notifications
// are sent to a queue first, then the consumer consume this queue.
// ObserveOn converts a push-based Observable into a pullable stream with backpressure capabilities.
//
// To schedule the upstream flow to a different goroutine, refer to SubscribeOn.
//
// When an Observable emits values faster than they can be consumed, ObserveOn buffers these values
// in a queue of specified capacity. This allows downstream consumers to pull values at their own pace
// while managing backpressure from upstream emissions.
//
// Note: Once the buffer reaches its capacity, upstream emissions will block until space becomes
// available, effectively implementing backpressure control.
//
// @TODO: add a backpressure policy ? drop vs block.
func ObserveOn[T any](bufferSize int) func(Observable[T]) Observable[T] {
	if bufferSize <= 0 {
		panic(ErrObserveOnWrongBufferSize)
	}

	return detachOn[T](bufferSize, false, true)
}

func detachOn[T any](bufferSize int, onUpstream, onDownstream bool) func(Observable[T]) Observable[T] {
	return func(source Observable[T]) Observable[T] {
		return NewObservableWithContext(func(subscriberCtx context.Context, destination Observer[T]) Teardown {
			ch := make(chan lo.Tuple2[context.Context, Notification[T]], bufferSize)

			once := sync.Once{}
			stop := func() {
				once.Do(func() {
					close(ch)
				})
			}

			subscriptions := NewSubscription(nil)

			consumeUpstream := func() {
				subscriptions.AddUnsubscribable(
					source.SubscribeWithContext(
						subscriberCtx,
						NewObserverWithContext(
							func(ctx context.Context, value T) {
								ch <- lo.T2(ctx, NewNotificationNext(value))
							},
							func(ctx context.Context, err error) {
								ch <- lo.T2(ctx, NewNotificationError[T](err))

								stop()
							},
							func(ctx context.Context) {
								ch <- lo.T2(ctx, NewNotificationComplete[T]())

								stop()
							},
						),
					),
				)
			}

			produceDownstream := func() {
				for notification := range ch {
					processNotificationWithContext(
						notification.A,
						notification.B,
						destination.NextWithContext,
						destination.ErrorWithContext,
						destination.CompleteWithContext,
					)
				}
			}

			// The goroutine could be used either on producer or consumer side.
			// 	* ObserveOn moves the goroutine on the consumer side.
			// 	* SubscribeOn moves the goroutine on the producer side.

			switch {
			case onUpstream:
				go recoverUnhandledError(func() {
					consumeUpstream()
				})

				produceDownstream()
			case onDownstream:
				go recoverUnhandledError(func() {
					produceDownstream()
				})

				consumeUpstream()
			default:
				panic(ErrDetachOnWrongMode)
			}

			return func() {
				subscriptions.Unsubscribe()
				stop()
			}
		})
	}
}
