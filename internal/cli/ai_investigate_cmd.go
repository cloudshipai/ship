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
	aiInvestigateCmd.Flags().String("log-level", "info", "Log level for Dagger engine")

	aiInvestigateCmd.MarkFlagRequired("prompt")
}

func runAIInvestigate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

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

	// Step 1: Generate investigation plan
	slog.Info("Creating investigation plan", "prompt", prompt)

	// Generate investigation plan using LLM if available
	var investigationSteps []modules.InvestigationStep
	if llmProvider != "" && model != "" {
		steps, err := llmModule.CreateInvestigationPlan(ctx, prompt, []string{provider})
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
		result, err := steampipeModule.RunQuery(ctx, provider, step.Query, credentials, "json")
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

	// Step 3: AI Analysis is removed in favor of direct query execution.
	// The raw results are now available for the user to analyze.

	slog.Info("Investigation finished.")

	fmt.Println("\n📝 Next Steps:")
	fmt.Println("- Run 'ship terraform-tools checkov-scan' for detailed security analysis")
	fmt.Println("- Use 'ship push' to analyze your infrastructure with Cloudship AI")

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
