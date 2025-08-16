package mcp

import (
	"fmt"
	"strings"
)

// GenerateMCPHelpText dynamically generates help text from the modular registry
func GenerateMCPHelpText() string {
	var helpText strings.Builder
	
	helpText.WriteString("Start an MCP server that exposes specific Ship CLI tools for AI assistants.\n\n")
	helpText.WriteString("Available tools:\n")
	
	// Category order and display names
	categoryOrder := []string{"terraform", "security", "kubernetes", "cloud", "aws", "supply-chain"}
	categoryNames := map[string]string{
		"terraform":    "# Terraform Tools",
		"security":     "# Security Tools", 
		"kubernetes":   "# Kubernetes Tools",
		"cloud":        "# Cloud & Infrastructure Tools",
		"aws":          "# AWS IAM Tools",
		"supply-chain": "# Supply Chain Tools",
	}
	
	// Generate tools by category from registry
	for _, category := range categoryOrder {
		if tools, exists := ToolRegistry[category]; exists && len(tools) > 0 {
			helpText.WriteString("  ")
			helpText.WriteString(categoryNames[category])
			helpText.WriteString("\n")
			
			for _, tool := range tools {
				helpText.WriteString(fmt.Sprintf("  %-16s - %s\n", tool.Name, tool.Description))
			}
			helpText.WriteString("\n")
		}
	}
	
	// Add collections
	helpText.WriteString("  # Collections\n")
	for category := range ToolRegistry {
		helpText.WriteString(fmt.Sprintf("  %-16s - All %s tools\n", category, category))
	}
	helpText.WriteString("  all              - All tools (default if no tool specified)\n\n")
	
	// External MCP Servers (generated from external servers module)
	helpText.WriteString("External MCP Servers:\n")
	helpText.WriteString(generateExternalServersHelpText())
	
	// Examples
	helpText.WriteString("Examples:\n")
	helpText.WriteString("  # Our Security & Infrastructure Tools\n")
	helpText.WriteString("  ship mcp gitleaks    # MCP server for just Gitleaks\n")
	helpText.WriteString("  ship mcp security    # MCP server for all security tools\n")
	helpText.WriteString("  ship mcp all         # MCP server for all tools\n")
	helpText.WriteString("\n")
	helpText.WriteString("  # External MCP Servers\n")
	helpText.WriteString("  ship mcp filesystem     # Proxy filesystem operations MCP server\n")
	helpText.WriteString("  ship mcp memory         # Proxy memory/knowledge storage MCP server\n")
	helpText.WriteString("  ship mcp brave-search --var BRAVE_API_KEY=your_api_key   # Proxy Brave search with API key\n")
	helpText.WriteString("  ship mcp steampipe                                       # Proxy Steampipe for cloud infrastructure queries\n")
	helpText.WriteString("\n")
	helpText.WriteString("  # AWS Labs Official MCP Servers (requires 'uv' and AWS credentials)\n")
	helpText.WriteString("  ship mcp aws-core --var AWS_PROFILE=default             # AWS core operations\n")
	helpText.WriteString("  ship mcp aws-iam --var AWS_REGION=us-west-2             # AWS IAM management\n")
	helpText.WriteString("  ship mcp aws-pricing                                     # AWS pricing queries\n")
	helpText.WriteString("  ship mcp aws-eks --var AWS_PROFILE=production           # EKS operations\n")
	helpText.WriteString("  ship mcp aws-ec2 --var AWS_REGION=eu-west-1             # EC2 operations\n")
	helpText.WriteString("  ship mcp aws-s3                                          # S3 operations\n")
	helpText.WriteString("\n")
	helpText.WriteString("  ship mcp gitleaks --var DEBUG=true  # Pass environment variables")
	
	return helpText.String()
}

// generateExternalServersHelpText generates help text for external MCP servers
func generateExternalServersHelpText() string {
	var helpText strings.Builder
	
	// Group external servers by type
	officialServers := []string{"filesystem", "memory", "brave-search", "steampipe"}
	awsLabsServers := []string{"aws-core", "aws-iam", "aws-pricing", "aws-eks", "aws-ec2", "aws-s3"}
	
	// Official ModelContextProtocol servers
	for _, serverName := range officialServers {
		if config, exists := ExternalMCPServers[serverName]; exists {
			description := getServerDescription(serverName)
			helpText.WriteString(fmt.Sprintf("  %-16s - %s\n", config.Name, description))
		}
	}
	helpText.WriteString("\n")
	
	// AWS Labs servers
	helpText.WriteString("  # AWS Labs Official MCP Servers\n")
	for _, serverName := range awsLabsServers {
		if config, exists := ExternalMCPServers[serverName]; exists {
			description := getServerDescription(serverName)
			helpText.WriteString(fmt.Sprintf("  %-16s - %s\n", config.Name, description))
		}
	}
	helpText.WriteString("\n")
	
	return helpText.String()
}

// getServerDescription returns a user-friendly description for external servers
func getServerDescription(serverName string) string {
	descriptions := map[string]string{
		"filesystem":   "Filesystem operations MCP server",
		"memory":       "Memory/knowledge storage MCP server",
		"brave-search": "Brave search MCP server",
		"steampipe":    "Cloud infrastructure queries MCP server",
		"aws-core":     "AWS core operations and general services",
		"aws-iam":      "AWS IAM operations and identity management",
		"aws-pricing":  "AWS pricing and cost estimation",
		"aws-eks":      "AWS EKS and Kubernetes operations",
		"aws-ec2":      "AWS EC2 compute operations",
		"aws-s3":       "AWS S3 storage operations",
	}
	
	if desc, exists := descriptions[serverName]; exists {
		return desc
	}
	return "External MCP server"
}

// GenerateToolsResourceContent dynamically generates tools resource content from the modular registry
func GenerateToolsResourceContent() string {
	var content strings.Builder
	
	content.WriteString("# Ship CLI Tools\n\n")
	content.WriteString("This document lists all available Ship CLI tools accessible through the MCP server.\n\n")
	
	// Category order and display names
	categoryOrder := []string{"terraform", "security", "kubernetes", "cloud", "aws", "supply-chain"}
	categoryNames := map[string]string{
		"terraform":    "Terraform Tools",
		"security":     "Security Tools", 
		"kubernetes":   "Kubernetes Tools",
		"cloud":        "Cloud & Infrastructure Tools",
		"aws":          "AWS IAM Tools",
		"supply-chain": "Supply Chain Tools",
	}
	
	// Generate tools by category from registry
	for _, category := range categoryOrder {
		if tools, exists := ToolRegistry[category]; exists && len(tools) > 0 {
			content.WriteString("## ")
			content.WriteString(categoryNames[category])
			content.WriteString("\n\n")
			
			for _, tool := range tools {
				content.WriteString(fmt.Sprintf("- **%s**: %s\n", tool.Name, tool.Description))
			}
			content.WriteString("\n")
		}
	}
	
	// Add usage examples
	content.WriteString("## Usage Examples\n\n")
	content.WriteString("### Tool Collections\n")
	content.WriteString("- `ship mcp security` - Start MCP server with all security tools\n")
	content.WriteString("- `ship mcp terraform` - Start MCP server with all Terraform tools\n")
	content.WriteString("- `ship mcp all` - Start MCP server with all tools\n\n")
	
	content.WriteString("### Individual Tools\n")
	content.WriteString("- `ship mcp gitleaks` - Start MCP server with only Gitleaks tools\n")
	content.WriteString("- `ship mcp trivy` - Start MCP server with only Trivy tools\n")
	content.WriteString("- `ship mcp checkov` - Start MCP server with only Checkov tools\n\n")
	
	content.WriteString("### External MCP Servers\n")
	content.WriteString("- `ship mcp filesystem` - Proxy filesystem operations MCP server\n")
	content.WriteString("- `ship mcp aws-core` - Proxy AWS Labs official MCP server\n")
	
	return content.String()
}