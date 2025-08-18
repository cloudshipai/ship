# Dependency Track

OWASP Dependency-Track for Software Bill of Materials (SBOM) analysis and component vulnerability management.

## Description

OWASP Dependency-Track is an intelligent Component Analysis platform that allows organizations to identify and reduce risk in the software supply chain. It consumes CycloneDX and SPDX BOMs and provides vulnerability intelligence and policy-based risk analysis.

## MCP Tools

### BOM Upload
- **`dependency_track_upload_bom`** - Upload SBOM to Dependency Track using dtrack-cli
- **`dependency_track_upload_bom_api`** - Upload SBOM via REST API using curl

### BOM Generation
- **`dependency_track_generate_bom`** - Generate CycloneDX BOMs using language-specific tools

## Real CLI Commands Used

### dtrack-cli Tool
- `dtrack-cli --bom-path bom.xml --project-name "MyProject" --project-version "1.0" --server https://dt.example.com --api-key xxx --auto-create true`

### API via curl
- `curl -X POST https://dt.example.com/api/v1/bom -H "X-Api-Key: xxx" -F "bom=@bom.xml" -F "projectName=MyProject" -F "autoCreate=true"`

### BOM Generation
- **NPM**: `cyclonedx-bom`
- **Maven**: `mvn org.cyclonedx:cyclonedx-maven-plugin:makeBom`
- **Gradle**: `gradle cyclonedxBom`
- **Python**: `cyclonedx-py`
- **PHP/Composer**: `cyclonedx-php composer`
- **.NET**: `cyclonedx dotnet`

## Use Cases

### Software Supply Chain Security
- Track and analyze software components
- Identify known vulnerabilities in dependencies
- Monitor license compliance
- Policy-based risk assessment

### Vulnerability Management
- Continuous monitoring of component vulnerabilities
- Automated vulnerability detection
- Risk scoring and prioritization
- Vulnerability remediation tracking

### Compliance & Governance
- Software composition analysis
- License compliance monitoring
- Security policy enforcement
- Audit trail for component usage

### DevSecOps Integration
- CI/CD pipeline integration
- Automated SBOM generation and upload
- Continuous security monitoring
- Shift-left security practices

## SBOM Formats Supported

- **CycloneDX** (JSON and XML)
- **SPDX** (JSON, YAML, and RDF)

## Integration Methods

### dtrack-cli Tool
A dedicated CLI tool for uploading BOMs to Dependency Track with features like:
- Auto-project creation
- Multi-format BOM support
- Authentication via API keys
- CI/CD pipeline integration

### REST API
Direct API integration using curl or other HTTP clients:
- Multipart file upload
- JSON payload support
- Authentication via API keys
- Project management capabilities

### CI/CD Pipeline Examples

**GitHub Actions:**
```yaml
- name: Generate and Upload BOM
  run: |
    cyclonedx-bom
    dtrack-cli --server ${{ secrets.DTRACK_URL }} \
      --api-key ${{ secrets.DTRACK_API_KEY }} \
      --bom-path bom.xml \
      --project-name ${{ github.repository }} \
      --project-version ${{ github.sha }} \
      --auto-create true
```

**GitLab CI:**
```yaml
dependency_track:
  script:
    - mvn org.cyclonedx:cyclonedx-maven-plugin:makeBom
    - dtrack-cli --server $DTRACK_URL --api-key $DTRACK_API_KEY --bom-path target/bom.xml --project-name $CI_PROJECT_NAME --project-version $CI_COMMIT_SHA --auto-create true
```

## Authentication

Requires API key authentication for both CLI and API access. API keys can be generated for users or teams with appropriate permissions (BOM_UPLOAD permission required).