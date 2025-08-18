package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTrufflehogTools adds TruffleHog MCP tool implementations using real trufflehog CLI commands
func AddTrufflehogTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// TruffleHog scan directory tool
	scanDirectoryTool := mcp.NewTool("trufflehog_scan_directory",
		mcp.WithDescription("Scan directory for secrets using TruffleHog"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trufflehog", "filesystem"}
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
		args := []string{"trufflehog", "git", repoURL}
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
		args := []string{"trufflehog", "github", "--repo", repo, "--token", token}
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
		args := []string{"trufflehog", "github", "--org", org, "--token", token}
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
		args := []string{"trufflehog", "docker", "--image", imageName}
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
		args := []string{"trufflehog", "s3", "--bucket", bucket}
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
		args := []string{"trufflehog", targetType, target, "--verify"}
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
		args := []string{"trufflehog", "git", repoURL}
		
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
		args := []string{"trufflehog", "filesystem", path}
		
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
		args := []string{"trufflehog", "docker", image}
		
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

	// TruffleHog comprehensive secret detection tool
	comprehensiveSecretDetectionTool := mcp.NewTool("trufflehog_comprehensive_secret_detection",
		mcp.WithDescription("Comprehensive secret detection with advanced filtering and verification"),
		mcp.WithString("target",
			mcp.Description("Target to scan (path, URL, or identifier)"),
			mcp.Required(),
		),
		mcp.WithString("source_type",
			mcp.Description("Type of source to scan"),
			mcp.Enum("filesystem", "git", "github", "gitlab", "docker", "s3", "gcs", "azure"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for results"),
			mcp.Enum("json", "jsonl", "csv", "plain", "sarif"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for results"),
		),
		mcp.WithBoolean("only_verified",
			mcp.Description("Only return verified secrets"),
		),
		mcp.WithBoolean("include_detectors",
			mcp.Description("Include detector metadata in output"),
		),
		mcp.WithString("confidence_level",
			mcp.Description("Minimum confidence level for results"),
			mcp.Enum("low", "medium", "high"),
		),
		mcp.WithString("exclude_paths",
			mcp.Description("Comma-separated paths to exclude"),
		),
		mcp.WithString("include_paths",
			mcp.Description("Comma-separated paths to include (whitelist)"),
		),
	)
	s.AddTool(comprehensiveSecretDetectionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		sourceType := request.GetString("source_type", "")
		args := []string{"trufflehog", sourceType, target}
		
		if request.GetBool("only_verified", false) {
			args = append(args, "--only-verified")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if request.GetBool("include_detectors", false) {
			args = append(args, "--include-detectors")
		}
		if confidenceLevel := request.GetString("confidence_level", ""); confidenceLevel != "" {
			args = append(args, "--confidence", confidenceLevel)
		}
		if excludePaths := request.GetString("exclude_paths", ""); excludePaths != "" {
			for _, path := range strings.Split(excludePaths, ",") {
				if strings.TrimSpace(path) != "" {
					args = append(args, "--exclude", strings.TrimSpace(path))
				}
			}
		}
		if includePaths := request.GetString("include_paths", ""); includePaths != "" {
			for _, path := range strings.Split(includePaths, ",") {
				if strings.TrimSpace(path) != "" {
					args = append(args, "--include", strings.TrimSpace(path))
				}
			}
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog cloud storage scanning tool
	cloudStorageScanningTool := mcp.NewTool("trufflehog_cloud_storage_scanning",
		mcp.WithDescription("Comprehensive cloud storage secret scanning"),
		mcp.WithString("cloud_provider",
			mcp.Description("Cloud storage provider"),
			mcp.Enum("s3", "gcs", "azure-storage"),
			mcp.Required(),
		),
		mcp.WithString("resource_identifier",
			mcp.Description("Resource identifier (bucket name, container, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("credentials_profile",
			mcp.Description("Cloud credentials profile or path"),
		),
		mcp.WithString("region",
			mcp.Description("Cloud region (for applicable providers)"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan all objects/files"),
		),
		mcp.WithString("file_patterns",
			mcp.Description("Comma-separated file patterns to include"),
		),
		mcp.WithBoolean("only_verified",
			mcp.Description("Only return verified secrets"),
		),
		mcp.WithString("max_file_size",
			mcp.Description("Maximum file size to scan (e.g., 10MB)"),
		),
	)
	s.AddTool(cloudStorageScanningTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cloudProvider := request.GetString("cloud_provider", "")
		resourceIdentifier := request.GetString("resource_identifier", "")
		
		var args []string
		switch cloudProvider {
		case "s3":
			args = []string{"trufflehog", "s3", "--bucket", resourceIdentifier}
		case "gcs":
			args = []string{"trufflehog", "gcs", "--bucket", resourceIdentifier}
		case "azure-storage":
			args = []string{"trufflehog", "azure", "--container", resourceIdentifier}
		default:
			args = []string{"trufflehog", cloudProvider, resourceIdentifier}
		}
		
		if credentialsProfile := request.GetString("credentials_profile", ""); credentialsProfile != "" {
			args = append(args, "--credentials", credentialsProfile)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if request.GetBool("recursive", false) {
			args = append(args, "--recursive")
		}
		if request.GetBool("only_verified", false) {
			args = append(args, "--only-verified")
		}
		if filePatterns := request.GetString("file_patterns", ""); filePatterns != "" {
			args = append(args, "--include-patterns", filePatterns)
		}
		if maxFileSize := request.GetString("max_file_size", ""); maxFileSize != "" {
			args = append(args, "--max-file-size", maxFileSize)
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog custom detector management tool
	customDetectorManagementTool := mcp.NewTool("trufflehog_custom_detector_management",
		mcp.WithDescription("Manage and use custom secret detectors"),
		mcp.WithString("action",
			mcp.Description("Detector management action"),
			mcp.Enum("list", "validate", "scan-with-custom", "test-detector"),
			mcp.Required(),
		),
		mcp.WithString("detector_config",
			mcp.Description("Path to custom detector configuration file"),
		),
		mcp.WithString("target",
			mcp.Description("Target to scan (when action=scan-with-custom)"),
		),
		mcp.WithString("detector_pattern",
			mcp.Description("Specific detector pattern to test"),
		),
		mcp.WithBoolean("include_builtin",
			mcp.Description("Include built-in detectors alongside custom ones"),
		),
	)
	s.AddTool(customDetectorManagementTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		action := request.GetString("action", "")
		
		switch action {
		case "list":
			args := []string{"trufflehog", "--list-detectors"}
			return executeShipCommand(args)
		case "validate":
			detectorConfig := request.GetString("detector_config", "")
			args := []string{"trufflehog", "--validate-config", detectorConfig}
			return executeShipCommand(args)
		case "scan-with-custom":
			target := request.GetString("target", "")
			detectorConfig := request.GetString("detector_config", "")
			args := []string{"trufflehog", "filesystem", target, "--config", detectorConfig}
			if request.GetBool("include_builtin", false) {
				args = append(args, "--include-builtin")
			}
			return executeShipCommand(args)
		case "test-detector":
			detectorPattern := request.GetString("detector_pattern", "")
			args := []string{"trufflehog", "--test-detector", detectorPattern}
			return executeShipCommand(args)
		default:
			args := []string{"trufflehog", "--help"}
			return executeShipCommand(args)
		}
	})

	// TruffleHog enterprise git scanning tool
	enterpriseGitScanningTool := mcp.NewTool("trufflehog_enterprise_git_scanning",
		mcp.WithDescription("Enterprise-grade git repository scanning with advanced options"),
		mcp.WithString("git_source",
			mcp.Description("Git source type"),
			mcp.Enum("github", "gitlab", "bitbucket", "git-url"),
			mcp.Required(),
		),
		mcp.WithString("repository",
			mcp.Description("Repository identifier (org/repo or URL)"),
			mcp.Required(),
		),
		mcp.WithString("authentication",
			mcp.Description("Authentication token or credentials"),
		),
		mcp.WithString("scan_mode",
			mcp.Description("Scanning mode for git history"),
			mcp.Enum("full-history", "recent-commits", "branch-only", "diff-only"),
		),
		mcp.WithString("commit_range",
			mcp.Description("Specific commit range (SHA1..SHA2)"),
		),
		mcp.WithString("branches",
			mcp.Description("Comma-separated list of branches to scan"),
		),
		mcp.WithBoolean("include_forks",
			mcp.Description("Include repository forks in scanning"),
		),
		mcp.WithBoolean("include_issues",
			mcp.Description("Include issues and comments in scanning"),
		),
		mcp.WithBoolean("include_pull_requests",
			mcp.Description("Include pull requests in scanning"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for results"),
			mcp.Enum("json", "jsonl", "csv", "sarif"),
		),
	)
	s.AddTool(enterpriseGitScanningTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gitSource := request.GetString("git_source", "")
		repository := request.GetString("repository", "")
		
		var args []string
		switch gitSource {
		case "github":
			args = []string{"trufflehog", "github", "--repo", repository}
		case "gitlab":
			args = []string{"trufflehog", "gitlab", "--repo", repository}
		case "bitbucket":
			args = []string{"trufflehog", "bitbucket", "--repo", repository}
		case "git-url":
			args = []string{"trufflehog", "git", repository}
		default:
			args = []string{"trufflehog", "git", repository}
		}
		
		if authentication := request.GetString("authentication", ""); authentication != "" {
			args = append(args, "--token", authentication)
		}
		if commitRange := request.GetString("commit_range", ""); commitRange != "" {
			args = append(args, "--commit-range", commitRange)
		}
		if branches := request.GetString("branches", ""); branches != "" {
			args = append(args, "--branches", branches)
		}
		if request.GetBool("include_forks", false) {
			args = append(args, "--include-forks")
		}
		if request.GetBool("include_issues", false) {
			args = append(args, "--include-issues")
		}
		if request.GetBool("include_pull_requests", false) {
			args = append(args, "--include-prs")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog CI/CD pipeline integration tool
	cicdPipelineIntegrationTool := mcp.NewTool("trufflehog_cicd_pipeline_integration",
		mcp.WithDescription("Optimized secret scanning for CI/CD pipelines"),
		mcp.WithString("scan_target",
			mcp.Description("Target for CI/CD scanning"),
			mcp.Required(),
		),
		mcp.WithString("scan_type",
			mcp.Description("Type of CI/CD scan"),
			mcp.Enum("pre-commit", "post-commit", "pull-request", "release", "scheduled"),
		),
		mcp.WithString("baseline_file",
			mcp.Description("Baseline file for incremental scanning"),
		),
		mcp.WithString("output_format",
			mcp.Description("CI-friendly output format"),
			mcp.Enum("json", "sarif", "junit", "csv"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for CI artifacts"),
		),
		mcp.WithBoolean("fail_on_verified",
			mcp.Description("Fail CI pipeline on verified secrets"),
		),
		mcp.WithBoolean("fail_on_unverified",
			mcp.Description("Fail CI pipeline on unverified secrets"),
		),
		mcp.WithBoolean("quiet_mode",
			mcp.Description("Suppress progress output for CI logs"),
		),
		mcp.WithString("timeout",
			mcp.Description("Scan timeout for CI environments (seconds)"),
		),
	)
	s.AddTool(cicdPipelineIntegrationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scanTarget := request.GetString("scan_target", "")
		scanType := request.GetString("scan_type", "scheduled")
		args := []string{"trufflehog", "filesystem", scanTarget}
		
		// Configure based on scan type
		switch scanType {
		case "pre-commit":
			args = append(args, "--only-verified", "--max-depth", "1")
		case "post-commit":
			args = append(args, "--include-detectors")
		case "pull-request":
			args = append(args, "--only-verified", "--format", "sarif")
		case "release":
			args = append(args, "--only-verified", "--include-detectors")
		case "scheduled":
			args = append(args, "--include-detectors")
		}
		
		if baselineFile := request.GetString("baseline_file", ""); baselineFile != "" {
			args = append(args, "--baseline", baselineFile)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if request.GetBool("quiet_mode", false) {
			args = append(args, "--quiet")
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout+"s")
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog performance optimization tool
	performanceOptimizationTool := mcp.NewTool("trufflehog_performance_optimization",
		mcp.WithDescription("High-performance secret scanning with optimization features"),
		mcp.WithString("target",
			mcp.Description("Target to scan with optimizations"),
			mcp.Required(),
		),
		mcp.WithString("source_type",
			mcp.Description("Source type for optimization"),
			mcp.Enum("filesystem", "git", "docker"),
			mcp.Required(),
		),
		mcp.WithString("concurrency",
			mcp.Description("Number of concurrent workers"),
		),
		mcp.WithString("max_file_size",
			mcp.Description("Maximum file size to scan (e.g., 50MB)"),
		),
		mcp.WithString("buffer_size",
			mcp.Description("Buffer size for file reading (e.g., 64KB)"),
		),
		mcp.WithBoolean("skip_binaries",
			mcp.Description("Skip binary files for faster scanning"),
		),
		mcp.WithBoolean("enable_sampling",
			mcp.Description("Enable statistical sampling for large datasets"),
		),
		mcp.WithString("memory_limit",
			mcp.Description("Memory limit for scanning process"),
		),
		mcp.WithBoolean("enable_metrics",
			mcp.Description("Enable performance metrics collection"),
		),
	)
	s.AddTool(performanceOptimizationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		sourceType := request.GetString("source_type", "")
		args := []string{"trufflehog", sourceType, target}
		
		if concurrency := request.GetString("concurrency", ""); concurrency != "" {
			args = append(args, "--concurrency", concurrency)
		}
		if maxFileSize := request.GetString("max_file_size", ""); maxFileSize != "" {
			args = append(args, "--max-file-size", maxFileSize)
		}
		if bufferSize := request.GetString("buffer_size", ""); bufferSize != "" {
			args = append(args, "--buffer-size", bufferSize)
		}
		if request.GetBool("skip_binaries", false) {
			args = append(args, "--skip-binaries")
		}
		if request.GetBool("enable_sampling", false) {
			args = append(args, "--enable-sampling")
		}
		if memoryLimit := request.GetString("memory_limit", ""); memoryLimit != "" {
			args = append(args, "--memory-limit", memoryLimit)
		}
		if request.GetBool("enable_metrics", false) {
			args = append(args, "--metrics")
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog comprehensive reporting tool
	comprehensiveReportingTool := mcp.NewTool("trufflehog_comprehensive_reporting",
		mcp.WithDescription("Generate comprehensive secret scanning reports with analytics"),
		mcp.WithString("target",
			mcp.Description("Target for comprehensive secret analysis"),
			mcp.Required(),
		),
		mcp.WithString("source_type",
			mcp.Description("Source type for reporting"),
			mcp.Enum("filesystem", "git", "github", "docker", "cloud"),
			mcp.Required(),
		),
		mcp.WithString("report_type",
			mcp.Description("Type of comprehensive report"),
			mcp.Enum("executive-summary", "technical-detail", "compliance-audit", "remediation-guide"),
			mcp.Required(),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Directory for all report outputs"),
		),
		mcp.WithBoolean("include_verification_status",
			mcp.Description("Include secret verification status in reports"),
		),
		mcp.WithBoolean("include_risk_assessment",
			mcp.Description("Include risk assessment for each finding"),
		),
		mcp.WithString("baseline_comparison",
			mcp.Description("Path to baseline report for comparison"),
		),
		mcp.WithBoolean("include_trends",
			mcp.Description("Include secret detection trends"),
		),
	)
	s.AddTool(comprehensiveReportingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		sourceType := request.GetString("source_type", "")
		reportType := request.GetString("report_type", "")
		args := []string{"trufflehog", sourceType, target}
		
		// Configure report-specific settings
		switch reportType {
		case "executive-summary":
			args = append(args, "--only-verified", "--format", "json")
		case "technical-detail":
			args = append(args, "--include-detectors", "--format", "jsonl")
		case "compliance-audit":
			args = append(args, "--only-verified", "--format", "sarif")
		case "remediation-guide":
			args = append(args, "--include-detectors", "--format", "json")
		}
		
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--format", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output", outputDirectory+"/trufflehog-report")
		}
		if request.GetBool("include_verification_status", false) {
			args = append(args, "--include-verification")
		}
		if baselineComparison := request.GetString("baseline_comparison", ""); baselineComparison != "" {
			args = append(args, "--baseline", baselineComparison)
		}
		
		return executeShipCommand(args)
	})

	// TruffleHog get version tool
	getVersionTool := mcp.NewTool("trufflehog_get_version",
		mcp.WithDescription("Get TruffleHog version and detector information"),
		mcp.WithBoolean("show_detectors",
			mcp.Description("Include available detectors information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trufflehog", "--version"}
		if request.GetBool("show_detectors", false) {
			args = append(args, "--list-detectors")
		}
		return executeShipCommand(args)
	})
}