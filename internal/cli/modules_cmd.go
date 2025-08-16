package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cloudshipai/ship/internal/modules"
	"github.com/spf13/cobra"
)

var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Manage Ship CLI modules",
	Long:  `Discover, list, install, and manage Ship CLI modules for extending functionality`,
}

var modulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available modules",
	Long:  `List all discovered modules from built-in, user, project, and git sources`,
	RunE:  runModulesList,
}

var modulesInfoCmd = &cobra.Command{
	Use:   "info [module-name]",
	Short: "Show detailed information about a module",
	Long:  `Display detailed information about a specific module including commands, flags, and metadata`,
	Args:  cobra.ExactArgs(1),
	RunE:  runModulesInfo,
}

var modulesNewCmd = &cobra.Command{
	Use:   "new [module-name]",
	Short: "Create a new module template",
	Long:  `Create a new module template with the specified name and type`,
	Args:  cobra.ExactArgs(1),
	RunE:  runModulesNew,
}

func init() {
	rootCmd.AddCommand(modulesCmd)

	modulesCmd.AddCommand(modulesListCmd)
	modulesCmd.AddCommand(modulesInfoCmd)
	modulesCmd.AddCommand(modulesNewCmd)

	// Flags for list command
	modulesListCmd.Flags().StringP("type", "t", "", "Filter by module type (docker, dagger)")
	modulesListCmd.Flags().StringP("source", "s", "", "Filter by source (builtin, user, project, git)")
	modulesListCmd.Flags().BoolP("trusted", "", false, "Show only trusted modules")

	// Flags for new command
	modulesNewCmd.Flags().StringP("type", "t", "docker", "Module type (docker, dagger)")
	modulesNewCmd.Flags().StringP("description", "d", "", "Module description")
	modulesNewCmd.Flags().StringP("author", "a", "", "Module author")
	modulesNewCmd.Flags().StringP("output", "o", "", "Output directory (default: ~/.ship/modules/[module-name])")
}

func runModulesList(cmd *cobra.Command, args []string) error {
	// Display available MCP tools and external MCP servers
	mcpTools := []struct {
		Name        string
		Description string
		Type        string
	}{
		// Terraform Tools
		{"lint", "TFLint for syntax and best practices", "terraform"},
		{"checkov", "Checkov security scanning", "terraform"},
		{"trivy", "Trivy security scanning", "terraform"},
		{"cost", "OpenInfraQuote cost analysis", "terraform"},
		{"docs", "terraform-docs documentation", "terraform"},
		{"diagram", "InfraMap diagram generation", "terraform"},
		{"terraform-docs", "Terraform documentation generation", "terraform"},
		{"tflint", "Terraform linting", "terraform"},
		{"terrascan", "Infrastructure as Code security scanning", "terraform"},
		{"openinfraquote", "Infrastructure cost estimation", "terraform"},
		{"aws-pricing-builtin", "Built-in AWS pricing information", "terraform"},

		// Security Tools (Core)
		{"gitleaks", "Secret detection with Gitleaks", "security"},
		{"grype", "Vulnerability scanning with Grype", "security"},
		{"syft", "SBOM generation with Syft", "security"},
		{"prowler", "Multi-cloud security assessment", "security"},
		{"trufflehog", "Verified secret detection", "security"},
		{"cosign", "Container signing and verification", "security"},

		// Security Tools (High-Priority Supply Chain)
		{"slsa-verifier", "SLSA provenance verification", "security"},
		{"in-toto", "Supply chain attestation", "security"},
		{"gatekeeper", "OPA Gatekeeper policy validation", "security"},
		{"kubescape", "Kubernetes security scanning", "security"},
		{"dockle", "Container image linting", "security"},
		{"sops", "Secrets management with Mozilla SOPS", "security"},

		// Security Tools (Supply Chain)
		{"dependency-track", "OWASP Dependency-Track SBOM analysis", "security"},
		{"guac", "GUAC supply chain analysis", "security"},
		{"sigstore-policy-controller", "Sigstore Policy Controller", "security"},

		// Security Tools (Additional)
		{"actionlint", "GitHub Actions workflow linting", "security"},
		{"semgrep", "Static analysis security scanning", "security"},
		{"hadolint", "Dockerfile security linting", "security"},
		{"cfn-nag", "CloudFormation security scanning", "security"},
		{"conftest", "OPA policy testing", "security"},
		{"git-secrets", "Git secrets scanning", "security"},
		{"kube-bench", "Kubernetes security benchmarks", "security"},
		{"kube-hunter", "Kubernetes penetration testing", "security"},
		{"zap", "Web application security testing", "security"},
		{"falco", "Runtime security monitoring", "security"},
		{"nikto", "Web server security scanning", "security"},
		{"openscap", "Security compliance scanning", "security"},
		{"ossf-scorecard", "Open Source Security Foundation scorecard", "security"},
		{"scout-suite", "Multi-cloud security auditing", "security"},
		{"powerpipe", "Infrastructure benchmarking", "security"},
		{"velero", "Kubernetes backup and disaster recovery", "security"},
		{"goldilocks", "Kubernetes resource recommendations", "security"},
		{"allstar", "GitHub security policy enforcement", "security"},
		{"rekor", "Software supply chain transparency", "security"},
		{"osv-scanner", "Open Source Vulnerability scanning", "security"},
		{"license-detector", "Software license detection", "security"},
		{"registry", "Container registry operations", "security"},
		{"cosign-golden", "Enhanced Cosign for golden images", "security"},
		{"history-scrub", "Git history cleaning and secret removal", "security"},
		{"trivy-golden", "Enhanced Trivy for golden images", "security"},
		{"iac-plan", "Infrastructure as Code planning", "security"},

		// Cloud & Infrastructure Tools
		{"cloudquery", "Cloud asset inventory", "cloud"},
		{"custodian", "Cloud governance engine", "cloud"},
		{"terraformer", "Infrastructure import and management", "cloud"},
		{"infracost", "Infrastructure cost estimation", "cloud"},
		{"inframap", "Infrastructure visualization", "cloud"},
		{"infrascan", "Infrastructure security scanning", "cloud"},
		{"aws-iam-rotation", "AWS IAM credential rotation", "cloud"},
		{"tfstate-reader", "Terraform state analysis", "cloud"},
		{"packer", "Machine image building", "cloud"},
		{"fleet", "GitOps for Kubernetes", "cloud"},
		{"kuttl", "Kubernetes testing framework", "cloud"},
		{"litmus", "Chaos engineering for Kubernetes", "cloud"},
		{"cert-manager", "Certificate management", "cloud"},
		{"step-ca", "Certificate authority operations", "cloud"},
		{"check-ssl-cert", "SSL certificate validation", "cloud"},
		{"k8s-network-policy", "Kubernetes network policy management", "cloud"},
		{"kyverno", "Kubernetes policy management", "cloud"},
		{"kyverno-multitenant", "Multi-tenant Kyverno policies", "cloud"},
		{"github-admin", "GitHub administration tools", "cloud"},
		{"github-packages", "GitHub Packages security", "cloud"},

		// AWS IAM Tools
		{"cloudsplaining", "AWS IAM security assessment", "aws-iam"},
		{"parliament", "AWS IAM policy linting", "aws-iam"},
		{"pmapper", "AWS IAM privilege mapping", "aws-iam"},
		{"policy-sentry", "AWS IAM policy generation", "aws-iam"},

		// Collections
		{"terraform", "All Terraform tools", "meta"},
		{"security", "All security tools", "meta"},
		{"aws-iam", "All AWS IAM tools", "meta"},
		{"cloud", "All cloud infrastructure tools", "meta"},
		{"kubernetes", "All Kubernetes tools", "meta"},
		{"all", "All tools combined", "meta"},

		// External MCP servers
		{"filesystem", "Filesystem operations MCP server", "mcp-external"},
		{"memory", "Memory/knowledge storage MCP server", "mcp-external"},
		{"brave-search", "Brave search MCP server", "mcp-external"},
		{"steampipe", "Cloud infrastructure queries MCP server", "mcp-external"},

		// AWS Labs Official MCP Servers
		{"aws-core", "AWS core operations and general services", "aws-mcp"},
		{"aws-iam", "AWS IAM operations and identity management", "aws-mcp"},
		{"aws-pricing", "AWS pricing and cost estimation", "aws-mcp"},
		{"aws-eks", "AWS EKS and Kubernetes operations", "aws-mcp"},
		{"aws-ec2", "AWS EC2 compute operations", "aws-mcp"},
		{"aws-s3", "AWS S3 storage operations", "aws-mcp"},
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tDESCRIPTION")
	fmt.Fprintln(w, "----\t----\t-----------")

	for _, tool := range mcpTools {
		fmt.Fprintf(w, "%s\t%s\t%s\n", tool.Name, tool.Type, tool.Description)
	}

	w.Flush()

	// Count tools by type
	terraformTools := 0
	securityTools := 0
	cloudTools := 0
	awsIamTools := 0
	metaTools := 0
	externalServers := 0
	awsMcpServers := 0

	for _, tool := range mcpTools {
		switch tool.Type {
		case "terraform":
			terraformTools++
		case "security":
			securityTools++
		case "cloud":
			cloudTools++
		case "aws-iam":
			awsIamTools++
		case "meta":
			metaTools++
		case "mcp-external":
			externalServers++
		case "aws-mcp":
			awsMcpServers++
		}
	}

	totalBuiltinTools := terraformTools + securityTools + cloudTools + awsIamTools

	fmt.Printf("\nTotal: %d tools available\n", len(mcpTools))
	fmt.Printf("  - Built-in tools: %d (%d terraform, %d security, %d cloud, %d aws-iam)\n",
		totalBuiltinTools, terraformTools, securityTools, cloudTools, awsIamTools)
	fmt.Printf("  - Collections: %d\n", metaTools)
	fmt.Printf("  - External MCP servers: %d\n", externalServers)
	fmt.Printf("  - AWS Labs MCP servers: %d\n", awsMcpServers)

	return nil
}

func runModulesInfo(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	moduleName := args[0]

	// Check if this is an external MCP server first
	if isExternalMCPServerModule(moduleName) {
		return showMCPServerInfo(moduleName)
	}

	// Check if this is a built-in MCP tool
	if isBuiltinMCPTool(moduleName) {
		return showBuiltinToolInfo(moduleName)
	}

	// Create module manager for custom modules
	manager := modules.NewManager(modules.ModuleConfig{
		AllowUntrusted: true,
	})

	// Load modules
	if err := manager.LoadModules(ctx); err != nil {
		return fmt.Errorf("failed to load modules: %w", err)
	}

	// Find module
	module, err := manager.GetModule(moduleName)
	if err != nil {
		return err
	}

	// Display module information
	fmt.Printf("Module: %s\n", module.Metadata.Name)
	fmt.Printf("Version: %s\n", module.Metadata.Version)
	fmt.Printf("Description: %s\n", module.Metadata.Description)
	fmt.Printf("Author: %s\n", module.Metadata.Author)
	fmt.Printf("Type: %s\n", module.Spec.Type)
	fmt.Printf("Source: %s\n", module.Source)
	fmt.Printf("Trusted: %t\n", module.Trusted)

	if module.Path != "" {
		fmt.Printf("Path: %s\n", module.Path)
	}

	if !module.LoadedAt.IsZero() {
		fmt.Printf("Loaded: %s\n", module.LoadedAt.Format(time.RFC3339))
	}

	// Tags
	if len(module.Metadata.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(module.Metadata.Tags, ", "))
	}

	// Labels
	if len(module.Metadata.Labels) > 0 {
		fmt.Printf("Labels:\n")
		for k, v := range module.Metadata.Labels {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	// Commands
	if len(module.Spec.Commands) > 0 {
		fmt.Printf("\nCommands:\n")
		for _, command := range module.Spec.Commands {
			fmt.Printf("  %s - %s\n", command.Name, command.Description)

			// Command flags
			if len(command.Flags) > 0 {
				fmt.Printf("    Flags:\n")
				for _, flag := range command.Flags {
					required := ""
					if flag.Required {
						required = " (required)"
					}

					enumInfo := ""
					if len(flag.Enum) > 0 {
						enumInfo = fmt.Sprintf(" [%s]", strings.Join(flag.Enum, ","))
					}

					fmt.Printf("      --%s (%s)%s%s - %s\n",
						flag.Name,
						flag.Type,
						enumInfo,
						required,
						flag.Description,
					)
				}
			}

			// Command examples
			if len(command.Examples) > 0 {
				fmt.Printf("    Examples:\n")
				for _, example := range command.Examples {
					fmt.Printf("      %s\n", example)
				}
			}
		}
	}

	// Dependencies
	if len(module.Spec.Dependencies) > 0 {
		fmt.Printf("\nDependencies:\n")
		for _, dep := range module.Spec.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}
	}

	// Permissions
	if len(module.Spec.Permissions) > 0 {
		fmt.Printf("\nPermissions:\n")
		for _, perm := range module.Spec.Permissions {
			fmt.Printf("  - %s\n", perm)
		}
	}

	// Module-specific configuration
	if module.Spec.Docker != nil {
		fmt.Printf("\nDocker Configuration:\n")
		fmt.Printf("  Image: %s\n", module.Spec.Docker.Image)
		if len(module.Spec.Docker.Entrypoint) > 0 {
			fmt.Printf("  Entrypoint: %s\n", strings.Join(module.Spec.Docker.Entrypoint, " "))
		}
		if module.Spec.Docker.WorkingDir != "" {
			fmt.Printf("  Working Directory: %s\n", module.Spec.Docker.WorkingDir)
		}
	}

	if module.Spec.Dagger != nil {
		fmt.Printf("\nDagger Configuration:\n")
		fmt.Printf("  Module: %s\n", module.Spec.Dagger.Module)
		if module.Spec.Dagger.Function != "" {
			fmt.Printf("  Function: %s\n", module.Spec.Dagger.Function)
		}
	}

	return nil
}

func runModulesNew(cmd *cobra.Command, args []string) error {
	moduleName := args[0]

	// Get flags
	moduleType, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")
	author, _ := cmd.Flags().GetString("author")
	outputDir, _ := cmd.Flags().GetString("output")

	// Set defaults
	if description == "" {
		description = fmt.Sprintf("Custom %s module", moduleName)
	}
	if author == "" {
		author = "Unknown"
	}
	if outputDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		outputDir = fmt.Sprintf("%s/.ship/modules/%s", homeDir, moduleName)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate module template
	if err := generateModuleTemplate(outputDir, moduleName, moduleType, description, author); err != nil {
		return fmt.Errorf("failed to generate module template: %w", err)
	}

	fmt.Printf("Module '%s' created successfully at: %s\n", moduleName, outputDir)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Edit %s/module.yaml to customize the module\n", outputDir)

	if moduleType == "docker" {
		fmt.Printf("2. Edit %s/Dockerfile to define the container\n", outputDir)
		fmt.Printf("3. Add your tool's logic to %s/entrypoint.sh\n", outputDir)
	} else {
		fmt.Printf("2. Initialize Dagger module: cd %s && dagger init\n", outputDir)
		fmt.Printf("3. Implement your Dagger functions\n")
	}

	fmt.Printf("4. Test your module: ship modules info %s\n", moduleName)

	return nil
}

func generateModuleTemplate(outputDir, name, moduleType, description, author string) error {
	// Generate module.yaml
	moduleYaml := fmt.Sprintf(`apiVersion: ship.cloudship.ai/v1
kind: Module
metadata:
  name: %s
  version: "1.0.0"
  description: "%s"
  author: "%s"
  tags:
    - custom
    - %s

spec:
  type: %s
`, name, description, author, moduleType, moduleType)

	if moduleType == "docker" {
		moduleYaml += fmt.Sprintf(`  docker:
    image: "%s:latest"
    entrypoint: ["./entrypoint.sh"]
    
  commands:
    - name: "run"
      description: "Run %s"
      flags:
        - name: "input"
          type: "string"
          required: true
          description: "Input file or directory"
        - name: "output" 
          type: "string"
          description: "Output file"
`, name, name)
	} else {
		moduleYaml += fmt.Sprintf(`  dagger:
    module: "."
    function: "main"
    
  commands:
    - name: "run"
      description: "Run %s Dagger function"
      flags:
        - name: "source"
          type: "string"
          description: "Source directory"
`, name)
	}

	moduleYaml += `
  dependencies:
    - "docker"
    
  permissions:
    - "filesystem:read"
    - "network"
`

	if err := os.WriteFile(fmt.Sprintf("%s/module.yaml", outputDir), []byte(moduleYaml), 0644); err != nil {
		return err
	}

	// Generate additional files based on type
	if moduleType == "docker" {
		// Generate Dockerfile
		dockerfile := fmt.Sprintf(`FROM alpine:latest

# Install dependencies
RUN apk add --no-cache bash curl

# Copy entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set working directory
WORKDIR /workspace

ENTRYPOINT ["/entrypoint.sh"]
`)
		if err := os.WriteFile(fmt.Sprintf("%s/Dockerfile", outputDir), []byte(dockerfile), 0644); err != nil {
			return err
		}

		// Generate entrypoint.sh
		entrypoint := fmt.Sprintf(`#!/bin/bash
set -e

# %s implementation
# This is where you implement your tool's logic

echo "Running %s..."
echo "Arguments: $@"

# Example: process input file
if [ ! -z "$INPUT" ]; then
    echo "Processing input: $INPUT"
    # Add your processing logic here
fi

# Example: generate output
if [ ! -z "$OUTPUT" ]; then
    echo "Generated output" > "$OUTPUT"
    echo "Output written to: $OUTPUT"
fi

echo "%s completed successfully"
`, name, name, name)
		if err := os.WriteFile(fmt.Sprintf("%s/entrypoint.sh", outputDir), []byte(entrypoint), 0755); err != nil {
			return err
		}
	}

	// Generate README.md
	readme := fmt.Sprintf("# %s\n\n%s\n\n## Usage\n\n```bash\nship %s run --input ./source --output ./result\n```\n\n## Development\n\n### Building\n```bash\n", name, description, name)

	if moduleType == "docker" {
		readme += fmt.Sprintf("docker build -t %s:latest .\n", name)
	} else {
		readme += "dagger init\ndagger develop\n"
	}

	readme += "```\n\n### Testing\n```bash\nship modules info " + name + "\n```\n\n## Configuration\n\nEdit `module.yaml` to customize:\n- Commands and flags\n- Dependencies\n- Permissions\n- Description and metadata\n"

	if err := os.WriteFile(fmt.Sprintf("%s/README.md", outputDir), []byte(readme), 0644); err != nil {
		return err
	}

	return nil
}

// isExternalMCPServerModule checks if the given name is an external MCP server
func isExternalMCPServerModule(serverName string) bool {
	externalServers := []string{
		"filesystem", "memory", "brave-search", "steampipe",
		"aws-core", "aws-iam", "aws-pricing", "aws-eks", "aws-ec2", "aws-s3",
	}
	for _, server := range externalServers {
		if server == serverName {
			return true
		}
	}
	return false
}

// isBuiltinMCPTool checks if the given name is a built-in MCP tool
func isBuiltinMCPTool(toolName string) bool {
	builtinTools := []string{"lint", "checkov", "trivy", "cost", "docs", "diagram", "all"}
	for _, tool := range builtinTools {
		if tool == toolName {
			return true
		}
	}
	return false
}

// showMCPServerInfo displays information about an external MCP server
func showMCPServerInfo(serverName string) error {
	// Get the actual hardcoded server configuration from the mcp_cmd.go file
	// We need to access the hardcodedMCPServers map from mcp_cmd.go
	serverConfigs := getHardcodedMCPServerConfigs()

	config, exists := serverConfigs[serverName]
	if !exists {
		return fmt.Errorf("external MCP server '%s' not found", serverName)
	}

	fmt.Printf("External MCP Server: %s\n", config.Name)
	fmt.Printf("Type: mcp-external\n")
	fmt.Printf("Description: %s\n", getServerDescription(serverName))
	fmt.Printf("Transport: %s\n", config.Transport)
	fmt.Printf("Command: %s %s\n", config.Command, strings.Join(config.Args, " "))
	fmt.Printf("Source: hardcoded\n")
	fmt.Printf("Trusted: true\n")

	// Display variables if any are defined
	if len(config.Variables) > 0 {
		fmt.Printf("\nVariables:\n")
		for _, variable := range config.Variables {
			required := ""
			if variable.Required {
				required = " (required)"
			}

			secret := ""
			if variable.Secret {
				secret = " (secret)"
			}

			defaultInfo := ""
			if variable.Default != "" {
				if variable.Secret {
					defaultInfo = " [default: <hidden>]"
				} else {
					defaultInfo = fmt.Sprintf(" [default: %s]", variable.Default)
				}
			}

			fmt.Printf("  %s%s%s%s\n", variable.Name, defaultInfo, required, secret)
			fmt.Printf("    %s\n", variable.Description)
		}
	}

	fmt.Printf("\nUsage:\n")
	fmt.Printf("  ship mcp %s    # Start MCP server proxy\n", serverName)

	// Show usage examples with variables
	if len(config.Variables) > 0 {
		fmt.Printf("\nExamples with variables:\n")

		// Generate example command with required variables
		requiredVars := []string{}
		for _, variable := range config.Variables {
			if variable.Required {
				if variable.Secret {
					requiredVars = append(requiredVars, fmt.Sprintf("--var %s=<your_%s>", variable.Name, strings.ToLower(variable.Name)))
				} else {
					requiredVars = append(requiredVars, fmt.Sprintf("--var %s=<value>", variable.Name))
				}
			}
		}

		if len(requiredVars) > 0 {
			fmt.Printf("  ship mcp %s %s\n", serverName, strings.Join(requiredVars, " "))
		}

		// Show example with optional variables
		if hasOptionalVariables(config.Variables) {
			optionalExample := getOptionalVariableExample(serverName, config.Variables)
			if optionalExample != "" {
				fmt.Printf("  ship mcp %s %s\n", serverName, optionalExample)
			}
		}
	}

	fmt.Printf("\nNote: This is an external MCP server that Ship proxies.\n")
	fmt.Printf("Tools are discovered dynamically when the server starts.\n")

	return nil
}

// showBuiltinToolInfo displays information about a built-in MCP tool
func showBuiltinToolInfo(toolName string) error {
	toolConfigs := map[string]struct {
		Name        string
		Description string
		Type        string
		Examples    []string
	}{
		"lint": {
			Name:        "lint",
			Description: "TFLint for Terraform syntax checking and best practices validation",
			Type:        "terraform",
			Examples:    []string{"ship tf lint", "ship tf lint --format json", "ship mcp lint"},
		},
		"checkov": {
			Name:        "checkov",
			Description: "Checkov security and compliance scanning for Terraform",
			Type:        "security",
			Examples:    []string{"ship tf checkov", "ship tf checkov --format sarif", "ship mcp checkov"},
		},
		"trivy": {
			Name:        "trivy",
			Description: "Trivy security scanning for Terraform configurations",
			Type:        "security",
			Examples:    []string{"ship tf trivy", "ship mcp trivy"},
		},
		"cost": {
			Name:        "cost",
			Description: "OpenInfraQuote cost analysis for Terraform infrastructure",
			Type:        "cost",
			Examples:    []string{"ship tf cost", "ship tf cost --region us-east-1", "ship mcp cost"},
		},
		"docs": {
			Name:        "docs",
			Description: "terraform-docs documentation generation for Terraform modules",
			Type:        "documentation",
			Examples:    []string{"ship tf docs", "ship tf docs --filename USAGE.md", "ship mcp docs"},
		},
		"diagram": {
			Name:        "diagram",
			Description: "InfraMap infrastructure diagram generation from Terraform",
			Type:        "visualization",
			Examples:    []string{"ship tf diagram", "ship tf diagram --format svg", "ship mcp diagram"},
		},
		"all": {
			Name:        "all",
			Description: "All built-in Ship tools combined in a single MCP server",
			Type:        "meta",
			Examples:    []string{"ship mcp all", "ship mcp"},
		},
	}

	config, exists := toolConfigs[toolName]
	if !exists {
		return fmt.Errorf("built-in tool '%s' not found", toolName)
	}

	fmt.Printf("Built-in Tool: %s\n", config.Name)
	fmt.Printf("Type: %s\n", config.Type)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("Source: built-in\n")
	fmt.Printf("Trusted: true\n")

	fmt.Printf("\nUsage:\n")
	for _, example := range config.Examples {
		fmt.Printf("  %s\n", example)
	}

	fmt.Printf("\nNote: This is a built-in Ship tool that runs in containerized environments via Dagger.\n")

	return nil
}

// getHardcodedMCPServerConfigs returns the hardcoded MCP server configurations
// This is a reference to the same data structure in mcp_cmd.go
func getHardcodedMCPServerConfigs() map[string]struct {
	Name        string
	Description string
	Command     string
	Args        []string
	Transport   string
	Variables   []struct {
		Name        string
		Description string
		Required    bool
		Default     string
		Secret      bool
	}
} {
	return map[string]struct {
		Name        string
		Description string
		Command     string
		Args        []string
		Transport   string
		Variables   []struct {
			Name        string
			Description string
			Required    bool
			Default     string
			Secret      bool
		}
	}{
		"filesystem": {
			Name:        "filesystem",
			Description: "Filesystem operations MCP server with tools for file and directory management",
			Command:     "npx",
			Args:        []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
			Transport:   "stdio",
			Variables: []struct {
				Name        string
				Description string
				Required    bool
				Default     string
				Secret      bool
			}{
				{
					Name:        "FILESYSTEM_ROOT",
					Description: "Root directory for filesystem operations (overrides /tmp default)",
					Required:    false,
					Default:     "/tmp",
					Secret:      false,
				},
			},
		},
		"memory": {
			Name:        "memory",
			Description: "Memory/knowledge storage MCP server for persistent data storage",
			Command:     "npx",
			Args:        []string{"-y", "@modelcontextprotocol/server-memory"},
			Transport:   "stdio",
			Variables: []struct {
				Name        string
				Description string
				Required    bool
				Default     string
				Secret      bool
			}{
				{
					Name:        "MEMORY_STORAGE_PATH",
					Description: "Path for persistent memory storage",
					Required:    false,
					Default:     "/tmp/mcp-memory",
					Secret:      false,
				},
				{
					Name:        "MEMORY_MAX_SIZE",
					Description: "Maximum memory storage size (e.g., 100MB)",
					Required:    false,
					Default:     "50MB",
					Secret:      false,
				},
			},
		},
		"brave-search": {
			Name:        "brave-search",
			Description: "Brave search MCP server for web search capabilities",
			Command:     "npx",
			Args:        []string{"-y", "@modelcontextprotocol/server-brave-search"},
			Transport:   "stdio",
			Variables: []struct {
				Name        string
				Description string
				Required    bool
				Default     string
				Secret      bool
			}{
				{
					Name:        "BRAVE_API_KEY",
					Description: "Brave Search API key for search functionality",
					Required:    true,
					Default:     "",
					Secret:      true,
				},
				{
					Name:        "BRAVE_SEARCH_COUNT",
					Description: "Number of search results to return (default: 10)",
					Required:    false,
					Default:     "10",
					Secret:      false,
				},
			},
		},
		"steampipe": {
			Name:        "steampipe",
			Description: "Cloud infrastructure queries MCP server with SQL-based tools for cloud resources",
			Command:     "npx",
			Args:        []string{"-y", "@turbot/steampipe-mcp"},
			Transport:   "stdio",
			Variables: []struct {
				Name        string
				Description string
				Required    bool
				Default     string
				Secret      bool
			}{
				{
					Name:        "STEAMPIPE_DATABASE_CONNECTIONS",
					Description: "Database connections configuration for Steampipe",
					Required:    false,
					Default:     "postgres://steampipe@localhost:9193/steampipe",
					Secret:      false,
				},
			},
		},
	}
}

// getServerDescription returns a user-friendly description for the server
func getServerDescription(serverName string) string {
	descriptions := map[string]string{
		"filesystem":   "Filesystem operations MCP server with tools for file and directory management",
		"memory":       "Memory/knowledge storage MCP server for persistent data storage",
		"brave-search": "Brave search MCP server for web search capabilities",
		"steampipe":    "Cloud infrastructure queries MCP server with SQL-based tools for cloud resources",
	}
	if desc, exists := descriptions[serverName]; exists {
		return desc
	}
	return "External MCP server"
}

// hasOptionalVariables checks if any variables are optional
func hasOptionalVariables(variables []struct {
	Name        string
	Description string
	Required    bool
	Default     string
	Secret      bool
}) bool {
	for _, variable := range variables {
		if !variable.Required {
			return true
		}
	}
	return false
}

// getOptionalVariableExample generates an example with optional variables
func getOptionalVariableExample(serverName string, _ []struct {
	Name        string
	Description string
	Required    bool
	Default     string
	Secret      bool
}) string {
	examples := map[string]string{
		"filesystem":   "--var FILESYSTEM_ROOT=/custom/path",
		"memory":       "--var MEMORY_STORAGE_PATH=/data --var MEMORY_MAX_SIZE=100MB",
		"brave-search": "--var BRAVE_SEARCH_COUNT=20",
		"aws-core":     "--var AWS_PROFILE=production --var AWS_REGION=us-west-2",
		"aws-iam":      "--var AWS_PROFILE=default --var AWS_REGION=eu-west-1",
		"aws-pricing":  "--var AWS_REGION=us-east-1",
		"aws-eks":      "--var AWS_PROFILE=k8s-admin --var AWS_REGION=us-west-2",
		"aws-ec2":      "--var AWS_PROFILE=compute --var AWS_REGION=eu-central-1",
		"aws-s3":       "--var AWS_PROFILE=storage --var AWS_REGION=ap-southeast-1",
	}
	if example, exists := examples[serverName]; exists {
		return example
	}
	return ""
}
