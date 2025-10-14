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


package roratelimit

import (
	"context"
	"time"

	"github.com/samber/ro"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func ExampleNewRateLimiter() {
	// Create a rate limiter with 5 requests per second
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  5,
	}
	limiter := limiter.New(store, rate)

	// Rate limit by user ID
	observable := ro.Pipe1(
		ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
		NewRateLimiter[string](limiter, func(userID string) string {
			return userID
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: user1
	// Next: user2
	// Next: user1
	// Next: user3
	// Next: user1
	// Next: user2
	// Next: user1
	// Completed
}

func ExampleNewRateLimiter_withStructs() {
	// Create a rate limiter with 3 requests per minute
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  3,
	}
	limiter := limiter.New(store, rate)

	type Request struct {
		UserID string
		Action string
		Data   string
	}

	// Rate limit by user ID
	observable := ro.Pipe1(
		ro.Just(
			Request{UserID: "user1", Action: "login", Data: "data1"},
			Request{UserID: "user2", Action: "login", Data: "data2"},
			Request{UserID: "user1", Action: "logout", Data: "data3"},
			Request{UserID: "user3", Action: "login", Data: "data4"},
			Request{UserID: "user1", Action: "update", Data: "data5"},
			Request{UserID: "user2", Action: "logout", Data: "data6"},
		),
		NewRateLimiter[Request](limiter, func(req Request) string {
			return req.UserID
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[Request]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: {user1 login data1}
	// Next: {user2 login data2}
	// Next: {user1 logout data3}
	// Next: {user3 login data4}
	// Next: {user1 update data5}
	// Next: {user2 logout data6}
	// Completed
}

func ExampleNewRateLimiter_withIPAddress() {
	// Create a rate limiter with 10 requests per minute
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  10,
	}
	limiter := limiter.New(store, rate)

	type APIRequest struct {
		IPAddress string
		Endpoint  string
		Method    string
	}

	// Rate limit by IP address
	observable := ro.Pipe1(
		ro.Just(
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/posts", Method: "POST"},
			APIRequest{IPAddress: "192.168.1.3", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/comments", Method: "GET"},
		),
		NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
			return req.IPAddress
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[APIRequest]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: {192.168.1.1 /api/users GET}
	// Next: {192.168.1.2 /api/users GET}
	// Next: {192.168.1.1 /api/posts POST}
	// Next: {192.168.1.3 /api/users GET}
	// Next: {192.168.1.1 /api/comments GET}
	// Completed
}

func ExampleNewRateLimiter_withEndpoint() {
	// Create a rate limiter with 2 requests per second
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  2,
	}
	limiter := limiter.New(store, rate)

	type APIRequest struct {
		IPAddress string
		Endpoint  string
		Method    string
	}

	// Rate limit by endpoint
	observable := ro.Pipe1(
		ro.Just(
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/posts", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "POST"},
			APIRequest{IPAddress: "192.168.1.3", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/comments", Method: "GET"},
		),
		NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
			return req.Endpoint
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[APIRequest]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: {192.168.1.1 /api/users GET}
	// Next: {192.168.1.2 /api/posts GET}
	// Next: {192.168.1.1 /api/users POST}
	// Next: {192.168.1.1 /api/comments GET}
	// Completed
}

func ExampleNewRateLimiter_withCompositeKey() {
	// Create a rate limiter with 5 requests per minute
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  5,
	}
	limiter := limiter.New(store, rate)

	type APIRequest struct {
		IPAddress string
		Endpoint  string
		Method    string
	}

	// Rate limit by IP + endpoint combination
	observable := ro.Pipe1(
		ro.Just(
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/users", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/posts", Method: "GET"},
			APIRequest{IPAddress: "192.168.1.1", Endpoint: "/api/users", Method: "POST"},
			APIRequest{IPAddress: "192.168.1.2", Endpoint: "/api/posts", Method: "GET"},
		),
		NewRateLimiter[APIRequest](limiter, func(req APIRequest) string {
			return req.IPAddress + ":" + req.Endpoint
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[APIRequest]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: {192.168.1.1 /api/users GET}
	// Next: {192.168.1.2 /api/users GET}
	// Next: {192.168.1.1 /api/posts GET}
	// Next: {192.168.1.1 /api/users POST}
	// Next: {192.168.1.2 /api/posts GET}
	// Completed
}

func ExampleNewRateLimiter_withErrorHandling() {
	// Create a rate limiter with 3 requests per second
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  3,
	}
	limiter := limiter.New(store, rate)

	// Rate limit with error handling
	observable := ro.Pipe1(
		ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
		NewRateLimiter[string](limiter, func(userID string) string {
			return userID
		}),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(value string) {
				// Handle successful rate-limited value
			},
			func(err error) {
				// Handle rate limiting error
				// This could be due to:
				// - Store errors
				// - Context cancellation
				// - Other limiter errors
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()
}

func ExampleNewRateLimiter_withContext() {
	// Create a rate limiter with 2 requests per second
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  2,
	}
	limiter := limiter.New(store, rate)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Rate limit with context
	observable := ro.Pipe1(
		ro.Just("user1", "user2", "user1", "user3", "user1", "user2", "user1"),
		NewRateLimiter[string](limiter, func(userID string) string {
			return userID
		}),
	)

	subscription := observable.SubscribeWithContext(
		ctx,
		ro.NewObserverWithContext(
			func(ctx context.Context, value string) {
				// Handle rate-limited value with context
			},
			func(ctx context.Context, err error) {
				// Handle error with context
			},
			func(ctx context.Context) {
				// Handle completion with context
			},
		),
	)
	defer subscription.Unsubscribe()
}
