package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Manager handles module discovery, loading, and execution
type Manager struct {
	discovery *DiscoveryManager
	executor  *ExecutorRegistry
	modules   []*Module
	config    ModuleConfig
}

// NewManager creates a new module manager
func NewManager(config ModuleConfig) *Manager {
	return &Manager{
		discovery: NewDiscoveryManager(config),
		executor:  NewExecutorRegistry(),
		config:    config,
	}
}

// LoadModules discovers and loads all available modules
func (m *Manager) LoadModules(ctx context.Context) error {
	modules, err := m.discovery.DiscoverAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover modules: %w", err)
	}
	
	m.modules = modules
	return nil
}

// GetModules returns all loaded modules
func (m *Manager) GetModules() []*Module {
	return m.modules
}

// GetModule returns a module by name
func (m *Manager) GetModule(name string) (*Module, error) {
	for _, module := range m.modules {
		if module.Metadata.Name == name {
			return module, nil
		}
	}
	return nil, fmt.Errorf("module not found: %s", name)
}

// ExecuteModule executes a module command
func (m *Manager) ExecuteModule(ctx context.Context, moduleName, command string, args []string, flags map[string]interface{}) (*ExecutionResult, error) {
	module, err := m.GetModule(moduleName)
	if err != nil {
		return nil, err
	}
	
	// Validate command exists
	if !m.hasCommand(module, command) {
		return nil, fmt.Errorf("command '%s' not found in module '%s'", command, moduleName)
	}
	
	// Security check
	if !module.Trusted && !m.config.AllowUntrusted {
		return nil, fmt.Errorf("module '%s' is not trusted and untrusted modules are disabled", moduleName)
	}
	
	return m.executor.Execute(ctx, module, command, args, flags)
}

// hasCommand checks if a module has a specific command
func (m *Manager) hasCommand(module *Module, command string) bool {
	for _, cmd := range module.Spec.Commands {
		if cmd.Name == command {
			return true
		}
	}
	return false
}

// RegisterDynamicCommands registers CLI commands for all discovered modules
func (m *Manager) RegisterDynamicCommands(rootCmd *cobra.Command) error {
	for _, module := range m.modules {
		if err := m.registerModuleCommands(rootCmd, module); err != nil {
			return fmt.Errorf("failed to register commands for module %s: %w", module.Metadata.Name, err)
		}
	}
	return nil
}

// registerModuleCommands registers CLI commands for a specific module
func (m *Manager) registerModuleCommands(rootCmd *cobra.Command, module *Module) error {
	for _, cmdSpec := range module.Spec.Commands {
		// Create command
		cmd := &cobra.Command{
			Use:   cmdSpec.Name,
			Short: cmdSpec.Description,
			Long:  fmt.Sprintf("%s\n\nModule: %s (%s)", cmdSpec.Description, module.Metadata.Name, module.Source),
			RunE:  m.createCommandRunner(module, cmdSpec),
		}
		
		// Add flags
		for _, flag := range cmdSpec.Flags {
			if err := m.addFlag(cmd, flag); err != nil {
				return fmt.Errorf("failed to add flag %s: %w", flag.Name, err)
			}
		}
		
		// Add examples
		if len(cmdSpec.Examples) > 0 {
			cmd.Example = strings.Join(cmdSpec.Examples, "\n")
		}
		
		// Register command
		rootCmd.AddCommand(cmd)
	}
	
	return nil
}

// createCommandRunner creates a command runner for a module command
func (m *Manager) createCommandRunner(module *Module, cmdSpec ModuleCommand) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		
		// Extract flags
		flags := make(map[string]interface{})
		for _, flagSpec := range cmdSpec.Flags {
			value, err := m.getFlagValue(cmd, flagSpec)
			if err != nil {
				return err
			}
			if value != nil {
				flags[flagSpec.Name] = value
			}
		}
		
		// Execute module
		result, err := m.ExecuteModule(ctx, module.Metadata.Name, cmdSpec.Name, args, flags)
		if err != nil {
			return err
		}
		
		// Output result
		if result.Stdout != "" {
			fmt.Print(result.Stdout)
		}
		if result.Stderr != "" {
			fmt.Fprint(cmd.ErrOrStderr(), result.Stderr)
		}
		
		// Set exit code
		if result.ExitCode != 0 {
			return fmt.Errorf("module execution failed with exit code %d", result.ExitCode)
		}
		
		return nil
	}
}

// addFlag adds a flag to a command based on the flag specification
func (m *Manager) addFlag(cmd *cobra.Command, flag ModuleFlag) error {
	switch flag.Type {
	case "string":
		defaultVal := ""
		if flag.Default != nil {
			if val, ok := flag.Default.(string); ok {
				defaultVal = val
			}
		}
		cmd.Flags().String(flag.Name, defaultVal, flag.Description)
		if flag.Short != "" {
			cmd.Flags().StringP(flag.Name, flag.Short, defaultVal, flag.Description)
		}
		if flag.Required {
			cmd.MarkFlagRequired(flag.Name)
		}
		
	case "bool":
		defaultVal := false
		if flag.Default != nil {
			if val, ok := flag.Default.(bool); ok {
				defaultVal = val
			}
		}
		cmd.Flags().Bool(flag.Name, defaultVal, flag.Description)
		if flag.Short != "" {
			cmd.Flags().BoolP(flag.Name, flag.Short, defaultVal, flag.Description)
		}
		
	case "int":
		defaultVal := 0
		if flag.Default != nil {
			if val, ok := flag.Default.(int); ok {
				defaultVal = val
			}
		}
		cmd.Flags().Int(flag.Name, defaultVal, flag.Description)
		if flag.Short != "" {
			cmd.Flags().IntP(flag.Name, flag.Short, defaultVal, flag.Description)
		}
		if flag.Required {
			cmd.MarkFlagRequired(flag.Name)
		}
		
	case "[]string":
		defaultVal := []string{}
		if flag.Default != nil {
			if val, ok := flag.Default.([]string); ok {
				defaultVal = val
			}
		}
		cmd.Flags().StringSlice(flag.Name, defaultVal, flag.Description)
		if flag.Short != "" {
			cmd.Flags().StringSliceP(flag.Name, flag.Short, defaultVal, flag.Description)
		}
		
	default:
		return fmt.Errorf("unsupported flag type: %s", flag.Type)
	}
	
	return nil
}

// getFlagValue gets the value of a flag from the command
func (m *Manager) getFlagValue(cmd *cobra.Command, flag ModuleFlag) (interface{}, error) {
	switch flag.Type {
	case "string":
		return cmd.Flags().GetString(flag.Name)
	case "bool":
		return cmd.Flags().GetBool(flag.Name)
	case "int":
		return cmd.Flags().GetInt(flag.Name)
	case "[]string":
		return cmd.Flags().GetStringSlice(flag.Name)
	default:
		return nil, fmt.Errorf("unsupported flag type: %s", flag.Type)
	}
}