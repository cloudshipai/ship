package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddStepCATools adds Step CA (Certificate Authority) MCP tool implementations
func AddStepCATools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Step CA init tool
	initTool := mcp.NewTool("step_ca_init",
		mcp.WithDescription("Initialize Step CA certificate authority"),
		mcp.WithString("ca_name",
			mcp.Description("Name of the certificate authority"),
			mcp.Required(),
		),
		mcp.WithString("dns_name",
			mcp.Description("DNS name for the CA"),
			mcp.Required(),
		),
		mcp.WithString("provisioner",
			mcp.Description("Default provisioner name"),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		caName := request.GetString("ca_name", "")
		dnsName := request.GetString("dns_name", "")
		args := []string{"security", "step-ca", "init", caName, dnsName}
		if provisioner := request.GetString("provisioner", ""); provisioner != "" {
			args = append(args, "--provisioner", provisioner)
		}
		return executeShipCommand(args)
	})

	// Step CA create certificate tool
	createCertificateTool := mcp.NewTool("step_ca_create_certificate",
		mcp.WithDescription("Create certificate using Step CA"),
		mcp.WithString("subject",
			mcp.Description("Certificate subject (CN)"),
			mcp.Required(),
		),
		mcp.WithString("cert_file",
			mcp.Description("Output certificate file path"),
			mcp.Required(),
		),
		mcp.WithString("key_file",
			mcp.Description("Output private key file path"),
			mcp.Required(),
		),
	)
	s.AddTool(createCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		subject := request.GetString("subject", "")
		certFile := request.GetString("cert_file", "")
		keyFile := request.GetString("key_file", "")
		args := []string{"security", "step-ca", "certificate", subject, certFile, keyFile}
		return executeShipCommand(args)
	})

	// Step CA create provisioner tool
	createProvisionerTool := mcp.NewTool("step_ca_create_provisioner",
		mcp.WithDescription("Create new provisioner for Step CA"),
		mcp.WithString("provisioner_name",
			mcp.Description("Name of the provisioner"),
			mcp.Required(),
		),
		mcp.WithString("provisioner_type",
			mcp.Description("Type of provisioner (jwk, oidc, x5c, acme)"),
			mcp.Required(),
		),
	)
	s.AddTool(createProvisionerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provisionerName := request.GetString("provisioner_name", "")
		provisionerType := request.GetString("provisioner_type", "")
		args := []string{"security", "step-ca", "provisioner", "add", provisionerName, "--type", provisionerType}
		return executeShipCommand(args)
	})

	// Step CA revoke certificate tool
	revokeCertificateTool := mcp.NewTool("step_ca_revoke_certificate",
		mcp.WithDescription("Revoke certificate using Step CA"),
		mcp.WithString("cert_file",
			mcp.Description("Path to certificate file to revoke"),
			mcp.Required(),
		),
		mcp.WithString("reason",
			mcp.Description("Revocation reason"),
		),
	)
	s.AddTool(revokeCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certFile := request.GetString("cert_file", "")
		args := []string{"security", "step-ca", "revoke", certFile}
		if reason := request.GetString("reason", ""); reason != "" {
			args = append(args, "--reason", reason)
		}
		return executeShipCommand(args)
	})

	// Step CA renew certificate tool
	renewCertificateTool := mcp.NewTool("step_ca_renew_certificate",
		mcp.WithDescription("Renew certificate using Step CA"),
		mcp.WithString("cert_file",
			mcp.Description("Path to certificate file to renew"),
			mcp.Required(),
		),
		mcp.WithString("key_file",
			mcp.Description("Path to private key file"),
			mcp.Required(),
		),
	)
	s.AddTool(renewCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certFile := request.GetString("cert_file", "")
		keyFile := request.GetString("key_file", "")
		args := []string{"security", "step-ca", "renew", certFile, keyFile}
		return executeShipCommand(args)
	})

	// Step CA get version tool
	getVersionTool := mcp.NewTool("step_ca_get_version",
		mcp.WithDescription("Get Step CA version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "step-ca", "--version"}
		return executeShipCommand(args)
	})
}