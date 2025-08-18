# Rekor MCP Tool

Rekor is a transparency log for storage and verification of software supply chain artifacts using real rekor-cli commands.

## Description

Rekor provides:
- Immutable transparency log for artifacts
- Cryptographic verification of entries
- Search capabilities across log entries
- Supply chain artifact attestation
- Public verification of software provenance

## MCP Functions

### `rekor_upload_artifact`
Upload artifact to Rekor transparency log using real rekor-cli.

**Parameters:**
- `artifact` (required): Path or URL to artifact file
- `signature` (required): Path or URL to signature file
- `public_key` (required): Path or URL to public key file
- `type`: Type of entry
- `pki_format`: Format of signature/public key (pgp, minisign, x509, ssh, tuf)

**CLI Command:** `rekor-cli upload --artifact <artifact> --signature <signature> --public-key <public_key> [--type <type>] [--pki-format <format>]`

### `rekor_search_log`
Search Rekor transparency log using real rekor-cli.

**Parameters:**
- `artifact`: Path or URL to artifact file
- `public_key`: Path or URL to public key file
- `sha`: SHA512, SHA256, or SHA1 sum of artifact
- `email`: Email associated with public key
- `pki_format`: Format of public key (required when using public-key)
- `operator`: Search operator (and, or)

**CLI Command:** `rekor-cli search [--artifact <artifact>] [--public-key <key>] [--sha <hash>] [--email <email>] [--pki-format <format>] [--operator <op>]`

### `rekor_get_by_uuid`
Get entry from Rekor transparency log by UUID using real rekor-cli.

**Parameters:**
- `uuid` (required): UUID of entry to retrieve
- `format`: Output format (default or tle)

**CLI Command:** `rekor-cli get --uuid <uuid> [--format <format>]`

### `rekor_get_by_index`
Get entry from Rekor transparency log by log index using real rekor-cli.

**Parameters:**
- `log_index` (required): Log index of entry to retrieve
- `format`: Output format (default or tle)

**CLI Command:** `rekor-cli get --log-index <index> [--format <format>]`

### `rekor_verify_by_uuid`
Verify entry in Rekor transparency log by UUID using real rekor-cli.

**Parameters:**
- `uuid` (required): UUID of entry to verify

**CLI Command:** `rekor-cli verify --uuid <uuid>`

### `rekor_verify_by_index`
Verify entry in Rekor transparency log by log index using real rekor-cli.

**Parameters:**
- `log_index` (required): Log index of entry to verify

**CLI Command:** `rekor-cli verify --log-index <index>`

### `rekor_verify_artifact`
Verify artifact in Rekor transparency log using real rekor-cli.

**Parameters:**
- `artifact` (required): Path or URL to artifact file
- `signature`: Path or URL to signature file
- `public_key`: Path or URL to public key file
- `type`: Type of entry to verify
- `pki_format`: Format of signature/public key (pgp, minisign, x509, ssh, tuf)

**CLI Command:** `rekor-cli verify --artifact <artifact> [--signature <signature>] [--public-key <key>] [--type <type>] [--pki-format <format>]`

## Common Use Cases

1. **Supply Chain Transparency**: Record software artifacts in public log
2. **Artifact Verification**: Verify authenticity of software packages
3. **Provenance Tracking**: Track software build and release history
4. **Compliance Auditing**: Audit software supply chain activities
5. **Security Monitoring**: Monitor for unauthorized artifact changes

## Integration with Ship CLI

All Rekor tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Upload artifacts to transparency log
- Search and verify log entries
- Track software supply chain provenance
- Implement transparency-based security

The tools use containerized execution via Dagger for consistent, isolated transparency log operations.