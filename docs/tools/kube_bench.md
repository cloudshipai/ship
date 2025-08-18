# Kube-bench

Kubernetes CIS benchmark tool for security assessment.

## Description

Kube-bench is a Go application that checks whether Kubernetes is deployed securely by running the checks documented in the CIS Kubernetes Benchmark. It provides comprehensive security assessments for Kubernetes clusters, helping organizations ensure their deployments follow security best practices. The tool automatically detects the Kubernetes version and maps it to the corresponding CIS Benchmark, providing detailed pass/fail results with remediation guidance.

## MCP Tools

### Core Benchmark Operations
- **`kube_bench_run`** - Run complete CIS Kubernetes benchmark with target selection
- **`kube_bench_run_checks`** - Run specific CIS benchmark checks by ID
- **`kube_bench_run_skip`** - Run benchmark while skipping specific checks or groups

### Advanced Configuration
- **`kube_bench_run_custom_output`** - Run benchmark with custom output options and filtering
- **`kube_bench_run_asff`** - Run benchmark and send results to AWS Security Hub

### Utility Operations
- **`kube_bench_version`** - Get kube-bench version information

## Real CLI Commands Used

### Core Commands
- `kube-bench` - Run default benchmark (auto-detects version and targets)
- `kube-bench run` - Explicitly run benchmark
- `kube-bench run --targets <targets>` - Run specific target components
- `kube-bench version` - Show version information

### Target Selection
- `kube-bench run --targets master` - Check master node components
- `kube-bench run --targets node` - Check worker node components
- `kube-bench run --targets etcd` - Check etcd configuration
- `kube-bench run --targets controlplane` - Check control plane
- `kube-bench run --targets policies` - Check Kubernetes policies
- `kube-bench run --targets master,node,etcd` - Multiple targets

### Check Filtering
- `kube-bench --check 1.1.1,1.1.2` - Run specific checks
- `kube-bench --skip 1.1,1.2.1` - Skip specific checks or groups
- `kube-bench --scored` - Run only scored checks
- `kube-bench --unscored` - Run only unscored checks

### Platform-Specific Benchmarks
- `kube-bench --benchmark gke-1.0` - Google Kubernetes Engine
- `kube-bench --benchmark ack-1.0` - Alibaba Cloud Kubernetes
- `kube-bench --benchmark rke-cis-1.7` - Rancher RKE
- `kube-bench --benchmark rh-1.0` - Red Hat OpenShift
- `kube-bench --benchmark tkgi-1.2.53` - VMware TKGI

### Output Options
- `kube-bench --json` - JSON output format
- `kube-bench --junit` - JUnit XML format
- `kube-bench --outputfile results.json` - Save to file
- `kube-bench --noremediations` - Hide remediation advice
- `kube-bench --noresults` - Hide detailed results
- `kube-bench --nosummary` - Hide summary section

### Configuration
- `kube-bench --config-dir /etc/kube-bench/cfg` - Custom config directory
- `kube-bench --config /path/to/config.yaml` - Custom config file
- `kube-bench --version 1.18` - Specify Kubernetes version
- `kube-bench --exit-code 42` - Custom exit code for failures

### AWS Integration
- `kube-bench --asff` - Send results to AWS Security Hub

## Use Cases

### Security Assessment
- **Compliance Auditing**: Verify Kubernetes deployment against CIS benchmarks
- **Security Baseline**: Establish security baseline for new clusters
- **Regular Scanning**: Continuous security monitoring of running clusters
- **Penetration Testing**: Identify security weaknesses before attackers do

### DevOps Integration
- **CI/CD Security Gates**: Integrate security checks into deployment pipelines
- **Infrastructure Validation**: Validate security during infrastructure provisioning
- **Configuration Drift**: Detect security configuration changes over time
- **Automated Remediation**: Generate actionable remediation guidance

### Cloud Platform Security
- **Multi-Cloud Assessment**: Assess security across different cloud providers
- **Managed Service Validation**: Verify security of managed Kubernetes services
- **Platform-Specific Checks**: Run benchmarks tailored to specific platforms
- **Hybrid Deployment**: Assess security across on-premises and cloud deployments

### Compliance Requirements
- **Regulatory Compliance**: Meet industry security standards and regulations
- **Security Documentation**: Generate reports for audit and compliance purposes
- **Risk Assessment**: Quantify security posture and identify high-risk areas
- **Governance**: Enforce security policies across multiple clusters

## Configuration Examples

### Basic Security Assessment
```bash
# Run default benchmark (auto-detects version and components)
kube-bench

# Run with explicit targets
kube-bench run --targets master,node,etcd

# Run for specific Kubernetes version
kube-bench --version 1.24

# Run with JSON output
kube-bench --json --outputfile security-report.json

# Run with custom config
kube-bench --config-dir /opt/kube-bench/cfg
```

### Platform-Specific Assessments
```bash
# Google Kubernetes Engine
kube-bench --benchmark gke-1.6.0

# Amazon EKS (using default detection)
kube-bench run --targets node

# Azure AKS (using default detection)
kube-bench run --targets node

# Red Hat OpenShift
kube-bench --benchmark rh-1.0

# Rancher RKE
kube-bench --benchmark rke-cis-1.7

# VMware TKGI
kube-bench --benchmark tkgi-1.2.53
```

### Targeted Security Checks
```bash
# Run specific checks
kube-bench --check 1.1.1,1.1.2,1.2.1

# Skip problematic checks
kube-bench --skip 1.1.12,1.2.6

# Run only scored checks
kube-bench --scored

# Run master node checks only
kube-bench run --targets master

# Run etcd checks only
kube-bench run --targets etcd
```

### Custom Output and Reporting
```bash
# Minimal output (just summary)
kube-bench --noresults --noremediations

# JSON output for automation
kube-bench --json --outputfile results.json

# JUnit XML for CI/CD integration
kube-bench --junit --outputfile results.xml

# Custom exit code for scripting
kube-bench --exit-code 42

# Hide specific sections
kube-bench --nosummary --nototals
```

## Advanced Usage

### Automated Security Scanning
```bash
#!/bin/bash
# automated-kube-bench-scan.sh

DATE=$(date +%Y%m%d)
RESULTS_DIR="kube-bench-results-$DATE"
mkdir -p $RESULTS_DIR

echo "Starting Kubernetes security assessment..."

# Run comprehensive benchmark
kube-bench run --targets master,node,etcd \
    --json \
    --outputfile $RESULTS_DIR/full-benchmark.json

# Run scored checks only for quick assessment
kube-bench --scored \
    --json \
    --outputfile $RESULTS_DIR/scored-only.json

# Generate human-readable report
kube-bench run --targets master,node,etcd \
    --outputfile $RESULTS_DIR/detailed-report.txt

# Check specific high-priority controls
kube-bench --check 1.1.1,1.1.2,1.1.3,1.2.1,1.2.2 \
    --outputfile $RESULTS_DIR/critical-checks.txt

echo "Security assessment complete. Results in $RESULTS_DIR/"

# Parse results for CI/CD decision making
if grep -q '"total_fail": 0' $RESULTS_DIR/scored-only.json; then
    echo "✅ Security assessment PASSED"
    exit 0
else
    echo "❌ Security assessment FAILED"
    echo "Review detailed report: $RESULTS_DIR/detailed-report.txt"
    exit 1
fi
```

### Multi-Cluster Security Assessment
```bash
#!/bin/bash
# multi-cluster-assessment.sh

CLUSTERS=("production" "staging" "development")
CONTEXTS=("prod-cluster" "staging-cluster" "dev-cluster")

for i in "${!CLUSTERS[@]}"; do
    cluster=${CLUSTERS[$i]}
    context=${CONTEXTS[$i]}
    
    echo "Assessing security for $cluster cluster..."
    
    # Switch to cluster context
    kubectl config use-context $context
    
    # Create results directory
    mkdir -p results/$cluster
    
    # Run assessment (note: kube-bench runs on nodes, not via kubectl)
    # This would typically be run as a DaemonSet or Job in the cluster
    
    # For demonstration, assuming kube-bench is run locally with access to cluster
    kube-bench run --targets node \
        --json \
        --outputfile results/$cluster/node-assessment.json
    
    kube-bench run --targets master \
        --json \
        --outputfile results/$cluster/master-assessment.json 2>/dev/null || \
        echo "Master assessment skipped (managed cluster)"
    
    # Generate summary report
    echo "=== $cluster Cluster Security Summary ===" > results/$cluster/summary.txt
    kube-bench run --targets node --noresults --nosummary >> results/$cluster/summary.txt
    
    echo "Assessment complete for $cluster cluster"
done

echo "Multi-cluster security assessment completed!"
echo "Results available in results/ directory"
```

### CI/CD Pipeline Integration
```yaml
# .github/workflows/security-assessment.yml
name: Kubernetes Security Assessment
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly on Monday at 2 AM
  workflow_dispatch:

jobs:
  security-scan:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        cluster: [staging, production]
        
    steps:
    - uses: actions/checkout@v2
    
    - name: Configure kubectl
      uses: azure/k8s-set-context@v1
      with:
        method: kubeconfig
        kubeconfig: ${{ secrets[format('KUBECONFIG_{0}', matrix.cluster)] }}
    
    - name: Deploy kube-bench Job
      run: |
        # Deploy kube-bench as a Kubernetes Job
        kubectl apply -f - <<EOF
        apiVersion: batch/v1
        kind: Job
        metadata:
          name: kube-bench-${{ matrix.cluster }}
        spec:
          template:
            spec:
              hostPID: true
              nodeSelector:
                kubernetes.io/os: linux
              tolerations:
              - operator: Exists
              containers:
              - name: kube-bench
                image: aquasec/kube-bench:latest
                command: ["kube-bench", "--json"]
                volumeMounts:
                - name: var-lib-etcd
                  mountPath: /var/lib/etcd
                  readOnly: true
                - name: var-lib-kubelet
                  mountPath: /var/lib/kubelet
                  readOnly: true
                - name: etc-kubernetes
                  mountPath: /etc/kubernetes
                  readOnly: true
              restartPolicy: Never
              volumes:
              - name: var-lib-etcd
                hostPath:
                  path: "/var/lib/etcd"
              - name: var-lib-kubelet
                hostPath:
                  path: "/var/lib/kubelet"
              - name: etc-kubernetes
                hostPath:
                  path: "/etc/kubernetes"
        EOF
    
    - name: Wait for Job Completion
      run: |
        kubectl wait --for=condition=complete job/kube-bench-${{ matrix.cluster }} --timeout=300s
    
    - name: Get Results
      run: |
        kubectl logs job/kube-bench-${{ matrix.cluster }} > kube-bench-results.json
        
        # Parse results for pass/fail
        FAILED_CHECKS=$(jq '.totals.total_fail' kube-bench-results.json)
        
        if [ "$FAILED_CHECKS" -gt 0 ]; then
          echo "❌ Security assessment failed with $FAILED_CHECKS failed checks"
          exit 1
        else
          echo "✅ Security assessment passed"
        fi
    
    - name: Upload Results
      uses: actions/upload-artifact@v2
      with:
        name: kube-bench-results-${{ matrix.cluster }}
        path: kube-bench-results.json
    
    - name: Cleanup
      if: always()
      run: kubectl delete job kube-bench-${{ matrix.cluster }}
```

### AWS Security Hub Integration
```bash
#!/bin/bash
# aws-security-hub-integration.sh

# Configure AWS credentials and region
export AWS_REGION="us-west-2"
export AWS_PROFILE="security"

echo "Running kube-bench with AWS Security Hub integration..."

# Run kube-bench and send findings to Security Hub
kube-bench run --targets master,node,etcd --asff

# Alternative: Run with custom filtering and send to Security Hub
kube-bench --scored --asff

# Generate local report for additional processing
kube-bench run --targets master,node,etcd \
    --json \
    --outputfile security-findings.json

# Process results for custom alerting
CRITICAL_FAILURES=$(jq '[.tests[].results[] | select(.result == "FAIL" and .scored == true)] | length' security-findings.json)

if [ "$CRITICAL_FAILURES" -gt 0 ]; then
    echo "Found $CRITICAL_FAILURES critical security failures"
    
    # Send custom notification
    aws sns publish \
        --topic-arn "arn:aws:sns:us-west-2:123456789012:security-alerts" \
        --message "Kubernetes security scan found $CRITICAL_FAILURES critical failures. Check AWS Security Hub for details." \
        --subject "Kubernetes Security Alert"
fi

echo "Security assessment complete. Check AWS Security Hub for detailed findings."
```

### Custom Configuration and Benchmarks
```bash
#!/bin/bash
# custom-benchmark-config.sh

# Create custom configuration directory
mkdir -p /opt/kube-bench/custom-cfg

# Download and customize benchmark configuration
wget -O /opt/kube-bench/custom-cfg/config.yaml \
    https://raw.githubusercontent.com/aquasecurity/kube-bench/main/cfg/config.yaml

# Customize checks for specific environment
cat > /opt/kube-bench/custom-cfg/master.yaml <<EOF
controls:
id: 1
text: "Master Node Security Configuration"
type: "master"
groups:
- id: 1.1
  text: "Master Node Configuration Files"
  checks:
  - id: 1.1.1
    text: "Ensure that the API server pod specification file permissions are set to 644 or more restrictive"
    audit: "stat -c %a /etc/kubernetes/manifests/kube-apiserver.yaml"
    tests:
      test_items:
      - flag: "644"
        compare:
          op: eq
          value: "644"
    remediation: "chmod 644 /etc/kubernetes/manifests/kube-apiserver.yaml"
    scored: true
EOF

# Run with custom configuration
kube-bench run --targets master \
    --config-dir /opt/kube-bench/custom-cfg \
    --json \
    --outputfile custom-benchmark-results.json

echo "Custom benchmark assessment completed"
```

## Integration Patterns

### Kubernetes Job Deployment
```yaml
# kube-bench-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: kube-bench
spec:
  template:
    spec:
      hostPID: true
      nodeSelector:
        kubernetes.io/os: linux
      tolerations:
      - operator: Exists
        effect: NoSchedule
      containers:
      - name: kube-bench
        image: aquasec/kube-bench:latest
        command: ["kube-bench"]
        args: ["run", "--targets", "node", "--json"]
        volumeMounts:
        - name: var-lib-etcd
          mountPath: /var/lib/etcd
          readOnly: true
        - name: var-lib-kubelet
          mountPath: /var/lib/kubelet
          readOnly: true
        - name: etc-kubernetes
          mountPath: /etc/kubernetes
          readOnly: true
        - name: usr-bin
          mountPath: /usr/local/mount-from-host/bin
          readOnly: true
      restartPolicy: Never
      volumes:
      - name: var-lib-etcd
        hostPath:
          path: "/var/lib/etcd"
      - name: var-lib-kubelet
        hostPath:
          path: "/var/lib/kubelet"
      - name: etc-kubernetes
        hostPath:
          path: "/etc/kubernetes"
      - name: usr-bin
        hostPath:
          path: "/usr/bin"
  backoffLimit: 2
```

### Terraform Integration
```hcl
# terraform/kube-bench-cronjob.tf
resource "kubernetes_cron_job" "kube_bench" {
  metadata {
    name      = "kube-bench-scan"
    namespace = "security"
  }

  spec {
    schedule = "0 2 * * 1"  # Weekly on Monday at 2 AM
    
    job_template {
      metadata {
        labels = {
          app = "kube-bench"
        }
      }
      
      spec {
        template {
          metadata {
            labels = {
              app = "kube-bench"
            }
          }
          
          spec {
            host_pid = true
            
            node_selector = {
              "kubernetes.io/os" = "linux"
            }
            
            toleration {
              operator = "Exists"
              effect   = "NoSchedule"
            }
            
            container {
              name  = "kube-bench"
              image = "aquasec/kube-bench:latest"
              
              command = ["kube-bench"]
              args    = ["run", "--targets", "node", "--json"]
              
              volume_mount {
                name       = "var-lib-kubelet"
                mount_path = "/var/lib/kubelet"
                read_only  = true
              }
              
              volume_mount {
                name       = "etc-kubernetes"
                mount_path = "/etc/kubernetes"
                read_only  = true
              }
            }
            
            volume {
              name = "var-lib-kubelet"
              host_path {
                path = "/var/lib/kubelet"
              }
            }
            
            volume {
              name = "etc-kubernetes"
              host_path {
                path = "/etc/kubernetes"
              }
            }
            
            restart_policy = "Never"
          }
        }
      }
    }
  }
}
```

## Best Practices

### Security Assessment Strategy
- **Regular Scanning**: Schedule weekly or monthly security assessments
- **Baseline Establishment**: Create security baselines for different environments
- **Trend Analysis**: Track security improvements over time
- **Priority Remediation**: Focus on scored checks and critical failures first

### Implementation Approach
- **Staged Rollout**: Start with development clusters before production
- **Custom Configurations**: Adapt benchmarks to organizational requirements
- **Automated Remediation**: Develop scripts for common remediation tasks
- **Documentation**: Maintain records of assessments and remediation actions

### Integration Strategy
- **CI/CD Gates**: Block deployments that fail critical security checks
- **Monitoring Integration**: Send results to security monitoring systems
- **Incident Response**: Integrate findings into security incident workflows
- **Compliance Reporting**: Generate reports for audit and compliance purposes

### Operational Considerations
- **Node Access**: Ensure kube-bench has necessary access to node filesystems
- **Performance Impact**: Schedule scans during low-traffic periods
- **Results Storage**: Implement secure storage for security assessment results
- **Access Control**: Limit access to security assessment tools and results

## Error Handling

### Common Issues
```bash
# Permission denied accessing config files
sudo kube-bench run --targets master
# Solution: Run with appropriate privileges

# No Kubernetes version detected
kube-bench --version 1.24
# Solution: Specify version manually

# Custom benchmark not found
kube-bench --benchmark custom-1.0 --config-dir /opt/custom-cfg
# Solution: Verify benchmark files exist

# Job fails in managed clusters
kube-bench run --targets node
# Solution: Only scan accessible components
```

### Troubleshooting
- **Debug Mode**: Use `-v 3` for detailed logging
- **Configuration**: Verify config files and permissions
- **Network Access**: Ensure connectivity to Kubernetes API
- **Platform Compatibility**: Use appropriate benchmarks for your platform

Kube-bench provides comprehensive Kubernetes security assessment capabilities through industry-standard CIS benchmarks, enabling organizations to maintain secure and compliant Kubernetes deployments.