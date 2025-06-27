# Dynamic Dagger Module Discovery System

## Overview

The Ship CLI should be able to dynamically discover and use both built-in and user-defined Dagger modules. This enables extensibility and allows users to add custom tools without modifying the core Ship CLI codebase.

## Design Goals

1. **Zero Configuration**: Modules should be discoverable without manual registration
2. **User Extensibility**: Users can add custom modules via Git repositories or local paths  
3. **Security**: Only trusted modules should be executed
4. **Performance**: Module discovery should be fast and cached when possible
5. **Compatibility**: Both Dagger Functions and traditional Docker-based modules should work

## Discovery Mechanisms

### 1. Built-in Modules
- Located in `internal/dagger/modules/`
- Compiled into the binary
- Always available and trusted

### 2. User Modules Directory
- `~/.ship/modules/` - Local user modules
- `~/.ship/dagger/` - Dagger modules (cloned from git)
- Auto-discovered on CLI startup

### 3. Project-level Modules
- `.ship/modules/` in current working directory  
- `dagger/` directory with dagger.json
- Higher priority than user modules

### 4. Git-based Modules
- Configured in `~/.ship/config.yaml`
- Automatically cloned and updated
- Cached locally for performance

## Module Structure

### Traditional Docker Module
```
my-custom-tool/
├── module.yaml          # Module metadata
├── Dockerfile           # Container definition
└── entrypoint.sh       # Execution logic
```

### Dagger Function Module
```
my-dagger-module/
├── dagger.json          # Dagger module config
├── dagger/              # Dagger module code
├── main.go             # Module implementation
└── README.md           # Documentation
```

## Configuration

### Ship Config (`~/.ship/config.yaml`)
```yaml
modules:
  # Git-based modules
  repositories:
    - url: "https://github.com/user/ship-modules.git"
      ref: "main"
      path: "modules"
      trusted: true
    - url: "https://github.com/cloudship/community-modules.git"
      ref: "v1.0.0"
      trusted: true
  
  # Local module directories
  directories:
    - "~/.ship/modules"
    - "./dagger"
  
  # Security settings
  allow_untrusted: false
  sandbox: true
```

### Module Metadata (`module.yaml`)
```yaml
apiVersion: ship.cloudship.ai/v1
kind: Module
metadata:
  name: my-custom-tool
  version: "1.0.0"
  description: "Custom security scanner"
  author: "user@company.com"
  
spec:
  type: docker | dagger
  
  # For Docker modules
  docker:
    image: "my-org/custom-scanner:latest"
    entrypoint: ["./scan.sh"]
    
  # For Dagger modules  
  dagger:
    module: "./dagger"
    function: "scan"
    
  # CLI integration
  commands:
    - name: "custom-scan"
      description: "Run custom security scan"
      flags:
        - name: "target"
          type: "string"
          required: true
          description: "Target to scan"
        - name: "format"
          type: "string"
          default: "json"
          enum: ["json", "yaml", "table"]
          
  # Dependencies
  dependencies:
    - "terraform >= 1.0"
    - "docker"
    
  # Security
  permissions:
    - "network"
    - "filesystem:read"
```

## Implementation Plan

### Phase 1: Basic Discovery
1. Scan built-in modules directory
2. Scan user modules directory (`~/.ship/modules/`)
3. Load module metadata
4. Register CLI commands dynamically

### Phase 2: Git Integration
1. Add git-based module support in config
2. Implement module caching and updates
3. Add `ship modules` management commands

### Phase 3: Dagger Functions
1. Integrate with Dagger SDK
2. Auto-discover Dagger modules in project
3. Support Dagger function execution

### Phase 4: Security & Sandboxing
1. Module signing and verification
2. Sandboxed execution
3. Permission system

## CLI Integration

### Dynamic Command Registration
```go
// Module discovery at startup
func discoverModules() ([]*Module, error) {
    var modules []*Module
    
    // 1. Built-in modules
    builtins := discoverBuiltinModules()
    modules = append(modules, builtins...)
    
    // 2. User modules
    userModules := discoverUserModules()
    modules = append(modules, userModules...)
    
    // 3. Project modules
    projectModules := discoverProjectModules()
    modules = append(modules, projectModules...)
    
    // 4. Git modules
    gitModules := discoverGitModules()
    modules = append(modules, gitModules...)
    
    return modules, nil
}

// Register commands dynamically
func registerModuleCommands(cmd *cobra.Command, modules []*Module) {
    for _, module := range modules {
        moduleCmd := &cobra.Command{
            Use:   module.Name,
            Short: module.Description,
            RunE:  createModuleRunner(module),
        }
        
        // Add flags from module spec
        for _, flag := range module.Spec.Commands[0].Flags {
            addFlagToCommand(moduleCmd, flag)
        }
        
        cmd.AddCommand(moduleCmd)
    }
}
```

### Module Execution
```go
func executeModule(module *Module, args []string) error {
    switch module.Spec.Type {
    case "docker":
        return executeDockerModule(module, args)
    case "dagger":
        return executeDaggerModule(module, args)
    default:
        return fmt.Errorf("unsupported module type: %s", module.Spec.Type)
    }
}
```

## Module Management Commands

```bash
# List available modules
ship modules list

# Install module from git
ship modules install https://github.com/user/ship-modules.git

# Update all modules
ship modules update

# Remove module
ship modules remove my-custom-tool

# Module information
ship modules info my-custom-tool

# Create new module template
ship modules new my-tool --type=docker
ship modules new my-tool --type=dagger
```

## Security Considerations

1. **Module Signing**: Cryptographic signatures for trusted modules
2. **Sandboxing**: Isolated execution environment
3. **Permissions**: Granular permission system
4. **Code Review**: Community modules require review
5. **Allow Lists**: Organization-level module allow lists

## Examples

### Custom Terraform Module
```yaml
# ~/.ship/modules/terraform-graph/module.yaml
apiVersion: ship.cloudship.ai/v1
kind: Module
metadata:
  name: terraform-graph
  description: "Generate Terraform dependency graphs"

spec:
  type: docker
  docker:
    image: "hashicorp/terraform:latest"
    entrypoint: ["terraform", "graph"]
    
  commands:
    - name: "terraform-graph"
      flags:
        - name: "output"
          type: "string"
          default: "graph.dot"
```

### Custom Dagger Function
```yaml
# ~/.ship/modules/custom-security/module.yaml  
apiVersion: ship.cloudship.ai/v1
kind: Module
metadata:
  name: custom-security
  description: "Custom security analysis"

spec:
  type: dagger
  dagger:
    function: "securityScan"
    
  commands:
    - name: "security-scan"
      flags:
        - name: "severity"
          type: "string"
          enum: ["low", "medium", "high", "critical"]
```

This design provides a flexible, secure, and extensible module system that can grow with user needs while maintaining the simplicity and reliability of Ship CLI.