# Kyverno MCP Tool

Kyverno is a policy engine designed for Kubernetes that can validate, mutate, and generate configurations using policies as Kubernetes resources.

## Description

Kyverno provides:
- Policy-based resource validation
- Automatic resource mutation
- Resource generation from policies
- Fine-grained RBAC controls
- Native Kubernetes integration (no new language to learn)

## MCP Functions

### `kyverno_apply_policy`
Apply Kyverno policy using real kyverno CLI.

**Parameters:**
- `policy_file` (required): Path to policy YAML file
- `namespace`: Target namespace
- `cluster`: Apply as cluster policy

**CLI Command:** `kyverno apply <policy_file> [--namespace <namespace>] [--cluster]`

### `kyverno_test_policy`
Test Kyverno policy using real kyverno CLI.

**Parameters:**
- `policy_file` (required): Path to policy YAML file
- `resource_file` (required): Path to resource YAML file
- `values_file`: Path to values file for variable substitution

**CLI Command:** `kyverno test <policy_file> --resource <resource_file> [--values <values_file>]`

### `kyverno_validate_policy`
Validate Kyverno policy syntax using real kyverno CLI.

**Parameters:**
- `policy_file` (required): Path to policy YAML file
- `cluster_resources`: Path to cluster resources file

**CLI Command:** `kyverno validate <policy_file> [--cluster-resources <cluster_resources>]`

### `kyverno_generate_report`
Generate compliance report using kubectl.

**Parameters:**
- `namespace`: Namespace to report on
- `output_format`: Output format (yaml, json)

**CLI Command:** `kubectl get policyreport,clusterpolicyreport [-n <namespace>] [-o <format>]`

### `kyverno_jp_query`
Execute JMESPath query for policy testing using kyverno CLI.

**Parameters:**
- `query` (required): JMESPath query string
- `input_file` (required): Input JSON/YAML file
- `expression_file`: File containing JMESPath expression

**CLI Command:** `kyverno jp query "<query>" -i <input_file> [-f <expression_file>]`

### `kyverno_version`
Get Kyverno version information.

**CLI Command:** `kyverno version`

### `kyverno_install`
Install Kyverno in cluster using kubectl.

**Parameters:**
- `namespace`: Installation namespace (default: kyverno)
- `version`: Kyverno version to install

**CLI Command:** `kubectl create -f https://github.com/kyverno/kyverno/releases/download/<version>/install.yaml`

### `kyverno_create_exception`
Create policy exception using kubectl.

**Parameters:**
- `exception_file` (required): Path to exception YAML file
- `namespace`: Target namespace

**CLI Command:** `kubectl apply -f <exception_file> [-n <namespace>]`

## Common Use Cases

1. **Resource Validation**: Enforce standards on Kubernetes resources
2. **Automatic Mutation**: Modify resources to meet requirements
3. **Resource Generation**: Auto-create resources based on policies
4. **Compliance Reporting**: Track policy violations
5. **Security Enforcement**: Apply security policies cluster-wide

## Integration with Ship CLI

All Kyverno tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Deploy and manage Kubernetes policies
- Validate resource configurations
- Generate compliance reports
- Test policies before deployment

The tools use containerized execution via Dagger for consistent, isolated policy management.