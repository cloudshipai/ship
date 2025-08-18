# OSV Scanner

Open Source Vulnerability scanner for comprehensive dependency vulnerability assessment and guided remediation across multiple programming languages and package ecosystems.

## Description

OSV Scanner is Google's official open-source vulnerability scanner that provides comprehensive security assessment for software dependencies across multiple programming languages, package managers, and ecosystems. Built on the Open Source Vulnerabilities (OSV) database—a distributed, standardized vulnerability database for open source software—OSV Scanner offers fast, accurate vulnerability detection with minimal false positives. The tool supports extensive ecosystem coverage including npm (Node.js), pip (Python), Maven/Gradle (Java), Go modules, Cargo (Rust), Composer (PHP), RubyGems, NuGet (.NET), and many others. OSV Scanner features container image scanning, SBOM (Software Bill of Materials) analysis, license compliance checking, guided remediation capabilities, and seamless integration with CI/CD pipelines, making it essential for modern software supply chain security, DevSecOps workflows, and regulatory compliance initiatives.

## MCP Tools

### Source Code and Directory Scanning
- **`osv_scanner_scan_source`** - Scan source directory for open source vulnerabilities using real osv-scanner CLI
- **`osv_scanner_verbose_scan`** - Run OSV Scanner with verbose logging using real osv-scanner CLI

### Container and Image Analysis
- **`osv_scanner_scan_image`** - Scan container image for vulnerabilities using real osv-scanner CLI

### Package and Dependency Scanning
- **`osv_scanner_scan_lockfile`** - Scan specific lockfile for vulnerabilities using real osv-scanner CLI
- **`osv_scanner_scan_manifest`** - Scan package manifest file for vulnerabilities using real osv-scanner CLI
- **`osv_scanner_scan_sbom`** - Scan SBOM file for vulnerabilities using real osv-scanner CLI

### License and Compliance Management
- **`osv_scanner_license_scan`** - Scan for license compliance using real osv-scanner CLI

### Offline and Enterprise Operations
- **`osv_scanner_offline_scan`** - Scan using offline vulnerability databases with real osv-scanner CLI

### Remediation and Fix Management
- **`osv_scanner_fix`** - Apply guided remediation for vulnerabilities using real osv-scanner CLI

### Reporting and Visualization
- **`osv_scanner_serve_report`** - Generate and serve HTML vulnerability report locally using real osv-scanner CLI

## Real CLI Commands Used

### Basic Scanning Commands
- `osv-scanner scan source [-r] <directory>` - Scan source directory recursively
- `osv-scanner scan image <image:tag>` - Scan container image
- `osv-scanner -L <lockfile>` - Scan specific lockfile
- `osv-scanner -M <manifest>` - Scan package manifest file
- `osv-scanner --sbom <sbom-file>` - Scan SBOM file

### License Scanning Commands
- `osv-scanner --licenses <path>` - Scan for licenses
- `osv-scanner --licenses="MIT,Apache-2.0" <path>` - Scan with allowed license list

### Output and Format Control
- `osv-scanner --output <file> --format <json|table|sarif> <path>` - Control output format
- `osv-scanner --serve [--port <port>] <path>` - Serve HTML report

### Offline and Configuration
- `osv-scanner --offline --download-offline-databases <path>` - Offline scanning
- `osv-scanner --config <config-file> <path>` - Use custom configuration
- `osv-scanner --verbosity <level> <path>` - Control logging verbosity

### Guided Remediation
- `osv-scanner fix -M <manifest> -L <lockfile>` - Interactive remediation
- `osv-scanner fix --strategy <strategy> --max-depth <depth> <path>` - Automated remediation

### Advanced Options
- `osv-scanner --all-packages <path>` - Output all packages in JSON format
- `osv-scanner --offline-vulnerabilities <db-path> <path>` - Use local vulnerability database

## Use Cases

### Software Supply Chain Security
- **Dependency Vulnerability Assessment**: Identify known vulnerabilities in open source dependencies
- **Container Security**: Scan container images for vulnerable packages and libraries
- **SBOM Analysis**: Analyze Software Bill of Materials for security compliance
- **License Compliance**: Ensure compliance with organizational license policies

### DevSecOps Integration
- **CI/CD Security Gates**: Integrate vulnerability scanning into build and deployment pipelines
- **Pre-commit Hooks**: Prevent vulnerable dependencies from entering version control
- **Security Testing**: Automated security testing as part of development workflows
- **Continuous Monitoring**: Ongoing vulnerability assessment of deployed applications

### Enterprise Risk Management
- **Risk Assessment**: Identify and prioritize security risks across software portfolios
- **Vulnerability Management**: Track and remediate security vulnerabilities systematically
- **Compliance Reporting**: Generate security reports for regulatory and audit requirements
- **Security Governance**: Implement organization-wide security policies and standards

### Development and Maintenance
- **Dependency Updates**: Guided remediation and dependency upgrade recommendations
- **Security Monitoring**: Continuous monitoring for new vulnerabilities in dependencies
- **Code Review**: Security assessment during code review processes
- **Release Management**: Ensure secure releases through pre-deployment scanning

## Configuration Examples

### Basic Vulnerability Scanning
```bash
# Scan current directory recursively
osv-scanner scan source -r .

# Scan specific directory
osv-scanner scan source -r /path/to/project

# Scan with JSON output
osv-scanner -r --format json --output results.json .

# Scan container image
osv-scanner scan image nginx:latest
osv-scanner scan image myapp:1.2.3

# Scan specific lockfiles
osv-scanner -L package-lock.json
osv-scanner -L go.mod
osv-scanner -L requirements.txt
osv-scanner -L Cargo.lock

# Scan package manifests
osv-scanner -M package.json
osv-scanner -M pom.xml
osv-scanner -M composer.json
```

### License Compliance Scanning
```bash
# Basic license scanning
osv-scanner --licenses .

# Scan with allowed license list
osv-scanner --licenses="MIT,Apache-2.0,BSD-3-Clause" .

# Combine vulnerability and license scanning
osv-scanner --licenses --format json --output compliance-report.json .

# License scanning with output to file
osv-scanner --licenses --output license-report.txt ./src
```

### SBOM and Supply Chain Analysis
```bash
# Scan SPDX SBOM
osv-scanner --sbom sbom.spdx.json

# Scan CycloneDX SBOM
osv-scanner --sbom sbom.cyclonedx.json

# Generate SBOM and scan (using external tools)
# Example with syft for SBOM generation
syft packages . -o spdx-json=sbom.spdx.json
osv-scanner --sbom sbom.spdx.json --format json --output sbom-scan-results.json

# Scan container image SBOM
syft packages nginx:latest -o cyclonedx-json=nginx-sbom.json
osv-scanner --sbom nginx-sbom.json
```

### Advanced Scanning Options
```bash
# Verbose scanning with debug output
osv-scanner --verbosity debug -r .

# Scan with custom configuration
osv-scanner --config custom-config.yml -r .

# Output all packages (useful for inventory)
osv-scanner --all-packages --format json --output inventory.json .

# Scan with SARIF output for security tools integration
osv-scanner --format sarif --output results.sarif -r .

# Generate and serve interactive HTML report
osv-scanner --serve --port 8080 -r .
```

### Offline and Enterprise Scanning
```bash
# Download offline vulnerability databases
osv-scanner --offline --download-offline-databases ./offline-db

# Perform offline scanning
osv-scanner --offline --offline-vulnerabilities ./offline-db -r .

# Offline scanning with recursive directory scan
osv-scanner --offline --offline-vulnerabilities ./offline-db -r /enterprise/projects
```

## Advanced Usage

### Comprehensive Security Assessment Pipeline
```bash
#!/bin/bash
# comprehensive-security-assessment.sh

PROJECT_DIR="$1"
ASSESSMENT_ID="security-$(date +%Y%m%d-%H%M%S)"
RESULTS_DIR="security-assessment-$ASSESSMENT_ID"

if [[ -z "$PROJECT_DIR" ]]; then
    echo "Usage: $0 <project-directory>"
    exit 1
fi

echo "Starting comprehensive security assessment for: $PROJECT_DIR"
mkdir -p "$RESULTS_DIR"

# Phase 1: Project structure analysis
echo "Phase 1: Analyzing project structure..."
find "$PROJECT_DIR" -name "package*.json" -o -name "go.mod" -o -name "requirements*.txt" \
    -o -name "Cargo.toml" -o -name "pom.xml" -o -name "composer.json" \
    > "$RESULTS_DIR/package-files.txt"

echo "Found $(wc -l < "$RESULTS_DIR/package-files.txt") package files"

# Phase 2: Dependency vulnerability scanning
echo "Phase 2: Scanning for dependency vulnerabilities..."
osv-scanner scan source -r "$PROJECT_DIR" \
    --format json \
    --output "$RESULTS_DIR/vulnerability-scan.json" \
    > "$RESULTS_DIR/scan-output.txt" 2>&1

VULN_EXIT_CODE=$?

# Phase 3: License compliance scanning
echo "Phase 3: Scanning for license compliance..."
osv-scanner --licenses "$PROJECT_DIR" \
    --output "$RESULTS_DIR/license-scan.txt" \
    > "$RESULTS_DIR/license-output.txt" 2>&1

# Phase 4: Generate detailed reports in multiple formats
echo "Phase 4: Generating detailed security reports..."

# SARIF format for security tools
osv-scanner scan source -r "$PROJECT_DIR" \
    --format sarif \
    --output "$RESULTS_DIR/security-results.sarif" \
    2>/dev/null || true

# Table format for human readability
osv-scanner scan source -r "$PROJECT_DIR" \
    --format table \
    --output "$RESULTS_DIR/vulnerability-report.txt" \
    2>/dev/null || true

# Phase 5: Package inventory
echo "Phase 5: Creating package inventory..."
osv-scanner --all-packages "$PROJECT_DIR" \
    --format json \
    --output "$RESULTS_DIR/package-inventory.json" \
    2>/dev/null || true

# Phase 6: Parse and analyze results
echo "Phase 6: Analyzing security findings..."
CRITICAL_VULNS=0
HIGH_VULNS=0
MEDIUM_VULNS=0
LOW_VULNS=0
TOTAL_VULNS=0

if [[ -f "$RESULTS_DIR/vulnerability-scan.json" ]]; then
    # Parse vulnerability counts from JSON (requires jq)
    if command -v jq > /dev/null; then
        TOTAL_VULNS=$(jq -r '.results[]?.packages[]?.vulnerabilities? | length' "$RESULTS_DIR/vulnerability-scan.json" 2>/dev/null | \
            awk '{sum += $1} END {print sum+0}')
        
        # Count by severity (if available in OSV data)
        CRITICAL_VULNS=$(jq -r '.results[]?.packages[]?.vulnerabilities[]? | select(.database_specific?.severity == "CRITICAL") | .id' \
            "$RESULTS_DIR/vulnerability-scan.json" 2>/dev/null | wc -l)
        HIGH_VULNS=$(jq -r '.results[]?.packages[]?.vulnerabilities[]? | select(.database_specific?.severity == "HIGH") | .id' \
            "$RESULTS_DIR/vulnerability-scan.json" 2>/dev/null | wc -l)
    else
        # Fallback: count unique vulnerability IDs
        TOTAL_VULNS=$(grep -o '"id":"[^"]*"' "$RESULTS_DIR/vulnerability-scan.json" 2>/dev/null | sort -u | wc -l || echo "0")
    fi
fi

# Phase 7: Generate executive summary
echo "Phase 7: Creating executive summary..."
cat > "$RESULTS_DIR/security-summary.md" <<EOF
# Security Assessment Summary

**Assessment ID**: $ASSESSMENT_ID  
**Project**: $PROJECT_DIR  
**Date**: $(date)  
**Scanner**: OSV Scanner (Google)

## Vulnerability Assessment Results

### Overall Security Status
$(if [[ $VULN_EXIT_CODE -eq 0 ]]; then
    echo "✅ **SECURE**: No known vulnerabilities detected"
elif [[ $TOTAL_VULNS -gt 50 ]]; then
    echo "🔴 **HIGH RISK**: $TOTAL_VULNS vulnerabilities found"
elif [[ $TOTAL_VULNS -gt 10 ]]; then
    echo "🟡 **MEDIUM RISK**: $TOTAL_VULNS vulnerabilities found"
else
    echo "🟠 **LOW RISK**: $TOTAL_VULNS vulnerabilities found"
fi)

### Vulnerability Breakdown
- **Total Vulnerabilities**: $TOTAL_VULNS
- **Critical Severity**: $CRITICAL_VULNS
- **High Severity**: $HIGH_VULNS
- **Medium Severity**: $MEDIUM_VULNS
- **Low Severity**: $LOW_VULNS

### Package Files Analyzed
$(cat "$RESULTS_DIR/package-files.txt" | sed 's/^/- /')

## License Compliance
$(if [[ -f "$RESULTS_DIR/license-scan.txt" ]]; then
    LICENSE_ISSUES=$(grep -c "non-compliant\|violation\|forbidden" "$RESULTS_DIR/license-scan.txt" 2>/dev/null || echo "0")
    if [[ $LICENSE_ISSUES -eq 0 ]]; then
        echo "✅ **COMPLIANT**: No license compliance issues detected"
    else
        echo "⚠️ **ISSUES FOUND**: $LICENSE_ISSUES license compliance issues"
    fi
else
    echo "ℹ️ License scanning completed (see detailed report)"
fi)

## Security Recommendations

### Immediate Actions
$(if [[ $CRITICAL_VULNS -gt 0 ]]; then
    echo "1. 🚨 **URGENT**: Address $CRITICAL_VULNS critical vulnerabilities immediately"
    echo "2. Review and update vulnerable dependencies"
    echo "3. Consider security patches and workarounds"
elif [[ $HIGH_VULNS -gt 0 ]]; then
    echo "1. 🔴 **HIGH PRIORITY**: Address $HIGH_VULNS high-severity vulnerabilities"
    echo "2. Plan dependency updates within 1-2 weeks"
    echo "3. Implement additional security controls if updates not available"
elif [[ $TOTAL_VULNS -gt 0 ]]; then
    echo "1. 📋 Review $TOTAL_VULNS vulnerabilities and prioritize fixes"
    echo "2. Update dependencies during next maintenance window"
    echo "3. Monitor for new vulnerabilities"
else
    echo "1. ✅ Continue current security practices"
    echo "2. Implement continuous vulnerability monitoring"
    echo "3. Regular dependency updates and security reviews"
fi)

### Long-term Security Strategy
1. **Automated Scanning**: Integrate OSV Scanner into CI/CD pipelines
2. **Dependency Management**: Implement automated dependency update workflows
3. **Security Monitoring**: Set up continuous vulnerability monitoring
4. **Supply Chain Security**: Implement SBOM generation and analysis
5. **License Governance**: Establish clear license compliance policies

## Generated Files
- \`vulnerability-scan.json\`: Detailed vulnerability scan results
- \`vulnerability-report.txt\`: Human-readable vulnerability report
- \`security-results.sarif\`: SARIF format for security tools integration
- \`license-scan.txt\`: License compliance analysis
- \`package-inventory.json\`: Complete package inventory
- \`package-files.txt\`: List of package management files found

## Next Steps
1. **Remediation Planning**: Prioritize fixes based on vulnerability severity
2. **Dependency Updates**: Plan and test dependency updates
3. **Security Testing**: Implement security testing in development workflows
4. **Continuous Monitoring**: Set up ongoing vulnerability monitoring
5. **Team Training**: Provide security awareness training for development teams

## Compliance and Audit
- **NIST Cybersecurity Framework**: Supports Identify and Protect functions
- **OWASP Top 10**: Addresses A06:2021 – Vulnerable and Outdated Components
- **SLSA Supply Chain**: Supports supply chain security requirements
- **ISO 27001**: Contributes to information security management
EOF

echo "Security assessment complete!"
echo "Results directory: $RESULTS_DIR/"
echo "Executive summary: $RESULTS_DIR/security-summary.md"
echo "Vulnerabilities found: $TOTAL_VULNS"

# Set exit code based on findings
if [[ $CRITICAL_VULNS -gt 0 ]]; then
    echo "⚠️ CRITICAL vulnerabilities found - immediate action required"
    exit 2
elif [[ $HIGH_VULNS -gt 5 ]]; then
    echo "⚠️ Multiple HIGH severity vulnerabilities - priority action required"
    exit 1
elif [[ $TOTAL_VULNS -gt 0 ]]; then
    echo "ℹ️ Vulnerabilities found - review and plan remediation"
    exit 0
else
    echo "✅ No vulnerabilities detected"
    exit 0
fi
```

### Guided Remediation Workflow
```bash
#!/bin/bash
# guided-remediation-workflow.sh

PROJECT_DIR="$1"
REMEDIATION_MODE="$2"  # interactive or automatic
MAX_DEPTH="${3:-3}"
MIN_SEVERITY="${4:-5}"

if [[ -z "$PROJECT_DIR" ]]; then
    echo "Usage: $0 <project-directory> [interactive|automatic] [max-depth] [min-severity]"
    echo "Example: $0 ./my-project interactive 3 7"
    exit 1
fi

REMEDIATION_MODE=${REMEDIATION_MODE:-interactive}
REMEDIATION_ID="remediation-$(date +%Y%m%d-%H%M%S)"
REMEDIATION_DIR="$PROJECT_DIR/security-remediation-$REMEDIATION_ID"

echo "Starting guided remediation workflow..."
echo "Project: $PROJECT_DIR"
echo "Mode: $REMEDIATION_MODE"
echo "Max dependency depth: $MAX_DEPTH"
echo "Min severity: $MIN_SEVERITY"

mkdir -p "$REMEDIATION_DIR"

# Phase 1: Initial vulnerability assessment
echo "Phase 1: Performing initial vulnerability scan..."
osv-scanner scan source -r "$PROJECT_DIR" \
    --format json \
    --output "$REMEDIATION_DIR/pre-fix-scan.json"

PRE_FIX_VULNS=$(grep -o '"id":"[^"]*"' "$REMEDIATION_DIR/pre-fix-scan.json" 2>/dev/null | wc -l || echo "0")
echo "Found $PRE_FIX_VULNS vulnerabilities before remediation"

if [[ $PRE_FIX_VULNS -eq 0 ]]; then
    echo "✅ No vulnerabilities found - no remediation needed"
    exit 0
fi

# Phase 2: Backup current state
echo "Phase 2: Creating backup of current state..."
BACKUP_DIR="$REMEDIATION_DIR/backup"
mkdir -p "$BACKUP_DIR"

# Backup package files
find "$PROJECT_DIR" -name "package*.json" -o -name "go.mod" -o -name "go.sum" \
    -o -name "requirements*.txt" -o -name "Cargo.toml" -o -name "Cargo.lock" \
    -o -name "pom.xml" -o -name "composer.json" -o -name "composer.lock" | \
while read -r file; do
    cp "$file" "$BACKUP_DIR/" 2>/dev/null || true
done

echo "Package files backed up to: $BACKUP_DIR/"

# Phase 3: Identify package management files
echo "Phase 3: Identifying package management files..."
MANIFEST_FILES=()
LOCKFILE_FILES=()

# Find manifest files
while IFS= read -r -d '' file; do
    MANIFEST_FILES+=("$file")
done < <(find "$PROJECT_DIR" -name "package.json" -o -name "go.mod" -o -name "pom.xml" \
    -o -name "composer.json" -o -name "Cargo.toml" -print0)

# Find lockfiles
while IFS= read -r -d '' file; do
    LOCKFILE_FILES+=("$file")
done < <(find "$PROJECT_DIR" -name "package-lock.json" -o -name "go.sum" -o -name "Cargo.lock" \
    -o -name "composer.lock" -print0)

echo "Found ${#MANIFEST_FILES[@]} manifest files and ${#LOCKFILE_FILES[@]} lockfiles"

# Phase 4: Apply guided remediation
echo "Phase 4: Applying guided remediation..."

for i in "${!MANIFEST_FILES[@]}"; do
    MANIFEST="${MANIFEST_FILES[$i]}"
    
    # Find corresponding lockfile
    LOCKFILE=""
    MANIFEST_DIR=$(dirname "$MANIFEST")
    MANIFEST_NAME=$(basename "$MANIFEST")
    
    case "$MANIFEST_NAME" in
        "package.json")
            LOCKFILE="$MANIFEST_DIR/package-lock.json"
            ;;
        "go.mod")
            LOCKFILE="$MANIFEST_DIR/go.sum"
            ;;
        "composer.json")
            LOCKFILE="$MANIFEST_DIR/composer.lock"
            ;;
        "Cargo.toml")
            LOCKFILE="$MANIFEST_DIR/Cargo.lock"
            ;;
    esac
    
    if [[ -f "$LOCKFILE" ]]; then
        echo "Processing: $MANIFEST with $LOCKFILE"
        
        if [[ "$REMEDIATION_MODE" == "interactive" ]]; then
            # Interactive mode
            echo "Running interactive remediation for $MANIFEST..."
            osv-scanner fix \
                -M "$MANIFEST" \
                -L "$LOCKFILE" \
                2>&1 | tee "$REMEDIATION_DIR/fix-output-$(basename "$MANIFEST").log"
        else
            # Automatic mode
            echo "Running automatic remediation for $MANIFEST..."
            osv-scanner fix \
                --strategy in-place \
                --max-depth "$MAX_DEPTH" \
                --min-severity "$MIN_SEVERITY" \
                --ignore-dev \
                -M "$MANIFEST" \
                -L "$LOCKFILE" \
                2>&1 | tee "$REMEDIATION_DIR/fix-output-$(basename "$MANIFEST").log"
        fi
    else
        echo "⚠️ No lockfile found for $MANIFEST - skipping guided remediation"
    fi
done

# Phase 5: Post-remediation assessment
echo "Phase 5: Performing post-remediation vulnerability scan..."
osv-scanner scan source -r "$PROJECT_DIR" \
    --format json \
    --output "$REMEDIATION_DIR/post-fix-scan.json"

POST_FIX_VULNS=$(grep -o '"id":"[^"]*"' "$REMEDIATION_DIR/post-fix-scan.json" 2>/dev/null | wc -l || echo "0")
VULNS_FIXED=$((PRE_FIX_VULNS - POST_FIX_VULNS))

echo "Vulnerabilities before remediation: $PRE_FIX_VULNS"
echo "Vulnerabilities after remediation: $POST_FIX_VULNS"
echo "Vulnerabilities fixed: $VULNS_FIXED"

# Phase 6: Generate remediation report
echo "Phase 6: Generating remediation report..."
cat > "$REMEDIATION_DIR/remediation-report.md" <<EOF
# Security Remediation Report

**Remediation ID**: $REMEDIATION_ID  
**Project**: $PROJECT_DIR  
**Date**: $(date)  
**Mode**: $REMEDIATION_MODE  
**Tool**: OSV Scanner Guided Remediation

## Remediation Summary

### Results Overview
- **Pre-Remediation Vulnerabilities**: $PRE_FIX_VULNS
- **Post-Remediation Vulnerabilities**: $POST_FIX_VULNS
- **Vulnerabilities Fixed**: $VULNS_FIXED
- **Fix Success Rate**: $(( PRE_FIX_VULNS > 0 ? VULNS_FIXED * 100 / PRE_FIX_VULNS : 0 ))%

### Configuration
- **Max Dependency Depth**: $MAX_DEPTH
- **Min Severity**: $MIN_SEVERITY
- **Strategy**: $(if [[ "$REMEDIATION_MODE" == "interactive" ]]; then echo "Interactive user selection"; else echo "Automatic in-place updates"; fi)

### Files Processed
$(for file in "${MANIFEST_FILES[@]}"; do echo "- $file"; done)

### Backup Information
Original package files backed up to: \`backup/\`

## Remediation Status
$(if [[ $VULNS_FIXED -gt 0 ]]; then
    echo "✅ **SUCCESS**: $VULNS_FIXED vulnerabilities successfully remediated"
elif [[ $POST_FIX_VULNS -lt $PRE_FIX_VULNS ]]; then
    echo "🟡 **PARTIAL**: Some vulnerabilities remediated ($VULNS_FIXED fixed)"
else
    echo "🔴 **LIMITED**: No vulnerabilities automatically fixable"
fi)

## Next Steps

### Immediate Actions
$(if [[ $POST_FIX_VULNS -eq 0 ]]; then
    echo "1. ✅ All vulnerabilities resolved"
    echo "2. Test application functionality thoroughly"
    echo "3. Deploy changes after validation"
    echo "4. Update security documentation"
elif [[ $VULNS_FIXED -gt 0 ]]; then
    echo "1. 🧪 Test application with updated dependencies"
    echo "2. Address remaining $POST_FIX_VULNS vulnerabilities manually"
    echo "3. Review fix logs for any breaking changes"
    echo "4. Consider alternative remediation strategies"
else
    echo "1. 📋 Review vulnerabilities requiring manual intervention"
    echo "2. Check for security patches or workarounds"
    echo "3. Consider upgrading to newer versions"
    echo "4. Implement additional security controls"
fi)

### Long-term Strategy
1. **Continuous Monitoring**: Set up automated vulnerability scanning
2. **Dependency Management**: Implement regular dependency updates
3. **Security Testing**: Add security tests to CI/CD pipeline
4. **Team Training**: Provide secure coding and dependency management training

## Files Generated
- \`pre-fix-scan.json\`: Vulnerability scan before remediation
- \`post-fix-scan.json\`: Vulnerability scan after remediation
- \`backup/\`: Backup of original package files
- \`fix-output-*.log\`: Detailed remediation logs for each package file

## Rollback Instructions
If issues are detected after remediation:
1. Stop the application
2. Restore original files: \`cp backup/* /path/to/original/locations\`
3. Reinstall dependencies: \`npm install\` / \`go mod download\` / etc.
4. Restart application and verify functionality

## Validation Checklist
- [ ] Application builds successfully
- [ ] All tests pass
- [ ] Functionality verification complete
- [ ] Performance impact assessed
- [ ] Security improvements validated
- [ ] Documentation updated
EOF

echo "Guided remediation complete!"
echo "Report: $REMEDIATION_DIR/remediation-report.md"
echo "Backup: $REMEDIATION_DIR/backup/"

if [[ $VULNS_FIXED -gt 0 ]]; then
    echo "✅ $VULNS_FIXED vulnerabilities fixed - please test your application"
    exit 0
elif [[ $POST_FIX_VULNS -gt 0 ]]; then
    echo "⚠️ $POST_FIX_VULNS vulnerabilities remain - manual intervention required"
    exit 1
else
    echo "ℹ️ No changes made - vulnerabilities may require manual remediation"
    exit 0
fi
```

### Continuous Integration Security Pipeline
```bash
#!/bin/bash
# ci-security-pipeline.sh

PROJECT_ROOT="$1"
BUILD_ID="$2"
SECURITY_THRESHOLD="$3"

PROJECT_ROOT=${PROJECT_ROOT:-.}
BUILD_ID=${BUILD_ID:-ci-$(date +%Y%m%d-%H%M%S)}
SECURITY_THRESHOLD=${SECURITY_THRESHOLD:-10}

echo "Starting CI/CD security pipeline..."
echo "Project: $PROJECT_ROOT"
echo "Build ID: $BUILD_ID"
echo "Security threshold: $SECURITY_THRESHOLD vulnerabilities"

CI_RESULTS_DIR="ci-security-$BUILD_ID"
mkdir -p "$CI_RESULTS_DIR"

# Phase 1: Pre-build security scan
echo "Phase 1: Pre-build security scanning..."
osv-scanner scan source -r "$PROJECT_ROOT" \
    --format json \
    --output "$CI_RESULTS_DIR/security-scan.json" \
    > "$CI_RESULTS_DIR/scan-output.log" 2>&1

SCAN_EXIT_CODE=$?

# Phase 2: Parse security results
echo "Phase 2: Analyzing security findings..."
VULNERABILITY_COUNT=0
CRITICAL_COUNT=0
HIGH_COUNT=0

if [[ -f "$CI_RESULTS_DIR/security-scan.json" ]]; then
    if command -v jq > /dev/null; then
        VULNERABILITY_COUNT=$(jq -r '.results[]?.packages[]?.vulnerabilities? | length' \
            "$CI_RESULTS_DIR/security-scan.json" 2>/dev/null | \
            awk '{sum += $1} END {print sum+0}')
        
        CRITICAL_COUNT=$(jq -r '.results[]?.packages[]?.vulnerabilities[]? | 
            select(.database_specific?.severity == "CRITICAL") | .id' \
            "$CI_RESULTS_DIR/security-scan.json" 2>/dev/null | wc -l)
        
        HIGH_COUNT=$(jq -r '.results[]?.packages[]?.vulnerabilities[]? | 
            select(.database_specific?.severity == "HIGH") | .id' \
            "$CI_RESULTS_DIR/security-scan.json" 2>/dev/null | wc -l)
    else
        VULNERABILITY_COUNT=$(grep -o '"id":"[^"]*"' "$CI_RESULTS_DIR/security-scan.json" 2>/dev/null | wc -l || echo "0")
    fi
fi

echo "Security scan results:"
echo "- Total vulnerabilities: $VULNERABILITY_COUNT"
echo "- Critical severity: $CRITICAL_COUNT"
echo "- High severity: $HIGH_COUNT"

# Phase 3: License compliance check
echo "Phase 3: License compliance scanning..."
osv-scanner --licenses "$PROJECT_ROOT" \
    --output "$CI_RESULTS_DIR/license-compliance.txt" \
    > "$CI_RESULTS_DIR/license-output.log" 2>&1

# Phase 4: Generate CI/CD artifacts
echo "Phase 4: Generating CI/CD security artifacts..."

# SARIF for security tools integration
osv-scanner scan source -r "$PROJECT_ROOT" \
    --format sarif \
    --output "$CI_RESULTS_DIR/security-results.sarif" \
    2>/dev/null || true

# Create security gate status
SECURITY_GATE_STATUS="PASS"
SECURITY_GATE_MESSAGE="Security scan passed"

if [[ $CRITICAL_COUNT -gt 0 ]]; then
    SECURITY_GATE_STATUS="FAIL"
    SECURITY_GATE_MESSAGE="CRITICAL vulnerabilities detected ($CRITICAL_COUNT)"
elif [[ $HIGH_COUNT -gt 5 ]]; then
    SECURITY_GATE_STATUS="FAIL"
    SECURITY_GATE_MESSAGE="Too many HIGH severity vulnerabilities ($HIGH_COUNT)"
elif [[ $VULNERABILITY_COUNT -gt $SECURITY_THRESHOLD ]]; then
    SECURITY_GATE_STATUS="FAIL"
    SECURITY_GATE_MESSAGE="Vulnerability count exceeds threshold ($VULNERABILITY_COUNT > $SECURITY_THRESHOLD)"
fi

# Phase 5: Generate CI/CD security report
echo "Phase 5: Creating CI/CD security report..."
cat > "$CI_RESULTS_DIR/ci-security-report.json" <<EOF
{
    "build_id": "$BUILD_ID",
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "project_root": "$PROJECT_ROOT",
    "security_gate": {
        "status": "$SECURITY_GATE_STATUS",
        "message": "$SECURITY_GATE_MESSAGE",
        "threshold": $SECURITY_THRESHOLD
    },
    "vulnerability_summary": {
        "total_vulnerabilities": $VULNERABILITY_COUNT,
        "critical_count": $CRITICAL_COUNT,
        "high_count": $HIGH_COUNT,
        "scan_exit_code": $SCAN_EXIT_CODE
    },
    "artifacts": {
        "security_scan": "security-scan.json",
        "sarif_results": "security-results.sarif",
        "license_compliance": "license-compliance.txt",
        "scan_logs": "scan-output.log"
    }
}
EOF

# Phase 6: Security gate decision
echo "Phase 6: Security gate evaluation..."
echo "Security Gate Status: $SECURITY_GATE_STATUS"
echo "Message: $SECURITY_GATE_MESSAGE"

# Generate summary for CI/CD systems
cat > "$CI_RESULTS_DIR/security-summary.txt" <<EOF
Security Gate: $SECURITY_GATE_STATUS
Vulnerabilities: $VULNERABILITY_COUNT (Critical: $CRITICAL_COUNT, High: $HIGH_COUNT)
Threshold: $SECURITY_THRESHOLD
Message: $SECURITY_GATE_MESSAGE
Build ID: $BUILD_ID
EOF

echo "CI/CD security pipeline complete!"
echo "Results: $CI_RESULTS_DIR/"
echo "Security report: $CI_RESULTS_DIR/ci-security-report.json"

# Exit with appropriate code for CI/CD systems
case "$SECURITY_GATE_STATUS" in
    "PASS")
        echo "✅ Security gate PASSED - build can proceed"
        exit 0
        ;;
    "FAIL")
        echo "❌ Security gate FAILED - build should be blocked"
        exit 1
        ;;
    *)
        echo "⚠️ Security gate status unknown - manual review required"
        exit 2
        ;;
esac
```

## Integration Patterns

### GitHub Actions Workflow
```yaml
# .github/workflows/security-scanning.yml
name: Security Vulnerability Scanning
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 2 * * 1'  # Weekly Monday 2 AM

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Install OSV Scanner
      run: |
        # Install OSV Scanner
        curl -L https://github.com/google/osv-scanner/releases/latest/download/osv-scanner_linux_amd64 \
          -o osv-scanner
        chmod +x osv-scanner
        sudo mv osv-scanner /usr/local/bin/
        
    - name: Verify OSV Scanner Installation
      run: |
        osv-scanner --help
        
    - name: Run Vulnerability Scan
      run: |
        echo "Scanning for vulnerabilities..."
        osv-scanner scan source -r . \
          --format json \
          --output vulnerability-results.json || true
          
        # Also generate SARIF for GitHub security tab
        osv-scanner scan source -r . \
          --format sarif \
          --output vulnerability-results.sarif || true
          
    - name: Run License Compliance Scan
      run: |
        echo "Scanning for license compliance..."
        osv-scanner --licenses . \
          --output license-results.txt || true
          
    - name: Parse Security Results
      id: security-results
      run: |
        if [[ -f vulnerability-results.json ]]; then
          if command -v jq > /dev/null; then
            VULN_COUNT=$(jq -r '.results[]?.packages[]?.vulnerabilities? | length' \
              vulnerability-results.json 2>/dev/null | \
              awk '{sum += $1} END {print sum+0}')
            
            CRITICAL_COUNT=$(jq -r '.results[]?.packages[]?.vulnerabilities[]? | 
              select(.database_specific?.severity == "CRITICAL") | .id' \
              vulnerability-results.json 2>/dev/null | wc -l)
          else
            VULN_COUNT=$(grep -c '"id":' vulnerability-results.json || echo "0")
            CRITICAL_COUNT=0
          fi
        else
          VULN_COUNT=0
          CRITICAL_COUNT=0
        fi
        
        echo "vulnerability_count=$VULN_COUNT" >> $GITHUB_OUTPUT
        echo "critical_count=$CRITICAL_COUNT" >> $GITHUB_OUTPUT
        
        echo "Found $VULN_COUNT vulnerabilities ($CRITICAL_COUNT critical)"
        
    - name: Upload SARIF Results
      if: always()
      uses: github/codeql-action/upload-sarif@v3
      with:
        sarif_file: vulnerability-results.sarif
        
    - name: Upload Security Artifacts
      if: always()
      uses: actions/upload-artifact@v4
      with:
        name: security-scan-results
        path: |
          vulnerability-results.json
          vulnerability-results.sarif
          license-results.txt
          
    - name: Security Gate
      run: |
        VULN_COUNT=${{ steps.security-results.outputs.vulnerability_count }}
        CRITICAL_COUNT=${{ steps.security-results.outputs.critical_count }}
        
        echo "Security Gate Evaluation:"
        echo "- Vulnerabilities: $VULN_COUNT"
        echo "- Critical: $CRITICAL_COUNT"
        
        if [[ $CRITICAL_COUNT -gt 0 ]]; then
          echo "❌ FAIL: Critical vulnerabilities detected"
          echo "::error::Found $CRITICAL_COUNT critical vulnerabilities"
          exit 1
        elif [[ $VULN_COUNT -gt 10 ]]; then
          echo "⚠️ WARNING: High vulnerability count ($VULN_COUNT)"
          echo "::warning::Consider addressing vulnerabilities before deployment"
        else
          echo "✅ PASS: Security scan acceptable"
        fi
```

### Docker Integration
```dockerfile
# Dockerfile.osv-scanner
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Install OSV Scanner
RUN curl -L https://github.com/google/osv-scanner/releases/latest/download/osv-scanner_linux_amd64 \
    -o /usr/local/bin/osv-scanner \
    && chmod +x /usr/local/bin/osv-scanner

# Verify installation
RUN osv-scanner --help

# Copy scanning scripts
COPY security-scan.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/security-scan.sh

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/security-scan.sh"]
```

### Terraform Security Module
```hcl
# terraform/security-scanning.tf
resource "null_resource" "osv_security_scan" {
  count = var.enable_security_scanning ? 1 : 0
  
  triggers = {
    source_hash = filemd5("${path.module}/../package.json")
  }

  provisioner "local-exec" {
    command = <<-EOF
      echo "Running OSV security scan..."
      
      # Install OSV Scanner if not present
      if ! command -v osv-scanner > /dev/null; then
        curl -L https://github.com/google/osv-scanner/releases/latest/download/osv-scanner_linux_amd64 \
          -o osv-scanner
        chmod +x osv-scanner
        sudo mv osv-scanner /usr/local/bin/
      fi
      
      # Run security scan
      osv-scanner scan source -r ${path.module}/.. \
        --format json \
        --output terraform-security-scan.json
      
      # Check for critical vulnerabilities
      if command -v jq > /dev/null; then
        CRITICAL_COUNT=$(jq -r '.results[]?.packages[]?.vulnerabilities[]? | 
          select(.database_specific?.severity == "CRITICAL") | .id' \
          terraform-security-scan.json 2>/dev/null | wc -l || echo "0")
        
        if [[ $CRITICAL_COUNT -gt 0 ]]; then
          echo "ERROR: Found $CRITICAL_COUNT critical vulnerabilities"
          exit 1
        fi
      fi
      
      echo "Security scan completed successfully"
    EOF
  }
}

variable "enable_security_scanning" {
  description = "Enable security vulnerability scanning"
  type        = bool
  default     = true
}

output "security_scan_results" {
  value = var.enable_security_scanning ? "terraform-security-scan.json" : "Security scanning disabled"
}
```

## Best Practices

### Scanning Strategy
- **Comprehensive Coverage**: Scan source code, containers, and SBOMs for complete visibility
- **Multi-Format Output**: Use JSON for automation, SARIF for security tools, table for humans
- **Regular Scheduling**: Implement continuous scanning in CI/CD and scheduled workflows
- **Ecosystem Support**: Leverage OSV Scanner's broad ecosystem coverage

### Integration Approach
- **CI/CD Integration**: Implement security gates with appropriate thresholds
- **Pre-commit Hooks**: Catch vulnerabilities before they enter version control
- **IDE Integration**: Provide real-time vulnerability feedback to developers
- **Container Scanning**: Include container images in security assessment workflows

### Remediation Management
- **Guided Remediation**: Use OSV Scanner's fix capabilities for automated updates
- **Risk-Based Prioritization**: Focus on critical and high-severity vulnerabilities first
- **Testing Requirements**: Thoroughly test applications after dependency updates
- **Rollback Procedures**: Maintain ability to quickly rollback problematic updates

### Reporting and Governance
- **Executive Dashboards**: Provide high-level security metrics for management
- **Developer Feedback**: Ensure actionable feedback reaches development teams
- **Compliance Reporting**: Generate reports for audit and regulatory requirements
- **Trend Analysis**: Track security improvements over time

## Error Handling

### Common Issues
```bash
# Network connectivity issues
osv-scanner --offline --download-offline-databases ./offline-db
# Solution: Use offline mode for environments without internet access

# Permission errors
sudo osv-scanner scan source -r /protected/directory
# Solution: Ensure proper file system permissions

# Large project timeouts
osv-scanner --verbosity error scan source -r ./large-project
# Solution: Reduce verbosity and consider scanning subdirectories separately

# Container scanning failures
docker pull image:tag && osv-scanner scan image image:tag
# Solution: Ensure image is available locally before scanning
```

### Troubleshooting
- **Scan Failures**: Check network connectivity and file permissions
- **Performance Issues**: Use offline mode or reduce scan scope for large projects
- **Integration Problems**: Verify OSV Scanner version compatibility with CI/CD tools
- **False Positives**: Review vulnerability context and consider suppression for confirmed false positives

OSV Scanner provides comprehensive open source vulnerability detection and remediation capabilities, enabling organizations to maintain secure software supply chains through automated scanning, guided remediation, and seamless integration with modern development workflows and security toolchains.