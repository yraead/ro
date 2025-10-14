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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/ro"
	"github.com/stretchr/testify/assert"
)

func TestNewFSListener(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fsnotify-test")
	is.Nil(err)
	defer os.RemoveAll(tempDir)

	// Create a temporary file to watch
	tempFile := filepath.Join(tempDir, "testfile.txt")
	f, err := os.Create(tempFile)
	is.Nil(err)
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(10 * time.Millisecond)
		_, err = f.WriteString("hello")
		is.Nil(err)
		err = f.Sync()
		is.Nil(err)
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	// Create a new FS listener
	obs := NewFSListener(tempFile)
	is.NotNil(obs)

	items, _, err := ro.CollectWithContext(ctx, obs)
	is.ErrorIs(err, context.Canceled)
	is.Len(items, 1)

	is.True(items[0].Op.Has(fsnotify.Write))
	is.Equal(tempFile, items[0].Name)
	is.Equal(fsnotify.Write, items[0].Op)
}

func TestNewFSListener_Error(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Create a new FS listener with an invalid path
	obs := NewFSListener("/invalid/path")
	is.NotNil(obs)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	items, _, err := ro.CollectWithContext(ctx, obs)
	is.Error(err)
	is.Len(items, 0)
}
