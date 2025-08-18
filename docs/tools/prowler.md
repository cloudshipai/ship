# Prowler MCP Tool

Prowler is a multi-cloud security assessment tool that performs security auditing and compliance checks across AWS, Azure, GCP, and Kubernetes environments.

## Description

Prowler provides comprehensive security assessments by:
- Scanning cloud infrastructure for security misconfigurations
- Checking compliance against various frameworks (CIS, PCI-DSS, NIST, etc.)
- Supporting multiple cloud providers and platforms
- Generating detailed reports in various formats
- Offering continuous monitoring capabilities

## MCP Functions

### `prowler_aws`
Scan AWS account for security issues using real prowler CLI.

**Parameters:**
- `profile`: AWS CLI profile to use
- `regions`: Comma-separated list of AWS regions to scan
- `checks`: Comma-separated list of specific checks to run
- `services`: Comma-separated list of AWS services to scan
- `compliance`: Compliance framework to check against
- `output_formats`: Output formats (csv, json-asff, json-ocsf, html)
- `output_directory`: Output directory for results

**CLI Command:** `prowler aws [--profile <profile>] [--regions <regions>] [options]`

### `prowler_azure`
Scan Azure subscription for security issues using real prowler CLI.

**Parameters:**
- `subscription_id`: Azure subscription ID to scan
- `checks`: Comma-separated list of specific checks to run
- `services`: Comma-separated list of Azure services to scan
- `compliance`: Compliance framework to check against
- `output_formats`: Output formats (csv, json-asff, json-ocsf, html)
- `output_directory`: Output directory for results

**CLI Command:** `prowler azure [--subscription-id <id>] [options]`

### `prowler_gcp`
Scan GCP project for security issues using real prowler CLI.

**Parameters:**
- `project_id`: GCP project ID to scan
- `checks`: Comma-separated list of specific checks to run
- `services`: Comma-separated list of GCP services to scan
- `compliance`: Compliance framework to check against
- `output_formats`: Output formats (csv, json-asff, json-ocsf, html)
- `output_directory`: Output directory for results

**CLI Command:** `prowler gcp [--project-id <id>] [options]`

### `prowler_kubernetes`
Scan Kubernetes cluster for security issues using real prowler CLI.

**Parameters:**
- `kubeconfig_path`: Path to kubeconfig file
- `context`: Kubernetes context to use
- `checks`: Comma-separated list of specific checks to run
- `compliance`: Compliance framework to check against
- `output_formats`: Output formats (csv, json-asff, json-ocsf, html)
- `output_directory`: Output directory for results

**CLI Command:** `prowler kubernetes [--kubeconfig <path>] [--context <context>] [options]`

### `prowler_list_checks`
List available security checks using real prowler CLI.

**Parameters:**
- `provider`: Cloud provider (aws, azure, gcp, kubernetes)
- `service`: Specific service to list checks for
- `compliance`: Filter by compliance framework

**CLI Command:** `prowler [<provider>] --list-checks [options]`

### `prowler_list_services`
List available services using real prowler CLI.

**Parameters:**
- `provider`: Cloud provider (aws, azure, gcp, kubernetes)

**CLI Command:** `prowler [<provider>] --list-services`

### `prowler_list_compliance`
List available compliance frameworks using real prowler CLI.

**Parameters:**
- `provider`: Cloud provider (aws, azure, gcp, kubernetes)

**CLI Command:** `prowler [<provider>] --list-compliance`

### `prowler_dashboard`
Launch Prowler dashboard using real prowler CLI.

**Parameters:**
- `port`: Port to run dashboard on (default: 11666)
- `host`: Host to bind dashboard to (default: 127.0.0.1)

**CLI Command:** `prowler dashboard [--port <port>] [--host <host>]`

## Supported Output Formats

- **csv**: Comma-separated values
- **json-asff**: AWS Security Finding Format
- **json-ocsf**: Open Cybersecurity Schema Framework
- **html**: HTML report format

## Common Use Cases

1. **Multi-Cloud Security Assessment**: Use provider-specific tools to scan AWS, Azure, GCP, and Kubernetes
2. **Compliance Auditing**: Run compliance checks against frameworks like CIS, PCI-DSS, NIST
3. **Service-Specific Scans**: Focus on specific cloud services or resources
4. **Continuous Monitoring**: Integrate into CI/CD pipelines for ongoing security validation
5. **Discovery**: Use list functions to explore available checks, services, and compliance frameworks
6. **Dashboard Monitoring**: Launch web-based dashboard for interactive analysis

## Integration with Ship CLI

All Prowler tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated multi-cloud security assessments
- Generate compliance reports across different frameworks
- Identify security misconfigurations and vulnerabilities
- Provide continuous security monitoring capabilities

The tools use containerized execution via Dagger for consistent, isolated security assessment environments.