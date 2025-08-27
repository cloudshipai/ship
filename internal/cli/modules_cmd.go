package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cloudshipai/ship/internal/cli/mcp"
	"github.com/cloudshipai/ship/internal/modules"
	"github.com/cloudshipai/ship/pkg/ship"
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
		Image       string
	}{
		// Terraform Tools
		{"lint", "TFLint for syntax and best practices", "terraform", "ghcr.io/terraform-linters/tflint:latest"},
		{"checkov", "Checkov security scanning", "terraform", "bridgecrew/checkov:latest"},
		{"trivy", "Trivy security scanning", "terraform", "aquasec/trivy:latest"},
		{"cost", "OpenInfraQuote cost analysis", "terraform", "ghcr.io/cycloidio/openinfraquote:latest"},
		{"docs", "terraform-docs documentation", "terraform", "quay.io/terraform-docs/terraform-docs:latest"},
		{"diagram", "InfraMap diagram generation", "terraform", "cycloid/inframap:latest"},
		{"terraform-docs", "Terraform documentation generation", "terraform", "quay.io/terraform-docs/terraform-docs:latest"},
		{"tflint", "Terraform linting", "terraform", "ghcr.io/terraform-linters/tflint:latest"},
		{"terrascan", "Infrastructure as Code security scanning", "terraform", "tenable/terrascan:latest"},
		{"openinfraquote", "Infrastructure cost estimation", "terraform", "ghcr.io/cycloidio/openinfraquote:latest"},
		{"aws-pricing-builtin", "Built-in AWS pricing information", "terraform", "N/A"},

		// Security Tools (Core)
		{"gitleaks", "Secret detection with Gitleaks", "security", "zricethezav/gitleaks:latest"},
		{"grype", "Vulnerability scanning with Grype", "security", "anchore/grype:latest"},
		{"syft", "SBOM generation with Syft", "security", "anchore/syft:latest"},
		{"prowler", "Multi-cloud security assessment", "security", "toniblyx/prowler:latest"},
		{"trufflehog", "Verified secret detection", "security", "trufflesecurity/trufflehog:latest"},
		{"cosign", "Container signing and verification", "security", "gcr.io/projectsigstore/cosign:latest"},

		// Security Tools (High-Priority Supply Chain)
		{"gatekeeper", "OPA Gatekeeper policy validation", "security", "openpolicyagent/gatekeeper:latest"},
		{"kubescape", "Kubernetes security scanning", "security", "quay.io/kubescape/kubescape:latest"},
		{"dockle", "Container image linting", "security", "goodwithtech/dockle:latest"},
		{"sops", "Secrets management with Mozilla SOPS", "security", "mozilla/sops:latest"},

		// Security Tools (Supply Chain)
		{"dependency-track", "OWASP Dependency-Track SBOM analysis", "security", "dependencytrack/bundler:latest"},

		// Security Tools (Additional)
		{"actionlint", "GitHub Actions workflow linting", "security", "rhymond/actionlint:latest"},
		{"semgrep", "Static analysis security scanning", "security", "returntocorp/semgrep:latest"},
		{"hadolint", "Dockerfile security linting", "security", "hadolint/hadolint:latest"},
		{"cfn-nag", "CloudFormation security scanning", "security", "stelligent/cfn_nag:latest"},
		{"conftest", "OPA policy testing", "security", "openpolicyagent/conftest:latest"},
		{"git-secrets", "Git secrets scanning", "security", "trufflesecurity/trufflehog:latest"},
		{"kube-bench", "Kubernetes security benchmarks", "security", "aquasec/kube-bench:latest"},
		{"kube-hunter", "Kubernetes penetration testing", "security", "aquasec/kube-hunter:latest"},
		{"zap", "Web application security testing", "security", "owasp/zap2docker-stable:latest"},
		{"falco", "Runtime security monitoring", "security", "falcosecurity/falco:latest"},
		{"nikto", "Web server security scanning", "security", "sullo/nikto:latest"},
		{"openscap", "Security compliance scanning", "security", "quay.io/compliance-operator/openscap:latest"},
		{"ossf-scorecard", "Open Source Security Foundation scorecard", "security", "gcr.io/openssf/scorecard:latest"},
		{"scout-suite", "Multi-cloud security auditing", "security", "nccgroup/scoutsuite:latest"},
		{"powerpipe", "Infrastructure benchmarking", "security", "turbot/powerpipe:latest"},
		{"velero", "Kubernetes backup and disaster recovery", "security", "velero/velero:latest"},
		{"goldilocks", "Kubernetes resource recommendations", "security", "us-docker.pkg.dev/fairwinds-ops/oss/goldilocks:latest"},
		{"osv-scanner", "Open Source Vulnerability scanning", "security", "gcr.io/osv-scanner/osv-scanner:latest"},
		{"license-detector", "Software license detection", "security", "licensefinder/license_finder:latest"},
		{"iac-plan", "Infrastructure as Code planning", "security", "hashicorp/terraform:latest"},

		// Cloud & Infrastructure Tools
		{"cloudquery", "Cloud asset inventory", "cloud", "cloudquery/cloudquery:latest"},
		{"custodian", "Cloud governance engine", "cloud", "cloudcustodian/c7n:latest"},
		{"terraformer", "Infrastructure import and management", "cloud", "quay.io/weaveworks/terraformer:latest"},
		{"inframap", "Infrastructure visualization", "cloud", "cycloidio/inframap:latest"},
		{"infrascan", "Infrastructure security scanning", "cloud", "bridgecrewio/checkov:latest"},
		{"aws-iam-rotation", "AWS IAM credential rotation", "cloud", "amazon/aws-cli:latest"},
		{"tfstate-reader", "Terraform state analysis", "cloud", "cycloidio/tfstate-lookup:latest"},
		{"packer", "Machine image building", "cloud", "hashicorp/packer:latest"},
		{"fleet", "GitOps for Kubernetes", "cloud", "rancher/fleet:latest"},
		{"kuttl", "Kubernetes testing framework", "cloud", "kudobuilder/kuttl:latest"},
		{"litmus", "Chaos engineering for Kubernetes", "cloud", "litmuschaos/litmus:latest"},
		{"cert-manager", "Certificate management", "cloud", "quay.io/jetstack/cert-manager-controller:latest"},
		{"k8s-network-policy", "Kubernetes network policy management", "cloud", "kinvolk/netfetch:latest"},
		{"kyverno", "Kubernetes policy management", "cloud", "ghcr.io/kyverno/kyverno:latest"},
		{"kyverno-multitenant", "Multi-tenant Kyverno policies", "cloud", "ghcr.io/kyverno/kyverno:latest"},
		{"github-admin", "GitHub administration tools", "cloud", "N/A (GitHub CLI)"},
		{"github-packages", "GitHub Packages security", "cloud", "N/A (GitHub CLI)"},

		// AWS IAM Tools
		{"cloudsplaining", "AWS IAM security assessment", "aws-iam", "bridgecrewio/cloudsplaining:latest"},
		{"parliament", "AWS IAM policy linting", "aws-iam", "duo-labs/parliament:latest"},
		{"pmapper", "AWS IAM privilege mapping", "aws-iam", "nccgroup/pmapper:latest"},
		{"policy-sentry", "AWS IAM policy generation", "aws-iam", "salesforce/policy-sentry:latest"},

		// Collections
		{"terraform", "All Terraform tools", "meta", "N/A (Collection)"},
		{"security", "All security tools", "meta", "N/A (Collection)"},
		{"aws-iam", "All AWS IAM tools", "meta", "N/A (Collection)"},
		{"cloud", "All cloud infrastructure tools", "meta", "N/A (Collection)"},
		{"kubernetes", "All Kubernetes tools", "meta", "N/A (Collection)"},
		{"all", "All tools combined", "meta", "N/A (Collection)"},

		// External MCP servers
		{"filesystem", "Filesystem operations MCP server", "mcp-external", "N/A (MCP Server)"},
		{"memory", "Memory/knowledge storage MCP server", "mcp-external", "N/A (MCP Server)"},
		{"brave-search", "Brave search MCP server", "mcp-external", "N/A (MCP Server)"},
		{"steampipe", "Cloud infrastructure queries MCP server", "mcp-external", "N/A (MCP Server)"},
		{"slack", "Slack workspace operations MCP server", "mcp-external", "N/A (MCP Server)"},
		{"github", "GitHub operations MCP server", "mcp-external", "N/A (MCP Server)"},
		{"desktop-commander", "Desktop operations MCP server", "mcp-external", "N/A (MCP Server)"},

		// AWS Labs Official MCP Servers
		{"aws-core", "AWS core operations and general services", "aws-mcp", "N/A (MCP Server)"},
		{"aws-iam", "AWS IAM operations and identity management", "aws-mcp", "N/A (MCP Server)"},
		{"aws-pricing", "AWS pricing and cost estimation", "aws-mcp", "N/A (MCP Server)"},
		{"aws-eks", "AWS EKS and Kubernetes operations", "aws-mcp", "N/A (MCP Server)"},
		{"aws-ec2", "AWS EC2 compute operations", "aws-mcp", "N/A (MCP Server)"},
		{"aws-s3", "AWS S3 storage operations", "aws-mcp", "N/A (MCP Server)"},
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tCONTAINER IMAGE\tDESCRIPTION")
	fmt.Fprintln(w, "----\t----\t---------------\t-----------")

	for _, tool := range mcpTools {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", tool.Name, tool.Type, tool.Image, tool.Description)
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
		dockerfile := `FROM alpine:latest

# Install dependencies
RUN apk add --no-cache bash curl

# Copy entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set working directory
WORKDIR /workspace

ENTRYPOINT ["/entrypoint.sh"]
`
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
	return mcp.IsExternalMCPServer(serverName)
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
	// Get the external server configuration from the proper mcp package
	config, exists := mcp.GetExternalMCPServer(serverName)
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

// getServerDescription returns a user-friendly description for the server
func getServerDescription(serverName string) string {
	// Check if this is an external MCP server
	if mcp.IsExternalMCPServer(serverName) {
		return fmt.Sprintf("External MCP server: %s", serverName)
	}
	return "External MCP server"
}

// hasOptionalVariables checks if any variables are optional
func hasOptionalVariables(variables []ship.Variable) bool {
	for _, variable := range variables {
		if !variable.Required {
			return true
		}
	}
	return false
}

// getOptionalVariableExample generates an example with optional variables
func getOptionalVariableExample(serverName string, _ []ship.Variable) string {
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
