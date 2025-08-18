# Litmus

Kubernetes-native chaos engineering platform for building resilient systems through controlled failure injection.

## Description

Litmus is a comprehensive chaos engineering platform designed specifically for Kubernetes environments. It enables teams to introduce controlled failures into their systems to identify weaknesses, improve resilience, and build confidence in system reliability. Litmus provides a complete framework for designing, executing, and monitoring chaos experiments through a combination of custom Kubernetes resources, a powerful CLI tool (litmusctl), and an intuitive web interface (ChaosCenter). The platform supports both declarative chaos experiments defined as YAML manifests and programmatic chaos workflows for complex scenarios.

## MCP Tools

### Platform Management
- **`litmus_install`** - Install Litmus chaos engineering platform using Helm
- **`litmus_config_set_account`** - Setup ChaosCenter account configuration
- **`litmus_version`** - Get litmusctl version information

### Infrastructure Management
- **`litmus_connect_chaos_infra`** - Connect chaos infrastructure using litmusctl
- **`litmus_get_chaos_infra`** - List chaos infrastructure

### Project Management
- **`litmus_create_project`** - Create a new project using litmusctl
- **`litmus_get_projects`** - List projects using litmusctl

### Experiment Management
- **`litmus_create_chaos_experiment`** - Create chaos experiment using litmusctl
- **`litmus_run_chaos_experiment`** - Run chaos experiment using litmusctl
- **`litmus_get_chaos_experiments`** - List chaos experiments using litmusctl
- **`litmus_apply_chaos_experiment`** - Apply chaos experiment manifest using kubectl
- **`litmus_get_chaos_results`** - Get chaos experiment results using kubectl

## Real CLI Commands Used

### Installation Commands
- `helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/` - Add Litmus Helm repository
- `helm repo update` - Update Helm repositories
- `helm install chaos litmuschaos/litmus --namespace litmus` - Install Litmus platform
- `kubectl create namespace litmus` - Create Litmus namespace

### litmusctl Commands
- `litmusctl version` - Show litmusctl version
- `litmusctl config set-account` - Configure ChaosCenter account
- `litmusctl connect chaos-infra` - Connect chaos infrastructure
- `litmusctl create project` - Create new project
- `litmusctl create chaos-experiment -f <file>` - Create experiment from manifest
- `litmusctl run chaos-experiment <id>` - Run chaos experiment
- `litmusctl get projects` - List all projects
- `litmusctl get chaos-experiment` - List experiments
- `litmusctl get chaos-infra` - List infrastructure

### kubectl Commands
- `kubectl apply -f <experiment.yaml>` - Apply chaos experiment manifest
- `kubectl get chaosresult` - Get experiment results
- `kubectl get chaosengine` - Get experiment execution status
- `kubectl get chaosexperiment` - List available experiments
- `kubectl describe chaosresult <name>` - Get detailed results

### Helm Management
- `helm list -n litmus` - List Litmus releases
- `helm upgrade chaos litmuschaos/litmus -n litmus` - Upgrade Litmus
- `helm uninstall chaos -n litmus` - Uninstall Litmus

## Use Cases

### Resilience Testing
- **Application Resilience**: Test application behavior under various failure conditions
- **Infrastructure Resilience**: Validate infrastructure stability and recovery capabilities
- **Network Resilience**: Test system behavior under network partitions and latency
- **Resource Exhaustion**: Validate application behavior under CPU, memory, and storage stress

### SRE and DevOps
- **Chaos Engineering**: Implement systematic chaos engineering practices
- **Production Testing**: Safely test production systems for unknown failure modes
- **Incident Response**: Validate incident response procedures and runbooks
- **Monitoring Validation**: Ensure monitoring and alerting systems detect failures

### Compliance and Governance
- **Disaster Recovery**: Test disaster recovery procedures and RTO/RPO objectives
- **SLA Validation**: Verify system meets availability and performance SLAs
- **Security Testing**: Test system resilience against security-related failures
- **Regulatory Compliance**: Demonstrate system resilience for compliance requirements

### Development and Testing
- **CI/CD Integration**: Integrate chaos testing into deployment pipelines
- **Pre-Production Testing**: Validate changes before production deployment
- **Load Testing**: Combine chaos engineering with performance testing
- **Quality Assurance**: Enhance QA processes with chaos testing scenarios

## Configuration Examples

### Basic Litmus Installation
```bash
# Add Litmus Helm repository
helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
helm repo update

# Create namespace for Litmus
kubectl create namespace litmus

# Install Litmus platform
helm install chaos litmuschaos/litmus --namespace litmus

# Verify installation
kubectl get pods -n litmus
kubectl get services -n litmus

# Access ChaosCenter (adjust based on service type)
kubectl port-forward svc/chaos-litmus-frontend-service 9091:9091 -n litmus
# Access at http://localhost:9091
```

### litmusctl Configuration
```bash
# Install litmusctl
curl -O https://litmusctl-bucket.s3-us-west-2.amazonaws.com/litmusctl-linux-amd64-master.tar.gz
tar -zxvf litmusctl-linux-amd64-master.tar.gz
chmod +x litmusctl
sudo mv litmusctl /usr/local/bin/

# Check version
litmusctl version

# Configure ChaosCenter account
litmusctl config set-account
# Follow interactive prompts to set:
# - ChaosCenter URL
# - Username/Password
# - Project details

# View current configuration
litmusctl config view

# Create a new project
litmusctl create project

# Connect chaos infrastructure
litmusctl connect chaos-infra
```

### Basic Chaos Experiment
```bash
# Create a simple pod-delete experiment manifest
cat <<EOF > pod-delete-experiment.yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: nginx-chaos
  namespace: default
spec:
  engineState: 'active'
  appinfo:
    appns: 'default'
    applabel: 'app=nginx'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-delete
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '30'
        - name: CHAOS_INTERVAL
          value: '10'
        - name: FORCE
          value: 'false'
EOF

# Apply the experiment
kubectl apply -f pod-delete-experiment.yaml

# Monitor experiment progress
kubectl get chaosengine nginx-chaos -w

# Check results
kubectl get chaosresult
kubectl describe chaosresult nginx-chaos-pod-delete
```

### Advanced Chaos Workflow
```bash
# Create a complex workflow with multiple experiments
cat <<EOF > complex-chaos-workflow.yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  name: chaos-workflow
  namespace: litmus
spec:
  entrypoint: chaos-tests
  serviceAccountName: argo-chaos
  templates:
  - name: chaos-tests
    steps:
    - - name: pod-delete-test
        template: pod-delete
    - - name: memory-hog-test
        template: memory-hog
    - - name: network-partition-test
        template: network-partition
        
  - name: pod-delete
    container:
      image: litmuschaos/litmus-checker:latest
      command: [sh, -c]
      args: ["kubectl apply -f /chaos/pod-delete-experiment.yaml && sleep 60"]
      
  - name: memory-hog
    container:
      image: litmuschaos/litmus-checker:latest
      command: [sh, -c]
      args: ["kubectl apply -f /chaos/memory-hog-experiment.yaml && sleep 120"]
      
  - name: network-partition
    container:
      image: litmuschaos/litmus-checker:latest
      command: [sh, -c]
      args: ["kubectl apply -f /chaos/network-partition-experiment.yaml && sleep 180"]
EOF

# Apply the workflow
kubectl apply -f complex-chaos-workflow.yaml

# Monitor workflow
kubectl get workflow chaos-workflow -w
```

## Advanced Usage

### Comprehensive Chaos Engineering Setup
```bash
#!/bin/bash
# setup-chaos-engineering.sh

echo "Setting up comprehensive chaos engineering environment..."

# Install Litmus platform
echo "Installing Litmus platform..."
helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
helm repo update

kubectl create namespace litmus
helm install chaos litmuschaos/litmus --namespace litmus --wait

# Wait for Litmus to be ready
echo "Waiting for Litmus to be ready..."
kubectl wait --for=condition=Ready pods --all -n litmus --timeout=300s

# Install chaos experiments
echo "Installing chaos experiments..."
kubectl apply -f https://hub.litmuschaos.io/api/chaos/master?file=charts/generic/experiments.yaml

# Create service account for experiments
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: litmus-admin
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: litmus-admin
rules:
- apiGroups: [""]
  resources: ["pods", "events", "configmaps", "secrets", "services"]
  verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "replicasets", "daemonsets"]
  verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
- apiGroups: ["litmuschaos.io"]
  resources: ["chaosengines", "chaosexperiments", "chaosresults"]
  verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: litmus-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: litmus-admin
subjects:
- kind: ServiceAccount
  name: litmus-admin
  namespace: default
EOF

# Install sample application for testing
echo "Installing sample application..."
kubectl create deployment nginx --image=nginx --replicas=3
kubectl expose deployment nginx --port=80 --type=ClusterIP
kubectl label deployment nginx app=nginx

# Create basic chaos experiments
echo "Creating basic chaos experiments..."

# Pod delete experiment
cat <<EOF | kubectl apply -f -
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: nginx-pod-delete
  namespace: default
spec:
  engineState: 'active'
  appinfo:
    appns: 'default'
    applabel: 'app=nginx'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-delete
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '30'
        - name: CHAOS_INTERVAL
          value: '10'
        - name: FORCE
          value: 'false'
EOF

# Memory hog experiment
cat <<EOF | kubectl apply -f -
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: nginx-memory-hog
  namespace: default
spec:
  engineState: 'stop'  # Start in stopped state
  appinfo:
    appns: 'default'
    applabel: 'app=nginx'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-memory-hog
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '60'
        - name: MEMORY_CONSUMPTION
          value: '500'
        - name: NUMBER_OF_WORKERS
          value: '1'
EOF

echo "Chaos engineering environment setup complete!"
echo "Access ChaosCenter: kubectl port-forward svc/chaos-litmus-frontend-service 9091:9091 -n litmus"
echo "Username: admin, Password: litmus"
```

### Automated Chaos Testing Pipeline
```bash
#!/bin/bash
# chaos-testing-pipeline.sh

NAMESPACE="$1"
APPLICATION="$2"
EXPERIMENT_TYPE="$3"

if [[ -z "$NAMESPACE" || -z "$APPLICATION" || -z "$EXPERIMENT_TYPE" ]]; then
    echo "Usage: $0 <namespace> <application> <experiment-type>"
    echo "Experiment types: pod-delete, memory-hog, cpu-hog, network-latency"
    exit 1
fi

echo "Starting chaos testing pipeline..."
echo "Namespace: $NAMESPACE"
echo "Application: $APPLICATION"
echo "Experiment: $EXPERIMENT_TYPE"

# Verify application is running
echo "Verifying application health..."
REPLICAS=$(kubectl get deployment $APPLICATION -n $NAMESPACE -o jsonpath='{.status.readyReplicas}')
if [[ "$REPLICAS" -eq 0 ]]; then
    echo "Application $APPLICATION has no ready replicas. Aborting."
    exit 1
fi

echo "Application has $REPLICAS ready replicas. Proceeding with chaos test."

# Create chaos experiment based on type
EXPERIMENT_NAME="${APPLICATION}-${EXPERIMENT_TYPE}-$(date +%s)"

case $EXPERIMENT_TYPE in
    "pod-delete")
        CHAOS_MANIFEST=$(cat <<EOF
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: $EXPERIMENT_NAME
  namespace: $NAMESPACE
spec:
  engineState: 'active'
  appinfo:
    appns: '$NAMESPACE'
    applabel: 'app=$APPLICATION'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-delete
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '30'
        - name: CHAOS_INTERVAL
          value: '10'
        - name: FORCE
          value: 'false'
        - name: PODS_AFFECTED_PERC
          value: '50'
EOF
)
        ;;
    "memory-hog")
        CHAOS_MANIFEST=$(cat <<EOF
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: $EXPERIMENT_NAME
  namespace: $NAMESPACE
spec:
  engineState: 'active'
  appinfo:
    appns: '$NAMESPACE'
    applabel: 'app=$APPLICATION'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-memory-hog
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '60'
        - name: MEMORY_CONSUMPTION
          value: '500'
        - name: NUMBER_OF_WORKERS
          value: '1'
        - name: PODS_AFFECTED_PERC
          value: '50'
EOF
)
        ;;
    "cpu-hog")
        CHAOS_MANIFEST=$(cat <<EOF
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: $EXPERIMENT_NAME
  namespace: $NAMESPACE
spec:
  engineState: 'active'
  appinfo:
    appns: '$NAMESPACE'
    applabel: 'app=$APPLICATION'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-cpu-hog
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '60'
        - name: CPU_CORES
          value: '1'
        - name: PODS_AFFECTED_PERC
          value: '50'
EOF
)
        ;;
    *)
        echo "Unknown experiment type: $EXPERIMENT_TYPE"
        exit 1
        ;;
esac

# Apply chaos experiment
echo "Applying chaos experiment..."
echo "$CHAOS_MANIFEST" | kubectl apply -f -

# Monitor experiment progress
echo "Monitoring experiment progress..."
timeout 300 kubectl wait --for=condition=ExperimentStatus=completed chaosengine $EXPERIMENT_NAME -n $NAMESPACE

# Get experiment results
echo "Getting experiment results..."
kubectl get chaosresult ${EXPERIMENT_NAME}-${EXPERIMENT_TYPE} -n $NAMESPACE -o yaml

# Verify application recovery
echo "Verifying application recovery..."
sleep 30
FINAL_REPLICAS=$(kubectl get deployment $APPLICATION -n $NAMESPACE -o jsonpath='{.status.readyReplicas}')

if [[ "$FINAL_REPLICAS" -eq "$REPLICAS" ]]; then
    echo "✅ Application recovered successfully. All replicas are ready."
    RESULT="PASSED"
else
    echo "❌ Application did not recover fully. Expected: $REPLICAS, Current: $FINAL_REPLICAS"
    RESULT="FAILED"
fi

# Generate report
echo "=== Chaos Test Report ===" > chaos-report-${EXPERIMENT_NAME}.txt
echo "Date: $(date)" >> chaos-report-${EXPERIMENT_NAME}.txt
echo "Namespace: $NAMESPACE" >> chaos-report-${EXPERIMENT_NAME}.txt
echo "Application: $APPLICATION" >> chaos-report-${EXPERIMENT_NAME}.txt
echo "Experiment: $EXPERIMENT_TYPE" >> chaos-report-${EXPERIMENT_NAME}.txt
echo "Result: $RESULT" >> chaos-report-${EXPERIMENT_NAME}.txt
echo "Initial Replicas: $REPLICAS" >> chaos-report-${EXPERIMENT_NAME}.txt
echo "Final Replicas: $FINAL_REPLICAS" >> chaos-report-${EXPERIMENT_NAME}.txt

echo "Chaos test completed. Report: chaos-report-${EXPERIMENT_NAME}.txt"

# Cleanup experiment
kubectl delete chaosengine $EXPERIMENT_NAME -n $NAMESPACE

exit $([[ "$RESULT" == "PASSED" ]] && echo 0 || echo 1)
```

### Multi-Environment Chaos Testing
```bash
#!/bin/bash
# multi-environment-chaos.sh

ENVIRONMENTS=("development" "staging" "production")
APPLICATIONS=("frontend" "backend" "database")
EXPERIMENTS=("pod-delete" "memory-hog")

DATE=$(date +%Y%m%d)
REPORT_DIR="chaos-test-results-$DATE"
mkdir -p $REPORT_DIR

echo "Starting multi-environment chaos testing..."

for env in "${ENVIRONMENTS[@]}"; do
    echo "Testing environment: $env"
    
    # Check if environment namespace exists
    if ! kubectl get namespace $env > /dev/null 2>&1; then
        echo "Namespace $env does not exist. Skipping."
        continue
    fi
    
    mkdir -p $REPORT_DIR/$env
    
    for app in "${APPLICATIONS[@]}"; do
        echo "  Testing application: $app"
        
        # Check if application exists
        if ! kubectl get deployment $app -n $env > /dev/null 2>&1; then
            echo "    Application $app not found in $env. Skipping."
            continue
        fi
        
        for experiment in "${EXPERIMENTS[@]}"; do
            echo "    Running experiment: $experiment"
            
            # Run chaos experiment
            EXPERIMENT_NAME="${app}-${experiment}-$(date +%s)"
            
            case $experiment in
                "pod-delete")
                    cat <<EOF | kubectl apply -f -
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: $EXPERIMENT_NAME
  namespace: $env
spec:
  engineState: 'active'
  appinfo:
    appns: '$env'
    applabel: 'app=$app'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-delete
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '30'
        - name: CHAOS_INTERVAL
          value: '10'
        - name: FORCE
          value: 'false'
EOF
                    ;;
                "memory-hog")
                    cat <<EOF | kubectl apply -f -
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: $EXPERIMENT_NAME
  namespace: $env
spec:
  engineState: 'active'
  appinfo:
    appns: '$env'
    applabel: 'app=$app'
    appkind: 'deployment'
  chaosServiceAccount: litmus-admin
  experiments:
  - name: pod-memory-hog
    spec:
      components:
        env:
        - name: TOTAL_CHAOS_DURATION
          value: '60'
        - name: MEMORY_CONSUMPTION
          value: '500'
EOF
                    ;;
            esac
            
            # Wait for experiment to complete
            timeout 180 kubectl wait --for=condition=ExperimentStatus=completed chaosengine $EXPERIMENT_NAME -n $env || true
            
            # Get results
            kubectl get chaosresult ${EXPERIMENT_NAME}-${experiment} -n $env -o yaml > $REPORT_DIR/$env/${app}-${experiment}-result.yaml 2>/dev/null || echo "No results found"
            
            # Cleanup
            kubectl delete chaosengine $EXPERIMENT_NAME -n $env
            
            sleep 30  # Wait between experiments
        done
    done
done

# Generate summary report
echo "=== Multi-Environment Chaos Testing Summary ===" > $REPORT_DIR/summary.txt
echo "Date: $(date)" >> $REPORT_DIR/summary.txt
echo "Environments: ${ENVIRONMENTS[*]}" >> $REPORT_DIR/summary.txt
echo "Applications: ${APPLICATIONS[*]}" >> $REPORT_DIR/summary.txt
echo "Experiments: ${EXPERIMENTS[*]}" >> $REPORT_DIR/summary.txt
echo "" >> $REPORT_DIR/summary.txt

for env in "${ENVIRONMENTS[@]}"; do
    if [[ -d $REPORT_DIR/$env ]]; then
        RESULT_COUNT=$(find $REPORT_DIR/$env -name "*.yaml" | wc -l)
        echo "$env: $RESULT_COUNT experiment results" >> $REPORT_DIR/summary.txt
    fi
done

echo "Multi-environment chaos testing complete!"
echo "Results available in: $REPORT_DIR/"
```

### Continuous Chaos Monitoring
```bash
#!/bin/bash
# continuous-chaos-monitoring.sh

SLACK_WEBHOOK="$1"
MONITORING_INTERVAL=3600  # 1 hour

if [[ -z "$SLACK_WEBHOOK" ]]; then
    echo "Usage: $0 <slack-webhook-url>"
    exit 1
fi

echo "Starting continuous chaos monitoring..."

while true; do
    echo "Running chaos health check at $(date)"
    
    # Check Litmus platform health
    LITMUS_PODS=$(kubectl get pods -n litmus --no-headers | grep -v Running | wc -l)
    
    if [[ $LITMUS_PODS -gt 0 ]]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"🚨 Litmus platform has $LITMUS_PODS unhealthy pods\"}" \
            $SLACK_WEBHOOK
    fi
    
    # Check for failed experiments
    FAILED_EXPERIMENTS=$(kubectl get chaosresult --all-namespaces --no-headers -o custom-columns="NAMESPACE:.metadata.namespace,NAME:.metadata.name,VERDICT:.status.experimentStatus.verdict" | grep -v Pass | wc -l)
    
    if [[ $FAILED_EXPERIMENTS -gt 0 ]]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"⚠️ Found $FAILED_EXPERIMENTS failed chaos experiments\"}" \
            $SLACK_WEBHOOK
    fi
    
    # Check experiment age
    OLD_EXPERIMENTS=$(kubectl get chaosengine --all-namespaces --no-headers -o custom-columns="AGE:.metadata.creationTimestamp" | while read timestamp; do
        if [[ -n "$timestamp" ]]; then
            AGE_SECONDS=$(( $(date +%s) - $(date -d "$timestamp" +%s) ))
            if [[ $AGE_SECONDS -gt 86400 ]]; then  # 24 hours
                echo "old"
            fi
        fi
    done | wc -l)
    
    if [[ $OLD_EXPERIMENTS -gt 0 ]]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"🧹 Found $OLD_EXPERIMENTS chaos experiments older than 24 hours - consider cleanup\"}" \
            $SLACK_WEBHOOK
    fi
    
    echo "Health check complete. Next check in $MONITORING_INTERVAL seconds."
    sleep $MONITORING_INTERVAL
done
```

## Integration Patterns

### CI/CD Pipeline Integration
```yaml
# .github/workflows/chaos-testing.yml
name: Chaos Engineering Tests
on:
  workflow_dispatch:
    inputs:
      environment:
        description: 'Environment to test'
        required: true
        type: choice
        options:
          - staging
          - production
      experiment_type:
        description: 'Type of chaos experiment'
        required: true
        type: choice
        options:
          - pod-delete
          - memory-hog
          - cpu-hog
          - network-latency

jobs:
  chaos-test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup kubectl
      uses: azure/setup-kubectl@v1
      
    - name: Configure cluster access
      run: |
        echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
        
    - name: Install litmusctl
      run: |
        curl -O https://litmusctl-bucket.s3-us-west-2.amazonaws.com/litmusctl-linux-amd64-master.tar.gz
        tar -zxvf litmusctl-linux-amd64-master.tar.gz
        chmod +x litmusctl
        sudo mv litmusctl /usr/local/bin/
        
    - name: Verify Litmus Installation
      run: |
        kubectl get pods -n litmus
        litmusctl version
        
    - name: Run Chaos Experiment
      env:
        ENVIRONMENT: ${{ github.event.inputs.environment }}
        EXPERIMENT: ${{ github.event.inputs.experiment_type }}
      run: |
        # Create experiment manifest
        cat > chaos-experiment.yaml <<EOF
        apiVersion: litmuschaos.io/v1alpha1
        kind: ChaosEngine
        metadata:
          name: ci-chaos-${{ github.run_id }}
          namespace: $ENVIRONMENT
        spec:
          engineState: 'active'
          appinfo:
            appns: '$ENVIRONMENT'
            applabel: 'app=frontend'
            appkind: 'deployment'
          chaosServiceAccount: litmus-admin
          experiments:
          - name: $EXPERIMENT
            spec:
              components:
                env:
                - name: TOTAL_CHAOS_DURATION
                  value: '60'
                - name: CHAOS_INTERVAL
                  value: '10'
        EOF
        
        # Apply experiment
        kubectl apply -f chaos-experiment.yaml
        
        # Wait for completion
        timeout 300 kubectl wait --for=condition=ExperimentStatus=completed chaosengine ci-chaos-${{ github.run_id }} -n $ENVIRONMENT
        
    - name: Collect Results
      run: |
        kubectl get chaosresult ci-chaos-${{ github.run_id }}-${{ github.event.inputs.experiment_type }} -n ${{ github.event.inputs.environment }} -o yaml > chaos-result.yaml
        
    - name: Validate Application Health
      run: |
        # Check if application recovered
        kubectl get deployment frontend -n ${{ github.event.inputs.environment }}
        kubectl wait --for=condition=Available deployment/frontend -n ${{ github.event.inputs.environment }} --timeout=300s
        
    - name: Upload Results
      uses: actions/upload-artifact@v2
      with:
        name: chaos-test-results
        path: chaos-result.yaml
        
    - name: Cleanup
      if: always()
      run: |
        kubectl delete chaosengine ci-chaos-${{ github.run_id }} -n ${{ github.event.inputs.environment }}
```

### Terraform Integration
```hcl
# terraform/litmus-setup.tf
resource "kubernetes_namespace" "litmus" {
  metadata {
    name = "litmus"
  }
}

resource "helm_release" "litmus" {
  name       = "chaos"
  repository = "https://litmuschaos.github.io/litmus-helm/"
  chart      = "litmus"
  namespace  = kubernetes_namespace.litmus.metadata[0].name

  values = [
    yamlencode({
      portal = {
        frontend = {
          service = {
            type = "LoadBalancer"
          }
        }
      }
      mongodb = {
        persistence = {
          size = "20Gi"
        }
      }
    })
  ]

  depends_on = [kubernetes_namespace.litmus]
}

# Service account for chaos experiments
resource "kubernetes_service_account" "litmus_admin" {
  metadata {
    name      = "litmus-admin"
    namespace = "default"
  }
}

resource "kubernetes_cluster_role" "litmus_admin" {
  metadata {
    name = "litmus-admin"
  }

  rule {
    api_groups = [""]
    resources  = ["pods", "events", "configmaps", "secrets", "services"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "statefulsets", "replicasets", "daemonsets"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }

  rule {
    api_groups = ["litmuschaos.io"]
    resources  = ["chaosengines", "chaosexperiments", "chaosresults"]
    verbs      = ["create", "delete", "get", "list", "patch", "update", "watch"]
  }
}

resource "kubernetes_cluster_role_binding" "litmus_admin" {
  metadata {
    name = "litmus-admin"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.litmus_admin.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.litmus_admin.metadata[0].name
    namespace = "default"
  }
}

# Install chaos experiments
resource "kubectl_manifest" "chaos_experiments" {
  yaml_body = file("${path.module}/chaos-experiments.yaml")
  
  depends_on = [helm_release.litmus]
}

output "chaos_center_url" {
  value = "Access ChaosCenter at the LoadBalancer IP on port 9091"
}
```

### ArgoCD Integration
```yaml
# argocd/litmus-application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: litmus
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://litmuschaos.github.io/litmus-helm/
    chart: litmus
    targetRevision: 3.0.0
    helm:
      values: |
        portal:
          frontend:
            service:
              type: ClusterIP
        mongodb:
          persistence:
            size: 20Gi
  destination:
    server: https://kubernetes.default.svc
    namespace: litmus
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true

---
# Chaos experiment application
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: chaos-experiments
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/company/chaos-experiments
    path: experiments
    targetRevision: HEAD
  destination:
    server: https://kubernetes.default.svc
    namespace: chaos-experiments
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

## Best Practices

### Experiment Design
- **Start Small**: Begin with low-impact experiments and gradually increase complexity
- **Hypothesis-Driven**: Define clear hypotheses and success criteria before experiments
- **Gradual Rollout**: Test in development → staging → production environments
- **Blast Radius Control**: Limit the scope and impact of chaos experiments

### Safety Measures
- **Circuit Breakers**: Implement automatic experiment termination on critical failures
- **Monitoring Integration**: Ensure comprehensive monitoring during experiments
- **Rollback Plans**: Have clear rollback procedures for each experiment
- **Business Hours**: Run initial experiments during business hours with full team availability

### Operational Excellence
- **Regular Execution**: Run chaos experiments regularly, not just during development
- **Incident Learning**: Use chaos engineering to validate incident response procedures
- **Documentation**: Document experiment results and lessons learned
- **Team Training**: Train teams on chaos engineering principles and tools

### Platform Management
- **Resource Limits**: Set appropriate resource limits for chaos experiment pods
- **RBAC**: Implement proper role-based access control for chaos operations
- **Audit Logging**: Enable audit logging for all chaos engineering activities
- **High Availability**: Ensure Litmus platform itself is highly available

## Error Handling

### Common Issues
```bash
# Litmus pods not starting
kubectl describe pods -n litmus
# Check resource constraints and node capacity

# Experiments failing to start
kubectl get chaosengine -A
kubectl describe chaosengine <name> -n <namespace>
# Check RBAC permissions and service account

# ChaosCenter not accessible
kubectl get svc -n litmus
kubectl port-forward svc/chaos-litmus-frontend-service 9091:9091 -n litmus
# Check service type and network policies

# litmusctl connection issues
litmusctl config view
# Verify ChaosCenter URL and credentials
```

### Troubleshooting
- **Permission Issues**: Verify service account has necessary RBAC permissions
- **Network Connectivity**: Check network policies and service accessibility
- **Resource Constraints**: Ensure sufficient cluster resources for experiments
- **Version Compatibility**: Verify Litmus version compatibility with Kubernetes

Litmus provides comprehensive chaos engineering capabilities for Kubernetes environments, enabling teams to build more resilient systems through systematic failure injection and testing methodologies.