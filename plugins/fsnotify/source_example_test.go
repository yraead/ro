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

package rofsnotify

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/ro"
)

func ExampleNewFSListener() {
	// Monitor file system events for a directory
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := NewFSListener(tempDir)

	// Create a test file to trigger events
	testFile := filepath.Join(tempDir, "test.txt")

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(event fsnotify.Event) {
				// Handle file system event
				switch event.Op {
				case fsnotify.Create:
					fmt.Println("File was created")
				case fsnotify.Write:
					fmt.Println("File was written to")
				case fsnotify.Remove:
					fmt.Println("File was removed")
				case fsnotify.Rename:
					fmt.Println("File was renamed")
				case fsnotify.Chmod:
					fmt.Println("File permissions changed")
				}
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Trigger a file creation event
	file, _ := os.Create(testFile)
	file.Close()
	time.Sleep(100 * time.Millisecond)

	// Output: File was created
}

func ExampleNewFSListener_withMultiplePaths() {
	// Monitor multiple directories
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "subdir")
	os.MkdirAll(subDir, 0755)

	paths := []string{
		tempDir,
		subDir,
	}

	observable := NewFSListener(paths...)

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(event fsnotify.Event) {
				// Handle file system event
				if filepath.Base(filepath.Dir(event.Name)) == "subdir" {
					fmt.Println("Event from: subdir")
				} else {
					fmt.Println("Event from: main directory")
				}
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Trigger events in both directories
	file1, _ := os.Create(filepath.Join(tempDir, "file1.txt"))
	file1.Close()
	file2, _ := os.Create(filepath.Join(subDir, "file2.txt"))
	file2.Close()
	time.Sleep(100 * time.Millisecond)

	// Output:
	// Event from: main directory
	// Event from: subdir
}

func ExampleNewFSListener_withFiltering() {
	// Monitor file system events with filtering
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := ro.Pipe1(
		NewFSListener(tempDir),
		ro.Filter(func(event fsnotify.Event) bool {
			// Only process .txt files
			return filepath.Ext(event.Name) == ".txt"
		}),
	)

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(event fsnotify.Event) {
				// Handle filtered file system event
				fmt.Println("Filtered event for:", filepath.Base(event.Name))
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Create files with different extensions
	file1, _ := os.Create(filepath.Join(tempDir, "test.txt"))
	file1.Close()
	file2, _ := os.Create(filepath.Join(tempDir, "test.log"))
	file2.Close()
	time.Sleep(100 * time.Millisecond)

	// Output: Filtered event for: test.txt
}

func ExampleNewFSListener_withEventTypeFiltering() {
	// Monitor specific types of file system events
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := ro.Pipe1(
		NewFSListener(tempDir),
		ro.Filter(func(event fsnotify.Event) bool {
			// Only process create and write events
			return event.Op&(fsnotify.Create) != 0
		}),
	)

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(event fsnotify.Event) {
				// Handle create and write events only
				fmt.Println("Event type:", event.Op.String())
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Create and write to a file
	file, _ := os.Create(filepath.Join(tempDir, "test.txt"))
	file.WriteString("hello")
	file.Sync()                       // Force sync to ensure write event
	time.Sleep(50 * time.Millisecond) // Wait for write event
	file.Close()
	time.Sleep(100 * time.Millisecond)

	// Output:
	// Event type: CREATE
}

func ExampleNewFSListener_withThrottling() {
	// Monitor file system events with throttling to avoid rapid successive events
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := ro.Pipe1(
		NewFSListener(tempDir),
		ro.ThrottleTime[fsnotify.Event](100*time.Millisecond),
	)

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(event fsnotify.Event) {
				// Handle throttled file system event
				fmt.Println("Throttled event:", filepath.Base(event.Name))
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Create multiple files rapidly
	for i := 0; i < 5; i++ {
		file, _ := os.Create(filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i)))
		file.Close()
	}
	time.Sleep(200 * time.Millisecond)

	// Output: Throttled event: file0.txt
}

func ExampleNewFSListener_withErrorHandling() {
	// Monitor file system events with comprehensive error handling
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := NewFSListener(tempDir)

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(event fsnotify.Event) {
				// Handle successful file system event
				fmt.Println("Event received:", event.Op.String())
			},
			func(err error) {
				// Handle file system monitoring error
				fmt.Println("Error occurred:", err.Error())
			},
			func() {
				// Handle completion (when monitoring stops)
				fmt.Println("Monitoring completed")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Trigger an event
	file, _ := os.Create(filepath.Join(tempDir, "test.txt"))
	file.Close()
	time.Sleep(100 * time.Millisecond)

	// Output: Event received: CREATE
}

func ExampleNewFSListener_withContext() {
	// Monitor file system events with context for cancellation
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := NewFSListener(tempDir)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Set up subscription first
	subscription := observable.SubscribeWithContext(
		ctx,
		ro.NewObserverWithContext(
			func(ctx context.Context, event fsnotify.Event) {
				// Handle file system event with context
				fmt.Println("Context event:", event.Op.String())
			},
			func(ctx context.Context, err error) {
				// Handle error with context
				fmt.Println("Context error:", err.Error())
			},
			func(ctx context.Context) {
				// Handle completion with context
				fmt.Println("Context completed")
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Trigger an event
	file, _ := os.Create(filepath.Join(tempDir, "test.txt"))
	file.Close()
	time.Sleep(100 * time.Millisecond)

	// Output: Context event: CREATE
}

func ExampleNewFSListener_withTransformation() {
	// Monitor file system events and transform them
	tempDir, err := os.MkdirTemp("", "fsnotify-example")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	observable := ro.Pipe1(
		NewFSListener(tempDir),
		ro.Map(func(event fsnotify.Event) string {
			// Transform event to string representation
			return filepath.Base(event.Name) + " - " + event.Op.String()
		}),
	)

	// Set up subscription first
	subscription := observable.Subscribe(
		ro.NewObserver(
			func(eventStr string) {
				fmt.Println("Transformed:", eventStr)
			},
			func(err error) {
				// Handle error
			},
			func() {
				// Handle completion
			},
		),
	)
	defer subscription.Unsubscribe()

	// Wait for watcher to be set up
	time.Sleep(100 * time.Millisecond)

	// Trigger an event
	file, _ := os.Create(filepath.Join(tempDir, "test.txt"))
	file.Close()
	time.Sleep(100 * time.Millisecond)

	// Output: Transformed: test.txt - CREATE
}
