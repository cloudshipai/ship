# Goldilocks

Kubernetes resource recommendations using Vertical Pod Autoscaler (VPA).

## Description

Goldilocks is a tool that creates a Vertical Pod Autoscaler (VPA) for each workload in a namespace and then queries the VPA for resource recommendations. It provides a dashboard to view recommendations for resource requests and limits, helping you find the "just right" resource settings for your deployments.

## MCP Tools

### Installation & Management
- **`goldilocks_install_helm`** - Install Goldilocks using Helm chart
- **`goldilocks_enable_namespace`** - Enable Goldilocks monitoring for a namespace
- **`goldilocks_dashboard`** - Port-forward to access Goldilocks dashboard
- **`goldilocks_uninstall`** - Uninstall Goldilocks using Helm

### Resource Analysis
- **`goldilocks_get_recommendations`** - Get VPA recommendations from enabled namespace

## Real CLI Commands Used

### Installation Commands
- `helm repo add fairwinds-stable https://charts.fairwinds.com/stable` - Add Fairwinds Helm repository
- `kubectl create namespace goldilocks --dry-run=client -o yaml | kubectl apply -f -` - Create namespace
- `helm install goldilocks --namespace goldilocks fairwinds-stable/goldilocks` - Install via Helm

### Configuration Commands
- `kubectl label ns <namespace> goldilocks.fairwinds.com/enabled=true --overwrite` - Enable namespace monitoring
- `kubectl -n goldilocks port-forward svc/goldilocks-dashboard 8080:80` - Access dashboard

### Analysis Commands
- `kubectl get vpa -n <namespace> -o yaml` - Get VPA recommendations
- `helm uninstall goldilocks --namespace goldilocks` - Uninstall

## Core Concepts

### Vertical Pod Autoscaler (VPA)
Goldilocks leverages Kubernetes VPA to analyze resource usage patterns:
- **Recommendation Mode**: VPA analyzes workloads without modifying them
- **Resource Monitoring**: Tracks CPU and memory usage over time
- **Right-sizing Suggestions**: Provides optimal request and limit recommendations

### Namespace Enablement
Enable monitoring by labeling namespaces:
```bash
kubectl label ns production goldilocks.fairwinds.com/enabled=true
```

### Dashboard Interface
Web-based dashboard provides:
- Visual representation of recommendations
- Current vs. recommended resource settings
- Historical usage patterns
- Cost impact analysis

## Use Cases

### Resource Optimization
- **Right-sizing Workloads**: Find optimal CPU/memory requests and limits
- **Cost Reduction**: Eliminate over-provisioned resources
- **Performance Improvement**: Prevent resource-constrained applications
- **Capacity Planning**: Understand cluster resource requirements

### Development Workflow
- **Development Environment**: Optimize staging environment resources
- **Production Tuning**: Fine-tune production workload efficiency
- **Migration Planning**: Size workloads for cluster migrations
- **Multi-environment Consistency**: Standardize resource settings

### Operational Excellence
- **Resource Governance**: Establish resource allocation standards
- **Monitoring Integration**: Track resource utilization trends
- **Automation Enablement**: Integrate with GitOps workflows
- **Team Collaboration**: Share recommendations across teams

## Prerequisites

### Cluster Requirements
- **Kubernetes v1.16+**: Minimum cluster version
- **Metrics Server**: Required for resource utilization data
- **VPA**: Vertical Pod Autoscaler must be installed
- **RBAC**: Appropriate permissions for VPA and metrics access

### Installation Prerequisites
```bash
# Install metrics-server (if not present)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Install VPA (if not present)
git clone https://github.com/kubernetes/autoscaler.git
cd autoscaler/vertical-pod-autoscaler/
./hack/vpa-up.sh
```

## Configuration Examples

### Namespace Enablement
```bash
# Enable multiple namespaces
kubectl label ns production goldilocks.fairwinds.com/enabled=true
kubectl label ns staging goldilocks.fairwinds.com/enabled=true
kubectl label ns development goldilocks.fairwinds.com/enabled=true
```

### Custom Installation
```yaml
# values.yaml for Helm installation
dashboard:
  enabled: true
  service:
    type: ClusterIP
    port: 80

vpa:
  enabled: true
  updater:
    enabled: false  # Recommendation mode only

controller:
  resources:
    requests:
      cpu: 25m
      memory: 32Mi
    limits:
      cpu: 100m
      memory: 128Mi
```

### VPA Configuration
```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: my-app-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  updatePolicy:
    updateMode: "Off"  # Recommendations only
```

## Dashboard Features

### Resource Recommendations
- **Current Settings**: View existing requests and limits
- **VPA Recommendations**: See optimal resource suggestions
- **Usage Patterns**: Historical utilization graphs
- **Cost Impact**: Estimated cost changes

### Workload Analysis
- **Pod-level Details**: Individual container recommendations
- **Deployment Overview**: Aggregate workload analysis
- **Resource Efficiency**: Over/under-provisioning identification
- **Trend Analysis**: Long-term usage patterns

### Export Capabilities
- **YAML Export**: Export recommendations as Kubernetes manifests
- **CSV Reports**: Resource data in spreadsheet format
- **API Access**: Programmatic access to recommendations
- **Integration Support**: Connect with external tools

## Integration Patterns

### GitOps Workflow
```bash
# Get recommendations and commit to GitOps repository
kubectl get vpa -n production -o yaml > vpa-recommendations.yaml
git add vpa-recommendations.yaml
git commit -m "Update VPA recommendations"
```

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Get Resource Recommendations
  run: |
    kubectl get vpa -n ${{ env.NAMESPACE }} -o yaml > recommendations.yaml
    # Process recommendations and update deployment manifests
```

### Monitoring Integration
```bash
# Export metrics for Prometheus
kubectl get vpa -n production -o json | \
  jq '.items[] | {name: .metadata.name, cpu: .status.recommendation.containerRecommendations[0].target.cpu}'
```

## Best Practices

### Implementation Strategy
- **Gradual Rollout**: Start with non-production environments
- **Baseline Establishment**: Collect data before making changes
- **Testing Validation**: Verify recommendations in staging
- **Incremental Updates**: Apply changes gradually

### Monitoring and Validation
- **Performance Tracking**: Monitor application performance after changes
- **Resource Utilization**: Verify actual usage matches recommendations
- **Cost Analysis**: Track cost impact of optimizations
- **Rollback Planning**: Maintain ability to revert changes

### Operational Considerations
- **Regular Reviews**: Periodically review and update recommendations
- **Seasonal Patterns**: Account for traffic and usage variations
- **Multi-cluster Consistency**: Standardize across environments
- **Documentation**: Maintain records of optimization decisions

Goldilocks provides data-driven insights for Kubernetes resource optimization, helping teams achieve better performance, cost efficiency, and operational excellence.