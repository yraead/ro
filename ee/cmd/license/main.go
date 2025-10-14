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


package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	rolicense "github.com/samber/ro/ee/pkg/license"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "license",
	Short: "License CLI Tool for generating keys and licenses",
	Long: `A CLI tool for managing licenses in the RO framework.
	
This tool provides commands to generate cryptographic key pairs and create signed licenses
with various parameters including type, organization, environment, and expiration.`,
}

var generateKeysCmd = &cobra.Command{
	Use:   "generate-keys",
	Short: "Generate a new private/public key pair",
	Long: `Generate a new ECDSA P521 key pair for license signing and verification.
	
The keys will be printed to stdout in PEM format. You can redirect the output to files:
- private.key (private key, should be kept secure)
- public.key (public key, can be shared)`,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate the key pair
		privateKey, publicKey, err := rolicense.GenerateKeys()
		if err != nil {
			fmt.Printf("Error generating keys: %v\n", err)
			os.Exit(1)
		}

		// Marshal the keys to PEM format
		privateKeyPEM, err := rolicense.MarshalPrivateKey(privateKey)
		if err != nil {
			fmt.Printf("Error marshaling private key: %v\n", err)
			os.Exit(1)
		}

		publicKeyPEM, err := rolicense.MarshalPublicKey(publicKey)
		if err != nil {
			fmt.Printf("Error marshaling public key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println("Keys generated successfully!")
		fmt.Println()
		fmt.Println(">>>> PRIVATE KEY:")
		fmt.Println("(Keep this secure and private)")
		fmt.Println()
		fmt.Println(string(privateKeyPEM))
		fmt.Println()
		fmt.Println(">>>> PUBLIC KEY:")
		fmt.Println("(This can be shared and embedded in applications)")
		fmt.Println()
		fmt.Println(string(publicKeyPEM))
		fmt.Println()
	},
}

var generateLicenseCmd = &cobra.Command{
	Use:   "generate-license",
	Short: "Generate a license with specified parameters",
	Long: `Generate a signed license with the provided parameters.
	
The license will be cryptographically signed using the specified private key and
output as a base64-encoded string that can be used to activate the application.

All parameters are required to ensure explicit license configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag values
		licenseType, _ := cmd.Flags().GetString("type")
		licenseID, _ := cmd.Flags().GetString("id")
		organizationID, _ := cmd.Flags().GetString("org")
		environment, _ := cmd.Flags().GetString("env")
		online, _ := cmd.Flags().GetBool("online")
		expiresAt, _ := cmd.Flags().GetString("expires")
		privateKey, _ := cmd.Flags().GetString("private-key")

		// Validate required parameters
		if licenseID == "" {
			fmt.Println("Error: license ID is required")
			cmd.Help()
			os.Exit(1)
		}

		if organizationID == "" {
			fmt.Println("Error: organization ID is required")
			cmd.Help()
			os.Exit(1)
		}

		if licenseType == "" {
			fmt.Println("Error: license type is required")
			cmd.Help()
			os.Exit(1)
		}

		if environment == "" {
			fmt.Println("Error: environment is required")
			cmd.Help()
			os.Exit(1)
		}

		if expiresAt == "" {
			fmt.Println("Error: expiration date is required")
			cmd.Help()
			os.Exit(1)
		}

		if privateKey == "" {
			fmt.Println("Error: private key is required")
			cmd.Help()
			os.Exit(1)
		}

		// Parse license type
		var licenseTypeEnum rolicense.Type
		switch licenseType {
		case "community":
			licenseTypeEnum = rolicense.TypeCommunity
		case "enterprise":
			licenseTypeEnum = rolicense.TypeEnterprise
		default:
			fmt.Printf("Error: invalid license type '%s'. Must be 'community' or 'enterprise'\n", licenseType)
			os.Exit(1)
		}

		// Parse environment
		var environmentEnum rolicense.Environment
		switch environment {
		case "development":
			environmentEnum = rolicense.EnvironmentDevelopment
		case "production":
			environmentEnum = rolicense.EnvironmentProduction
		default:
			fmt.Printf("Error: invalid environment '%s'. Must be 'development' or 'production'\n", environment)
			os.Exit(1)
		}

		// Parse expiration date
		var expiresAtTime time.Time
		var err error
		expiresAtTime, err = time.Parse("2006-01-02T15:04:05", expiresAt)
		if err != nil {
			fmt.Printf("Error: invalid expiration date format '%s'. Use YYYY-MM-DDTHH:MM:SS format\n", expiresAt)
			os.Exit(1)
		}

		// Create license
		licenseObj := &rolicense.License{
			Type:           licenseTypeEnum,
			LicenseID:      licenseID,
			OrganizationID: organizationID,
			Environment:    environmentEnum,
			Online:         online,
			ExpiresAt:      expiresAtTime,
		}

		// Generate the license string
		licenseString, err := rolicense.GenerateLicense(licenseObj, []byte(privateKey))
		if err != nil {
			fmt.Printf("Error generating license: %v\n", err)
			os.Exit(1)
		}

		// Output the license
		fmt.Println()
		fmt.Println("License generated successfully!")
		fmt.Println()
		fmt.Println("License string:")
		fmt.Println(licenseString)
		fmt.Println()

		// Also output as JSON for reference
		licenseJSON, err := json.MarshalIndent(licenseObj, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling license to JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nLicense details (JSON):")
		fmt.Println(string(licenseJSON))
	},
}

func init() {
	// Add commands to root
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(generateKeysCmd)
	rootCmd.AddCommand(generateLicenseCmd)

	// Add flags to generate-license command
	generateLicenseCmd.Flags().String("type", "", "License type (community|enterprise) (required)")
	generateLicenseCmd.Flags().String("id", "", "License ID (required)")
	generateLicenseCmd.Flags().String("org", "", "Organization ID (required)")
	generateLicenseCmd.Flags().String("env", "", "Environment (development|production) (required)")
	generateLicenseCmd.Flags().Bool("online", false, "Online license (required)")
	generateLicenseCmd.Flags().String("expires", "", "Expiration date (YYYY-MM-DDTHH:MM:SS format) (required)")
	generateLicenseCmd.Flags().String("private-key", "", "Path to private key file (required)")

	// Mark all flags as required
	generateLicenseCmd.MarkFlagRequired("type")
	generateLicenseCmd.MarkFlagRequired("id")
	generateLicenseCmd.MarkFlagRequired("org")
	generateLicenseCmd.MarkFlagRequired("env")
	generateLicenseCmd.MarkFlagRequired("online")
	generateLicenseCmd.MarkFlagRequired("expires")
	generateLicenseCmd.MarkFlagRequired("private-key")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
