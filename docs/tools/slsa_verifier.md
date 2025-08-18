# SLSA Verifier MCP Tool

SLSA Verifier is a tool for verifying SLSA provenance attestations for software artifacts using real slsa-verifier CLI commands.

## Description

SLSA Verifier provides:
- SLSA provenance verification for artifacts
- Container image provenance validation
- NPM package attestation verification (experimental)
- Supply chain security through provenance checking
- Support for GitHub Actions-generated attestations

## MCP Functions

### `slsa_verifier_verify_artifact`
Verify SLSA provenance for artifact using real slsa-verifier CLI.

**Parameters:**
- `artifact` (required): Path to artifact to verify
- `provenance_path` (required): Path to SLSA provenance file
- `source_uri` (required): Expected source URI (e.g., github.com/owner/repo)
- `source_tag`: Expected source tag for verification
- `source_branch`: Expected source branch for verification
- `builder_id`: Unique builder ID for verification
- `print_provenance`: Output verified provenance

**CLI Command:** `slsa-verifier verify-artifact <artifact> --provenance-path <provenance> --source-uri <uri> [--source-tag <tag>] [--source-branch <branch>] [--builder-id <id>] [--print-provenance]`

### `slsa_verifier_verify_image`
Verify SLSA provenance for container image using real slsa-verifier CLI.

**Parameters:**
- `image` (required): Container image digest to verify
- `source_uri` (required): Expected source URI (e.g., github.com/owner/repo)
- `source_tag`: Expected source tag for verification
- `source_branch`: Expected source branch for verification
- `builder_id`: Unique builder ID for verification
- `print_provenance`: Output verified provenance

**CLI Command:** `slsa-verifier verify-image <image> --source-uri <uri> [--source-tag <tag>] [--source-branch <branch>] [--builder-id <id>] [--print-provenance]`

### `slsa_verifier_verify_npm_package`
Verify SLSA provenance for npm package using real slsa-verifier CLI (experimental).

**Parameters:**
- `package_tarball` (required): Path to npm package tarball to verify
- `attestations_path` (required): Path to attestations file
- `package_name` (required): NPM package name
- `package_version` (required): NPM package version
- `source_uri`: Expected source URI (e.g., github.com/owner/repo)
- `print_provenance`: Output verified provenance

**CLI Command:** `slsa-verifier verify-npm-package <tarball> --attestations-path <attestations> --package-name <name> --package-version <version> [--source-uri <uri>] [--print-provenance]`

**Note:** Requires `SLSA_VERIFIER_EXPERIMENTAL=1` environment variable.

### `slsa_verifier_version`
Get SLSA Verifier version information using real slsa-verifier CLI.

**CLI Command:** `slsa-verifier version`

## Common Use Cases

1. **Artifact Verification**: Verify software artifacts have valid SLSA provenance
2. **Container Security**: Validate container images were built securely
3. **Supply Chain Validation**: Ensure artifacts come from expected sources
4. **NPM Package Security**: Verify npm packages have valid attestations
5. **CI/CD Integration**: Automate provenance verification in pipelines

## Integration with Ship CLI

All SLSA Verifier tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Verify SLSA provenance for various artifact types
- Validate supply chain security
- Check container image attestations
- Implement secure software distribution practices

The tools use containerized execution via Dagger for consistent, isolated verification operations.