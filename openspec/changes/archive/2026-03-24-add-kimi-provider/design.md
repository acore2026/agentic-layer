## Context

The `adk-go` framework currently lacks a built-in OpenAI model provider in the version we are using. To integrate Kimi K2.5, we must implement the `model.LLM` interface from `google.golang.org/adk/model`.

## Goals / Non-Goals

**Goals:**
- **Custom Model Implementation**: Create a struct that implements `model.LLM` using standard Go HTTP and OpenAI's protocol.
- **Provider Parity**: Support text generation and tool-calling (function calling) to match Gemini's capabilities.
- **Dynamic Switching**: Allow the system to choose between Gemini and Kimi via the `AGENTIC_LLM_PROVIDER` env var.

**Non-Goals:**
- **Streaming Support**: For this initial implementation, we will focus on non-streaming responses for simplicity.
- **Full OpenAI SDK Integration**: We will implement a minimal HTTP-based client rather than importing a large external SDK if possible, or use a lightweight one.

## Decisions

- **Custom HTTP Implementation**: We will build a lightweight OpenAI-compatible wrapper around `net/http` to satisfy the `model.LLM` interface.
  - *Rationale*: Keeps the dependency tree small and allows precise control over the mapping between `adk-go` events and OpenAI payloads.
- **New Package `internal/agent/openai`**: House the custom provider code here.
- **Configuration Update**: 
  - `AGENTIC_KIMI_API_KEY`: API Key for Moonshot.
  - `AGENTIC_KIMI_BASE_URL`: `https://api.moonshot.cn/v1`.
  - `AGENTIC_KIMI_MODEL`: `kimi-k2.5`.
  - `AGENTIC_LLM_PROVIDER`: `gemini` | `kimi`.

## Risks / Trade-offs

- **[Risk] Function Calling Differences**: OpenAI and Gemini have slightly different tool-calling schemas.
  - *Mitigation*: The custom provider must carefully map `adk-go` tool definitions to the OpenAI format.
- **[Trade-off] Maintenance**: We are now responsible for the OpenAI provider code since it's not a part of the core `adk-go` library.
  - *Mitigation*: Keep the implementation simple and focused on the `model.LLM` interface methods (`GenerateContent`).
