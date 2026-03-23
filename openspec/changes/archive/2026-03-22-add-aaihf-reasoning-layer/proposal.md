## Why

Standard network control requires hard-coded logic for every intent. This change introduces the **AAIHF (Agentic AI Host Function)**, which uses LLM-based reasoning via `adk-go` to dynamically resolve natural language intents into executable network skills by interacting with the ACRF and A-IGW.

## What Changes

- **Implement AAIHF Microservice:** Create the `cmd/aaihf` entry point.
- **Integrate `adk-go` Framework:** Use Google's `adk-go` to define an autonomous agent.
- **Implement Agentic Tools:** 
    - `SearchSkill`: Tool for semantic discovery via ACRF.
    - `ExecuteSkill`: Tool for deterministic invocation via A-IGW.
- **Intent Processing Endpoint:** Expose an HTTP POST `/intent` endpoint to receive natural language prompts.

## Capabilities

### New Capabilities
- `intent-resolution`: The ability for the system to process natural language prompts and map them to skills.
- `agent-tooling`: Standardized tools for the AI agent to interact with internal 6G services (ACRF/IGW).

### Modified Capabilities
- None

## Impact

- **New Service:** `aaihf` (port 18082).
- **Dependencies:** `github.com/google/adk-go`.
- **Infrastructure:** Requires `GEMINI_API_KEY` for LLM inference. 
