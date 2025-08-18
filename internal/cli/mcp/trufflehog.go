package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTrufflehogTools adds TruffleHog MCP tool implementations using direct Dagger calls
func AddTrufflehogTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTrufflehogToolsDirect(s)
}

// addTrufflehogToolsDirect adds TruffleHog tools using direct Dagger module calls
func addTrufflehogToolsDirect(s *server.MCPServer) {
	// TruffleHog scan directory tool
	scanDirectoryTool := mcp.NewTool("trufflehog_scan_directory",
		mcp.WithDescription("Scan directory for secrets using TruffleHog"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		dir := request.GetString("directory", ".")

		// Scan directory
		output, err := module.ScanDirectory(ctx, dir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog directory scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		repoURL := request.GetString("repo_url", "")
		if repoURL == "" {
			return mcp.NewToolResultError("repo_url is required"), nil
		}

		// Scan git repo
		output, err := module.ScanGitRepo(ctx, repoURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog git repo scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		repo := request.GetString("repo", "")
		token := request.GetString("token", "")

		if repo == "" {
			return mcp.NewToolResultError("repo is required"), nil
		}
		if token == "" {
			return mcp.NewToolResultError("token is required"), nil
		}

		// Scan GitHub repo
		output, err := module.ScanGitHub(ctx, repo, token)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog GitHub scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		org := request.GetString("org", "")
		token := request.GetString("token", "")

		if org == "" {
			return mcp.NewToolResultError("org is required"), nil
		}
		if token == "" {
			return mcp.NewToolResultError("token is required"), nil
		}

		// Scan GitHub org
		output, err := module.ScanGitHubOrg(ctx, org, token)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog GitHub org scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		imageName := request.GetString("image_name", "")
		if imageName == "" {
			return mcp.NewToolResultError("image_name is required"), nil
		}

		// Scan Docker image
		output, err := module.ScanDockerImage(ctx, imageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog Docker scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		bucket := request.GetString("bucket", "")
		if bucket == "" {
			return mcp.NewToolResultError("bucket is required"), nil
		}

		// Scan S3 bucket
		output, err := module.ScanS3(ctx, bucket)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog S3 scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		target := request.GetString("target", "")
		targetType := request.GetString("target_type", "")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if targetType == "" {
			return mcp.NewToolResultError("target_type is required"), nil
		}

		// Scan with verification
		output, err := module.ScanWithVerification(ctx, target, targetType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog verification scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		repoURL := request.GetString("repo_url", "")
		if repoURL == "" {
			return mcp.NewToolResultError("repo_url is required"), nil
		}

		branch := request.GetString("branch", "")
		sinceDate := request.GetString("since_date", "")
		untilDate := request.GetString("until_date", "")
		onlyVerified := request.GetBool("only_verified", false)
		outputFormat := request.GetString("output_format", "")
		excludePathsStr := request.GetString("exclude_paths", "")

		// Parse exclude paths
		var excludePaths []string
		if excludePathsStr != "" {
			for _, path := range strings.Split(excludePathsStr, ",") {
				if strings.TrimSpace(path) != "" {
					excludePaths = append(excludePaths, strings.TrimSpace(path))
				}
			}
		}

		// Scan git advanced
		output, err := module.ScanGitAdvanced(ctx, repoURL, branch, sinceDate, untilDate, onlyVerified, outputFormat, excludePaths)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog advanced git scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		onlyVerified := request.GetBool("only_verified", false)
		outputFormat := request.GetString("output_format", "")
		maxDepth := request.GetString("max_depth", "")
		excludePathsStr := request.GetString("exclude_paths", "")

		// Parse exclude paths
		var excludePaths []string
		if excludePathsStr != "" {
			for _, excludePath := range strings.Split(excludePathsStr, ",") {
				if strings.TrimSpace(excludePath) != "" {
					excludePaths = append(excludePaths, strings.TrimSpace(excludePath))
				}
			}
		}

		// Scan filesystem advanced
		output, err := module.ScanFilesystemAdvanced(ctx, path, onlyVerified, excludePaths, outputFormat, maxDepth)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog advanced filesystem scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		image := request.GetString("image", "")
		if image == "" {
			return mcp.NewToolResultError("image is required"), nil
		}

		onlyVerified := request.GetBool("only_verified", false)
		outputFormat := request.GetString("output_format", "")
		layers := request.GetString("layers", "")

		// Scan docker advanced
		output, err := module.ScanDockerAdvanced(ctx, image, onlyVerified, outputFormat, layers)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog advanced docker scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		target := request.GetString("target", "")
		sourceType := request.GetString("source_type", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if sourceType == "" {
			return mcp.NewToolResultError("source_type is required"), nil
		}

		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		onlyVerified := request.GetBool("only_verified", false)
		includeDetectors := request.GetBool("include_detectors", false)
		confidenceLevel := request.GetString("confidence_level", "")

		// Parse exclude and include paths
		var excludePaths, includePaths []string
		if excludePathsStr := request.GetString("exclude_paths", ""); excludePathsStr != "" {
			for _, path := range strings.Split(excludePathsStr, ",") {
				if strings.TrimSpace(path) != "" {
					excludePaths = append(excludePaths, strings.TrimSpace(path))
				}
			}
		}
		if includePathsStr := request.GetString("include_paths", ""); includePathsStr != "" {
			for _, path := range strings.Split(includePathsStr, ",") {
				if strings.TrimSpace(path) != "" {
					includePaths = append(includePaths, strings.TrimSpace(path))
				}
			}
		}

		// Comprehensive secret detection
		output, err := module.ComprehensiveSecretDetection(ctx, target, sourceType, outputFormat, outputFile, onlyVerified, includeDetectors, confidenceLevel, excludePaths, includePaths)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog comprehensive detection failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		cloudProvider := request.GetString("cloud_provider", "")
		resourceIdentifier := request.GetString("resource_identifier", "")
		if cloudProvider == "" {
			return mcp.NewToolResultError("cloud_provider is required"), nil
		}
		if resourceIdentifier == "" {
			return mcp.NewToolResultError("resource_identifier is required"), nil
		}

		credentialsProfile := request.GetString("credentials_profile", "")
		region := request.GetString("region", "")
		recursive := request.GetBool("recursive", false)
		filePatterns := request.GetString("file_patterns", "")
		onlyVerified := request.GetBool("only_verified", false)
		maxFileSize := request.GetString("max_file_size", "")

		// Cloud storage scanning
		output, err := module.CloudStorageScanning(ctx, cloudProvider, resourceIdentifier, credentialsProfile, region, recursive, filePatterns, onlyVerified, maxFileSize)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog cloud storage scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		action := request.GetString("action", "")
		detectorConfig := request.GetString("detector_config", "")
		target := request.GetString("target", "")
		detectorPattern := request.GetString("detector_pattern", "")
		includeBuiltin := request.GetBool("include_builtin", false)

		if action == "" {
			return mcp.NewToolResultError("action is required"), nil
		}

		// Custom detector management
		output, err := module.CustomDetectorManagement(ctx, action, detectorConfig, target, detectorPattern, includeBuiltin)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog custom detector management failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		gitSource := request.GetString("git_source", "")
		repository := request.GetString("repository", "")
		if gitSource == "" {
			return mcp.NewToolResultError("git_source is required"), nil
		}
		if repository == "" {
			return mcp.NewToolResultError("repository is required"), nil
		}

		authentication := request.GetString("authentication", "")
		scanMode := request.GetString("scan_mode", "")
		commitRange := request.GetString("commit_range", "")
		branches := request.GetString("branches", "")
		includeForks := request.GetBool("include_forks", false)
		includeIssues := request.GetBool("include_issues", false)
		includePullRequests := request.GetBool("include_pull_requests", false)
		outputFormat := request.GetString("output_format", "")

		// Enterprise git scanning
		output, err := module.EnterpriseGitScanning(ctx, gitSource, repository, authentication, scanMode, commitRange, branches, includeForks, includeIssues, includePullRequests, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog enterprise git scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		scanTarget := request.GetString("scan_target", "")
		if scanTarget == "" {
			return mcp.NewToolResultError("scan_target is required"), nil
		}

		scanType := request.GetString("scan_type", "scheduled")
		baselineFile := request.GetString("baseline_file", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		failOnVerified := request.GetBool("fail_on_verified", false)
		failOnUnverified := request.GetBool("fail_on_unverified", false)
		quietMode := request.GetBool("quiet_mode", false)
		timeout := request.GetString("timeout", "")

		// CI/CD pipeline integration
		output, err := module.CICDPipelineIntegration(ctx, scanTarget, scanType, baselineFile, outputFormat, outputFile, failOnVerified, failOnUnverified, quietMode, timeout)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog CI/CD integration failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		target := request.GetString("target", "")
		sourceType := request.GetString("source_type", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if sourceType == "" {
			return mcp.NewToolResultError("source_type is required"), nil
		}

		concurrency := request.GetString("concurrency", "")
		maxFileSize := request.GetString("max_file_size", "")
		bufferSize := request.GetString("buffer_size", "")
		skipBinaries := request.GetBool("skip_binaries", false)
		enableSampling := request.GetBool("enable_sampling", false)
		memoryLimit := request.GetString("memory_limit", "")
		enableMetrics := request.GetBool("enable_metrics", false)

		// Performance optimization
		output, err := module.PerformanceOptimization(ctx, target, sourceType, concurrency, maxFileSize, bufferSize, skipBinaries, enableSampling, memoryLimit, enableMetrics)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog performance optimization failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Get parameters
		target := request.GetString("target", "")
		sourceType := request.GetString("source_type", "")
		reportType := request.GetString("report_type", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if sourceType == "" {
			return mcp.NewToolResultError("source_type is required"), nil
		}
		if reportType == "" {
			return mcp.NewToolResultError("report_type is required"), nil
		}

		outputFormats := request.GetString("output_formats", "")
		outputDirectory := request.GetString("output_directory", "")
		includeVerificationStatus := request.GetBool("include_verification_status", false)
		includeRiskAssessment := request.GetBool("include_risk_assessment", false)
		baselineComparison := request.GetString("baseline_comparison", "")
		includeTrends := request.GetBool("include_trends", false)

		// Comprehensive reporting
		output, err := module.ComprehensiveReporting(ctx, target, sourceType, reportType, outputFormats, outputDirectory, includeVerificationStatus, includeRiskAssessment, baselineComparison, includeTrends)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog comprehensive reporting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// TruffleHog get version tool
	getVersionTool := mcp.NewTool("trufflehog_get_version",
		mcp.WithDescription("Get TruffleHog version and detector information"),
		mcp.WithBoolean("show_detectors",
			mcp.Description("Include available detectors information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTruffleHogModule(client)

		// Note: show_detectors parameter not directly supported in Dagger GetVersion function
		if request.GetBool("show_detectors", false) {
			return mcp.NewToolResultError("Warning: show_detectors parameter not supported with direct Dagger calls"), nil
		}

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TruffleHog get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}