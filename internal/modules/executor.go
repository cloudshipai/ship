package modules

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ExecutorRegistry manages module executors
type ExecutorRegistry struct {
	executors map[ModuleType]ModuleExecutor
}

// NewExecutorRegistry creates a new executor registry
func NewExecutorRegistry() *ExecutorRegistry {
	registry := &ExecutorRegistry{
		executors: make(map[ModuleType]ModuleExecutor),
	}

	// Register built-in executors
	registry.Register(ModuleTypeDocker, &DockerExecutor{})
	registry.Register(ModuleTypeDagger, &DaggerExecutor{})

	return registry
}

// Register registers an executor for a module type
func (r *ExecutorRegistry) Register(moduleType ModuleType, executor ModuleExecutor) {
	r.executors[moduleType] = executor
}

// Execute executes a module using the appropriate executor
func (r *ExecutorRegistry) Execute(ctx context.Context, module *Module, command string, args []string, flags map[string]interface{}) (*ExecutionResult, error) {
	executor, exists := r.executors[module.Spec.Type]
	if !exists {
		return nil, fmt.Errorf("no executor found for module type: %s", module.Spec.Type)
	}

	if !executor.CanExecute(module) {
		return nil, fmt.Errorf("executor cannot execute module: %s", module.Metadata.Name)
	}

	return executor.Execute(ctx, module, command, args, flags)
}

// DockerExecutor executes Docker-based modules
type DockerExecutor struct{}

func (e *DockerExecutor) CanExecute(module *Module) bool {
	return module.Spec.Type == ModuleTypeDocker && module.Spec.Docker != nil
}

func (e *DockerExecutor) Execute(ctx context.Context, module *Module, command string, args []string, flags map[string]interface{}) (*ExecutionResult, error) {
	start := time.Now()

	// For Docker modules, we'll delegate to the existing Ship CLI commands
	// This is a bridge until we have full Docker execution
	result, err := e.executeViaShipCLI(ctx, module, command, args, flags)

	result.Duration = time.Since(start)
	return result, err
}

func (e *DockerExecutor) executeViaShipCLI(ctx context.Context, module *Module, command string, args []string, flags map[string]interface{}) (*ExecutionResult, error) {
	// Map module commands to Ship CLI commands
	var shipArgs []string

	switch module.Metadata.Name {
	case "terraform-tools":
		shipArgs = append(shipArgs, "tf", command)
	default:
		return nil, fmt.Errorf("unknown built-in module: %s", module.Metadata.Name)
	}

	// Add remaining args
	shipArgs = append(shipArgs, args...)

	// Execute the Ship CLI command
	cmd := exec.CommandContext(ctx, "ship", shipArgs...)
	output, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &ExecutionResult{
		ExitCode: exitCode,
		Stdout:   string(output),
		Stderr:   "",
	}, nil
}

// DaggerExecutor executes Dagger-based modules
type DaggerExecutor struct{}

func (e *DaggerExecutor) CanExecute(module *Module) bool {
	return module.Spec.Type == ModuleTypeDagger && module.Spec.Dagger != nil
}

func (e *DaggerExecutor) Execute(ctx context.Context, module *Module, command string, args []string, flags map[string]interface{}) (*ExecutionResult, error) {
	start := time.Now()

	// For now, return a placeholder
	// In a real implementation, this would use the Dagger SDK to execute functions
	result := &ExecutionResult{
		ExitCode: 0,
		Stdout:   fmt.Sprintf("Dagger module %s executed successfully (placeholder)", module.Metadata.Name),
		Stderr:   "",
		Duration: time.Since(start),
		Metadata: map[string]string{
			"module":   module.Spec.Dagger.Module,
			"function": module.Spec.Dagger.Function,
		},
	}

	return result, nil
}
