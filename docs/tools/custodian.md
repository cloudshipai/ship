# Cloud Custodian

Cloud governance engine for managing and securing cloud resources at scale.

## Description

Cloud Custodian (c7n) is a rules engine for managing cloud services and resources. It provides a unified set of policies to manage public cloud infrastructure, enabling organizations to enforce security, compliance, and cost optimization policies across AWS, Azure, and GCP.

## MCP Tools

### Policy Execution
- **`custodian_run_policy`** - Execute Cloud Custodian policies against cloud resources
- **`custodian_dry_run`** - Preview policy effects without making changes (dry run mode)

### Policy Management
- **`custodian_validate_policy`** - Validate policy syntax and structure
- **`custodian_schema`** - Get policy schema information for resource types

### Utility
- **`custodian_get_version`** - Get Cloud Custodian version information

## Real CLI Commands Used

- `custodian run -s out policy.yml` - Execute policies with output directory
- `custodian run --dryrun -s out policy.yml` - Preview policy execution
- `custodian validate policy.yml` - Validate policy syntax
- `custodian schema [resource-type]` - Get schema for resource types
- `custodian version` - Show version information

## Use Cases

### Cloud Governance
- Enforce security policies across cloud resources
- Automated compliance remediation
- Cost optimization through resource lifecycle management
- Security posture monitoring and enforcement

### Resource Management
- Automated cleanup of unused resources
- Tagging enforcement and standardization
- Security group and access control management
- Backup and retention policy enforcement

### Compliance
- CIS benchmark enforcement
- SOC 2 compliance automation
- PCI DSS security controls
- GDPR data protection policies

### Cost Optimization
- Unused resource identification and cleanup
- Right-sizing recommendations and enforcement
- Reserved instance optimization
- Spot instance lifecycle management

## Policy Examples

### Security Enforcement
- Remove default VPC security groups
- Encrypt unencrypted EBS volumes
- Disable public read access on S3 buckets
- Terminate instances without proper tags

### Cost Management
- Stop development instances after hours
- Delete unattached EBS volumes
- Remove unused elastic IPs
- Archive old S3 objects to cheaper storage

## Integration

Works with cloud provider APIs (AWS, Azure, GCP) and integrates with CI/CD pipelines for automated policy enforcement. Supports various execution modes including Lambda functions, Docker containers, and local execution.