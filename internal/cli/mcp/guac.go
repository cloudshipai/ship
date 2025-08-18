package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGuacTools adds GUAC (Graph for Understanding Artifact Composition) MCP tool implementations using real CLI tools
func AddGuacTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// GUAC collect and ingest tool
	collectTool := mcp.NewTool("guac_collect",
		mcp.WithDescription("Run GUAC collector to gather supply chain metadata"),
		mcp.WithString("source_type",
			mcp.Description("Type of source to collect from"),
			mcp.Enum("files", "deps_dev", "osv", "scorecard"),
			mcp.Required(),
		),
		mcp.WithString("path",
			mcp.Description("Path to files or directory (for files source type)"),
		),
	)
	s.AddTool(collectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceType := request.GetString("source_type", "")
		args := []string{"guacone", "collector", sourceType}
		
		if path := request.GetString("path", ""); path != "" && sourceType == "files" {
			args = append(args, "--path", path)
		}
		
		return executeShipCommand(args)
	})

	// GUAC certifier tool
	certifyTool := mcp.NewTool("guac_certify",
		mcp.WithDescription("Run GUAC certifier to add metadata attestations"),
		mcp.WithString("certifier_type",
			mcp.Description("Type of certifier to run"),
			mcp.Enum("osv", "scorecard", "vulns"),
			mcp.Required(),
		),
		mcp.WithBoolean("poll",
			mcp.Description("Enable polling mode for continuous certification"),
		),
	)
	s.AddTool(certifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certifierType := request.GetString("certifier_type", "")
		args := []string{"guacone", "certify", certifierType}
		
		if request.GetBool("poll", false) {
			args = append(args, "--poll")
		}
		
		return executeShipCommand(args)
	})

	// GUAC query tool
	queryTool := mcp.NewTool("guac_query",
		mcp.WithDescription("Run GUAC canned queries for analysis"),
		mcp.WithString("query_type",
			mcp.Description("Type of query to run"),
			mcp.Enum("vulnerabilities", "dependencies", "packages"),
			mcp.Required(),
		),
		mcp.WithString("subject",
			mcp.Description("Subject to query (package name, image, etc.)"),
		),
	)
	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		queryType := request.GetString("query_type", "")
		args := []string{"guacone", "query", queryType}
		
		if subject := request.GetString("subject", ""); subject != "" {
			args = append(args, "--subject", subject)
		}
		
		return executeShipCommand(args)
	})

	// GUAC GraphQL server tool
	startGraphQLTool := mcp.NewTool("guac_start_graphql",
		mcp.WithDescription("Start GUAC GraphQL server"),
		mcp.WithString("backend",
			mcp.Description("Backend type to use"),
			mcp.Enum("keyvalue", "postgresql", "redis", "tikv"),
		),
		mcp.WithString("port",
			mcp.Description("Port for GraphQL server"),
		),
	)
	s.AddTool(startGraphQLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"guacgql"}
		
		if backend := request.GetString("backend", ""); backend != "" {
			args = append(args, "--gql-backend", backend)
		}
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "--gql-listen-port", port)
		}
		
		return executeShipCommand(args)
	})

	// GUAC ingest tool
	ingestTool := mcp.NewTool("guac_ingest",
		mcp.WithDescription("Run GUAC ingestor service connected to NATS and GraphQL"),
		mcp.WithString("nats_addr",
			mcp.Description("NATS server address"),
		),
		mcp.WithString("gql_addr",
			mcp.Description("GraphQL server address"),
		),
	)
	s.AddTool(ingestTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"guacingest"}
		
		if natsAddr := request.GetString("nats_addr", ""); natsAddr != "" {
			args = append(args, "--nats-addr", natsAddr)
		}
		if gqlAddr := request.GetString("gql_addr", ""); gqlAddr != "" {
			args = append(args, "--gql-addr", gqlAddr)
		}
		
		return executeShipCommand(args)
	})

	// GUAC collector-subscriber service tool
	collectorSubTool := mcp.NewTool("guac_collector_subscriber",
		mcp.WithDescription("Run GUAC collector-subscriber service"),
		mcp.WithString("nats_addr",
			mcp.Description("NATS server address"),
		),
		mcp.WithString("source_type",
			mcp.Description("Source type to collect from"),
			mcp.Enum("files", "deps_dev", "osv", "scorecard"),
		),
	)
	s.AddTool(collectorSubTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"guaccsub"}
		
		if natsAddr := request.GetString("nats_addr", ""); natsAddr != "" {
			args = append(args, "--nats-addr", natsAddr)
		}
		if sourceType := request.GetString("source_type", ""); sourceType != "" {
			args = append(args, "--source", sourceType)
		}
		
		return executeShipCommand(args)
	})
}