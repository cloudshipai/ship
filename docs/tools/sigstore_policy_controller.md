# Sigstore Policy Controller MCP Tool

Sigstore Policy Controller is a Kubernetes admission controller that enforces image signature policies using real policy-controller-tester and kubectl commands.

## Description

Sigstore Policy Controller provides:
- Kubernetes admission control for container images
- Image signature verification enforcement
- Policy-based security controls
- ClusterImagePolicy resource management
- Namespace-level policy enforcement configuration

## MCP Functions

### `sigstore_test_policy`
Test Sigstore policy against container image using real policy-controller-tester.

**Parameters:**
- `policy` (required): Path to ClusterImagePolicy file or URL
- `image` (required): Container image to test against policy
- `resource`: Path to Kubernetes resource file
- `trustroot`: Path to Kubernetes TrustRoot resource
- `log_level`: Log level for output (debug, info, warn, error)

**CLI Command:** `policy-controller-tester -policy <policy> -image <image> [-resource <resource>] [-trustroot <trustroot>] [-log-level <level>]`

### `sigstore_tester_version`
Get policy-controller-tester version.

**CLI Command:** `policy-controller-tester -version`

### `sigstore_create_policy`
Create ClusterImagePolicy using kubectl.

**Parameters:**
- `policy_file` (required): Path to ClusterImagePolicy YAML file
- `namespace`: Namespace for the policy (optional)

**CLI Command:** `kubectl apply -f <policy_file> [-n <namespace>]`

### `sigstore_list_policies`
List ClusterImagePolicies using kubectl.

**Parameters:**
- `output`: Output format (yaml, json, wide, name)

**CLI Command:** `kubectl get clusterimagepolicy [-o <output>]`

### `sigstore_delete_policy`
Delete ClusterImagePolicy using kubectl.

**Parameters:**
- `policy_name` (required): Name of the ClusterImagePolicy to delete

**CLI Command:** `kubectl delete clusterimagepolicy <policy_name>`

### `sigstore_enable_namespace`
Enable Sigstore policy enforcement for namespace using kubectl.

**Parameters:**
- `namespace` (required): Namespace to enable policy enforcement

**CLI Command:** `kubectl label namespace <namespace> policy.sigstore.dev/include=true`

### `sigstore_disable_namespace`
Disable Sigstore policy enforcement for namespace using kubectl.

**Parameters:**
- `namespace` (required): Namespace to disable policy enforcement

**CLI Command:** `kubectl label namespace <namespace> policy.sigstore.dev/exclude=true`

### `sigstore_get_namespace_status`
Get Sigstore policy enforcement status for namespace using kubectl.

**Parameters:**
- `namespace` (required): Namespace to check policy enforcement status

**CLI Command:** `kubectl get namespace <namespace> --show-labels`

### `sigstore_describe_policy`
Describe ClusterImagePolicy using kubectl.

**Parameters:**
- `policy_name` (required): Name of the ClusterImagePolicy to describe

**CLI Command:** `kubectl describe clusterimagepolicy <policy_name>`

## Common Use Cases

1. **Image Verification**: Enforce signature verification on container images
2. **Policy Testing**: Test policies before deployment using policy-controller-tester
3. **Namespace Management**: Configure policy enforcement per namespace
4. **Compliance Enforcement**: Ensure only signed images are deployed
5. **Policy Management**: Create, update, and delete image verification policies

## Integration with Ship CLI

All Sigstore Policy Controller tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Test image policies before deployment
- Manage ClusterImagePolicy resources
- Configure namespace-level enforcement
- Verify policy compliance

The tools use containerized execution via Dagger for consistent, isolated policy operations.