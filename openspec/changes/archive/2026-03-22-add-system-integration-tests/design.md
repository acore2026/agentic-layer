## Context

The current system relies on three separate microservices that communicate over HTTP. To verify the full flow, we need a test environment that can manage the lifecycle of these services and simulate the interaction between them.

## Goals / Non-Goals

**Goals:**
- **Deterministic Orchestration:** Automate the startup and teardown of ACRF, A-IGW, and AAIHF for testing.
- **Mock-Based Reasoning:** Use a `MockCoreAgent` during integration tests to avoid LLM dependency, cost, and latency.
- **Service Verification:** Assert that the A-IGW correctly registers its skill with ACRF on boot.
- **End-to-End Validation:** Confirm that calling the AAIHF `/intent` endpoint results in the expected signaling sequence from the IGW.

**Non-Goals:**
- **External 5G Core Integration:** We will not use a real 5G core (e.g., free5gc) in this test suite; we will assert based on IGW logs/stubs.
- **UI Testing:** This framework focuses on backend microservice interaction.

## Decisions

- **Go `testing` Package:** Use standard Go `testing` and `httptest` where appropriate.
- **In-Process Service Lifecycle:** Run the services as goroutines within the test process to allow for easy cleanup and shared memory for stubs if needed.
- **Environment Manipulation:** Dynamically set `AGENTIC_ACRF_URL`, `AGENTIC_IGW_URL`, and `AGENTIC_AAIHF_PORT` during test setup.
- **Shared Test Utilities:** Create a `testutil` package to avoid code duplication across different integration scenarios.

## Risks / Trade-offs

- **[Risk] Port Collisions:** Parallel tests might try to bind to the same ports.
  - *Mitigation:* Use random available ports or run integration tests sequentially.
- **[Trade-off] Realism vs. Cost:** Using a Mock Agent instead of Gemini means we aren't testing the actual LLM "reasoning."
  - *Mitigation:* The goal of *integration* tests is to verify the plumbing and protocol (Skill URI, JSON mapping). The reasoning quality is better tested via specialized "eval" datasets.
