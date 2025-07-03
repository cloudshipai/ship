package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/spf13/cobra"
)

var aiInvestigateCmd = &cobra.Command{
	Use:   "ai-investigate",
	Short: "Run AI-powered infrastructure investigation",
	Long: `Use natural language to investigate your cloud infrastructure with AI assistance.

Examples of complex investigations you can run:

Security & Compliance:
  - "Find all security groups allowing inbound traffic from 0.0.0.0/0"
  - "Show me IAM users without MFA enabled"
  - "List S3 buckets with public access or no encryption"
  - "Find RDS instances that are publicly accessible"
  - "Show EC2 instances with no tags or missing required tags"

Cost Optimization:
  - "Find unused EBS volumes and calculate their monthly cost"
  - "List EC2 instances that have been stopped for more than 30 days"
  - "Show me oversized instances with low CPU utilization"
  - "Find unattached Elastic IPs costing money"

Operations & Monitoring:
  - "List all Lambda functions with errors in the last 24 hours"
  - "Show EC2 instances without proper backup tags"
  - "Find load balancers with unhealthy targets"
  - "List VPCs with overlapping CIDR blocks"

Multi-Resource Investigations:
  - "Show me all resources in the us-west-2 region"
  - "Find resources created in the last 7 days"
  - "List all resources associated with a specific application tag"
  - "Show dependencies between EC2 instances and their security groups"`,
	RunE: runAIInvestigate,
}

func init() {
	rootCmd.AddCommand(aiInvestigateCmd)

	aiInvestigateCmd.Flags().String("prompt", "", "Natural language investigation prompt")
	aiInvestigateCmd.Flags().String("provider", "aws", "Cloud provider (aws, azure, gcp)")
	aiInvestigateCmd.Flags().String("llm-provider", "openai", "LLM provider (openai, anthropic, ollama)")
	aiInvestigateCmd.Flags().String("model", "gpt-4", "LLM model to use")
	aiInvestigateCmd.Flags().Bool("execute", false, "Execute the generated queries")
	aiInvestigateCmd.Flags().String("aws-profile", "", "AWS profile to use (from ~/.aws/config)")
	aiInvestigateCmd.Flags().String("aws-region", "", "AWS region to use")
	aiInvestigateCmd.Flags().String("log-level", "info", "Log level for Dagger engine")
	aiInvestigateCmd.Flags().Int("timeout", 15, "Timeout in minutes for the investigation")

	aiInvestigateCmd.MarkFlagRequired("prompt")
}

func runAIInvestigate(cmd *cobra.Command, args []string) error {
	// Get timeout from flags
	timeoutMinutes, _ := cmd.Flags().GetInt("timeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMinutes)*time.Minute)
	defer cancel()

	if err := checkDockerRunning(); err != nil {
		return err
	}

	prompt, _ := cmd.Flags().GetString("prompt")
	provider, _ := cmd.Flags().GetString("provider")
	llmProvider, _ := cmd.Flags().GetString("llm-provider")
	model, _ := cmd.Flags().GetString("model")
	execute, _ := cmd.Flags().GetBool("execute")
	awsProfile, _ := cmd.Flags().GetString("aws-profile")
	awsRegion, _ := cmd.Flags().GetString("aws-region")
	logLevel, _ := cmd.Flags().GetString("log-level")

	// Initialize Dagger engine
	slog.Info("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx, logLevel)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create modules
	// Use native Dagger LLM module if available, otherwise fall back to custom implementation
	var llmModule interface {
		CreateInvestigationPlan(ctx context.Context, objective string, providers []string) ([]modules.InvestigationStep, error)
		AnalyzeSteampipeResults(ctx context.Context, queryResults string, queryContext string) (string, error)
	}

	// Try native Dagger LLM first
	if llmProvider != "" {
		llmModule = modules.NewDaggerLLMModule(engine.GetClient(), model)
	} else {
		// Fall back to our custom LLM module
		llmModule = modules.NewLLMModule(engine.GetClient(), llmProvider, model)
	}

	steampipeModule := engine.NewSteampipeModule()

	// Step 0: Query table schemas for accurate column names
	slog.Info("Discovering table schemas...")
	
	// Prepare credentials for schema discovery
	credentials := getProviderCredentials(provider)
	if provider == "aws" {
		if awsProfile != "" {
			credentials["AWS_PROFILE"] = awsProfile
		}
		// Set AWS region, defaulting to us-east-1 if not specified
		if awsRegion != "" {
			credentials["AWS_REGION"] = awsRegion
		} else if credentials["AWS_REGION"] == "" {
			credentials["AWS_REGION"] = "us-east-1"
		}
	}
	
	// Get schema information for tables relevant to the prompt
	schemaInfo := ""
	relevantTables := make(map[string]bool)
	
	if provider == "aws" {
		// Always include some core tables
		relevantTables["aws_account"] = true
		
		// Add tables based on prompt keywords
		promptLower := strings.ToLower(prompt)
		
		// IAM-related
		if strings.Contains(promptLower, "iam") || strings.Contains(promptLower, "user") || 
		   strings.Contains(promptLower, "role") || strings.Contains(promptLower, "permission") ||
		   strings.Contains(promptLower, "security") || strings.Contains(promptLower, "access") {
			relevantTables["aws_iam_user"] = true
			relevantTables["aws_iam_role"] = true
			relevantTables["aws_iam_policy"] = true
			relevantTables["aws_iam_access_key"] = true
			relevantTables["aws_iam_group"] = true
		}
		
		// EC2-related
		if strings.Contains(promptLower, "ec2") || strings.Contains(promptLower, "instance") ||
		   strings.Contains(promptLower, "compute") || strings.Contains(promptLower, "server") {
			relevantTables["aws_ec2_instance"] = true
			relevantTables["aws_vpc_security_group"] = true
		}
		
		// S3-related
		if strings.Contains(promptLower, "s3") || strings.Contains(promptLower, "bucket") ||
		   strings.Contains(promptLower, "storage") {
			relevantTables["aws_s3_bucket"] = true
		}
		
		// RDS-related
		if strings.Contains(promptLower, "rds") || strings.Contains(promptLower, "database") ||
		   strings.Contains(promptLower, "db") {
			relevantTables["aws_rds_db_instance"] = true
			relevantTables["aws_rds_db_cluster"] = true
		}
		
		// VPC-related
		if strings.Contains(promptLower, "vpc") || strings.Contains(promptLower, "network") ||
		   strings.Contains(promptLower, "security group") {
			relevantTables["aws_vpc"] = true
			relevantTables["aws_vpc_security_group"] = true
			relevantTables["aws_vpc_subnet"] = true
		}
		
		// Lambda-related
		if strings.Contains(promptLower, "lambda") || strings.Contains(promptLower, "function") {
			relevantTables["aws_lambda_function"] = true
		}
		
		// Query actual column names from Steampipe
		for table := range relevantTables {
			columns, err := steampipeModule.GetTableColumns(ctx, provider, table, credentials)
			if err == nil && len(columns) > 0 {
				schemaInfo += fmt.Sprintf("\nTable %s has columns:\n", table)
				// Group columns for better readability
				for i, col := range columns {
					if i > 0 {
						schemaInfo += ", "
					}
					schemaInfo += col
				}
				schemaInfo += "\n"
			} else {
				slog.Debug("Could not get schema for table", "table", table, "error", err)
			}
		}
	}
	
	// Step 1: Generate investigation plan
	slog.Info("Creating investigation plan", "prompt", prompt)

	// Generate investigation plan using LLM if available
	var investigationSteps []modules.InvestigationStep
	if llmProvider != "" && model != "" {
		// Pass schema info with the prompt if available
		enhancedPrompt := prompt
		if schemaInfo != "" {
			slog.Info("Schema info discovered", "schema", schemaInfo)
			enhancedPrompt = fmt.Sprintf("%s\n\nAvailable table schemas:\n%s", prompt, schemaInfo)
		} else {
			slog.Warn("No schema info discovered")
		}
		steps, err := llmModule.CreateInvestigationPlan(ctx, enhancedPrompt, []string{provider})
		if err != nil {
			errorMsg := "AI failed to generate an investigation plan"
			if strings.Contains(err.Error(), "authentication") {
				errorMsg = "AI authentication failed. Please check your LLM provider API key and configuration."
			} else if strings.Contains(err.Error(), "model not found") {
				errorMsg = fmt.Sprintf("Model '%s' not found. You may not have access to it.", model)
				if model == "gpt-4" {
					errorMsg += " Try using '--model gpt-3.5-turbo' as an alternative."
				}
			}
			return fmt.Errorf("%s\n\nOriginal error: %w\n\n💡 You can also try rephrasing your prompt or use the 'ship query' command for direct execution.", errorMsg, err)
		}

		if len(steps) == 0 {
			return fmt.Errorf("AI generated an empty or invalid investigation plan. Please try rephrasing your prompt")
		}
		investigationSteps = steps
		slog.Info("Using AI-generated investigation plan")
	} else {
		return fmt.Errorf("LLM provider and model must be specified to generate an investigation plan. Use --llm-provider and --model flags, or use 'ship query' for direct execution.")
	}

	// Display the plan
	slog.Info("Investigation Plan:")
	for _, step := range investigationSteps {
		slog.Info("step", "number", step.StepNumber, "description", step.Description)
		slog.Info("  ", "provider", step.Provider)
		slog.Info("  ", "insights", step.ExpectedInsights)
		if !execute {
			slog.Info("  ", "query", truncateQuery(step.Query))
		} else {
			slog.Debug("Full query", "query", step.Query)
		}
	}

	if !execute {
		slog.Info("To execute this investigation, add the --execute flag")
		return nil
	}

	// Step 2: Execute queries
	slog.Info("Executing investigation...")
	
	// Fix common query mistakes before execution
	for i, step := range investigationSteps {
		fixedQuery := modules.FixSteampipeQuery(step.Query, provider)
		if fixedQuery != step.Query {
			slog.Info("Fixed query", "original", step.Query, "fixed", fixedQuery)
			investigationSteps[i].Query = fixedQuery
		}
	}

	// Use the credentials we already prepared during schema discovery

	allResults := make(map[string]interface{})

	for _, step := range investigationSteps {
		slog.Info("Executing step", "number", step.StepNumber, "description", step.Description)
		slog.Debug("Full query to execute", "query", step.Query)

		// Execute query with retry logic
		var result string
		var err error
		maxRetries := 3
		
		for attempt := 1; attempt <= maxRetries; attempt++ {
			result, err = steampipeModule.RunQuery(ctx, provider, step.Query, credentials, "json")
			if err == nil {
				break // Success!
			}
			
			// Check if we can fix the query
			errorMsg := err.Error()
			if attempt < maxRetries && strings.Contains(errorMsg, "column") && strings.Contains(errorMsg, "does not exist") {
				originalQuery := step.Query
				
				// Try to fix based on error
				if strings.Contains(errorMsg, `column "running"`) {
					step.Query = strings.ReplaceAll(step.Query, "WHERE running", "WHERE instance_state = 'running'")
				} else if strings.Contains(errorMsg, `"state_name"`) {
					step.Query = strings.ReplaceAll(step.Query, "state_name", "instance_state")
				} else if strings.Contains(errorMsg, `"state"`) && strings.Contains(step.Query, "aws_ec2") {
					step.Query = strings.ReplaceAll(step.Query, " state ", " instance_state ")
				} else if strings.Contains(errorMsg, "sg.group_id") {
					step.Query = strings.ReplaceAll(step.Query, "sg.group_id", "sg->>'GroupId'")
					step.Query = strings.ReplaceAll(step.Query, "sg.group_name", "sg->>'GroupName'")
				}
				
				if step.Query != originalQuery {
					slog.Info("Retrying with fixed query", "attempt", attempt+1, "fixed", step.Query)
					continue
				}
			}
			
			slog.Error("error executing query", "error", err, "attempt", attempt)
			if attempt >= maxRetries {
				continue // Move to next step
			}
		}

		// Parse results
		var queryResults []map[string]interface{}
		if err := json.Unmarshal([]byte(result), &queryResults); err == nil {
			allResults[fmt.Sprintf("step_%d", step.StepNumber)] = queryResults
			slog.Info("query finished", "found", len(queryResults))

			// Show sample results
			if len(queryResults) > 0 && len(queryResults[0]) > 0 {
				slog.Info("Sample finding:")
				for k, v := range queryResults[0] {
					slog.Info("  ", k, v)
					if len(fmt.Sprintf("   - %s: %v\n", k, v)) > 3 {
						break // Show only first 3 fields
					}
				}
			}
		}
	}

	// Step 3: Generate final summary
	slog.Info("Generating investigation summary...")
	
	// Create a summary of findings
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🔍 INVESTIGATION SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\nObjective: %s\n", prompt)
	fmt.Printf("Provider: %s (Region: %s)\n", provider, credentials["AWS_REGION"])
	
	// Summarize each step's results
	fmt.Println("\n📊 Findings:")
	for _, step := range investigationSteps {
		stepKey := fmt.Sprintf("step_%d", step.StepNumber)
		if results, ok := allResults[stepKey]; ok {
			fmt.Printf("\n%d. %s\n", step.StepNumber, step.Description)
			
			// Type assert to []map[string]interface{}
			if resultArray, ok := results.([]map[string]interface{}); ok {
				if len(resultArray) == 0 {
					fmt.Println("   ❌ No results found")
				} else if len(resultArray) == 1 && len(resultArray[0]) == 1 {
					// Single value result (like COUNT)
					for k, v := range resultArray[0] {
						fmt.Printf("   ✅ %s: %v\n", k, v)
					}
				} else {
					// Multiple results
					fmt.Printf("   ✅ Found %d items\n", len(resultArray))
					// Show first few items
					for i, item := range resultArray {
						if i >= 3 {
							fmt.Printf("   ... and %d more\n", len(resultArray)-3)
							break
						}
						fmt.Printf("   - Item %d:\n", i+1)
						for k, v := range item {
							fmt.Printf("     %s: %v\n", k, v)
						}
					}
				}
			}
		} else {
			fmt.Printf("\n%d. %s\n", step.StepNumber, step.Description)
			fmt.Println("   ⚠️  Query failed or no data returned")
		}
	}
	
	// Generate natural language summary
	fmt.Println("\n💡 Summary:")
	if totalCount, ok := allResults["step_1"].([]map[string]interface{}); ok && len(totalCount) > 0 {
		if count, ok := totalCount[0]["instance_count"].(float64); ok {
			if count == 0 {
				fmt.Println("You have no EC2 instances in your AWS account.")
			} else {
				fmt.Printf("You have %.0f EC2 instance(s) total.\n", count)
				
				// Check running instances
				if runningCount, ok := allResults["step_2"].([]map[string]interface{}); ok && len(runningCount) > 0 {
					if rCount, ok := runningCount[0]["count"].(float64); ok {
						if rCount == 0 {
							fmt.Println("None of your EC2 instances are currently running.")
						} else {
							fmt.Printf("%.0f of them are currently running.\n", rCount)
						}
					}
				}
			}
		}
	}
	
	// Use LLM for deeper analysis if available and results exist
	if llmProvider != "" && len(allResults) > 0 {
		fmt.Println("\n🤖 AI Analysis:")
		
		// Prepare results for AI analysis
		resultsJSON, _ := json.MarshalIndent(allResults, "", "  ")
		analysisPrompt := fmt.Sprintf("Based on this AWS infrastructure investigation for '%s', provide a brief analysis:\n\nResults:\n%s", prompt, string(resultsJSON))
		
		if analysis, err := llmModule.AnalyzeSteampipeResults(ctx, string(resultsJSON), analysisPrompt); err == nil {
			fmt.Println(analysis)
		} else {
			slog.Debug("AI analysis failed", "error", err)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("\n📝 Next Steps:")
	fmt.Println("- Run 'ship terraform-tools checkov-scan' for detailed security analysis")
	fmt.Println("- Use 'ship push' to analyze your infrastructure with Cloudship AI")
	fmt.Println("- Try more complex queries like:")
	fmt.Println("  - ship ai-investigate --prompt \"show me all security groups with open ports\" --execute")
	fmt.Println("  - ship ai-investigate --prompt \"find S3 buckets with public access\" --execute")
	fmt.Println("  - ship ai-investigate --prompt \"list Lambda functions and their costs\" --execute")

	return nil
}

func truncateQuery(query string) string {
	if len(query) > 80 {
		return query[:77] + "..."
	}
	return query
}

func getProviderCredentials(provider string) map[string]string {
	creds := make(map[string]string)

	switch provider {
	case "aws":
		// AWS credentials from environment
		if v := getEnvVar("AWS_ACCESS_KEY_ID"); v != "" {
			creds["AWS_ACCESS_KEY_ID"] = v
			slog.Info("Using AWS credentials from environment variables")
		}
		if v := getEnvVar("AWS_SECRET_ACCESS_KEY"); v != "" {
			creds["AWS_SECRET_ACCESS_KEY"] = v
		}
		if v := getEnvVar("AWS_SESSION_TOKEN"); v != "" {
			creds["AWS_SESSION_TOKEN"] = v
		}
		if v := getEnvVar("AWS_REGION"); v != "" {
			creds["AWS_REGION"] = v
		} else {
			creds["AWS_REGION"] = "us-east-1" // Default region
		}
		
		// If no environment credentials, try to load from ~/.aws/credentials
		if creds["AWS_ACCESS_KEY_ID"] == "" {
			if awsCreds := loadAWSCredentials(); awsCreds != nil {
				for k, v := range awsCreds {
					creds[k] = v
				}
				slog.Info("Using AWS credentials from ~/.aws/credentials")
			}
		}
	case "azure":
		// Azure credentials
		if v := getEnvVar("AZURE_TENANT_ID"); v != "" {
			creds["AZURE_TENANT_ID"] = v
		}
		if v := getEnvVar("AZURE_CLIENT_ID"); v != "" {
			creds["AZURE_CLIENT_ID"] = v
		}
		if v := getEnvVar("AZURE_CLIENT_SECRET"); v != "" {
			creds["AZURE_CLIENT_SECRET"] = v
		}
	case "gcp":
		// GCP credentials
		if v := getEnvVar("GOOGLE_APPLICATION_CREDENTIALS"); v != "" {
			creds["GOOGLE_APPLICATION_CREDENTIALS"] = v
		}
	}

	return creds
}

func getEnvVar(key string) string {
	return os.Getenv(key)
}

func loadAWSCredentials() map[string]string {
	creds := make(map[string]string)
	
	// Try to read from credentials file
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		credFile := filepath.Join(homeDir, ".aws", "credentials")
		if content, err := os.ReadFile(credFile); err == nil {
			// Simple parsing for default profile
			lines := strings.Split(string(content), "\n")
			inDefaultProfile := false
			slog.Debug("Parsing AWS credentials file")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "[default]" {
					inDefaultProfile = true
					continue
				}
				if strings.HasPrefix(line, "[") && line != "[default]" {
					inDefaultProfile = false
					continue
				}
				if inDefaultProfile {
					if strings.HasPrefix(line, "aws_access_key_id") {
						parts := strings.SplitN(line, "=", 2)
						if len(parts) == 2 {
							creds["AWS_ACCESS_KEY_ID"] = strings.TrimSpace(parts[1])
							slog.Debug("Found AWS access key in credentials file")
						}
					}
					if strings.HasPrefix(line, "aws_secret_access_key") {
						parts := strings.SplitN(line, "=", 2)
						if len(parts) == 2 {
							creds["AWS_SECRET_ACCESS_KEY"] = strings.TrimSpace(parts[1])
						}
					}
					if strings.HasPrefix(line, "region") {
						parts := strings.SplitN(line, "=", 2)
						if len(parts) == 2 {
							creds["AWS_REGION"] = strings.TrimSpace(parts[1])
						}
					}
				}
			}
		}
	}
	
	return creds
}
