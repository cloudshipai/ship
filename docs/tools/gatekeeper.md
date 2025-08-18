# Gatekeeper

OPA Gatekeeper for Kubernetes policy enforcement and governance.

## Description

Gatekeeper is a validating admission webhook that enforces Custom Resource Definitions (CRDs) based policies executed by Open Policy Agent (OPA). It provides policy-based control for Kubernetes clusters, allowing administrators to define and enforce security, compliance, and governance policies.

## MCP Tools

### Installation & Management
- **`gatekeeper_install`** - Install Gatekeeper using kubectl or Helm
- **`gatekeeper_uninstall`** - Uninstall Gatekeeper from cluster
- **`gatekeeper_get_status`** - Get Gatekeeper system status

### Policy Management
- **`gatekeeper_apply_constraint_template`** - Apply constraint template using kubectl
- **`gatekeeper_apply_constraint`** - Apply constraint using kubectl
- **`gatekeeper_get_constraint_templates`** - List constraint templates
- **`gatekeeper_get_constraints`** - List constraints by type

## Real CLI Commands Used

### Installation Commands
- `kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/v3.20.0/deploy/gatekeeper.yaml` - Install via kubectl
- `helm install gatekeeper gatekeeper/gatekeeper --namespace gatekeeper-system --create-namespace` - Install via Helm
- `kubectl delete -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/v3.20.0/deploy/gatekeeper.yaml` - Uninstall via kubectl
- `helm delete gatekeeper --namespace gatekeeper-system` - Uninstall via Helm

### Policy Management Commands
- `kubectl apply -f constraint-template.yaml` - Apply constraint template
- `kubectl apply -f constraint.yaml` - Apply constraint
- `kubectl get constrainttemplates` - List constraint templates
- `kubectl get <constraint-type>` - List specific constraint type
- `kubectl get pods -n gatekeeper-system` - Check Gatekeeper status

## Core Concepts

### Constraint Templates
Define the schema and logic for policy constraints:
```yaml
apiVersion: templates.gatekeeper.sh/v1beta1
kind: ConstraintTemplate
metadata:
  name: k8srequiredlabels
spec:
  crd:
    spec:
      names:
        kind: K8sRequiredLabels
      validation:
        properties:
          labels:
            type: array
            items:
              type: string
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package k8srequiredlabels
        
        violation[{"msg": msg}] {
          required := input.parameters.labels
          provided := input.review.object.metadata.labels
          missing := required[_]
          not provided[missing]
          msg := sprintf("Missing required label: %v", [missing])
        }
```

### Constraints
Instances of constraint templates applied to specific resources:
```yaml
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sRequiredLabels
metadata:
  name: must-have-env
spec:
  match:
    kinds:
      - apiGroups: ["apps"]
        kinds: ["Deployment"]
  parameters:
    labels: ["environment", "team"]
```

## Use Cases

### Security Governance
- Enforce security policies across clusters
- Prevent privileged container deployment
- Ensure proper RBAC configurations
- Block insecure image sources

### Compliance Management
- Implement regulatory compliance rules
- Enforce organizational standards
- Audit policy violations
- Generate compliance reports

### Resource Management
- Enforce resource quotas and limits
- Ensure proper labeling and tagging
- Control namespace usage
- Manage service mesh policies

### Development Standards
- Enforce coding and deployment standards
- Ensure proper configuration management
- Control container image policies
- Manage secret handling practices

## Policy Examples

### Required Labels
Ensure all deployments have required labels:
```yaml
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sRequiredLabels
metadata:
  name: deployment-must-have-labels
spec:
  match:
    kinds:
      - apiGroups: ["apps"]
        kinds: ["Deployment"]
  parameters:
    labels: ["app", "version", "environment"]
```

### Container Security
Prevent privileged containers:
```yaml
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sPSPPrivileged
metadata:
  name: psp-privileged-container
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
  parameters:
    runAsUser:
      rule: "MustRunAsNonRoot"
```

### Resource Limits
Enforce resource limits on containers:
```yaml
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: K8sRequiredResources
metadata:
  name: container-must-have-limits
spec:
  match:
    kinds:
      - apiGroups: [""]
        kinds: ["Pod"]
  parameters:
    limits: ["memory", "cpu"]
    requests: ["memory", "cpu"]
```

## Integration Features

### Admission Control
- Validates resources at creation/update time
- Blocks non-compliant resources
- Provides detailed violation messages
- Supports dry-run mode for testing

### Audit and Monitoring
- Continuous compliance monitoring
- Violation tracking and reporting
- Integration with monitoring systems
- Policy effectiveness metrics

### Multi-Cluster Management
- Consistent policy enforcement
- Centralized policy management
- Cross-cluster compliance reporting
- Federated governance

## Installation Requirements

- Kubernetes v1.16+
- Cluster admin permissions
- OPA runtime (included with Gatekeeper)
- Admission webhook capabilities

Works with any Kubernetes distribution and integrates with CI/CD pipelines, monitoring systems, and security platforms for comprehensive policy governance.