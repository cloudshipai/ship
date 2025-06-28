# Native Dagger LLM Integration in Ship CLI

## Overview

Ship CLI now uses **Dagger's native LLM features** for AI-powered infrastructure investigation. This provides better integration, caching, and performance compared to custom API implementations.

## Architecture

### 1. Native Dagger LLM Module (`llm_dagger.go`)
Uses Dagger's built-in LLM primitives:
```go
llm := client.LLM(dagger.LLMOpts{
    Model: "gpt-4",
})

response := llm.
    WithSystemPrompt("You are a Steampipe expert...").
    WithPrompt(userPrompt).
    Sync(ctx).
    LastReply(ctx)
```

### 2. Fallback Custom Module (`llm.go`) 
Legacy implementation using curl in containers for when Dagger LLM is unavailable.

## Key Features

### Native Dagger LLM Benefits
- **Automatic Caching**: Dagger caches LLM responses for identical prompts
- **Token Usage Tracking**: Built-in `TokenUsage()` method for cost monitoring
- **Provider Abstraction**: Seamless switching between OpenAI, Anthropic, etc.
- **Error Handling**: Better retry logic and error messages
- **Streaming Support**: Can stream responses (when needed)

### Available Methods
```go
// Core LLM methods in Dagger
llm.WithPrompt(string)           // Set user prompt
llm.WithSystemPrompt(string)     // Set system prompt
llm.WithPromptFile(*File)        // Load prompt from file
llm.Sync(ctx)                    // Execute the LLM
llm.LastReply(ctx)              // Get the response
llm.History(ctx)                // Get conversation history
llm.TokenUsage()                // Track token usage
llm.Model(ctx)                  // Get model name
llm.Provider(ctx)               // Get provider name
```

## Usage Examples

### Basic Usage
```bash
# Use default rule-based system (no API needed)
ship ai-investigate --prompt "Find security issues" --execute

# Use Dagger LLM with OpenAI
export OPENAI_API_KEY=your-key
ship ai-investigate --prompt "Find security issues" \
  --llm-provider openai --model gpt-4 --execute

# Use Dagger LLM with Anthropic
export ANTHROPIC_API_KEY=your-key  
ship ai-investigate --prompt "Analyze cost optimization" \
  --llm-provider anthropic --model claude-3 --execute
```

### How It Works

1. **Query Generation**
   ```go
   llm := client.LLM(dagger.LLMOpts{Model: "gpt-4"})
   llm.WithPrompt("Convert to SQL: " + naturalLanguage).Sync(ctx)
   ```

2. **Results Analysis**
   ```go
   llm.WithSystemPrompt("You are an infrastructure expert").
       WithPrompt("Analyze: " + queryResults).
       Sync(ctx)
   ```

3. **Investigation Planning**
   ```go
   llm.WithPrompt("Create investigation plan for: " + objective).
       Sync(ctx)
   ```

## Implementation Status

✅ **Completed**:
- Native Dagger LLM module (`DaggerLLMModule`)
- Integration with ai-investigate command
- Fallback to rule-based system
- System prompts for better responses

⚠️ **In Progress**:
- JSON response parsing for structured data
- Token usage tracking and reporting
- Streaming response support

❌ **TODO**:
- Cache management for responses
- Multi-turn conversations
- Custom model fine-tuning support
- Prompt templates library

## Environment Variables

```bash
# OpenAI
export OPENAI_API_KEY=sk-...

# Anthropic  
export ANTHROPIC_API_KEY=sk-ant-...

# Azure OpenAI
export AZURE_OPENAI_API_KEY=...
export AZURE_OPENAI_ENDPOINT=...

# Custom endpoint
export LLM_ENDPOINT=http://localhost:11434  # For Ollama
```

## Advanced Features

### Token Usage Tracking
```go
tokenUsage := llm.TokenUsage()
promptTokens, _ := tokenUsage.PromptTokens(ctx)
completionTokens, _ := tokenUsage.CompletionTokens(ctx)
totalTokens, _ := tokenUsage.TotalTokens(ctx)
```

### Conversation History
```go
history, _ := llm.History(ctx)
for _, message := range history {
    fmt.Println(message)
}
```

### Custom System Prompts
```go
llm.WithSystemPrompt(`
You are a Steampipe SQL expert specializing in AWS infrastructure.
Always return valid SQL queries that follow Steampipe conventions.
Include comments explaining what each query does.
`)
```

## Debugging

Enable Dagger debug logs:
```bash
export DAGGER_LOG_LEVEL=debug
ship ai-investigate --prompt "test" --llm-provider openai
```

Check LLM configuration:
```go
model, _ := llm.Model(ctx)
provider, _ := llm.Provider(ctx)
fmt.Printf("Using %s with %s\n", provider, model)
```

## Performance Tips

1. **Use specific prompts** - More specific prompts = better results
2. **Leverage caching** - Dagger caches identical prompts automatically  
3. **Batch operations** - Group related queries together
4. **Monitor tokens** - Track usage to control costs

## Future Enhancements

1. **Prompt Library**: Pre-built prompts for common tasks
2. **Fine-tuned Models**: Custom models for Steampipe SQL
3. **Multi-Agent**: Multiple LLMs working together
4. **Feedback Loop**: Learn from user corrections

## Example: Full Investigation Flow

```go
// 1. Create investigation plan
plan := llm.WithPrompt("Plan investigation for: " + objective).Sync(ctx)

// 2. Generate SQL queries  
queries := llm.WithPrompt("Convert to Steampipe SQL: " + plan).Sync(ctx)

// 3. Analyze results
insights := llm.WithPrompt("Analyze security issues in: " + results).Sync(ctx)

// 4. Generate recommendations
fixes := llm.WithPrompt("Suggest fixes for: " + insights).Sync(ctx)
```

This native integration makes Ship CLI's AI features more powerful, reliable, and easier to extend.