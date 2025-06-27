package modules

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DiscoveryManager manages multiple module discovery sources
type DiscoveryManager struct {
	discoveries []Discovery
	config      ModuleConfig
}

// NewDiscoveryManager creates a new discovery manager
func NewDiscoveryManager(config ModuleConfig) *DiscoveryManager {
	dm := &DiscoveryManager{
		config: config,
	}
	
	// Add built-in discovery sources
	dm.addDiscovery(&BuiltinDiscovery{})
	dm.addDiscovery(&UserDirectoryDiscovery{})
	dm.addDiscovery(&ProjectDiscovery{})
	
	// Add git-based discovery if configured
	if len(config.Repositories) > 0 {
		dm.addDiscovery(&GitDiscovery{})
	}
	
	return dm
}

// addDiscovery adds a discovery source
func (dm *DiscoveryManager) addDiscovery(d Discovery) {
	d.SetConfig(dm.config)
	dm.discoveries = append(dm.discoveries, d)
}

// DiscoverAll discovers modules from all sources
func (dm *DiscoveryManager) DiscoverAll(ctx context.Context) ([]*Module, error) {
	var allModules []*Module
	moduleNames := make(map[string]bool) // Track duplicates
	
	for _, discovery := range dm.discoveries {
		modules, err := discovery.DiscoverModules(ctx)
		if err != nil {
			// Log error but continue with other sources
			fmt.Fprintf(os.Stderr, "Warning: discovery error from %s: %v\n", discovery.GetSourceType(), err)
			continue
		}
		
		// Add modules, handling duplicates (higher priority sources win)
		for _, module := range modules {
			if !moduleNames[module.Metadata.Name] {
				allModules = append(allModules, module)
				moduleNames[module.Metadata.Name] = true
			}
		}
	}
	
	return allModules, nil
}

// BuiltinDiscovery discovers built-in modules
type BuiltinDiscovery struct {
	config ModuleConfig
}

func (d *BuiltinDiscovery) SetConfig(config ModuleConfig) {
	d.config = config
}

func (d *BuiltinDiscovery) GetSourceType() string {
	return "builtin"
}

func (d *BuiltinDiscovery) DiscoverModules(ctx context.Context) ([]*Module, error) {
	// Built-in modules are hardcoded for now
	// In a real implementation, these would be discovered from internal/dagger/modules/
	return []*Module{
		{
			APIVersion: "ship.cloudship.ai/v1",
			Kind:       "Module",
			Metadata: ModuleMetadata{
				Name:        "terraform-tools",
				Version:     "1.0.0",
				Description: "Terraform analysis and documentation tools",
				Author:      "CloudshipAI",
			},
			Spec: ModuleSpec{
				Type: ModuleTypeDocker,
				Commands: []ModuleCommand{
					{
						Name:        "lint",
						Description: "Run TFLint on Terraform code",
					},
					{
						Name:        "checkov-scan",
						Description: "Run Checkov security scan",
					},
					{
						Name:        "cost-estimate",
						Description: "Estimate infrastructure costs",
					},
				},
			},
			Source:  "builtin",
			Trusted: true,
		},
		{
			APIVersion: "ship.cloudship.ai/v1",
			Kind:       "Module",
			Metadata: ModuleMetadata{
				Name:        "ai-investigate",
				Version:     "1.0.0",
				Description: "AI-powered infrastructure investigation",
				Author:      "CloudshipAI",
			},
			Spec: ModuleSpec{
				Type: ModuleTypeDocker,
				Commands: []ModuleCommand{
					{
						Name:        "investigate",
						Description: "Investigate infrastructure using natural language",
						Flags: []ModuleFlag{
							{
								Name:        "prompt",
								Type:        "string",
								Required:    true,
								Description: "Natural language investigation prompt",
							},
							{
								Name:        "provider",
								Type:        "string",
								Default:     "aws",
								Description: "Cloud provider",
								Enum:        []string{"aws", "azure", "gcp"},
							},
						},
					},
				},
			},
			Source:  "builtin",
			Trusted: true,
		},
	}, nil
}

// UserDirectoryDiscovery discovers modules from user directory
type UserDirectoryDiscovery struct {
	config ModuleConfig
}

func (d *UserDirectoryDiscovery) SetConfig(config ModuleConfig) {
	d.config = config
}

func (d *UserDirectoryDiscovery) GetSourceType() string {
	return "user"
}

func (d *UserDirectoryDiscovery) DiscoverModules(ctx context.Context) ([]*Module, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	userModulesDir := filepath.Join(homeDir, ".ship", "modules")
	return d.discoverModulesInDirectory(userModulesDir, "user")
}

func (d *UserDirectoryDiscovery) discoverModulesInDirectory(dir, source string) ([]*Module, error) {
	var modules []*Module
	
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return modules, nil
	}
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		modulePath := filepath.Join(dir, entry.Name())
		moduleYaml := filepath.Join(modulePath, "module.yaml")
		
		if _, err := os.Stat(moduleYaml); os.IsNotExist(err) {
			continue
		}
		
		module, err := d.loadModuleFromFile(moduleYaml, modulePath, source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load module from %s: %v\n", moduleYaml, err)
			continue
		}
		
		modules = append(modules, module)
	}
	
	return modules, nil
}

func (d *UserDirectoryDiscovery) loadModuleFromFile(yamlPath, modulePath, source string) (*Module, error) {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read module.yaml: %w", err)
	}
	
	var module Module
	if err := yaml.Unmarshal(data, &module); err != nil {
		return nil, fmt.Errorf("failed to parse module.yaml: %w", err)
	}
	
	// Set runtime fields
	module.Path = modulePath
	module.Source = source
	module.Trusted = source == "builtin" || (source == "user" && d.config.AllowUntrusted)
	
	return &module, nil
}

// ProjectDiscovery discovers modules from current project
type ProjectDiscovery struct {
	config ModuleConfig
}

func (d *ProjectDiscovery) SetConfig(config ModuleConfig) {
	d.config = config
}

func (d *ProjectDiscovery) GetSourceType() string {
	return "project"
}

func (d *ProjectDiscovery) DiscoverModules(ctx context.Context) ([]*Module, error) {
	var modules []*Module
	
	// Check for .ship/modules directory
	if _, err := os.Stat(".ship/modules"); err == nil {
		userDiscovery := &UserDirectoryDiscovery{config: d.config}
		projectModules, err := userDiscovery.discoverModulesInDirectory(".ship/modules", "project")
		if err != nil {
			return nil, err
		}
		modules = append(modules, projectModules...)
	}
	
	// Check for dagger.json (Dagger module)
	if _, err := os.Stat("dagger.json"); err == nil {
		daggerModule, err := d.discoverDaggerModule()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load Dagger module: %v\n", err)
		} else if daggerModule != nil {
			modules = append(modules, daggerModule)
		}
	}
	
	return modules, nil
}

func (d *ProjectDiscovery) discoverDaggerModule() (*Module, error) {
	// For now, create a basic Dagger module representation
	// In a real implementation, this would parse dagger.json and discover functions
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	
	projectName := filepath.Base(cwd)
	
	return &Module{
		APIVersion: "ship.cloudship.ai/v1",
		Kind:       "Module",
		Metadata: ModuleMetadata{
			Name:        fmt.Sprintf("dagger-%s", strings.ToLower(projectName)),
			Version:     "1.0.0",
			Description: fmt.Sprintf("Dagger module from %s", projectName),
			Author:      "Project",
		},
		Spec: ModuleSpec{
			Type: ModuleTypeDagger,
			Dagger: &DaggerModuleSpec{
				Module: ".",
			},
			Commands: []ModuleCommand{
				{
					Name:        "dagger-functions",
					Description: "Run Dagger functions from this project",
				},
			},
		},
		Path:    cwd,
		Source:  "project",
		Trusted: true, // Project modules are trusted
	}, nil
}

// GitDiscovery discovers modules from Git repositories
type GitDiscovery struct {
	config ModuleConfig
}

func (d *GitDiscovery) SetConfig(config ModuleConfig) {
	d.config = config
}

func (d *GitDiscovery) GetSourceType() string {
	return "git"
}

func (d *GitDiscovery) DiscoverModules(ctx context.Context) ([]*Module, error) {
	var modules []*Module
	
	// For now, return empty - git discovery would require git clone logic
	// This would be implemented in Phase 2 of the roadmap
	
	return modules, nil
}