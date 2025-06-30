package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudshipai/ship/internal/cloudship"
	"github.com/cloudshipai/ship/internal/config"
	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var terraformToolsCmd = &cobra.Command{
	Use:   "terraform-tools",
	Short: "Run Terraform analysis tools",
	Long:  `Run various Terraform analysis tools including cost estimation, security scanning, and documentation generation`,
}

var costAnalysisCmd = &cobra.Command{
	Use:   "cost-analysis [plan-file]",
	Short: "Analyze Terraform costs using OpenInfraQuote",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCostAnalysis,
}

var securityScanCmd = &cobra.Command{
	Use:   "security-scan [directory]",
	Short: "Scan Terraform code for security issues using InfraScan",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSecurityScan,
}

var generateDocsCmd = &cobra.Command{
	Use:   "generate-docs [directory]",
	Short: "Generate documentation for Terraform modules",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGenerateDocs,
}

var lintCmd = &cobra.Command{
	Use:   "lint [directory]",
	Short: "Lint Terraform code using TFLint",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runLint,
}

var checkovScanCmd = &cobra.Command{
	Use:   "checkov-scan [directory]",
	Short: "Scan Terraform code for security issues using Checkov",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCheckovScan,
}

var infracostCmd = &cobra.Command{
	Use:   "cost-estimate [directory]",
	Short: "Estimate cloud costs using Infracost",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInfracost,
}

var infraMapCmd = &cobra.Command{
	Use:   "generate-diagram [input]",
	Short: "Generate infrastructure diagram using InfraMap",
	Long: `Generate visual infrastructure diagrams from Terraform state files or HCL configurations.
InfraMap creates human-readable graphs showing infrastructure dependencies and relationships.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInfraMap,
}

func init() {
	rootCmd.AddCommand(terraformToolsCmd)
	terraformToolsCmd.AddCommand(costAnalysisCmd)
	terraformToolsCmd.AddCommand(securityScanCmd)
	terraformToolsCmd.AddCommand(generateDocsCmd)
	terraformToolsCmd.AddCommand(lintCmd)
	terraformToolsCmd.AddCommand(checkovScanCmd)
	terraformToolsCmd.AddCommand(infracostCmd)
	terraformToolsCmd.AddCommand(infraMapCmd)

	// Add common push flags to all terraform-tools subcommands
	for _, cmd := range []*cobra.Command{
		costAnalysisCmd,
		securityScanCmd,
		generateDocsCmd,
		lintCmd,
		checkovScanCmd,
		infracostCmd,
		infraMapCmd,
	} {
		cmd.Flags().Bool("push", false, "Automatically push results to CloudShip")
		cmd.Flags().String("push-fleet-id", "", "Fleet ID for push (overrides config/env)")
		cmd.Flags().StringSlice("push-tags", []string{}, "Tags for the pushed artifact")
		cmd.Flags().StringToString("push-metadata", map[string]string{}, "Additional metadata as key=value pairs")
	}

	// Add output file flags
	generateDocsCmd.Flags().StringP("output", "o", "", "Output file to save documentation (default: print to stdout)")
	generateDocsCmd.Flags().StringP("filename", "f", "README.md", "Filename to save documentation as")

	costAnalysisCmd.Flags().StringP("output", "o", "", "Output file to save cost analysis (default: print to stdout)")
	costAnalysisCmd.Flags().StringP("format", "", "json", "Output format: json, table")
	costAnalysisCmd.Flags().StringP("region", "r", "us-east-1", "AWS region for pricing (e.g., us-east-1, us-west-2)")

	securityScanCmd.Flags().StringP("output", "o", "", "Output file to save security scan results (default: print to stdout)")
	securityScanCmd.Flags().StringP("format", "", "json", "Output format: json, table, sarif")

	lintCmd.Flags().StringP("output", "o", "", "Output file to save lint results (default: print to stdout)")
	lintCmd.Flags().StringP("format", "", "default", "Output format: default, json, compact")

	checkovScanCmd.Flags().StringP("output", "o", "", "Output file to save scan results (default: print to stdout)")
	checkovScanCmd.Flags().StringP("format", "", "cli", "Output format: cli, json, junit, sarif")

	infracostCmd.Flags().StringP("output", "o", "", "Output file to save cost estimation (default: print to stdout)")
	infracostCmd.Flags().StringP("format", "", "table", "Output format: json, table, html")

	infraMapCmd.Flags().StringP("output", "o", "", "Output file to save diagram (default: print to stdout)")
	infraMapCmd.Flags().StringP("format", "", "png", "Output format: png, svg, pdf, dot")
	infraMapCmd.Flags().Bool("hcl", false, "Generate from HCL files instead of state file")
	infraMapCmd.Flags().Bool("raw", false, "Show all resources without InfraMap logic")
	infraMapCmd.Flags().Bool("no-clean", false, "Don't remove unconnected nodes")
	infraMapCmd.Flags().String("provider", "", "Filter by specific provider (aws, google, azurerm)")
}

// saveOrPrintOutput saves output to a file or prints to stdout
func saveOrPrintOutput(output, outputFile string, successMsg string) error {
	if outputFile != "" {
		// Ensure output directory exists
		if dir := filepath.Dir(outputFile); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		// Write to file
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		green := color.New(color.FgGreen)
		green.Printf("✓ %s\n", successMsg)
		fmt.Printf("Output saved to: %s\n", outputFile)
	} else {
		// Print to stdout
		green := color.New(color.FgGreen)
		green.Printf("\n✓ %s\n", successMsg)
		fmt.Printf("\n%s\n", output)
	}
	return nil
}

func runCostAnalysis(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Check if path exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create OpenInfraQuote module
	module := engine.NewOpenInfraQuoteModule()

	// Get region flag
	region, _ := cmd.Flags().GetString("region")

	// Check if it's a file or directory
	fileInfo, _ := os.Stat(path)
	var output string

	if fileInfo.IsDir() {
		fmt.Printf("Analyzing Terraform directory: %s\n", path)
		fmt.Printf("Using AWS region: %s\n", region)
		output, err = module.AnalyzeDirectory(ctx, path, region)
	} else {
		fmt.Printf("Analyzing Terraform plan file: %s\n", path)
		fmt.Printf("Using AWS region: %s\n", region)
		output, err = module.AnalyzePlan(ctx, path, region)
	}

	if err != nil {
		return fmt.Errorf("cost analysis failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "cost_analysis", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Cost analysis completed!")
}

func runSecurityScan(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create InfraScan module
	module := engine.NewInfraScanModule()

	fmt.Printf("Scanning Terraform code in: %s\n", dir)
	output, err := module.ScanDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("security scan failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "security_scan", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Security scan completed!")
}

func runGenerateDocs(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Get flags
	outputFile, _ := cmd.Flags().GetString("output")
	filename, _ := cmd.Flags().GetString("filename")

	// If output flag is not set but we have a filename, use it
	if outputFile == "" && filename != "README.md" {
		outputFile = filename
	} else if outputFile == "" && filename == "README.md" {
		// Auto-generate output file path
		outputFile = filepath.Join(dir, "README.md")
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create terraform-docs module
	module := engine.NewTerraformDocsModule()

	fmt.Printf("Generating documentation for: %s\n", dir)
	output, err := module.GenerateMarkdown(ctx, dir)
	if err != nil {
		return fmt.Errorf("documentation generation failed: %w", err)
	}

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "terraform_docs", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Documentation generated!")
}

func runLint(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create TFLint module
	module := engine.NewTFLintModule()

	fmt.Printf("Linting Terraform code in: %s\n", dir)
	output, err := module.LintDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("linting failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "lint_results", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Linting completed!")
}

func runCheckovScan(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create Checkov module
	module := engine.NewCheckovModule()

	fmt.Printf("Scanning Terraform code with Checkov in: %s\n", dir)
	output, err := module.ScanDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("Checkov scan failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "checkov_scan", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Checkov scan completed!")
}

func runInfracost(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create Infracost module
	module := engine.NewInfracostModule()

	// Check if INFRACOST_API_KEY is set
	if os.Getenv("INFRACOST_API_KEY") == "" {
		fmt.Println("Warning: INFRACOST_API_KEY is not set. Using free tier limits.")
	}

	fmt.Printf("Estimating costs for Terraform code in: %s\n", dir)
	output, err := module.GenerateTableReport(ctx, dir)
	if err != nil {
		return fmt.Errorf("cost estimation failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "infracost_estimate", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Cost estimation completed!")
}

func runInfraMap(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to terraform.tfstate or current directory
	input := "terraform.tfstate"
	if len(args) > 0 {
		input = args[0]
	}

	// Get flags
	format, _ := cmd.Flags().GetString("format")
	isHCL, _ := cmd.Flags().GetBool("hcl")
	raw, _ := cmd.Flags().GetBool("raw")
	noClean, _ := cmd.Flags().GetBool("no-clean")
	provider, _ := cmd.Flags().GetString("provider")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create InfraMap module
	module := engine.NewInfraMapModule()

	// Generate the diagram
	var output string

	if isHCL {
		// For HCL, input is a directory
		if _, err := os.Stat(input); err != nil {
			return fmt.Errorf("directory does not exist: %s", input)
		}
		fmt.Printf("Generating infrastructure diagram from HCL in: %s\n", input)
		output, err = module.GenerateFromHCL(ctx, input, format)
	} else if raw || noClean || provider != "" {
		// Use custom options
		options := modules.InfraMapOptions{
			Raw:      raw,
			Clean:    !noClean,
			Provider: provider,
			Format:   format,
		}
		fmt.Printf("Generating infrastructure diagram from: %s\n", input)
		output, err = module.GenerateWithOptions(ctx, input, options)
	} else {
		// Simple state file generation
		if _, err := os.Stat(input); err != nil {
			return fmt.Errorf("state file does not exist: %s", input)
		}
		fmt.Printf("Generating infrastructure diagram from state: %s\n", input)
		output, err = module.GenerateFromState(ctx, input, format)
	}

	if err != nil {
		return fmt.Errorf("diagram generation failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Handle push flag
	shouldPush, _ := cmd.Flags().GetBool("push")
	if shouldPush {
		if err := pushToCloudShip(cmd, output, "infrastructure_diagram", outputFile); err != nil {
			return fmt.Errorf("failed to push to CloudShip: %w", err)
		}
	}

	// For binary formats, we need to write as binary
	if format != "dot" && outputFile != "" {
		// Write binary output
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		color.Green("✓ Infrastructure diagram generated!")
		fmt.Printf("Output saved to: %s\n", outputFile)
		return nil
	}

	// For text formats or stdout
	return saveOrPrintOutput(output, outputFile, "Infrastructure diagram generated!")
}

// pushToCloudShip handles pushing results to CloudShip
func pushToCloudShip(cmd *cobra.Command, output string, scanType string, outputFile string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check for API key
	if cfg.APIKey == "" {
		return fmt.Errorf("not authenticated - run 'ship auth --api-key YOUR_KEY' first")
	}

	// Get fleet ID from flag, config, or env
	fleetID, _ := cmd.Flags().GetString("push-fleet-id")
	if fleetID == "" {
		fleetID = cfg.FleetID
	}
	if fleetID == "" {
		return fmt.Errorf("fleet ID required - use --push-fleet-id flag or set CLOUDSHIP_FLEET_ID")
	}

	// Get tags and metadata
	tags, _ := cmd.Flags().GetStringSlice("push-tags")
	metadata, _ := cmd.Flags().GetStringToString("push-metadata")

	// Create client
	client := cloudship.NewClient(cfg.APIKey)

	// Prepare metadata
	meta := map[string]interface{}{
		"scan_type":      scanType,
		"scan_timestamp": time.Now().UTC().Format(time.RFC3339),
		"source":         fmt.Sprintf("ship-cli/v%s", getVersion()),
		"tags":           tags,
		"command":        cmd.CommandPath(),
	}

	// Add custom metadata
	for k, v := range metadata {
		meta[k] = v
	}

	// Add ship ID if available
	if shipID := os.Getenv("SHIP_ID"); shipID != "" {
		meta["ship_id"] = shipID
	}

	// Add execution ID if available
	if execID := os.Getenv("EXECUTION_ID"); execID != "" {
		meta["execution_id"] = execID
	}

	// Determine filename
	fileName := outputFile
	if fileName == "" {
		fileName = fmt.Sprintf("%s-%s.json", scanType, time.Now().Format("20060102-150405"))
	} else {
		fileName = filepath.Base(fileName)
	}

	// Create upload request
	req := &cloudship.UploadArtifactRequest{
		FleetID:  fleetID,
		FileName: fileName,
		FileType: "application/json",
		Content:  base64.StdEncoding.EncodeToString([]byte(output)),
		Metadata: meta,
	}

	// Upload artifact
	fmt.Printf("\nPushing results to CloudShip...\n")
	
	resp, err := client.UploadArtifact(req)
	if err != nil {
		return err
	}

	// Success
	green := color.New(color.FgGreen)
	green.Printf("✓ Successfully pushed to CloudShip!\n")
	fmt.Printf("Artifact ID: %s\n", resp.ArtifactID)
	
	return nil
}

