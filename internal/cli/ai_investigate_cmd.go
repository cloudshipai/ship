package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/spf13/cobra"
)

var aiInvestigateCmd = &cobra.Command{
	Use:   "ai-investigate",
	Short: "Run AI-powered infrastructure investigation",
	Long:  `Use natural language to investigate your cloud infrastructure with AI assistance`,
	RunE:  runAIInvestigate,
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

	aiInvestigateCmd.MarkFlagRequired("prompt")
}

func runAIInvestigate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	prompt, _ := cmd.Flags().GetString("prompt")
	provider, _ := cmd.Flags().GetString("provider")
	llmProvider, _ := cmd.Flags().GetString("llm-provider")
	model, _ := cmd.Flags().GetString("model")
	execute, _ := cmd.Flags().GetBool("execute")
	awsProfile, _ := cmd.Flags().GetString("aws-profile")
	awsRegion, _ := cmd.Flags().GetString("aws-region")

	// Initialize Dagger engine
	slog.Info("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
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

	// Step 1: Generate investigation plan
	slog.Info("Creating investigation plan", "prompt", prompt)

	// Generate investigation plan using LLM if available
	var investigationSteps []modules.InvestigationStep
	if llmProvider != "" && model != "" {
		// Try to use LLM module for smarter plan generation
		steps, err := llmModule.CreateInvestigationPlan(ctx, prompt, []string{provider})
		if err == nil && len(steps) > 0 {
			investigationSteps = steps
			slog.Info("Using AI-generated investigation plan")
		} else {
			// Fallback to hardcoded logic
			investigationSteps = GenerateInvestigationPlan(ctx, prompt, provider)
			slog.Info("Using rule-based investigation plan")
		}
	} else {
		// Use hardcoded logic when no LLM configured
		investigationSteps = GenerateInvestigationPlan(ctx, prompt, provider)
	}

	// Display the plan
	slog.Info("Investigation Plan:")
	for _, step := range investigationSteps {
		slog.Info("step", "number", step.StepNumber, "description", step.Description)
		slog.Info("  ", "provider", step.Provider)
		slog.Info("  ", "insights", step.ExpectedInsights)
		if !execute {
			slog.Info("  ", "query", truncateQuery(step.Query))
		}
	}

	if !execute {
		slog.Info("To execute this investigation, add the --execute flag")
		return nil
	}

	// Step 2: Execute queries
	slog.Info("Executing investigation...")

	// Prepare credentials with profile and region
	credentials := getProviderCredentials(provider)
	if provider == "aws" {
		if awsProfile != "" {
			credentials["AWS_PROFILE"] = awsProfile
		}
		if awsRegion != "" {
			credentials["AWS_REGION"] = awsRegion
		}
	}

	allResults := make(map[string]interface{})

	for _, step := range investigationSteps {
		slog.Info("Executing step", "number", step.StepNumber, "description", step.Description)

		// Execute query
		result, err := steampipeModule.RunQuery(ctx, provider, step.Query, credentials)
		if err != nil {
			slog.Error("error executing query", "error", err)
			continue
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

	// Step 3: AI Analysis
	slog.Info("AI Analysis:")

	// Try to use LLM for deeper analysis if available
	var insights string
	if llmProvider != "" && model != "" && len(allResults) > 0 {
		// Convert results to JSON for LLM analysis
		resultsJSON, _ := json.Marshal(allResults)
		llmInsights, err := llmModule.AnalyzeSteampipeResults(ctx, string(resultsJSON), prompt)
		if err == nil && llmInsights != "" {
			insights = llmInsights
			slog.Info("Using AI-powered analysis")
		} else {
			// Fallback to rule-based analysis
			insights = ParseQueryResults(allResults, prompt)
			slog.Info("Using rule-based analysis")
		}
	} else {
		// Use hardcoded analysis when no LLM configured
		insights = ParseQueryResults(allResults, prompt)
	}

	if insights != "" {
		slog.Info("Summary of Findings:")
		slog.Info(insights)
	} else {
		slog.Info("Investigation completed. No significant issues found based on your query.")
	}

	// Provide contextual recommendations based on the prompt
	slog.Info("Recommendations:")
	if strings.Contains(strings.ToLower(prompt), "s3") {
		slog.Info("- Review S3 bucket policies and access controls")
		slog.Info("- Enable versioning and encryption for sensitive buckets")
		slog.Info("- Consider implementing S3 lifecycle policies")
	} else if strings.Contains(strings.ToLower(prompt), "security") {
		slog.Info("- Review and restrict security group rules")
		slog.Info("- Enable encryption for all data at rest")
		slog.Info("- Implement least privilege access policies")
	} else if strings.Contains(strings.ToLower(prompt), "cost") {
		slog.Info("- Remove unused resources to reduce costs")
		slog.Info("- Consider using reserved instances for long-running workloads")
		slog.Info("- Implement auto-scaling to optimize resource usage")
	}

	slog.Info("Next Steps:")
	slog.Info("- Run 'ship terraform-tools checkov-scan' for detailed security analysis")
	slog.Info("- Use 'ship push' to analyze your infrastructure with Cloudship AI")

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
