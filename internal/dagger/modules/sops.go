package modules

import (
	"context"
	"strings"

	"dagger.io/dagger"
)

type SOPSModule struct {
	client *dagger.Client
}

func NewSOPSModule(client *dagger.Client) *SOPSModule {
	return &SOPSModule{
		client: client,
	}
}

// EncryptFile encrypts a file using SOPS
func (m *SOPSModule) EncryptFile(ctx context.Context, filePath string, opts ...SOPSOption) (*dagger.Container, error) {
	config := &SOPSConfig{
		SOPSVersion: "v3.9.0",
		Format:      "yaml",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("mozilla/sops:" + config.SOPSVersion).
		WithWorkdir("/workspace")

	// Mount input file
	if filePath != "" {
		container = container.WithMountedFile("/workspace/input", m.client.Host().File(filePath))
	}

	// Mount KMS/AWS credentials if provided
	if config.AWSProfile != "" {
		container = container.WithEnvVariable("AWS_PROFILE", config.AWSProfile)
	}

	// Mount GPG keys if provided
	if config.GPGKeyDir != "" {
		container = container.WithMountedDirectory("/root/.gnupg", m.client.Host().Directory(config.GPGKeyDir))
	}

	// Mount age keys if provided
	if config.AgeKeyFile != "" {
		container = container.WithMountedFile("/root/.config/sops/age/keys.txt", m.client.Host().File(config.AgeKeyFile))
	}

	args := []string{"--encrypt"}

	// Add KMS ARN if provided
	if config.KMSARN != "" {
		args = append(args, "--kms", config.KMSARN)
	}

	// Add GPG fingerprint if provided
	if config.GPGFingerprint != "" {
		args = append(args, "--pgp", config.GPGFingerprint)
	}

	// Add age public key if provided
	if config.AgePublicKey != "" {
		args = append(args, "--age", config.AgePublicKey)
	}

	// Add Azure Key Vault if provided
	if config.AzureKeyVault != "" {
		args = append(args, "--azure-kv", config.AzureKeyVault)
	}

	// Add GCP KMS if provided
	if config.GCPKMS != "" {
		args = append(args, "--gcp-kms", config.GCPKMS)
	}

	// Add in-place flag if specified
	if config.InPlace {
		args = append(args, "--in-place")
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add input file
	args = append(args, "/workspace/input")

	return container.WithExec(args), nil
}

// DecryptFile decrypts a SOPS-encrypted file
func (m *SOPSModule) DecryptFile(ctx context.Context, filePath string, opts ...SOPSOption) (*dagger.Container, error) {
	config := &SOPSConfig{
		SOPSVersion: "v3.9.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("mozilla/sops:" + config.SOPSVersion).
		WithWorkdir("/workspace")

	// Mount encrypted file
	if filePath != "" {
		container = container.WithMountedFile("/workspace/encrypted", m.client.Host().File(filePath))
	}

	// Mount KMS/AWS credentials if provided
	if config.AWSProfile != "" {
		container = container.WithEnvVariable("AWS_PROFILE", config.AWSProfile)
	}

	// Mount GPG keys if provided
	if config.GPGKeyDir != "" {
		container = container.WithMountedDirectory("/root/.gnupg", m.client.Host().Directory(config.GPGKeyDir))
	}

	// Mount age keys if provided
	if config.AgeKeyFile != "" {
		container = container.WithMountedFile("/root/.config/sops/age/keys.txt", m.client.Host().File(config.AgeKeyFile))
	}

	args := []string{"--decrypt"}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add input file
	args = append(args, "/workspace/encrypted")

	return container.WithExec(args), nil
}

// RotateKeys rotates encryption keys for SOPS files
func (m *SOPSModule) RotateKeys(ctx context.Context, filePath string, opts ...SOPSOption) (*dagger.Container, error) {
	config := &SOPSConfig{
		SOPSVersion: "v3.9.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("mozilla/sops:" + config.SOPSVersion).
		WithWorkdir("/workspace")

	// Mount encrypted file
	if filePath != "" {
		container = container.WithMountedFile("/workspace/encrypted", m.client.Host().File(filePath))
	}

	// Mount GPG keys if provided
	if config.GPGKeyDir != "" {
		container = container.WithMountedDirectory("/root/.gnupg", m.client.Host().Directory(config.GPGKeyDir))
	}

	// Mount age keys if provided
	if config.AgeKeyFile != "" {
		container = container.WithMountedFile("/root/.config/sops/age/keys.txt", m.client.Host().File(config.AgeKeyFile))
	}

	args := []string{"--rotate"}

	// Add new KMS ARN if provided
	if config.KMSARN != "" {
		args = append(args, "--add-kms", config.KMSARN)
	}

	// Add new GPG fingerprint if provided
	if config.GPGFingerprint != "" {
		args = append(args, "--add-pgp", config.GPGFingerprint)
	}

	// Add new age public key if provided
	if config.AgePublicKey != "" {
		args = append(args, "--add-age", config.AgePublicKey)
	}

	// Remove old keys if specified
	for _, rmKMS := range config.RemoveKMS {
		args = append(args, "--rm-kms", rmKMS)
	}

	for _, rmPGP := range config.RemovePGP {
		args = append(args, "--rm-pgp", rmPGP)
	}

	for _, rmAge := range config.RemoveAge {
		args = append(args, "--rm-age", rmAge)
	}

	// Add in-place flag
	if config.InPlace {
		args = append(args, "--in-place")
	}

	// Add input file
	args = append(args, "/workspace/encrypted")

	return container.WithExec(args), nil
}

// EditFile opens a SOPS file for editing
func (m *SOPSModule) EditFile(ctx context.Context, filePath string, opts ...SOPSOption) (*dagger.Container, error) {
	config := &SOPSConfig{
		SOPSVersion: "v3.9.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("mozilla/sops:" + config.SOPSVersion).
		WithWorkdir("/workspace")

	// Mount encrypted file
	if filePath != "" {
		container = container.WithMountedFile("/workspace/encrypted", m.client.Host().File(filePath))
	}

	// Mount GPG keys if provided
	if config.GPGKeyDir != "" {
		container = container.WithMountedDirectory("/root/.gnupg", m.client.Host().Directory(config.GPGKeyDir))
	}

	// Mount age keys if provided
	if config.AgeKeyFile != "" {
		container = container.WithMountedFile("/root/.config/sops/age/keys.txt", m.client.Host().File(config.AgeKeyFile))
	}

	// Note: Interactive editing is limited in containerized environments
	// This command will show the decrypted content
	args := []string{"--decrypt", "/workspace/encrypted"}

	return container.WithExec(args), nil
}

// GenerateConfig creates a SOPS configuration file
func (m *SOPSModule) GenerateConfig(ctx context.Context, opts ...SOPSOption) (*dagger.Container, error) {
	config := &SOPSConfig{}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "yq"}).
		WithWorkdir("/workspace")

	// Create a basic SOPS configuration
	sopsConfig := `# SOPS Configuration
creation_rules:
  # Default rule for YAML files
  - path_regex: \.yaml$
    encrypted_regex: '^(data|stringData)$'
    kms: ""
    age: ""
    pgp: ""
    
  # Rule for environment files
  - path_regex: \.env$
    kms: ""
    age: ""
    pgp: ""
    
  # Rule for JSON files
  - path_regex: \.json$
    kms: ""
    age: ""
    pgp: ""
`

	if config.KMSARN != "" {
		sopsConfig = strings.ReplaceAll(sopsConfig, `kms: ""`, `kms: "`+config.KMSARN+`"`)
	}

	if config.AgePublicKey != "" {
		sopsConfig = strings.ReplaceAll(sopsConfig, `age: ""`, `age: "`+config.AgePublicKey+`"`)
	}

	if config.GPGFingerprint != "" {
		sopsConfig = strings.ReplaceAll(sopsConfig, `pgp: ""`, `pgp: "`+config.GPGFingerprint+`"`)
	}

	container = container.
		WithNewFile("/workspace/.sops.yaml", sopsConfig).
		WithExec([]string{"cat", "/workspace/.sops.yaml"})

	return container, nil
}

// ValidateFile validates a SOPS-encrypted file structure
func (m *SOPSModule) ValidateFile(ctx context.Context, filePath string, opts ...SOPSOption) (*dagger.Container, error) {
	config := &SOPSConfig{
		SOPSVersion: "v3.9.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("mozilla/sops:" + config.SOPSVersion).
		WithWorkdir("/workspace")

	// Mount encrypted file
	if filePath != "" {
		container = container.WithMountedFile("/workspace/encrypted", m.client.Host().File(filePath))
	}

	// Validate by attempting to decrypt and checking structure
	args := []string{"--decrypt", "--output-type", "json", "/workspace/encrypted"}

	return container.WithExec(args), nil
}

type SOPSConfig struct {
	SOPSVersion    string
	Format         string
	Output         string
	KMSARN         string
	GPGFingerprint string
	AgePublicKey   string
	AgeKeyFile     string
	GPGKeyDir      string
	AzureKeyVault  string
	GCPKMS         string
	AWSProfile     string
	InPlace        bool
	RemoveKMS      []string
	RemovePGP      []string
	RemoveAge      []string
}

type SOPSOption func(*SOPSConfig)

func WithSOPSVersion(version string) SOPSOption {
	return func(c *SOPSConfig) {
		c.SOPSVersion = version
	}
}

func WithSOPSFormat(format string) SOPSOption {
	return func(c *SOPSConfig) {
		c.Format = format
	}
}

func WithSOPSOutput(output string) SOPSOption {
	return func(c *SOPSConfig) {
		c.Output = output
	}
}

func WithKMSARN(arn string) SOPSOption {
	return func(c *SOPSConfig) {
		c.KMSARN = arn
	}
}

func WithGPGFingerprint(fingerprint string) SOPSOption {
	return func(c *SOPSConfig) {
		c.GPGFingerprint = fingerprint
	}
}

func WithAgePublicKey(key string) SOPSOption {
	return func(c *SOPSConfig) {
		c.AgePublicKey = key
	}
}

func WithAgeKeyFile(file string) SOPSOption {
	return func(c *SOPSConfig) {
		c.AgeKeyFile = file
	}
}

func WithGPGKeyDir(dir string) SOPSOption {
	return func(c *SOPSConfig) {
		c.GPGKeyDir = dir
	}
}

func WithAzureKeyVault(vault string) SOPSOption {
	return func(c *SOPSConfig) {
		c.AzureKeyVault = vault
	}
}

func WithGCPKMS(kms string) SOPSOption {
	return func(c *SOPSConfig) {
		c.GCPKMS = kms
	}
}

func WithAWSProfile(profile string) SOPSOption {
	return func(c *SOPSConfig) {
		c.AWSProfile = profile
	}
}

func WithInPlace(inPlace bool) SOPSOption {
	return func(c *SOPSConfig) {
		c.InPlace = inPlace
	}
}

func WithRemoveKMS(arns []string) SOPSOption {
	return func(c *SOPSConfig) {
		c.RemoveKMS = arns
	}
}

func WithRemovePGP(fingerprints []string) SOPSOption {
	return func(c *SOPSConfig) {
		c.RemovePGP = fingerprints
	}
}

func WithRemoveAge(keys []string) SOPSOption {
	return func(c *SOPSConfig) {
		c.RemoveAge = keys
	}
}
