# OpenCode AI Coding Assistant

OpenCode is an AI-powered coding assistant that provides interactive chat, code generation, analysis, review, refactoring, testing, and documentation capabilities in containerized environments.

## Overview

OpenCode integrates with multiple AI providers to offer comprehensive coding assistance through both CLI and MCP interfaces. It supports session persistence and can create, modify, and analyze code files in your project workspace.

## Usage

### CLI Commands

#### Chat with OpenCode
```bash
# Basic chat
ship opencode chat "explain this code structure"

# Chat with specific model
ship opencode chat "generate a REST API handler" --model "openai/gpt-4o-mini"

# Chat with session persistence
ship opencode chat "My name is John" --continue --model "openai/gpt-4o-mini"
ship opencode chat "What's my name?" --continue --model "openai/gpt-4o-mini"

# Chat in specific directory
ship opencode chat "analyze the main.go file" --workdir ./src

# Ephemeral mode (don't persist files)
ship opencode chat "show me example code" --ephemeral
```

#### Generate Code
```bash
# Generate code with specific output
ship opencode generate "create a user model struct" --output models/user.go

# Generate in specific directory
ship opencode generate "REST API endpoints" --workdir ./api
```

#### Code Analysis
```bash
# Analyze specific file
ship opencode analyze main.go "explain the main function"

# Analyze with questions
ship opencode analyze config.yaml "what does this configuration do?"
```

#### Code Review
```bash
# Review current changes
ship opencode review

# Review specific target
ship opencode review --target src/handlers
```

#### Refactor Code
```bash
# Refactor with instructions
ship opencode refactor "extract common validation logic"

# Refactor specific files
ship opencode refactor "improve error handling" --files main.go,handlers.go
```

#### Generate Tests
```bash
# Generate tests
ship opencode test

# Generate specific test types with coverage
ship opencode test --type unit --coverage
```

#### Generate Documentation
```bash
# Generate markdown documentation
ship opencode document

# Generate HTML documentation
ship opencode document --format html --output-dir ./docs
```

#### Interactive Session
```bash
# Start interactive session
ship opencode interactive

# Interactive with specific model
ship opencode interactive --model "anthropic/claude-3-sonnet"
```

#### Batch Processing
```bash
# Process multiple files
ship opencode batch "*.go" "add error handling"
```

#### Version Information
```bash
# Get OpenCode version
ship opencode version
```

### MCP Integration

OpenCode is available through the MCP server interface with these tools:

#### opencode_run
Main tool for all OpenCode operations:

```json
{
  "name": "opencode_run",
  "arguments": {
    "message": "analyze this code and suggest improvements",
    "workdir": ".",
    "persist_files": true,
    "model": "openai/gpt-4o-mini",
    "session": "my-coding-session",
    "continue": false,
    "share": false,
    "agent": "code-reviewer"
  }
}
```

#### opencode_version
Get version information:

```json
{
  "name": "opencode_version",
  "arguments": {}
}
```

## Configuration

### Environment Variables

OpenCode automatically uses AI provider API keys from environment variables:

```bash
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export CLAUDE_API_KEY="sk-..."
export GEMINI_API_KEY="..."
export GROQ_API_KEY="gsk_..."
export OPENROUTER_API_KEY="sk-or-..."
```

### Available Models

Specify models using the format `provider/model`:

- **OpenAI**: `openai/gpt-4o-mini`, `openai/gpt-4o`, `openai/gpt-3.5-turbo`
- **Anthropic**: `anthropic/claude-3-sonnet`, `anthropic/claude-3-haiku`
- **Google**: `google/gemini-pro`, `google/gemini-flash`
- **Groq**: `groq/llama-70b`, `groq/mixtral-8x7b`

## Flags and Options

### Global Flags
- `--workdir, -w`: Working directory for operations (default: current directory)
- `--ephemeral`: Run in ephemeral mode - don't persist files to host
- `--session, -s`: Session ID for multi-turn conversations
- `--continue, -c`: Continue the last session

### Command-Specific Flags

**generate:**
- `--output, -o`: Output file for generated code

**refactor:**
- `--files, -f`: Specific files to refactor

**test:**
- `--type, -t`: Test type (unit, integration, e2e)
- `--coverage, -c`: Enable test coverage analysis

**document:**
- `--format`: Documentation format (markdown, html, pdf)
- `--output-dir, -d`: Output directory for documentation

**interactive:**
- `--model, -m`: Specific AI model to use

**chat:**
- `--model, -m`: AI model to use

## Session Persistence

OpenCode supports session persistence through:

- **Session Directory**: `~/.local/share/opencode/` (mounted in containers)
- **Continue Sessions**: Use `--continue` flag to resume last conversation
- **Named Sessions**: Use `--session "session-name"` for specific sessions

**Note**: Session memory functionality has limitations - conversations don't maintain context between calls.

## File Persistence

By default, OpenCode persists all file changes back to the host:

- **Persistent Mode** (default): Files created/modified in container are exported to host
- **Ephemeral Mode** (`--ephemeral`): Files stay in container only

## Examples

### Basic Code Generation
```bash
# Generate a web server
ship opencode generate "create a simple HTTP server in Go with health check endpoint" \
  --output server/main.go \
  --model "openai/gpt-4o-mini"
```

### Code Review Workflow
```bash
# Review current changes
ship opencode review --workdir ./src

# Get specific analysis
ship opencode analyze server.go "identify potential security issues"
```

### Documentation Generation
```bash
# Generate project documentation
ship opencode document \
  --format markdown \
  --output-dir ./docs \
  --workdir .
```

### Interactive Development
```bash
# Start interactive session with specific model
ship opencode interactive --model "anthropic/claude-3-sonnet" --workdir ./project
```

## Troubleshooting

### Model Selection Issues
If you encounter "model not found" errors, explicitly specify a model:
```bash
ship opencode chat "help me" --model "openai/gpt-4o-mini"
```

### Session Issues
- Named sessions (`--session "name"`) may not work reliably
- Use `--continue` for basic session continuation
- Session memory doesn't persist conversation context

### File Persistence Issues
- Ensure proper permissions in working directory
- Check Docker daemon is running (required for Dagger)
- Use `--ephemeral` flag if file persistence isn't needed

## Technical Details

- **Container**: Uses `opencode-ai` npm package in Node.js 18 container
- **Dagger Integration**: Full containerized execution with volume mounting
- **Session Storage**: Maps host `~/.local/share/opencode/` to container
- **Environment Variables**: Automatically passed to container for AI provider access