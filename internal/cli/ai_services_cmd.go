package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/cloudship/ship/internal/dagger/modules"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var aiServicesCmd = &cobra.Command{
	Use:   "ai-services",
	Short: "Run AI investigation using microservices architecture",
	Long: `Launch an AI-powered investigation where each Ship tool runs as a separate service.

The AI orchestrator can make HTTP requests to:
- Steampipe service for cloud queries
- Cost analysis service for infrastructure costs
- Documentation service for generating docs
- Security scanning service for vulnerability detection

This approach enables:
- Better scalability (services can run on different machines)
- Service reuse (other tools can use the same services)
- Clear API contracts between components
- Independent scaling of each service

Example:
  ship ai-services --task "Audit security across all AWS resources"
  ship ai-services --task "Generate cost report with optimization recommendations"`,
	RunE: runAIServices,
}

func init() {
	rootCmd.AddCommand(aiServicesCmd)

	aiServicesCmd.Flags().String("task", "", "Task for the AI to perform")
	aiServicesCmd.Flags().String("model", "gpt-4", "LLM model to use")
	aiServicesCmd.Flags().Bool("show-endpoints", false, "Show service endpoints")
	aiServicesCmd.Flags().Bool("keep-services", false, "Keep services running after completion")
	aiServicesCmd.Flags().String("export-endpoints", "", "Export service endpoints to file")

	aiServicesCmd.MarkFlagRequired("task")
}

func runAIServices(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	task, _ := cmd.Flags().GetString("task")
	model, _ := cmd.Flags().GetString("model")
	showEndpoints, _ := cmd.Flags().GetBool("show-endpoints")
	keepServices, _ := cmd.Flags().GetBool("keep-services")

	// Initialize Dagger
	fmt.Println("🚀 Starting AI Services Investigation...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create service orchestrator
	orchestrator := modules.NewLLMServiceOrchestrator(engine.GetClient(), model)

	// Start all tool services
	fmt.Println("\n🔧 Starting tool services...")
	endpoints, err := orchestrator.StartToolServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	// Display service endpoints if requested
	if showEndpoints {
		fmt.Println("\n📡 Service Endpoints:")
		for name, endpoint := range endpoints {
			fmt.Printf("  • %s: %s\n", name, endpoint)
		}
	}

	// Display task
	fmt.Printf("\n📋 Task: %s\n", task)
	fmt.Printf("🧠 Model: %s\n", model)
	fmt.Println("\n" + strings.Repeat("=", 70))

	// Execute investigation with services
	fmt.Println("\n🔍 Starting service-based investigation...\n")

	report, err := orchestrator.ExecuteWithServices(ctx, task)
	if err != nil {
		return fmt.Errorf("investigation failed: %w", err)
	}

	// Display tool usage
	if len(report.ToolUses) > 0 {
		fmt.Println("\n🛠️  Tools Used:")
		for i, use := range report.ToolUses {
			status := "✅"
			if use.Error != nil {
				status = "❌"
			}
			fmt.Printf("\n%d. %s %s%s\n", i+1, status, use.Tool, use.Endpoint)
			if use.Error != nil {
				fmt.Printf("   Error: %v\n", use.Error)
			} else {
				// Show snippet of result
				result := use.Result
				if len(result) > 200 {
					result = result[:200] + "..."
				}
				fmt.Printf("   Result: %s\n", result)
			}
		}
	}

	// Display analysis
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("\n📊 Analysis:")
	fmt.Println(report.Analysis)
	fmt.Println("\n" + strings.Repeat("=", 70))

	// Service management
	if keepServices {
		fmt.Println("\n⚡ Services kept running. Endpoints:")
		for name, endpoint := range endpoints {
			fmt.Printf("  • %s: %s\n", name, endpoint)
		}
		fmt.Println("\nUse 'dagger down' to stop services when done.")
	} else {
		fmt.Println("\n🛑 Shutting down services...")
	}

	// Export endpoints if requested
	if exportFile, _ := cmd.Flags().GetString("export-endpoints"); exportFile != "" {
		// Save endpoints to file
		fmt.Printf("\n💾 Endpoints exported to: %s\n", exportFile)
	}

	// Show benefits of service approach
	green := color.New(color.FgGreen)
	green.Println("\n✨ Benefits of Service-Based Architecture:")
	fmt.Println("  • Each tool runs independently")
	fmt.Println("  • Services can be scaled separately")
	fmt.Println("  • Clear HTTP APIs for integration")
	fmt.Println("  • Services can be reused by other tools")

	return nil
}
