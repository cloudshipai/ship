package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKyvernoTools adds Kyverno policy management MCP tool implementations using real CLI commands
func AddKyvernoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kyverno install with Helm tool
	installTool := mcp.NewTool("kyverno_install",
		mcp.WithDescription("Install Kyverno in Kubernetes cluster using Helm"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for Kyverno installation (default: kyverno)"),
		),
		mcp.WithString("version",
			mcp.Description("Kyverno Helm chart version to install"),
		),
		mcp.WithString("values_file",
			mcp.Description("Path to Helm values file for customization"),
		),
		mcp.WithBoolean("create_namespace",
			mcp.Description("Create namespace if it doesn't exist"),
		),
		mcp.WithBoolean("wait",
			mcp.Description("Wait for installation to complete"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// First add Kyverno Helm repository
		namespace := request.GetString("namespace", "kyverno")
		args := []string{"sh", "-c", "helm repo add kyverno https://kyverno.github.io/kyverno/ && helm repo update"}
		
		// Execute repo setup first, then install
		_, err := executeShipCommand(args)
		if err != nil {
			return nil, err
		}
		
		// Now install Kyverno
		installArgs := []string{"helm", "install", "kyverno", "kyverno/kyverno", "--namespace", namespace}
		
		if request.GetBool("create_namespace", false) {
			installArgs = append(installArgs, "--create-namespace")
		}
		if version := request.GetString("version", ""); version != "" {
			installArgs = append(installArgs, "--version", version)
		}
		if valuesFile := request.GetString("values_file", ""); valuesFile != "" {
			installArgs = append(installArgs, "--values", valuesFile)
		}
		if request.GetBool("wait", false) {
			installArgs = append(installArgs, "--wait")
		}
		
		return executeShipCommand(installArgs)
	})

	// Kyverno apply policy using CLI
	applyPolicyTool := mcp.NewTool("kyverno_apply",
		mcp.WithDescription("Test policies against resources using kyverno apply command"),
		mcp.WithString("policy",
			mcp.Description("Path to policy file or directory"),
			mcp.Required(),
		),
		mcp.WithString("resource",
			mcp.Description("Path to resource file or directory to test against"),
		),
		mcp.WithBoolean("cluster",
			mcp.Description("Apply policies against existing cluster resources"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace to apply policies in"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "yaml"),
		),
		mcp.WithBoolean("auto_generate_rules",
			mcp.Description("Enable auto-generation of rules"),
		),
	)
	s.AddTool(applyPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policy := request.GetString("policy", "")
		args := []string{"kyverno", "apply", policy}
		
		if resource := request.GetString("resource", ""); resource != "" {
			args = append(args, "--resource", resource)
		}
		if request.GetBool("cluster", false) {
			args = append(args, "--cluster")
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if request.GetBool("auto_generate_rules", false) {
			args = append(args, "--auto-gen-rules")
		}
		
		return executeShipCommand(args)
	})

	// Kyverno test policy using CLI
	testPolicyTool := mcp.NewTool("kyverno_test",
		mcp.WithDescription("Test Kyverno policies using kyverno test command"),
		mcp.WithString("policy",
			mcp.Description("Path to policy file or directory"),
			mcp.Required(),
		),
		mcp.WithString("resource",
			mcp.Description("Path to resource file or directory to test against"),
		),
		mcp.WithString("values",
			mcp.Description("Path to file containing values for policy variables"),
		),
		mcp.WithString("user_info",
			mcp.Description("Path to file containing user information for admission context"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "yaml"),
		),
		mcp.WithBoolean("audit",
			mcp.Description("Run policies in audit mode"),
		),
	)
	s.AddTool(testPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policy := request.GetString("policy", "")
		args := []string{"kyverno", "test", policy}
		
		if resource := request.GetString("resource", ""); resource != "" {
			args = append(args, "--resource", resource)
		}
		if values := request.GetString("values", ""); values != "" {
			args = append(args, "--values", values)
		}
		if userInfo := request.GetString("user_info", ""); userInfo != "" {
			args = append(args, "--user-info", userInfo)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if request.GetBool("audit", false) {
			args = append(args, "--audit")
		}
		
		return executeShipCommand(args)
	})

	// Kyverno create cluster role
	createClusterRoleTool := mcp.NewTool("kyverno_create_cluster_role",
		mcp.WithDescription("Create cluster role for Kyverno using kyverno create cluster-role"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json"),
		),
	)
	s.AddTool(createClusterRoleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kyverno", "create", "cluster-role"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		return executeShipCommand(args)
	})

	// Kyverno version
	versionTool := mcp.NewTool("kyverno_version",
		mcp.WithDescription("Get Kyverno CLI version using kyverno version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kyverno", "version"}
		return executeShipCommand(args)
	})

	// List Kyverno policies using kubectl
	listPoliciesTool := mcp.NewTool("kyverno_list_policies",
		mcp.WithDescription("List Kyverno policies using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to list policies from (empty for cluster policies)"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("table", "yaml", "json", "wide"),
		),
	)
	s.AddTool(listPoliciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		
		var args []string
		if namespace != "" {
			// List namespaced policies
			args = []string{"kubectl", "get", "policy", "-n", namespace}
		} else {
			// List cluster policies
			args = []string{"kubectl", "get", "clusterpolicy"}
		}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
	})

	// Apply Kyverno policy using kubectl
	applyPolicyFileTool := mcp.NewTool("kyverno_apply_policy_file",
		mcp.WithDescription("Apply Kyverno policy file to cluster using kubectl"),
		mcp.WithString("file",
			mcp.Description("Path to Kyverno policy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace for namespaced policies"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform dry run without applying"),
		),
	)
	s.AddTool(applyPolicyFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file := request.GetString("file", "")
		args := []string{"kubectl", "apply", "-f", file}
		
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run=client")
		}
		
		return executeShipCommand(args)
	})

	// Get PolicyReports using kubectl
	getPolicyReportsTool := mcp.NewTool("kyverno_get_policy_reports",
		mcp.WithDescription("Get Kyverno PolicyReports using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to get reports from (empty for ClusterPolicyReports)"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("table", "yaml", "json", "wide"),
		),
	)
	s.AddTool(getPolicyReportsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		
		var args []string
		if namespace != "" {
			// Get namespaced PolicyReports
			args = []string{"kubectl", "get", "policyreport", "-n", namespace}
		} else {
			// Get ClusterPolicyReports
			args = []string{"kubectl", "get", "clusterpolicyreport"}
		}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
	})

	// Check Kyverno status using kubectl
	statusTool := mcp.NewTool("kyverno_status",
		mcp.WithDescription("Check Kyverno installation status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kyverno namespace (default: kyverno)"),
		),
	)
	s.AddTool(statusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "kyverno")
		args := []string{"kubectl", "get", "pods", "-n", namespace, "-l", "app.kubernetes.io/name=kyverno"}
		
		return executeShipCommand(args)
	})
}