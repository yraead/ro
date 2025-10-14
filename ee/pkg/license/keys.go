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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Curve is the elliptic curve used for key generation
var Curve = elliptic.P521()

// GenerateKeys generates a new ECDSA key pair for P521
func GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// Generate ECDSA key pair directly
	privateKey, err := ecdsa.GenerateKey(Curve, rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// UnmarshalPrivateKey parses a PEM-encoded ECDSA private key and validates it uses P521 curve
func UnmarshalPrivateKey(privateKeyPEM []byte) (*ecdsa.PrivateKey, error) {
	// Parse the private key
	pk, err := base64.StdEncoding.DecodeString(string(privateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key")
	}

	privateKey, err := x509.ParseECPrivateKey(pk)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Ensure the key is using P521 curve
	if privateKey.Curve != Curve {
		return nil, fmt.Errorf("private key must use P521 curve")
	}

	return privateKey, nil
}

// UnmarshalPublicKey parses a PEM-encoded ECDSA public key and validates it uses P521 curve
func UnmarshalPublicKey(publicKeyPEM []byte) (*ecdsa.PublicKey, error) {
	// Parse the public key
	pk, err := base64.StdEncoding.DecodeString(string(publicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(pk)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not ECDSA")
	}

	// Ensure the public key is using P521 curve
	if ecdsaPublicKey.Curve != Curve {
		return nil, fmt.Errorf("public key must use P521 curve")
	}

	return ecdsaPublicKey, nil
}

// MarshalPrivateKey encodes an ECDSA private key to base64 format
func MarshalPrivateKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	return []byte(base64.StdEncoding.EncodeToString(privateKeyBytes)), nil
}

// MarshalPublicKey encodes an ECDSA public key to base64 format
func MarshalPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	return []byte(base64.StdEncoding.EncodeToString(publicKeyBytes)), nil
}

// signedContainer represents the structure of a signed payload
type signedContainer struct {
	Data      []byte `json:"d"`
	Signature []byte `json:"s"`
}

func EncodeDataWithPrivateKey(data []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	// Sign the data
	container, err := signData(data, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	// Serialize the container
	containerData, err := json.Marshal(container)
	if err != nil {
		return "", fmt.Errorf("failed to marshal container: %w", err)
	}

	return base64.StdEncoding.EncodeToString(containerData), nil
}

func DecodeDataWithPublicKey(data string, publicKey []byte) ([]byte, error) {
	// Decode from base64
	containerData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Parse the container
	var container signedContainer
	if err := json.Unmarshal(containerData, &container); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container: %w", err)
	}

	// Verify the data
	if err := verifyDataSignature(&container, publicKey); err != nil {
		return nil, fmt.Errorf("failed to verify payload: %w", err)
	}

	return container.Data, nil
}

// signData signs data using the provided private key
func signData(data []byte, privateKey *ecdsa.PrivateKey) (*signedContainer, error) {
	// Hash the data before signing
	hashedData := sha256.Sum256(data)

	// Sign the hashed data using ECDSA
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hashedData[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	// Create a container with data and signature
	container := &signedContainer{
		Data:      data,
		Signature: signature,
	}

	return container, nil
}

// verifyDataSignature verifies data using the provided public key
func verifyDataSignature(container *signedContainer, publicKey []byte) error {
	// Parse the embedded public key using the keys package
	ecdsaPublicKey, err := UnmarshalPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to parse embedded public key: %w", err)
	}

	// Hash the data for verification
	hashedData := sha256.Sum256(container.Data)

	// Verify the signature using ECDSA
	if !ecdsa.VerifyASN1(ecdsaPublicKey, hashedData[:], container.Signature) {
		return fmt.Errorf("failed to verify signature")
	}

	return nil
}
