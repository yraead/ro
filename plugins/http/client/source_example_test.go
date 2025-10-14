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

package rohttpclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/samber/ro"
)

// customHTTPObserver formats HTTP responses in the expected test output format
func customHTTPObserver() ro.Observer[*http.Response] {
	return ro.NewObserverWithContext(
		func(ctx context.Context, resp *http.Response) {
			fmt.Printf("Next: &http.Response{Status: \"%s\", StatusCode: %d, ...}\n", resp.Status, resp.StatusCode)
		},
		func(ctx context.Context, err error) {
			fmt.Printf("Error: %s\n", err.Error())
		},
		func(ctx context.Context) {
			fmt.Printf("Completed\n")
		},
	)
}

func ExampleHTTPRequest_basic() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, World!"}`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request
	observable := HTTPRequest(req, nil)

	subscription := observable.Subscribe(customHTTPObserver())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
	// Completed
}

func ExampleHTTPRequest_withHeaders() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the User-Agent header
		userAgent := r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User-Agent: %s", userAgent)))
	}))
	defer server.Close()

	// Create HTTP request with custom headers
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("User-Agent", "ro-http-client/1.0")
	req.Header.Set("Accept", "application/json")

	// Send HTTP request
	observable := HTTPRequest(req, nil)

	subscription := observable.Subscribe(customHTTPObserver())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
	// Completed
}

func ExampleHTTPRequest_withTimeout() {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate slow response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delayed response"))
	}))
	defer server.Close()

	// Create HTTP request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequest("GET", server.URL, nil)
	req = req.WithContext(ctx)

	// Send HTTP request
	observable := HTTPRequest(req, nil)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(resp *http.Response) {
				// Handle successful response
			},
			func(err error) {
				// Handle timeout or other errors
				fmt.Printf("Error due to timeout\n")
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output: Error due to timeout
}

func ExampleHTTPRequest_withCustomClient() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Custom client response"))
	}))
	defer server.Close()

	// Create custom HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request with custom client
	observable := HTTPRequest(req, client)

	subscription := observable.Subscribe(customHTTPObserver())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
	// Completed
}

func ExampleHTTPRequest_processingResponse() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "data": "test"}`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request and process response
	observable := ro.Pipe1(
		HTTPRequest(req, nil),
		ro.Map(func(resp *http.Response) string {
			// Read response body (in real code, you'd handle errors)
			// For this example, we'll just return status
			return fmt.Sprintf("Status: %s", resp.Status)
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: Status: 200 OK
	// Completed
}

func ExampleHTTPRequest_multipleRequests() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Response"))
	}))
	defer server.Close()

	// Create multiple requests
	requests := []*http.Request{
		func() *http.Request { req, _ := http.NewRequest("GET", server.URL, nil); return req }(),
		func() *http.Request { req, _ := http.NewRequest("GET", server.URL, nil); return req }(),
		func() *http.Request { req, _ := http.NewRequest("GET", server.URL, nil); return req }(),
	}

	// Process multiple requests
	observable := ro.Pipe1(
		ro.FromSlice(requests),
		ro.MergeMap(func(req *http.Request) ro.Observable[*http.Response] {
			return HTTPRequest(req, nil)
		}),
	)

	subscription := observable.Subscribe(customHTTPObserver())
	defer subscription.Unsubscribe()

	// Wait for async operations to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
	// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
	// Next: &http.Response{Status: "200 OK", StatusCode: 200, ...}
	// Completed
}

func ExampleHTTPRequest_errorHandling() {
	// Create a test server that returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request
	observable := HTTPRequest(req, nil)

	subscription := observable.Subscribe(customHTTPObserver())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: &http.Response{Status: "500 Internal Server Error", StatusCode: 500, ...}
	// Completed
}

// customJSONObserver formats JSON responses in the expected test output format
func customJSONObserver[T any]() ro.Observer[T] {
	return ro.NewObserverWithContext(
		func(ctx context.Context, value T) {
			fmt.Printf("Next: %+v\n", value)
		},
		func(ctx context.Context, err error) {
			fmt.Printf("Error: %s\n", err.Error())
		},
		func(ctx context.Context) {
			fmt.Printf("Completed\n")
		},
	)
}

func ExampleHTTPRequestJSON_basicString() {
	// Create a test server that returns JSON string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"Hello, World!"`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request and parse JSON response
	observable := HTTPRequestJSON[string](req, nil)

	subscription := observable.Subscribe(customJSONObserver[string]())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: Hello, World!
	// Completed
}

func ExampleHTTPRequestJSON_struct() {
	// Define a struct for the JSON response
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Create a test server that returns JSON object
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "name": "John Doe", "email": "john@example.com"}`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request and parse JSON response into struct
	observable := HTTPRequestJSON[User](req, nil)

	subscription := observable.Subscribe(customJSONObserver[User]())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: {ID:1 Name:John Doe Email:john@example.com}
	// Completed
}

func ExampleHTTPRequestJSON_slice() {
	// Create a test server that returns JSON array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`["apple", "banana", "cherry"]`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request and parse JSON response into slice
	observable := HTTPRequestJSON[[]string](req, nil)

	subscription := observable.Subscribe(customJSONObserver[[]string]())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: [apple banana cherry]
	// Completed
}

func ExampleHTTPRequestJSON_withErrorHandling() {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json content`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request and parse JSON response
	observable := HTTPRequestJSON[string](req, nil)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(value string) {
				fmt.Printf("Next: %s\n", value)
			},
			func(err error) {
				fmt.Printf("JSON parsing error occurred\n")
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output: JSON parsing error occurred
}

func ExampleHTTPRequestJSON_withCustomClient() {
	// Define a response struct
	type APIResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "Custom client used"}`))
	}))
	defer server.Close()

	// Create custom HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Send HTTP request with custom client and parse JSON
	observable := HTTPRequestJSON[APIResponse](req, client)

	subscription := observable.Subscribe(customJSONObserver[APIResponse]())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: {Status:success Message:Custom client used}
	// Completed
}

func ExampleHTTPRequestJSON_processingPipeline() {
	// Define response structs
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 42, "name": "Alice"}`))
	}))
	defer server.Close()

	// Create HTTP request
	req, _ := http.NewRequest("GET", server.URL, nil)

	// Create processing pipeline
	observable := ro.Pipe1(
		HTTPRequestJSON[User](req, nil),
		ro.Map(func(user User) string {
			return fmt.Sprintf("User %s has ID %d", user.Name, user.ID)
		}),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Wait for async operation to complete
	time.Sleep(50 * time.Millisecond)

	// Output:
	// Next: User Alice has ID 42
	// Completed
}
