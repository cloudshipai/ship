package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCertManagerTools adds cert-manager MCP tool implementations
func AddCertManagerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cert-manager install tool (using kubectl)
	installTool := mcp.NewTool("cert_manager_install",
		mcp.WithDescription("Install cert-manager using kubectl apply"),
		mcp.WithString("version",
			mcp.Description("Cert-manager version to install (e.g., v1.18.2)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform dry run without actually installing"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version := request.GetString("version", "v1.18.2")
		manifestUrl := fmt.Sprintf("https://github.com/cert-manager/cert-manager/releases/download/%s/cert-manager.yaml", version)
		args := []string{"kubectl", "apply", "-f", manifestUrl}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run=client")
		}
		return executeShipCommand(args)
	})

	// Cert-manager check installation tool
	checkInstallTool := mcp.NewTool("cert_manager_check_installation",
		mcp.WithDescription("Check cert-manager installation status"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to check (default: cert-manager)"),
		),
	)
	s.AddTool(checkInstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "cert-manager")
		args := []string{"kubectl", "get", "pods", "-n", namespace}
		return executeShipCommand(args)
	})

	// Cert-manager create certificate request tool
	createCertRequestTool := mcp.NewTool("cert_manager_create_certificate_request",
		mcp.WithDescription("Create a CertificateRequest using cmctl"),
		mcp.WithString("name",
			mcp.Description("Name of the CertificateRequest"),
			mcp.Required(),
		),
		mcp.WithString("from_certificate_file",
			mcp.Description("Path to certificate file to create request from"),
		),
		mcp.WithBoolean("fetch_certificate",
			mcp.Description("Fetch the certificate once issued"),
		),
		mcp.WithString("timeout",
			mcp.Description("Timeout for the operation (e.g., 20m)"),
		),
	)
	s.AddTool(createCertRequestTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"cmctl", "create", "certificaterequest", name}
		if fromCertFile := request.GetString("from_certificate_file", ""); fromCertFile != "" {
			args = append(args, "--from-certificate-file", fromCertFile)
		}
		if request.GetBool("fetch_certificate", false) {
			args = append(args, "--fetch-certificate")
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout)
		}
		return executeShipCommand(args)
	})

	// Cert-manager list certificates tool
	listCertificatesTool := mcp.NewTool("cert_manager_list_certificates",
		mcp.WithDescription("List certificates using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to list certificates from"),
		),
		mcp.WithBoolean("all_namespaces",
			mcp.Description("List certificates from all namespaces"),
		),
	)
	s.AddTool(listCertificatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "certificates"}
		if request.GetBool("all_namespaces", false) {
			args = append(args, "--all-namespaces")
		} else if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		return executeShipCommand(args)
	})

	// Cert-manager renew certificate tool
	renewCertificateTool := mcp.NewTool("cert_manager_renew_certificate",
		mcp.WithDescription("Mark certificate for manual renewal using cmctl"),
		mcp.WithString("cert_name",
			mcp.Description("Name of the certificate to renew"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithBoolean("all",
			mcp.Description("Renew all certificates in namespace"),
		),
	)
	s.AddTool(renewCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cmctl", "renew"}
		if request.GetBool("all", false) {
			args = append(args, "--all")
		} else {
			certName := request.GetString("cert_name", "")
			args = append(args, certName)
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		return executeShipCommand(args)
	})

	// Cert-manager status tool
	statusTool := mcp.NewTool("cert_manager_status",
		mcp.WithDescription("Get status of certificate using cmctl"),
		mcp.WithString("certificate_name",
			mcp.Description("Name of the certificate to check"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(statusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certName := request.GetString("certificate_name", "")
		args := []string{"cmctl", "status", "certificate", certName}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		return executeShipCommand(args)
	})

	// Cert-manager get version tool
	getVersionTool := mcp.NewTool("cert_manager_get_version",
		mcp.WithDescription("Get cmctl version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cmctl", "version"}
		return executeShipCommand(args)
	})
}