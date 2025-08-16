package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGoldilocksTools adds Goldilocks (Kubernetes resource recommendations) MCP tool implementations
func AddGoldilocksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Goldilocks analyze tool
	analyzeTool := mcp.NewTool("goldilocks_analyze",
		mcp.WithDescription("Analyze Kubernetes deployments for resource recommendations"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to analyze"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(analyzeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "goldilocks", "analyze"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Goldilocks install tool
	installTool := mcp.NewTool("goldilocks_install",
		mcp.WithDescription("Install Goldilocks in Kubernetes cluster"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for installation"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "goldilocks", "install"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Goldilocks dashboard tool
	dashboardTool := mcp.NewTool("goldilocks_dashboard",
		mcp.WithDescription("Access Goldilocks dashboard for resource recommendations"),
		mcp.WithNumber("port",
			mcp.Description("Port for dashboard"),
		),
	)
	s.AddTool(dashboardTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "goldilocks", "dashboard"}
		if port := request.GetInt("port", 0); port > 0 {
			args = append(args, "--port", string(rune(port)))
		}
		return executeShipCommand(args)
	})

	// Goldilocks get version tool
	getVersionTool := mcp.NewTool("goldilocks_get_version",
		mcp.WithDescription("Get Goldilocks version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "goldilocks", "--version"}
		return executeShipCommand(args)
	})
}