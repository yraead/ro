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
	"github.com/samber/lo"
	"github.com/samber/ro"
)

var bridgeInstance *bridge

func initStreams() {
	bridgeInstance = newEchange()
}

func closeStreams() {
	bridgeInstance.Close()
}

func newEchange() *bridge {
	e := &bridge{
		upstream:      ro.NewPublishSubject[lo.Tuple2[string, string]](),
		downstream:    ro.NewReplaySubject[lo.Tuple2[string, string]](10_000),
		subscriptions: ro.NewSubscription(nil),
	}

	// websocket->redis (downstream)
	e.subscriptions.AddUnsubscribable(
		e.upstream.
			Subscribe(
				ro.OnNext(func(msg lo.Tuple2[string, string]) {
					publishSink(msg.A, msg.B)
				}),
			),
	)

	// redis->websocket (upstream)
	e.subscriptions.AddUnsubscribable(
		ro.NewObservable(subscribeSource).
			Subscribe(e.downstream),
	)

	return e
}

type bridge struct {
	upstream      ro.Subject[lo.Tuple2[string, string]]
	downstream    ro.Subject[lo.Tuple2[string, string]]
	subscriptions ro.Subscription
}

func (e *bridge) Publish(roomID string, msg string) {
	e.upstream.Next(lo.T2(roomID, msg))
}

func (e *bridge) Subscribe(roomID string, destination ro.Observer[string]) ro.Subscription {
	sub := ro.Pipe2(
		e.downstream.AsObservable(),
		ro.Filter(func(msg lo.Tuple2[string, string]) bool {
			// exclude messages from other rooms
			return msg.A == roomID
		}),
		ro.Map(func(msg lo.Tuple2[string, string]) string {
			return msg.B
		}),
	).Subscribe(destination)

	e.subscriptions.AddUnsubscribable(sub)

	return sub
}

func (e *bridge) Close() {
	e.upstream.Complete()
	e.downstream.Complete()
	e.subscriptions.Unsubscribe()
}
