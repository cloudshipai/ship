// Package all provides convenience functions for registering collections of Ship tools
package all

import (
	"github.com/cloudshipai/ship/internal/ship"
	"github.com/cloudshipai/ship/pkg/tools"
)

// TerraformRegistry returns a registry with all Terraform-related tools
func TerraformRegistry() *ship.Registry {
	registry := ship.NewRegistry()

	// Add TFLint tool
	registry.RegisterTool(tools.NewTFLintTool())

	// TODO: Add more Terraform tools as they're converted to the framework:
	// registry.RegisterTool(tools.NewCheckovTool())
	// registry.RegisterTool(tools.NewTrivyTool())
	// registry.RegisterTool(tools.NewCostAnalysisTool())
	// registry.RegisterTool(tools.NewDocsTool())
	// registry.RegisterTool(tools.NewDiagramTool())

	return registry
}

// SecurityRegistry returns a registry with all security-related tools
func SecurityRegistry() *ship.Registry {
	registry := ship.NewRegistry()

	// TODO: Add security tools as they're converted:
	// registry.RegisterTool(tools.NewCheckovTool())
	// registry.RegisterTool(tools.NewTrivyTool())

	return registry
}

// DocsRegistry returns a registry with all documentation-related tools
func DocsRegistry() *ship.Registry {
	registry := ship.NewRegistry()

	// TODO: Add documentation tools as they're converted:
	// registry.RegisterTool(tools.NewDocsTool())
	// registry.RegisterTool(tools.NewDiagramTool())

	return registry
}

// AllRegistry returns a registry with all Ship tools
func AllRegistry() *ship.Registry {
	registry := ship.NewRegistry()

	// Import all tool collections
	registry.ImportFrom(TerraformRegistry())
	registry.ImportFrom(SecurityRegistry())
	registry.ImportFrom(DocsRegistry())

	return registry
}

// AddTerraformTools is a convenience function to add all Terraform tools to a server builder
func AddTerraformTools(builder *ship.ServerBuilder) *ship.ServerBuilder {
	return builder.ImportRegistry(TerraformRegistry())
}

// AddSecurityTools is a convenience function to add all security tools to a server builder
func AddSecurityTools(builder *ship.ServerBuilder) *ship.ServerBuilder {
	return builder.ImportRegistry(SecurityRegistry())
}

// AddDocsTools is a convenience function to add all documentation tools to a server builder
func AddDocsTools(builder *ship.ServerBuilder) *ship.ServerBuilder {
	return builder.ImportRegistry(DocsRegistry())
}

// AddAllTools is a convenience function to add all Ship tools to a server builder
func AddAllTools(builder *ship.ServerBuilder) *ship.ServerBuilder {
	return builder.ImportRegistry(AllRegistry())
}
