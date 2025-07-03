# Ship CLI Terraform Tools Demos

✅ **Successfully created demo GIFs for working Terraform tools!**

## Created Demos

1. **terraform-generate-docs.gif** - Generates Terraform documentation using terraform-docs
2. **terraform-lint.gif** - Lints Terraform code using TFLint
3. **terraform-security-scan.gif** - Scans for security issues using Trivy
4. **terraform-checkov-scan.gif** - Comprehensive security scanning with Checkov
5. **terraform-tools-demo.gif** - Combined demo showing multiple tools (may be truncated)

## Working Tools

These tools work without requiring AWS credentials or terraform plan:
- `generate-docs` - Creates README documentation from Terraform code
- `lint` - Validates Terraform syntax and best practices
- `security-scan` - Identifies security vulnerabilities
- `checkov-scan` - Comprehensive policy-as-code scanning

## Tools Requiring AWS Credentials

These tools need `terraform init` and `terraform plan` to work:
- `cost-analysis` - Uses OpenInfraQuote (requires terraform plan)
- `cost-estimate` - Uses Infracost (requires terraform plan)
- `generate-diagram` - Uses InfraMap (requires terraform state)

## Usage Examples

```bash
# Generate documentation
ship terraform-tools generate-docs /path/to/terraform

# Lint code
ship terraform-tools lint /path/to/terraform

# Security scan
ship terraform-tools security-scan /path/to/terraform

# Checkov scan
ship terraform-tools checkov-scan /path/to/terraform

# Push results to CloudShip
ship terraform-tools security-scan /path/to/terraform --push
```

## Notes

- All demos use the `examples/terraform/easy-s3-bucket` module
- The `--push` flag can be added to any command to upload results to CloudShip
- Tools run in containerized environments via Dagger, no local installation needed