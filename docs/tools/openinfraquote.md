# OpenInfraQuote MCP Tool

OpenInfraQuote (oiq) is a command-line tool for matching Terraform plans with cloud provider pricing to estimate infrastructure costs.

## Description

OpenInfraQuote provides:
- Terraform plan cost estimation
- Multi-cloud pricing support
- Detailed cost breakdowns
- JSON output for integration
- Custom pricing sheet support

## MCP Functions

### `openinfraquote_estimate_cost`
Estimate cost for Terraform plan using real oiq CLI.

**Parameters:**
- `tfplan_json` (required): Path to Terraform plan JSON file
- `pricesheet`: Custom pricing sheet to use
- `region`: Cloud region for pricing
- `provider`: Cloud provider (aws, gcp, azure)

**CLI Command:** `oiq match --pricesheet <pricesheet> <tfplan_json>`

### `openinfraquote_generate_plan`
Generate cost estimate from Terraform directory using real oiq CLI.

**Parameters:**
- `terraform_dir` (required): Path to Terraform configuration directory
- `pricesheet`: Custom pricing sheet to use
- `region`: Cloud region for pricing
- `provider`: Cloud provider (aws, gcp, azure)
- `var_file`: Terraform variables file

**CLI Command:** Uses terraform plan + oiq match pipeline

### `openinfraquote_compare_costs`
Compare costs between different configurations using real oiq CLI.

**Parameters:**
- `baseline_plan` (required): Path to baseline Terraform plan JSON
- `comparison_plan` (required): Path to comparison Terraform plan JSON
- `pricesheet`: Custom pricing sheet to use
- `output_format`: Output format (json, table)

**CLI Command:** `oiq match <baseline_plan> && oiq match <comparison_plan>`

### `openinfraquote_validate_plan`
Validate Terraform plan for cost estimation using real oiq CLI.

**Parameters:**
- `tfplan_json` (required): Path to Terraform plan JSON file
- `max_cost`: Maximum allowed cost threshold
- `fail_on_increase`: Fail if cost increases

**CLI Command:** `oiq match <tfplan_json>` with validation logic

## Common Use Cases

1. **Cost Estimation**: Estimate infrastructure costs from Terraform plans
2. **Budget Validation**: Ensure deployments stay within budget
3. **Cost Comparison**: Compare costs between different configurations
4. **CI/CD Integration**: Automated cost checks in pipelines
5. **Multi-Cloud Pricing**: Get pricing across different cloud providers

## Integration with Ship CLI

All OpenInfraQuote tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Estimate infrastructure costs automatically
- Validate budget constraints
- Compare deployment costs
- Integrate cost checks into workflows

The tools use containerized execution via Dagger for consistent, isolated cost estimation.