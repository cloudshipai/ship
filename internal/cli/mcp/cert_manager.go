package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCertManagerTools adds cert-manager MCP tool implementations
func AddCertManagerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCertManagerToolsDirect(s)
}

// addCertManagerToolsDirect implements direct Dagger calls for cert-manager tools
func addCertManagerToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		version := request.GetString("version", "v1.18.2")
		dryRun := request.GetBool("dry_run", false)

		// Create cert-manager module and install
		certManagerModule := modules.NewCertManagerModule(client)
		result, err := certManagerModule.Install(ctx, version, dryRun)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cert-manager install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Cert-manager check installation tool
	checkInstallTool := mcp.NewTool("cert_manager_check_installation",
		mcp.WithDescription("Check cert-manager installation status"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to check (default: cert-manager)"),
		),
	)
	s.AddTool(checkInstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		namespace := request.GetString("namespace", "cert-manager")

		// Create cert-manager module and check installation
		certManagerModule := modules.NewCertManagerModule(client)
		result, err := certManagerModule.CheckInstallation(ctx, namespace, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("check installation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		name := request.GetString("name", "")
		fromCertFile := request.GetString("from_certificate_file", "")
		fetchCertificate := request.GetBool("fetch_certificate", false)
		timeout := request.GetString("timeout", "")

		// Create cert-manager module and create certificate request
		certManagerModule := modules.NewCertManagerModule(client)
		result, err := certManagerModule.CreateCertificateRequest(ctx, name, fromCertFile, fetchCertificate, timeout, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create certificate request failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		namespace := request.GetString("namespace", "")
		allNamespaces := request.GetBool("all_namespaces", false)

		// Create cert-manager module and list certificates
		certManagerModule := modules.NewCertManagerModule(client)
		result, err := certManagerModule.ListCertificates(ctx, namespace, allNamespaces, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list certificates failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		certName := request.GetString("cert_name", "")
		namespace := request.GetString("namespace", "")
		all := request.GetBool("all", false)

		// Create cert-manager module and renew certificate(s)
		certManagerModule := modules.NewCertManagerModule(client)
		var result string
		var renewErr error
		if all {
			result, renewErr = certManagerModule.RenewAllCertificates(ctx, namespace, "")
		} else {
			result, renewErr = certManagerModule.RenewCertificate(ctx, certName, namespace, "")
		}
		if renewErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("renew certificate failed: %v", renewErr)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		certName := request.GetString("certificate_name", "")
		namespace := request.GetString("namespace", "")

		// Create cert-manager module and check certificate status
		certManagerModule := modules.NewCertManagerModule(client)
		result, err := certManagerModule.CheckCertificate(ctx, certName, namespace, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("check certificate status failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Cert-manager get version tool
	getVersionTool := mcp.NewTool("cert_manager_get_version",
		mcp.WithDescription("Get cmctl version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create cert-manager module and get version
		certManagerModule := modules.NewCertManagerModule(client)
		result, err := certManagerModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get cert-manager version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}