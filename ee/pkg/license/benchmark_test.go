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
	"testing"
	"time"
)

func BenchmarkGetLicense(b *testing.B) {
	currentLicense.Store(&License{
		Type:           TypeEnterprise,
		LicenseID:      "test-license",
		OrganizationID: "test-org",
		Environment:    EnvironmentProduction,
		Online:         true,
		ExpiresAt:      time.Now().Add(1 * time.Hour), // expires in 1 hour
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		GetLicense()
	}
}

func BenchmarkIsEnterpriseEnabled(b *testing.B) {
	currentLicense.Store(&License{
		Type:           TypeEnterprise,
		LicenseID:      "test-license",
		OrganizationID: "test-org",
		Environment:    EnvironmentProduction,
		Online:         true,
		ExpiresAt:      time.Now().Add(1 * time.Hour), // expires in 1 hour
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		IsEnterpriseEnabled()
	}
}
