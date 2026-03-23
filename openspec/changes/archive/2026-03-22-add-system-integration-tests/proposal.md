## Why

The 6G Skill-Based Agentic Core involves complex asynchronous interactions between three microservices (AAIHF, ACRF, A-IGW). Manual verification is error-prone and slow. This change introduces a system-wide integration testing framework to ensure deterministic network signaling follows natural language intents reliably.

## What Changes

- **Establish `tests/` Directory:** A dedicated space for cross-service integration tests.
- **Implement System Integration Test:** A Go-based test that spins up all three services and verifies the "Fleet Wake-Up" end-to-end flow.
- **Create Test Utilities:** Shared functions for service orchestration, mock agent creation, and HTTP assertion.
- **Define CI Verification Mandate:** Establish the rule that tests must be run before any future implementation.

## Capabilities

### New Capabilities
- `automated-system-verification`: The ability to verify the entire "Intent -> Skill -> Signaling" pipeline automatically.

### Modified Capabilities
- None

## Impact

- **Project Structure:** Adds `/tests` directory.
- **Workflow:** Future tasks will require successful test runs as a prerequisite for completion.
- **Reliability:** Ensures that changes to one service (e.g., registry logic) do not break the end-to-end intent resolution.
