# Cosign

Container signing and verification tool for supply chain security.

## Description

Cosign is a tool for signing and verifying container images and other artifacts using a variety of signature formats. It's part of the Sigstore project and provides keyless signing capabilities.

## MCP Tools

### Image Operations
- **`cosign_sign_image`** - Sign container images with cryptographic signatures
- **`cosign_verify_image`** - Verify container image signatures
- **`cosign_copy_image`** - Copy images between registries

### Blob Operations
- **`cosign_sign_blob`** - Sign arbitrary files and artifacts
- **`cosign_verify_blob`** - Verify blob signatures

### Attestations
- **`cosign_attest`** - Create and sign attestations for container images
- **`cosign_verify_attestation`** - Verify attestations attached to container images

### Key Management
- **`cosign_generate_key`** - Generate cryptographic key pairs for signing

### Artifact Upload
- **`cosign_upload_blob`** - Upload generic artifacts as blobs to OCI registries
- **`cosign_upload_wasm`** - Upload WebAssembly modules to registries

### WebAssembly
- **`cosign_sign_wasm`** - Sign WebAssembly artifacts

### Utility
- **`cosign_get_version`** - Get Cosign version information

## Real CLI Commands Used

- `cosign sign` - Sign container images or artifacts
- `cosign verify` - Verify signatures on images or artifacts
- `cosign generate-key-pair` - Generate signing key pairs
- `cosign attest` - Create attestations with predicates
- `cosign verify-attestation` - Verify attestations
- `cosign sign-blob` - Sign arbitrary files
- `cosign verify-blob` - Verify blob signatures
- `cosign upload blob` - Upload artifacts to registries
- `cosign upload wasm` - Upload WebAssembly modules
- `cosign copy` - Copy images between registries
- `cosign version` - Show version information

## Use Cases

- Container image signing for supply chain security
- Artifact verification in CI/CD pipelines
- Software bill of materials (SBOM) attestations
- Keyless signing with OIDC providers
- WebAssembly module signing
- Compliance with software supply chain requirements

## Integration

Works with OCI-compliant registries and supports both key-based and keyless signing workflows using OIDC providers like GitHub Actions, GitLab CI, and others.