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
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLicense_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		license  *License
		expected bool
	}{
		{
			name: "expired license",
			license: &License{
				Type:           TypeCommunity,
				LicenseID:      "test-license",
				OrganizationID: "test-org",
				Environment:    EnvironmentDevelopment,
				Online:         false,
				ExpiresAt:      time.Now().Add(-1 * time.Hour), // expired 1 hour ago
			},
			expected: true,
		},
		{
			name: "valid license",
			license: &License{
				Type:           TypeEnterprise,
				LicenseID:      "test-license",
				OrganizationID: "test-org",
				Environment:    EnvironmentProduction,
				Online:         true,
				ExpiresAt:      time.Now().Add(1 * time.Hour), // expires in 1 hour
			},
			expected: false,
		},
		{
			name: "license with zero expiration (no expiration)",
			license: &License{
				Type:           TypeCommunity,
				LicenseID:      "test-license",
				OrganizationID: "test-org",
				Environment:    EnvironmentDevelopment,
				Online:         false,
				ExpiresAt:      time.Time{}, // zero value means no expiration
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.license.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateLicense(t *testing.T) {
	// Generate test keys
	privateKey, _, err := GenerateKeys()
	require.NoError(t, err)

	privateKeyPEM, err := MarshalPrivateKey(privateKey)
	require.NoError(t, err)

	tests := []struct {
		name    string
		license *License
		wantErr bool
	}{
		{
			name: "valid community license",
			license: &License{
				Type:           TypeCommunity,
				LicenseID:      "test-community-license",
				OrganizationID: "test-org",
				Environment:    EnvironmentDevelopment,
				Online:         false,
				ExpiresAt:      time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "valid enterprise license",
			license: &License{
				Type:           TypeEnterprise,
				LicenseID:      "test-enterprise-license",
				OrganizationID: "test-org",
				Environment:    EnvironmentProduction,
				Online:         true,
				ExpiresAt:      time.Now().Add(365 * 24 * time.Hour),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			licenseString, err := GenerateLicense(tt.license, privateKeyPEM)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, licenseString)
			assert.True(t, strings.HasPrefix(licenseString, v0LicensePrefix))
		})
	}
}

func TestParseLicense_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "invalid prefix",
			key:     "invalid-prefix-license",
			wantErr: true,
		},
		{
			name:    "empty string",
			key:     "",
			wantErr: true,
		},
		{
			name:    "malformed base64",
			key:     v0LicensePrefix + "invalid-base64!@#",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			license, err := ParseLicense(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, license)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, license)
			}
		})
	}
}

func TestGetLicense(t *testing.T) {
	t.Run("no license set", func(t *testing.T) {
		// Clear any existing license
		currentLicense.Store(nil)

		license := GetLicense()
		assert.Equal(t, defaultLicense, license)
	})

	t.Run("expired license", func(t *testing.T) {
		expiredLicense := &License{
			Type:           TypeCommunity,
			LicenseID:      "expired-license",
			OrganizationID: "test-org",
			Environment:    EnvironmentDevelopment,
			Online:         false,
			ExpiresAt:      time.Now().Add(-1 * time.Hour),
		}
		currentLicense.Store(expiredLicense)
		defer currentLicense.Store(nil)

		license := GetLicense()
		assert.Equal(t, defaultLicense, license)
	})

	t.Run("valid license", func(t *testing.T) {
		validLicense := &License{
			Type:           TypeEnterprise,
			LicenseID:      "valid-license",
			OrganizationID: "test-org",
			Environment:    EnvironmentProduction,
			Online:         true,
			ExpiresAt:      time.Now().Add(1 * time.Hour),
		}
		currentLicense.Store(validLicense)
		defer currentLicense.Store(nil)

		license := GetLicense()
		require.NotNil(t, license)
		assert.Equal(t, validLicense.LicenseID, license.LicenseID)
	})
}

func TestSetLicense_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "invalid license format",
			key:     "invalid-license-format",
			wantErr: true,
		},
		{
			name:    "empty string",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any existing license
			currentLicense.Store(nil)

			err := SetLicense(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				license := GetLicense()
				assert.Equal(t, defaultLicense, license)
			} else {
				assert.NoError(t, err)
				license := GetLicense()
				assert.NotEqual(t, defaultLicense, license)
			}
		})
	}
}

func TestLicenseConstants(t *testing.T) {
	assert.Equal(t, "ro-", licensePrefix)
	assert.Equal(t, "ro-00", v0LicensePrefix)
	assert.Equal(t, Type("community"), TypeCommunity)
	assert.Equal(t, Type("enterprise"), TypeEnterprise)
	assert.Equal(t, Environment("development"), EnvironmentDevelopment)
	assert.Equal(t, Environment("production"), EnvironmentProduction)
}

func TestLicenseStructure(t *testing.T) {
	license := &License{
		Type:           TypeEnterprise,
		LicenseID:      "test-structure-license",
		OrganizationID: "test-org",
		Environment:    EnvironmentProduction,
		Online:         true,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	// Test JSON marshaling
	data, err := json.Marshal(license)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test JSON unmarshaling
	var unmarshaledLicense License
	err = json.Unmarshal(data, &unmarshaledLicense)
	require.NoError(t, err)

	// Verify the unmarshaled license matches the original
	assert.Equal(t, license.Type, unmarshaledLicense.Type)
	assert.Equal(t, license.LicenseID, unmarshaledLicense.LicenseID)
	assert.Equal(t, license.OrganizationID, unmarshaledLicense.OrganizationID)
	assert.Equal(t, license.Environment, unmarshaledLicense.Environment)
	assert.Equal(t, license.Online, unmarshaledLicense.Online)
}

func TestIsEnterpriseEnabled(t *testing.T) {
	t.Run("no license set", func(t *testing.T) {
		currentLicense.Store(nil)
		assert.False(t, IsEnterpriseEnabled())
	})

	t.Run("community license", func(t *testing.T) {
		communityLicense := &License{
			Type:           TypeCommunity,
			LicenseID:      "community-license",
			OrganizationID: "test-org",
			Environment:    EnvironmentDevelopment,
			Online:         false,
			ExpiresAt:      time.Now().Add(24 * time.Hour),
		}
		currentLicense.Store(communityLicense)
		defer currentLicense.Store(nil)

		assert.False(t, IsEnterpriseEnabled())
	})

	t.Run("enterprise license", func(t *testing.T) {
		enterpriseLicense := &License{
			Type:           TypeEnterprise,
			LicenseID:      "enterprise-license",
			OrganizationID: "test-org",
			Environment:    EnvironmentProduction,
			Online:         true,
			ExpiresAt:      time.Now().Add(24 * time.Hour),
		}
		currentLicense.Store(enterpriseLicense)
		defer currentLicense.Store(nil)

		assert.True(t, IsEnterpriseEnabled())
	})

	t.Run("expired enterprise license", func(t *testing.T) {
		expiredEnterpriseLicense := &License{
			Type:           TypeEnterprise,
			LicenseID:      "expired-enterprise-license",
			OrganizationID: "test-org",
			Environment:    EnvironmentProduction,
			Online:         true,
			ExpiresAt:      time.Now().Add(-1 * time.Hour),
		}
		currentLicense.Store(expiredEnterpriseLicense)
		defer currentLicense.Store(nil)

		assert.False(t, IsEnterpriseEnabled())
	})
}
