## Why

The current Gemini API integration is unstable, frequently returning TLS handshake timeouts in the current network environment. To ensure the reliability of the 6G Agentic Core reasoning layer, we need an alternative LLM provider. Kimi K2.5 offers an OpenAI-compatible API that can serve as a robust fallback or primary reasoning engine.

## What Changes

- **Implement OpenAI-Compatible Provider for `adk-go`**: Create a new `model.LLM` implementation that supports OpenAI-compatible endpoints (specifically Moonshot/Kimi).
- **Add Kimi Configuration to `internal/config`**: Include `AGENTIC_KIMI_API_KEY` and Moonshot base URL in the centralized configuration.
- **Update AAIHF to Support Provider Selection**: Modify `internal/agent/agent.go` to allow switching between Gemini and Kimi based on environment variables.
- **Update `.env.example`**: Add placeholders for Kimi-related credentials.

## Capabilities

### New Capabilities
- `openai-compatible-llm-provider`: A generic LLM provider for `adk-go` that can connect to any OpenAI-compatible API (Kimi, DeepSeek, etc.).

### Modified Capabilities
- `intent-resolution`: Update the AAIHF to leverage the new provider for more stable reasoning.

## Impact

- **Dependencies**: May require an OpenAI Go client (e.g., `sashabaranov/go-openai`) or a custom HTTP implementation within the provider.
- **Reliability**: Improves system uptime by providing a secondary reasoning path.
- **Configuration**: Introduces new environment variables for Kimi integration.
