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

	"github.com/samber/lo"
)

// Observable is the producer of values. It is the source of values that are
// emitted to Observers.
// Observable is a representation of any set of values over any amount of time.
//
// The primary method of an Observable is subscribe, which is used to attach an
// Observer to the Observable. Once an Observer is subscribed, the Observable
// may begin to emit items to the Observer. An Observable may emit any number
// of items (including zero items), then may either complete or error, but not
// both. Upon completion or error, the Observable will not emit any more items.
//
// An Observable may call an Observer's methods synchronously or asynchronously.
//
// An Observable is not a stream. It is a factory for streams.
type Observable[T any] interface {
	// Subscribe subscribes an Observer to the Observable. The Observer will begin
	// to receive items emitted by the Observable. The Observer may receive any
	// number of items (including zero items), then may either complete or error,
	// but not both. Upon completion or error, the Observer will not receive any
	// more items.
	//
	// The Subscribe method returns a Subscription that can be used to unsubscribe
	// the Observer from the Observable. The Subscription may be used to cancel the
	// subscription, and to wait for the subscription to complete.
	//
	// The Subscription might be already disposed when the Subscribe method returns.
	// In this case, the Teardown function is not called.
	//
	// The Subscribe method may call the Observer's methods synchronously or
	// asynchronously. The Observer is responsible for handling concurrency and
	// synchronization.
	Subscribe(destination Observer[T]) Subscription
	SubscribeWithContext(ctx context.Context, destination Observer[T]) Subscription
}

var _ Observable[int] = (*observableImpl[int])(nil)

// NewObservable creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is not safe for concurrent use.
func NewObservable[T any](subscribe func(destination Observer[T]) Teardown) Observable[T] {
	return NewSafeObservable(subscribe)
}

// NewSafeObservable creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is not safe for concurrent use.
func NewSafeObservable[T any](subscribe func(destination Observer[T]) Teardown) Observable[T] {
	return NewObservableWithConcurrencyMode(
		func(ctx context.Context, destination Observer[T]) Teardown {
			return subscribe(destination)
		},
		ConcurrencyModeSafe,
	)
}

// NewUnsafeObservable creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is not safe for concurrent use.
func NewUnsafeObservable[T any](subscribe func(destination Observer[T]) Teardown) Observable[T] {
	return NewObservableWithConcurrencyMode(
		func(ctx context.Context, destination Observer[T]) Teardown {
			return subscribe(destination)
		},
		ConcurrencyModeUnsafe,
	)
}

// NewEventualySafeObservable creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is safe for concurrent use, but concurrent messages are dropped.
func NewEventualySafeObservable[T any](subscribe func(destination Observer[T]) Teardown) Observable[T] {
	return NewObservableWithConcurrencyMode(
		func(ctx context.Context, destination Observer[T]) Teardown {
			return subscribe(destination)
		},
		ConcurrencyModeEventualySafe,
	)
}

// NewObservableWithContext creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is not safe for concurrent use.
func NewObservableWithContext[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown) Observable[T] {
	return NewSafeObservableWithContext(subscribe)
}

// NewSafeObservableWithContext creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is not safe for concurrent use.
func NewSafeObservableWithContext[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown) Observable[T] {
	return NewObservableWithConcurrencyMode(subscribe, ConcurrencyModeSafe)
}

// NewUnsafeObservableWithContext creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is not safe for concurrent use.
func NewUnsafeObservableWithContext[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown) Observable[T] {
	return NewObservableWithConcurrencyMode(subscribe, ConcurrencyModeUnsafe)
}

// NewEventualySafeObservableWithContext creates a new Observable. The subscribe function is called when
// the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error,
// but not both. Upon completion or error, the Observable will not emit any more
// items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// This method is safe for concurrent use, but concurrent messages are dropped.
func NewEventualySafeObservableWithContext[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown) Observable[T] {
	return NewObservableWithConcurrencyMode(subscribe, ConcurrencyModeEventualySafe)
}

// NewObservableWithConcurrencyMode creates a new Observable with the given concurrency mode.
// The subscribe function is called when the Observable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error, but not both. Upon completion or error, the Observable will not emit any more items.
//
// The subscribe function should return a Teardown function that will be called
// when the Subscription is unsubscribed. The Teardown function should clean up
// any resources created during the subscription.
//
// The subscribe function may return a Teardown function that does nothing, if
// no cleanup is necessary. In this case, the Teardown function should return nil.
//
// The Observable will use the given concurrency mode.
//
// It is rarely used as a public API.
func NewObservableWithConcurrencyMode[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown, mode ConcurrencyMode) Observable[T] {
	return &observableImpl[T]{
		mode:      mode,
		subscribe: subscribe,
	}
}

type observableImpl[T any] struct {
	mode      ConcurrencyMode
	subscribe func(ctx context.Context, destination Observer[T]) Teardown
}

// Subscribe subscribes an Observer to the Observable. The Observer will begin
// to receive items emitted by the Observable. The Observer may receive any
// number of items (including zero items), then may either complete or error,
// but not both. Upon completion or error, the Observer will not receive any
// more items.
//
// The Subscribe method returns a Subscription that can be used to unsubscribe
// the Observer from the Observable. The Subscription may be used to cancel the
// subscription, and to wait for the subscription to complete.
//
// The Subscription might be already disposed when the Subscribe method returns.
// In this case, the Teardown function is not called.
//
// The Subscribe method may call the Observer's methods synchronously or
// asynchronously. The Observer is responsible for handling concurrency and
// synchronization.
func (s *observableImpl[T]) Subscribe(destination Observer[T]) Subscription {
	return s.SubscribeWithContext(context.Background(), destination)
}

// SubscribeWithContext subscribes an Observer to the Observable. The Observer will begin
// to receive items emitted by the Observable. The Observer may receive any
// number of items (including zero items), then may either complete or error,
// but not both. Upon completion or error, the Observer will not receive any
// more items.
//
// The Subscribe method returns a Subscription that can be used to unsubscribe
// the Observer from the Observable. The Subscription may be used to cancel the
// subscription, and to wait for the subscription to complete.
//
// The Subscription might be already disposed when the Subscribe method returns.
// In this case, the Teardown function is not called.
//
// The Subscribe method may call the Observer's methods synchronously or
// asynchronously. The Observer is responsible for handling concurrency and
// synchronization.
func (s *observableImpl[T]) SubscribeWithContext(ctx context.Context, destination Observer[T]) Subscription {
	subscription := NewSubscriberWithConcurrencyMode(destination, s.mode)

	lo.TryCatchWithErrorValue(
		func() error {
			// Warning: here, we are catching panic in subscription.Add.
			// I'm not sure if it's a good idea.
			subscription.Add(s.subscribe(ctx, subscription))
			return nil
		},
		func(e any) {
			err := recoverValueToError(e)
			subscription.ErrorWithContext(ctx, newObservableError(err))
			subscription.Unsubscribe()
		},
	)

	return subscription
}

// Collect collects all values emitted by the source Observable and returns them
// as a slice. It waits for the source Observable to complete before returning.
// If the source Observable emits an error, the error is returned along with the
// values collected so far.
func Collect[T any](obs Observable[T]) ([]T, error) {
	v, _, err := CollectWithContext(context.Background(), obs)
	return v, err
}

// CollectWithContext collects all values emitted by the source Observable and returns them
// as a slice. It waits for the source Observable to complete before returning.
// If the source Observable emits an error, the error is returned along with the
// values collected so far.
// @TODO: return more values, such as (isCanceled bool) or (duration time.Duration) ?
func CollectWithContext[T any](ctx context.Context, obs Observable[T]) ([]T, context.Context, error) {
	values := []T{}

	var lastCtx context.Context
	var err error

	sub := obs.SubscribeWithContext(
		ctx,
		NewObserverWithContext(
			func(ctx context.Context, value T) {
				values = append(values, value)
			},
			func(ctx context.Context, thrown error) {
				err = thrown
				lastCtx = ctx
			},
			func(ctx context.Context) {
				lastCtx = ctx
			},
		),
	)

	sub.Wait() // Note: using .Wait() is not recommended.

	return values, lastCtx, err
}
