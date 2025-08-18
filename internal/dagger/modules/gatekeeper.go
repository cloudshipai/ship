package modules

import (
	"context"

	"dagger.io/dagger"
)

type GatekeeperModule struct {
	client *dagger.Client
}

func NewGatekeeperModule(client *dagger.Client) *GatekeeperModule {
	return &GatekeeperModule{
		client: client,
	}
}

// ValidateConstraints validates Kubernetes resources against OPA Gatekeeper constraints
func (m *GatekeeperModule) ValidateConstraints(ctx context.Context, resourcesDir string, opts ...GatekeeperOption) (*dagger.Container, error) {
	config := &GatekeeperConfig{
		GatekeeperVersion: "v3.17.1",
		RegoVersion:       "v0.67.1",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("openpolicyagent/opa:" + config.RegoVersion).
		WithWorkdir("/workspace")

	// Mount resources directory
	if resourcesDir != "" {
		container = container.WithMountedDirectory("/workspace/resources", m.client.Host().Directory(resourcesDir))
	}

	// Mount constraints if provided
	if config.ConstraintsDir != "" {
		container = container.WithMountedDirectory("/workspace/constraints", m.client.Host().Directory(config.ConstraintsDir))
	}

	// Mount constraint templates if provided
	if config.TemplatesDir != "" {
		container = container.WithMountedDirectory("/workspace/templates", m.client.Host().Directory(config.TemplatesDir))
	}

	args := []string{"eval"}

	if config.Format != "" {
		args = append(args, "--format", config.Format)
	} else {
		args = append(args, "--format", "pretty")
	}

	// Add data directories
	if config.ConstraintsDir != "" {
		args = append(args, "--data", "/workspace/constraints")
	}
	if config.TemplatesDir != "" {
		args = append(args, "--data", "/workspace/templates")
	}

	// Add input directory
	args = append(args, "--input", "/workspace/resources")

	// Add query
	if config.Query != "" {
		args = append(args, config.Query)
	} else {
		args = append(args, "data.gatekeeper.violations")
	}

	return container.WithExec(args), nil
}

// TestConstraints runs tests for Gatekeeper constraints
func (m *GatekeeperModule) TestConstraints(ctx context.Context, testsDir string, opts ...GatekeeperOption) (*dagger.Container, error) {
	config := &GatekeeperConfig{
		RegoVersion: "v0.67.1",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("openpolicyagent/opa:" + config.RegoVersion).
		WithWorkdir("/workspace")

	// Mount tests directory
	if testsDir != "" {
		container = container.WithMountedDirectory("/workspace/tests", m.client.Host().Directory(testsDir))
	}

	args := []string{"test", "/workspace/tests"}

	if config.Verbose {
		args = append(args, "--verbose")
	}

	if config.Coverage {
		args = append(args, "--coverage")
	}

	return container.WithExec(args), nil
}

// GenerateConstraintTemplate creates a new constraint template
func (m *GatekeeperModule) GenerateConstraintTemplate(ctx context.Context, templateName string, opts ...GatekeeperOption) (*dagger.Container, error) {
	config := &GatekeeperConfig{}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "yq"}).
		WithWorkdir("/workspace")

	templateYAML := `apiVersion: templates.gatekeeper.sh/v1beta1
kind: ConstraintTemplate
metadata:
  name: ` + templateName + `
spec:
  crd:
    spec:
      names:
        kind: ` + templateName + `
      validation:
        type: object
        properties:
          message:
            type: string
          labels:
            type: array
            items:
              type: string
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package ` + templateName + `

        violation[{"msg": msg}] {
          required := input.parameters.labels
          provided := input.review.object.metadata.labels
          missing := required[_]
          not provided[missing]
          msg := sprintf("You must provide labels: %v", [missing])
        }
`

	container = container.
		WithNewFile("/workspace/"+templateName+".yaml", templateYAML).
		WithExec([]string{"cat", "/workspace/" + templateName + ".yaml"})

	return container, nil
}

// SyncConstraints syncs Gatekeeper constraints with cluster state
func (m *GatekeeperModule) SyncConstraints(ctx context.Context, opts ...GatekeeperOption) (*dagger.Container, error) {
	config := &GatekeeperConfig{
		GatekeeperVersion: "v3.17.1",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithWorkdir("/workspace")

	// Mount kubeconfig if provided
	if config.KubeconfigPath != "" {
		container = container.WithMountedFile("/root/.kube/config", m.client.Host().File(config.KubeconfigPath))
	}

	args := []string{"get", "constrainttemplates,constraints"}

	if config.Namespace != "" {
		args = append(args, "--namespace", config.Namespace)
	} else {
		args = append(args, "--all-namespaces")
	}

	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	return container.WithExec(args), nil
}

// AnalyzeViolations analyzes constraint violations in the cluster
func (m *GatekeeperModule) AnalyzeViolations(ctx context.Context, opts ...GatekeeperOption) (*dagger.Container, error) {
	config := &GatekeeperConfig{}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithWorkdir("/workspace")

	// Mount kubeconfig if provided
	if config.KubeconfigPath != "" {
		container = container.WithMountedFile("/root/.kube/config", m.client.Host().File(config.KubeconfigPath))
	}

	// Create analysis script
	analysisScript := `#!/bin/bash
echo "=== Gatekeeper Constraint Violations Analysis ==="
echo

echo "1. Constraint Templates:"
kubectl get constrainttemplates --no-headers | wc -l | xargs echo "   Total:"

echo
echo "2. Active Constraints:"
kubectl get constraints --all-namespaces --no-headers | wc -l | xargs echo "   Total:"

echo
echo "3. Violations by Constraint:"
for constraint in $(kubectl get constraints --all-namespaces --no-headers -o custom-columns=":metadata.name"); do
  violations=$(kubectl get $constraint --all-namespaces -o jsonpath='{.status.violations[*].message}' 2>/dev/null | wc -w)
  if [ "$violations" -gt 0 ]; then
    echo "   $constraint: $violations violations"
  fi
done

echo
echo "4. Recent Audit Results:"
kubectl get events --all-namespaces --field-selector reason=ConstraintViolation --sort-by='.lastTimestamp' | tail -10
`

	container = container.
		WithNewFile("/workspace/analyze.sh", analysisScript).
		WithExec([]string{"chmod", "+x", "/workspace/analyze.sh"}).
		WithExec([]string{"/workspace/analyze.sh"})

	return container, nil
}

// GetVersion returns the version of Gatekeeper
func (m *GatekeeperModule) GetVersion(ctx context.Context) (*dagger.Container, error) {
	container := m.client.Container().
		From("openpolicyagent/gatekeeper:v3.17.1").
		WithExec([]string{"manager", "--version"})

	return container, nil
}

type GatekeeperConfig struct {
	GatekeeperVersion string
	RegoVersion       string
	ConstraintsDir    string
	TemplatesDir      string
	KubeconfigPath    string
	Namespace         string
	Format            string
	Output            string
	Query             string
	Verbose           bool
	Coverage          bool
}

type GatekeeperOption func(*GatekeeperConfig)

func WithGatekeeperVersion(version string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.GatekeeperVersion = version
	}
}

func WithRegoVersion(version string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.RegoVersion = version
	}
}

func WithConstraintsDir(dir string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.ConstraintsDir = dir
	}
}

func WithTemplatesDir(dir string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.TemplatesDir = dir
	}
}

func WithKubeconfigPath(path string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.KubeconfigPath = path
	}
}

func WithNamespace(namespace string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.Namespace = namespace
	}
}

func WithFormat(format string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.Format = format
	}
}

func WithOutput(output string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.Output = output
	}
}

func WithQuery(query string) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.Query = query
	}
}

func WithVerbose(verbose bool) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.Verbose = verbose
	}
}

func WithCoverage(coverage bool) GatekeeperOption {
	return func(c *GatekeeperConfig) {
		c.Coverage = coverage
	}
}
