# Kubernetes Network Policy Management

Manage and analyze Kubernetes network policies using kubectl and specialized tools.

## Description

Kubernetes Network Policy Management provides comprehensive tools for managing, analyzing, and visualizing network policies in Kubernetes clusters. It combines standard kubectl operations for basic policy management with specialized tools like netfetch for security scanning and netpol-analyzer for connectivity analysis. These tools help teams implement proper network segmentation, ensure policy compliance, and understand network connectivity patterns in their Kubernetes environments.

## MCP Tools

### Basic Policy Management
- **`k8s_network_policy_kubectl`** - Manage network policies using kubectl (CRUD operations)

### Policy Analysis & Scanning
- **`k8s_network_policy_netfetch_scan`** - Scan network policies for security gaps using netfetch
- **`k8s_network_policy_netfetch_dash`** - Launch netfetch dashboard for visualization

### Connectivity Analysis
- **`k8s_network_policy_netpol_eval`** - Evaluate connectivity between pods using netpol-analyzer
- **`k8s_network_policy_netpol_list`** - List all allowed connections using netpol-analyzer
- **`k8s_network_policy_netpol_diff`** - Compare network policies between environments

## Real CLI Commands Used

### kubectl Commands
- `kubectl get networkpolicy` - List network policies
- `kubectl get networkpolicy <name>` - Get specific policy
- `kubectl describe networkpolicy <name>` - Describe policy details
- `kubectl create -f <policy.yaml>` - Create policy from file
- `kubectl apply -f <policy.yaml>` - Apply policy changes
- `kubectl delete networkpolicy <name>` - Delete policy
- `kubectl get networkpolicy -n <namespace>` - Namespace-specific operations
- `kubectl get networkpolicy -o json/yaml` - Output in different formats

### Netfetch Commands
- `netfetch scan` - Scan entire cluster for policy gaps
- `netfetch scan <namespace>` - Scan specific namespace
- `netfetch scan --dryrun` - Run without making changes
- `netfetch scan --cilium` - Scan Cilium network policies
- `netfetch scan --target <policy>` - Scan specific policy
- `netfetch dash` - Launch dashboard
- `netfetch dash --port 8081` - Dashboard on custom port

### Netpol-analyzer Commands
- `netpol-analyzer eval --dirpath <dir> -s <source> -d <dest>` - Test connectivity
- `netpol-analyzer list --dirpath <dir>` - List all connections
- `netpol-analyzer diff --dir1 <dir1> --dir2 <dir2>` - Compare policies
- `netpol-analyzer eval --dirpath <dir> -s <source> -d <dest> -p <port>` - Test specific port
- `netpol-analyzer list --dirpath <dir> -v` - Verbose output

## Use Cases

### Network Policy Management
- **Policy Creation**: Create and apply network policies for workload isolation
- **Policy Updates**: Modify existing policies to adapt to changing requirements
- **Policy Cleanup**: Remove obsolete or unused network policies
- **Policy Validation**: Ensure policies are correctly formatted and deployable

### Security Assessment
- **Gap Analysis**: Identify workloads without network policy protection
- **Compliance Scanning**: Verify adherence to security policies
- **Risk Assessment**: Understand exposure of unprotected services
- **Security Scoring**: Get quantified security metrics for clusters

### Connectivity Analysis
- **Traffic Flow Analysis**: Understand allowed traffic patterns
- **Connectivity Testing**: Verify network policies work as expected
- **Impact Assessment**: Analyze effects of policy changes
- **Troubleshooting**: Debug connectivity issues in restricted environments

### Environment Comparison
- **Policy Drift Detection**: Compare policies between environments
- **Configuration Consistency**: Ensure consistent security across clusters
- **Change Impact Analysis**: Understand differences between policy versions
- **Compliance Auditing**: Verify policy consistency across environments

## Configuration Examples

### Basic kubectl Network Policy Management
```bash
# List all network policies
kubectl get networkpolicy --all-namespaces

# Get network policies in specific namespace
kubectl get networkpolicy -n production

# Describe a specific policy
kubectl describe networkpolicy web-policy -n production

# Create policy from file
kubectl create -f network-policy.yaml

# Apply policy changes
kubectl apply -f updated-policy.yaml

# Delete policy
kubectl delete networkpolicy old-policy -n staging

# Get policy in YAML format
kubectl get networkpolicy web-policy -n production -o yaml
```

### Network Policy Security Scanning
```bash
# Scan entire cluster for security gaps
netfetch scan

# Scan specific namespace
netfetch scan production

# Dry run scan (no changes)
netfetch scan --dryrun

# Scan Cilium network policies
netfetch scan --cilium

# Scan specific policy
netfetch scan --target production-deny-all

# Scan with custom kubeconfig
netfetch scan --kubeconfig ~/.kube/prod-config

# Scan namespace with Cilium support
netfetch scan production --cilium
```

### Network Policy Dashboard
```bash
# Launch default dashboard
netfetch dash

# Launch dashboard on custom port
netfetch dash --port 8081

# Access dashboard in browser
# Navigate to http://localhost:8080 (or custom port)
```

### Connectivity Analysis
```bash
# Test connectivity between pods
netpol-analyzer eval --dirpath ./k8s-manifests -s web-pod -d db-pod

# Test specific port connectivity
netpol-analyzer eval --dirpath ./k8s-manifests -s frontend -d backend -p 8080

# Verbose connectivity evaluation
netpol-analyzer eval --dirpath ./k8s-manifests -s app -d service -v

# List all allowed connections
netpol-analyzer list --dirpath ./k8s-manifests

# Verbose connection listing
netpol-analyzer list --dirpath ./k8s-manifests -v

# Quiet mode (minimal output)
netpol-analyzer list --dirpath ./k8s-manifests -q
```

### Policy Comparison
```bash
# Compare policies between environments
netpol-analyzer diff --dir1 ./staging-manifests --dir2 ./prod-manifests

# Verbose policy comparison
netpol-analyzer diff --dir1 ./old-policies --dir2 ./new-policies -v

# Compare current vs proposed changes
netpol-analyzer diff --dir1 ./current --dir2 ./proposed
```

## Advanced Usage

### Comprehensive Security Assessment
```bash
#!/bin/bash
# comprehensive-netpol-audit.sh

NAMESPACES=("production" "staging" "development")
DATE=$(date +%Y%m%d)
REPORT_DIR="network-policy-audit-$DATE"

mkdir -p $REPORT_DIR

echo "Starting comprehensive network policy audit..."

# Scan all namespaces
for ns in "${NAMESPACES[@]}"; do
    echo "Scanning namespace: $ns"
    
    # Netfetch scan
    netfetch scan $ns > $REPORT_DIR/$ns-netfetch-scan.txt
    
    # Export network policies
    kubectl get networkpolicy -n $ns -o yaml > $REPORT_DIR/$ns-policies.yaml
    
    # List policy details
    kubectl get networkpolicy -n $ns -o wide > $REPORT_DIR/$ns-policy-list.txt
done

# Generate cluster-wide scan
echo "Performing cluster-wide scan..."
netfetch scan > $REPORT_DIR/cluster-wide-scan.txt

# Launch dashboard for interactive analysis
echo "Launching dashboard for interactive analysis..."
netfetch dash --port 8080 &

echo "Audit complete! Results in $REPORT_DIR/"
echo "Dashboard available at http://localhost:8080"
```

### Policy Validation Pipeline
```bash
#!/bin/bash
# validate-network-policies.sh

POLICY_DIR="$1"
if [[ -z "$POLICY_DIR" ]]; then
    echo "Usage: $0 <policy-directory>"
    exit 1
fi

echo "Validating network policies in: $POLICY_DIR"

# Validate YAML syntax
echo "Checking YAML syntax..."
for file in $POLICY_DIR/*.yaml; do
    if ! kubectl apply --dry-run=client -f "$file" > /dev/null 2>&1; then
        echo "ERROR: Invalid YAML syntax in $file"
        exit 1
    fi
done

# Analyze connectivity impact
echo "Analyzing connectivity impact..."
netpol-analyzer list --dirpath $POLICY_DIR > connectivity-analysis.txt

# Check for overly permissive policies
echo "Checking for security gaps..."
netfetch scan --dryrun > security-gaps.txt

echo "Validation complete!"
echo "- Connectivity analysis: connectivity-analysis.txt"
echo "- Security gaps: security-gaps.txt"
```

### Multi-Environment Policy Management
```bash
#!/bin/bash
# manage-multi-env-policies.sh

ENVIRONMENTS=("dev" "staging" "prod")
BASE_DIR="/policies"

for env in "${ENVIRONMENTS[@]}"; do
    echo "Processing $env environment..."
    
    # Apply environment-specific policies
    kubectl apply -f $BASE_DIR/$env/ --namespace=$env
    
    # Verify policies are applied
    kubectl get networkpolicy -n $env
    
    # Scan for security gaps
    netfetch scan $env > $env-scan-results.txt
    
    # Test critical connections
    netpol-analyzer eval --dirpath $BASE_DIR/$env -s web -d api -p 8080
    netpol-analyzer eval --dirpath $BASE_DIR/$env -s api -d database -p 5432
    
    echo "$env environment processed."
done

# Generate comparison report
echo "Generating environment comparison..."
netpol-analyzer diff --dir1 $BASE_DIR/staging --dir2 $BASE_DIR/prod > staging-prod-diff.txt

echo "Multi-environment policy management complete!"
```

### Continuous Monitoring Script
```bash
#!/bin/bash
# monitor-network-policies.sh

SLACK_WEBHOOK="$1"
THRESHOLD_SCORE=80

if [[ -z "$SLACK_WEBHOOK" ]]; then
    echo "Usage: $0 <slack-webhook-url>"
    exit 1
fi

# Run security scan
SCAN_OUTPUT=$(netfetch scan)
SCORE=$(echo "$SCAN_OUTPUT" | grep -o 'Score: [0-9]*' | cut -d' ' -f2)

echo "Current security score: $SCORE"

if [[ $SCORE -lt $THRESHOLD_SCORE ]]; then
    # Send alert to Slack
    curl -X POST -H 'Content-type: application/json' \
        --data "{
            \"text\": \"🚨 Network Policy Alert\",
            \"attachments\": [{
                \"color\": \"danger\",
                \"fields\": [{
                    \"title\": \"Security Score\",
                    \"value\": \"$SCORE/$THRESHOLD_SCORE\",
                    \"short\": true
                }, {
                    \"title\": \"Status\",
                    \"value\": \"Below Threshold\",
                    \"short\": true
                }],
                \"text\": \"Network policy security score has dropped below threshold. Please review unprotected workloads.\"
            }]
        }" \
        $SLACK_WEBHOOK
    
    echo "Alert sent: Security score below threshold"
else
    echo "Security score acceptable: $SCORE >= $THRESHOLD_SCORE"
fi

# Generate detailed report
netfetch scan > daily-security-scan.txt
echo "Daily scan completed: daily-security-scan.txt"
```

## Integration Patterns

### GitOps Workflow
```yaml
# .github/workflows/network-policy-validation.yml
name: Network Policy Validation
on:
  pull_request:
    paths:
      - 'k8s/network-policies/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup kubectl
      uses: azure/setup-kubectl@v1
      
    - name: Install netpol-analyzer
      run: |
        wget -O netpol-analyzer https://github.com/np-guard/netpol-analyzer/releases/latest/download/netpol-analyzer-linux-amd64
        chmod +x netpol-analyzer
        sudo mv netpol-analyzer /usr/local/bin/
    
    - name: Validate Policy Syntax
      run: |
        for file in k8s/network-policies/*.yaml; do
          kubectl apply --dry-run=client -f "$file"
        done
    
    - name: Analyze Connectivity Impact
      run: |
        netpol-analyzer list --dirpath k8s/network-policies/ > connectivity-report.txt
        cat connectivity-report.txt
    
    - name: Compare with Main Branch
      run: |
        git checkout main
        netpol-analyzer list --dirpath k8s/network-policies/ > main-connectivity.txt
        git checkout -
        netpol-analyzer diff --dir1 . --dir2 . > policy-diff.txt || true
        cat policy-diff.txt
```

### Terraform Integration
```hcl
# terraform/network-policies.tf
resource "kubernetes_network_policy" "webapp_policy" {
  metadata {
    name      = "webapp-network-policy"
    namespace = var.namespace
  }

  spec {
    pod_selector {
      match_labels = {
        app = "webapp"
      }
    }

    policy_types = ["Ingress", "Egress"]

    ingress {
      from {
        pod_selector {
          match_labels = {
            app = "frontend"
          }
        }
      }
      ports {
        port     = "8080"
        protocol = "TCP"
      }
    }

    egress {
      to {
        pod_selector {
          match_labels = {
            app = "database"
          }
        }
      }
      ports {
        port     = "5432"
        protocol = "TCP"
      }
    }
  }
}

# Null resource to run validation after apply
resource "null_resource" "validate_policies" {
  depends_on = [kubernetes_network_policy.webapp_policy]
  
  provisioner "local-exec" {
    command = <<-EOT
      kubectl get networkpolicy -n ${var.namespace}
      netfetch scan ${var.namespace}
    EOT
  }
}
```

### Helm Chart Integration
```yaml
# helm/templates/network-policy.yaml
{{- if .Values.networkPolicy.enabled }}
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ include "app.fullname" . }}-netpol
  namespace: {{ .Release.Namespace }}
spec:
  podSelector:
    matchLabels:
      {{- include "app.selectorLabels" . | nindent 6 }}
  policyTypes:
  {{- if .Values.networkPolicy.ingress }}
  - Ingress
  {{- end }}
  {{- if .Values.networkPolicy.egress }}
  - Egress
  {{- end }}
  {{- with .Values.networkPolicy.ingress }}
  ingress:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.networkPolicy.egress }}
  egress:
    {{- toYaml . | nindent 4 }}
  {{- end }}
{{- end }}

# helm/templates/tests/network-policy-test.yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "app.fullname" . }}-netpol-test"
  annotations:
    "helm.sh/hook": test
spec:
  restartPolicy: Never
  containers:
  - name: test
    image: alpine:latest
    command:
    - /bin/sh
    - -c
    - |
      # Test that network policy is applied
      if kubectl get networkpolicy {{ include "app.fullname" . }}-netpol -n {{ .Release.Namespace }}; then
        echo "Network policy found"
        exit 0
      else
        echo "Network policy not found"
        exit 1
      fi
```

## Best Practices

### Policy Design
- **Least Privilege**: Start with deny-all policies and add specific allow rules
- **Namespace Isolation**: Use namespace selectors for environment separation
- **Application Grouping**: Group related services with consistent labeling
- **Gradual Rollout**: Implement policies incrementally to avoid service disruption

### Security Implementation
- **Default Deny**: Implement default deny policies for all namespaces
- **Ingress Control**: Carefully control ingress traffic to critical services
- **Egress Restrictions**: Limit egress traffic to prevent data exfiltration
- **Regular Auditing**: Continuously monitor and audit network policies

### Operational Management
- **Documentation**: Document policy intent and business requirements
- **Testing**: Test policies in non-production environments first
- **Monitoring**: Monitor policy effectiveness and adjust as needed
- **Automation**: Automate policy deployment and validation

### Tool Usage
- **Regular Scanning**: Use netfetch for regular security assessments
- **Connectivity Testing**: Use netpol-analyzer for change impact analysis
- **Dashboard Monitoring**: Use netfetch dashboard for ongoing visibility
- **Integration**: Integrate tools into CI/CD pipelines for continuous validation

## Error Handling

### Common Issues
```bash
# Policy not taking effect
kubectl describe networkpolicy <policy-name> -n <namespace>
# Check CNI plugin supports network policies

# Connectivity blocked unexpectedly
netpol-analyzer eval --dirpath . -s source-pod -d dest-pod -v
# Use verbose mode to understand blocking rules

# Netfetch scan failures
netfetch scan --kubeconfig ~/.kube/config
# Verify kubeconfig and cluster access

# Missing tools
which netfetch netpol-analyzer
# Install missing tools
```

### Troubleshooting
- **CNI Support**: Verify CNI plugin supports network policies
- **Label Matching**: Ensure pod labels match policy selectors
- **Namespace Scope**: Check policy and pod namespaces align
- **Tool Permissions**: Verify cluster access permissions for tools

Kubernetes Network Policy Management provides comprehensive tools for implementing, analyzing, and maintaining network security in Kubernetes environments through proven CLI tools and established best practices.