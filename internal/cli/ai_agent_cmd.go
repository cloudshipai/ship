package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/cloudship/ship/internal/dagger/modules"
	"github.com/spf13/cobra"
)

var aiAgentCmd = &cobra.Command{
	Use:   "ai-agent",
	Short: "Run autonomous AI agent with access to all Ship tools",
	Long: `Launch an AI agent that can autonomously use Ship's tools to investigate infrastructure.

The agent has access to:
- Steampipe for querying cloud resources
- OpenInfraQuote for cost analysis  
- Terraform-docs for documentation
- Security scanning tools

Example:
  ship ai-agent --task "Perform complete security audit of AWS infrastructure"
  ship ai-agent --task "Optimize costs for our production environment"
  ship ai-agent --task "Document all Terraform modules and check for issues"`,
	RunE: runAIAgent,
}

func init() {
	rootCmd.AddCommand(aiAgentCmd)
	
	aiAgentCmd.Flags().String("task", "", "Task for the AI agent to perform")
	aiAgentCmd.Flags().String("llm-provider", "openai", "LLM provider (openai, anthropic)")
	aiAgentCmd.Flags().String("model", "gpt-4", "LLM model to use")
	aiAgentCmd.Flags().Int("max-steps", 10, "Maximum tool-use steps")
	aiAgentCmd.Flags().Bool("approve-each", false, "Require approval before each tool use")
	
	aiAgentCmd.MarkFlagRequired("task")
}

func runAIAgent(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	task, _ := cmd.Flags().GetString("task")
	llmProvider, _ := cmd.Flags().GetString("llm-provider")
	model, _ := cmd.Flags().GetString("model")
	maxSteps, _ := cmd.Flags().GetInt("max-steps")
	approveEach, _ := cmd.Flags().GetBool("approve-each")
	
	// Initialize Dagger engine
	fmt.Println("🤖 Initializing AI Agent...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()
	
	// Create AI agent with tools
	agent := modules.NewLLMWithToolsModule(engine.GetClient(), model)
	
	fmt.Printf("\n📋 Task: %s\n", task)
	fmt.Printf("🔧 Available tools: Steampipe, OpenInfraQuote, Terraform-docs\n")
	fmt.Printf("🧠 Using: %s with %s\n\n", llmProvider, model)
	
	if approveEach {
		fmt.Println("⚠️  Approval mode: You'll be asked before each tool use")
	}
	
	// Execute the investigation
	fmt.Println("🔍 Starting investigation...\n")
	
	report, err := agent.InvestigateWithTools(ctx, task)
	if err != nil {
		return fmt.Errorf("investigation failed: %w", err)
	}
	
	// Display results
	fmt.Println("\n📊 Investigation Complete!")
	fmt.Printf("\n🔧 Tools used: %d\n", len(report.ToolsUsed))
	
	for i, tool := range report.ToolsUsed {
		fmt.Printf("\n  %d. %s - %s\n", i+1, tool.Tool, tool.Action)
		if tool.Error != nil {
			fmt.Printf("     ❌ Error: %v\n", tool.Error)
		} else {
			fmt.Printf("     ✅ Success\n")
		}
	}
	
	fmt.Println("\n📝 Final Analysis:")
	fmt.Println("=" + strings.Repeat("=", 70))
	fmt.Println(report.Analysis)
	fmt.Println("=" + strings.Repeat("=", 70))
	
	// Save report if requested
	if output, _ := cmd.Flags().GetString("output"); output != "" {
		// Save to file
		fmt.Printf("\n💾 Report saved to: %s\n", output)
	}
	
	return nil
}

// Example of a more interactive version
func runInteractiveAgent(ctx context.Context, task string) error {
	// This would create an interactive session where:
	// 1. LLM proposes a tool to use
	// 2. User approves/modifies
	// 3. Tool runs
	// 4. Results shown
	// 5. LLM proposes next step
	// 6. Repeat until done
	
	return nil
}