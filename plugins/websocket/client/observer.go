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

type WebsocketObserverConfig[In any] struct {
	URL        string
	Headers    map[string]string
	Serializer Serializer[In]
	Dialer     *websocket.Dialer
}

// NewWebsocketObserver creates a websocket observer that can send messages to a websocket endpoint.
func NewWebsocketObserver[In any](config WebsocketObserverConfig[In]) ro.Observer[In] {
	return NewWebsocketSubject(WebsocketSubjectConfig[In, struct{}]{
		URL:             config.URL,
		Headers:         config.Headers,
		Serializer:      config.Serializer,
		Deserializer:    func([]byte) (struct{}, error) { return struct{}{}, nil },
		Dialer:          config.Dialer,
		OutputConnector: ro.NewPublishSubject[struct{}],
	})
}
