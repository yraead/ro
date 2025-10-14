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
	"context"
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/samber/ro"
	rolicense "github.com/samber/ro/ee/pkg/license"
	rootel "github.com/samber/ro/ee/plugins/otel"
	rocsv "github.com/samber/ro/plugins/encoding/csv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
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
var pipeline, collector = rootel.Pipe7(
	rootel.CollectorConfig{
		EnableLogging: false,
		EnableMetrics: false,
		EnableTracing: true,
		TraceAttributes: []attribute.KeyValue{
			attribute.String("test-mode", "tracing"),
		},
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

	ctx := context.Background()

	// Initialize OpenTelemetry trace exporter
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ee-otel-tracing"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(ctx)

	// Subscribe to the pipeline
	subscription := pipeline.Subscribe(
		ro.OnError[int](func(err error) {
			log.Println(err.Error())
		}),
	)
	defer subscription.Unsubscribe()

	log.Println("Processing completed!")
}
