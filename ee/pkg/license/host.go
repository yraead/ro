// Copyright 2025 samber.
//
// Licensed as an Enterprise License (the "License"); you may not use
// this file except in compliance with the License. You may obtain
// a copy of the License at:
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.ee.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package rolicense

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

var (
	// instanceID is a unique identifier for this instance
	instanceID     string
	instanceIDOnce sync.Once
)

// getInstanceID returns the unique identifier for this instance.
// The ID is generated once and cached for subsequent calls.
func getInstanceID() string {
	instanceIDOnce.Do(func() {
		// Generate a random 16-byte identifier
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			// Fallback to a timestamp-based ID if random generation fails
			instanceID = fmt.Sprintf("instance-%d", time.Now().UnixNano())
		} else {
			instanceID = hex.EncodeToString(bytes)
		}
	})
	return instanceID
}

// GetInstanceID returns the unique identifier for this instance.
// This is the public API for accessing the instance ID.
func GetInstanceID() string {
	return getInstanceID()
}
