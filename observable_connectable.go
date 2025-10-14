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
)

// ConnectableObservable is an Observable that can be connected and disconnected.
// When connected, it will emit values to its observers.
//
// ConnectableObservable is useful when you want to share a single subscription to an Observable
// among multiple observers. This is useful when you want to multicast the values of an Observable.
type ConnectableObservable[T any] interface {
	Observable[T]

	// Connect connects the ConnectableObservable. When connected, the ConnectableObservable
	// will emit values to its observers. If the ConnectableObservable is already connected,
	// this method creates a new subscription and starts emitting values to its observers.
	//
	// The Connect method returns a Subscription that can be used to disconnect the
	// ConnectableObservable. The Subscription may be used to cancel the connection,
	// and to wait for the connection to complete.
	//
	// The Subscription might be already disposed when the Connect method returns.
	Connect() Subscription
	ConnectWithContext(ctx context.Context) Subscription
}

var (
	_ ConnectableObservable[int] = (*connectableObservableImpl[int])(nil)
	_ Observable[int]            = (*connectableObservableImpl[int])(nil)
)

// ConnectableConfig is the configuration for a ConnectableObservable.
type ConnectableConfig[T any] struct {
	Connector         func() Subject[T]
	ResetOnDisconnect bool
}

func defaultConnector[T any]() Subject[T] {
	return NewPublishSubject[T]()
}

// NewConnectableObservable creates a new ConnectableObservable. The subscribe function is called when
// the ConnectableObservable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error, but not both.
// Upon completion or error, the ConnectableObservable will not emit any more items.
//
// The ConnectableObservable will use the default connector, which is a PublishSubject.
// The ConnectableObservable will reset the source when disconnected. This means that
// when the ConnectableObservable is disconnected, it will create a new source when
// reconnected.
//
// If you want to use a different connector or change the reset behavior, use
// NewConnectableObservableWithConfig.
func NewConnectableObservable[T any](subscribe func(destination Observer[T]) Teardown) ConnectableObservable[T] {
	return newConnectableObservableImpl(
		NewObservable(subscribe),
		ConnectableConfig[T]{
			Connector:         defaultConnector[T],
			ResetOnDisconnect: true,
		},
	)
}

// NewConnectableObservableWithContext creates a new ConnectableObservable. The subscribe function is called when
// the ConnectableObservable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error, but not both.
// Upon completion or error, the ConnectableObservable will not emit any more items.
//
// The ConnectableObservable will use the default connector, which is a PublishSubject.
// The ConnectableObservable will reset the source when disconnected. This means that
// when the ConnectableObservable is disconnected, it will create a new source when
// reconnected.
//
// If you want to use a different connector or change the reset behavior, use
// NewConnectableObservableWithConfig.
func NewConnectableObservableWithContext[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown) ConnectableObservable[T] {
	return newConnectableObservableImpl(
		NewObservableWithContext(subscribe),
		ConnectableConfig[T]{
			Connector:         defaultConnector[T],
			ResetOnDisconnect: true,
		},
	)
}

// NewConnectableObservableWithConfig creates a new ConnectableObservable. The subscribe function is called when
// the ConnectableObservable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error, but not both.
// Upon completion or error, the ConnectableObservable will not emit any more items.
//
// The ConnectableObservable will use the given connector. The ConnectableObservable will reset
// the source when disconnected if ResetOnDisconnect is true. This means that when the
// ConnectableObservable is disconnected, it will create a new source when reconnected.
func NewConnectableObservableWithConfig[T any](subscribe func(destination Observer[T]) Teardown, config ConnectableConfig[T]) ConnectableObservable[T] {
	return newConnectableObservableImpl(
		NewObservable(subscribe),
		config,
	)
}

// NewConnectableObservableWithConfigAndContext creates a new ConnectableObservable. The subscribe function is called when
// the ConnectableObservable is subscribed to. The subscribe function is given an Observer,
// to which it may emit any number of items, then may either complete or error, but not both.
// Upon completion or error, the ConnectableObservable will not emit any more items.
//
// The ConnectableObservable will use the given connector. The ConnectableObservable will reset
// the source when disconnected if ResetOnDisconnect is true. This means that when the
// ConnectableObservable is disconnected, it will create a new source when reconnected.
func NewConnectableObservableWithConfigAndContext[T any](subscribe func(ctx context.Context, destination Observer[T]) Teardown, config ConnectableConfig[T]) ConnectableObservable[T] {
	return newConnectableObservableImpl(
		NewObservableWithContext(subscribe),
		config,
	)
}

// Connectable creates a new ConnectableObservable from an Observable. The ConnectableObservable
// will use the default connector, which is a PublishSubject. The ConnectableObservable will reset
// the source when disconnected. This means that when the ConnectableObservable is disconnected,
// it will create a new source when reconnected.
//
// If you want to use a different connector or change the reset behavior, use ConnectableWithConfig.
func Connectable[T any](source Observable[T]) ConnectableObservable[T] {
	return newConnectableObservableImpl(
		source,
		ConnectableConfig[T]{
			Connector:         defaultConnector[T],
			ResetOnDisconnect: true,
		},
	)
}

// ConnectableWithConfig creates a new ConnectableObservable from an Observable. The ConnectableObservable
// will use the given connector. The ConnectableObservable will reset the source when disconnected
// if ResetOnDisconnect is true. This means that when the ConnectableObservable is disconnected,
// it will create a new source when reconnected.
func ConnectableWithConfig[T any](source Observable[T], config ConnectableConfig[T]) ConnectableObservable[T] {
	return newConnectableObservableImpl(
		source,
		config,
	)
}

func newConnectableObservableImpl[T any](source Observable[T], config ConnectableConfig[T]) ConnectableObservable[T] {
	if config.Connector == nil {
		panic(ErrConnectableObservableMissingConnectorFactory)
	}

	return &connectableObservableImpl[T]{
		config:       config,
		source:       source,
		subject:      config.Connector(),
		subscription: nil,
	}
}

type connectableObservableImpl[T any] struct {
	mu           sync.Mutex
	config       ConnectableConfig[T]
	source       Observable[T]
	subject      Subject[T]
	subscription Subscription
}

// Connect connects the ConnectableObservable. When connected, the ConnectableObservable
// will emit values to its observers. If the ConnectableObservable is already connected,
// this method creates a new subscription and starts emitting values to its observers.
//
// The Connect method returns a Subscription that can be used to disconnect the
// ConnectableObservable. The Subscription may be used to cancel the connection,
// and to wait for the connection to complete.
//
// The Subscription might be already disposed when the Connect method returns.
func (s *connectableObservableImpl[T]) Connect() Subscription {
	return s.ConnectWithContext(context.Background())
}

// ConnectWithContext connects the ConnectableObservable. When connected, the ConnectableObservable
// will emit values to its observers. If the ConnectableObservable is already connected,
// this method creates a new subscription and starts emitting values to its observers.
//
// The Connect method returns a Subscription that can be used to disconnect the
// ConnectableObservable. The Subscription may be used to cancel the connection,
// and to wait for the connection to complete.
//
// The Subscription might be already disposed when the Connect method returns.
func (s *connectableObservableImpl[T]) ConnectWithContext(ctx context.Context) Subscription {
	s.mu.Lock()
	if s.subscription == nil || s.subscription.IsClosed() {
		s.subscription = s.source.SubscribeWithContext(ctx, s.subject)
		s.mu.Unlock()
		s.subscription.Add(func() {
			if s.config.ResetOnDisconnect {
				s.subject = s.config.Connector()
			}
		})
	} else {
		s.mu.Unlock()
	}

	return s.subscription
}

func (s *connectableObservableImpl[T]) Subscribe(observer Observer[T]) Subscription {
	return s.SubscribeWithContext(context.Background(), observer)
}

func (s *connectableObservableImpl[T]) SubscribeWithContext(ctx context.Context, observer Observer[T]) Subscription {
	return s.subject.SubscribeWithContext(ctx, observer)
}
