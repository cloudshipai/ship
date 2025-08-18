package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type KubescapeModule struct {
	client *dagger.Client
}

func NewKubescapeModule(client *dagger.Client) *KubescapeModule {
	return &KubescapeModule{
		client: client,
	}
}

// ScanCluster scans a Kubernetes cluster for security issues
func (m *KubescapeModule) ScanCluster(ctx context.Context, opts ...KubescapeOption) (*dagger.Container, error) {
	config := &KubescapeConfig{
		KubescapeVersion: "v3.0.15",
		Framework:        "nsa",
		Format:           "pretty-printer",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("quay.io/kubescape/kubescape:" + config.KubescapeVersion).
		WithWorkdir("/workspace")

	// Mount kubeconfig if provided
	if config.KubeconfigPath != "" {
		container = container.WithMountedFile("/root/.kube/config", m.client.Host().File(config.KubeconfigPath))
	}

	args := []string{"kubescape", "scan"}

	// Add framework
	if config.Framework != "" {
		args = append(args, "framework", config.Framework)
	}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add severity threshold
	if config.SeverityThreshold != "" {
		args = append(args, "--severity-threshold", config.SeverityThreshold)
	}

	// Add compliance threshold
	if config.ComplianceThreshold > 0 {
		args = append(args, "--compliance-threshold", fmt.Sprintf("%.2f", config.ComplianceThreshold))
	}

	// Add namespace filter
	if config.Namespace != "" {
		args = append(args, "--include-namespaces", config.Namespace)
	}

	// Add resource filter
	if len(config.IncludeResources) > 0 {
		for _, resource := range config.IncludeResources {
			args = append(args, "--include-resources", resource)
		}
	}

	// Exclude kube-system by default unless specified
	if !config.IncludeKubeSystem {
		args = append(args, "--exclude-namespaces", "kube-system")
	}

	// Enable verbose output
	if config.Verbose {
		args = append(args, "--verbose")
	}

	return container.WithExec(args), nil
}

// ScanManifests scans Kubernetes manifest files
func (m *KubescapeModule) ScanManifests(ctx context.Context, manifestsDir string, opts ...KubescapeOption) (*dagger.Container, error) {
	config := &KubescapeConfig{
		KubescapeVersion: "v3.0.15",
		Framework:        "nsa",
		Format:           "pretty-printer",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("quay.io/kubescape/kubescape:" + config.KubescapeVersion).
		WithWorkdir("/workspace")

	// Mount manifests directory
	if manifestsDir != "" {
		container = container.WithMountedDirectory("/workspace/manifests", m.client.Host().Directory(manifestsDir))
	}

	args := []string{"scan", "framework", config.Framework}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add severity threshold
	if config.SeverityThreshold != "" {
		args = append(args, "--severity-threshold", config.SeverityThreshold)
	}

	// Add compliance threshold
	if config.ComplianceThreshold > 0 {
		args = append(args, "--compliance-threshold", fmt.Sprintf("%.2f", config.ComplianceThreshold))
	}

	// Enable verbose output
	if config.Verbose {
		args = append(args, "--verbose")
	}

	// Add manifests directory
	args = append(args, "/workspace/manifests")

	return container.WithExec(args), nil
}

// ScanHelm scans Helm charts for security issues
func (m *KubescapeModule) ScanHelm(ctx context.Context, chartPath string, opts ...KubescapeOption) (*dagger.Container, error) {
	config := &KubescapeConfig{
		KubescapeVersion: "v3.0.15",
		Framework:        "nsa",
		Format:           "pretty-printer",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("quay.io/kubescape/kubescape:" + config.KubescapeVersion).
		WithWorkdir("/workspace")

	// Mount chart directory
	if chartPath != "" {
		container = container.WithMountedDirectory("/workspace/chart", m.client.Host().Directory(chartPath))
	}

	args := []string{"scan", "framework", config.Framework}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add severity threshold
	if config.SeverityThreshold != "" {
		args = append(args, "--severity-threshold", config.SeverityThreshold)
	}

	// Enable verbose output
	if config.Verbose {
		args = append(args, "--verbose")
	}

	// Add chart path
	args = append(args, "/workspace/chart")

	return container.WithExec(args), nil
}

// ScanRepository scans a Git repository for security issues
func (m *KubescapeModule) ScanRepository(ctx context.Context, repoPath string, opts ...KubescapeOption) (*dagger.Container, error) {
	config := &KubescapeConfig{
		KubescapeVersion: "v3.0.15",
		Framework:        "nsa",
		Format:           "pretty-printer",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("quay.io/kubescape/kubescape:" + config.KubescapeVersion).
		WithWorkdir("/workspace")

	// Mount repository directory
	if repoPath != "" {
		container = container.WithMountedDirectory("/workspace/repo", m.client.Host().Directory(repoPath))
	}

	args := []string{"scan", "framework", config.Framework}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add severity threshold
	if config.SeverityThreshold != "" {
		args = append(args, "--severity-threshold", config.SeverityThreshold)
	}

	// Enable verbose output
	if config.Verbose {
		args = append(args, "--verbose")
	}

	// Enable repository scanning
	args = append(args, "--enable-host-scan")

	// Add repository path
	args = append(args, "/workspace/repo")

	return container.WithExec(args), nil
}

// GenerateReport generates a comprehensive security report
func (m *KubescapeModule) GenerateReport(ctx context.Context, opts ...KubescapeOption) (*dagger.Container, error) {
	config := &KubescapeConfig{
		KubescapeVersion: "v3.0.15",
		Framework:        "allframeworks",
		Format:           "html",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("quay.io/kubescape/kubescape:" + config.KubescapeVersion).
		WithWorkdir("/workspace")

	// Mount kubeconfig if provided
	if config.KubeconfigPath != "" {
		container = container.WithMountedFile("/root/.kube/config", m.client.Host().File(config.KubeconfigPath))
	}

	args := []string{"scan", "framework", config.Framework}

	// Add format for comprehensive report
	args = append(args, "--format", config.Format)

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	} else {
		args = append(args, "--output", "/workspace/kubescape-report.html")
	}

	// Include all severity levels for comprehensive report
	args = append(args, "--severity-threshold", "low")

	// Enable verbose output
	args = append(args, "--verbose")

	return container.WithExec(args), nil
}

type KubescapeConfig struct {
	KubescapeVersion    string
	Framework           string
	Format              string
	Output              string
	SeverityThreshold   string
	ComplianceThreshold float64
	Namespace           string
	IncludeResources    []string
	IncludeKubeSystem   bool
	KubeconfigPath      string
	Verbose             bool
}

type KubescapeOption func(*KubescapeConfig)

func WithKubescapeVersion(version string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.KubescapeVersion = version
	}
}

func WithFramework(framework string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.Framework = framework
	}
}

func WithKubescapeFormat(format string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.Format = format
	}
}

func WithKubescapeOutput(output string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.Output = output
	}
}

func WithSeverityThreshold(threshold string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.SeverityThreshold = threshold
	}
}

func WithComplianceThreshold(threshold float64) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.ComplianceThreshold = threshold
	}
}

func WithKubescapeNamespace(namespace string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.Namespace = namespace
	}
}

func WithIncludeResources(resources []string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.IncludeResources = resources
	}
}

func WithIncludeKubeSystem(include bool) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.IncludeKubeSystem = include
	}
}

func WithKubescapeKubeconfig(path string) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.KubeconfigPath = path
	}
}

func WithKubescapeVerbose(verbose bool) KubescapeOption {
	return func(c *KubescapeConfig) {
		c.Verbose = verbose
	}
}
