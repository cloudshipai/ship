# Cosign Golden (Advanced) MCP Tool

Advanced Cosign workflows for container signing, attestation, and verification using real cosign CLI capabilities.

## Description

Cosign Golden provides advanced container supply chain security features:
- Keyless signing with OIDC
- Certificate-based identity verification
- eBPF program signing
- Advanced attestation with predicate types
- Offline bundle verification
- Signature triangulation and cleanup

## MCP Functions

### `cosign_advanced_sign_keyless`
Sign container image using keyless signing with OIDC.

**Parameters:**
- `image_ref` (required): Container image reference to sign
- `identity_regex`: Identity regex for verification
- `issuer`: OIDC issuer URL

**CLI Command:** `cosign sign <image_ref>`

### `cosign_advanced_verify_identity`
Verify container image signature with certificate identity.

**Parameters:**
- `image_ref` (required): Container image reference to verify
- `certificate_identity`: Certificate identity to verify
- `certificate_identity_regexp`: Certificate identity regex pattern
- `certificate_oidc_issuer`: Certificate OIDC issuer

**CLI Command:** `cosign verify [--certificate-identity <identity>] [--certificate-identity-regexp <regexp>] [--certificate-oidc-issuer <issuer>] <image_ref>`

### `cosign_advanced_upload_ebpf`
Upload eBPF program to OCI registry.

**Parameters:**
- `ebpf_path` (required): Path to eBPF program file
- `registry_url` (required): OCI registry URL

**CLI Command:** `cosign upload blob -f <ebpf_path> <registry_url>`

### `cosign_advanced_attest_type`
Create attestation with specific predicate type.

**Parameters:**
- `image_ref` (required): Container image reference
- `predicate_type` (required): Predicate type URI (e.g., https://slsa.dev/provenance/v0.2)
- `predicate_file` (required): Path to predicate JSON file
- `key`: Path to signing key

**CLI Command:** `cosign attest --predicate <predicate_file> --type <predicate_type> [--key <key>] <image_ref>`

### `cosign_advanced_verify_attestation`
Verify attestation with specific type and policy.

**Parameters:**
- `image_ref` (required): Container image reference
- `type`: Attestation type to verify
- `policy`: Policy file path for verification
- `key`: Public key for verification

**CLI Command:** `cosign verify-attestation [--type <type>] [--policy <policy>] [--key <key>] <image_ref>`

### `cosign_advanced_verify_offline`
Verify signatures using offline bundle.

**Parameters:**
- `image_ref` (required): Container image reference
- `bundle` (required): Path to offline bundle file
- `certificate_identity`: Certificate identity to verify
- `certificate_oidc_issuer`: Certificate OIDC issuer

**CLI Command:** `cosign verify --bundle <bundle> [--certificate-identity <identity>] [--certificate-oidc-issuer <issuer>] <image_ref>`

### `cosign_advanced_triangulate`
Get signature image reference for a given image.

**Parameters:**
- `image_ref` (required): Container image reference

**CLI Command:** `cosign triangulate <image_ref>`

### `cosign_advanced_clean`
Clean signatures from a given image.

**Parameters:**
- `image_ref` (required): Container image reference
- `type`: Type of signatures to clean (signature, attestation, sbom, all)

**CLI Command:** `cosign clean [--type <type>] <image_ref>`

## Common Use Cases

1. **Keyless Signing**: Sign without managing keys using OIDC
2. **Identity Verification**: Verify signatures with certificate identities
3. **Supply Chain Attestations**: Create and verify SLSA attestations
4. **eBPF Security**: Sign and verify eBPF programs
5. **Offline Verification**: Verify signatures without network access
6. **Signature Management**: Clean and manage image signatures

## Integration with Ship CLI

All Cosign Golden tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Implement advanced supply chain security workflows
- Manage keyless signing and verification
- Create and verify attestations
- Handle offline verification scenarios

The tools use containerized execution via Dagger for consistent, isolated signing operations.