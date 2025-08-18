# Fleet

GitOps for Kubernetes multi-cluster deployments.

## Description

Fleet is a GitOps toolkit for managing Kubernetes deployments from Git repositories. It enables declarative cluster management and multi-cluster operations using custom resources and kubectl commands. Fleet is designed to handle large-scale Kubernetes deployments across multiple clusters.

## MCP Tools

### GitRepo Management
- **`fleet_apply_gitrepo`** - Apply Fleet GitRepo configuration using kubectl
- **`fleet_get_gitrepos`** - Get Fleet GitRepos status
- **`fleet_describe_gitrepo`** - Describe Fleet GitRepo details

### Bundle Management
- **`fleet_get_bundles`** - Get Fleet Bundles status across clusters
- **`fleet_get_bundledeployments`** - Get Fleet BundleDeployments status

### Installation
- **`fleet_install`** - Install Fleet using Helm

## Real CLI Commands Used

### kubectl Commands
- `kubectl apply -f gitrepo.yaml` - Deploy GitRepo configuration
- `kubectl get gitrepo -n fleet-local` - List GitRepos
- `kubectl get bundles --all-namespaces` - List all bundles
- `kubectl get bundledeployments --all-namespaces` - List bundle deployments
- `kubectl describe gitrepo <name> -n fleet-local` - Describe specific GitRepo

### Helm Commands
- `helm -n cattle-fleet-system install --create-namespace --wait fleet-crd <chart-url>` - Install Fleet CRDs
- `helm -n cattle-fleet-system install --create-namespace --wait fleet <chart-url>` - Install Fleet

## Core Concepts

### GitRepo
Defines a Git repository to monitor for Kubernetes manifests:
```yaml
apiVersion: fleet.cattle.io/v1alpha1
kind: GitRepo
metadata:
  name: sample
  namespace: fleet-local
spec:
  repo: https://github.com/rancher/fleet-examples
  paths:
  - simple
```

### Bundles
Packaged sets of resources from GitRepos, deployed to target clusters.

### BundleDeployments
Instances of bundles deployed to specific clusters.

### Clusters
Target Kubernetes clusters managed by Fleet.

## Use Cases

### Multi-Cluster GitOps
- Deploy applications across multiple Kubernetes clusters
- Centralized configuration management
- Consistent application rollouts
- Cross-cluster policy enforcement

### Edge Computing
- Manage thousands of edge clusters
- Lightweight cluster agents
- Disconnected operation support
- Efficient resource distribution

### CI/CD Integration
- Git-based deployment workflows
- Automated application delivery
- Environment promotion pipelines
- Configuration drift detection

### Compliance & Security
- Declarative security policies
- Audit trail for all changes
- Compliance verification
- Security posture management

## Fleet Architecture

### Control Plane
- Fleet controller in management cluster
- GitRepo monitoring and processing
- Bundle generation and distribution
- Cluster registration and management

### Agents
- Lightweight agents on target clusters
- Bundle deployment and reconciliation
- Status reporting to control plane
- Local resource management

## Configuration Examples

### Simple Application Deployment
```yaml
apiVersion: fleet.cattle.io/v1alpha1
kind: GitRepo
metadata:
  name: webapp
  namespace: fleet-local
spec:
  repo: https://github.com/company/webapp-config
  branch: production
  paths:
  - manifests/
  targets:
  - name: production
    clusterSelector:
      matchLabels:
        env: production
```

### Multi-Environment Deployment
```yaml
apiVersion: fleet.cattle.io/v1alpha1
kind: GitRepo
metadata:
  name: microservice
  namespace: fleet-local
spec:
  repo: https://github.com/company/microservice
  paths:
  - deploy/base
  - deploy/environments/production
  targets:
  - name: production-us
    clusterSelector:
      matchLabels:
        env: production
        region: us-east-1
  - name: production-eu
    clusterSelector:
      matchLabels:
        env: production
        region: eu-west-1
```

## Monitoring and Troubleshooting

### Status Monitoring
- GitRepo status shows sync state
- Bundle status shows deployment progress
- BundleDeployment status shows cluster-specific state
- Event logs for troubleshooting

### Common Operations
- Check GitRepo sync status
- Monitor bundle deployment progress
- Troubleshoot deployment failures
- View cluster-specific resource state

## Integration

Works with any Kubernetes cluster and integrates with CI/CD pipelines, monitoring systems, and security tools. Supports Helm charts, Kustomize, and raw Kubernetes manifests.