# In-Toto

Supply chain integrity and attestation framework using cryptographic verification.

## Description

In-Toto is a framework that protects the integrity of software supply chains by providing cryptographic evidence for the operations performed on code during the software development lifecycle. It verifies that each step in the supply chain was performed by authorized personnel and that no unauthorized changes were made to the code. In-Toto uses cryptographic signatures to create tamper-evident records of supply chain activities.

## MCP Tools

### Step Execution & Recording
- **`in_toto_run_step`** - Run supply chain step with attestation using in-toto-run
- **`in_toto_record`** - Record step metadata using in-toto-record (start/stop)

### Verification & Validation
- **`in_toto_verify`** - Verify supply chain integrity using in-toto-verify

### Metadata Management
- **`in_toto_sign`** - Sign in-toto metadata using in-toto-sign

## Real CLI Commands Used

### Core Commands
- `in-toto-run -n <step> --signing-key <key> -- <command>` - Execute and attest step
- `in-toto-run -n <step> -g <gpg-keyid> -- <command>` - Execute with GPG signing
- `in-toto-record start -n <step> --signing-key <key>` - Start recording step
- `in-toto-record stop -n <step> --signing-key <key>` - Stop recording step
- `in-toto-verify -l <layout> -k <keys>` - Verify supply chain
- `in-toto-sign -f <metadata> --signing-key <key>` - Sign metadata

### Advanced Options
- `-m <materials>` - Record materials (input files)
- `-p <products>` - Record products (output files)
- `-s` - Record command stdout/stderr
- `-x` - Generate metadata without executing command
- `--base-path <path>` - Set base path for relative artifacts
- `-d <link-dir>` - Directory containing link metadata
- `-v` - Verbose verification output

## Core Concepts

### Supply Chain Steps
In-Toto models the software supply chain as a series of steps:
- **Materials**: Input files/data required for the step
- **Products**: Output files/data created by the step
- **Command**: The actual operation performed
- **Functionary**: The person/system authorized to perform the step

### Cryptographic Attestation
Each step generates a signed link file containing:
- **Metadata**: Step name, timestamp, functionary identity
- **Materials Hash**: Cryptographic hash of input files
- **Products Hash**: Cryptographic hash of output files
- **Command Record**: Executed command and environment
- **Signature**: Cryptographic proof of authenticity

### Layout and Verification
- **Layout**: Defines the expected supply chain structure
- **Inspections**: Automated checks for compliance
- **Verification**: Validates all steps follow the layout

## Use Cases

### Software Development
- **Build Attestation**: Prove builds came from specific source code
- **Test Verification**: Attest that tests were run on correct artifacts
- **Code Review**: Record that code was reviewed by authorized personnel
- **Release Management**: Verify release artifacts match approved builds

### CI/CD Pipeline Security
- **Pipeline Integrity**: Ensure each CI/CD step is properly executed
- **Artifact Provenance**: Track origin and transformations of artifacts
- **Access Control**: Verify only authorized systems performed operations
- **Compliance Auditing**: Provide cryptographic proof for audit requirements

### Supply Chain Security
- **Third-party Dependencies**: Verify integrity of external components
- **Multi-organization Workflows**: Secure handoffs between organizations
- **Compliance Requirements**: Meet regulatory attestation requirements
- **Incident Response**: Trace compromised artifacts to their source

### DevSecOps Integration
- **Security Gates**: Ensure security scans were performed
- **Policy Enforcement**: Verify compliance with security policies
- **Vulnerability Management**: Track remediation of security issues
- **Risk Assessment**: Understand supply chain risk exposure

## Configuration Examples

### Basic Step Execution
```bash
# Execute build step with attestation
in-toto-run -n build --signing-key build.key -m src/ -p dist/ -- make build

# Execute test step with GPG signing
in-toto-run -n test -g test@company.com -m dist/ -p test-results.xml -- npm test

# Manual review without command execution
in-toto-run -n review --signing-key review.key -m code-changes.diff -x
```

### Multi-part Step Recording
```bash
# Start recording a long-running step
in-toto-record start -n deployment --signing-key deploy.key -m app.tar.gz

# Perform deployment operations...
kubectl apply -f deployment.yaml

# Stop recording and capture products
in-toto-record stop -n deployment --signing-key deploy.key -p deployment.yaml
```

### Supply Chain Verification
```bash
# Verify complete supply chain
in-toto-verify -l layout.json -k owner.pub,developer.pub

# Verify with specific link directory
in-toto-verify -l layout.json -k keys/ -d links/ -v

# Verify in current directory
in-toto-verify -l supply-chain-layout.json -k project-keys.pub
```

### Metadata Signing
```bash
# Sign existing metadata
in-toto-sign -f build.link --signing-key signer.key

# Sign with output to specific file
in-toto-sign -f layout.json --signing-key layout.key -o signed-layout.json

# Sign with GPG
in-toto-sign -f metadata.link -g signer@company.com
```

## Integration Patterns

### GitHub Actions
```yaml
name: In-Toto Attestation
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Install In-Toto
      run: pip install in-toto
      
    - name: Attested Build
      run: |
        in-toto-run -n build \
          --signing-key ${{ secrets.BUILD_KEY }} \
          -m src/ \
          -p dist/ \
          -- make build
          
    - name: Attested Test
      run: |
        in-toto-run -n test \
          --signing-key ${{ secrets.TEST_KEY }} \
          -m dist/ \
          -p test-results.xml \
          -- npm test
          
    - name: Upload Link Files
      uses: actions/upload-artifact@v2
      with:
        name: in-toto-links
        path: "*.link"
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Attested Build') {
            steps {
                sh '''
                in-toto-run -n jenkins-build \
                  --signing-key /secrets/jenkins.key \
                  -m src/ \
                  -p target/ \
                  -- mvn clean package
                '''
                archiveArtifacts artifacts: '*.link'
            }
        }
        stage('Verify Supply Chain') {
            steps {
                sh '''
                in-toto-verify \
                  -l supply-chain-layout.json \
                  -k /keys/project-keys.pub
                '''
            }
        }
    }
}
```

### Docker Integration
```dockerfile
# Multi-stage build with attestation
FROM python:3.9 as builder
RUN pip install in-toto

COPY src/ /src/
COPY build.key /keys/

RUN in-toto-run -n docker-build \
    --signing-key /keys/build.key \
    -m /src/ \
    -p /app/ \
    -- python setup.py build

FROM python:3.9-slim
COPY --from=builder /app/ /app/
COPY --from=builder /*.link /attestations/
```

### Supply Chain Layout
```json
{
  "signed": {
    "_type": "Layout",
    "spec_version": "1.0.0",
    "expires": "2025-12-31T23:59:59Z",
    "steps": [
      {
        "name": "build",
        "expected_materials": [
          ["MATCH", "src/*", "WITH", "PRODUCTS", "FROM", "checkout"]
        ],
        "expected_products": [
          ["CREATE", "dist/*"]
        ],
        "pubkeys": ["build-key-id"],
        "expected_command": ["make", "build"],
        "threshold": 1
      },
      {
        "name": "test",
        "expected_materials": [
          ["MATCH", "dist/*", "WITH", "PRODUCTS", "FROM", "build"]
        ],
        "expected_products": [
          ["CREATE", "test-results.xml"]
        ],
        "pubkeys": ["test-key-id"],
        "expected_command": ["npm", "test"],
        "threshold": 1
      }
    ],
    "inspections": [
      {
        "name": "verify-build-outputs",
        "expected_materials": [
          ["MATCH", "dist/*", "WITH", "PRODUCTS", "FROM", "build"]
        ],
        "expected_products": [],
        "run": ["sha256sum", "dist/*"]
      }
    ]
  },
  "signatures": []
}
```

## Best Practices

### Key Management
- **Separate Keys**: Use different keys for different steps
- **Key Rotation**: Regularly rotate signing keys
- **Secure Storage**: Store private keys securely
- **Access Control**: Limit key access to authorized personnel

### Step Design
- **Atomic Operations**: Keep steps focused and atomic
- **Clear Boundaries**: Define clear inputs and outputs
- **Repeatable**: Ensure steps can be verified independently
- **Auditable**: Include sufficient detail for auditing

### Verification Strategy
- **Regular Verification**: Verify supply chain regularly
- **Automated Checks**: Integrate verification into CI/CD
- **Layout Updates**: Keep layouts current with process changes
- **Threshold Policies**: Use appropriate signature thresholds

### Security Considerations
- **Environment Security**: Secure the execution environment
- **Network Security**: Protect against network attacks
- **Insider Threats**: Implement defense against insider threats
- **Supply Chain Attacks**: Guard against upstream compromises

## Error Handling

### Common Issues
```bash
# Missing materials
in-toto-run -n build --signing-key key.pem -m nonexistent/ -- make
# Solution: Ensure all material paths exist

# Invalid signatures
in-toto-verify -l layout.json -k wrong.pub
# Solution: Use correct public keys for verification

# Expired layout
in-toto-verify -l expired-layout.json -k keys.pub
# Solution: Update layout with new expiration date

# Missing link files
in-toto-verify -l layout.json -k keys.pub -d empty-dir/
# Solution: Ensure all required link files are present
```

### Troubleshooting
- **Debug Mode**: Use `-v` flag for verbose output
- **Key Verification**: Verify public keys match private keys
- **Path Issues**: Use absolute paths when possible
- **Permission Problems**: Ensure proper file permissions

In-Toto provides comprehensive supply chain security through cryptographic attestation, enabling organizations to verify the integrity and authenticity of their software development processes.