package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTrufflehogTools adds TruffleHog MCP tool implementations
func AddTrufflehogTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// TruffleHog scan directory tool
	scanDirectoryTool := mcp.NewTool("trufflehog_scan_directory",
		mcp.WithDescription("Scan directory for secrets using TruffleHog"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trufflehog", "filesystem"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// TruffleHog scan git repository tool
	scanGitRepoTool := mcp.NewTool("trufflehog_scan_git_repo",
		mcp.WithDescription("Scan git repository for secrets using TruffleHog"),
		mcp.WithString("repo_url",
			mcp.Description("Git repository URL to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanGitRepoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		args := []string{"security", "trufflehog", "git", repoURL}
		return executeShipCommand(args)
	})

	// TruffleHog scan GitHub repository tool
	scanGitHubTool := mcp.NewTool("trufflehog_scan_github",
		mcp.WithDescription("Scan GitHub repository for secrets using TruffleHog"),
		mcp.WithString("repo",
			mcp.Description("GitHub repository (owner/repo)"),
			mcp.Required(),
		),
		mcp.WithString("token",
			mcp.Description("GitHub personal access token"),
			mcp.Required(),
		),
	)
	s.AddTool(scanGitHubTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repo := request.GetString("repo", "")
		token := request.GetString("token", "")
		args := []string{"security", "trufflehog", "github", "--repo", repo, "--token", token}
		return executeShipCommand(args)
	})

	// TruffleHog scan GitHub organization tool
	scanGitHubOrgTool := mcp.NewTool("trufflehog_scan_github_org",
		mcp.WithDescription("Scan GitHub organization for secrets using TruffleHog"),
		mcp.WithString("org",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("token",
			mcp.Description("GitHub personal access token"),
			mcp.Required(),
		),
	)
	s.AddTool(scanGitHubOrgTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		org := request.GetString("org", "")
		token := request.GetString("token", "")
		args := []string{"security", "trufflehog", "github", "--org", org, "--token", token}
		return executeShipCommand(args)
	})

	// TruffleHog scan Docker image tool
	scanDockerImageTool := mcp.NewTool("trufflehog_scan_docker",
		mcp.WithDescription("Scan Docker image for secrets using TruffleHog"),
		mcp.WithString("image_name",
			mcp.Description("Docker image name to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanDockerImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"security", "trufflehog", "docker", "--image", imageName}
		return executeShipCommand(args)
	})

	// TruffleHog scan S3 bucket tool
	scanS3Tool := mcp.NewTool("trufflehog_scan_s3",
		mcp.WithDescription("Scan S3 bucket for secrets using TruffleHog"),
		mcp.WithString("bucket",
			mcp.Description("S3 bucket name to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanS3Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bucket := request.GetString("bucket", "")
		args := []string{"security", "trufflehog", "s3", "--bucket", bucket}
		return executeShipCommand(args)
	})

	// TruffleHog scan with verification tool
	scanWithVerificationTool := mcp.NewTool("trufflehog_scan_verified",
		mcp.WithDescription("Scan target with secret verification using TruffleHog"),
		mcp.WithString("target",
			mcp.Description("Target to scan"),
			mcp.Required(),
		),
		mcp.WithString("target_type",
			mcp.Description("Type of target (filesystem, git, github, docker, s3)"),
			mcp.Required(),
			mcp.Enum("filesystem", "git", "github", "docker", "s3"),
		),
	)
	s.AddTool(scanWithVerificationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		targetType := request.GetString("target_type", "")
		args := []string{"security", "trufflehog", targetType, target, "--verify"}
		return executeShipCommand(args)
	})

	// TruffleHog scan git with advanced options tool
	scanGitAdvancedTool := mcp.NewTool("trufflehog_scan_git_advanced",
		mcp.WithDescription("Scan git repository with advanced filtering options"),
		mcp.WithString("repo_url",
			mcp.Description("Git repository URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("branch",
			mcp.Description("Specific branch to scan"),
		),
		mcp.WithString("since_date",
			mcp.Description("Only scan commits since this date (YYYY-MM-DD)"),
		),
		mcp.WithString("until_date",
			mcp.Description("Only scan commits until this date (YYYY-MM-DD)"),
		),
		mcp.WithBoolean("only_verified",
			mcp.Description("Only return verified secrets"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "plain"),
		),
		mcp.WithString("exclude_paths",
			mcp.Description("Comma-separated paths to exclude"),
		),
	)
	s.AddTool(scanGitAdvancedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		args := []string{"security", "trufflehog", "git", repoURL}
		
		if branch := request.GetString("branch", ""); branch != "" {
			args = append(args, "--branch", branch)
		}
		if sinceDate := request.GetString("since_date", ""); sinceDate != "" {
			args = append(args, "--since", sinceDate)
		}
		if untilDate := request.GetString("until_date", ""); untilDate != "" {
			args = append(args, "--until", untilDate)
		}
		if request.GetBool("only_verified", false) {
			args = append(args, "--only-verified")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if excludePaths := request.GetString("exclude_paths", ""); excludePaths != "" {
			for _, path := range strings.Split(excludePaths, ",") {
				if strings.TrimSpace(path) != "" {
					args = append(args, "--exclude", strings.TrimSpace(path))
				}
			}
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog scan filesystem with exclusions tool
	scanFilesystemAdvancedTool := mcp.NewTool("trufflehog_scan_filesystem_advanced",
		mcp.WithDescription("Scan filesystem with advanced options and exclusions"),
		mcp.WithString("path",
			mcp.Description("Path to scan (file or directory)"),
			mcp.Required(),
		),
		mcp.WithBoolean("only_verified",
			mcp.Description("Only return verified secrets"),
		),
		mcp.WithString("exclude_paths",
			mcp.Description("Comma-separated paths to exclude from scanning"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "plain"),
		),
		mcp.WithString("max_depth",
			mcp.Description("Maximum directory depth to scan"),
		),
	)
	s.AddTool(scanFilesystemAdvancedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.GetString("path", "")
		args := []string{"security", "trufflehog", "filesystem", path}
		
		if request.GetBool("only_verified", false) {
			args = append(args, "--only-verified")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if maxDepth := request.GetString("max_depth", ""); maxDepth != "" {
			args = append(args, "--max-depth", maxDepth)
		}
		if excludePaths := request.GetString("exclude_paths", ""); excludePaths != "" {
			for _, excludePath := range strings.Split(excludePaths, ",") {
				if strings.TrimSpace(excludePath) != "" {
					args = append(args, "--exclude", strings.TrimSpace(excludePath))
				}
			}
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog scan docker with advanced options tool
	scanDockerAdvancedTool := mcp.NewTool("trufflehog_scan_docker_advanced",
		mcp.WithDescription("Scan Docker image with advanced verification options"),
		mcp.WithString("image",
			mcp.Description("Docker image name or ID to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("only_verified",
			mcp.Description("Only return verified secrets"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "plain"),
		),
		mcp.WithString("layers",
			mcp.Description("Specific layers to scan (comma-separated)"),
		),
	)
	s.AddTool(scanDockerAdvancedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		args := []string{"security", "trufflehog", "docker", image}
		
		if request.GetBool("only_verified", false) {
			args = append(args, "--only-verified")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if layers := request.GetString("layers", ""); layers != "" {
			args = append(args, "--layers", layers)
		}
		
		return executeShipCommand(args)
	})
}