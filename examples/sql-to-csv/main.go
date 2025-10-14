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
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/samber/ro"
	rocsv "github.com/samber/ro/plugins/encoding/csv"
)

type User struct {
	ID   string
	Name string
}

var mockUsers = []User{
	{ID: "1", Name: "Alice"},
	{ID: "2", Name: "Bob"},
	{ID: "3", Name: "Charlie"},
	{ID: "4", Name: "David"},
	{ID: "5", Name: "Eve"},
	{ID: "6", Name: "Frank"},
	{ID: "7", Name: "Grace"},
	{ID: "8", Name: "Helen"},
	{ID: "9", Name: "Ivy"},
	{ID: "10", Name: "Jack"},
}

func getUsers(index int64) ([]User, error) {
	const batchSize = 5

	offset := batchSize * index
	end := batchSize*index + batchSize
	return lo.Slice(mockUsers, int(offset), int(end)), nil
}

// Define a pipeline to query users from a database and write them to stdout as csv.
var pipeline = ro.Pipe7(
	// Up to 100 page retrieval in database.
	ro.Range(0, 100),
	// Query database.
	// Convert database error into observer.Error() notification.
	ro.MapErr(getUsers),
	// Retry on database error.
	ro.RetryWithConfig[[]User](ro.RetryConfig{
		MaxRetries: 2,
		Delay:      5 * time.Second,
	}),
	// Complete stream when database response is empty.
	ro.TakeWhile(func(users []User) bool {
		return len(users) > 0
	}),
	// Convert Observable of []User into Observable of User.
	ro.Flatten[User](),
	// Convert User into a csv row ([]string).
	ro.Map(func(user User) []string {
		return []string{user.ID, user.Name}
	}),
	// Add csv header.
	ro.StartWith([]string{"ID", "Name"}),
	// Write csv row to stdout.
	rocsv.NewCSVWriter(
		csv.NewWriter(os.Stdout),
	),
)

func main() {
	subscription := pipeline.Subscribe(
		// Log SQL and writer errors.
		ro.OnError[int](func(err error) {
			log.Println(err.Error())
		}),
	)

	// Optional, since the pipeline will complete itself, in a blocking way.
	defer subscription.Unsubscribe()
}
