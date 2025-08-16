package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSemgrepTools adds Semgrep (static analysis) MCP tool implementations
func AddSemgrepTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Semgrep scan directory tool
	scanDirectoryTool := mcp.NewTool("semgrep_scan_directory",
		mcp.WithDescription("Scan directory for security issues using Semgrep"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Enum("INFO", "WARNING", "ERROR"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif", "junit-xml"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "semgrep", directory}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Semgrep scan with ruleset tool
	scanWithRulesetTool := mcp.NewTool("semgrep_scan_with_ruleset",
		mcp.WithDescription("Scan with specific Semgrep ruleset"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("ruleset",
			mcp.Description("Semgrep ruleset to use"),
			mcp.Required(),
			mcp.Enum("auto", "p/security-audit", "p/owasp-top-ten", "p/cwe-top-25", "p/python", "p/javascript", "p/go"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif", "junit-xml"),
		),
	)
	s.AddTool(scanWithRulesetTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		ruleset := request.GetString("ruleset", "")
		args := []string{"security", "semgrep", directory, "--config", ruleset}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Semgrep scan specific file tool
	scanFileTool := mcp.NewTool("semgrep_scan_file",
		mcp.WithDescription("Scan specific file using Semgrep"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to scan"),
			mcp.Required(),
		),
		mcp.WithString("ruleset",
			mcp.Description("Semgrep ruleset to use"),
			mcp.Enum("auto", "p/security-audit", "p/owasp-top-ten", "p/cwe-top-25"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif", "junit-xml"),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "semgrep", "--file", filePath}
		if ruleset := request.GetString("ruleset", ""); ruleset != "" {
			args = append(args, "--config", ruleset)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Semgrep scan with custom rules tool
	scanWithCustomRulesTool := mcp.NewTool("semgrep_scan_with_custom_rules",
		mcp.WithDescription("Scan using custom Semgrep rules file"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("rules_file",
			mcp.Description("Path to custom rules file"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif", "junit-xml"),
		),
	)
	s.AddTool(scanWithCustomRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		rulesFile := request.GetString("rules_file", "")
		args := []string{"security", "semgrep", directory, "--config", rulesFile}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Semgrep comprehensive scan tool
	comprehensiveScanTool := mcp.NewTool("semgrep_comprehensive_scan",
		mcp.WithDescription("Run comprehensive Semgrep SAST scan with advanced options"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("rules",
			mcp.Description("Comma-separated rules (p/security-audit,p/secrets,p/owasp-top-ten)"),
		),
		mcp.WithString("exclude_patterns",
			mcp.Description("Comma-separated exclude patterns (node_modules,.git,vendor)"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "sarif", "text", "gitlab-sast", "gitlab-secrets"),
		),
		mcp.WithString("severity_filter",
			mcp.Description("Comma-separated severity levels to include"),
			mcp.Enum("ERROR", "WARNING", "INFO"),
		),
		mcp.WithString("timeout",
			mcp.Description("Timeout in seconds for the scan"),
		),
		mcp.WithBoolean("fail_on_findings",
			mcp.Description("Return error if findings are detected"),
		),
	)
	s.AddTool(comprehensiveScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "semgrep", directory}
		
		if rules := request.GetString("rules", ""); rules != "" {
			for _, rule := range strings.Split(rules, ",") {
				if strings.TrimSpace(rule) != "" {
					args = append(args, "--config", strings.TrimSpace(rule))
				}
			}
		}
		
		if excludePatterns := request.GetString("exclude_patterns", ""); excludePatterns != "" {
			for _, pattern := range strings.Split(excludePatterns, ",") {
				if strings.TrimSpace(pattern) != "" {
					args = append(args, "--exclude", strings.TrimSpace(pattern))
				}
			}
		}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		
		if severity := request.GetString("severity_filter", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		
		if timeoutStr := request.GetString("timeout", ""); timeoutStr != "" {
			args = append(args, "--timeout", timeoutStr)
		}
		
		if request.GetBool("fail_on_findings", false) {
			args = append(args, "--error")
		}
		
		return executeShipCommand(args)
	})

	// Semgrep secrets scanning tool
	scanSecretsTool := mcp.NewTool("semgrep_scan_secrets",
		mcp.WithDescription("Specialized secrets scanning using Semgrep"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for secrets"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "sarif", "gitlab-secrets"),
		),
		mcp.WithString("exclude_patterns",
			mcp.Description("Comma-separated exclude patterns"),
		),
	)
	s.AddTool(scanSecretsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "semgrep", directory, "--config", "p/secrets"}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		
		if excludePatterns := request.GetString("exclude_patterns", ""); excludePatterns != "" {
			for _, pattern := range strings.Split(excludePatterns, ",") {
				if strings.TrimSpace(pattern) != "" {
					args = append(args, "--exclude", strings.TrimSpace(pattern))
				}
			}
		}
		
		return executeShipCommand(args)
	})

	// Semgrep OWASP Top 10 scan tool
	scanOWASPTool := mcp.NewTool("semgrep_scan_owasp_top10",
		mcp.WithDescription("Scan for OWASP Top 10 vulnerabilities using Semgrep"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "sarif", "text"),
		),
		mcp.WithString("language_focus",
			mcp.Description("Focus on specific language"),
			mcp.Enum("python", "javascript", "java", "go", "php", "ruby"),
		),
	)
	s.AddTool(scanOWASPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "semgrep", directory, "--config", "p/owasp-top-ten"}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		
		if language := request.GetString("language_focus", ""); language != "" {
			args = append(args, "--config", "p/"+language)
		}
		
		return executeShipCommand(args)
	})

	// Semgrep get version tool
	getVersionTool := mcp.NewTool("semgrep_get_version",
		mcp.WithDescription("Get Semgrep version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "semgrep", "--version"}
		return executeShipCommand(args)
	})
}