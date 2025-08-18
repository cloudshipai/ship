# Conftest MCP Tool

Conftest is a utility to help you write tests against structured configuration data using Open Policy Agent (OPA) Rego language.

## Description

Conftest enables:
- Policy testing for Kubernetes, Terraform, Dockerfile, and more
- Custom policy enforcement across configuration files
- Integration with CI/CD pipelines
- Support for multiple input formats (YAML, JSON, HCL, TOML, etc.)
- Shareable policy libraries

## MCP Functions

### `conftest_verify`
Verify configuration files against policies using real conftest CLI.

**Parameters:**
- `path` (required): Path to file or directory to test
- `policy`: Path to directory containing Rego policies
- `namespace`: Policy namespace to use
- `output`: Output format (stdout, json, tap, table)

**CLI Command:** `conftest verify <path> [--policy <policy>] [--namespace <namespace>] [--output <output>]`

### `conftest_test`
Test configuration files against policies using real conftest CLI.

**Parameters:**
- `path` (required): Path to file or directory to test
- `policy`: Path to directory containing Rego policies
- `namespace`: Policy namespace to use
- `output`: Output format (stdout, json, tap, table)
- `fail_on_warn`: Exit with non-zero code on warnings

**CLI Command:** `conftest test <path> [--policy <policy>] [--namespace <namespace>] [--output <output>] [--fail-on-warn]`

### `conftest_pull`
Pull policies from registry using real conftest CLI.

**Parameters:**
- `url` (required): URL of policy bundle to pull
- `policy`: Directory to save policies

**CLI Command:** `conftest pull <url> [--policy <policy>]`

### `conftest_push`
Push policies to registry using real conftest CLI.

**Parameters:**
- `url` (required): URL of registry to push to
- `policy`: Directory containing policies to push

**CLI Command:** `conftest push <url> [--policy <policy>]`

### `conftest_parse`
Parse configuration file and output JSON using real conftest CLI.

**Parameters:**
- `file` (required): File to parse
- `parser`: Parser to use (auto-detected if not specified)

**CLI Command:** `conftest parse <file> [--parser <parser>]`

## Common Use Cases

1. **Policy Enforcement**: Enforce organizational policies on configurations
2. **Security Validation**: Check for security misconfigurations
3. **Compliance Testing**: Validate compliance requirements
4. **CI/CD Integration**: Automated policy checks in pipelines
5. **Multi-format Support**: Test various configuration formats

## Integration with Ship CLI

All Conftest tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated policy testing
- Validate configuration compliance
- Share and manage policy libraries
- Integrate with development workflows

The tools use containerized execution via Dagger for consistent, isolated policy testing.