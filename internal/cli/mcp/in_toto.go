package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddInTotoTools adds in-toto (supply chain attestation) MCP tool implementations using direct Dagger calls
func AddInTotoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addInTotoToolsDirect(s)
}

// addInTotoToolsDirect adds in-toto tools using direct Dagger module calls
func addInTotoToolsDirect(s *server.MCPServer) {
	// in-toto-run tool
	runStepTool := mcp.NewTool("in_toto_run_step",
		mcp.WithDescription("Run in-toto supply chain step with attestation using in-toto-run"),
		mcp.WithString("step_name",
			mcp.Description("Name for link metadata file"),
			mcp.Required(),
		),
		mcp.WithString("command",
			mcp.Description("Command to execute for this step"),
			mcp.Required(),
		),
		mcp.WithString("signing_key",
			mcp.Description("Path to signing key in PKCS8/PEM format"),
		),
		mcp.WithString("gpg_keyid",
			mcp.Description("GPG keyid to sign metadata"),
		),
		mcp.WithString("materials",
			mcp.Description("Paths to files/directories to record before command"),
		),
		mcp.WithString("products",
			mcp.Description("Paths to files/directories to record after command"),
		),
		mcp.WithBoolean("record_streams",
			mcp.Description("Capture stdout/stderr"),
		),
		mcp.WithBoolean("no_command",
			mcp.Description("Generate metadata without executing command"),
		),
	)
	s.AddTool(runStepTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInTotoModule(client)

		// Get parameters
		stepName := request.GetString("step_name", "")
		if stepName == "" {
			return mcp.NewToolResultError("step_name is required"), nil
		}

		command := request.GetString("command", "")
		if command == "" && !request.GetBool("no_command", false) {
			return mcp.NewToolResultError("command is required unless no_command is true"), nil
		}

		// Prepare options
		var opts []modules.InTotoOption
		
		if signingKey := request.GetString("signing_key", ""); signingKey != "" {
			opts = append(opts, modules.WithKeyPath(signingKey))
		}
		
		if materials := request.GetString("materials", ""); materials != "" {
			materialsList := strings.Split(materials, ",")
			opts = append(opts, modules.WithMaterials(materialsList))
		}
		
		if products := request.GetString("products", ""); products != "" {
			productsList := strings.Split(products, ",")
			opts = append(opts, modules.WithProducts(productsList))
		}

		// Execute step or record metadata
		var container *dagger.Container
		if request.GetBool("no_command", false) {
			// Record metadata without running command
			container, err = module.RecordMetadata(ctx, stepName, opts...)
		} else {
			// Run step with command
			commandParts := strings.Fields(command)
			container, err = module.RunStep(ctx, stepName, commandParts, opts...)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to run step: %v", err)), nil
		}

		// Get output
		output, err := container.Stdout(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get output: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// in-toto-verify tool
	verifySupplyChainTool := mcp.NewTool("in_toto_verify",
		mcp.WithDescription("Verify in-toto supply chain integrity using in-toto-verify"),
		mcp.WithString("layout_path",
			mcp.Description("Path to supply chain layout file"),
			mcp.Required(),
		),
		mcp.WithString("layout_key_paths",
			mcp.Description("Comma-separated paths to layout verification keys"),
		),
		mcp.WithString("link_dir",
			mcp.Description("Directory containing link metadata"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
	)
	s.AddTool(verifySupplyChainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInTotoModule(client)

		// Get parameters
		layoutPath := request.GetString("layout_path", "")
		if layoutPath == "" {
			return mcp.NewToolResultError("layout_path is required"), nil
		}

		// Prepare options
		var opts []modules.InTotoOption
		
		if layoutKeyPaths := request.GetString("layout_key_paths", ""); layoutKeyPaths != "" {
			keysList := strings.Split(layoutKeyPaths, ",")
			opts = append(opts, modules.WithPublicKeys(keysList))
		}
		
		if linkDir := request.GetString("link_dir", ""); linkDir != "" {
			opts = append(opts, modules.WithLinkDir(linkDir))
		}

		// Verify supply chain
		container, err := module.VerifySupplyChain(ctx, layoutPath, opts...)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to verify supply chain: %v", err)), nil
		}

		// Get output
		output, err := container.Stdout(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get output: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// in-toto-record tool
	recordTool := mcp.NewTool("in_toto_record",
		mcp.WithDescription("Record in-toto step metadata using in-toto-record"),
		mcp.WithString("step_name",
			mcp.Description("Name for link metadata file"),
			mcp.Required(),
		),
		mcp.WithString("operation",
			mcp.Description("Record operation"),
			mcp.Required(),
			mcp.Enum("start", "stop"),
		),
		mcp.WithString("signing_key",
			mcp.Description("Path to signing key in PKCS8/PEM format"),
		),
		mcp.WithString("gpg_keyid",
			mcp.Description("GPG keyid to sign metadata"),
		),
		mcp.WithString("materials",
			mcp.Description("Paths to files/directories to record as materials"),
		),
		mcp.WithString("products",
			mcp.Description("Paths to files/directories to record as products"),
		),
	)
	s.AddTool(recordTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInTotoModule(client)

		// Get parameters
		stepName := request.GetString("step_name", "")
		if stepName == "" {
			return mcp.NewToolResultError("step_name is required"), nil
		}

		operation := request.GetString("operation", "")
		if operation == "" {
			return mcp.NewToolResultError("operation is required"), nil
		}

		// Prepare options
		var opts []modules.InTotoOption
		
		if signingKey := request.GetString("signing_key", ""); signingKey != "" {
			opts = append(opts, modules.WithKeyPath(signingKey))
		}
		
		if materials := request.GetString("materials", ""); materials != "" {
			materialsList := strings.Split(materials, ",")
			if operation == "start" {
				opts = append(opts, modules.WithMaterials(materialsList))
			}
		}
		
		if products := request.GetString("products", ""); products != "" {
			productsList := strings.Split(products, ",")
			if operation == "stop" {
				opts = append(opts, modules.WithProducts(productsList))
			}
		}

		// Record metadata
		container, err := module.RecordMetadata(ctx, fmt.Sprintf("%s-%s", stepName, operation), opts...)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to record metadata: %v", err)), nil
		}

		// Get output
		output, err := container.Stdout(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get output: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// in-toto-sign tool (or generate layout)
	signTool := mcp.NewTool("in_toto_sign",
		mcp.WithDescription("Sign in-toto metadata or generate layout"),
		mcp.WithString("metadata_file",
			mcp.Description("Path to metadata file to sign (or empty to generate layout)"),
		),
		mcp.WithString("signing_key",
			mcp.Description("Path to signing key in PKCS8/PEM format"),
		),
		mcp.WithString("gpg_keyid",
			mcp.Description("GPG keyid to sign metadata"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for signed metadata"),
		),
	)
	s.AddTool(signTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInTotoModule(client)

		// Prepare options
		var opts []modules.InTotoOption
		
		if signingKey := request.GetString("signing_key", ""); signingKey != "" {
			opts = append(opts, modules.WithKeyPath(signingKey))
		}

		// If no metadata file, generate layout
		if metadataFile := request.GetString("metadata_file", ""); metadataFile == "" {
			// Generate layout
			container, err := module.GenerateLayout(ctx, opts...)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to generate layout: %v", err)), nil
			}

			// Get output
			output, err := container.Stdout(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get output: %v", err)), nil
			}

			if outputFile := request.GetString("output_file", ""); outputFile != "" {
				output += fmt.Sprintf("\n\nLayout should be saved to: %s", outputFile)
			}

			return mcp.NewToolResultText(output), nil
		}

		// Sign existing metadata - use record with the file as material
		opts = append(opts, modules.WithMaterials([]string{request.GetString("metadata_file", "")}))
		container, err := module.RecordMetadata(ctx, "sign-metadata", opts...)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to sign metadata: %v", err)), nil
		}

		// Get output
		output, err := container.Stdout(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get output: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}