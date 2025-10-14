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
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/samber/ro/internal/xatomic"
)

//go:embed public.key
var embeddedPublicKey []byte

// licensePrefix is a prefix on license strings to make them easily recognized.
const licensePrefix = "ro-"
const v0LicensePrefix = licensePrefix + "00"

type Type string

const (
	TypeCommunity  Type = "community"
	TypeEnterprise Type = "enterprise"
)

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

var defaultLicense = &License{
	Type:           TypeCommunity,
	LicenseID:      "",
	OrganizationID: "",
	Environment:    EnvironmentDevelopment,
	ExpiresAt:      time.Time{},
	Online:         false,
}

type License struct {
	Type           Type        `json:"t"`
	LicenseID      string      `json:"id"`
	OrganizationID string      `json:"own"`
	Environment    Environment `json:"env"`
	ExpiresAt      time.Time   `json:"exp"` // zero value -> no expiration

	// Online or offline: when online, the license is validated remotely.
	// For future use.
	Online bool `json:"o"`
}

// IsExpired checks if a license has expired
func (l *License) IsExpired() bool {
	return !l.ExpiresAt.IsZero() && time.Now().After(l.ExpiresAt)
}

// GenerateLicense creates a signed license string using the provided private key
func GenerateLicense(license *License, privateKeyPEM []byte) (string, error) {
	// Parse the private key using the keys package
	privateKey, err := UnmarshalPrivateKey(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Serialize the license to JSON
	licenseData, err := json.Marshal(license)
	if err != nil {
		return "", fmt.Errorf("failed to marshal license: %w", err)
	}

	// Encode the container
	encodedLicense, err := EncodeDataWithPrivateKey(licenseData, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to encode license container: %w", err)
	}

	// Add the license prefix
	return v0LicensePrefix + encodedLicense, nil
}

// ParseLicense decodes and verifies a license string using the embedded public key
func ParseLicense(key string) (*License, error) {
	// Check if the license has the correct prefix
	if !strings.HasPrefix(key, v0LicensePrefix) {
		return nil, fmt.Errorf("invalid license format: missing prefix")
	}

	// Remove the prefix and decode the container
	key = key[len(v0LicensePrefix):]

	// Decode the license data
	data, err := DecodeDataWithPublicKey(key, embeddedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode license container: %w", err)
	}

	// Unmarshal the license
	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return nil, fmt.Errorf("failed to unmarshal license: %w", err)
	}

	return &license, nil
}

// On expiration, the license will be set to nil on the next call to GetLicense().
var currentLicense = xatomic.NewPointer[License](nil)

// SetLicense sets the current license.
// On invalid license or expired license, an error is returned. This error might be
// ignored by the caller, since the license is not critical to the application in some cases.
func SetLicense(key string) (err error) {
	license, err := ParseLicense(key)
	if err != nil {
		return fmt.Errorf("failed to parse license: %w", err)
	}

	if license.IsExpired() {
		return fmt.Errorf("license has expired")
	}

	currentLicense.Store(license)
	time.AfterFunc(time.Until(license.ExpiresAt), func() {
		currentLicense.Store(defaultLicense)
	})

	return nil
}

// GetLicense returns the current license.
// If the license is expired, it will be set to nil on the next call to GetLicense().
// Warning: calling this function too often might cause performance issues.
func GetLicense() *License {
	license := currentLicense.Load()
	if license == nil {
		return defaultLicense
	}

	// Commented, because it's too expensive to call IsExpired() too often.
	// if license.IsExpired() {
	// 	currentLicense.Store(defaultLicense)
	// 	return defaultLicense
	// }

	return license
}

// IsEnterpriseEnabled returns true if the current license is an enterprise license.
// Warning: calling this function too often might cause performance issues.
func IsEnterpriseEnabled() bool {
	license := GetLicense()
	return license != nil && license.Type == TypeEnterprise
}
