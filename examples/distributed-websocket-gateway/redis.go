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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/samber/ro"
)

var (
	rdb           *redis.Client
	pubsubChannel string = "exchange"
)

func init() {
	fmt.Println("Connecting to Redis...")

	url := "redis://@localhost:6379/0"
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opts)
}

func publishSink(roomID string, msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res := rdb.Publish(ctx, pubsubChannel+"."+roomID, msg)
	return res.Err()
}

func subscribeSource(destination ro.Observer[lo.Tuple2[string, string]]) ro.Teardown {
	res := rdb.PSubscribe(context.Background(), pubsubChannel+".*")
	ch := res.Channel()

	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					destination.Complete()
					return
				}

				roomID := strings.Replace(msg.Channel, pubsubChannel+".", "", 1)
				destination.Next(lo.T2(roomID, msg.Payload))
			}
		}
	}()

	return func() {
		res.Close()
		rdb.Close()
	}
}
