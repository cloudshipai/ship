# InfraMap

Read Terraform and generate visual infrastructure graphs.

## Description

InfraMap reads Terraform state files and HCL configurations to generate visual infrastructure graphs. It creates human-readable diagrams showing the relationships between cloud resources, helping teams understand complex infrastructure layouts. InfraMap supports both Terraform state files (versions 3 and 4) and HCL configuration files, automatically detecting the input type and generating graphs in various formats.

## MCP Tools

### Infrastructure Visualization
- **`inframap_generate`** - Generate infrastructure graph from Terraform state or HCL
- **`inframap_prune`** - Remove unnecessary information from state/HCL files

## Real CLI Commands Used

### Core Commands
- `inframap generate <input>` - Generate graph from file or directory
- `inframap generate state.tfstate` - Generate from Terraform state
- `inframap generate config.tf` - Generate from HCL file
- `inframap generate ./my-module/` - Generate from HCL directory
- `inframap prune <input>` - Prune unnecessary information

### Input Type Control
- `inframap generate --hcl <input>` - Force HCL input type
- `inframap generate --tfstate <input>` - Force Terraform state input type

### Graph Customization
- `inframap generate --connections=false <input>` - Disable connections
- `inframap generate --raw <input>` - Show unprocessed configuration
- `inframap generate --clean=false <input>` - Keep unconnected nodes

## Use Cases

### Infrastructure Documentation
- **Visual Documentation**: Create diagrams for architecture documentation
- **Team Onboarding**: Help new team members understand infrastructure
- **Architecture Reviews**: Visual aids for design discussions
- **Compliance Documentation**: Generate diagrams for audit requirements

### Infrastructure Analysis
- **Dependency Mapping**: Understand resource relationships
- **Impact Analysis**: Visualize change impacts before modifications
- **Resource Discovery**: Find all resources in complex configurations
- **Optimization Planning**: Identify optimization opportunities

### Development Workflow
- **Code Reviews**: Include infrastructure diagrams in PR reviews
- **Debugging**: Visualize infrastructure during troubleshooting
- **Planning**: Design infrastructure changes visually
- **Communication**: Share infrastructure designs with stakeholders

### Education and Training
- **Learning Tool**: Understand Terraform configurations visually
- **Best Practices**: Demonstrate good architecture patterns
- **Workshops**: Create visual examples for training
- **Knowledge Transfer**: Document tribal knowledge visually

## Configuration Examples

### Basic Graph Generation
```bash
# Generate graph from Terraform state
inframap generate terraform.tfstate > graph.dot

# Generate graph from HCL file
inframap generate main.tf > graph.dot

# Generate graph from module directory
inframap generate ./modules/vpc/ > graph.dot

# Force input type detection
inframap generate --tfstate state.json > graph.dot
inframap generate --hcl config.tf > graph.dot
```

### Graph Customization
```bash
# Disable connection visualization
inframap generate --connections=false terraform.tfstate > graph.dot

# Show raw configuration without processing
inframap generate --raw ./terraform/ > graph.dot

# Keep unconnected nodes in graph
inframap generate --clean=false state.tfstate > graph.dot

# Combine multiple options
inframap generate --raw --connections=false --clean=false ./module/ > graph.dot
```

### Output Processing
```bash
# Generate PNG using Graphviz
inframap generate state.tfstate | dot -Tpng > infrastructure.png

# Generate SVG for web display
inframap generate main.tf | dot -Tsvg > infrastructure.svg

# Generate PDF for documentation
inframap generate ./terraform/ | dot -Tpdf > infrastructure.pdf

# Use graph-easy for ASCII output
inframap generate state.tfstate | graph-easy --as=ascii
```

### File Processing
```bash
# Prune unnecessary information
inframap prune large-state.tfstate > pruned-state.tfstate

# Prune HCL configuration
inframap prune complex-config.tf > simplified-config.tf
```

## Advanced Usage

### Docker Integration
```bash
# Run InfraMap in Docker
docker run --rm -v ${PWD}:/opt cycloid/inframap generate /opt/terraform.tfstate

# Generate PNG with Docker
docker run --rm -v ${PWD}:/opt cycloid/inframap generate /opt/terraform.tfstate | dot -Tpng > graph.png

# Process multiple files
docker run --rm -v ${PWD}:/opt cycloid/inframap generate /opt/modules/ | dot -Tsvg > modules.svg
```

### Automation Scripts
```bash
#!/bin/bash
# generate-diagrams.sh

# Generate diagrams for all environments
for env in dev staging prod; do
    echo "Generating diagram for $env environment..."
    inframap generate "environments/$env/terraform.tfstate" | \
        dot -Tpng > "docs/diagrams/$env-infrastructure.png"
done

# Generate module diagrams
for module in modules/*/; do
    module_name=$(basename "$module")
    echo "Generating diagram for $module_name module..."
    inframap generate "$module" | \
        dot -Tsvg > "docs/modules/$module_name.svg"
done
```

### CI/CD Integration
```yaml
# GitHub Actions
name: Generate Infrastructure Diagrams
on: [push, pull_request]
jobs:
  diagrams:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Install dependencies
      run: |
        sudo apt-get update
        sudo apt-get install -y graphviz
        wget -O inframap https://github.com/cycloidio/inframap/releases/latest/download/inframap-linux-amd64
        chmod +x inframap
        sudo mv inframap /usr/local/bin/
    
    - name: Generate diagrams
      run: |
        mkdir -p diagrams
        ./inframap generate terraform.tfstate | dot -Tpng > diagrams/infrastructure.png
        ./inframap generate modules/ | dot -Tsvg > diagrams/modules.svg
    
    - name: Upload diagrams
      uses: actions/upload-artifact@v2
      with:
        name: infrastructure-diagrams
        path: diagrams/
```

### Makefile Integration
```makefile
# Makefile
.PHONY: diagrams clean-diagrams

DIAGRAMS_DIR := docs/diagrams
STATE_FILES := $(wildcard environments/*/terraform.tfstate)
MODULE_DIRS := $(wildcard modules/*/)

diagrams: $(DIAGRAMS_DIR) $(STATE_FILES:.tfstate=.png) $(MODULE_DIRS:modules/%/=%.svg)

$(DIAGRAMS_DIR):
	mkdir -p $(DIAGRAMS_DIR)

$(DIAGRAMS_DIR)/%.png: environments/%/terraform.tfstate
	inframap generate $< | dot -Tpng > $@

$(DIAGRAMS_DIR)/%.svg: modules/%/
	inframap generate $< | dot -Tsvg > $@

clean-diagrams:
	rm -rf $(DIAGRAMS_DIR)

# Generate all diagrams
update-docs: diagrams
	git add $(DIAGRAMS_DIR)
	git commit -m "Update infrastructure diagrams"
```

## Integration Patterns

### Documentation Workflow
```bash
# Generate comprehensive documentation
#!/bin/bash

# Create documentation structure
mkdir -p docs/{diagrams,architecture}

# Generate overview diagram
inframap generate terraform.tfstate | dot -Tpng > docs/diagrams/overview.png

# Generate detailed module diagrams
for module in modules/*/; do
    name=$(basename "$module")
    inframap generate --raw "$module" | dot -Tsvg > "docs/diagrams/$name-detailed.svg"
    inframap generate "$module" | dot -Tpng > "docs/diagrams/$name-overview.png"
done

# Generate simplified architecture view
inframap generate --connections=false terraform.tfstate | \
    dot -Tpng > docs/architecture/components.png

echo "Documentation diagrams generated in docs/ directory"
```

### Multi-Environment Visualization
```bash
# compare-environments.sh
#!/bin/bash

environments=("dev" "staging" "prod")

echo "Generating environment comparison diagrams..."

for env in "${environments[@]}"; do
    # Standard diagram
    inframap generate "environments/$env/terraform.tfstate" | \
        dot -Tpng > "docs/environments/$env.png"
    
    # Raw resource view
    inframap generate --raw "environments/$env/terraform.tfstate" | \
        dot -Tsvg > "docs/environments/$env-raw.svg"
    
    # Simplified view
    inframap generate --connections=false "environments/$env/terraform.tfstate" | \
        dot -Tpng > "docs/environments/$env-simple.png"
done
```

### Module Development
```bash
# module-workflow.sh
#!/bin/bash

MODULE_DIR="$1"
if [[ -z "$MODULE_DIR" ]]; then
    echo "Usage: $0 <module-directory>"
    exit 1
fi

echo "Generating diagrams for module: $MODULE_DIR"

# Generate module overview
inframap generate "$MODULE_DIR" | dot -Tpng > "$MODULE_DIR/README-diagram.png"

# Generate detailed view
inframap generate --raw "$MODULE_DIR" | dot -Tsvg > "$MODULE_DIR/architecture.svg"

# Generate dependency view
inframap generate --clean=false "$MODULE_DIR" | dot -Tpng > "$MODULE_DIR/dependencies.png"

echo "Module diagrams generated:"
echo "  - README-diagram.png (overview)"
echo "  - architecture.svg (detailed)"
echo "  - dependencies.png (with dependencies)"
```

## Best Practices

### Graph Generation
- **Regular Updates**: Generate diagrams on infrastructure changes
- **Version Control**: Include diagrams in version control
- **Multiple Views**: Create different diagram types for different audiences
- **Consistent Naming**: Use consistent file naming for easy navigation

### Documentation Integration
- **README Files**: Include overview diagrams in module README files
- **Architecture Docs**: Use detailed diagrams in architecture documentation
- **Change Reviews**: Include diagrams in infrastructure change reviews
- **Onboarding**: Use diagrams for team onboarding materials

### Automation Strategy
- **CI/CD Integration**: Automatically generate diagrams in pipelines
- **Git Hooks**: Update diagrams on commit or push
- **Scheduled Generation**: Regularly update diagrams for evolving infrastructure
- **Quality Gates**: Verify diagram generation doesn't fail

### Output Management
- **Format Selection**: Choose appropriate formats for different use cases
- **File Organization**: Organize diagrams in logical directory structures
- **Size Optimization**: Use appropriate formats to manage file sizes
- **Accessibility**: Ensure diagrams are accessible to all team members

## Error Handling

### Common Issues
```bash
# Invalid input file
inframap generate nonexistent.tfstate
# Solution: Verify file exists and is readable

# Unsupported Terraform version
inframap generate old-state.tfstate
# Solution: Upgrade state file or use compatible version

# Memory issues with large states
inframap generate --raw huge-state.tfstate
# Solution: Use pruning or process in smaller chunks

# Missing Graphviz for output processing
inframap generate state.tfstate | dot -Tpng
# Solution: Install graphviz package
```

### Troubleshooting
- **Validation**: Verify input files are valid Terraform state or HCL
- **Dependencies**: Ensure required tools (graphviz) are installed
- **Memory**: Monitor memory usage with large infrastructure states
- **Permissions**: Check file permissions for input and output files

InfraMap provides essential visualization capabilities for Terraform infrastructure, enabling teams to understand, document, and communicate complex cloud architectures through clear visual diagrams.