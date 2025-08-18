package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddStepCATools adds Step CA (Certificate Authority) MCP tool implementations using direct Dagger calls
func AddStepCATools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addStepCAToolsDirect(s)
}

// addStepCAToolsDirect adds Step CA tools using direct Dagger module calls
func addStepCAToolsDirect(s *server.MCPServer) {
	// Step CA init tool
	initTool := mcp.NewTool("step_ca_init",
		mcp.WithDescription("Initialize Step CA certificate authority using real Step CA CLI"),
		mcp.WithString("ca_name",
			mcp.Description("Name of the certificate authority"),
			mcp.Required(),
		),
		mcp.WithString("dns_name",
			mcp.Description("DNS name for the CA"),
			mcp.Required(),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewStepCAModule(client)

		// Get parameters
		caName := request.GetString("ca_name", "")
		dnsName := request.GetString("dns_name", "")

		// Initialize CA
		output, err := module.InitCA(ctx, caName, dnsName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Step CA init failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Step CA create certificate tool
	createCertificateTool := mcp.NewTool("step_ca_create_certificate",
		mcp.WithDescription("Create certificate using real Step CA CLI"),
		mcp.WithString("subject",
			mcp.Description("Certificate subject (CN)"),
			mcp.Required(),
		),
		mcp.WithString("ca_url",
			mcp.Description("Certificate authority URL"),
			mcp.Required(),
		),
		mcp.WithString("root_cert",
			mcp.Description("Root certificate path"),
			mcp.Required(),
		),
	)
	s.AddTool(createCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewStepCAModule(client)

		// Get parameters
		subject := request.GetString("subject", "")
		caURL := request.GetString("ca_url", "")
		rootCert := request.GetString("root_cert", "")

		// Create certificate
		output, err := module.CreateCertificate(ctx, subject, caURL, rootCert)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Step CA create certificate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Step CA renew certificate tool
	renewCertificateTool := mcp.NewTool("step_ca_renew_certificate",
		mcp.WithDescription("Renew certificate using real Step CA CLI"),
		mcp.WithString("cert_path",
			mcp.Description("Certificate file path"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Private key file path"),
			mcp.Required(),
		),
		mcp.WithString("ca_url",
			mcp.Description("Certificate authority URL"),
			mcp.Required(),
		),
	)
	s.AddTool(renewCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewStepCAModule(client)

		// Get parameters
		certPath := request.GetString("cert_path", "")
		keyPath := request.GetString("key_path", "")
		caURL := request.GetString("ca_url", "")

		// Renew certificate
		output, err := module.RenewCertificate(ctx, certPath, keyPath, caURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Step CA renew certificate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Step CA add provisioner tool
	addProvisionerTool := mcp.NewTool("step_ca_add_provisioner",
		mcp.WithDescription("Add provisioner to Step CA using real Step CA CLI"),
		mcp.WithString("provisioner_name",
			mcp.Description("Provisioner name"),
			mcp.Required(),
		),
		mcp.WithString("provisioner_type",
			mcp.Description("Provisioner type (e.g., JWK, OIDC, ACME)"),
			mcp.Required(),
		),
		mcp.WithString("ca_config",
			mcp.Description("CA configuration file path"),
			mcp.Required(),
		),
	)
	s.AddTool(addProvisionerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewStepCAModule(client)

		// Get parameters
		name := request.GetString("provisioner_name", "")
		provisionerType := request.GetString("provisioner_type", "")
		caConfig := request.GetString("ca_config", "")

		// Add provisioner
		output, err := module.AddProvisioner(ctx, name, provisionerType, caConfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Step CA add provisioner failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Step CA revoke certificate tool
	revokeCertificateTool := mcp.NewTool("step_ca_revoke_certificate",
		mcp.WithDescription("Revoke certificate using real Step CA CLI"),
		mcp.WithString("cert_path",
			mcp.Description("Certificate file path"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Private key file path"),
			mcp.Required(),
		),
		mcp.WithString("ca_url",
			mcp.Description("Certificate authority URL"),
			mcp.Required(),
		),
		mcp.WithString("reason",
			mcp.Description("Revocation reason"),
		),
	)
	s.AddTool(revokeCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewStepCAModule(client)

		// Get parameters
		certPath := request.GetString("cert_path", "")
		keyPath := request.GetString("key_path", "")
		caURL := request.GetString("ca_url", "")
		reason := request.GetString("reason", "")

		// Revoke certificate
		output, err := module.RevokeCertificate(ctx, certPath, keyPath, caURL, reason)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Step CA revoke certificate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Step CA version tool
	versionTool := mcp.NewTool("step_ca_version",
		mcp.WithDescription("Get Step CA version information using real Step CA CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewStepCAModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Step CA get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}