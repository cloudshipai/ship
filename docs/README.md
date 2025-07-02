# Ship CLI Documentation

Welcome to the Ship CLI documentation! This folder contains comprehensive documentation including architecture diagrams, API references, and development guides.

## 📚 Documentation Structure

### Architecture & Design
- **[architecture-diagrams.md](./architecture-diagrams.md)** - Comprehensive Mermaid diagrams showing:
  - High-level system architecture
  - Component relationships
  - User workflows
  - Module system design
  - CloudShip integration
  - Container orchestration

- **[ai-system-architecture.md](./ai-system-architecture.md)** - Detailed AI system diagrams:
  - AI component overview
  - Investigation flow
  - Hallucination prevention
  - LLM integration patterns
  - Prompt engineering architecture

- **[data-flow-diagrams.md](./data-flow-diagrams.md)** - Data flow and state management:
  - Overall data flow
  - Command execution states
  - Configuration precedence
  - Container lifecycle
  - Error propagation

### Developer Guides
- **[developer-guide-diagrams.md](./developer-guide-diagrams.md)** - Visual guides for developers:
  - Adding new commands
  - Creating Dagger modules
  - Testing strategies
  - Contributing workflow
  - Performance optimization

### API & CLI Reference
- **[api-reference.md](./api-reference.md)** - CloudShip API documentation
- **[cli-reference.md](./cli-reference.md)** - Complete CLI command reference

### Product Documentation
- **[PRD.md](./PRD.md)** - Product Requirements Document
- **[technical-spec.md](./technical-spec.md)** - Technical specifications
- **[implementation-plan.md](./implementation-plan.md)** - Phased development plan
- **[development-tasks.md](./development-tasks.md)** - Sprint task breakdown

## 🎯 Quick Navigation

### For Users
1. Start with [cli-reference.md](./cli-reference.md) to learn commands
2. Check [architecture-diagrams.md](./architecture-diagrams.md) for system overview
3. Review [api-reference.md](./api-reference.md) for CloudShip integration

### For Developers
1. Read [developer-guide-diagrams.md](./developer-guide-diagrams.md) for contribution guide
2. Study [data-flow-diagrams.md](./data-flow-diagrams.md) to understand internals
3. Reference [ai-system-architecture.md](./ai-system-architecture.md) for AI components

### For AI/LLM Integration
1. See [ai-system-architecture.md](./ai-system-architecture.md) for AI architecture
2. Check `../llms.txt` for LLM-specific instructions
3. Review MCP integration in [architecture-diagrams.md](./architecture-diagrams.md)

## 🔍 Key Concepts Explained

### Dagger Integration
Ship CLI uses Dagger for container orchestration, ensuring all tools run in isolated environments without requiring local installation.

### AI Investigation
Natural language prompts are converted to Steampipe SQL queries through LLM integration, with built-in table knowledge to prevent hallucination.

### CloudShip Push
Results from any tool can be automatically uploaded to CloudShip for centralized analysis using the `--push` flag.

### Module System
Extensible architecture allows adding new tools as Dagger modules, with consistent interfaces for execution and output handling.

## 📊 Diagram Rendering

All diagrams are written in Mermaid and can be rendered in:
- GitHub (automatic rendering)
- VS Code (with Mermaid extension)
- Online at [mermaid.live](https://mermaid.live/)
- Any Markdown viewer with Mermaid support

## 🚀 Getting Started

1. **New to Ship CLI?** Start with the [CLI Reference](./cli-reference.md)
2. **Want to contribute?** Check the [Developer Guide](./developer-guide-diagrams.md)
3. **Building integrations?** See the [Architecture Diagrams](./architecture-diagrams.md)
4. **Working with AI features?** Review the [AI System Architecture](./ai-system-architecture.md)

---

For the latest updates and more information, visit the [Ship CLI GitHub repository](https://github.com/cloudshipai/ship).