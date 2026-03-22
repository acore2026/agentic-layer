## ADDED Requirements

### Requirement: Intent Processing API
The AAIHF SHALL expose an HTTP POST `/intent` endpoint that accepts a natural language prompt in JSON format.

#### Scenario: Successful Intent Resolution
- **WHEN** a POST request is made to `/intent` with the prompt "Wake up my fleet"
- **THEN** the system SHALL resolve the intent and trigger the corresponding `fleet-update` skill.

### Requirement: Autonomous Reasoner
The AAIHF SHALL utilize an `adk-go` agent that identifies the necessary sequence of tool calls (Discovery then Invocation) based on the user's prompt.

#### Scenario: Agent Discovery and Invocation Loop
- **WHEN** the agent receives a prompt for which it lacks a skill profile
- **THEN** it SHALL call the `SearchSkill` tool before calling the `ExecuteSkill` tool.
