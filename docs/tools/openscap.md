# OpenSCAP

Open-source security compliance scanning framework for comprehensive SCAP content evaluation and system compliance assessment.

## Description

OpenSCAP is a comprehensive open-source implementation of the Security Content Automation Protocol (SCAP) specification, providing robust security compliance scanning, vulnerability assessment, and configuration management capabilities. Developed as a NIST-certified SCAP 1.3 and 1.2 validated tool, OpenSCAP enables organizations to evaluate systems against industry-standard security benchmarks including CIS (Center for Internet Security), DISA STIG (Defense Information Systems Agency Security Technical Implementation Guides), NIST 800-53, PCI DSS, and other compliance frameworks. The tool supports XCCDF (Extensible Configuration Checklist Description Format) for compliance checking, OVAL (Open Vulnerability and Assessment Language) for detailed system state evaluation, and CPE (Common Platform Enumeration) for platform identification, making it essential for automated security compliance, vulnerability management, and regulatory compliance initiatives across enterprise environments.

## MCP Tools

### XCCDF Content Evaluation
- **`openscap_xccdf_eval`** - Evaluate XCCDF content for security compliance using real oscap CLI
- **`openscap_xccdf_generate_report`** - Generate HTML report from XCCDF results using real oscap CLI
- **`openscap_xccdf_generate_guide`** - Generate HTML guide from XCCDF content using real oscap CLI
- **`openscap_xccdf_remediate`** - Apply remediation based on XCCDF results using real oscap CLI

### OVAL Definitions Processing
- **`openscap_oval_eval`** - Evaluate OVAL definitions using real oscap CLI
- **`openscap_oval_generate_report`** - Generate report from OVAL results using real oscap CLI

### DataStream Management
- **`openscap_ds_validate`** - Validate Source DataStream file using real oscap CLI
- **`openscap_ds_split`** - Split DataStream into component files using real oscap CLI

### Content Validation and Information
- **`openscap_validate`** - Validate SCAP content (XCCDF, OVAL, CPE, CVE) using real oscap CLI
- **`openscap_info`** - Display information about SCAP content using real oscap CLI

## Real CLI Commands Used

### XCCDF Commands
- `oscap xccdf eval [options] <xccdf_file>` - Evaluate XCCDF content
- `oscap xccdf generate report <results_file>` - Generate HTML report from results
- `oscap xccdf generate guide [--profile <profile>] <xccdf_file>` - Generate HTML guide
- `oscap xccdf remediate <results_file>` - Apply remediation based on results

### OVAL Commands
- `oscap oval eval [options] <oval_file>` - Evaluate OVAL definitions
- `oscap oval generate report <oval_results_file>` - Generate OVAL report

### DataStream Commands
- `oscap ds sds-validate <datastream_file>` - Validate Source DataStream
- `oscap ds sds-split [--output-dir <dir>] <datastream_file>` - Split DataStream

### Content Information and Validation
- `oscap info <scap_file>` - Display content information
- `oscap [xccdf|oval|cpe|cve] validate <content_file>` - Validate specific content types

### Common XCCDF Evaluation Options
- `--profile <profile_id>` - Specify security profile to evaluate
- `--results <results_file>` - Save evaluation results to XML file
- `--report <report_file>` - Generate HTML report during evaluation
- `--cpe <cpe_file>` - Use CPE dictionary for platform identification
- `--fetch-remote-resources` - Download remote resources during evaluation
- `--oval-results` - Include OVAL results in evaluation

### Common OVAL Evaluation Options
- `--results <results_file>` - Save OVAL evaluation results
- `--variables <variables_file>` - Use external OVAL variables
- `--id <definition_id>` - Evaluate specific OVAL definition

## Use Cases

### Security Compliance Assessment
- **Benchmark Evaluation**: Assess systems against CIS, STIG, NIST, and other security benchmarks
- **Vulnerability Scanning**: Identify security vulnerabilities using OVAL definitions
- **Configuration Compliance**: Verify system configurations meet security requirements
- **Regulatory Compliance**: Support SOX, HIPAA, PCI DSS, and other regulatory frameworks

### Enterprise Security Management
- **Automated Auditing**: Implement continuous compliance monitoring
- **Risk Assessment**: Identify and prioritize security risks across infrastructure
- **Remediation Management**: Generate and apply automated security fixes
- **Compliance Reporting**: Produce detailed compliance reports for stakeholders

### DevSecOps Integration
- **CI/CD Security Gates**: Integrate security compliance checks in deployment pipelines
- **Infrastructure as Code**: Validate security configurations before deployment
- **Container Security**: Assess container images and runtime configurations
- **Cloud Security**: Evaluate cloud infrastructure compliance

### Government and Defense
- **STIG Compliance**: Ensure compliance with Defense Information Systems Agency guidelines
- **FedRAMP Assessment**: Support Federal Risk and Authorization Management Program requirements
- **FISMA Compliance**: Meet Federal Information Security Management Act requirements
- **Common Criteria**: Support Common Criteria evaluation and certification processes

## Configuration Examples

### Basic XCCDF Evaluation
```bash
# Evaluate system against OSPP (Operating System Protection Profile)
oscap xccdf eval --profile ospp --results results.xml --report report.html datastream.xml

# Evaluate specific CIS benchmark
oscap xccdf eval --profile xccdf_org.cisecurity.benchmarks_profile_Level_1_-_Server \
  --results cis-results.xml \
  --report cis-report.html \
  cis-benchmark.xml

# Evaluate DISA STIG benchmark
oscap xccdf eval --profile stig-rhel8-server \
  --results stig-results.xml \
  --report stig-report.html \
  --fetch-remote-resources \
  stig-datastream.xml

# Generate compliance guide
oscap xccdf generate guide --profile ospp datastream.xml > compliance-guide.html

# Apply remediation based on evaluation results
oscap xccdf remediate results.xml
```

### OVAL Definition Evaluation
```bash
# Evaluate OVAL definitions for vulnerability assessment
oscap oval eval --results oval-results.xml vulnerability-definitions.xml

# Evaluate specific OVAL definition
oscap oval eval --id oval:com.example:def:123 \
  --results specific-oval-results.xml \
  definitions.xml

# Generate OVAL evaluation report
oscap oval generate report oval-results.xml > oval-report.html

# Evaluate OVAL with external variables
oscap oval eval --variables oval-variables.xml \
  --results oval-results.xml \
  definitions.xml
```

### DataStream Management
```bash
# Validate Source DataStream integrity
oscap ds sds-validate scap-datastream.xml

# Split DataStream into component files
oscap ds sds-split --output-dir ./components scap-datastream.xml

# Get information about SCAP content
oscap info scap-content.xml

# Validate XCCDF content
oscap xccdf validate benchmark.xml

# Validate OVAL content with Schematron
oscap oval validate --schematron definitions.xml
```

### Content Validation and Information
```bash
# Display comprehensive information about SCAP content
oscap info datastream.xml

# Validate different content types
oscap xccdf validate xccdf-benchmark.xml
oscap oval validate oval-definitions.xml
oscap cpe validate cpe-dictionary.xml
oscap cve validate cve-feed.xml

# Get detailed content metadata
oscap info --fetch-remote-resources datastream.xml
```

## Advanced Usage

### Comprehensive Security Compliance Assessment
```bash
#!/bin/bash
# comprehensive-compliance-assessment.sh

TARGET_SYSTEM="$1"
BENCHMARK_TYPE="$2"
ASSESSMENT_ID="assessment-$(date +%Y%m%d-%H%M%S)"
RESULTS_DIR="compliance-results-$ASSESSMENT_ID"

if [[ -z "$TARGET_SYSTEM" || -z "$BENCHMARK_TYPE" ]]; then
    echo "Usage: $0 <target-system> <benchmark-type>"
    echo "Benchmark types: cis, stig, ospp, pci-dss"
    exit 1
fi

echo "Starting comprehensive compliance assessment for: $TARGET_SYSTEM"
echo "Benchmark: $BENCHMARK_TYPE"
mkdir -p "$RESULTS_DIR"

# Phase 1: Content information and validation
echo "Phase 1: Validating SCAP content..."
case "$BENCHMARK_TYPE" in
    "cis")
        CONTENT_FILE="cis-$TARGET_SYSTEM-benchmark.xml"
        PROFILE_ID="xccdf_org.cisecurity.benchmarks_profile_Level_1_-_Server"
        ;;
    "stig")
        CONTENT_FILE="stig-$TARGET_SYSTEM-datastream.xml"
        PROFILE_ID="stig-$TARGET_SYSTEM-server"
        ;;
    "ospp")
        CONTENT_FILE="ospp-$TARGET_SYSTEM-datastream.xml"
        PROFILE_ID="ospp"
        ;;
    "pci-dss")
        CONTENT_FILE="pci-dss-$TARGET_SYSTEM-datastream.xml"
        PROFILE_ID="pci-dss"
        ;;
    *)
        echo "Unknown benchmark type: $BENCHMARK_TYPE"
        exit 1
        ;;
esac

# Validate content integrity
oscap info "$CONTENT_FILE" > "$RESULTS_DIR/content-info.txt"
oscap ds sds-validate "$CONTENT_FILE" > "$RESULTS_DIR/content-validation.txt" 2>&1

# Phase 2: System evaluation
echo "Phase 2: Performing security compliance evaluation..."
oscap xccdf eval \
    --profile "$PROFILE_ID" \
    --results "$RESULTS_DIR/compliance-results.xml" \
    --report "$RESULTS_DIR/compliance-report.html" \
    --oval-results \
    --fetch-remote-resources \
    "$CONTENT_FILE" > "$RESULTS_DIR/evaluation-output.txt" 2>&1

EVAL_EXIT_CODE=$?

# Phase 3: Generate comprehensive reports
echo "Phase 3: Generating compliance documentation..."
oscap xccdf generate guide \
    --profile "$PROFILE_ID" \
    "$CONTENT_FILE" > "$RESULTS_DIR/compliance-guide.html"

# Generate remediation scripts if evaluation found issues
if [[ $EVAL_EXIT_CODE -ne 0 ]]; then
    echo "Phase 4: Generating remediation guidance..."
    oscap xccdf generate fix \
        --profile "$PROFILE_ID" \
        --fix-type bash \
        "$RESULTS_DIR/compliance-results.xml" > "$RESULTS_DIR/remediation-script.sh"
    
    chmod +x "$RESULTS_DIR/remediation-script.sh"
fi

# Phase 5: Parse and summarize results
echo "Phase 5: Analyzing compliance results..."
PASS_COUNT=$(grep -c 'result="pass"' "$RESULTS_DIR/compliance-results.xml" 2>/dev/null || echo "0")
FAIL_COUNT=$(grep -c 'result="fail"' "$RESULTS_DIR/compliance-results.xml" 2>/dev/null || echo "0")
ERROR_COUNT=$(grep -c 'result="error"' "$RESULTS_DIR/compliance-results.xml" 2>/dev/null || echo "0")
UNKNOWN_COUNT=$(grep -c 'result="unknown"' "$RESULTS_DIR/compliance-results.xml" 2>/dev/null || echo "0")
TOTAL_RULES=$((PASS_COUNT + FAIL_COUNT + ERROR_COUNT + UNKNOWN_COUNT))

COMPLIANCE_PERCENTAGE=0
if [[ $TOTAL_RULES -gt 0 ]]; then
    COMPLIANCE_PERCENTAGE=$((PASS_COUNT * 100 / TOTAL_RULES))
fi

# Phase 6: Generate executive summary
echo "Phase 6: Creating executive summary..."
cat > "$RESULTS_DIR/compliance-summary.md" <<EOF
# Security Compliance Assessment Summary

**Assessment ID**: $ASSESSMENT_ID  
**Target System**: $TARGET_SYSTEM  
**Benchmark**: $BENCHMARK_TYPE  
**Profile**: $PROFILE_ID  
**Date**: $(date)  

## Compliance Results

### Overall Compliance Score
**${COMPLIANCE_PERCENTAGE}%** compliant ($PASS_COUNT/$TOTAL_RULES rules passed)

### Detailed Results
- **Passed**: $PASS_COUNT rules
- **Failed**: $FAIL_COUNT rules
- **Errors**: $ERROR_COUNT rules  
- **Unknown**: $UNKNOWN_COUNT rules
- **Total Evaluated**: $TOTAL_RULES rules

### Risk Assessment
$(if [[ $FAIL_COUNT -gt 20 ]]; then
    echo "🔴 **HIGH RISK**: Significant compliance failures detected"
elif [[ $FAIL_COUNT -gt 5 ]]; then
    echo "🟡 **MEDIUM RISK**: Moderate compliance issues identified"
else
    echo "🟢 **LOW RISK**: Minor or no compliance issues"
fi)

## Files Generated
- \`compliance-results.xml\`: Detailed XCCDF evaluation results
- \`compliance-report.html\`: Interactive HTML compliance report
- \`compliance-guide.html\`: Security configuration guide
- \`content-info.txt\`: SCAP content information and metadata
$(if [[ $EVAL_EXIT_CODE -ne 0 ]]; then
    echo "- \`remediation-script.sh\`: Automated remediation script"
fi)

## Recommendations

### Immediate Actions
$(if [[ $FAIL_COUNT -gt 0 ]]; then
    echo "1. Review failed compliance checks in the detailed HTML report"
    echo "2. Prioritize high-severity security configuration issues"
    echo "3. Test and apply automated remediation scripts in non-production environment"
else
    echo "1. System demonstrates good security compliance posture"
    echo "2. Continue regular compliance monitoring and assessment"
fi)

### Long-term Improvements
1. Implement automated compliance monitoring
2. Integrate security assessments into CI/CD pipelines
3. Establish regular compliance review cycles
4. Provide security awareness training to system administrators

## Next Steps
1. **Technical Review**: System administrators should review detailed findings
2. **Risk Evaluation**: Security team should assess business impact of failures
3. **Remediation Planning**: Develop timeline for addressing compliance gaps
4. **Validation Testing**: Re-run assessment after implementing fixes
5. **Documentation Update**: Update security documentation and procedures

## Compliance Framework Information
- **Framework**: $BENCHMARK_TYPE
- **Standard**: $(echo $PROFILE_ID | sed 's/_/ /g')
- **Assessment Method**: NIST SCAP 1.3 compliant evaluation
- **Tool Version**: $(oscap --version 2>/dev/null | head -1 || echo "OpenSCAP CLI")
EOF

echo "Compliance assessment complete!"
echo "Results directory: $RESULTS_DIR/"
echo "Executive summary: $RESULTS_DIR/compliance-summary.md"
echo "Detailed report: $RESULTS_DIR/compliance-report.html"
echo "Compliance score: ${COMPLIANCE_PERCENTAGE}% ($PASS_COUNT/$TOTAL_RULES)"

if [[ $COMPLIANCE_PERCENTAGE -lt 80 ]]; then
    echo "⚠️ WARNING: Compliance score below 80% - immediate attention required"
    exit 1
elif [[ $COMPLIANCE_PERCENTAGE -lt 95 ]]; then
    echo "✅ ACCEPTABLE: Good compliance posture with room for improvement"
    exit 0
else
    echo "🏆 EXCELLENT: Outstanding compliance posture"
    exit 0
fi
```

### Automated Remediation Workflow
```bash
#!/bin/bash
# automated-remediation-workflow.sh

RESULTS_FILE="$1"
REMEDIATION_TYPE="$2"
DRY_RUN="$3"

if [[ -z "$RESULTS_FILE" ]]; then
    echo "Usage: $0 <xccdf-results-file> [remediation-type] [dry-run]"
    echo "Remediation types: bash, ansible, puppet, chef"
    echo "Dry run: true/false (default: true)"
    exit 1
fi

REMEDIATION_TYPE=${REMEDIATION_TYPE:-bash}
DRY_RUN=${DRY_RUN:-true}
REMEDIATION_ID="remediation-$(date +%Y%m%d-%H%M%S)"
REMEDIATION_DIR="remediation-$REMEDIATION_ID"

echo "Starting automated remediation workflow..."
echo "Results file: $RESULTS_FILE"
echo "Remediation type: $REMEDIATION_TYPE"
echo "Dry run mode: $DRY_RUN"

mkdir -p "$REMEDIATION_DIR"

# Phase 1: Analyze compliance results
echo "Phase 1: Analyzing compliance failures..."
FAILED_RULES=$(grep 'result="fail"' "$RESULTS_FILE" | wc -l)
echo "Found $FAILED_RULES failed compliance rules"

if [[ $FAILED_RULES -eq 0 ]]; then
    echo "✅ No compliance failures found - no remediation needed"
    exit 0
fi

# Phase 2: Generate remediation content
echo "Phase 2: Generating remediation scripts..."
case "$REMEDIATION_TYPE" in
    "bash")
        oscap xccdf generate fix \
            --fix-type bash \
            "$RESULTS_FILE" > "$REMEDIATION_DIR/remediation.sh"
        chmod +x "$REMEDIATION_DIR/remediation.sh"
        ;;
    "ansible")
        oscap xccdf generate fix \
            --fix-type ansible \
            "$RESULTS_FILE" > "$REMEDIATION_DIR/remediation.yml"
        ;;
    "puppet")
        oscap xccdf generate fix \
            --fix-type puppet \
            "$RESULTS_FILE" > "$REMEDIATION_DIR/remediation.pp"
        ;;
    "chef")
        oscap xccdf generate fix \
            --fix-type chef \
            "$RESULTS_FILE" > "$REMEDIATION_DIR/remediation.rb"
        ;;
    *)
        echo "Unsupported remediation type: $REMEDIATION_TYPE"
        exit 1
        ;;
esac

# Phase 3: Backup current configuration
echo "Phase 3: Creating system backup..."
BACKUP_DIR="$REMEDIATION_DIR/system-backup"
mkdir -p "$BACKUP_DIR"

# Backup critical configuration files
cp -p /etc/passwd "$BACKUP_DIR/" 2>/dev/null || true
cp -p /etc/shadow "$BACKUP_DIR/" 2>/dev/null || true
cp -p /etc/group "$BACKUP_DIR/" 2>/dev/null || true
cp -p /etc/ssh/sshd_config "$BACKUP_DIR/" 2>/dev/null || true
cp -p /etc/sudoers "$BACKUP_DIR/" 2>/dev/null || true
cp -rp /etc/security "$BACKUP_DIR/" 2>/dev/null || true

# Phase 4: Execute or validate remediation
if [[ "$DRY_RUN" == "true" ]]; then
    echo "Phase 4: Dry run - validating remediation script..."
    case "$REMEDIATION_TYPE" in
        "bash")
            bash -n "$REMEDIATION_DIR/remediation.sh"
            if [[ $? -eq 0 ]]; then
                echo "✅ Bash remediation script syntax is valid"
            else
                echo "❌ Bash remediation script has syntax errors"
                exit 1
            fi
            ;;
        "ansible")
            ansible-playbook --syntax-check "$REMEDIATION_DIR/remediation.yml" || \
                echo "⚠️ Ansible syntax validation not available"
            ;;
        *)
            echo "ℹ️ Syntax validation not implemented for $REMEDIATION_TYPE"
            ;;
    esac
    
    echo "📋 Remediation preview:"
    head -50 "$REMEDIATION_DIR/remediation.$( [[ $REMEDIATION_TYPE == 'bash' ]] && echo 'sh' || echo 'yml' )"
    
else
    echo "Phase 4: Applying remediation..."
    case "$REMEDIATION_TYPE" in
        "bash")
            bash "$REMEDIATION_DIR/remediation.sh" > "$REMEDIATION_DIR/remediation-output.log" 2>&1
            REMEDIATION_EXIT_CODE=$?
            ;;
        "ansible")
            ansible-playbook "$REMEDIATION_DIR/remediation.yml" > "$REMEDIATION_DIR/remediation-output.log" 2>&1
            REMEDIATION_EXIT_CODE=$?
            ;;
        *)
            echo "❌ Automated execution not implemented for $REMEDIATION_TYPE"
            echo "Please execute the remediation script manually: $REMEDIATION_DIR/remediation.*"
            exit 1
            ;;
    esac
    
    if [[ $REMEDIATION_EXIT_CODE -eq 0 ]]; then
        echo "✅ Remediation applied successfully"
    else
        echo "❌ Remediation failed with exit code: $REMEDIATION_EXIT_CODE"
        echo "Check logs: $REMEDIATION_DIR/remediation-output.log"
        exit 1
    fi
fi

# Phase 5: Generate remediation report
echo "Phase 5: Generating remediation report..."
cat > "$REMEDIATION_DIR/remediation-report.md" <<EOF
# Security Remediation Report

**Remediation ID**: $REMEDIATION_ID  
**Date**: $(date)  
**Source Results**: $RESULTS_FILE  
**Remediation Type**: $REMEDIATION_TYPE  
**Execution Mode**: $(if [[ "$DRY_RUN" == "true" ]]; then echo "Dry Run"; else echo "Live Execution"; fi)

## Summary
- **Failed Rules**: $FAILED_RULES
- **Remediation Status**: $(if [[ "$DRY_RUN" == "true" ]]; then echo "Validated"; elif [[ ${REMEDIATION_EXIT_CODE:-0} -eq 0 ]]; then echo "Applied Successfully"; else echo "Failed"; fi)

## Files Generated
- \`remediation.$( [[ $REMEDIATION_TYPE == 'bash' ]] && echo 'sh' || echo 'yml' )\`: Remediation script/playbook
- \`system-backup/\`: System configuration backup
$(if [[ "$DRY_RUN" != "true" ]]; then echo "- \`remediation-output.log\`: Execution log"; fi)

## Next Steps
$(if [[ "$DRY_RUN" == "true" ]]; then
    echo "1. Review remediation script for appropriateness"
    echo "2. Test in non-production environment"
    echo "3. Execute with dry-run=false when ready"
    echo "4. Re-run compliance assessment to validate fixes"
else
    echo "1. Re-run compliance assessment to validate remediation"
    echo "2. Monitor system for any stability issues"
    echo "3. Update security documentation"
    echo "4. Schedule follow-up compliance reviews"
fi)

## Backup Information
System configuration backup available in: \`system-backup/\`
Use this backup to restore original configuration if needed.
EOF

echo "Remediation workflow complete!"
echo "Report: $REMEDIATION_DIR/remediation-report.md"
echo "Backup: $REMEDIATION_DIR/system-backup/"

if [[ "$DRY_RUN" == "true" ]]; then
    echo "🔍 DRY RUN: Review remediation script before applying"
    echo "Execute with: $0 $RESULTS_FILE $REMEDIATION_TYPE false"
else
    echo "🔄 Re-run compliance assessment to validate remediation"
fi
```

### Continuous Compliance Monitoring
```bash
#!/bin/bash
# continuous-compliance-monitoring.sh

CONFIG_FILE="$1"
ALERT_WEBHOOK="$2"
MONITORING_INTERVAL=3600  # 1 hour

if [[ -z "$CONFIG_FILE" || -z "$ALERT_WEBHOOK" ]]; then
    echo "Usage: $0 <config-file> <webhook-url>"
    echo "Config file format: JSON with systems and benchmarks"
    exit 1
fi

echo "Starting continuous compliance monitoring..."
echo "Configuration: $CONFIG_FILE"
echo "Monitoring interval: $MONITORING_INTERVAL seconds"

while true; do
    echo "$(date): Starting compliance monitoring cycle..."
    
    CYCLE_ID="monitor-$(date +%Y%m%d-%H%M%S)"
    MONITOR_DIR="compliance-monitoring-$CYCLE_ID"
    mkdir -p "$MONITOR_DIR"
    
    TOTAL_SYSTEMS=0
    COMPLIANT_SYSTEMS=0
    CRITICAL_FAILURES=0
    
    # Process each system configuration
    while IFS= read -r system_config; do
        SYSTEM_NAME=$(echo "$system_config" | jq -r '.name')
        CONTENT_FILE=$(echo "$system_config" | jq -r '.content_file')
        PROFILE_ID=$(echo "$system_config" | jq -r '.profile_id')
        THRESHOLD=$(echo "$system_config" | jq -r '.compliance_threshold // 90')
        
        [[ "$SYSTEM_NAME" == "null" ]] && continue
        ((TOTAL_SYSTEMS++))
        
        echo "Monitoring: $SYSTEM_NAME"
        SYSTEM_DIR="$MONITOR_DIR/$SYSTEM_NAME"
        mkdir -p "$SYSTEM_DIR"
        
        # Perform compliance evaluation
        oscap xccdf eval \
            --profile "$PROFILE_ID" \
            --results "$SYSTEM_DIR/results.xml" \
            --report "$SYSTEM_DIR/report.html" \
            "$CONTENT_FILE" > "$SYSTEM_DIR/evaluation.log" 2>&1
        
        # Analyze results
        if [[ -f "$SYSTEM_DIR/results.xml" ]]; then
            PASS_COUNT=$(grep -c 'result="pass"' "$SYSTEM_DIR/results.xml" 2>/dev/null || echo "0")
            FAIL_COUNT=$(grep -c 'result="fail"' "$SYSTEM_DIR/results.xml" 2>/dev/null || echo "0")
            TOTAL_RULES=$((PASS_COUNT + FAIL_COUNT))
            
            COMPLIANCE_PCT=0
            if [[ $TOTAL_RULES -gt 0 ]]; then
                COMPLIANCE_PCT=$((PASS_COUNT * 100 / TOTAL_RULES))
            fi
            
            echo "  Compliance: ${COMPLIANCE_PCT}% ($PASS_COUNT/$TOTAL_RULES)"
            
            if [[ $COMPLIANCE_PCT -ge $THRESHOLD ]]; then
                ((COMPLIANT_SYSTEMS++))
                echo "  Status: ✅ COMPLIANT"
            else
                echo "  Status: ❌ NON-COMPLIANT"
                ((CRITICAL_FAILURES++))
                
                # Send immediate alert for critical failures
                ALERT_MSG="🚨 COMPLIANCE ALERT: $SYSTEM_NAME failed compliance check ($COMPLIANCE_PCT% < $THRESHOLD%)"
                curl -X POST -H 'Content-type: application/json' \
                    --data "{\"text\":\"$ALERT_MSG\"}" \
                    "$ALERT_WEBHOOK" 2>/dev/null || true
            fi
            
            # Store monitoring data
            cat > "$SYSTEM_DIR/monitoring-data.json" <<EOF
{
    "system": "$SYSTEM_NAME",
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "compliance_percentage": $COMPLIANCE_PCT,
    "threshold": $THRESHOLD,
    "passed_rules": $PASS_COUNT,
    "failed_rules": $FAIL_COUNT,
    "total_rules": $TOTAL_RULES,
    "status": "$(if [[ $COMPLIANCE_PCT -ge $THRESHOLD ]]; then echo "COMPLIANT"; else echo "NON_COMPLIANT"; fi)"
}
EOF
        else
            echo "  Status: ⚠️ EVALUATION FAILED"
            ((CRITICAL_FAILURES++))
        fi
        
    done < <(jq -c '.systems[]' "$CONFIG_FILE")
    
    # Generate monitoring summary
    cat > "$MONITOR_DIR/monitoring-summary.json" <<EOF
{
    "cycle_id": "$CYCLE_ID",
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "summary": {
        "total_systems": $TOTAL_SYSTEMS,
        "compliant_systems": $COMPLIANT_SYSTEMS,
        "critical_failures": $CRITICAL_FAILURES,
        "overall_compliance_rate": $(( TOTAL_SYSTEMS > 0 ? COMPLIANT_SYSTEMS * 100 / TOTAL_SYSTEMS : 0 ))
    }
}
EOF
    
    # Send cycle summary
    SUMMARY_MSG="📊 Compliance Monitoring Summary: $COMPLIANT_SYSTEMS/$TOTAL_SYSTEMS systems compliant"
    if [[ $CRITICAL_FAILURES -gt 0 ]]; then
        SUMMARY_MSG="$SUMMARY_MSG ($CRITICAL_FAILURES critical failures)"
    fi
    
    curl -X POST -H 'Content-type: application/json' \
        --data "{\"text\":\"$SUMMARY_MSG\"}" \
        "$ALERT_WEBHOOK" 2>/dev/null || true
    
    echo "Monitoring cycle complete: $MONITOR_DIR/"
    echo "Compliance rate: $(( TOTAL_SYSTEMS > 0 ? COMPLIANT_SYSTEMS * 100 / TOTAL_SYSTEMS : 0 ))%"
    
    # Cleanup old monitoring data (keep last 7 days)
    find . -name "compliance-monitoring-*" -type d -mtime +7 -exec rm -rf {} \; 2>/dev/null || true
    
    echo "Next monitoring cycle in $MONITORING_INTERVAL seconds..."
    sleep $MONITORING_INTERVAL
done
```

## Integration Patterns

### GitHub Actions Workflow
```yaml
# .github/workflows/security-compliance.yml
name: Security Compliance Assessment
on:
  push:
    branches: [main, develop]
  schedule:
    - cron: '0 6 * * 1'  # Weekly Monday 6 AM
  workflow_dispatch:
    inputs:
      benchmark_type:
        description: 'Security benchmark to evaluate'
        required: true
        default: 'cis'
        type: choice
        options:
        - cis
        - stig
        - ospp

jobs:
  security-compliance:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Install OpenSCAP
      run: |
        sudo apt-get update
        sudo apt-get install -y libopenscap8 openscap-utils
        
    - name: Verify OpenSCAP Installation
      run: |
        oscap --version
        
    - name: Download Security Content
      env:
        BENCHMARK: ${{ github.event.inputs.benchmark_type || 'cis' }}
      run: |
        case "$BENCHMARK" in
          "cis")
            wget https://github.com/ComplianceAsCode/content/releases/latest/download/scap-security-guide-latest.zip
            unzip scap-security-guide-latest.zip
            CONTENT_FILE=$(find . -name "*cis*datastream.xml" | head -1)
            ;;
          "stig")
            wget https://public.cyber.mil/stigs/scap/
            # Download appropriate STIG content
            ;;
          "ospp")
            wget https://github.com/ComplianceAsCode/content/releases/latest/download/scap-security-guide-latest.zip
            unzip scap-security-guide-latest.zip
            CONTENT_FILE=$(find . -name "*ospp*datastream.xml" | head -1)
            ;;
        esac
        echo "CONTENT_FILE=$CONTENT_FILE" >> $GITHUB_ENV
        
    - name: Validate Security Content
      run: |
        echo "Validating SCAP content: $CONTENT_FILE"
        oscap info "$CONTENT_FILE"
        oscap ds sds-validate "$CONTENT_FILE"
        
    - name: Run Compliance Assessment
      run: |
        echo "Running compliance assessment..."
        
        # Determine profile based on benchmark
        case "${{ github.event.inputs.benchmark_type || 'cis' }}" in
          "cis")
            PROFILE="xccdf_org.ssgproject.content_profile_cis"
            ;;
          "stig")
            PROFILE="xccdf_org.ssgproject.content_profile_stig"
            ;;
          "ospp")
            PROFILE="xccdf_org.ssgproject.content_profile_ospp"
            ;;
        esac
        
        # Note: In real environment, this would scan the actual system
        # For demo purposes, we'll generate a guide and validate content
        oscap xccdf generate guide \
          --profile "$PROFILE" \
          "$CONTENT_FILE" > compliance-guide.html
          
        echo "PROFILE=$PROFILE" >> $GITHUB_ENV
        
    - name: Generate Compliance Documentation
      run: |
        # Create compliance documentation
        cat > compliance-summary.md <<EOF
        # Security Compliance Assessment
        
        **Benchmark**: ${{ github.event.inputs.benchmark_type || 'cis' }}
        **Profile**: $PROFILE
        **Date**: $(date)
        **Content**: $CONTENT_FILE
        
        ## Assessment Results
        Compliance guide generated successfully.
        Security content validated and ready for system evaluation.
        
        ## Next Steps
        1. Deploy this content to target systems
        2. Run full compliance evaluation
        3. Review and remediate any findings
        4. Implement continuous monitoring
        EOF
        
    - name: Upload Compliance Artifacts
      uses: actions/upload-artifact@v4
      with:
        name: compliance-assessment-${{ github.event.inputs.benchmark_type || 'cis' }}
        path: |
          compliance-guide.html
          compliance-summary.md
          ${{ env.CONTENT_FILE }}
```

### Terraform Security Module
```hcl
# terraform/security-compliance.tf
resource "null_resource" "openscap_compliance_check" {
  count = var.enable_compliance_scanning ? 1 : 0
  
  triggers = {
    instance_ids = join(",", var.instance_ids)
    benchmark   = var.security_benchmark
  }

  provisioner "local-exec" {
    command = <<-EOF
      echo "Downloading security content for ${var.security_benchmark}..."
      
      # Download appropriate security content
      case "${var.security_benchmark}" in
        "cis")
          CONTENT_URL="https://github.com/ComplianceAsCode/content/releases/latest/download/scap-security-guide-latest.zip"
          PROFILE="xccdf_org.ssgproject.content_profile_cis"
          ;;
        "stig")
          CONTENT_URL="https://public.cyber.mil/stigs/scap/"
          PROFILE="xccdf_org.ssgproject.content_profile_stig"
          ;;
        *)
          echo "Unsupported benchmark: ${var.security_benchmark}"
          exit 1
          ;;
      esac
      
      # Validate security content
      wget -O security-content.zip "$CONTENT_URL"
      unzip security-content.zip
      CONTENT_FILE=$(find . -name "*datastream.xml" | head -1)
      
      oscap info "$CONTENT_FILE"
      oscap ds sds-validate "$CONTENT_FILE"
      
      # Generate compliance guide
      oscap xccdf generate guide \
        --profile "$PROFILE" \
        "$CONTENT_FILE" > terraform-compliance-guide.html
        
      echo "Compliance content ready for deployment"
    EOF
  }
}

variable "enable_compliance_scanning" {
  description = "Enable security compliance scanning"
  type        = bool
  default     = true
}

variable "security_benchmark" {
  description = "Security benchmark to use"
  type        = string
  default     = "cis"
  
  validation {
    condition     = contains(["cis", "stig", "ospp"], var.security_benchmark)
    error_message = "Security benchmark must be one of: cis, stig, ospp."
  }
}

variable "instance_ids" {
  description = "List of instance IDs to assess"
  type        = list(string)
  default     = []
}

output "compliance_guide_path" {
  value = var.enable_compliance_scanning ? "terraform-compliance-guide.html" : "Compliance scanning disabled"
}
```

### Docker Integration
```dockerfile
# Dockerfile.openscap-scanner
FROM ubuntu:22.04

# Install OpenSCAP and dependencies
RUN apt-get update && apt-get install -y \
    libopenscap8 \
    openscap-utils \
    wget \
    unzip \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Verify OpenSCAP installation
RUN oscap --version

# Create scanning directory
WORKDIR /scans

# Copy scanning scripts
COPY openscap-assessment.sh /usr/local/bin/
COPY compliance-monitoring.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

# Set up entrypoint
ENTRYPOINT ["/usr/local/bin/openscap-assessment.sh"]
```

## Best Practices

### Content Management
- **Certified Content**: Use NIST-certified SCAP content from trusted sources
- **Content Validation**: Always validate SCAP content before evaluation
- **Version Control**: Track security content versions and updates
- **Profile Selection**: Choose appropriate security profiles for your environment

### Assessment Strategy
- **Baseline Establishment**: Create security configuration baselines
- **Regular Evaluation**: Schedule periodic compliance assessments
- **Incremental Scanning**: Focus on changed configurations
- **Risk Prioritization**: Address high-risk findings first

### Remediation Management
- **Testing Environment**: Test remediation scripts in non-production environments
- **Backup Strategy**: Create system backups before applying remediation
- **Phased Deployment**: Implement remediation in phases
- **Validation Process**: Re-assess systems after remediation

### Reporting and Documentation
- **Executive Summaries**: Provide high-level compliance status for management
- **Technical Details**: Include detailed findings for technical teams
- **Trend Analysis**: Track compliance trends over time
- **Audit Trails**: Maintain comprehensive audit documentation

## Error Handling

### Common Issues
```bash
# Content validation errors
oscap ds sds-validate content.xml
# Solution: Download fresh content from official sources

# Profile not found
oscap info content.xml | grep -A5 "Profiles:"
# Solution: Use correct profile ID from content information

# Permission denied during evaluation
sudo oscap xccdf eval --profile profile_id content.xml
# Solution: Run with appropriate privileges for system access

# Missing OVAL results
oscap xccdf eval --oval-results --profile profile_id content.xml
# Solution: Include --oval-results flag for detailed assessment
```

### Troubleshooting
- **Content Issues**: Verify SCAP content integrity and version compatibility
- **Permission Problems**: Ensure OpenSCAP has necessary system access privileges
- **Profile Errors**: Confirm profile IDs match those available in content
- **Performance Issues**: Consider using specific rule selection for faster evaluation

OpenSCAP provides comprehensive security compliance scanning and assessment capabilities, enabling organizations to maintain robust security posture through automated evaluation, detailed reporting, and systematic remediation of security configuration issues across enterprise environments.