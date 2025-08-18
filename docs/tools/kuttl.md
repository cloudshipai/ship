# KUTTL

Kubernetes Test Tool (KUTTL) for declarative testing of Kubernetes operators and controllers.

## Description

KUTTL (KUbernetes Test TooL) is a declarative testing framework designed specifically for testing Kubernetes operators and controllers. It provides a simple way to write tests for Kubernetes workloads by allowing developers to define test cases using YAML files that describe the expected behavior and state of resources. KUTTL excels at testing complex scenarios like operator upgrades, failure conditions, and multi-step workflows, making it an essential tool for ensuring the reliability of Kubernetes applications and operators.

## MCP Tools

### Core Testing
- **`kuttl_test`** - Run Kubernetes tests using kubectl kuttl test with comprehensive configuration options
- **`kuttl_test_kind`** - Run KUTTL tests with automated KIND cluster setup for isolated testing

### Information and Help
- **`kuttl_version`** - Get KUTTL version information
- **`kuttl_help`** - Get KUTTL help information for commands

## Real CLI Commands Used

### Primary Commands
- `kubectl kuttl test <path>` - Run KUTTL test suite
- `kubectl kuttl version` - Show KUTTL version
- `kubectl kuttl help [command]` - Show general or command-specific help

### Test Configuration Options
- `kubectl kuttl test --config <file>` - Use specific test settings file
- `kubectl kuttl test --artifacts-dir <dir>` - Directory for test output and logs
- `kubectl kuttl test --crd-dir <dir>` - Apply CustomResourceDefinitions before tests
- `kubectl kuttl test --manifest-dir <dir>` - Apply manifests before tests
- `kubectl kuttl test --parallel <n>` - Set maximum concurrent tests (default 8)

### Cluster Management
- `kubectl kuttl test --start-kind` - Start KIND cluster for testing
- `kubectl kuttl test --start-control-plane` - Start local Kubernetes control plane

### Logging and Debugging
- `kubectl kuttl test -v` - Verbose logging
- `kubectl kuttl test -vv` - Extra verbose logging

### Combined Usage Examples
- `kubectl kuttl test --start-kind --artifacts-dir ./results ./test-suite/`
- `kubectl kuttl test --config kuttl-test.yaml --parallel 4 ./tests/`

## Use Cases

### Operator Testing
- **Operator Lifecycle**: Test operator installation, upgrades, and uninstallation
- **Custom Resource Validation**: Verify custom resources are processed correctly
- **Controller Logic**: Test complex controller behaviors and reconciliation loops
- **Error Handling**: Validate operator behavior under failure conditions

### CI/CD Integration
- **Automated Testing**: Run tests as part of continuous integration pipelines
- **Release Validation**: Ensure operator releases work correctly before deployment
- **Regression Testing**: Catch regressions in operator behavior
- **Quality Gates**: Block releases that fail critical tests

### Development Workflows
- **Local Testing**: Test operator changes during development
- **End-to-End Testing**: Validate complete operator workflows
- **Integration Testing**: Test interactions between multiple operators
- **Performance Testing**: Verify operator performance under load

### Kubernetes Application Testing
- **Application Deployment**: Test complex application deployment scenarios
- **Configuration Management**: Validate configuration changes and updates
- **Scaling Behavior**: Test horizontal and vertical scaling scenarios
- **Disaster Recovery**: Test backup, restore, and failover scenarios

## Configuration Examples

### Basic Test Execution
```bash
# Run basic test suite
kubectl kuttl test ./tests/

# Run with verbose output
kubectl kuttl test -v ./tests/

# Run with extra verbose output
kubectl kuttl test -vv ./tests/

# Get version information
kubectl kuttl version

# Get help
kubectl kuttl help
kubectl kuttl help test
```

### Test Configuration
```bash
# Run with custom configuration file
kubectl kuttl test --config kuttl-test.yaml ./tests/

# Specify artifacts directory for logs
kubectl kuttl test --artifacts-dir ./test-results ./tests/

# Apply CRDs before testing
kubectl kuttl test --crd-dir ./crds ./tests/

# Apply manifests before testing
kubectl kuttl test --manifest-dir ./prereqs ./tests/

# Control test parallelism
kubectl kuttl test --parallel 2 ./tests/
```

### Cluster Management
```bash
# Start KIND cluster for testing
kubectl kuttl test --start-kind ./tests/

# Start local control plane
kubectl kuttl test --start-control-plane ./tests/

# Combined KIND with configuration
kubectl kuttl test --start-kind --config kuttl-test.yaml --artifacts-dir ./results ./tests/
```

### Complex Test Scenarios
```bash
# Full-featured test run
kubectl kuttl test \
  --start-kind \
  --config kuttl-test.yaml \
  --artifacts-dir ./test-artifacts \
  --crd-dir ./crds \
  --manifest-dir ./setup \
  --parallel 4 \
  -v \
  ./tests/

# Test with existing cluster
kubectl kuttl test \
  --config kuttl-test.yaml \
  --artifacts-dir ./results \
  --parallel 8 \
  ./operator-tests/
```

## Advanced Usage

### Comprehensive Operator Testing Suite
```bash
#!/bin/bash
# comprehensive-operator-test.sh

OPERATOR_VERSION="$1"
TEST_SUITE="$2"

if [[ -z "$OPERATOR_VERSION" || -z "$TEST_SUITE" ]]; then
    echo "Usage: $0 <operator-version> <test-suite-path>"
    exit 1
fi

DATE=$(date +%Y%m%d_%H%M%S)
RESULTS_DIR="operator-test-results-$DATE"
mkdir -p $RESULTS_DIR

echo "Starting comprehensive operator testing for version $OPERATOR_VERSION..."

# Phase 1: Basic functionality tests
echo "Phase 1: Basic functionality tests..."
kubectl kuttl test \
    --start-kind \
    --config configs/basic-test.yaml \
    --artifacts-dir $RESULTS_DIR/basic \
    --manifest-dir manifests/basic \
    --parallel 4 \
    -v \
    $TEST_SUITE/basic/

# Phase 2: Upgrade tests
echo "Phase 2: Upgrade scenario tests..."
kubectl kuttl test \
    --config configs/upgrade-test.yaml \
    --artifacts-dir $RESULTS_DIR/upgrade \
    --manifest-dir manifests/upgrade \
    --parallel 2 \
    -v \
    $TEST_SUITE/upgrade/

# Phase 3: Failure and recovery tests
echo "Phase 3: Failure and recovery tests..."
kubectl kuttl test \
    --config configs/failure-test.yaml \
    --artifacts-dir $RESULTS_DIR/failure \
    --manifest-dir manifests/failure \
    --parallel 1 \
    -vv \
    $TEST_SUITE/failure/

# Phase 4: Performance tests
echo "Phase 4: Performance and scale tests..."
kubectl kuttl test \
    --config configs/performance-test.yaml \
    --artifacts-dir $RESULTS_DIR/performance \
    --manifest-dir manifests/performance \
    --parallel 1 \
    -v \
    $TEST_SUITE/performance/

# Generate summary report
echo "Generating test summary..."
echo "=== Operator Test Summary ===" > $RESULTS_DIR/summary.txt
echo "Operator Version: $OPERATOR_VERSION" >> $RESULTS_DIR/summary.txt
echo "Test Date: $(date)" >> $RESULTS_DIR/summary.txt
echo "Test Suite: $TEST_SUITE" >> $RESULTS_DIR/summary.txt
echo "" >> $RESULTS_DIR/summary.txt

# Check results
FAILED_TESTS=0
for phase in basic upgrade failure performance; do
    if [[ -d $RESULTS_DIR/$phase ]]; then
        FAILURES=$(find $RESULTS_DIR/$phase -name "*.log" -exec grep -l "FAIL" {} \; | wc -l)
        echo "$phase tests failed: $FAILURES" >> $RESULTS_DIR/summary.txt
        FAILED_TESTS=$((FAILED_TESTS + FAILURES))
    fi
done

echo "Total failed tests: $FAILED_TESTS" >> $RESULTS_DIR/summary.txt

if [[ $FAILED_TESTS -eq 0 ]]; then
    echo "✅ All operator tests passed!"
    exit 0
else
    echo "❌ $FAILED_TESTS test(s) failed. Check $RESULTS_DIR/ for details."
    exit 1
fi
```

### Multi-Environment Testing
```bash
#!/bin/bash
# multi-environment-kuttl-test.sh

ENVIRONMENTS=("dev" "staging" "prod")
TEST_SUITES=("smoke" "integration" "e2e")

DATE=$(date +%Y%m%d)
REPORT_DIR="multi-env-test-$DATE"
mkdir -p $REPORT_DIR

for env in "${ENVIRONMENTS[@]}"; do
    echo "Testing environment: $env"
    
    # Switch to environment context
    kubectl config use-context $env-cluster
    
    mkdir -p $REPORT_DIR/$env
    
    for suite in "${TEST_SUITES[@]}"; do
        echo "  Running $suite tests in $env..."
        
        # Run test suite
        kubectl kuttl test \
            --config configs/$env-config.yaml \
            --artifacts-dir $REPORT_DIR/$env/$suite \
            --manifest-dir manifests/$env \
            --parallel 4 \
            -v \
            tests/$suite/ || echo "FAILED: $suite in $env"
    done
    
    # Generate environment summary
    echo "=== $env Environment Test Results ===" > $REPORT_DIR/$env/summary.txt
    echo "Test Date: $(date)" >> $REPORT_DIR/$env/summary.txt
    
    for suite in "${TEST_SUITES[@]}"; do
        if [[ -d $REPORT_DIR/$env/$suite ]]; then
            SUITE_FAILURES=$(find $REPORT_DIR/$env/$suite -name "*.log" -exec grep -l "FAIL" {} \; 2>/dev/null | wc -l)
            echo "$suite suite failures: $SUITE_FAILURES" >> $REPORT_DIR/$env/summary.txt
        fi
    done
    
    echo "Environment $env testing complete"
done

echo "Multi-environment testing finished!"
echo "Results available in $REPORT_DIR/"
```

### Continuous Integration Integration
```bash
#!/bin/bash
# ci-kuttl-integration.sh

# CI/CD script for KUTTL testing
BUILD_ID="$1"
COMMIT_SHA="$2"

if [[ -z "$BUILD_ID" || -z "$COMMIT_SHA" ]]; then
    echo "Usage: $0 <build-id> <commit-sha>"
    exit 1
fi

ARTIFACTS_DIR="ci-artifacts-$BUILD_ID"
mkdir -p $ARTIFACTS_DIR

echo "Starting CI KUTTL tests for build $BUILD_ID (commit: $COMMIT_SHA)"

# Quick smoke tests
echo "Running smoke tests..."
if kubectl kuttl test \
    --start-kind \
    --config ci/smoke-test.yaml \
    --artifacts-dir $ARTIFACTS_DIR/smoke \
    --parallel 8 \
    -v \
    tests/smoke/; then
    echo "✅ Smoke tests passed"
else
    echo "❌ Smoke tests failed"
    exit 1
fi

# Integration tests
echo "Running integration tests..."
if kubectl kuttl test \
    --config ci/integration-test.yaml \
    --artifacts-dir $ARTIFACTS_DIR/integration \
    --crd-dir crds/ \
    --manifest-dir ci/setup \
    --parallel 4 \
    -v \
    tests/integration/; then
    echo "✅ Integration tests passed"
else
    echo "❌ Integration tests failed"
    # Upload artifacts for debugging
    tar -czf integration-test-artifacts.tar.gz $ARTIFACTS_DIR/integration/
    echo "Integration test artifacts: integration-test-artifacts.tar.gz"
    exit 1
fi

# End-to-end tests (if not PR)
if [[ "$GITHUB_EVENT_NAME" != "pull_request" ]]; then
    echo "Running end-to-end tests..."
    if kubectl kuttl test \
        --config ci/e2e-test.yaml \
        --artifacts-dir $ARTIFACTS_DIR/e2e \
        --parallel 2 \
        -v \
        tests/e2e/; then
        echo "✅ End-to-end tests passed"
    else
        echo "❌ End-to-end tests failed"
        tar -czf e2e-test-artifacts.tar.gz $ARTIFACTS_DIR/e2e/
        echo "E2E test artifacts: e2e-test-artifacts.tar.gz"
        exit 1
    fi
fi

echo "All CI tests passed for build $BUILD_ID!"

# Generate test report
echo "=== CI Test Report ===" > $ARTIFACTS_DIR/test-report.txt
echo "Build ID: $BUILD_ID" >> $ARTIFACTS_DIR/test-report.txt
echo "Commit: $COMMIT_SHA" >> $ARTIFACTS_DIR/test-report.txt
echo "Date: $(date)" >> $ARTIFACTS_DIR/test-report.txt
echo "Status: PASSED" >> $ARTIFACTS_DIR/test-report.txt

# Archive artifacts
tar -czf ci-test-artifacts-$BUILD_ID.tar.gz $ARTIFACTS_DIR/
echo "Test artifacts archived: ci-test-artifacts-$BUILD_ID.tar.gz"
```

### Performance Testing Script
```bash
#!/bin/bash
# performance-kuttl-testing.sh

SCALE_LEVELS=(10 50 100 500)
PERFORMANCE_DIR="performance-results-$(date +%Y%m%d)"
mkdir -p $PERFORMANCE_DIR

echo "Starting KUTTL performance testing..."

for scale in "${SCALE_LEVELS[@]}"; do
    echo "Testing scale level: $scale resources"
    
    # Generate test configuration for this scale
    cat > configs/perf-test-$scale.yaml <<EOF
apiVersion: kuttl.dev/v1beta1
kind: TestSuite
metadata:
  name: performance-test-$scale
spec:
  startKIND: true
  testDirs:
  - tests/performance/
  parallel: 4
  timeout: 600
  env:
  - name: SCALE_LEVEL
    value: "$scale"
EOF
    
    # Run performance test
    echo "  Starting test run for scale $scale..."
    START_TIME=$(date +%s)
    
    if kubectl kuttl test \
        --config configs/perf-test-$scale.yaml \
        --artifacts-dir $PERFORMANCE_DIR/scale-$scale \
        --parallel 1 \
        -v \
        tests/performance/; then
        
        END_TIME=$(date +%s)
        DURATION=$((END_TIME - START_TIME))
        echo "  ✅ Scale $scale completed in ${DURATION}s"
        
        # Record results
        echo "scale_$scale:${DURATION}s:PASSED" >> $PERFORMANCE_DIR/results.txt
    else
        END_TIME=$(date +%s)
        DURATION=$((END_TIME - START_TIME))
        echo "  ❌ Scale $scale failed after ${DURATION}s"
        
        # Record failure
        echo "scale_$scale:${DURATION}s:FAILED" >> $PERFORMANCE_DIR/results.txt
    fi
done

# Generate performance report
echo "=== KUTTL Performance Test Report ===" > $PERFORMANCE_DIR/report.txt
echo "Test Date: $(date)" >> $PERFORMANCE_DIR/report.txt
echo "" >> $PERFORMANCE_DIR/report.txt
echo "Results:" >> $PERFORMANCE_DIR/report.txt
cat $PERFORMANCE_DIR/results.txt >> $PERFORMANCE_DIR/report.txt

echo "Performance testing complete! Results in $PERFORMANCE_DIR/"
```

## Integration Patterns

### GitHub Actions Integration
```yaml
# .github/workflows/kuttl-tests.yml
name: KUTTL Tests
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  kuttl-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        test-suite: [smoke, integration, e2e]
        
    steps:
    - uses: actions/checkout@v2
    
    - name: Install kubectl
      uses: azure/setup-kubectl@v1
      
    - name: Install KUTTL
      run: |
        kubectl krew install kuttl
        
    - name: Run KUTTL Tests
      env:
        TEST_SUITE: ${{ matrix.test-suite }}
      run: |
        echo "Running $TEST_SUITE tests..."
        
        # Configure test based on type
        case $TEST_SUITE in
          smoke)
            PARALLEL=8
            TIMEOUT=300
            ;;
          integration)
            PARALLEL=4
            TIMEOUT=600
            ;;
          e2e)
            PARALLEL=2
            TIMEOUT=1200
            ;;
        esac
        
        # Run tests
        kubectl kuttl test \
          --start-kind \
          --config configs/$TEST_SUITE.yaml \
          --artifacts-dir artifacts/$TEST_SUITE \
          --parallel $PARALLEL \
          -v \
          tests/$TEST_SUITE/
          
    - name: Upload Test Artifacts
      if: always()
      uses: actions/upload-artifact@v2
      with:
        name: kuttl-${{ matrix.test-suite }}-artifacts
        path: artifacts/${{ matrix.test-suite }}/
        
    - name: Test Summary
      if: always()
      run: |
        echo "### KUTTL Test Results - ${{ matrix.test-suite }}" >> $GITHUB_STEP_SUMMARY
        if [[ -d artifacts/${{ matrix.test-suite }} ]]; then
          TOTAL_TESTS=$(find artifacts/${{ matrix.test-suite }} -name "*.yaml" | wc -l)
          FAILED_TESTS=$(find artifacts/${{ matrix.test-suite }} -name "*.log" -exec grep -l "FAIL" {} \; 2>/dev/null | wc -l)
          PASSED_TESTS=$((TOTAL_TESTS - FAILED_TESTS))
          
          echo "**Total Tests:** $TOTAL_TESTS" >> $GITHUB_STEP_SUMMARY
          echo "**Passed:** $PASSED_TESTS" >> $GITHUB_STEP_SUMMARY
          echo "**Failed:** $FAILED_TESTS" >> $GITHUB_STEP_SUMMARY
          echo "**Success Rate:** $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%" >> $GITHUB_STEP_SUMMARY
        fi
```

### Kubernetes CronJob for Regular Testing
```yaml
# kuttl-scheduled-tests.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: kuttl-nightly-tests
  namespace: testing
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: kuttl-runner
          containers:
          - name: kuttl
            image: kudobuilder/kuttl:latest
            command:
            - /bin/bash
            - -c
            - |
              # Clone test repository
              git clone https://github.com/company/operator-tests.git /tests
              
              # Run nightly test suite
              kubectl kuttl test \
                --config /tests/configs/nightly.yaml \
                --artifacts-dir /artifacts \
                --parallel 4 \
                -v \
                /tests/nightly/
              
              # Upload results (implement based on your storage)
              # tar -czf /artifacts/nightly-results-$(date +%Y%m%d).tar.gz /artifacts/
            
            volumeMounts:
            - name: artifacts
              mountPath: /artifacts
            - name: kubeconfig
              mountPath: /root/.kube
              
          volumes:
          - name: artifacts
            persistentVolumeClaim:
              claimName: test-artifacts
          - name: kubeconfig
            secret:
              secretName: kuttl-kubeconfig
              
          restartPolicy: OnFailure
  
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kuttl-runner
  namespace: testing

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kuttl-runner
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kuttl-runner
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kuttl-runner
subjects:
- kind: ServiceAccount
  name: kuttl-runner
  namespace: testing
```

### Terraform Integration
```hcl
# terraform/kuttl-testing.tf
resource "kubernetes_namespace" "testing" {
  metadata {
    name = "kuttl-testing"
  }
}

resource "kubernetes_config_map" "kuttl_config" {
  metadata {
    name      = "kuttl-test-config"
    namespace = kubernetes_namespace.testing.metadata[0].name
  }

  data = {
    "kuttl-test.yaml" = yamlencode({
      apiVersion = "kuttl.dev/v1beta1"
      kind       = "TestSuite"
      metadata = {
        name = "terraform-managed-tests"
      }
      spec = {
        startKIND = false
        testDirs  = ["tests/"]
        parallel  = 4
        timeout   = 600
      }
    })
  }
}

resource "kubernetes_job" "kuttl_test" {
  count = var.run_tests ? 1 : 0
  
  metadata {
    name      = "kuttl-test-${random_id.test_id.hex}"
    namespace = kubernetes_namespace.testing.metadata[0].name
  }

  spec {
    template {
      metadata {
        labels = {
          app = "kuttl-test"
        }
      }
      
      spec {
        container {
          name  = "kuttl"
          image = "kudobuilder/kuttl:latest"
          
          command = ["kubectl", "kuttl", "test"]
          args = [
            "--config", "/config/kuttl-test.yaml",
            "--artifacts-dir", "/artifacts",
            "-v",
            "/tests"
          ]
          
          volume_mount {
            name       = "config"
            mount_path = "/config"
          }
          
          volume_mount {
            name       = "tests"
            mount_path = "/tests"
          }
          
          volume_mount {
            name       = "artifacts"
            mount_path = "/artifacts"
          }
        }
        
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.kuttl_config.metadata[0].name
          }
        }
        
        volume {
          name = "tests"
          git_repo {
            repository = var.test_repository
            revision   = var.test_branch
          }
        }
        
        volume {
          name = "artifacts"
          empty_dir {}
        }
        
        restart_policy = "Never"
      }
    }
  }
}

resource "random_id" "test_id" {
  byte_length = 4
}

variable "run_tests" {
  description = "Whether to run KUTTL tests"
  type        = bool
  default     = false
}

variable "test_repository" {
  description = "Git repository containing KUTTL tests"
  type        = string
}

variable "test_branch" {
  description = "Git branch to test"
  type        = string
  default     = "main"
}
```

## Best Practices

### Test Organization
- **Structured Directories**: Organize tests by functionality or component
- **Descriptive Names**: Use clear, descriptive names for test cases
- **Test Isolation**: Ensure tests don't interfere with each other
- **Resource Cleanup**: Clean up resources after test completion

### Configuration Management
- **Environment-Specific Configs**: Use different configurations for different environments
- **Parameterized Tests**: Use environment variables for test customization
- **Timeout Settings**: Set appropriate timeouts for different test types
- **Parallel Execution**: Balance speed with resource usage

### CI/CD Integration
- **Fast Feedback**: Run quick smoke tests first, longer tests later
- **Artifact Collection**: Always collect test artifacts for debugging
- **Conditional Testing**: Run appropriate tests based on changes
- **Clear Reporting**: Provide clear test results and summaries

### Development Workflow
- **Local Testing**: Test locally before pushing changes
- **Incremental Testing**: Test changes incrementally during development
- **Test-Driven Development**: Write tests before implementing features
- **Regular Updates**: Keep tests updated with application changes

## Error Handling

### Common Issues
```bash
# kubectl plugin not found
kubectl krew install kuttl
# Solution: Install KUTTL via krew plugin manager

# Permission denied
kubectl auth can-i "*" "*" --as=system:serviceaccount:default:default
# Solution: Ensure proper RBAC permissions

# Test timeout
kubectl kuttl test --timeout 600 ./tests/
# Solution: Increase timeout for slow tests

# KIND cluster issues
kubectl kuttl test --start-kind --parallel 1 ./tests/
# Solution: Reduce parallelism or use existing cluster
```

### Troubleshooting
- **Test Artifacts**: Always check artifacts directory for detailed logs
- **Verbose Output**: Use -v or -vv flags for detailed debugging information
- **Resource Issues**: Monitor cluster resources during test execution
- **Networking**: Verify cluster networking is functioning correctly

KUTTL provides powerful declarative testing capabilities for Kubernetes environments, enabling comprehensive testing of operators, controllers, and complex Kubernetes applications through simple YAML-based test definitions.