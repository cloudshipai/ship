package cli

import (
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var opencodeCmd = &cobra.Command{
	Use:   "opencode",
	Short: "AI-powered coding assistant with comprehensive development tools",
	Long: `OpenCode is an AI coding agent that provides interactive chat, code generation, 
analysis, review, refactoring, testing, and documentation capabilities in containerized environments.`,
}

var opencodeChatCmd = &cobra.Command{
	Use:   "chat <message>",
	Short: "Start an interactive chat session with OpenCode",
	Long:  `Start an interactive chat session with OpenCode AI assistant for coding questions and guidance.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runOpencodeChat,
}

var opencodeGenerateCmd = &cobra.Command{
	Use:   "generate <prompt>",
	Short: "Generate code based on a prompt",
	Long:  `Generate code based on natural language prompts using OpenCode AI.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runOpencodeGenerate,
}

var opencodeAnalyzeCmd = &cobra.Command{
	Use:   "analyze <file> <question>",
	Short: "Analyze a specific file with OpenCode",
	Long:  `Analyze a specific file and ask questions about its implementation, structure, or functionality.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runOpencodeAnalyze,
}

var opencodeReviewCmd = &cobra.Command{
	Use:   "review [target]",
	Short: "Perform code review on changes",
	Long:  `Perform comprehensive code review on changes or specific targets using OpenCode AI.`,
	RunE:  runOpencodeReview,
}

var opencodeRefactorCmd = &cobra.Command{
	Use:   "refactor <instructions>",
	Short: "Perform code refactoring based on instructions",
	Long:  `Refactor code based on natural language instructions using OpenCode AI.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runOpencodeRefactor,
}

var opencodeTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Generate and run tests for code",
	Long:  `Generate comprehensive tests for your code and optionally run them with coverage analysis.`,
	RunE:  runOpencodeTest,
}

var opencodeDocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Generate documentation for code",
	Long:  `Generate comprehensive documentation for your codebase in various formats.`,
	RunE:  runOpencodeDocument,
}

var opencodeVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get OpenCode version information",
	Long:  `Display the version of OpenCode AI coding agent.`,
	RunE:  runOpencodeVersion,
}

var opencodeInteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive OpenCode session",
	Long:  `Start an interactive session with OpenCode for extended coding conversations.`,
	RunE:  runOpencodeInteractive,
}

var opencodeBatchCmd = &cobra.Command{
	Use:   "batch <pattern> <operation>",
	Short: "Process multiple files with OpenCode",
	Long:  `Process multiple files matching a pattern with a specific operation using OpenCode AI.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runOpencodeBatch,
}

func init() {
	rootCmd.AddCommand(opencodeCmd)
	
	// Add subcommands
	opencodeCmd.AddCommand(opencodeChatCmd)
	opencodeCmd.AddCommand(opencodeGenerateCmd)
	opencodeCmd.AddCommand(opencodeAnalyzeCmd)
	opencodeCmd.AddCommand(opencodeReviewCmd)
	opencodeCmd.AddCommand(opencodeRefactorCmd)
	opencodeCmd.AddCommand(opencodeTestCmd)
	opencodeCmd.AddCommand(opencodeDocumentCmd)
	opencodeCmd.AddCommand(opencodeVersionCmd)
	opencodeCmd.AddCommand(opencodeInteractiveCmd)
	opencodeCmd.AddCommand(opencodeBatchCmd)

	// Flags for generate command
	opencodeGenerateCmd.Flags().StringP("output", "o", "", "Output file for generated code")
	
	// Flags for refactor command  
	opencodeRefactorCmd.Flags().StringSliceP("files", "f", []string{}, "Specific files to refactor")
	
	// Flags for test command
	opencodeTestCmd.Flags().StringP("type", "t", "", "Test type (unit, integration, e2e)")
	opencodeTestCmd.Flags().BoolP("coverage", "c", false, "Enable test coverage analysis")
	
	// Flags for document command
	opencodeDocumentCmd.Flags().StringP("format", "", "markdown", "Documentation format (markdown, html, pdf)")
	opencodeDocumentCmd.Flags().StringP("output-dir", "d", "docs", "Output directory for documentation")
	
	// Flags for chat command
	opencodeChatCmd.Flags().StringP("model", "m", "", "AI model to use (format: provider/model, e.g., 'openai/gpt-4o-mini', 'anthropic/claude-3-sonnet')")
	
	// Flags for interactive command
	opencodeInteractiveCmd.Flags().StringP("model", "m", "", "Specific AI model to use for interaction")
	
	// Global flags for work directory and persistence
	opencodeCmd.PersistentFlags().StringP("workdir", "w", ".", "Working directory for OpenCode operations")
	opencodeCmd.PersistentFlags().Bool("ephemeral", false, "Run in ephemeral mode - do not persist files to host")
	
	// Session support flags
	opencodeCmd.PersistentFlags().StringP("session", "s", "", "Session ID for multi-turn conversations")
	opencodeCmd.PersistentFlags().BoolP("continue", "c", false, "Continue the last session")
}

func runOpencodeChat(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	ephemeral, _ := cmd.Flags().GetBool("ephemeral")
	sessionID, _ := cmd.Flags().GetString("session")
	continueSession, _ := cmd.Flags().GetBool("continue")
	model, _ := cmd.Flags().GetString("model")
	message := strings.Join(args, " ")

	if sessionID != "" {
		fmt.Printf("Starting OpenCode chat session (session: %s)...\n", sessionID)
	} else if continueSession {
		fmt.Printf("Continuing OpenCode chat session...\n")
	} else {
		fmt.Printf("Starting OpenCode chat session...\n")
	}
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.ChatWithSessionAndModel(ctx, workDir, message, !ephemeral, sessionID, continueSession, model)
	if err != nil {
		return fmt.Errorf("failed to run opencode chat: %w", err)
	}

	fmt.Printf("OpenCode Response:\n%s\n", output)
	return nil
}

func runOpencodeGenerate(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	outputFile, _ := cmd.Flags().GetString("output")
	prompt := strings.Join(args, " ")

	fmt.Printf("Generating code with OpenCode...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.Generate(ctx, workDir, prompt, outputFile)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	if outputFile != "" {
		green := color.New(color.FgGreen)
		green.Printf("✓ Code generated successfully to: %s\n", outputFile)
	}
	
	fmt.Printf("Generated Code:\n%s\n", output)
	return nil
}

func runOpencodeAnalyze(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	filePath := args[0]
	question := args[1]

	fmt.Printf("Analyzing file '%s' with OpenCode...\n", filePath)
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.AnalyzeFile(ctx, filePath, question)
	if err != nil {
		return fmt.Errorf("failed to analyze file: %w", err)
	}

	fmt.Printf("Analysis Result:\n%s\n", output)
	return nil
}

func runOpencodeReview(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	var target string
	if len(args) > 0 {
		target = args[0]
	}

	fmt.Printf("Performing code review with OpenCode...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.Review(ctx, workDir, target)
	if err != nil {
		return fmt.Errorf("failed to review code: %w", err)
	}

	fmt.Printf("Code Review:\n%s\n", output)
	return nil
}

func runOpencodeRefactor(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	files, _ := cmd.Flags().GetStringSlice("files")
	instructions := strings.Join(args, " ")

	fmt.Printf("Refactoring code with OpenCode...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.Refactor(ctx, workDir, instructions, files)
	if err != nil {
		return fmt.Errorf("failed to refactor code: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Printf("✓ Code refactoring completed\n")
	fmt.Printf("Refactoring Result:\n%s\n", output)
	return nil
}

func runOpencodeTest(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	testType, _ := cmd.Flags().GetString("type")
	coverage, _ := cmd.Flags().GetBool("coverage")

	fmt.Printf("Generating and running tests with OpenCode...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.Test(ctx, workDir, testType, coverage)
	if err != nil {
		return fmt.Errorf("failed to run tests: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Printf("✓ Test generation and execution completed\n")
	fmt.Printf("Test Results:\n%s\n", output)
	return nil
}

func runOpencodeDocument(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	format, _ := cmd.Flags().GetString("format")
	outputDir, _ := cmd.Flags().GetString("output-dir")

	fmt.Printf("Generating documentation with OpenCode...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.Document(ctx, workDir, format, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Printf("✓ Documentation generated successfully in: %s\n", outputDir)
	fmt.Printf("Documentation Result:\n%s\n", output)
	return nil
}

func runOpencodeVersion(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	fmt.Printf("Getting OpenCode version information...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get opencode version: %w", err)
	}

	fmt.Printf("OpenCode Version:\n%s\n", output)
	return nil
}

func runOpencodeInteractive(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	model, _ := cmd.Flags().GetString("model")

	fmt.Printf("Starting interactive OpenCode session...\n")
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.Interactive(ctx, workDir, model)
	if err != nil {
		return fmt.Errorf("failed to start interactive session: %w", err)
	}

	fmt.Printf("Interactive Session:\n%s\n", output)
	return nil
}

func runOpencodeBatch(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	workDir, _ := cmd.Flags().GetString("workdir")
	pattern := args[0]
	operation := args[1]

	fmt.Printf("Processing files matching '%s' with operation '%s'...\n", pattern, operation)
	
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	opencode := engine.OpenCode()
	output, err := opencode.BatchProcess(ctx, workDir, pattern, operation)
	if err != nil {
		return fmt.Errorf("failed to batch process: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Printf("✓ Batch processing completed\n")
	fmt.Printf("Batch Processing Result:\n%s\n", output)
	return nil
}