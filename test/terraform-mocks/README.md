# Terraform Mock Infrastructure

This directory contains example Terraform code designed to demonstrate Ship CLI's terraform-tools functionality. The infrastructure includes common AWS resources with **intentional security issues** for testing security scanning tools.

## Infrastructure Overview

- **EC2 Instance** (`aws_instance.example`) - t3.medium with encrypted EBS volume
- **S3 Bucket** (`aws_s3_bucket.example`) - With public-read ACL (security issue)
- **RDS Database** (`aws_db_instance.example`) - PostgreSQL with unencrypted storage (security issue)
- **Security Group** (`aws_security_group.db`) - Allows 0.0.0.0/0 access to PostgreSQL (security issue)

## Security Issues (Intentional)

This example includes several security misconfigurations for testing:

1. **Hardcoded AWS credentials** in EC2 user_data
2. **Public S3 bucket ACL** allowing public-read access
3. **Unencrypted RDS storage** (`storage_encrypted = false`)
4. **Overly permissive security group** allowing 0.0.0.0/0 access to PostgreSQL
5. **Hardcoded database password** in plain text

## Ship CLI Usage Examples

### Prerequisites

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Create variables file (already provided as `terraform.tfvars`)

### Terraform Tools Commands

#### 1. Security Scanning ✅
```bash
ship terraform-tools security-scan
```
**Output**: [security-scan-summary.txt](../../outputs/terraform-tools/security-scan-summary.txt)

**Result**: Identifies 21 security misconfigurations including:
- Instance metadata service vulnerabilities
- RDS encryption issues  
- S3 public access problems
- Missing logging configurations

#### 2. Generate Documentation ✅
```bash
ship terraform-tools generate-docs
```
**Output**: [generate-docs.txt](../../outputs/terraform-tools/generate-docs.txt)

**Result**: Successfully generates comprehensive Terraform documentation (see this README!)

#### 3. Code Linting ✅
```bash
ship terraform-tools lint
```
**Output**: [lint.txt](../../outputs/terraform-tools/lint.txt)

**Result**: Clean validation - no syntax errors or lint issues found

#### 4. Checkov Security Analysis ✅
```bash
ship terraform-tools checkov-scan
```
**Output**: [checkov-scan-summary.txt](../../outputs/terraform-tools/checkov-scan-summary.txt)

**Result**: Comprehensive security policy validation with detailed findings and remediation guidance

### Additional Commands (Currently Not Working)

The following commands are available but currently have issues:

#### Cost Analysis & Infrastructure Diagramming
- `ship terraform-tools cost-analysis` - ❌ Container cannot access .terraform directory
- `ship terraform-tools generate-diagram` - ❌ Requires terraform.tfstate from actual deployment
- `ship terraform-tools cost-estimate` (Infracost) - ❌ Removed due to API key requirements and container issues

#### AI-Powered Investigation ✅  
```bash
ship investigate --prompt "analyze my AWS infrastructure security"
```
**Note**: Requires AWS credentials as environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION)

**Result**: AI agent successfully analyzes infrastructure and provides detailed recommendations

These features need additional development work to function properly.

## AI Agent Capabilities

The Ship CLI AI agent successfully provides:

- ✅ **Infrastructure Analysis**: Real-time SQL queries via Steampipe 
- ✅ **Security Assessments**: Identifies misconfigurations and vulnerabilities
- ✅ **Cost Optimization**: Provides recommendations for resource optimization
- ✅ **Compliance Checks**: Analyzes against best practices and policies
- ❌ **Cannot execute**: terraform-tools commands (security-scan, lint, etc.)
- ❌ **Cannot choose**: Between different CLI commands automatically

**Working**: AI investigations with `ship investigate --prompt "your question"`
**Manual**: terraform-tools commands for static code analysis and documentation

## Troubleshooting

### AWS Credentials for AI Investigation

If you have AWS credentials in `~/.aws/credentials` but the investigate command isn't working:

```bash
# Extract credentials from your AWS credentials file
grep -A 3 "\[default\]" ~/.aws/credentials

# Set them as environment variables
export AWS_ACCESS_KEY_ID="your-key-from-file"
export AWS_SECRET_ACCESS_KEY="your-secret-from-file"
export AWS_REGION="us-east-1"

# Verify they're set
echo $AWS_ACCESS_KEY_ID
```

### Current Known Issues

1. **Cost Analysis**: Container cannot access .terraform directory (needs Dagger container fix)
2. **Generate Diagram**: Requires terraform.tfstate from actual deployment
3. **Credential Inconsistency**: AI investigation uses environment variables while terraform-tools use credential files

## Best Practices

1. **Never commit** this example to production repositories due to intentional security issues
2. **Use Ship CLI tools** in CI/CD pipelines for automated security scanning
3. **Combine multiple tools** for comprehensive analysis:
   - Use `security-scan` + `checkov-scan` for thorough security analysis
   - Use `generate-docs` for maintaining documentation
4. **Set AWS credentials as environment variables** for AI investigations
5. **Test AI investigations** against real infrastructure for operational insights

## Resource Cleanup

To avoid AWS charges, destroy the test infrastructure:

```bash
terraform destroy
```

**Warning**: Only run this if you actually deployed the infrastructure with `terraform apply`.

---

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_db_instance.example](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance) | resource |
| [aws_instance.example](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance) | resource |
| [aws_s3_bucket.example](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_acl.example](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl) | resource |
| [aws_security_group.db](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_ami_id"></a> [ami\_id](#input\_ami\_id) | AMI ID for EC2 instances | `string` | `"ami-0c02fb55956c7d316"` | no |
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | AWS region to deploy resources | `string` | `"us-east-1"` | no |
| <a name="input_db_instance_class"></a> [db\_instance\_class](#input\_db\_instance\_class) | RDS instance class | `string` | `"db.t3.micro"` | no |
| <a name="input_environment"></a> [environment](#input\_environment) | Environment name (dev, staging, prod) | `string` | `"dev"` | no |
| <a name="input_instance_type"></a> [instance\_type](#input\_instance\_type) | EC2 instance type | `string` | `"t3.medium"` | no |
| <a name="input_project_name"></a> [project\_name](#input\_project\_name) | Name of the project | `string` | `"ship-test"` | no |
| <a name="input_root_volume_size"></a> [root\_volume\_size](#input\_root\_volume\_size) | Size of the root EBS volume in GB | `number` | `30` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_db_instance_address"></a> [db\_instance\_address](#output\_db\_instance\_address) | Address of the RDS instance |
| <a name="output_db_instance_endpoint"></a> [db\_instance\_endpoint](#output\_db\_instance\_endpoint) | Connection endpoint for the RDS instance |
| <a name="output_db_instance_port"></a> [db\_instance\_port](#output\_db\_instance\_port) | Port of the RDS instance |
| <a name="output_instance_id"></a> [instance\_id](#output\_instance\_id) | ID of the EC2 instance |
| <a name="output_instance_private_ip"></a> [instance\_private\_ip](#output\_instance\_private\_ip) | Private IP address of the EC2 instance |
| <a name="output_instance_public_ip"></a> [instance\_public\_ip](#output\_instance\_public\_ip) | Public IP address of the EC2 instance |
| <a name="output_s3_bucket_arn"></a> [s3\_bucket\_arn](#output\_s3\_bucket\_arn) | ARN of the S3 bucket |
| <a name="output_s3_bucket_name"></a> [s3\_bucket\_name](#output\_s3\_bucket\_name) | Name of the S3 bucket |
