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

	"github.com/cloudshipai/ship/internal/agent"
	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/spf13/cobra"
)

var investigateCmd = &cobra.Command{
	Use:   "investigate",
	Short: "AI-powered infrastructure investigation using Eino framework",
	Long: `Use natural language to investigate your cloud infrastructure with Eino-powered AI assistance.
This new reliable version uses the Eino framework for better accuracy and consistency.

Examples of investigations you can run:

Security & Compliance:
  - "Find all security groups allowing inbound traffic from 0.0.0.0/0"
  - "Show me IAM users without MFA enabled"
  - "List S3 buckets with public access or no encryption"
  - "Find RDS instances that are publicly accessible"

Cost Optimization:
  - "Find unused EBS volumes and calculate their monthly cost"
  - "List EC2 instances that have been stopped for more than 30 days"
  - "Show me oversized instances with low CPU utilization"

Operations & Monitoring:
  - "List all Lambda functions with errors in the last 24 hours"
  - "Show EC2 instances without proper backup tags"
  - "Find load balancers with unhealthy targets"`,
	RunE: runInvestigate,
}

func init() {
	rootCmd.AddCommand(investigateCmd)

	investigateCmd.Flags().String("prompt", "", "Natural language investigation prompt")
	investigateCmd.Flags().String("provider", "aws", "Cloud provider (aws, azure, gcp)")
	investigateCmd.Flags().String("region", "", "Cloud region to focus on")
	investigateCmd.Flags().String("openai-key", "", "OpenAI API key (or use OPENAI_API_KEY env var)")
	investigateCmd.Flags().String("memory-path", "", "Path to save agent memory (optional)")
	investigateCmd.Flags().String("log-level", "info", "Log level for Dagger engine")
	investigateCmd.Flags().Int("timeout", 10, "Timeout in minutes for the investigation")
	investigateCmd.Flags().Bool("save-results", false, "Save investigation results to file")

	investigateCmd.MarkFlagRequired("prompt")
}

func runInvestigate(cmd *cobra.Command, args []string) error {
	// Get timeout from flags
	timeoutMinutes, _ := cmd.Flags().GetInt("timeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMinutes)*time.Minute)
	defer cancel()

	if err := checkDockerRunning(); err != nil {
		return err
	}

	// Parse flags
	prompt, _ := cmd.Flags().GetString("prompt")
	provider, _ := cmd.Flags().GetString("provider")
	region, _ := cmd.Flags().GetString("region")
	openaiKey, _ := cmd.Flags().GetString("openai-key")
	memoryPath, _ := cmd.Flags().GetString("memory-path")
	logLevel, _ := cmd.Flags().GetString("log-level")
	saveResults, _ := cmd.Flags().GetBool("save-results")

	// Get OpenAI API key from environment if not provided
	if openaiKey == "" {
		openaiKey = os.Getenv("OPENAI_API_KEY")
		if openaiKey == "" {
			return fmt.Errorf("OpenAI API key is required. Use --openai-key flag or set OPENAI_API_KEY environment variable")
		}
	}

	// Set default memory path if not provided
	if memoryPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			memoryPath = filepath.Join(homeDir, ".ship", "agent_memory.json")
		}
	}

	// Initialize Dagger engine
	slog.Info("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx, logLevel)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create Eino investigation agent
	slog.Info("Creating Eino investigation agent...")
	investigationAgent, err := agent.NewEinoInvestigationAgent(ctx, engine.GetClient(), openaiKey, memoryPath)
	if err != nil {
		return fmt.Errorf("failed to create Eino agent: %w", err)
	}

	// Prepare credentials
	credentials := getProviderCredentials(provider)
	
	// Set region if provided
	if region != "" {
		credentials["AWS_REGION"] = region
	} else if credentials["AWS_REGION"] == "" {
		credentials["AWS_REGION"] = "us-east-1" // Default region
	}

	// Create investigation request
	var request agent.InvestigationRequest
	request.Prompt = prompt
	request.Provider = provider
	request.Region = region
	request.Credentials = credentials

	// Execute investigation
	slog.Info("Starting Eino-powered investigation", "prompt", prompt, "provider", provider)
	
	result, err := investigationAgent.Investigate(ctx, request)
	if err != nil {
		return fmt.Errorf("investigation failed: %w", err)
	}

	// Display results
	displayInvestigationResult(result, prompt, provider, credentials["AWS_REGION"])

	// Save results if requested
	if saveResults {
		if err := saveInvestigationResults(result, prompt); err != nil {
			slog.Warn("Failed to save investigation results", "error", err)
		}
	}

	return nil
}

// displayInvestigationResult shows the investigation results in a user-friendly format
func displayInvestigationResult(result *agent.InvestigationResult, prompt, provider, region string) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🤖 EINO AGENT INVESTIGATION RESULTS")
	fmt.Println(strings.Repeat("=", 70))
	
	fmt.Printf("\n📋 Investigation: %s\n", prompt)
	fmt.Printf("🔍 Provider: %s", provider)
	if region != "" {
		fmt.Printf(" (Region: %s)", region)
	}
	fmt.Println()
	fmt.Printf("⏱️ Duration: %s\n", result.Duration)
	fmt.Printf("📊 Confidence: %.1f%%\n", result.Confidence*100)
	fmt.Printf("🔢 Queries executed: %d\n", result.QueryCount)

	// Show investigation steps if available
	if len(result.Steps) > 0 {
		fmt.Println("\n📋 Investigation Steps:")
		for _, step := range result.Steps {
			status := "✅"
			if !step.Success {
				status = "❌"
			}
			fmt.Printf("  %s Step %d: %s\n", status, step.StepNumber, step.Description)
			if step.Error != "" {
				fmt.Printf("    Error: %s\n", step.Error)
			} else if len(step.Results) > 0 {
				fmt.Printf("    Found: %d results\n", len(step.Results))
			}
		}
	}

	// Show insights
	if len(result.Insights) > 0 {
		fmt.Println("\n🔍 Key Insights:")
		for _, insight := range result.Insights {
			severity := getSeverityIcon(insight.Severity)
			fmt.Printf("  %s %s: %s\n", severity, insight.Type, insight.Title)
			fmt.Printf("    %s\n", insight.Description)
			if insight.Recommendation != "" {
				fmt.Printf("    💡 Recommendation: %s\n", insight.Recommendation)
			}
		}
	}

	// Show main analysis
	fmt.Println("\n📝 Analysis:")
	fmt.Println(result.Summary)

	fmt.Println("\n" + strings.Repeat("=", 70))
	
	// Show next steps
	fmt.Println("\n💡 What's Next:")
	fmt.Println("- Run 'ship terraform-tools checkov-scan' for detailed security analysis")
	fmt.Println("- Use 'ship push' to analyze your infrastructure with Cloudship AI")
	fmt.Println("- Try more complex queries with the new Eino agent")
}

// getSeverityIcon returns an icon for the severity level
func getSeverityIcon(severity string) string {
	switch severity {
	case "critical":
		return "🚨"
	case "high":
		return "⚠️"
	case "medium":
		return "🔶"
	case "low":
		return "ℹ️"
	default:
		return "📄"
	}
}

// saveInvestigationResults saves the results to a JSON file
func saveInvestigationResults(result *agent.InvestigationResult, prompt string) error {
	// Create results directory
	resultsDir := "investigation_results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	// Generate filename based on timestamp and prompt
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	promptHash := fmt.Sprintf("%x", prompt)[:8] // First 8 chars of hash
	filename := fmt.Sprintf("investigation_%s_%s.json", timestamp, promptHash)
	filepath := filepath.Join(resultsDir, filename)

	// Marshal results to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %w", err)
	}

	fmt.Printf("📁 Results saved to: %s\n", filepath)
	return nil
}

// getProviderCredentials returns a map of environment variables for the specified provider
func getProviderCredentials(provider string) map[string]string {
	credentials := make(map[string]string)

	// Copy relevant environment variables based on provider
	switch provider {
	case "aws":
		// AWS credentials from environment
		if val := os.Getenv("AWS_ACCESS_KEY_ID"); val != "" {
			credentials["AWS_ACCESS_KEY_ID"] = val
		}
		if val := os.Getenv("AWS_SECRET_ACCESS_KEY"); val != "" {
			credentials["AWS_SECRET_ACCESS_KEY"] = val
		}
		if val := os.Getenv("AWS_SESSION_TOKEN"); val != "" {
			credentials["AWS_SESSION_TOKEN"] = val
		}
		if val := os.Getenv("AWS_REGION"); val != "" {
			credentials["AWS_REGION"] = val
		}
		if val := os.Getenv("AWS_PROFILE"); val != "" {
			credentials["AWS_PROFILE"] = val
		}

	case "azure":
		// Azure credentials from environment
		if val := os.Getenv("AZURE_CLIENT_ID"); val != "" {
			credentials["AZURE_CLIENT_ID"] = val
		}
		if val := os.Getenv("AZURE_CLIENT_SECRET"); val != "" {
			credentials["AZURE_CLIENT_SECRET"] = val
		}
		if val := os.Getenv("AZURE_TENANT_ID"); val != "" {
			credentials["AZURE_TENANT_ID"] = val
		}
		if val := os.Getenv("AZURE_SUBSCRIPTION_ID"); val != "" {
			credentials["AZURE_SUBSCRIPTION_ID"] = val
		}

	case "gcp":
		// GCP credentials from environment
		if val := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); val != "" {
			credentials["GOOGLE_APPLICATION_CREDENTIALS"] = val
		}
		if val := os.Getenv("GOOGLE_CLOUD_PROJECT"); val != "" {
			credentials["GOOGLE_CLOUD_PROJECT"] = val
		}
		if val := os.Getenv("GCLOUD_PROJECT"); val != "" {
			credentials["GCLOUD_PROJECT"] = val
		}
	}

	return credentials
}
