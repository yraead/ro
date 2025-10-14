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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectableObservable(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a := []int{}
	b := []string{}

	source := func(destination Observer[int]) Teardown {
		destination.Next(1)
		destination.Next(2)
		destination.Next(3)
		destination.Complete()

		return nil
	}

	connectable, ok := NewConnectableObservable(source).(*connectableObservableImpl[int])

	is.True(ok)

	is.True(connectable.config.ResetOnDisconnect)
	is.NotNil(connectable.config.Connector)
	is.NotNil(connectable.source)
	is.Nil(connectable.subscription)

	sub1 := connectable.Subscribe(OnNext(func(item int) {
		a = append(a, item)
	}))
	sub2 := connectable.Subscribe(OnNext(func(item int) {
		b = append(b, strconv.Itoa(item))
	}))

	is.Nil(connectable.subscription)
	is.False(sub1.IsClosed())
	is.False(sub2.IsClosed())

	sub := connectable.Connect()
	is.True(connectable.subscription.IsClosed())
	is.True(sub.IsClosed())
	is.True(sub1.IsClosed())
	is.True(sub2.IsClosed())

	is.Equal([]int{1, 2, 3}, a)
	is.Equal([]string{"1", "2", "3"}, b)
}

func TestNewConnectableObservableWithConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a := []int{}
	b := []string{}

	source := func(destination Observer[int]) Teardown {
		destination.Next(1)
		destination.Next(2)
		destination.Next(3)
		destination.Complete()

		return nil
	}

	config := ConnectableConfig[int]{
		Connector:         NewSubject[int],
		ResetOnDisconnect: true,
	}
	connectable, ok := NewConnectableObservableWithConfig(source, config).(*connectableObservableImpl[int])

	is.True(ok)

	is.True(connectable.config.ResetOnDisconnect)
	is.NotNil(connectable.config.Connector)
	is.NotNil(connectable.source)
	is.Nil(connectable.subscription)

	sub1 := connectable.Subscribe(OnNext(func(item int) {
		a = append(a, item)
	}))
	sub2 := connectable.Subscribe(OnNext(func(item int) {
		b = append(b, strconv.Itoa(item))
	}))

	is.Nil(connectable.subscription)
	is.False(sub1.IsClosed())
	is.False(sub2.IsClosed())

	sub := connectable.Connect()
	is.True(connectable.subscription.IsClosed())
	is.True(sub.IsClosed())
	is.True(sub1.IsClosed())
	is.True(sub2.IsClosed())

	is.Equal([]int{1, 2, 3}, a)
	is.Equal([]string{"1", "2", "3"}, b)
}

func TestConnectable(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a := []int{}
	b := []int{}
	c := []string{}

	source := TapOnNext(func(value int) {
		a = append(a, value*2)
	})(Of(1, 2, 3))

	connectable, ok := Connectable(source).(*connectableObservableImpl[int])

	is.True(ok)

	is.True(connectable.config.ResetOnDisconnect)
	is.NotNil(connectable.config.Connector)
	is.NotNil(connectable.source)
	is.Nil(connectable.subscription)

	sub1 := connectable.Subscribe(OnNext(func(item int) {
		b = append(b, item)
	}))
	sub2 := connectable.Subscribe(OnNext(func(item int) {
		c = append(c, strconv.Itoa(item))
	}))

	is.Nil(connectable.subscription)
	is.False(sub1.IsClosed())
	is.False(sub2.IsClosed())

	sub := connectable.Connect()
	is.True(connectable.subscription.IsClosed())
	is.True(sub.IsClosed())
	is.True(sub1.IsClosed())
	is.True(sub2.IsClosed())

	is.Equal([]int{2, 4, 6}, a)
	is.Equal([]int{1, 2, 3}, b)
	is.Equal([]string{"1", "2", "3"}, c)
}

func TestConnectableWithConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a := []int{}
	b := []int{}
	c := []string{}

	source := TapOnNext(func(value int) {
		a = append(a, value*2)
	})(Of(1, 2, 3))

	config := ConnectableConfig[int]{
		Connector:         NewSubject[int],
		ResetOnDisconnect: true,
	}
	connectable, ok := ConnectableWithConfig(source, config).(*connectableObservableImpl[int])

	is.True(ok)

	is.True(connectable.config.ResetOnDisconnect)
	is.NotNil(connectable.config.Connector)
	is.NotNil(connectable.source)
	is.Nil(connectable.subscription)

	sub1 := connectable.Subscribe(OnNext(func(item int) {
		b = append(b, item)
	}))
	sub2 := connectable.Subscribe(OnNext(func(item int) {
		c = append(c, strconv.Itoa(item))
	}))

	is.Nil(connectable.subscription)
	is.False(sub1.IsClosed())
	is.False(sub2.IsClosed())

	sub := connectable.Connect()
	is.True(connectable.subscription.IsClosed())
	is.True(sub.IsClosed())
	is.True(sub1.IsClosed())
	is.True(sub2.IsClosed())

	is.Equal([]int{2, 4, 6}, a)
	is.Equal([]int{1, 2, 3}, b)
	is.Equal([]string{"1", "2", "3"}, c)
}
