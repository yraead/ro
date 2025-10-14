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

package rowebsocketclient

import (
	"github.com/gorilla/websocket"
	"github.com/samber/ro"
)

type WebsocketObservableConfig[In any, Out any] struct {
	URL          string
	Headers      map[string]string
	Deserializer Deserializer[Out]
	Dialer       *websocket.Dialer

	// Connector is a function that returns a Subject[Out].
	// This is useful when you want to use a different Subject implementation.
	// For example, you could use a ReplaySubject to replay the last N messages.
	OutputConnector func() ro.Subject[Out]
	// ResetOnError    bool
	// ResetOnComplete bool
}

// NewWebsocketObservable creates a websocket observable that receives messages from a websocket endpoint.
func NewWebsocketObservable[Out any](config WebsocketObservableConfig[struct{}, Out]) ro.Observable[Out] {
	return NewWebsocketSubject(WebsocketSubjectConfig[struct{}, Out]{
		URL:             config.URL,
		Headers:         config.Headers,
		Serializer:      func(value struct{}) ([]byte, error) { return []byte{}, nil },
		Deserializer:    config.Deserializer,
		Dialer:          config.Dialer,
		OutputConnector: config.OutputConnector,
	})
}
