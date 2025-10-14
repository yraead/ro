# Ro WebSocket Client Plugin

This plugin provides reactive WebSocket client functionality for the [Ro](https://github.com/samber/ro) reactive programming library. It allows you to:

- Create WebSocket connections as reactive subjects
- Send and receive WebSocket messages through reactive streams
- Apply reactive operators to WebSocket message processing
- Handle WebSocket lifecycle events reactively
- Build real-time applications with bidirectional communication

## Installation

```bash
go get github.com/samber/ro/plugins/websocket/client
```

## Requirements

- [Ro](https://github.com/samber/ro) reactive programming library
- [gorilla/websocket](https://github.com/gorilla/websocket) for WebSocket protocol support
- Go 1.18 or later

## Quick Start

```go
package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/websocket/client"
)

func main() {
	// Create a WebSocket subject
	ws := rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[string, string]{
		URL: "ws://localhost:8080/ws",
		Serializer: func(msg string) ([]byte, error) {
			return json.Marshal(msg)
		},
		Deserializer: func(data []byte) (string, error) {
			var msg string
			err := json.Unmarshal(data, &msg)
			return msg, err
		},
	})
	
	// Subscribe to incoming messages
	subscription := ws.Subscribe(
		ro.NewObserver(
			func(message string) {
				fmt.Printf("Received: %s\n", message)
			},
			func(err error) {
				fmt.Printf("Error: %v\n", err)
			},
			func() {
				fmt.Println("WebSocket closed")
			},
		),
	)
	defer subscription.Unsubscribe()
	
	// Send a message
	ws.Next("Hello, WebSocket!")
	
	// Keep the application running
	select {}
}
```

## API Reference

### Core Types

#### `WebsocketSubjectConfig[In any, Out any]`

Configuration for creating WebSocket subjects.

```go
type WebsocketSubjectConfig[In any, Out any] struct {
	URL          string                    // WebSocket server URL
	Headers      map[string]string        // Custom headers for connection
	Serializer   Serializer[In]          // Function to serialize outgoing messages
	Deserializer Deserializer[Out]       // Function to deserialize incoming messages
	Dialer       *websocket.Dialer        // Custom WebSocket dialer
	OutputConnector func() ro.Subject[Out] // Custom subject factory for output
}
```

#### `Serializer[T any]`

Function type for serializing messages before sending.

```go
type Serializer[T any] func(T) ([]byte, error)
```

#### `Deserializer[T any]`

Function type for deserializing received messages.

```go
type Deserializer[T any] func([]byte) (T, error)
```

### Constructors

#### `NewWebsocketSubject[In any, Out any](config WebsocketSubjectConfig[In, Out]) *WebsocketSubject[In, Out]`

Creates a new WebSocket subject that can both send and receive messages.

```go
ws := rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[string, string]{
	URL: "ws://localhost:8080/ws",
	Serializer: func(msg string) ([]byte, error) {
		return json.Marshal(msg)
	},
	Deserializer: func(data []byte) (string, error) {
		var msg string
		return msg, json.Unmarshal(data, &msg)
	},
})
```

#### `NewWebsocketObservable[Out any](config WebsocketObservableConfig[struct{}, Out]) ro.Observable[Out]`

Creates a read-only WebSocket observable for receiving messages only.

```go
observable := rowebsocket.NewWebsocketObservable[ChatMessage](rowebsocket.WebsocketObservableConfig[struct{}, ChatMessage]{
	URL: "ws://localhost:8080/chat",
	Headers: map[string]string{
		"Authorization": "Bearer your-token",
	},
	Deserializer: func(data []byte) (ChatMessage, error) {
		var msg ChatMessage
		return msg, json.Unmarshal(data, &msg)
	},
})
```

#### `NewWebsocketObserver[In any](config WebsocketObserverConfig[In]) ro.Observer[In]`

Creates a write-only WebSocket observer for sending messages only.

```go
observer := rowebsocket.NewWebsocketObserver[Command](rowebsocket.WebsocketObserverConfig[Command]{
	URL: "ws://localhost:8080/commands",
	Serializer: func(cmd Command) ([]byte, error) {
		return json.Marshal(cmd)
	},
})
```

## Usage Examples

### Basic Chat Client

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/websocket/client"
)

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Timestamp int64 `json:"timestamp"`
}

func main() {
	// Get username from command line or use default
	username := "Anonymous"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}
	
	// Create WebSocket subject
	ws := rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[ChatMessage, ChatMessage]{
		URL: "ws://localhost:8080/chat",
		Headers: map[string]string{
			"X-Username": username,
		},
		Serializer: func(msg ChatMessage) ([]byte, error) {
			return json.Marshal(msg)
		},
		Deserializer: func(data []byte) (ChatMessage, error) {
			var msg ChatMessage
			err := json.Unmarshal(data, &msg)
			return msg, err
		},
	})
	
	// Subscribe to incoming messages
	subscription := ws.Subscribe(
		ro.NewObserver(
			func(message ChatMessage) {
				fmt.Printf("[%s] %s: %s\n", 
					formatTime(message.Timestamp), 
					message.Username, 
					message.Message)
			},
			func(err error) {
				log.Printf("WebSocket error: %v", err)
			},
			func() {
				fmt.Println("Disconnected from chat server")
			},
		),
	)
	defer subscription.Unsubscribe()
	
	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Read user input in a separate goroutine
	go func() {
		var input string
		for {
			fmt.Scanln(&input)
			if input == "/quit" {
				ws.Complete()
				return
			}
			
			// Send message
			ws.Next(ChatMessage{
				Username: username,
				Message:  input,
				Timestamp: time.Now().Unix(),
			})
		}
	}()
	
	// Wait for shutdown signal
	<-c
	fmt.Println("\nShutting down...")
	ws.Complete()
}

func formatTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("15:04:05")
}
```

### Real-time Data Stream with Processing

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/websocket/client"
)

type SensorData struct {
	SensorID    string  `json:"sensor_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   int64   `json:"timestamp"`
}

type ProcessedData struct {
	SensorID    string  `json:"sensor_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Status      string  `json:"status"`
	Timestamp   int64   `json:"timestamp"`
}

func main() {
	// Create WebSocket observable for receiving sensor data
	observable := rowebsocket.NewWebsocketObservable[SensorData](rowebsocket.WebsocketObservableConfig[struct{}, SensorData]{
		URL: "ws://localhost:8080/sensors",
		Deserializer: func(data []byte) (SensorData, error) {
			var sensor SensorData
			err := json.Unmarshal(data, &sensor)
			return sensor, err
		},
	})
	
	// Process the data stream with reactive operators
	processed := ro.Pipe4(
		observable,
		// Filter out invalid readings
		ro.Filter(func(data SensorData) bool {
			return data.Temperature > -50 && data.Temperature < 100 &&
				data.Humidity >= 0 && data.Humidity <= 100
		}),
		// Transform and add status
		ro.Map(func(data SensorData) ProcessedData {
			status := "normal"
			if data.Temperature > 30 {
				status = "hot"
			} else if data.Temperature < 10 {
				status = "cold"
			}
			if data.Humidity > 70 {
				status = "humid"
			}
			
			return ProcessedData{
				SensorID:    data.SensorID,
				Temperature: data.Temperature,
				Humidity:    data.Humidity,
				Status:      status,
				Timestamp:   data.Timestamp,
			}
		}),
		// Group readings by sensor
		ro.GroupBy(func(data ProcessedData) string {
			return data.SensorID
		}),
		// Calculate moving averages for each sensor
		ro.Map(func(grouped ro.Observable[ProcessedData]) ro.Observable[ProcessedData] {
			return ro.Pipe1(
				grouped,
				ro.BufferCount[[]ProcessedData](5),
				ro.Map(func(readings []ProcessedData) ProcessedData {
					if len(readings) == 0 {
						return ProcessedData{}
					}
					
					// Calculate averages
					var tempSum, humiditySum float64
					for _, reading := range readings {
						tempSum += reading.Temperature
						humiditySum += reading.Humidity
					}
					
					last := readings[len(readings)-1]
					return ProcessedData{
						SensorID:    last.SensorID,
						Temperature: tempSum / float64(len(readings)),
						Humidity:    humiditySum / float64(len(readings)),
						Status:      "averaged",
						Timestamp:   last.Timestamp,
					}
				}),
			)
		}),
		ro.Merge[ProcessedData](),
	)
	
	// Subscribe to processed data
	subscription := processed.Subscribe(
		ro.NewObserver(
			func(data ProcessedData) {
				fmt.Printf("Sensor %s: %.1f°C, %.1f%% humidity (%s)\n",
					data.SensorID, data.Temperature, data.Humidity, data.Status)
			},
			func(err error) {
				log.Printf("Stream error: %v", err)
			},
			func() {
				fmt.Println("Data stream ended")
			},
		),
	)
	defer subscription.Unsubscribe()
	
	// Keep the application running
	select {}
}
```

### Bidirectional Command Interface

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/websocket/client"
)

type Command struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	ID      string      `json:"id"`
}

type Response struct {
	CommandID string      `json:"command_id"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Error     string      `json:"error,omitempty"`
}

func main() {
	// Create WebSocket subject for bidirectional communication
	ws := rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[Command, Response]{
		URL: "ws://localhost:8080/api",
		Headers: map[string]string{
			"Authorization": "Bearer your-api-token",
		},
		Serializer: func(cmd Command) ([]byte, error) {
			return json.Marshal(cmd)
		},
		Deserializer: func(data []byte) (Response, error) {
			var resp Response
			err := json.Unmarshal(data, &resp)
			return resp, err
		},
	})
	
	// Subscribe to responses
	subscription := ws.Subscribe(
		ro.NewObserver(
			func(response Response) {
				if response.Success {
					fmt.Printf("✓ Command %s succeeded: %v\n", response.CommandID, response.Data)
				} else {
					fmt.Printf("✗ Command %s failed: %s\n", response.CommandID, response.Error)
				}
			},
			func(err error) {
				log.Printf("WebSocket error: %v", err)
			},
			func() {
				fmt.Println("Connection closed")
			},
		),
	)
	defer subscription.Unsubscribe()
	
	// Send some example commands
	go func() {
		time.Sleep(1 * time.Second) // Wait for connection
		
		// Send ping command
		ws.Next(Command{
			Type:    "ping",
			Payload: nil,
			ID:      generateID(),
		})
		
		time.Sleep(2 * time.Second)
		
		// Send data request
		ws.Next(Command{
			Type: "get_data",
			Payload: map[string]interface{}{
				"source": "sensors",
				"limit":  10,
			},
			ID: generateID(),
		})
		
		time.Sleep(2 * time.Second)
		
		// Send configuration update
		ws.Next(Command{
			Type: "update_config",
			Payload: map[string]interface{}{
				"threshold": 25.5,
				"mode":      "auto",
			},
			ID: generateID(),
		})
	}()
	
	// Keep the application running
	select {}
}

func generateID() string {
	return fmt.Sprintf("cmd_%d", time.Now().UnixNano())
}
```

### Reconnection Logic with Backoff

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/websocket/client"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func main() {
	// Create reconnection logic
	connect := func() ro.Subject[Message] {
		return rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[Message, Message]{
			URL: "ws://localhost:8080/ws",
			Serializer: func(msg Message) ([]byte, error) {
				return json.Marshal(msg)
			},
			Deserializer: func(data []byte) (Message, error) {
				var msg Message
				err := json.Unmarshal(data, &msg)
				return msg, err
			},
		})
	}
	
	// Create retry logic with exponential backoff
	retryLogic := func(attempt int) time.Duration {
		return time.Duration(attempt*attempt) * time.Second
	}
	
	// Main connection loop
	go func() {
		for attempt := 1; attempt <= 10; attempt++ {
			fmt.Printf("Attempting to connect (attempt %d)...\n", attempt)
			
			ws := connect()
			
			// Subscribe to messages
			subscription := ws.Subscribe(
				ro.NewObserver(
					func(msg Message) {
						fmt.Printf("Received: %+v\n", msg)
						attempt = 1 // Reset attempt counter on successful message
					},
					func(err error) {
						log.Printf("Connection error: %v", err)
					},
					func() {
						fmt.Println("Connection closed, will retry...")
					},
				),
			)
			
			// Send initial message
			go func() {
				time.Sleep(1 * time.Second)
				ws.Next(Message{
					Type: "hello",
					Data: map[string]string{
						"client": "reconnecting-client",
					},
				})
			}()
			
			// Wait for disconnection
			<-ws.IsCompleted()
			subscription.Unsubscribe()
			
			// Wait before retrying
			if attempt < 10 {
				waitTime := retryLogic(attempt)
				fmt.Printf("Waiting %v before reconnecting...\n", waitTime)
				time.Sleep(waitTime)
			}
		}
		
		fmt.Println("Max retry attempts reached, giving up")
	}()
	
	// Keep the application running
	select {}
}
```

### Multi-Client Chat Room Manager

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	
	"github.com/samber/ro"
	"github.com/samber/ro/plugins/websocket/client"
)

type ChatMessage struct {
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type ChatRoom struct {
	ID       string
	Messages []ChatMessage
	Clients  map[string]*WebSocketClient
	mu       sync.RWMutex
}

type WebSocketClient struct {
	Subject   *rowebsocket.WebsocketSubject[ChatMessage, ChatMessage]
	Username  string
	UserID    string
	LastSeen  time.Time
}

func main() {
	// Create multiple chat rooms
	rooms := map[string]*ChatRoom{
		"general": {ID: "general", Clients: make(map[string]*WebSocketClient)},
		"random":  {ID: "random", Clients: make(map[string]*WebSocketClient)},
		"tech":    {ID: "tech", Clients: make(map[string]*WebSocketClient)},
	}
	
	// Simulate multiple clients joining different rooms
	for i := 1; i <= 5; i++ {
		for roomID := range rooms {
			go func(clientID, roomID string) {
				client := createClient(clientID, roomID)
				rooms[roomID].mu.Lock()
				rooms[roomID].Clients[clientID] = client
				rooms[roomID].mu.Unlock()
				
				// Handle messages
				subscription := client.Subject.Subscribe(
					ro.NewObserver(
						func(msg ChatMessage) {
							// Store message in room
							rooms[roomID].mu.Lock()
							rooms[roomID].Messages = append(rooms[roomID].Messages, msg)
							fmt.Printf("[%s] %s: %s\n", roomID, msg.Username, msg.Message)
							rooms[roomID].mu.Unlock()
							
							// Update last seen
							client.LastSeen = time.Now()
						},
						func(err error) {
							log.Printf("Client %s in room %s error: %v", clientID, roomID, err)
						},
						func() {
							fmt.Printf("Client %s left room %s\n", clientID, roomID)
							
							// Remove client from room
							rooms[roomID].mu.Lock()
							delete(rooms[roomID].Clients, clientID)
							rooms[roomID].mu.Unlock()
						},
					),
				)
				defer subscription.Unsubscribe()
				
				// Send periodic messages
				ticker := time.NewTicker(10 * time.Second)
				defer ticker.Stop()
				
				for {
					select {
					case <-ticker.C:
						// Send a message to the room
						client.Subject.Next(ChatMessage{
							RoomID:    roomID,
							UserID:    clientID,
							Username:  client.Username,
							Message:   fmt.Sprintf("Hello from %s at %s!", client.Username, time.Now().Format("15:04:05")),
							Timestamp: time.Now().Unix(),
						})
					}
				}
			}(fmt.Sprintf("client_%d", i), roomID)
		}
	}
	
	// Monitor rooms
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				for roomID, room := range rooms {
					room.mu.RLock()
					fmt.Printf("Room %s: %d clients, %d messages\n", 
						roomID, len(room.Clients), len(room.Messages))
					
					// Clean up inactive clients
					for clientID, client := range room.Clients {
						if time.Since(client.LastSeen) > 2*time.Minute {
							fmt.Printf("Removing inactive client %s from room %s\n", clientID, roomID)
							client.Subject.Complete()
						}
					}
					room.mu.RUnlock()
				}
			}
		}
	}()
	
	// Keep the application running
	select {}
}

func createClient(userID, roomID string) *WebSocketClient {
	return &WebSocketClient{
		Subject: rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[ChatMessage, ChatMessage]{
			URL: fmt.Sprintf("ws://localhost:8080/room/%s", roomID),
			Headers: map[string]string{
				"X-User-ID": userID,
				"X-Room-ID": roomID,
			},
			Serializer: func(msg ChatMessage) ([]byte, error) {
				return json.Marshal(msg)
			},
			Deserializer: func(data []byte) (ChatMessage, error) {
				var msg ChatMessage
				err := json.Unmarshal(data, &msg)
				return msg, err
			},
		}),
		UserID:   userID,
		Username: fmt.Sprintf("User_%s", userID[len(userID)-1:]),
		LastSeen: time.Now(),
	}
}
```

## Configuration Options

### Custom Dialer

```go
// Create a custom dialer with timeout
dialer := &websocket.Dialer{
	HandshakeTimeout: 10 * time.Second,
	NetDial: func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, 5*time.Second)
	},
}

ws := rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[string, string]{
	URL:    "wss://echo.websocket.org",
	Dialer: dialer,
	Serializer: func(msg string) ([]byte, error) {
		return []byte(msg), nil
	},
	Deserializer: func(data []byte) (string, error) {
		return string(data), nil
	},
})
```

### Custom Output Subject

```go
// Use a ReplaySubject to replay last 5 messages for new subscribers
ws := rowebsocket.NewWebsocketSubject(rowebsocket.WebsocketSubjectConfig[string, string]{
	URL: "ws://localhost:8080/ws",
	Serializer: func(msg string) ([]byte, error) {
		return json.Marshal(msg)
	},
	Deserializer: func(data []byte) (string, error) {
		var msg string
		return msg, json.Unmarshal(data, &msg)
	},
	OutputConnector: func() ro.Subject[string] {
		return ro.NewReplaySubject[string](5)
	},
})
```

## Best Practices

1. **Connection Management**: Always handle connection lifecycle events properly
2. **Error Handling**: Implement proper error handling and reconnection logic
3. **Message Validation**: Validate incoming messages before processing
4. **Backpressure**: Use reactive operators to handle high-frequency messages
5. **Resource Cleanup**: Always unsubscribe and complete connections properly
6. **Security**: Use WSS (WebSocket Secure) for production applications
7. **Authentication**: Include authentication tokens in headers or initial messages
8. **Heartbeat**: Implement ping/pong mechanism for connection health monitoring

## Performance Considerations

1. **Message Buffering**: Use appropriate buffering strategies for high-frequency data
2. **Serialization**: Choose efficient serialization formats (JSON, Protobuf, MessagePack)
3. **Connection Pooling**: Reuse connections when possible
4. **Memory Management**: Be mindful of message accumulation in subjects
5. **Concurrency**: Leverage Go's concurrency features for multiple connections

## Security Considerations

1. **Use WSS**: Always use secure WebSocket connections in production
2. **Authentication**: Implement proper authentication and authorization
3. **Input Validation**: Validate all incoming messages
4. **Rate Limiting**: Implement rate limiting to prevent abuse
5. **Origin Validation**: Validate the Origin header on the server side

## License

Apache 2.0 - See [LICENSE](https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md) for details.