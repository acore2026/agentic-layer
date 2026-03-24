## MODIFIED Requirements

### Requirement: Autonomous Reasoner
The AAIHF SHALL utilize an `adk-go` agent that identifies the necessary sequence of tool calls (Discovery then Invocation) based on the user's prompt. It SHALL support multiple LLM providers (Gemini, Kimi) for increased reliability.

#### Scenario: Provider Selection
- **WHEN** `AGENTIC_LLM_PROVIDER` is set to "kimi"
- **THEN** the system SHALL initialize the `NewCoreAgent` using the Kimi provider instead of Gemini.

#### Scenario: Agent Discovery and Invocation Loop (Kimi)
- **WHEN** the Kimi-based agent receives a prompt for which it lacks a skill profile
- **THEN** it SHALL call the `SearchSkill` tool before calling the `ExecuteSkill` tool, exactly as the Gemini-based agent does.
