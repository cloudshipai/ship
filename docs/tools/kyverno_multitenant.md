# Kyverno Multi-Tenant

Multi-tenant Kubernetes policy management using Kyverno policies, kubectl, and standard Kubernetes resources.

## Description

Kyverno Multi-Tenant provides comprehensive tools for implementing multi-tenancy in Kubernetes clusters using Kyverno policies combined with standard Kubernetes primitives like namespaces, ResourceQuotas, and RBAC. This approach enables organizations to safely share Kubernetes clusters among multiple tenants while ensuring proper isolation, resource limits, and security policies. The implementation leverages Kyverno's policy engine to automatically enforce tenant boundaries and generate necessary resources.

## MCP Tools

### Tenant Namespace Management
- **`kyverno_multitenant_create_namespace`** - Create tenant namespace with appropriate labels
- **`kyverno_multitenant_list_namespaces`** - List namespaces for a specific tenant

### Resource Management  
- **`kyverno_multitenant_create_quota`** - Create ResourceQuota for tenant namespace
- **`kyverno_multitenant_namespace_isolation`** - Apply namespace isolation policies

### Policy Management
- **`kyverno_multitenant_generate_policy`** - Apply Kyverno generate policies for automatic resource creation
- **`kyverno_multitenant_get_policies`** - Get Kyverno policies affecting tenant namespaces

## Real CLI Commands Used

### Namespace Management
- `kubectl create namespace <namespace>` - Create tenant namespace
- `kubectl label namespace <namespace> tenant=<name>` - Label namespace with tenant info
- `kubectl get namespaces -l tenant=<name>` - List tenant namespaces

### Resource Quota Management
- `echo '<yaml>' | kubectl apply -f -` - Apply inline ResourceQuota YAML
- `kubectl get resourcequota -n <namespace>` - Check resource quotas

### Policy Management
- `kubectl apply -f <policy-file>` - Apply Kyverno policies
- `kubectl get policy -n <namespace>` - List namespace policies
- `kubectl get clusterpolicy` - List cluster-wide policies
- `kyverno apply <policy> --dry-run` - Validate policies before applying

### Resource Validation
- `kyverno apply <policy-file> --dry-run` - Validate Kyverno policies
- `kubectl apply -f <file> --dry-run=client` - Validate Kubernetes resources

## Use Cases

### Tenant Isolation
- **Namespace Isolation**: Separate tenant workloads into dedicated namespaces
- **Resource Boundaries**: Enforce CPU, memory, and storage limits per tenant
- **Network Isolation**: Implement network policies for tenant separation
- **Security Boundaries**: Apply security policies specific to tenant requirements

### Resource Management
- **Quota Enforcement**: Automatically apply resource quotas to tenant namespaces
- **Limit Ranges**: Set default and maximum resource limits for containers
- **Storage Classes**: Provide tenant-specific storage classes and policies
- **Compute Resources**: Manage CPU, memory, and GPU allocation per tenant

### Policy Automation
- **Automatic Policy Application**: Use Kyverno generate rules to create resources
- **Compliance Enforcement**: Ensure tenant workloads comply with organizational policies
- **Security Standards**: Apply consistent security policies across tenants
- **Configuration Drift Prevention**: Prevent unauthorized changes to tenant configurations

### Operational Management
- **Tenant Onboarding**: Streamline process of adding new tenants
- **Resource Monitoring**: Track resource usage per tenant
- **Policy Auditing**: Audit policy compliance across all tenants
- **Lifecycle Management**: Manage tenant creation, updates, and deletion

## Configuration Examples

### Basic Tenant Setup
```bash
# Create tenant namespace with labels
kubectl create namespace tenant-acme
kubectl label namespace tenant-acme tenant=acme
kubectl label namespace tenant-acme kyverno.io/tenant=acme

# Create resource quota for tenant
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ResourceQuota
metadata:
  name: tenant-quota
  namespace: tenant-acme
spec:
  hard:
    requests.cpu: "4"
    limits.cpu: "8"
    requests.memory: 8Gi
    limits.memory: 16Gi
    pods: "10"
EOF

# List tenant namespaces
kubectl get namespaces -l tenant=acme
```

### Namespace Isolation Policies
```bash
# Apply network isolation policy
cat <<EOF | kubectl apply -f -
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: tenant-namespace-isolation
spec:
  validationFailureAction: enforce
  background: true
  rules:
  - name: deny-cross-tenant-access
    match:
      any:
      - resources:
          kinds:
          - Pod
    validate:
      message: "Pods cannot access resources in other tenant namespaces"
      pattern:
        spec:
          =(serviceAccount): "!tenant-*"
EOF

# Apply pod security standards per tenant
cat <<EOF | kubectl apply -f -
apiVersion: kyverno.io/v1
kind: Policy
metadata:
  name: tenant-pod-security
  namespace: tenant-acme
spec:
  validationFailureAction: enforce
  rules:
  - name: require-non-root
    match:
      any:
      - resources:
          kinds:
          - Pod
    validate:
      message: "Pods must run as non-root user"
      pattern:
        spec:
          securityContext:
            runAsNonRoot: true
EOF
```

### Automatic Resource Generation
```bash
# Apply generate policy for automatic NetworkPolicy creation
cat <<EOF | kubectl apply -f -
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: generate-tenant-network-policy
spec:
  rules:
  - name: generate-netpol
    match:
      any:
      - resources:
          kinds:
          - Namespace
    generate:
      kind: NetworkPolicy
      name: deny-all
      namespace: "{{request.object.metadata.name}}"
      synchronize: true
      data:
        spec:
          podSelector: {}
          policyTypes:
          - Ingress
          - Egress
EOF

# Apply generate policy for automatic LimitRange
cat <<EOF | kubectl apply -f -
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: generate-tenant-limit-range
spec:
  rules:
  - name: generate-limit-range
    match:
      any:
      - resources:
          kinds:
          - Namespace
          name: "tenant-*"
    generate:
      kind: LimitRange
      name: tenant-limits
      namespace: "{{request.object.metadata.name}}"
      data:
        spec:
          limits:
          - default:
              cpu: 100m
              memory: 128Mi
            defaultRequest:
              cpu: 50m
              memory: 64Mi
            type: Container
EOF
```

## Advanced Usage

### Comprehensive Multi-Tenant Setup
```bash
#!/bin/bash
# setup-multi-tenant-cluster.sh

TENANTS=("acme" "globex" "initech")

echo "Setting up multi-tenant Kubernetes cluster with Kyverno..."

# Install Kyverno if not present
if ! kubectl get namespace kyverno > /dev/null 2>&1; then
    echo "Installing Kyverno..."
    helm repo add kyverno https://kyverno.github.io/kyverno/
    helm repo update
    helm install kyverno kyverno/kyverno --namespace kyverno --create-namespace --wait
fi

# Apply base cluster policies
echo "Applying base cluster policies..."
cat <<EOF | kubectl apply -f -
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: tenant-namespace-isolation
spec:
  validationFailureAction: enforce
  background: true
  rules:
  - name: require-tenant-label
    match:
      any:
      - resources:
          kinds:
          - Namespace
          names:
          - "tenant-*"
    validate:
      message: "Tenant namespaces must have tenant label"
      pattern:
        metadata:
          labels:
            tenant: "?*"
  - name: deny-cross-tenant-service-access
    match:
      any:
      - resources:
          kinds:
          - Service
    validate:
      message: "Services cannot expose endpoints outside tenant namespace"
      pattern:
        spec:
          =(externalName): "!*.tenant-*"
EOF

# Set up each tenant
for tenant in "${TENANTS[@]}"; do
    echo "Setting up tenant: $tenant"
    
    # Create tenant namespace
    kubectl create namespace "tenant-$tenant" --dry-run=client -o yaml | kubectl apply -f -
    
    # Label namespace
    kubectl label namespace "tenant-$tenant" tenant="$tenant" --overwrite
    kubectl label namespace "tenant-$tenant" kyverno.io/tenant="$tenant" --overwrite
    
    # Create resource quota
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ResourceQuota
metadata:
  name: tenant-quota
  namespace: tenant-$tenant
spec:
  hard:
    requests.cpu: "2"
    limits.cpu: "4"
    requests.memory: 4Gi
    limits.memory: 8Gi
    pods: "20"
    services: "10"
    persistentvolumeclaims: "5"
EOF
    
    # Create tenant-specific RBAC
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tenant-admin
  namespace: tenant-$tenant
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: tenant-admin
  namespace: tenant-$tenant
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: tenant-admin
  namespace: tenant-$tenant
subjects:
- kind: ServiceAccount
  name: tenant-admin
  namespace: tenant-$tenant
roleRef:
  kind: Role
  name: tenant-admin
  apiGroup: rbac.authorization.k8s.io
EOF
    
    # Create tenant-specific policies
    cat <<EOF | kubectl apply -f -
apiVersion: kyverno.io/v1
kind: Policy
metadata:
  name: tenant-pod-security
  namespace: tenant-$tenant
spec:
  validationFailureAction: enforce
  rules:
  - name: require-security-context
    match:
      any:
      - resources:
          kinds:
          - Pod
    validate:
      message: "Pods must define security context"
      pattern:
        spec:
          securityContext:
            runAsNonRoot: true
            runAsUser: ">0"
  - name: deny-privileged
    match:
      any:
      - resources:
          kinds:
          - Pod
    validate:
      message: "Privileged pods are not allowed"
      pattern:
        spec:
          =(securityContext):
            =(privileged): "false"
          containers:
          - securityContext:
              =(privileged): "false"
EOF
    
    echo "Tenant $tenant setup complete"
done

echo "Multi-tenant cluster setup completed!"
echo "Tenants created: ${TENANTS[*]}"
```

### Tenant Resource Monitoring
```bash
#!/bin/bash
# monitor-tenant-resources.sh

echo "=== Multi-Tenant Resource Monitoring ==="
echo "Date: $(date)"
echo ""

# Get all tenant namespaces
TENANTS=$(kubectl get namespaces -l tenant --no-headers -o custom-columns=":metadata.labels.tenant" | sort | uniq)

for tenant in $TENANTS; do
    echo "=== Tenant: $tenant ==="
    
    # Get tenant namespaces
    NAMESPACES=$(kubectl get namespaces -l tenant=$tenant --no-headers -o custom-columns=":metadata.name")
    echo "Namespaces: $NAMESPACES"
    
    for ns in $NAMESPACES; do
        echo "  Namespace: $ns"
        
        # Resource quota status
        echo "    Resource Quotas:"
        kubectl get resourcequota -n $ns -o custom-columns="NAME:.metadata.name,CPU-USED:.status.used.requests\.cpu,CPU-LIMIT:.status.hard.requests\.cpu,MEMORY-USED:.status.used.requests\.memory,MEMORY-LIMIT:.status.hard.requests\.memory,PODS:.status.used.pods" --no-headers 2>/dev/null | sed 's/^/      /' || echo "      No resource quotas"
        
        # Pod count and status
        POD_COUNT=$(kubectl get pods -n $ns --no-headers 2>/dev/null | wc -l)
        RUNNING_PODS=$(kubectl get pods -n $ns --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l)
        echo "    Pods: $RUNNING_PODS/$POD_COUNT running"
        
        # Service count
        SERVICE_COUNT=$(kubectl get services -n $ns --no-headers 2>/dev/null | wc -l)
        echo "    Services: $SERVICE_COUNT"
        
        # PVC count
        PVC_COUNT=$(kubectl get pvc -n $ns --no-headers 2>/dev/null | wc -l)
        echo "    PVCs: $PVC_COUNT"
        
        echo ""
    done
    
    # Check policy violations
    echo "  Policy Violations:"
    kubectl get events -A --field-selector reason=PolicyViolation,involvedObject.namespace=$ns-* --no-headers 2>/dev/null | wc -l | sed 's/^/    /' || echo "    0"
    
    echo ""
done

# Cluster-wide policy status
echo "=== Cluster Policy Status ==="
CLUSTER_POLICIES=$(kubectl get clusterpolicy --no-headers -o custom-columns=":metadata.name" 2>/dev/null)
for policy in $CLUSTER_POLICIES; do
    echo "Policy: $policy"
    kubectl get clusterpolicy $policy -o jsonpath='{.status.ready}' 2>/dev/null | sed 's/^/  Status: /'
    echo ""
done
```

### Tenant Policy Validation
```bash
#!/bin/bash
# validate-tenant-policies.sh

TENANT="$1"
if [[ -z "$TENANT" ]]; then
    echo "Usage: $0 <tenant-name>"
    exit 1
fi

echo "Validating policies for tenant: $TENANT"

# Get tenant namespaces
NAMESPACES=$(kubectl get namespaces -l tenant=$TENANT --no-headers -o custom-columns=":metadata.name")

if [[ -z "$NAMESPACES" ]]; then
    echo "No namespaces found for tenant: $TENANT"
    exit 1
fi

echo "Tenant namespaces: $NAMESPACES"
echo ""

for ns in $NAMESPACES; do
    echo "=== Validating namespace: $ns ==="
    
    # Check resource quota
    echo "Resource Quota Check:"
    if kubectl get resourcequota -n $ns --no-headers > /dev/null 2>&1; then
        echo "  ✅ ResourceQuota present"
        kubectl get resourcequota -n $ns -o yaml | grep -A 10 "spec:" | sed 's/^/    /'
    else
        echo "  ❌ No ResourceQuota found"
    fi
    
    # Check limit range
    echo "LimitRange Check:"
    if kubectl get limitrange -n $ns --no-headers > /dev/null 2>&1; then
        echo "  ✅ LimitRange present"
    else
        echo "  ❌ No LimitRange found"
    fi
    
    # Check network policies
    echo "NetworkPolicy Check:"
    NETPOL_COUNT=$(kubectl get networkpolicy -n $ns --no-headers 2>/dev/null | wc -l)
    if [[ $NETPOL_COUNT -gt 0 ]]; then
        echo "  ✅ $NETPOL_COUNT NetworkPolicy(ies) present"
    else
        echo "  ❌ No NetworkPolicies found"
    fi
    
    # Check namespace policies
    echo "Kyverno Policies Check:"
    POLICY_COUNT=$(kubectl get policy -n $ns --no-headers 2>/dev/null | wc -l)
    if [[ $POLICY_COUNT -gt 0 ]]; then
        echo "  ✅ $POLICY_COUNT Kyverno Policy(ies) present"
        kubectl get policy -n $ns --no-headers -o custom-columns="NAME:.metadata.name,READY:.status.ready" | sed 's/^/    /'
    else
        echo "  ❌ No Kyverno Policies found"
    fi
    
    # Check pod compliance
    echo "Pod Compliance Check:"
    PODS=$(kubectl get pods -n $ns --no-headers -o custom-columns=":metadata.name" 2>/dev/null)
    if [[ -n "$PODS" ]]; then
        for pod in $PODS; do
            # Check if pod runs as non-root
            NON_ROOT=$(kubectl get pod $pod -n $ns -o jsonpath='{.spec.securityContext.runAsNonRoot}' 2>/dev/null)
            if [[ "$NON_ROOT" == "true" ]]; then
                echo "  ✅ Pod $pod runs as non-root"
            else
                echo "  ❌ Pod $pod may run as root"
            fi
        done
    else
        echo "  ℹ️  No pods to check"
    fi
    
    echo ""
done

echo "Validation complete for tenant: $TENANT"
```

## Integration Patterns

### Terraform Integration
```hcl
# terraform/multi-tenant-setup.tf
variable "tenants" {
  description = "List of tenant configurations"
  type = list(object({
    name         = string
    cpu_limit    = string
    memory_limit = string
    pod_limit    = string
  }))
  default = [
    {
      name         = "acme"
      cpu_limit    = "4"
      memory_limit = "8Gi"
      pod_limit    = "20"
    },
    {
      name         = "globex"
      cpu_limit    = "2"
      memory_limit = "4Gi"
      pod_limit    = "10"
    }
  ]
}

# Create tenant namespaces
resource "kubernetes_namespace" "tenant_namespaces" {
  for_each = { for tenant in var.tenants : tenant.name => tenant }
  
  metadata {
    name = "tenant-${each.value.name}"
    labels = {
      tenant             = each.value.name
      "kyverno.io/tenant" = each.value.name
    }
  }
}

# Create resource quotas for each tenant
resource "kubernetes_resource_quota" "tenant_quotas" {
  for_each = { for tenant in var.tenants : tenant.name => tenant }
  
  metadata {
    name      = "tenant-quota"
    namespace = kubernetes_namespace.tenant_namespaces[each.key].metadata[0].name
  }
  
  spec {
    hard = {
      "requests.cpu"                = each.value.cpu_limit
      "limits.cpu"                  = each.value.cpu_limit
      "requests.memory"             = each.value.memory_limit
      "limits.memory"               = each.value.memory_limit
      "pods"                        = each.value.pod_limit
      "services"                    = "10"
      "persistentvolumeclaims"     = "5"
    }
  }
}

# Apply Kyverno cluster policies
resource "kubectl_manifest" "tenant_isolation_policy" {
  yaml_body = yamlencode({
    apiVersion = "kyverno.io/v1"
    kind       = "ClusterPolicy"
    metadata = {
      name = "tenant-isolation"
    }
    spec = {
      validationFailureAction = "enforce"
      background              = true
      rules = [
        {
          name = "require-tenant-label"
          match = {
            any = [
              {
                resources = {
                  kinds = ["Namespace"]
                  names = ["tenant-*"]
                }
              }
            ]
          }
          validate = {
            message = "Tenant namespaces must have tenant label"
            pattern = {
              metadata = {
                labels = {
                  tenant = "?*"
                }
              }
            }
          }
        }
      ]
    }
  })
}
```

### GitOps Workflow
```yaml
# .github/workflows/tenant-management.yml
name: Multi-Tenant Management
on:
  push:
    paths:
      - 'tenants/**'
  workflow_dispatch:
    inputs:
      tenant_name:
        description: 'Tenant name to manage'
        required: true
      action:
        description: 'Action to perform'
        required: true
        type: choice
        options:
          - create
          - update
          - validate
          - delete

jobs:
  manage-tenant:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup kubectl
      uses: azure/setup-kubectl@v1
      
    - name: Setup Helm
      uses: azure/setup-helm@v1
      
    - name: Configure cluster access
      run: |
        echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
        
    - name: Install Kyverno CLI
      run: |
        curl -LO https://github.com/kyverno/kyverno/releases/latest/download/kyverno-cli_linux_x86_64.tar.gz
        tar -xzf kyverno-cli_linux_x86_64.tar.gz
        sudo mv kyverno /usr/local/bin/
        
    - name: Validate Tenant Policies
      run: |
        for policy in tenants/${{ github.event.inputs.tenant_name || 'all' }}/*.yaml; do
          echo "Validating policy: $policy"
          kyverno apply $policy --dry-run
        done
        
    - name: Apply Tenant Configuration
      if: github.event.inputs.action == 'create' || github.event.inputs.action == 'update'
      run: |
        TENANT="${{ github.event.inputs.tenant_name }}"
        
        # Create namespace
        kubectl create namespace "tenant-$TENANT" --dry-run=client -o yaml | kubectl apply -f -
        
        # Label namespace
        kubectl label namespace "tenant-$TENANT" tenant="$TENANT" --overwrite
        
        # Apply tenant policies
        kubectl apply -f tenants/$TENANT/
        
    - name: Validate Tenant Setup
      run: |
        TENANT="${{ github.event.inputs.tenant_name }}"
        
        # Check namespace exists
        kubectl get namespace "tenant-$TENANT"
        
        # Check resource quota
        kubectl get resourcequota -n "tenant-$TENANT"
        
        # Check policies
        kubectl get policy -n "tenant-$TENANT"
        
        # Check policy reports
        kubectl get policyreport -n "tenant-$TENANT" || echo "No policy violations"
```

## Best Practices

### Tenant Design
- **Namespace Strategy**: Use consistent naming patterns (e.g., tenant-{name})
- **Label Standards**: Apply consistent labeling for tenant identification
- **Resource Boundaries**: Set appropriate resource limits based on tenant requirements
- **Security Isolation**: Implement proper network and security policies

### Policy Management
- **Policy as Code**: Store all policies in version control
- **Graduated Enforcement**: Start with audit mode, then enforce
- **Regular Review**: Periodically review and update policies
- **Documentation**: Document policy intent and business requirements

### Operational Excellence
- **Monitoring**: Implement comprehensive tenant resource monitoring
- **Alerting**: Set up alerts for quota violations and policy failures
- **Automation**: Automate tenant onboarding and offboarding
- **Backup**: Regularly backup tenant configurations and data

### Security Considerations
- **Least Privilege**: Apply principle of least privilege to tenant access
- **Network Segmentation**: Implement proper network isolation
- **Secret Management**: Use proper secret management for tenant credentials
- **Audit Logging**: Enable comprehensive audit logging for tenant activities

## Error Handling

### Common Issues
```bash
# Policy validation failures
kyverno apply policy.yaml --dry-run
# Solution: Check policy syntax and resource constraints

# Resource quota exceeded
kubectl describe resourcequota -n tenant-namespace
# Solution: Adjust quotas or optimize resource usage

# Namespace creation conflicts
kubectl get namespace tenant-name
# Solution: Use unique namespace names or clean up existing resources

# Policy not applying
kubectl get clusterpolicy policy-name -o yaml
# Solution: Check policy selectors and conditions
```

### Troubleshooting
- **Policy Issues**: Use kyverno CLI to validate policies before applying
- **Resource Conflicts**: Check existing resources before creating new ones
- **Permission Errors**: Verify RBAC permissions for tenant operations
- **Quota Problems**: Monitor and adjust resource quotas based on usage

Kyverno Multi-Tenant provides comprehensive multi-tenancy capabilities for Kubernetes clusters through policy automation, resource management, and security enforcement, enabling safe and efficient cluster sharing across multiple tenants.