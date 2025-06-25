# Changelog

All notable changes to Ship CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Ship CLI
- Authentication system with secure token storage
- Terraform analysis tools integration via Dagger:
  - TFLint for Terraform linting
  - Checkov for security scanning
  - Infracost for cost estimation
  - Trivy for vulnerability scanning
  - terraform-docs for documentation generation
  - OpenInfraQuote for cost analysis
- Dagger SDK integration for containerized execution
- Comprehensive documentation and examples
- CI/CD ready with GitHub Actions examples
- Support for AWS, Azure, and GCP credentials

### Infrastructure
- GoReleaser configuration for multi-platform builds
- Automated changelog generation
- Docker image publication
- Homebrew tap support
- DEB/RPM/APK package generation