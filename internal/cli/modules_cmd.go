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
	// Display available MCP tools instead of Docker modules
	mcpTools := []struct {
		Name        string
		Description string
		Type        string
	}{
		{"lint", "TFLint for syntax and best practices", "terraform"},
		{"checkov", "Checkov security scanning", "security"},
		{"trivy", "Trivy security scanning", "security"},
		{"cost", "OpenInfraQuote cost analysis", "cost"},
		{"docs", "terraform-docs documentation", "documentation"},
		{"diagram", "InfraMap diagram generation", "visualization"},
		{"all", "All tools combined", "meta"},
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tDESCRIPTION")
	fmt.Fprintln(w, "----\t----\t-----------")

	for _, tool := range mcpTools {
		fmt.Fprintf(w, "%s\t%s\t%s\n", tool.Name, tool.Type, tool.Description)
	}

	w.Flush()

	fmt.Printf("\nTotal: %d MCP tools available\n", len(mcpTools))

	return nil
}

func runModulesInfo(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	moduleName := args[0]

	// Create module manager
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
