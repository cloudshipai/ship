# CloudQuery MCP Tool

CloudQuery is an open-source cloud asset inventory tool that extracts, transforms, and loads cloud infrastructure configuration and metadata.

## Description

CloudQuery provides:
- Multi-cloud asset inventory and analysis
- Infrastructure data extraction to various destinations
- SQL-based querying of cloud resources
- Compliance and security posture assessment
- Asset relationship mapping

## MCP Functions

### `cloudquery_sync`
Sync cloud resources using real cloudquery CLI.

**Parameters:**
- `config` (required): Path to CloudQuery configuration file
- `log_level`: Logging level (trace, debug, info, warn, error)
- `log_format`: Log format (text, json)

**CLI Command:** `cloudquery sync <config> [--log-level <level>] [--log-format <format>]`

### `cloudquery_sync_aws`
Sync AWS resources using real cloudquery CLI with AWS source.

**Parameters:**
- `config` (required): Path to CloudQuery configuration file
- `tables`: Comma-separated list of tables to sync
- `skip_tables`: Comma-separated list of tables to skip

**CLI Command:** `cloudquery sync <config> [--tables <tables>] [--skip-tables <skip_tables>]`

### `cloudquery_sync_azure`
Sync Azure resources using real cloudquery CLI with Azure source.

**Parameters:**
- `config` (required): Path to CloudQuery configuration file
- `tables`: Comma-separated list of tables to sync
- `skip_tables`: Comma-separated list of tables to skip

**CLI Command:** `cloudquery sync <config> [--tables <tables>] [--skip-tables <skip_tables>]`

### `cloudquery_sync_gcp`
Sync GCP resources using real cloudquery CLI with GCP source.

**Parameters:**
- `config` (required): Path to CloudQuery configuration file
- `tables`: Comma-separated list of tables to sync
- `skip_tables`: Comma-separated list of tables to skip

**CLI Command:** `cloudquery sync <config> [--tables <tables>] [--skip-tables <skip_tables>]`

### `cloudquery_sync_kubernetes`
Sync Kubernetes resources using real cloudquery CLI with K8s source.

**Parameters:**
- `config` (required): Path to CloudQuery configuration file
- `tables`: Comma-separated list of tables to sync

**CLI Command:** `cloudquery sync <config> [--tables <tables>]`

### `cloudquery_migrate`
Run database migrations using real cloudquery CLI.

**Parameters:**
- `config` (required): Path to CloudQuery configuration file
- `force`: Force migration even if losing data

**CLI Command:** `cloudquery migrate <config> [--force]`

### `cloudquery_plugin_install`
Install CloudQuery plugin using real cloudquery CLI.

**Parameters:**
- `path` (required): Plugin path (e.g., cloudquery/aws)
- `version`: Plugin version

**CLI Command:** `cloudquery plugin install <path> [--version <version>]`

### `cloudquery_plugin_list`
List installed CloudQuery plugins using real cloudquery CLI.

**CLI Command:** `cloudquery plugin list`

### `cloudquery_tables`
List available tables for a source using real cloudquery CLI.

**Parameters:**
- `source` (required): Source plugin name
- `version`: Source plugin version
- `format`: Output format (json, csv, markdown)

**CLI Command:** `cloudquery tables <source> [--version <version>] [--format <format>]`

### `cloudquery_version`
Get CloudQuery version information.

**CLI Command:** `cloudquery --version`

## Common Use Cases

1. **Cloud Asset Inventory**: Complete inventory of cloud resources
2. **Compliance Monitoring**: Track compliance across cloud environments
3. **Security Analysis**: Identify security misconfigurations
4. **Cost Analysis**: Analyze cloud resource usage and costs
5. **Multi-cloud Management**: Unified view across cloud providers

## Integration with Ship CLI

All CloudQuery tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated cloud asset discovery
- Generate infrastructure inventories
- Analyze cloud security posture
- Export data to various destinations

The tools use containerized execution via Dagger for consistent, isolated cloud querying.