package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddGuacTools adds GUAC (Graph for Understanding Artifact Composition) MCP tool implementations using direct Dagger calls
func AddGuacTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGuacToolsDirect(s)
}

// addGuacToolsDirect adds GUAC tools using direct Dagger module calls
func addGuacToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGuacModule(client)

		// Get source type and path
		sourceType := request.GetString("source_type", "")
		path := request.GetString("path", ".")
		
		if sourceType == "files" && path != "" {
			// Collect files
			output, err := module.CollectFiles(ctx, path)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to collect files: %v", err)), nil
			}
			return mcp.NewToolResultText(output), nil
		}

		// For other source types, use generic analysis
		output, err := module.AnalyzeArtifact(ctx, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to analyze artifact: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// GUAC certifier tool
	certifyTool := mcp.NewTool("guac_certify",
		mcp.WithDescription("Run GUAC certifier to add metadata attestations"),
		mcp.WithString("certifier_type",
			mcp.Description("Type of certifier to run"),
			mcp.Enum("osv", "scorecard", "vulns"),
			mcp.Required(),
		),
		mcp.WithString("attestation_path",
			mcp.Description("Path to attestation file to validate"),
		),
		mcp.WithBoolean("poll",
			mcp.Description("Enable polling mode for continuous certification"),
		),
	)
	s.AddTool(certifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGuacModule(client)

		// Get attestation path if provided
		attestationPath := request.GetString("attestation_path", "")
		
		if attestationPath != "" {
			// Validate attestation
			output, err := module.ValidateAttestation(ctx, attestationPath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to validate attestation: %v", err)), nil
			}
			return mcp.NewToolResultText(output), nil
		}

		// Default to artifact analysis
		output, err := module.AnalyzeArtifact(ctx, ".")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to certify: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGuacModule(client)

		// Get query type and subject
		queryType := request.GetString("query_type", "")
		subject := request.GetString("subject", "")
		
		if subject == "" {
			return mcp.NewToolResultError("subject is required"), nil
		}

		var output string
		switch queryType {
		case "vulnerabilities":
			output, err = module.QueryVulnerabilities(ctx, subject)
		case "dependencies":
			output, err = module.QueryDependencies(ctx, subject)
		default:
			// For packages and other queries, use dependencies
			output, err = module.QueryDependencies(ctx, subject)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to query %s: %v", queryType, err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// GUAC GraphQL server tool (informational)
	startGraphQLTool := mcp.NewTool("guac_start_graphql",
		mcp.WithDescription("Get information about starting GUAC GraphQL server"),
		mcp.WithString("backend",
			mcp.Description("Backend type to use"),
			mcp.Enum("keyvalue", "postgresql", "redis", "tikv"),
		),
		mcp.WithString("port",
			mcp.Description("Port for GraphQL server"),
		),
	)
	s.AddTool(startGraphQLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		backend := request.GetString("backend", "keyvalue")
		port := request.GetString("port", "8080")
		
		info := fmt.Sprintf("To start GUAC GraphQL server:\n")
		info += fmt.Sprintf("guacgql --gql-backend %s --gql-listen-port %s\n\n", backend, port)
		info += fmt.Sprintf("This will start the GraphQL API server on port %s using %s backend.\n", port, backend)
		info += fmt.Sprintf("Access the GraphQL playground at: http://localhost:%s/graphql", port)

		return mcp.NewToolResultText(info), nil
	})

	// GUAC ingest tool
	ingestTool := mcp.NewTool("guac_ingest",
		mcp.WithDescription("Ingest SBOM or other supply chain data into GUAC"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file to ingest"),
			mcp.Required(),
		),
		mcp.WithString("nats_addr",
			mcp.Description("NATS server address"),
		),
		mcp.WithString("gql_addr",
			mcp.Description("GraphQL server address"),
		),
	)
	s.AddTool(ingestTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGuacModule(client)

		// Get SBOM path
		sbomPath := request.GetString("sbom_path", "")
		if sbomPath == "" {
			return mcp.NewToolResultError("sbom_path is required"), nil
		}

		// Ingest SBOM
		output, err := module.IngestSBOM(ctx, sbomPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to ingest SBOM: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// GUAC collector-subscriber service tool
	collectorSubTool := mcp.NewTool("guac_collector_subscriber",
		mcp.WithDescription("Generate supply chain graph for a project"),
		mcp.WithString("project_path",
			mcp.Description("Path to project to analyze"),
			mcp.Required(),
		),
		mcp.WithString("source_type",
			mcp.Description("Source type to collect from"),
			mcp.Enum("files", "deps_dev", "osv", "scorecard"),
		),
	)
	s.AddTool(collectorSubTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGuacModule(client)

		// Get project path
		projectPath := request.GetString("project_path", ".")
		
		// Generate graph
		output, err := module.GenerateGraph(ctx, projectPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate graph: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}