package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCertManagerTools adds cert-manager MCP tool implementations
func AddCertManagerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cert-manager install tool
	installTool := mcp.NewTool("cert_manager_install",
		mcp.WithDescription("Install cert-manager for automated certificate management"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for cert-manager installation"),
		),
		mcp.WithString("version",
			mcp.Description("Cert-manager version to install"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "cert-manager", "install"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if version := request.GetString("version", ""); version != "" {
			args = append(args, "--version", version)
		}
		return executeShipCommand(args)
	})

	// Cert-manager create issuer tool
	createIssuerTool := mcp.NewTool("cert_manager_create_issuer",
		mcp.WithDescription("Create certificate issuer (Let's Encrypt, CA, etc.)"),
		mcp.WithString("issuer_name",
			mcp.Description("Name of the certificate issuer"),
			mcp.Required(),
		),
		mcp.WithString("issuer_type",
			mcp.Description("Type of issuer (letsencrypt, ca, selfsigned)"),
			mcp.Required(),
		),
		mcp.WithString("email",
			mcp.Description("Email address for Let's Encrypt notifications"),
		),
	)
	s.AddTool(createIssuerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		issuerName := request.GetString("issuer_name", "")
		issuerType := request.GetString("issuer_type", "")
		args := []string{"kubernetes", "cert-manager", "create-issuer", issuerName, "--type", issuerType}
		if email := request.GetString("email", ""); email != "" {
			args = append(args, "--email", email)
		}
		return executeShipCommand(args)
	})

	// Cert-manager request certificate tool
	requestCertificateTool := mcp.NewTool("cert_manager_request_certificate",
		mcp.WithDescription("Request SSL/TLS certificate for domain"),
		mcp.WithString("cert_name",
			mcp.Description("Name of the certificate"),
			mcp.Required(),
		),
		mcp.WithString("domain",
			mcp.Description("Domain name for the certificate"),
			mcp.Required(),
		),
		mcp.WithString("issuer",
			mcp.Description("Certificate issuer to use"),
			mcp.Required(),
		),
	)
	s.AddTool(requestCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certName := request.GetString("cert_name", "")
		domain := request.GetString("domain", "")
		issuer := request.GetString("issuer", "")
		args := []string{"kubernetes", "cert-manager", "request-cert", certName, "--domain", domain, "--issuer", issuer}
		return executeShipCommand(args)
	})

	// Cert-manager list certificates tool
	listCertificatesTool := mcp.NewTool("cert_manager_list_certificates",
		mcp.WithDescription("List managed certificates and their status"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to list certificates from"),
		),
	)
	s.AddTool(listCertificatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "cert-manager", "list-certs"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Cert-manager renew certificate tool
	renewCertificateTool := mcp.NewTool("cert_manager_renew_certificate",
		mcp.WithDescription("Force renewal of existing certificate"),
		mcp.WithString("cert_name",
			mcp.Description("Name of the certificate to renew"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(renewCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certName := request.GetString("cert_name", "")
		args := []string{"kubernetes", "cert-manager", "renew", certName}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Cert-manager get version tool
	getVersionTool := mcp.NewTool("cert_manager_get_version",
		mcp.WithDescription("Get cert-manager version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "cert-manager", "--version"}
		return executeShipCommand(args)
	})
}