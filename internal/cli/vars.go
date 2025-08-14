// Variables export implementation for Ship CLI

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var varsCmd = &cobra.Command{
	Use:   "vars [tool]",
	Short: "Export tool variables to variables.yml or print to console",
	Long: `Export tool variables to variables.yml file or print to console.

Available tools:
  lint       - TFLint variables
  checkov    - Checkov variables
  trivy      - Trivy variables
  cost       - Cost analysis variables
  docs       - Documentation variables
  diagram    - Diagram generation variables
  export     - Export all tool variables

Examples:
  ship vars lint                    # Export lint variables to ./variables.yml
  ship vars lint --print           # Print lint variables to console
  ship vars lint --prod            # Export to ~/.config/station/environments/prod/variables.yml
  ship vars export                 # Export all tools to variables.yml
  ship vars export --staging       # Export all to staging environment`,
	Args: cobra.ExactArgs(1),
	RunE: runVarsCommand,
}

func init() {
	rootCmd.AddCommand(varsCmd)

	varsCmd.Flags().Bool("print", false, "Print variables to console instead of file")
	varsCmd.Flags().Bool("dev", false, "Use dev environment (~/.config/station/environments/dev/variables.yml)")
	varsCmd.Flags().Bool("staging", false, "Use staging environment (~/.config/station/environments/staging/variables.yml)")
	varsCmd.Flags().Bool("prod", false, "Use prod environment (~/.config/station/environments/prod/variables.yml)")
}

func runVarsCommand(cmd *cobra.Command, args []string) error {
	toolName := args[0]
	printFlag, _ := cmd.Flags().GetBool("print")
	devFlag, _ := cmd.Flags().GetBool("dev")
	stagingFlag, _ := cmd.Flags().GetBool("staging")
	prodFlag, _ := cmd.Flags().GetBool("prod")

	// Validate tool name
	validTools := []string{"lint", "checkov", "trivy", "cost", "docs", "diagram", "export"}
	if !contains(validTools, toolName) {
		return fmt.Errorf("unknown tool: %s. Available: %s", toolName, strings.Join(validTools, ", "))
	}

	// Get variables for the tool(s)
	var variables map[string]interface{}
	var err error

	if toolName == "export" {
		variables, err = getAllToolVariables()
	} else {
		variables, err = getToolVariables(toolName)
	}
	if err != nil {
		return err
	}

	// Handle print flag
	if printFlag {
		return printVariables(variables)
	}

	// Determine output file path
	outputPath, err := getOutputPath(devFlag, stagingFlag, prodFlag)
	if err != nil {
		return err
	}

	// Write variables to file
	return writeVariablesToFile(variables, outputPath)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getToolVariables(toolName string) (map[string]interface{}, error) {
	switch toolName {
	case "lint":
		return getLintVariables(), nil
	case "checkov":
		return getCheckovVariables(), nil
	case "trivy":
		return getTrivyVariables(), nil
	case "cost":
		return getCostVariables(), nil
	case "docs":
		return getDocsVariables(), nil
	case "diagram":
		return getDiagramVariables(), nil
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

func getAllToolVariables() (map[string]interface{}, error) {
	allVars := make(map[string]interface{})

	tools := []string{"lint", "checkov", "trivy", "cost", "docs", "diagram"}
	for _, tool := range tools {
		vars, err := getToolVariables(tool)
		if err != nil {
			return nil, err
		}
		allVars[tool] = vars
	}

	return allVars, nil
}

func getLintVariables() map[string]interface{} {
	return map[string]interface{}{
		"lint": map[string]interface{}{
			"format":    "default",
			"output":    "",
			"directory": ".",
			"config":    ".tflint.hcl",
			"recursive": true,
			"module":    true,
		},
	}
}

func getCheckovVariables() map[string]interface{} {
	return map[string]interface{}{
		"checkov": map[string]interface{}{
			"format":      "cli",
			"output":      "",
			"directory":   ".",
			"config_file": ".checkov.yml",
			"framework":   "terraform",
			"check":       "",
			"skip_check":  "",
			"compact":     false,
			"quiet":       false,
		},
	}
}

func getTrivyVariables() map[string]interface{} {
	return map[string]interface{}{
		"trivy": map[string]interface{}{
			"directory":     ".",
			"scanners":      "misconfig",
			"format":        "json",
			"severity":      "HIGH,CRITICAL,MEDIUM,LOW",
			"config_policy": "",
			"policy":        "",
		},
	}
}

func getCostVariables() map[string]interface{} {
	return map[string]interface{}{
		"cost": map[string]interface{}{
			"directory":  ".",
			"region":     "us-east-1",
			"format":     "table",
			"output":     "",
			"currency":   "USD",
			"usage_file": "",
		},
	}
}

func getDocsVariables() map[string]interface{} {
	return map[string]interface{}{
		"docs": map[string]interface{}{
			"directory": ".",
			"filename":  "README.md",
			"output":    "",
			"config":    ".terraform-docs.yml",
			"sort_by":   "name",
			"recursive": false,
		},
	}
}

func getDiagramVariables() map[string]interface{} {
	return map[string]interface{}{
		"diagram": map[string]interface{}{
			"input":     ".",
			"format":    "png",
			"output":    "",
			"hcl":       true,
			"provider":  "",
			"show_plan": false,
		},
	}
}

func getOutputPath(devFlag, stagingFlag, prodFlag bool) (string, error) {
	// Check that only one environment flag is set
	envFlags := []bool{devFlag, stagingFlag, prodFlag}
	envCount := 0
	var envName string

	for i, flag := range envFlags {
		if flag {
			envCount++
			switch i {
			case 0:
				envName = "dev"
			case 1:
				envName = "staging"
			case 2:
				envName = "prod"
			}
		}
	}

	if envCount > 1 {
		return "", fmt.Errorf("only one environment flag can be specified")
	}

	// If no environment flag, use current directory
	if envCount == 0 {
		return "./variables.yml", nil
	}

	// Use XDG config path for station environments
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	stationPath := filepath.Join(homeDir, ".config", "station", "environments", envName, "variables.yml")

	// Ensure directory exists
	dir := filepath.Dir(stationPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return stationPath, nil
}

func printVariables(variables map[string]interface{}) error {
	yamlData, err := yaml.Marshal(variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables to YAML: %w", err)
	}

	fmt.Print(string(yamlData))
	return nil
}

func writeVariablesToFile(variables map[string]interface{}, outputPath string) error {
	// Check if file exists
	var existingVars map[string]interface{}
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, read and merge
		existingData, err := os.ReadFile(outputPath)
		if err != nil {
			return fmt.Errorf("failed to read existing file %s: %w", outputPath, err)
		}

		if err := yaml.Unmarshal(existingData, &existingVars); err != nil {
			return fmt.Errorf("failed to parse existing YAML: %w", err)
		}

		// Merge variables
		if existingVars == nil {
			existingVars = make(map[string]interface{})
		}

		for key, value := range variables {
			existingVars[key] = value
		}

		variables = existingVars
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
	}

	fmt.Printf("Variables exported to: %s\n", outputPath)
	return nil
}
