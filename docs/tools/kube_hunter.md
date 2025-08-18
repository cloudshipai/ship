# Kube-hunter

Kubernetes penetration testing tool for identifying security vulnerabilities.

## Description

Kube-hunter is a penetration testing tool designed to identify security weaknesses in Kubernetes clusters. It hunts for security issues by attacking Kubernetes clusters from the inside and outside perspectives, helping security teams understand their attack surface. The tool can perform passive reconnaissance to gather information about the cluster or active testing to actually exploit vulnerabilities. It's particularly valuable for security audits, compliance assessments, and proactive security testing.

## MCP Tools

### Remote Scanning
- **`kube_hunter_remote_scan`** - Scan specific IP addresses or DNS names for Kubernetes vulnerabilities
- **`kube_hunter_cidr_scan`** - Scan IP ranges using CIDR notation for comprehensive network analysis

### Network Discovery
- **`kube_hunter_interface_scan`** - Scan all local network interfaces for Kubernetes services
- **`kube_hunter_pod_scan`** - Scan from within a Kubernetes pod for internal vulnerabilities

### Testing Management
- **`kube_hunter_list_tests`** - List all available security tests and hunters
- **`kube_hunter_custom_hunters`** - Run specific security tests by hunter class names

## Real CLI Commands Used

### Basic Scanning
- `kube-hunter --remote <target>` - Scan specific IP or hostname
- `kube-hunter --cidr <cidr>` - Scan IP range (e.g., 192.168.0.0/24)
- `kube-hunter --interface` - Scan all local network interfaces
- `kube-hunter --pod` - Scan from within a Kubernetes pod

### Active Testing
- `kube-hunter --active` - Enable active hunting (exploit testing)
- `kube-hunter --remote <target> --active` - Active scan of remote target
- `kube-hunter --pod --active` - Active scan from pod perspective

### Output Options
- `kube-hunter --report json` - Output results in JSON format
- `kube-hunter --log DEBUG` - Set logging level (DEBUG, INFO, WARNING)
- `kube-hunter --dispatch stdout` - Output method (stdout, http)

### Advanced Options
- `kube-hunter --mapping` - Show only network mapping of Kubernetes nodes
- `kube-hunter --quick` - Limit subnet scanning to /24 CIDR
- `kube-hunter --k8s-auto-discover-nodes` - Auto-discover and scan all nodes
- `kube-hunter --kubeconfig <path>` - Use specific kubeconfig file
- `kube-hunter --service-account-token <token>` - Use service account JWT token

### Testing Control
- `kube-hunter --list` - List available tests
- `kube-hunter --list --active` - Include active hunting tests
- `kube-hunter --raw-hunter-names` - Show raw hunter class names
- `kube-hunter --custom <hunters>` - Run specific hunters

## Use Cases

### Security Assessment
- **Vulnerability Discovery**: Identify security weaknesses in Kubernetes clusters
- **Attack Surface Analysis**: Understand exposed services and potential entry points
- **Penetration Testing**: Simulate real-world attacks on Kubernetes infrastructure
- **Security Posture Evaluation**: Assess overall security stance of clusters

### Compliance and Auditing
- **Security Audits**: Perform comprehensive security assessments for compliance
- **Risk Assessment**: Identify and quantify security risks in Kubernetes deployments
- **Regulatory Compliance**: Support compliance with security frameworks and standards
- **Documentation**: Generate security reports for audit trails

### DevSecOps Integration
- **CI/CD Security Gates**: Integrate security testing into deployment pipelines
- **Continuous Security**: Regular automated security testing of clusters
- **Security Validation**: Verify security controls are working as expected
- **Shift-Left Security**: Early detection of security issues in development

### Red Team Exercises
- **Offensive Security**: Simulate advanced persistent threat scenarios
- **Attack Simulation**: Test incident response and detection capabilities
- **Security Training**: Train security teams on Kubernetes attack vectors
- **Purple Team Activities**: Collaborative testing between red and blue teams

## Configuration Examples

### Basic Vulnerability Scanning
```bash
# Scan specific target
kube-hunter --remote 192.168.1.100

# Scan IP range
kube-hunter --cidr 192.168.1.0/24

# Scan all local interfaces
kube-hunter --interface

# Scan from pod perspective
kube-hunter --pod
```

### Active Security Testing
```bash
# Active scan of remote target
kube-hunter --remote api.cluster.local --active

# Active scan with detailed logging
kube-hunter --remote 10.0.1.10 --active --log DEBUG

# Active pod-based scan with node discovery
kube-hunter --pod --active --k8s-auto-discover-nodes

# Quick active scan (limited to /24)
kube-hunter --cidr 192.168.0.0/16 --active --quick
```

### Reconnaissance and Mapping
```bash
# Network mapping only
kube-hunter --cidr 10.0.0.0/8 --mapping

# Quick network discovery
kube-hunter --interface --quick --mapping

# Passive reconnaissance
kube-hunter --remote cluster.example.com --log INFO

# Interface scan with JSON output
kube-hunter --interface --report json
```

### Custom Testing
```bash
# List available tests
kube-hunter --list

# List active hunting tests
kube-hunter --list --active

# Show raw hunter class names
kube-hunter --list --raw-hunter-names

# Run specific hunters
kube-hunter --custom "KubeletExposure ApiServerExposure"

# Custom hunters with active testing
kube-hunter --custom "EtcdExposure" --active --remote 192.168.1.50
```

## Advanced Usage

### Comprehensive Security Assessment
```bash
#!/bin/bash
# comprehensive-k8s-security-scan.sh

TARGET_NETWORK="$1"
if [[ -z "$TARGET_NETWORK" ]]; then
    echo "Usage: $0 <target-network-cidr>"
    exit 1
fi

DATE=$(date +%Y%m%d)
RESULTS_DIR="kube-hunter-results-$DATE"
mkdir -p $RESULTS_DIR

echo "Starting comprehensive Kubernetes security assessment..."

# Passive reconnaissance
echo "Phase 1: Passive reconnaissance..."
kube-hunter --cidr $TARGET_NETWORK --mapping \
    --report json > $RESULTS_DIR/network-mapping.json

# Basic vulnerability scan
echo "Phase 2: Basic vulnerability scanning..."
kube-hunter --cidr $TARGET_NETWORK \
    --report json > $RESULTS_DIR/passive-scan.json

# Active vulnerability testing (if authorized)
read -p "Perform active testing? This may impact services (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Phase 3: Active vulnerability testing..."
    kube-hunter --cidr $TARGET_NETWORK --active \
        --report json > $RESULTS_DIR/active-scan.json
    
    echo "Phase 4: Focused active testing..."
    kube-hunter --cidr $TARGET_NETWORK --active --quick \
        --report json > $RESULTS_DIR/focused-active.json
fi

# Generate summary report
echo "Generating summary report..."
echo "=== Kubernetes Security Assessment Report ===" > $RESULTS_DIR/summary.txt
echo "Date: $(date)" >> $RESULTS_DIR/summary.txt
echo "Target: $TARGET_NETWORK" >> $RESULTS_DIR/summary.txt
echo "" >> $RESULTS_DIR/summary.txt

# Parse vulnerabilities from JSON
if [[ -f $RESULTS_DIR/passive-scan.json ]]; then
    VULNERABILITIES=$(jq -r '.vulnerabilities | length' $RESULTS_DIR/passive-scan.json 2>/dev/null || echo "0")
    echo "Vulnerabilities found (passive): $VULNERABILITIES" >> $RESULTS_DIR/summary.txt
fi

if [[ -f $RESULTS_DIR/active-scan.json ]]; then
    ACTIVE_VULNS=$(jq -r '.vulnerabilities | length' $RESULTS_DIR/active-scan.json 2>/dev/null || echo "0")
    echo "Vulnerabilities found (active): $ACTIVE_VULNS" >> $RESULTS_DIR/summary.txt
fi

echo "Assessment complete! Results in $RESULTS_DIR/"
```

### Pod-Based Internal Assessment
```bash
#!/bin/bash
# internal-k8s-assessment.sh

echo "Starting internal Kubernetes security assessment..."

# Get service account information
echo "Service Account Information:"
cat /var/run/secrets/kubernetes.io/serviceaccount/namespace
echo ""

# Run internal pod scan
echo "Phase 1: Internal pod perspective scan..."
kube-hunter --pod --log INFO \
    --report json > /tmp/internal-scan.json

# Active internal testing (if authorized)
read -p "Perform active internal testing? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Phase 2: Active internal testing..."
    kube-hunter --pod --active \
        --k8s-auto-discover-nodes \
        --report json > /tmp/active-internal.json
fi

# List available tests for reference
echo "Available hunters:"
kube-hunter --list --raw-hunter-names

# Test specific high-value targets
echo "Phase 3: Targeted testing..."
kube-hunter --pod --custom "ApiServerExposure KubeletExposure EtcdExposure" \
    --report json > /tmp/targeted-scan.json

echo "Internal assessment complete!"
echo "Results stored in /tmp/*-scan.json"
```

### Multi-Cluster Security Testing
```bash
#!/bin/bash
# multi-cluster-security-test.sh

CLUSTERS=(
    "prod-cluster:10.1.0.0/16"
    "staging-cluster:10.2.0.0/16"
    "dev-cluster:10.3.0.0/16"
)

DATE=$(date +%Y%m%d)
REPORT_DIR="multi-cluster-security-$DATE"
mkdir -p $REPORT_DIR

for cluster_info in "${CLUSTERS[@]}"; do
    cluster_name=$(echo $cluster_info | cut -d: -f1)
    cluster_cidr=$(echo $cluster_info | cut -d: -f2)
    
    echo "Testing cluster: $cluster_name ($cluster_cidr)"
    
    # Create cluster-specific directory
    mkdir -p $REPORT_DIR/$cluster_name
    
    # Passive scan
    echo "  Running passive scan..."
    kube-hunter --cidr $cluster_cidr \
        --report json > $REPORT_DIR/$cluster_name/passive.json
    
    # Network mapping
    echo "  Generating network map..."
    kube-hunter --cidr $cluster_cidr --mapping \
        --report json > $REPORT_DIR/$cluster_name/mapping.json
    
    # Quick active scan (if approved)
    if [[ "$ENABLE_ACTIVE" == "true" ]]; then
        echo "  Running quick active scan..."
        kube-hunter --cidr $cluster_cidr --active --quick \
            --report json > $REPORT_DIR/$cluster_name/active.json
    fi
    
    # Generate cluster summary
    echo "=== $cluster_name Security Summary ===" > $REPORT_DIR/$cluster_name/summary.txt
    echo "CIDR: $cluster_cidr" >> $REPORT_DIR/$cluster_name/summary.txt
    echo "Scan Date: $(date)" >> $REPORT_DIR/$cluster_name/summary.txt
    
    if [[ -f $REPORT_DIR/$cluster_name/passive.json ]]; then
        VULN_COUNT=$(jq -r '.vulnerabilities | length' $REPORT_DIR/$cluster_name/passive.json 2>/dev/null || echo "0")
        echo "Vulnerabilities: $VULN_COUNT" >> $REPORT_DIR/$cluster_name/summary.txt
    fi
    
    echo "Cluster $cluster_name assessment complete"
done

echo "Multi-cluster security assessment finished!"
echo "Results in $REPORT_DIR/"
```

### Continuous Security Monitoring
```bash
#!/bin/bash
# continuous-k8s-security-monitor.sh

SLACK_WEBHOOK="$1"
CLUSTER_CIDR="$2"
THRESHOLD_VULNS=5

if [[ -z "$SLACK_WEBHOOK" || -z "$CLUSTER_CIDR" ]]; then
    echo "Usage: $0 <slack-webhook> <cluster-cidr>"
    exit 1
fi

echo "Starting continuous Kubernetes security monitoring..."

# Run security scan
SCAN_OUTPUT=$(kube-hunter --cidr $CLUSTER_CIDR --report json)
VULN_COUNT=$(echo "$SCAN_OUTPUT" | jq -r '.vulnerabilities | length' 2>/dev/null || echo "0")

echo "Vulnerabilities found: $VULN_COUNT"

# Save results with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
echo "$SCAN_OUTPUT" > "security-scan-$TIMESTAMP.json"

if [[ $VULN_COUNT -gt $THRESHOLD_VULNS ]]; then
    # High vulnerability count - send alert
    curl -X POST -H 'Content-type: application/json' \
        --data "{
            \"text\": \"🚨 Kubernetes Security Alert\",
            \"attachments\": [{
                \"color\": \"danger\",
                \"fields\": [{
                    \"title\": \"Vulnerabilities Found\",
                    \"value\": \"$VULN_COUNT\",
                    \"short\": true
                }, {
                    \"title\": \"Threshold\",
                    \"value\": \"$THRESHOLD_VULNS\",
                    \"short\": true
                }, {
                    \"title\": \"Cluster\",
                    \"value\": \"$CLUSTER_CIDR\",
                    \"short\": true
                }, {
                    \"title\": \"Scan Time\",
                    \"value\": \"$(date)\",
                    \"short\": true
                }],
                \"text\": \"Security scan detected $VULN_COUNT vulnerabilities in Kubernetes cluster. Immediate attention required.\"
            }]
        }" \
        $SLACK_WEBHOOK
    
    echo "Alert sent: High vulnerability count detected"
elif [[ $VULN_COUNT -gt 0 ]]; then
    # Some vulnerabilities - send info
    curl -X POST -H 'Content-type: application/json' \
        --data "{
            \"text\": \"ℹ️ Kubernetes Security Report\",
            \"attachments\": [{
                \"color\": \"warning\",
                \"text\": \"Security scan found $VULN_COUNT vulnerabilities in cluster $CLUSTER_CIDR. Review recommended.\"
            }]
        }" \
        $SLACK_WEBHOOK
    
    echo "Info sent: Vulnerabilities detected within threshold"
else
    echo "No vulnerabilities detected - cluster secure"
fi

echo "Security monitoring completed: security-scan-$TIMESTAMP.json"
```

## Integration Patterns

### CI/CD Pipeline Integration
```yaml
# .github/workflows/kubernetes-security-scan.yml
name: Kubernetes Security Scan
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly Monday 2 AM
  workflow_dispatch:
    inputs:
      target:
        description: 'Target CIDR or IP'
        required: true
        default: '10.0.0.0/8'
      active_testing:
        description: 'Enable active testing'
        type: boolean
        default: false

jobs:
  security-scan:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        cluster: [staging, production]
        
    steps:
    - uses: actions/checkout@v2
    
    - name: Install kube-hunter
      run: |
        pip3 install kube-hunter
        
    - name: Run Security Scan
      env:
        TARGET: ${{ inputs.target || '10.0.0.0/8' }}
        ACTIVE: ${{ inputs.active_testing || false }}
      run: |
        echo "Scanning cluster: ${{ matrix.cluster }}"
        
        # Basic passive scan
        kube-hunter --cidr $TARGET --report json > scan-results.json
        
        # Active scan if enabled
        if [[ "$ACTIVE" == "true" ]]; then
          kube-hunter --cidr $TARGET --active --quick --report json > active-scan-results.json
        fi
        
        # Parse results
        VULN_COUNT=$(jq -r '.vulnerabilities | length' scan-results.json)
        echo "Vulnerabilities found: $VULN_COUNT"
        
        # Fail if critical vulnerabilities found
        if [[ $VULN_COUNT -gt 10 ]]; then
          echo "❌ Critical vulnerability count exceeded"
          exit 1
        else
          echo "✅ Security scan passed"
        fi
    
    - name: Upload Results
      uses: actions/upload-artifact@v2
      with:
        name: security-scan-${{ matrix.cluster }}
        path: '*-results.json'
    
    - name: Security Summary
      run: |
        echo "### Kubernetes Security Scan Results" >> $GITHUB_STEP_SUMMARY
        echo "**Cluster:** ${{ matrix.cluster }}" >> $GITHUB_STEP_SUMMARY
        echo "**Vulnerabilities:** $(jq -r '.vulnerabilities | length' scan-results.json)" >> $GITHUB_STEP_SUMMARY
        echo "**Scan Date:** $(date)" >> $GITHUB_STEP_SUMMARY
```

### Kubernetes Job Deployment
```yaml
# kube-hunter-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: kube-hunter-security-scan
  namespace: security
spec:
  template:
    spec:
      serviceAccountName: kube-hunter
      containers:
      - name: kube-hunter
        image: aquasec/kube-hunter:latest
        command: ["kube-hunter"]
        args: ["--pod", "--active", "--report", "json"]
        env:
        - name: LOG_LEVEL
          value: "INFO"
        volumeMounts:
        - name: results
          mountPath: /results
      volumes:
      - name: results
        emptyDir: {}
      restartPolicy: Never
  backoffLimit: 2

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-hunter
  namespace: security

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kube-hunter
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services"]
  verbs: ["get", "list"]
- apiGroups: ["extensions", "apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kube-hunter
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-hunter
subjects:
- kind: ServiceAccount
  name: kube-hunter
  namespace: security
```

### Terraform Integration
```hcl
# terraform/security-scanning.tf
resource "kubernetes_namespace" "security" {
  metadata {
    name = "security"
  }
}

resource "kubernetes_cron_job" "kube_hunter" {
  metadata {
    name      = "kube-hunter-scan"
    namespace = kubernetes_namespace.security.metadata[0].name
  }

  spec {
    schedule = "0 2 * * 1"  # Weekly Monday 2 AM
    
    job_template {
      metadata {
        labels = {
          app = "kube-hunter"
        }
      }
      
      spec {
        template {
          metadata {
            labels = {
              app = "kube-hunter"
            }
          }
          
          spec {
            service_account_name = kubernetes_service_account.kube_hunter.metadata[0].name
            
            container {
              name  = "kube-hunter"
              image = "aquasec/kube-hunter:latest"
              
              command = ["kube-hunter"]
              args    = ["--pod", "--log", "INFO", "--report", "json"]
              
              env {
                name  = "CLUSTER_NAME"
                value = var.cluster_name
              }
            }
            
            restart_policy = "OnFailure"
          }
        }
      }
    }
  }
}

resource "kubernetes_service_account" "kube_hunter" {
  metadata {
    name      = "kube-hunter"
    namespace = kubernetes_namespace.security.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "kube_hunter" {
  metadata {
    name = "kube-hunter"
  }

  rule {
    api_groups = [""]
    resources  = ["nodes", "pods", "services", "endpoints"]
    verbs      = ["get", "list"]
  }

  rule {
    api_groups = ["apps", "extensions"]
    resources  = ["deployments", "replicasets", "daemonsets"]
    verbs      = ["get", "list"]
  }
}

resource "kubernetes_cluster_role_binding" "kube_hunter" {
  metadata {
    name = "kube-hunter"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.kube_hunter.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.kube_hunter.metadata[0].name
    namespace = kubernetes_namespace.security.metadata[0].name
  }
}
```

## Best Practices

### Security Testing Strategy
- **Progressive Testing**: Start with passive reconnaissance before active testing
- **Authorization**: Always get explicit authorization before running active tests
- **Scope Definition**: Clearly define test scope and boundaries
- **Impact Assessment**: Understand potential impact of active testing on services

### Operational Security
- **Scheduled Scanning**: Regular automated security assessments
- **Baseline Establishment**: Create security baselines for comparison
- **Incident Response**: Have procedures for addressing discovered vulnerabilities
- **Documentation**: Maintain detailed records of testing activities

### Integration Approach
- **CI/CD Integration**: Include security testing in deployment pipelines
- **Monitoring Integration**: Send results to security monitoring systems
- **Alerting**: Configure appropriate alerting for critical vulnerabilities
- **Reporting**: Generate executive summaries for stakeholders

### Risk Management
- **Vulnerability Prioritization**: Focus on critical and high-risk vulnerabilities first
- **Regular Updates**: Keep kube-hunter updated for latest vulnerability checks
- **False Positive Management**: Validate findings to reduce false positives
- **Remediation Tracking**: Track remediation efforts and verify fixes

## Error Handling

### Common Issues
```bash
# Connection refused to target
kube-hunter --remote target.example.com --log DEBUG
# Check network connectivity and firewall rules

# No vulnerabilities found (suspicious)
kube-hunter --list
# Verify tool is working and up to date

# Permission denied in pod scan
kubectl auth can-i get nodes --as=system:serviceaccount:default:default
# Check service account permissions

# Timeout during scanning
kube-hunter --cidr 10.0.0.0/8 --quick
# Use quick scan for large networks
```

### Troubleshooting
- **Network Connectivity**: Verify network access to target systems
- **Permissions**: Ensure appropriate RBAC permissions for pod-based scans
- **Tool Updates**: Keep kube-hunter updated for latest vulnerability signatures
- **Resource Limits**: Monitor resource usage during large network scans

Kube-hunter provides comprehensive Kubernetes penetration testing capabilities, enabling security teams to proactively identify and address vulnerabilities in their Kubernetes infrastructure through realistic attack simulations.