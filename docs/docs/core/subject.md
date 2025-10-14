---
title: Subject
description: Learn about Subject - the bridge between Observable and Observer in samber/ro
sidebar_position: 5
---

# 🎯 Subject

A **Subject** is a special type that extends both `Observable` and `Observer`, acting as a bridge or proxy that can multicast values to multiple observers. Subjects are the key to implementing hot observables and enabling event broadcasting patterns.

## What is a Subject?

A `Subject` is:
- **An Observable**: Can be subscribed to like any other `Observable`
- **An Observer**: Can receive values through `Next`, `Error`, and `Complete` methods
- **A multicaster**: Can broadcast values to multiple subscribers
- **Hot by nature**: Values are shared among all subscribers

The Subject interface combines both Observable and Observer:

```go
type Subject[T any] interface {
    Observable[T]
    // Subscribe(destination Observer[T]) Subscription
    // SubscribeWithContext(ctx context.Context, destination Observer[T]) Subscription

    Observer[T]
    // Next(value T)
    // NextWithContext(ctx context.Context, value T)
    // Error(err error)
    // ErrorWithContext(ctx context.Context, err error)
    // Complete()
    // CompleteWithContext(ctx context.Context)
    // IsClosed() bool
    // HasThrown() bool
    // IsCompleted() bool


    HasObserver() bool
    CountObservers() int
}
```

## Subject Types

`samber/ro` provides four types of subjects, each with different behaviors:

### 1. PublishSubject

PublishSubject emits only the values that were sent after the subscription. This is the most basic subject type, perfect for simple event broadcasting.

```go
// Create a PublishSubject
subject := ro.NewPublishSubject[int]()

// Subscriber 1 - gets all future values
subject.Subscribe(ro.OnNext(func(n int) {
    fmt.Println("Subscriber 1:", n)
}))

// Emit values
subject.Next(1)
subject.Next(2)

// Subscriber 2 - gets only values sent after subscription
subject.Subscribe(ro.OnNext(func(n int) {
    fmt.Println("Subscriber 2:", n)
}))

subject.Next(3)
subject.Next(4)
subject.Complete()

// Output:
// Subscriber 1: 1
// Subscriber 1: 2
// Subscriber 1: 3
// Subscriber 2: 3
// Subscriber 1: 4
// Subscriber 2: 4
// Subscriber 1: Completed
// Subscriber 2: Completed
```

**Use cases for PublishSubject:**
- Event broadcasting systems
- Real-time notifications
- Chat applications
- Live data streams where only new values matter

### 2. BehaviorSubject

BehaviorSubject emits the last value and all subsequent values to new subscribers. It requires an initial value and is ideal for state management scenarios.

```go
// Create a BehaviorSubject with initial value
subject := ro.NewBehaviorSubject(42)

// Subscriber 1 - immediately gets the current value
subject.Subscribe(ro.NewObserver(
    func(value int) {
        fmt.Println("Subscriber 1 received:", value)
    },
    func(err error) {
        fmt.Println("Subscriber 1 error:", err)
    },
    func() {
        fmt.Println("Subscriber 1 completed")
    },
))

// Emit new values
subject.Next(100)
subject.Next(200)

// Subscriber 2 - immediately gets the latest value
subject.Subscribe(ro.NewObserver(
    func(value int) {
        fmt.Println("Subscriber 2 received:", value)
    },
    func(err error) {
        fmt.Println("Subscriber 2 error:", err)
    },
    func() {
        fmt.Println("Subscriber 2 completed")
    },
))

subject.Next(300)
subject.Complete()

// Output:
// Subscriber 1 received: 42
// Subscriber 1 received: 100
// Subscriber 1 received: 200
// Subscriber 2 received: 200
// Subscriber 1 received: 300
// Subscriber 2 received: 300
// Subscriber 1 completed
// Subscriber 2 completed
```

**Use cases for BehaviorSubject:**
- State management
- Configuration settings
- Current user information
- Any scenario where new subscribers need the current state

### 3. ReplaySubject

ReplaySubject emits a specified number of previous values and all subsequent values to new subscribers. Use it when you need to provide context or history to new subscribers.

```go
// Create a ReplaySubject with buffer size of 3
subject := ro.NewReplaySubject[string](3)

// Emit values
subject.Next("first")
subject.Next("second")
subject.Next("third")
subject.Next("fourth")

// Subscriber 1 - gets last 3 values
subject.Subscribe(ro.OnNext(func(s string) {
    fmt.Println("Subscriber 1:", s)
}))

subject.Next("fifth")

// Subscriber 2 - gets last 3 values
subject.Subscribe(ro.OnNext(func(s string) {
    fmt.Println("Subscriber 2:", s)
}))

subject.Complete()

// Output:
// Subscriber 1: second
// Subscriber 1: third
// Subscriber 1: fourth
// Subscriber 1: fifth
// Subscriber 2: third
// Subscriber 2: fourth
// Subscriber 2: fifth
// Subscriber 1: Completed
// Subscriber 2: Completed
```

**Use cases for ReplaySubject:**
- Chat history
- Stock price updates
- Game replay systems
- Notification history
- Any scenario where new users need context

### 4. AsyncSubject

AsyncSubject emits only the last value and only when the sequence completes. This is perfect for async operations that produce a single final result.

```go
// Create an AsyncSubject
subject := ro.NewAsyncSubject[float64]()

// Subscriber 1 - will receive nothing until completion
subject.Subscribe(ro.NewObserver(
    func(value float64) {
        fmt.Println("Subscriber 1 received:", value)
    },
    func(err error) {
        fmt.Println("Subscriber 1 error:", err)
    },
    func() {
        fmt.Println("Subscriber 1 completed")
    },
))

// Emit multiple values
subject.Next(1.0)
subject.Next(2.0)
subject.Next(3.0)

// Subscriber 2 - will also receive only the final value
subject.Subscribe(ro.NewObserver(
    func(value float64) {
        fmt.Println("Subscriber 2 received:", value)
    },
    func(err error) {
        fmt.Println("Subscriber 2 error:", err)
    },
    func() {
        fmt.Println("Subscriber 2 completed")
    },
))

// Complete to trigger emission
subject.Complete()

// Output:
// Subscriber 1 received: 3.0
// Subscriber 1 completed
// Subscriber 2 received: 3.0
// Subscriber 2 completed
```

**Use cases for AsyncSubject:**
- Asynchronous operations that return a single result
- HTTP requests
- Database queries
- File operations
- Any computation that produces one final result

## Subject Lifecycle Management

### Checking Subject State

Monitor subject activity to debug issues or implement conditional logic based on subscriber presence.

```go
subject := ro.NewPublishSubject[int]()

// Check if there are observers
fmt.Println("Has observers:", subject.HasObserver()) // false
fmt.Println("Observer count:", subject.CountObservers()) // 0

// Add subscribers
subject.Subscribe(ro.OnNext(func(n int) { fmt.Println(n) }))

// Check state again
fmt.Println("Has observers:", subject.HasObserver()) // true
fmt.Println("Observer count:", subject.CountObservers()) // 1
```

### Unsubscribing from Subjects

Manage individual subscriptions to control which observers receive values. This is useful for fine-grained resource management.

```go
subject := ro.NewPublishSubject[string]()

// Subscribe and keep the subscription
subscription1 := subject.Subscribe(ro.OnNext(func(s string) {
    fmt.Println("Subscriber 1:", s)
}))

subscription2 := subject.Subscribe(ro.OnNext(func(s string) {
    fmt.Println("Subscriber 2:", s)
}))

subject.Next("hello")
subject.Next("world")

// Unsubscribe one observer
subscription1.Unsubscribe()

subject.Next("again")

// Output:
// Subscriber 1: hello
// Subscriber 2: hello
// Subscriber 1: world
// Subscriber 2: world
// Subscriber 2: again
```

## Subject vs Observable

```go
// Observable (cold) - each subscription gets independent values
observable := ro.Just(1, 2, 3)

// first Subscription
observable.Subscribe(ro.OnNext(func(n int) {
    fmt.Println("Observable 1:", n)
}))

// second Subscription
observable.Subscribe(ro.OnNext(func(n int) {
    fmt.Println("Observable 2:", n)
}))

// Output:
// Observable 1: 1
// Observable 1: 2
// Observable 1: 3
// Observable 2: 1
// Observable 2: 2
// Observable 2: 3

// Subject (hot) - subscriptions share the same values concurrently
subject := ro.NewPublishSubject[int]()

subject.Subscribe(ro.OnNext(func(n int) {
    fmt.Println("Subject 1:", n)
}))

subject.Subscribe(ro.OnNext(func(n int) {
    fmt.Println("Subject 2:", n)
}))

subject.Next(1)
subject.Next(2)
subject.Next(3)

// Output:
// Subject 1: 1
// Subject 2: 1
// Subject 1: 2
// Subject 2: 2
// Subject 1: 3
// Subject 2: 3
```

## Advanced Subject Patterns

### State Management with BehaviorSubject

```go
// Simple state management system
type AppState struct {
    User    string
    Counter int
}

// Create state subject
stateSubject := ro.NewBehaviorSubject(AppState{
    User:    "guest",
    Counter: 0,
})

// Component that listens to state changes
stateSubject.Subscribe(ro.OnNext(func(state AppState) {
    fmt.Printf("UI Update: User=%s, Counter=%d\n", state.User, state.Counter)
}))

// Update state
updateState := func(user string, counter int) {
    current := stateSubject.LastValue() // Requires implementation
    stateSubject.Next(AppState{
        User:    user,
        Counter: counter,
    })
}

// Simulate state changes
updateState("alice", 1)
updateState("alice", 2)
updateState("bob", 1)

// Output:
// UI Update: User=guest, Counter=0
// UI Update: User=alice, Counter=1
// UI Update: User=alice, Counter=2
// UI Update: User=bob, Counter=1
```

### Event Bus with PublishSubject

```go
// Simple event bus
type EventBus struct {
    events map[string]ro.Subject[string]
    mutex  sync.Mutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        events: make(map[string]ro.Subject[string]),
    }
}

func (eb *EventBus) Publish(eventType, data string) {
    eb.mutex.Lock()
    defer eb.mutex.Unlock()

    if subject, exists := eb.events[eventType]; exists {
        subject.Next(data)
    }
}

func (eb *EventBus) Subscribe(eventType string, observer ro.Observer[string]) ro.Subscription {
    eb.mutex.Lock()
    defer eb.mutex.Unlock()

    if _, exists := eb.events[eventType]; !exists {
        eb.events[eventType] = ro.NewPublishSubject[string]()
    }

    return eb.events[eventType].Subscribe(observer)
}

// Usage
eventBus := NewEventBus()

// Subscribe to events
userSub := eventBus.Subscribe("user.login", ro.OnNext(func(data string) {
    fmt.Println("User logged in:", data)
}))

orderSub := eventBus.Subscribe("order.created", ro.OnNext(func(data string) {
    fmt.Println("Order created:", data)
}))

// Publish events
eventBus.Publish("user.login", "alice")
eventBus.Publish("order.created", "order-123")
eventBus.Publish("user.login", "bob")

// Output:
// User logged in: alice
// Order created: order-123
// User logged in: bob
```

### Chat System with ReplaySubject

```go
// Chat room with history
type ChatRoom struct {
    messages ro.Subject[string]
    users    map[string]ro.Subject[string]
    mutex    sync.Mutex
}

func NewChatRoom() *ChatRoom {
    return &ChatRoom{
        messages: ro.NewReplaySubject[string](100), // Keep last 100 messages
        users:    make(map[string]ro.Subject[string]),
    }
}

func (cr *ChatRoom) Join(username string) ro.Subject[string] {
    cr.mutex.Lock()
    defer cr.mutex.Unlock()

    if _, exists := cr.users[username]; !exists {
        cr.users[username] = ro.NewPublishSubject[string]()
        cr.messages.Next(fmt.Sprintf("%s joined the chat", username))
    }

    return cr.users[username]
}

func (cr *ChatRoom) SendMessage(username, message string) {
    fullMessage := fmt.Sprintf("%s: %s", username, message)
    cr.messages.Next(fullMessage)
}

func (cr *ChatRoom) GetMessages() ro.Observable[string] {
    return cr.messages
}

// Usage
chatRoom := NewChatRoom()

// User joins and gets message history
user1Messages := chatRoom.Join("alice")
user1Messages.Subscribe(ro.OnNext(func(msg string) {
    fmt.Println("Alice sees:", msg)
}))

// Send some messages
chatRoom.SendMessage("alice", "Hello!")
chatRoom.SendMessage("alice", "Is anyone there?")

// Another user joins later and gets history
user2Messages := chatRoom.Join("bob")
user2Messages.Subscribe(ro.OnNext(func(msg string) {
    fmt.Println("Bob sees:", msg)
}))

chatRoom.SendMessage("bob", "Hi Alice!")
```

## Subject Best Practices

### 1. Choose the Right Subject Type

```go
// Good: Use appropriate subject for the use case
currentState := ro.NewBehaviorSubject(initialState)         // For current state
eventStream := ro.NewPublishSubject[Event]()                // For new events only
messageHistory := ro.NewReplaySubject[Message](1000)        // For history/caching
finalResult := ro.NewAsyncSubject[Result]()                 // For single async result
```

### 2. Manage Subject Lifecycle

```go
// Good: Clean up subjects when done
func processEvents() {
    subject := ro.NewPublishSubject[Event]()
    defer subject.Complete()

    subscription := subject.Subscribe(observer)
    defer subscription.Unsubscribe()

    // Process events...
}

// Risky: Subject might outlive its usefulness
var globalSubject ro.Subject[Event] // Could leak memory
```

### 3. Handle Errors Gracefully

```go
// Good: Handle errors in subject streams
subject := ro.NewPublishSubject[Data]()

subject.Subscribe(ro.NewObserver(
    func(data Data) { /* process data */ },
    func(err error) {
        log.Printf("Subject error: %v", err)
        // Implement recovery logic
    },
    func() { /* handle completion */ },
))

// Don't let errors go unhandled
subject.Error(fmt.Errorf("processing failed"))
```

### 4. Avoid Memory Leaks

```go
// Good: Use bounded replay subjects
history := ro.NewReplaySubject[Event](1_000) // Reasonable buffer size

// Risky: Unbounded buffer could exhaust memory
unbounded := ro.NewReplaySubject[Event](1_000_000) // Too large
```

Subjects provide a powerful, reactive way to implement event-driven systems with built-in multicasting, lifecycle management, and composition with other reactive operators. They are essential for building complex, real-time applications in Go.
