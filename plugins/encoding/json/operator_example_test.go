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


package rojson

import (
	"fmt"

	"github.com/samber/ro"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ExampleMarshal() {
	// Marshal structs to JSON
	observable := ro.Pipe1(
		ro.Just(
			User{ID: 1, Name: "Alice", Age: 30},
			User{ID: 2, Name: "Bob", Age: 25},
			User{ID: 3, Name: "Charlie", Age: 35},
		),
		Marshal[User](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(data []byte) {
				fmt.Printf("Next: %s\n", string(data))
			},
			func(err error) {
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output:
	// Next: {"id":1,"name":"Alice","age":30}
	// Next: {"id":2,"name":"Bob","age":25}
	// Next: {"id":3,"name":"Charlie","age":35}
	// Completed
}

func ExampleMarshal_map() {
	// Marshal maps to JSON
	observable := ro.Pipe1(
		ro.Just(
			map[string]interface{}{"name": "Alice", "age": 30},
			map[string]interface{}{"name": "Bob", "age": 25},
			map[string]interface{}{"name": "Charlie", "age": 35},
		),
		Marshal[map[string]interface{}](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(data []byte) {
				fmt.Printf("Next: %s\n", string(data))
			},
			func(err error) {
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output:
	// Next: {"age":30,"name":"Alice"}
	// Next: {"age":25,"name":"Bob"}
	// Next: {"age":35,"name":"Charlie"}
	// Completed
}

func ExampleUnmarshal() {
	// Unmarshal JSON to structs
	observable := ro.Pipe1(
		ro.Just(
			[]byte(`{"id":1,"name":"Alice","age":30}`),
			[]byte(`{"id":2,"name":"Bob","age":25}`),
			[]byte(`{"id":3,"name":"Charlie","age":35}`),
		),
		Unmarshal[User](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(user User) {
				fmt.Printf("Next: {ID:%d Name:%s Age:%d}\n", user.ID, user.Name, user.Age)
			},
			func(err error) {
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output:
	// Next: {ID:1 Name:Alice Age:30}
	// Next: {ID:2 Name:Bob Age:25}
	// Next: {ID:3 Name:Charlie Age:35}
	// Completed
}

func ExampleUnmarshal_map() {
	// Unmarshal JSON to maps
	observable := ro.Pipe1(
		ro.Just(
			[]byte(`{"name":"Alice","age":30}`),
			[]byte(`{"name":"Bob","age":25}`),
			[]byte(`{"name":"Charlie","age":35}`),
		),
		Unmarshal[map[string]interface{}](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(data map[string]interface{}) {
				fmt.Printf("Next: map[age:%v name:%s]\n", data["age"], data["name"])
			},
			func(err error) {
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output:
	// Next: map[age:30 name:Alice]
	// Next: map[age:25 name:Bob]
	// Next: map[age:35 name:Charlie]
	// Completed
}

func ExampleMarshal_withError() {
	// Marshal with potential errors (e.g., channels cannot be marshaled)
	type InvalidStruct struct {
		Data chan int `json:"data"`
	}

	invalid := InvalidStruct{
		Data: make(chan int),
	}

	observable := ro.Pipe1(
		ro.Just(invalid),
		Marshal[InvalidStruct](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(data []byte) {
				// Handle successful marshaling
				fmt.Printf("Next: %s\n", string(data))
			},
			func(err error) {
				// Handle marshaling error
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				// Handle completion
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output: Error: json: unsupported type: chan int
}

func ExampleUnmarshal_withError() {
	// Unmarshal with invalid JSON
	observable := ro.Pipe1(
		ro.Just(
			[]byte(`{"id":1,"name":"Alice","age":30}`),   // Valid JSON
			[]byte(`{"id":2,"name":"Bob",`),              // Invalid JSON (truncated)
			[]byte(`{"id":3,"name":"Charlie","age":35}`), // Valid JSON
		),
		Unmarshal[User](),
	)

	subscription := observable.Subscribe(
		ro.NewObserver(
			func(user User) {
				// Handle successful unmarshaling
				fmt.Printf("Next: {ID:%d Name:%s Age:%d}\n", user.ID, user.Name, user.Age)
			},
			func(err error) {
				// Handle unmarshaling error
				fmt.Printf("Error: %s\n", err.Error())
			},
			func() {
				// Handle completion
				fmt.Printf("Completed\n")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Output:
	// Next: {ID:1 Name:Alice Age:30}
	// Error: unexpected end of JSON input
}
