## MODIFIED Requirements

### Requirement: Autonomous Reasoner
The AAIHF SHALL utilize an `adk-go` agent that identifies the necessary sequence of tool calls (Discovery then Invocation) based on the user's prompt. It SHALL support multiple LLM providers (Gemini, Kimi) for increased reliability. The reasoner SHALL emit real-time events to the SSE broker at key lifecycle stages.

#### Scenario: Provider Selection
- **WHEN** `AGENTIC_LLM_PROVIDER` is set to "kimi"
- **THEN** the system SHALL initialize the `NewCoreAgent` using the Kimi provider instead of Gemini.

#### Scenario: Agent Discovery and Invocation Loop
- **WHEN** the agent receives a prompt for which it lacks a skill profile (e.g., "Wake up the fleet")
- **THEN** it SHALL call the `SearchSkill` tool to discover the URI from ACRF before calling the `ExecuteSkill` tool.

#### Scenario: Agent Discovery and Invocation Loop (Kimi)
- **WHEN** the Kimi-based agent receives a prompt for which it lacks a skill profile
- **THEN** it SHALL call the `SearchSkill` tool before calling the `ExecuteSkill` tool, exactly as the Gemini-based agent does.

#### Scenario: Real-time Event Emission
- **WHEN** the reasoning process starts
- **THEN** the agent SHALL emit a `reasoning_started` event containing the user's prompt.
- **WHEN** the reasoning process completes
- **THEN** the agent SHALL emit a `reasoning_completed` event containing the final response.
