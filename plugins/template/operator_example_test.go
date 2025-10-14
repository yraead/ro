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


package rotemplate

import (
	"github.com/samber/ro"
)

type User struct {
	Name string
	Age  int
	City string
}

func ExampleTextTemplate() {
	// Process data with text template
	observable := ro.Pipe1(
		ro.Just(
			User{Name: "Alice", Age: 30, City: "New York"},
			User{Name: "Bob", Age: 25, City: "Los Angeles"},
			User{Name: "Charlie", Age: 35, City: "Chicago"},
		),
		TextTemplate[User]("Hello {{.Name}}, you are {{.Age}} years old from {{.City}}."),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Hello Alice, you are 30 years old from New York.
	// Next: Hello Bob, you are 25 years old from Los Angeles.
	// Next: Hello Charlie, you are 35 years old from Chicago.
	// Completed
}

func ExampleHTMLTemplate() {
	// Process data with HTML template
	observable := ro.Pipe1(
		ro.Just(
			User{Name: "Alice", Age: 30, City: "New York"},
			User{Name: "Bob", Age: 25, City: "Los Angeles"},
			User{Name: "Charlie", Age: 35, City: "Chicago"},
		),
		HTMLTemplate[User](`<div class="user">
  <h2>{{.Name}}</h2>
  <p>Age: {{.Age}}</p>
  <p>City: {{.City}}</p>
</div>`),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: <div class="user">
	//   <h2>Alice</h2>
	//   <p>Age: 30</p>
	//   <p>City: New York</p>
	// </div>
	// Next: <div class="user">
	//   <h2>Bob</h2>
	//   <p>Age: 25</p>
	//   <p>City: Los Angeles</p>
	// </div>
	// Next: <div class="user">
	//   <h2>Charlie</h2>
	//   <p>Age: 35</p>
	//   <p>City: Chicago</p>
	// </div>
	// Completed
}

func ExampleTextTemplate_withSimpleTypes() {
	// Process simple types with text template
	observable := ro.Pipe1(
		ro.Just("Alice", "Bob", "Charlie"),
		TextTemplate[string]("Hello {{.}}!"),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Hello Alice!
	// Next: Hello Bob!
	// Next: Hello Charlie!
	// Completed
}

func ExampleHTMLTemplate_withSimpleTypes() {
	// Process simple types with HTML template
	observable := ro.Pipe1(
		ro.Just("Alice", "Bob", "Charlie"),
		HTMLTemplate[string](`<span class="name">{{.}}</span>`),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: <span class="name">Alice</span>
	// Next: <span class="name">Bob</span>
	// Next: <span class="name">Charlie</span>
	// Completed
}

func ExampleTextTemplate_withMaps() {
	// Process maps with text template
	observable := ro.Pipe1(
		ro.Just(
			map[string]interface{}{"name": "Alice", "age": 30, "city": "New York"},
			map[string]interface{}{"name": "Bob", "age": 25, "city": "Los Angeles"},
		),
		TextTemplate[map[string]interface{}]("User {{.name}} is {{.age}} years old from {{.city}}."),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: User Alice is 30 years old from New York.
	// Next: User Bob is 25 years old from Los Angeles.
	// Completed
}

func ExampleHTMLTemplate_withMaps() {
	// Process maps with HTML template
	observable := ro.Pipe1(
		ro.Just(
			map[string]interface{}{"name": "Alice", "age": 30, "city": "New York"},
			map[string]interface{}{"name": "Bob", "age": 25, "city": "Los Angeles"},
		),
		HTMLTemplate[map[string]interface{}](`<div class="user">
  <h2>{{.name}}</h2>
  <p>Age: {{.age}}</p>
  <p>City: {{.city}}</p>
</div>`),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: <div class="user">
	//   <h2>Alice</h2>
	//   <p>Age: 30</p>
	//   <p>City: New York</p>
	// </div>
	// Next: <div class="user">
	//   <h2>Bob</h2>
	//   <p>Age: 25</p>
	//   <p>City: Los Angeles</p>
	// </div>
	// Completed
}

func ExampleTextTemplate_withConditionals() {
	// Process data with conditional templates
	type Person struct {
		Name string
		Age  int
	}

	observable := ro.Pipe1(
		ro.Just(
			Person{Name: "Alice", Age: 30},
			Person{Name: "Bob", Age: 17},
			Person{Name: "Charlie", Age: 25},
		),
		TextTemplate[Person](`{{.Name}} is {{if ge .Age 18}}an adult{{else}}a minor{{end}} ({{.Age}} years old).`),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: Alice is an adult (30 years old).
	// Next: Bob is a minor (17 years old).
	// Next: Charlie is an adult (25 years old).
	// Completed
}

func ExampleHTMLTemplate_withConditionals() {
	// Process data with conditional HTML templates
	type Person struct {
		Name string
		Age  int
	}

	observable := ro.Pipe1(
		ro.Just(
			Person{Name: "Alice", Age: 30},
			Person{Name: "Bob", Age: 17},
			Person{Name: "Charlie", Age: 25},
		),
		HTMLTemplate[Person](`<div class="person {{if ge .Age 18}}adult{{else}}minor{{end}}">
  <h2>{{.Name}}</h2>
  <p>Age: {{.Age}}</p>
  <p>Status: {{if ge .Age 18}}Adult{{else}}Minor{{end}}</p>
</div>`),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Next: <div class="person adult">
	//   <h2>Alice</h2>
	//   <p>Age: 30</p>
	//   <p>Status: Adult</p>
	// </div>
	// Next: <div class="person minor">
	//   <h2>Bob</h2>
	//   <p>Age: 17</p>
	//   <p>Status: Minor</p>
	// </div>
	// Next: <div class="person adult">
	//   <h2>Charlie</h2>
	//   <p>Age: 25</p>
	//   <p>Status: Adult</p>
	// </div>
	// Completed
}

func ExampleTextTemplate_withError() {
	// Process data with potential template errors
	observable := ro.Pipe1(
		ro.Just(
			User{Name: "Alice", Age: 30, City: "New York"},
			User{Name: "Bob", Age: 25, City: "Los Angeles"},
		),
		TextTemplate[User]("Hello {{.Name}}, you are {{.Age}} years old from {{.InvalidField}}."),
	)

	subscription := observable.Subscribe(ro.PrintObserver[string]())
	defer subscription.Unsubscribe()

	// Output:
	// Error: template: Hello {{.Name}}, you are {{.Age}} years old from {{.InvalidField}}.:1:51: executing "Hello {{.Name}}, you are {{.Age}} years old from {{.InvalidField}}." at <.InvalidField>: can't evaluate field InvalidField in type rotemplate.User
}
