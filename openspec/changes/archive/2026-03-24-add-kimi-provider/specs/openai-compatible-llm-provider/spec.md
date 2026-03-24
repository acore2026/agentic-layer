## ADDED Requirements

### Requirement: model.LLM Implementation
The system SHALL provide a struct that implements the `model.LLM` interface for OpenAI-compatible REST APIs.

#### Scenario: Successful Text Generation
- **WHEN** `GenerateContent` is called with a simple text prompt
- **THEN** the provider SHALL return a valid `model.LLMResponse` containing the generated text.

### Requirement: Tool-Calling Support
The OpenAI provider SHALL map `adk-go` tool definitions to the OpenAI `tools` format and handle tool execution events.

#### Scenario: Successful Tool Call
- **WHEN** the agent decides to use a tool (e.g., SearchSkill)
- **THEN** the OpenAI provider SHALL generate a response with the correct `tool_calls` payload.

### Requirement: Configurable Base URL
The system SHALL allow the base URL to be overridden via environment variables to support different providers like Moonshot/Kimi.

#### Scenario: Moonshot API Connection
- **WHEN** `AGENTIC_KIMI_BASE_URL` is set to Moonshot's endpoint
- **THEN** the provider SHALL send requests to that URL with the appropriate `Authorization` header.
