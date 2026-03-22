## Context

The ACRF and A-IGW are already implemented. We now need a reasoning entity (AAIHF) that can take a natural language intent like "Wake up the fleet" and use these services as tools. We'll use the `adk-go` framework to build this logic.

## Goals / Non-Goals

**Goals:**
- **Intent-Driven Interaction:** Enable processing of natural language.
- **Service Integration:** Successfully use ACRF for discovery and A-IGW for execution.
- **Adoption of `adk-go`:** Follow idiomatic patterns for tool definition and agent loops.

**Non-Goals:**
- **Persistent Sessions:** We won't store conversation history across service restarts for this prototype.
- **Custom Model Support:** We will focus on Gemini (native to `adk-go`) rather than implementing a custom OpenAI adapter for now.

## Decisions

- **AAIHF as a Microservice:** Build a standalone HTTP server that acts as the entry point for intents.
  - *Rationale:* Decouples AI reasoning from external user interfaces (like a CLI or Dashboard).
- **Tool-First Design:** Map ACRF and A-IGW endpoints directly to `adk-go` tools.
  - *Rationale:* The LLM should decide *when* to search and *when* to execute based on its reasoning.
- **Environment-Based Config:** Use `ACRF_URL`, `IGW_URL`, and `GEMINI_API_KEY` env vars.
  - *Rationale:* Ensures portability and keeps secrets out of code.

## Risks / Trade-offs

- **[Risk] Hallucination:** The agent might try to invoke non-existent skills.
  - *Mitigation:* The `SearchSkill` tool must return a definitive schema, and the agent must be instructed to only execute skills it has discovered.
- **[Trade-off] Performance:** LLM reasoning adds latency to network control.
  - *Mitigation:* This is acceptable for a prototype; in production, we would use specialized, smaller models or cached reasoning.
