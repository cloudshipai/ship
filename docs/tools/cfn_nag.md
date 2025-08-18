# CFN Nag - CloudFormation Security Linter

CFN Nag is a static analysis tool for finding potential security issues in CloudFormation templates. It looks for patterns in CloudFormation templates that may indicate insecure infrastructure.

## Overview

CFN Nag scans CloudFormation templates for patterns that may indicate insecure infrastructure. It has over 140 predefined rules checking for issues like:
- IAM rules that are too permissive (wildcards)
- Security group rules that are too permissive
- Access logs that aren't enabled
- Encryption that isn't enabled

## Available MCP Functions

### 1. `cfn_nag_scan`
**Description**: Scan CloudFormation templates for security issues using cfn_nag_scan

**Parameters**:
- `input_path` (required): Path to CloudFormation template file or directory
- `output_format` (optional): Output format (json or text)
- `debug` (optional): Dump information about rule loading

**Example Usage**:
```bash
# Scan a single template
cfn_nag_scan(input_path="template.yaml")

# Scan directory with JSON output
cfn_nag_scan(input_path="./cloudformation", output_format="json")

# Scan with debug information
cfn_nag_scan(input_path="template.yaml", debug=true)
```

### 2. `cfn_nag_scan_with_profile`
**Description**: Scan CloudFormation template with specific rule profile

**Parameters**:
- `input_path` (required): Path to CloudFormation template file or directory
- `profile_path` (optional): Path to profile file specifying rules to apply
- `deny_list_path` (optional): Path to deny list file specifying rules to ignore

**Example Usage**:
```bash
# Scan with custom profile
cfn_nag_scan_with_profile(
  input_path="template.yaml",
  profile_path="security-profile.yaml"
)

# Scan with deny list to ignore specific rules
cfn_nag_scan_with_profile(
  input_path="template.yaml",
  deny_list_path="rules-to-ignore.yaml"
)
```

### 3. `cfn_nag_scan_with_parameters`
**Description**: Scan CloudFormation template with parameter values

**Parameters**:
- `input_path` (required): Path to CloudFormation template file or directory
- `parameter_values_path` (optional): Path to JSON file with parameter values
- `condition_values_path` (optional): Path to JSON file with condition values
- `rule_arguments` (optional): Custom rule thresholds (e.g., spcm_threshold:100)

**Example Usage**:
```bash
# Scan with parameter values
cfn_nag_scan_with_parameters(
  input_path="template.yaml",
  parameter_values_path="parameters.json"
)

# Scan with custom rule thresholds
cfn_nag_scan_with_parameters(
  input_path="template.yaml",
  rule_arguments="spcm_threshold:100"
)
```

### 4. `cfn_nag_list_rules`
**Description**: List all available CFN Nag rules

**Parameters**: None

**Example Usage**:
```bash
cfn_nag_list_rules()
```

### 5. `cfn_nag_spcm_scan`
**Description**: Generate Stelligent Policy Complexity Metrics report

**Parameters**:
- `input_path` (required): Path to CloudFormation template file or directory
- `output_format` (optional): Output format for report (json or html)

**Example Usage**:
```bash
# Generate JSON metrics report
cfn_nag_spcm_scan(input_path="template.yaml", output_format="json")

# Generate HTML metrics report
cfn_nag_spcm_scan(input_path="./templates", output_format="html")
```

### 6. `cfn_nag_get_version`
**Description**: Get cfn_nag version information

**Parameters**: None

**Example Usage**:
```bash
cfn_nag_get_version()
```

## Real CLI Capabilities

All MCP functions are based on actual cfn_nag commands:

### Main Commands
```bash
# Scan templates
cfn_nag_scan --input-path template.yaml
cfn_nag_scan --input-path ./templates --output-format json

# List rules
cfn_nag_rules

# SPCM metrics
spcm_scan --input-path template.yaml --output-format html
```

### Advanced Options
```bash
# With profile
cfn_nag_scan --input-path template.yaml --profile-path profile.yaml

# With deny list
cfn_nag_scan --input-path template.yaml --deny-list-path deny.yaml

# With parameters
cfn_nag_scan --input-path template.yaml --parameter-values-path params.json

# With condition values
cfn_nag_scan --input-path template.yaml --condition-values-path conditions.json

# With custom rule arguments
cfn_nag_scan --input-path template.yaml --rule-arguments spcm_threshold:100
```

## Prerequisites

### Installation
```bash
# Install via Ruby gem
gem install cfn-nag

# Verify installation
cfn_nag_scan --version
```

### Requirements
- Ruby 2.5 or later
- RubyGems package manager

## File Types Supported

CFN Nag automatically processes these file extensions:
- `.json` - JSON CloudFormation templates
- `.template` - CloudFormation template files
- `.yml` - YAML CloudFormation templates
- `.yaml` - YAML CloudFormation templates

When scanning directories, all files with these extensions are processed, including subdirectories.

## Rule Categories

### IAM Rules
- Check for overly permissive policies
- Detect wildcard permissions
- Identify missing resource constraints

### Security Group Rules
- Detect overly broad ingress rules
- Check for unrestricted egress
- Identify missing protocol specifications

### Encryption Rules
- Check for unencrypted storage
- Verify encryption at rest
- Detect missing KMS configurations

### Logging Rules
- Verify access logging is enabled
- Check for audit trail configuration
- Detect missing CloudTrail setup

### Network Rules
- Check VPC configurations
- Verify subnet settings
- Detect public access issues

## Profile Configuration

### Profile File Format
```yaml
# profile.yaml
RulesToApply:
  - W1
  - W2
  - W3
  - F1
  - F2
```

### Deny List Format
```yaml
# deny-list.yaml
RulesToSuppress:
  - W4
  - W5
  - F3
```

## Parameter Values

### Parameter File Format
```json
{
  "InstanceType": "t3.micro",
  "KeyName": "my-key",
  "SubnetId": "subnet-12345"
}
```

### Condition Values Format
```json
{
  "CreateProdResources": true,
  "EnableLogging": false
}
```

## Exit Codes

- **0**: No violations found
- **Non-zero**: Violations detected (failing violations)

## Best Practices

### Template Development
- Run cfn_nag early in development
- Fix security issues before deployment
- Use profiles for consistent checks
- Document suppressed rules

### CI/CD Integration
- Add cfn_nag to pipeline checks
- Fail builds on violations
- Use JSON output for parsing
- Track metrics over time

### Rule Management
- Create organization-specific profiles
- Document why rules are suppressed
- Review suppressed rules regularly
- Update rules as needed

## Suppressing Rules

### In-Template Suppression
```yaml
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Metadata:
      cfn_nag:
        rules_to_suppress:
          - id: W35
            reason: "Access logging not required for this development bucket"
```

### Command-Line Suppression
Use deny list files to suppress rules across all templates.

## SPCM Metrics

The Stelligent Policy Complexity Metrics provide insights into:
- Policy complexity scores
- Resource relationship analysis
- Security posture metrics
- Compliance indicators

## Troubleshooting

### Common Issues

1. **Ruby Not Found**
   - Install Ruby 2.5 or later
   - Update PATH environment variable

2. **Gem Installation Fails**
   - Update RubyGems: `gem update --system`
   - Use sudo if needed: `sudo gem install cfn-nag`

3. **Template Parse Errors**
   - Validate JSON/YAML syntax
   - Check for CloudFormation-specific syntax

4. **Rules Not Applied**
   - Verify profile file format
   - Check profile file path

### Debug Mode
Use `--debug` flag to see:
- Rule loading details
- Template parsing information
- Rule evaluation process

## Integration with Ship CLI

These MCP functions integrate with Ship CLI's containerized execution:
- cfn_nag runs in a Ruby container via Dagger
- Templates are mounted into the container
- Results are returned to Ship CLI

## Custom Rules

CFN Nag supports custom rules written in Ruby:
```ruby
# custom_rule.rb
require 'cfn-nag/custom_rule'

class MyCustomRule < CustomRule
  def rule_id
    'C1'
  end
  
  def rule_type
    Violation::FAILING_VIOLATION
  end
  
  def rule_text
    'Custom rule description'
  end
  
  def audit(cfn_model)
    # Rule logic here
  end
end
```

## References

- **Official Repository**: https://github.com/stelligent/cfn_nag
- **Rule Documentation**: https://github.com/stelligent/cfn_nag/blob/master/doc/cfn_nag_rules.md
- **CloudFormation Best Practices**: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/best-practices.html
- **Stelligent Blog**: https://stelligent.com/category/cfn-nag/