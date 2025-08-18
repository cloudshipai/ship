# Cosign Advanced

Advanced container signing and verification workflows using real cosign CLI features.

## Description

This module provides advanced cosign functionality for enterprise and production use cases, including keyless signing, certificate identity verification, offline verification, and specialized attestation workflows.

## MCP Tools

### Advanced Signing
- **`cosign_advanced_sign_keyless`** - Keyless signing with OIDC integration
- **`cosign_advanced_attest_type`** - Create attestations with specific predicate types (SLSA, SPDX, etc.)

### Advanced Verification
- **`cosign_advanced_verify_identity`** - Verify signatures with certificate identity matching
- **`cosign_advanced_verify_attestation`** - Advanced attestation verification with policy support
- **`cosign_advanced_verify_offline`** - Offline verification using bundle files

### Enterprise Features
- **`cosign_advanced_upload_ebpf`** - Upload eBPF programs as OCI artifacts
- **`cosign_advanced_triangulate`** - Get signature repository references
- **`cosign_advanced_clean`** - Remove signatures and attestations from images

## Real CLI Commands Used

- `cosign sign` - Keyless and key-based signing
- `cosign verify --certificate-identity` - Certificate identity verification
- `cosign verify --certificate-identity-regexp` - Regex-based identity matching
- `cosign verify --certificate-oidc-issuer` - OIDC issuer verification
- `cosign attest --type --predicate` - Typed attestations with predicates
- `cosign verify-attestation --type --policy` - Policy-based attestation verification
- `cosign verify --bundle` - Offline verification with bundles
- `cosign upload blob` - Upload arbitrary artifacts (eBPF, etc.)
- `cosign triangulate` - Get signature repository references
- `cosign clean --type` - Clean specific signature types

## Use Cases

### Enterprise Signing Workflows
- Keyless signing with corporate OIDC providers
- Certificate identity verification for compliance
- Multi-stage signing with different identities
- Policy-based attestation verification

### Supply Chain Security
- SLSA provenance attestations
- SPDX software bill of materials
- Vulnerability attestations
- Build metadata attestations

### Air-Gapped Environments
- Offline verification with bundle files
- Pre-generated signature verification
- Detached signature workflows

### Artifact Management
- eBPF program signing and distribution
- WebAssembly module verification
- Multi-artifact signing workflows
- Signature cleanup and maintenance

## Advanced Features

- **Certificate Identity Matching**: Verify signatures against specific OIDC identities
- **Regex-based Verification**: Flexible identity matching with regular expressions
- **Policy-based Attestations**: Verify attestations against OPA policies
- **Offline Verification**: Work in air-gapped environments with bundle files
- **Artifact Cleanup**: Manage signature lifecycle and storage

## Integration

Works with enterprise OIDC providers (GitHub Actions, GitLab CI, Google Cloud Build, etc.) and supports air-gapped deployment scenarios with offline verification capabilities.