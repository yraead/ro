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


package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/samber/lo"
	"github.com/samber/ro"
	rolicense "github.com/samber/ro/ee/pkg/license"
	roprometheus "github.com/samber/ro/ee/plugins/prometheus"
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
var pipeline, collector = roprometheus.Pipe7(
	roprometheus.CollectorConfig{
		Namespace: "example",
	},
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
	err := rolicense.SetLicense("ro-00eyJkIjoiZXlKMElqb2laVzUwWlhKd2NtbHpaU0lzSW1sa0lqb2liWGt0YVdRaUxDSnZkMjRpT2lKelkzSmxaV0lpTENKbGJuWWlPaUp3Y205a2RXTjBhVzl1SWl3aVpYaHdJam9pTWpBeU5TMHhNaTB6TVZRd01Eb3dNRG93TUZvaUxDSnZJanAwY25WbGZRPT0iLCJzIjoiTUlHSUFrSUEvSDVhbWlpQVA0SXNTWFhNczVhV205WjVJUHlKaExubmxpd25obG0vdTRHNkNsclVUaGdoVXo4Q3BwRit5b2VDT2hJL0orbGxiOHpZMWlIcUE5SWtXV0FDUWdFTHUxT0NKRU1LUlNhSkduQzFTQWdlYTVHUEl5cFFwb0FRSlB2ZVcxMmQvVDdrbnF0RU52WlZnWXdZNDRSNzV5anJhTXFhZ016bzNHZHllT3oyL0taZmhBPT0ifQ==")
	if err != nil {
		log.Fatalf("Failed to set license: %v", err)
	}
	prometheus.MustRegister(collector)

	time.Sleep(10 * time.Second)

	// Start HTTP server for Prometheus metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting Prometheus metrics server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	subscription := pipeline.Subscribe(
		// Log SQL and writer errors.
		ro.OnError[int](func(err error) {
			log.Println(err.Error())
		}),
	)

	// Optional, since the pipeline will complete itself, in a blocking way.
	defer subscription.Unsubscribe()

	// Keep the main goroutine alive to serve metrics
	select {}
}
