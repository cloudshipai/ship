package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/cloudship/ship/internal/dagger/modules"
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
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create modules
	_ = modules.NewLLMModule(engine.GetClient(), llmProvider, model) // Will be used for actual LLM calls in production
	steampipeModule := engine.NewSteampipeModule()

	// Step 1: Generate investigation plan
	fmt.Printf("\n🤖 AI: Creating investigation plan for: %s\n", prompt)

	// Generate dynamic investigation plan based on the prompt
	investigationSteps := GenerateInvestigationPlan(ctx, prompt, provider)

	// Display the plan
	fmt.Println("\n📋 Investigation Plan:")
	for _, step := range investigationSteps {
		fmt.Printf("\n%d. %s\n", step.StepNumber, step.Description)
		fmt.Printf("   Provider: %s\n", step.Provider)
		fmt.Printf("   Expected insights: %s\n", step.ExpectedInsights)
		if !execute {
			fmt.Printf("   Query: %s\n", truncateQuery(step.Query))
		}
	}

	if !execute {
		fmt.Println("\n💡 To execute this investigation, add the --execute flag")
		return nil
	}

	// Step 2: Execute queries
	fmt.Println("\n🔍 Executing investigation...")

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
		fmt.Printf("\nStep %d: %s\n", step.StepNumber, step.Description)

		// Execute query
		result, err := steampipeModule.RunQuery(ctx, provider, step.Query, credentials)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		// Parse results
		var queryResults []map[string]interface{}
		if err := json.Unmarshal([]byte(result), &queryResults); err == nil {
			allResults[fmt.Sprintf("step_%d", step.StepNumber)] = queryResults
			fmt.Printf("✓ Found %d results\n", len(queryResults))

			// Show sample results
			if len(queryResults) > 0 && len(queryResults[0]) > 0 {
				fmt.Println("   Sample finding:")
				for k, v := range queryResults[0] {
					fmt.Printf("   - %s: %v\n", k, v)
					if len(fmt.Sprintf("   - %s: %v\n", k, v)) > 3 {
						break // Show only first 3 fields
					}
				}
			}
		}
	}

	// Step 3: AI Analysis
	fmt.Println("\n🧠 AI Analysis:")

	// Generate insights based on actual results
	insights := ParseQueryResults(allResults, prompt)

	if insights != "" {
		fmt.Println("\n📊 Summary of Findings:")
		fmt.Println(insights)
	} else {
		fmt.Println("\n✅ Investigation completed. No significant issues found based on your query.")
	}

	// Provide contextual recommendations based on the prompt
	fmt.Println("\n💡 Recommendations:")
	if strings.Contains(strings.ToLower(prompt), "s3") {
		fmt.Println("- Review S3 bucket policies and access controls")
		fmt.Println("- Enable versioning and encryption for sensitive buckets")
		fmt.Println("- Consider implementing S3 lifecycle policies")
	} else if strings.Contains(strings.ToLower(prompt), "security") {
		fmt.Println("- Review and restrict security group rules")
		fmt.Println("- Enable encryption for all data at rest")
		fmt.Println("- Implement least privilege access policies")
	} else if strings.Contains(strings.ToLower(prompt), "cost") {
		fmt.Println("- Remove unused resources to reduce costs")
		fmt.Println("- Consider using reserved instances for long-running workloads")
		fmt.Println("- Implement auto-scaling to optimize resource usage")
	}

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
			fmt.Println("✓ Using AWS credentials from environment variables")
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
