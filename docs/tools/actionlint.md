# Actionlint - GitHub Actions Workflow Linter

Actionlint is a static checker for GitHub Actions workflow files that catches errors and validates workflow syntax.

## Overview

Actionlint analyzes GitHub Actions workflow files to detect:
- Syntax validation for workflow YAML files
- Type checking for `${{ }}` expressions  
- Actions usage verification
- Security vulnerability detection
- Integration with shellcheck and pyflakes for script validation

## Available MCP Functions

The MCP server provides discrete functions that map directly to CLI capabilities:

### 1. `actionlint_scan_workflows`
**Description**: Scan GitHub Actions workflow files for issues

**Parameters**:
- `workflow_files` (optional): Comma-separated list of workflow file paths (leave empty to scan all)
- `format_template` (optional): Go template for formatting output
- `ignore_patterns` (optional): Comma-separated regex patterns to ignore errors
- `color` (boolean, optional): Enable colored output

**MCP Usage**:
```bash
# Scan all workflows in repository
actionlint_scan_workflows()

# Scan specific workflow files
actionlint_scan_workflows(workflow_files=".github/workflows/ci.yml,.github/workflows/deploy.yml")

# Scan with custom format and ignore patterns
actionlint_scan_workflows(
  format_template="{{json .}}",
  ignore_patterns="SC2086,unused variable.*",
  color=true
)
```

**Equivalent CLI Command**:
```bash
ship security actionlint . --format "{{json .}}" --ignore "SC2086,unused variable.*" --color
```

### 2. `actionlint_scan_with_external_tools`
**Description**: Scan workflows with shellcheck and pyflakes integration

**Parameters**:
- `workflow_files` (optional): Comma-separated list of workflow file paths (leave empty to scan all)
- `shellcheck_path` (optional): Path to shellcheck executable
- `pyflakes_path` (optional): Path to pyflakes executable  
- `color` (boolean, optional): Enable colored output

**MCP Usage**:
```bash
# Scan with external tools using default paths
actionlint_scan_with_external_tools()

# Scan with custom tool paths
actionlint_scan_with_external_tools(
  shellcheck_path="/usr/local/bin/shellcheck",
  pyflakes_path="/usr/local/bin/pyflakes",
  color=true
)
```

**Equivalent CLI Command**:
```bash
ship security actionlint . --shellcheck /usr/local/bin/shellcheck --pyflakes /usr/local/bin/pyflakes --color
```

### 3. `actionlint_get_version`
**Description**: Get Actionlint version information

**Parameters**: None

**MCP Usage**:
```bash
actionlint_get_version()
```

**Equivalent CLI Command**:
```bash
ship security actionlint --version
```

## Real CLI Capabilities

All MCP functions are based on the actual actionlint CLI capabilities, and the Ship CLI provides enhanced access to these features:

### Ship CLI Usage
```bash
# Basic workflow scanning
ship security actionlint [directory]

# Advanced scanning with format template
ship security actionlint . --format "{{.Path}}: {{.Message}}"

# Scanning with ignore patterns
ship security actionlint . --ignore ".*test.*,.*example.*"

# Enable colored output
ship security actionlint . --color

# Integration with external tools
ship security actionlint . --shellcheck /usr/bin/shellcheck --pyflakes /usr/bin/pyflakes

# Save output to file
ship security actionlint . --output results.txt
```

### Ship CLI Flags
- `-o, --output`: Output file to save results (default: print to stdout)
- `-c, --config`: Path to actionlint configuration file
- `-f, --format`: Go template for formatting output
- `-i, --ignore`: Comma-separated regex patterns to ignore errors
- `--color`: Enable colored output
- `--shellcheck`: Path to shellcheck executable for shell script validation
- `--pyflakes`: Path to pyflakes executable for Python validation

### Direct actionlint CLI Usage
```bash
# Check all workflows in repository
actionlint

# Check specific workflow files
actionlint .github/workflows/ci.yml .github/workflows/deploy.yml

# Read from stdin
cat workflow.yml | actionlint -
```

### Supported actionlint Flags
- `-ignore <regex>`: Filter errors by message using regular expression (repeatable)
- `-format <template>`: Format output using Go template syntax
- `-color`: Enable colored output
- `-shellcheck [path]`: Specify shellcheck executable path
- `-pyflakes [path]`: Specify pyflakes executable path
- `-version`: Display actionlint version

### Output Formats
Via `-format` flag, actionlint supports:
- JSON format
- Markdown format
- JSON Lines format
- GitHub Actions annotations
- SARIF format
- Custom Go templates

## Installation

Install actionlint using one of these methods:

```bash
# Using go install
go install github.com/rhysd/actionlint/cmd/actionlint@latest

# Using Homebrew
brew install actionlint

# Download binary from releases
# See: https://github.com/rhysd/actionlint/releases
```

## Integration Notes

- Works with GitHub Actions workflows in `.github/workflows/` directory
- Supports YAML workflow files (`.yml` or `.yaml`)
- Can integrate with shellcheck for shell script validation
- Can integrate with pyflakes for Python script validation
- Provides detailed error messages with file locations

## Architecture

The Ship CLI provides a unified interface that combines:

1. **Dagger Modules**: Containerized execution of actionlint with advanced options
2. **CLI Commands**: Direct access to enhanced functionality via command-line flags
3. **MCP Server**: AI assistant integration with discrete tool functions

### Dagger Module Functions
- `ScanDirectory()`: Basic directory scanning
- `ScanDirectoryWithOptions()`: Advanced scanning with format, ignore patterns, and color
- `ScanWithExternalTools()`: Integration with shellcheck and pyflakes
- `ScanSpecificFiles()`: Scan specific workflow files with options
- `GetVersion()`: Get actionlint version information

This architecture ensures consistent behavior whether using the CLI directly, through MCP tools, or in automated workflows.

## Benefits of Unified Architecture

1. **Consistency**: Same functionality available through CLI, MCP, and programmatic access
2. **Containerization**: All tools run in isolated Docker containers via Dagger
3. **Flexibility**: Choose between direct CLI usage or AI-assisted MCP tools
4. **Maintainability**: Single source of truth for tool functionality
5. **Integration**: Seamless workflow between CLI commands and MCP server

## Exit Codes

- `0`: No problems found
- `1`: Problems found  
- `2`: Invalid command line option
- `3`: Fatal error during execution

## References

- **Official Repository**: https://github.com/rhysd/actionlint
- **Documentation**: https://github.com/rhysd/actionlint/blob/main/docs/usage.md
- **Project Website**: https://rhysd.github.io/actionlint/