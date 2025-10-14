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
	"fmt"
)

// Kind represents the kind of a Notification.
// It can be Next, Error, or Complete.
type Kind uint8

// String returns the string representation of a Kind.
func (k Kind) String() string {
	switch k {
	case KindNext:
		return "Next"
	case KindError:
		return "Error"
	case KindComplete:
		return "Complete"
	}

	panic("you shall not pass")
}

// Kind constants.
const (
	KindNext Kind = iota
	KindError
	KindComplete
)

// Notification represents a value emitted by an Observable. It can be a Next
// value, an Error, or a Complete signal. It is used to communicate between
// Observables and Observers. It is a generic type, so it can hold any value.
type Notification[T any] struct {
	Kind  Kind
	Value T
	Err   error
}

func (n Notification[T]) String() string {
	switch n.Kind {
	case KindNext:
		return fmt.Sprintf("Next(%+v)", n.Value)
	case KindError:
		if n.Err == nil {
			return "Error(nil)"
		}

		return fmt.Sprintf("Error(%s)", n.Err.Error())
	case KindComplete:
		return "Complete()"
	}

	panic("you shall not pass")
}

// NewNotificationNext creates a new Notification with a Next value.
func NewNotificationNext[T any](value T) Notification[T] {
	return Notification[T]{
		Kind:  KindNext,
		Value: value,
	}
}

// NewNotificationError creates a new Notification with an Error.
func NewNotificationError[T any](err error) Notification[T] {
	return Notification[T]{
		Kind: KindError,
		Err:  err,
	}
}

// NewNotificationComplete creates a new Notification with a Complete signal.
func NewNotificationComplete[T any]() Notification[T] {
	return Notification[T]{
		Kind: KindComplete,
	}
}

func processNotification[T any](n Notification[T], onNext func(T), onError func(error), onComplete func()) bool {
	switch n.Kind {
	case KindNext:
		onNext(n.Value)
		return true
	case KindError:
		onError(n.Err)
		return false
	case KindComplete:
		onComplete()
		return false
	}

	panic("you shall not pass")
}

func processNotificationWithContext[T any](ctx context.Context, n Notification[T], onNext func(context.Context, T), onError func(context.Context, error), onComplete func(context.Context)) bool {
	switch n.Kind {
	case KindNext:
		onNext(ctx, n.Value)
		return true
	case KindError:
		onError(ctx, n.Err)
		return false
	case KindComplete:
		onComplete(ctx)
		return false
	}

	panic("you shall not pass")
}

func processNotificationWithObserver[T any](n Notification[T], destination Observer[T]) bool {
	return processNotificationWithContext(
		context.Background(),
		n,
		destination.NextWithContext,
		destination.ErrorWithContext,
		destination.CompleteWithContext,
	)
}

func processNotificationWithObserverAndContext[T any](ctx context.Context, n Notification[T], destination Observer[T]) bool {
	return processNotificationWithContext(
		ctx,
		n,
		destination.NextWithContext,
		destination.ErrorWithContext,
		destination.CompleteWithContext,
	)
}
