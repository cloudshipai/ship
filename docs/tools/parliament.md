# Parliament MCP Tool

Parliament is an AWS IAM policy linter that identifies security misconfigurations and potential privilege escalation paths in IAM policies.

## Description

Parliament performs static analysis of AWS IAM policies to detect:
- Overprivileged permissions
- Missing condition statements
- Potential privilege escalation risks
- AWS service-specific security best practices violations

## MCP Functions

### `parliament_lint_file`
Lint AWS IAM policy file using real parliament CLI.

**Parameters:**
- `policy_path` (required): Path to IAM policy JSON file
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --file <policy_path> [--config <config>] [--json]`

### `parliament_lint_directory`
Lint AWS IAM policy files in directory using real parliament CLI.

**Parameters:**
- `directory` (required): Directory containing IAM policy files
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format
- `include_policy_extension`: File extension to include (e.g., json)
- `exclude_pattern`: Pattern to exclude (regex)

**CLI Command:** `parliament --directory <directory> [--config <config>] [--json] [--include_policy_extension <extension>] [--exclude_pattern <pattern>]`

### `parliament_lint_string`
Lint AWS IAM policy JSON string using real parliament CLI.

**Parameters:**
- `policy_json` (required): IAM policy JSON string
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --string <policy_json> [--config <config>] [--json]`

### `parliament_lint_community`
Lint AWS IAM policy with community auditors using real parliament CLI.

**Parameters:**
- `policy_path` (required): Path to IAM policy JSON file
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --file <policy_path> --include-community-auditors [--config <config>] [--json]`

### `parliament_lint_private`
Lint AWS IAM policy with private auditors using real parliament CLI.

**Parameters:**
- `policy_path` (required): Path to IAM policy JSON file
- `private_auditors` (required): Path to private auditors directory
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --file <policy_path> --private_auditors <private_auditors> [--config <config>] [--json]`

### `parliament_lint_aws_managed`
Lint AWS managed policies using real parliament CLI.

**Parameters:**
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --aws-managed-policies [--config <config>] [--json]`

### `parliament_lint_auth_details`
Lint AWS IAM authorization details file using real parliament CLI.

**Parameters:**
- `auth_details_file` (required): Path to AWS IAM authorization details file
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --auth-details-file <auth_details_file> [--config <config>] [--json]`

### `parliament_comprehensive_analysis`
Comprehensive IAM policy analysis with all auditors using real parliament CLI.

**Parameters:**
- `policy_path` (required): Path to IAM policy JSON file
- `private_auditors`: Path to private auditors directory
- `config`: Path to custom configuration file
- `json_output`: Output results in JSON format

**CLI Command:** `parliament --file <policy_path> --include-community-auditors [--private_auditors <private_auditors>] [--config <config>] [--json]`

### `parliament_batch_directory_analysis`
Batch analysis of multiple policy directories using real parliament CLI.

**Parameters:**
- `base_directory` (required): Base directory containing policy subdirectories
- `config`: Path to custom configuration file
- `private_auditors`: Path to private auditors directory
- `json_output`: Output results in JSON format
- `include_policy_extension`: File extension to include (default: json)
- `exclude_pattern`: Pattern to exclude from analysis

**CLI Command:** `parliament --directory <base_directory> --include-community-auditors [--private_auditors <private_auditors>] [--config <config>] [--json] [--include_policy_extension <extension>] [--exclude_pattern <pattern>]`

## Common Use Cases

1. **Single Policy Analysis**: Use `parliament_lint_file` to analyze individual IAM policy files
2. **Directory Scanning**: Use `parliament_lint_directory` to analyze all policies in a directory
3. **Policy String Validation**: Use `parliament_lint_string` to validate policy JSON directly
4. **Enhanced Auditing**: Use `parliament_lint_community` for additional community-based checks
5. **Custom Auditors**: Use `parliament_lint_private` with custom auditing rules
6. **AWS Managed Policies**: Use `parliament_lint_aws_managed` to analyze AWS-provided policies
7. **Authorization Details**: Use `parliament_lint_auth_details` for comprehensive account analysis

## Integration with Ship CLI

All Parliament tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated IAM policy security analysis
- Identify privilege escalation vulnerabilities
- Validate policy configurations against security best practices
- Generate security reports in multiple formats

The tools use containerized execution via Dagger for consistent, isolated analysis environments.