package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInTotoTools adds in-toto (supply chain attestation) MCP tool implementations
func AddInTotoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// in-toto run step tool
	runStepTool := mcp.NewTool("in_toto_run_step",
		mcp.WithDescription("Run in-toto supply chain step with attestation"),
		mcp.WithString("step_name",
			mcp.Description("Name of the supply chain step"),
			mcp.Required(),
		),
		mcp.WithString("command",
			mcp.Description("Command to execute for this step"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to signing key"),
		),
		mcp.WithString("materials",
			mcp.Description("Comma-separated list of material files"),
		),
		mcp.WithString("products",
			mcp.Description("Comma-separated list of product files"),
		),
	)
	s.AddTool(runStepTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stepName := request.GetString("step_name", "")
		command := request.GetString("command", "")
		args := []string{"security", "in-toto", "--run-step", stepName, "--command", command}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		if materials := request.GetString("materials", ""); materials != "" {
			args = append(args, "--materials", materials)
		}
		if products := request.GetString("products", ""); products != "" {
			args = append(args, "--products", products)
		}
		return executeShipCommand(args)
	})

	// in-toto verify supply chain tool
	verifySupplyChainTool := mcp.NewTool("in_toto_verify_supply_chain",
		mcp.WithDescription("Verify in-toto supply chain integrity"),
		mcp.WithString("layout_path",
			mcp.Description("Path to supply chain layout file"),
			mcp.Required(),
		),
		mcp.WithString("layout_key_path",
			mcp.Description("Path to layout verification key"),
		),
		mcp.WithString("link_dir",
			mcp.Description("Directory containing link metadata"),
		),
	)
	s.AddTool(verifySupplyChainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		layoutPath := request.GetString("layout_path", "")
		args := []string{"security", "in-toto", "--verify", layoutPath}
		if layoutKeyPath := request.GetString("layout_key_path", ""); layoutKeyPath != "" {
			args = append(args, "--layout-key", layoutKeyPath)
		}
		if linkDir := request.GetString("link_dir", ""); linkDir != "" {
			args = append(args, "--link-dir", linkDir)
		}
		return executeShipCommand(args)
	})

	// in-toto generate layout tool
	generateLayoutTool := mcp.NewTool("in_toto_generate_layout",
		mcp.WithDescription("Generate in-toto supply chain layout"),
		mcp.WithString("output_path",
			mcp.Description("Output path for layout file"),
		),
		mcp.WithString("steps",
			mcp.Description("JSON string defining supply chain steps"),
		),
		mcp.WithString("inspections",
			mcp.Description("JSON string defining inspections"),
		),
	)
	s.AddTool(generateLayoutTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "in-toto", "--generate-layout"}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		if steps := request.GetString("steps", ""); steps != "" {
			args = append(args, "--steps", steps)
		}
		if inspections := request.GetString("inspections", ""); inspections != "" {
			args = append(args, "--inspections", inspections)
		}
		return executeShipCommand(args)
	})

	// in-toto record metadata tool
	recordMetadataTool := mcp.NewTool("in_toto_record_metadata",
		mcp.WithDescription("Record in-toto step metadata"),
		mcp.WithString("step_name",
			mcp.Description("Name of the supply chain step"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to signing key"),
		),
		mcp.WithString("materials",
			mcp.Description("Comma-separated list of material files"),
		),
		mcp.WithString("products",
			mcp.Description("Comma-separated list of product files"),
		),
	)
	s.AddTool(recordMetadataTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stepName := request.GetString("step_name", "")
		args := []string{"security", "in-toto", "--record", stepName}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		if materials := request.GetString("materials", ""); materials != "" {
			args = append(args, "--materials", materials)
		}
		if products := request.GetString("products", ""); products != "" {
			args = append(args, "--products", products)
		}
		return executeShipCommand(args)
	})
}