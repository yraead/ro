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
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samber/ro"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	fmt.Println("New WS with roomID:", roomID)

	// upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	sub := ro.NewSubscription(nil)

	// close connection and unsubscribe from all streams
	onClose := func() {
		conn.Close()
		sub.Unsubscribe()
	}

	// websocket->redis (downstream)
	sub.AddUnsubscribable(
		bridgeInstance.Subscribe(
			roomID,
			ro.NewObserver(
				func(msg string) {
					if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
						log.Println(err)
						onClose()
					}
				},
				func(err error) {
					log.Println(err)
					onClose()
				},
				onClose,
			),
		),
	)

	// redis->websocket (upstream)
	go func() {
		defer onClose()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// log.Printf("received: %s", msg)
			bridgeInstance.Publish(roomID, string(msg))
		}
	}()
}
