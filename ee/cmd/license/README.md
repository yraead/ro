# License CLI Tool

A command-line tool for managing licenses in the RO framework. This tool provides commands to generate cryptographic key pairs and create signed licenses with various parameters.

## Installation

```bash
go install github.com/samber/ro/ee/cmd/license@latest
```

## Commands

### Generate Keys

Generate a new ECDSA P521 key pair for license signing and verification.

```bash
license generate-keys
```

### Generate License

Generate a signed license with the provided parameters.

```bash
license generate-license [flags]
```

**All flags are required:**
- `--type` - License type: `community` or `enterprise`
- `--id` - License ID
- `--org` - Organization ID
- `--env` - Environment: `development` or `production`
- `--online` - Online license: `true` or `false`
- `--expires` - Expiration date in YYYY-MM-DDTHH:MM:SS format
- `--private-key` - Private key

**Examples:**

Generate a community license for development:
```bash
license generate-license \
  --type community \
  --id "my-license" \
  --org "my-organization" \
  --env development \
  --online false \
  --expires "2024-12-31T23:59:59" \
  --private-key "private.key"
```

Generate an enterprise license for production:
```bash
license generate-license \
  --type enterprise \
  --id "enterprise-license" \
  --org "my-company" \
  --env production \
  --online true \
  --expires "2024-12-31T23:59:59" \
  --private-key "private.key"
```

**Output format:**
```
License generated successfully!
License string:
ro-00<base64-encoded-license>

License details (JSON):
{
  "t": "enterprise",
  "id": "enterprise-license",
  "own": "my-company",
  "env": "production",
  "o": true,
  "exp": "2024-12-31T23:59:59Z"
}
```

## License Types

- **Community**: Basic license with limited features
- **Enterprise**: Full license with all features enabled

## Environments

- **Development**: For development and testing environments
- **Production**: For production environments

## Security Notes

- Keep the private key secure and never share it
- The public key can be embedded in applications for license verification
- Licenses are cryptographically signed and cannot be tampered with
- Expired licenses are automatically rejected by the application

## Integration

The generated license string can be used to activate the RO framework:

```go
import "github.com/samber/ro/ee/pkg/license"

// Set the license
err := license.SetLicense(licenseString)
if err != nil {
    // Handle error
}

// Check if enterprise features are enabled
if license.IsEnterpriseEnabled() {
    // Use enterprise features
}
```
