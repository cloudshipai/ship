# Terraform Docs MCP Tool

terraform-docs is a utility to generate documentation from Terraform modules in various output formats.

## Description

terraform-docs provides:
- Automatic documentation generation from Terraform code
- Multiple output formats (markdown, json, yaml, xml, etc.)
- Customizable documentation templates
- Module input/output documentation
- Provider and resource documentation

## MCP Functions

### `terraform_docs_markdown`
Generate Markdown documentation using real terraform-docs CLI.

**Parameters:**
- `module_path` (required): Path to Terraform module directory
- `output_file`: Output file path
- `sort_by`: Sort inputs/outputs by (name, required, type)
- `show_header`: Show file header in output

**CLI Command:** `terraform-docs markdown <module_path> [--output-file <output_file>] [--sort-by <sort_by>] [--header-from <module_path>]`

### `terraform_docs_json`
Generate JSON documentation using real terraform-docs CLI.

**Parameters:**
- `module_path` (required): Path to Terraform module directory
- `output_file`: Output file path
- `sort_by`: Sort inputs/outputs by (name, required, type)

**CLI Command:** `terraform-docs json <module_path> [--output-file <output_file>] [--sort-by <sort_by>]`

### `terraform_docs_yaml`
Generate YAML documentation using real terraform-docs CLI.

**Parameters:**
- `module_path` (required): Path to Terraform module directory
- `output_file`: Output file path
- `sort_by`: Sort inputs/outputs by (name, required, type)

**CLI Command:** `terraform-docs yaml <module_path> [--output-file <output_file>] [--sort-by <sort_by>]`

### `terraform_docs_xml`
Generate XML documentation using real terraform-docs CLI.

**Parameters:**
- `module_path` (required): Path to Terraform module directory
- `output_file`: Output file path
- `sort_by`: Sort inputs/outputs by (name, required, type)

**CLI Command:** `terraform-docs xml <module_path> [--output-file <output_file>] [--sort-by <sort_by>]`

## Common Use Cases

1. **Documentation Generation**: Auto-generate module documentation
2. **README Updates**: Keep README.md files updated with module info
3. **CI/CD Integration**: Automated documentation in pipelines
4. **Module Registry**: Generate docs for Terraform module registries
5. **Team Documentation**: Standardize module documentation format

## Integration with Ship CLI

All terraform-docs tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Generate comprehensive module documentation
- Keep documentation synchronized with code
- Support multiple output formats
- Automate documentation workflows

The tools use containerized execution via Dagger for consistent, isolated documentation generation.