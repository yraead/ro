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


package rosentry

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

// StdoutTransport is a custom Sentry transport that writes events to stdout
// instead of sending them to the remote Sentry platform. This is useful for testing.
type StdoutTransport struct{}

// Configure implements the sentry.Transport interface
func (t *StdoutTransport) Configure(options sentry.ClientOptions) {
	// No configuration needed for stdout transport
}

// SendEvent implements the sentry.Transport interface
func (t *StdoutTransport) SendEvent(event *sentry.Event) {
	// Format the event for stdout output
	output := fmt.Sprintf("Sentry event: level=%s message=\"%s\"", event.Level, event.Message)

	// Add extra fields if present
	if len(event.Extra) > 0 {
		extraStr := fmt.Sprintf(" extra=%v", event.Extra)
		output += extraStr
	}

	// Add exception info if present
	if len(event.Exception) > 0 {
		output += " exception=error"
	}

	fmt.Fprintln(os.Stdout, output)
}

// Flush implements the sentry.Transport interface
func (t *StdoutTransport) Flush(timeout time.Duration) bool {
	return true
}

// createStdoutHub creates a Sentry hub that outputs to stdout instead of remote platform
// This is useful for testing purposes when you want to see Sentry events locally
// without sending them to the actual Sentry service.
func createStdoutHub() *sentry.Hub {
	client, _ := sentry.NewClient(sentry.ClientOptions{
		Transport: &StdoutTransport{},
		// Disable DSN to prevent remote sending
		Dsn: "",
	})

	hub := sentry.NewHub(client, sentry.NewScope())
	return hub
}
