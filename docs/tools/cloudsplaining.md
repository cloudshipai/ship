# Cloudsplaining MCP Tool

Cloudsplaining is an AWS IAM security assessment tool that identifies violations of least privilege and generates a risk-prioritized HTML report.

## Description

Cloudsplaining analyzes AWS IAM policies to identify:
- Overly permissive policies
- Privilege escalation risks
- Resource exposure issues
- Violations of least privilege principles
- Infrastructure modification permissions

## MCP Functions

### `cloudsplaining_scan`
Scan AWS account for IAM security issues using real cloudsplaining CLI.

**Parameters:**
- `input_file` (required): Path to AWS account authorization details JSON file
- `output_directory`: Output directory for report (default: output)
- `skip_open_report`: Skip opening report in browser

**CLI Command:** `cloudsplaining scan --input-file <input_file> [--output-directory <output_directory>] [--skip-open-report]`

### `cloudsplaining_scan_policy_file`
Scan specific IAM policy file using real cloudsplaining CLI.

**Parameters:**
- `input_file` (required): Path to IAM policy JSON file
- `output_directory`: Output directory for report

**CLI Command:** `cloudsplaining scan-policy-file --input-file <input_file> [--output-directory <output_directory>]`

### `cloudsplaining_download`
Download AWS account authorization details using real cloudsplaining CLI.

**Parameters:**
- `profile`: AWS CLI profile name
- `output_file`: Output file path (default: account-auth-details.json)

**CLI Command:** `cloudsplaining download [--profile <profile>] [--output-file <output_file>]`

### `cloudsplaining_create_exclusions`
Create exclusions template using real cloudsplaining CLI.

**Parameters:**
- `output_file`: Output file path for exclusions template

**CLI Command:** `cloudsplaining create-exclusions-file [--output-file <output_file>]`

### `cloudsplaining_scan_with_exclusions`
Scan with exclusions file using real cloudsplaining CLI.

**Parameters:**
- `input_file` (required): Path to AWS account authorization details JSON file
- `exclusions_file` (required): Path to exclusions file
- `output_directory`: Output directory for report

**CLI Command:** `cloudsplaining scan --input-file <input_file> --exclusions-file <exclusions_file> [--output-directory <output_directory>]`

### `cloudsplaining_version`
Get Cloudsplaining version information.

**CLI Command:** `cloudsplaining --version`

## Common Use Cases

1. **IAM Security Assessment**: Identify overly permissive policies
2. **Privilege Escalation Detection**: Find potential escalation paths
3. **Compliance Validation**: Ensure least privilege compliance
4. **Risk Prioritization**: Focus on high-risk permissions
5. **Policy Remediation**: Generate actionable recommendations

## Integration with Ship CLI

All Cloudsplaining tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated IAM security assessments
- Generate risk-prioritized reports
- Identify privilege escalation paths
- Provide remediation recommendations

The tools use containerized execution via Dagger for consistent, isolated IAM analysis.