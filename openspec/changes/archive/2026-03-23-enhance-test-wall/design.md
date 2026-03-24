## Context

We currently rely heavily on a bash script (`scripts/verify_system.sh`) and a basic Go integration test (`tests/system_integration_test.go`) to verify the system. As we add more complex logic, such as Temporal workflows and BDD specifications, we need a robust testing strategy that covers multiple layers: unit tests for Agent tools, isolated workflow tests for Temporal, and Behavior-Driven Development (BDD) tests that directly map to our OpenSpec files.

## Goals / Non-Goals

**Goals:**
- Implement BDD testing using `github.com/cucumber/godog` to map our `.md` specs to tests.
- Implement Temporal Workflow isolated testing using `go.temporal.io/sdk/testsuite`.
- Expand Go integration tests to cover failure paths and edge cases (e.g. LLM hallucination handling, missing skills).
- Add unit tests for `internal/agent/tools.go` to test schema injection and parameter handling with a mock LLM.

**Non-Goals:**
- Fully rewriting existing integration tests (we will just expand them).
- Deploying a CI/CD pipeline (this change focuses only on adding the tests to the codebase).

## Decisions

- **Temporal Workflow Test Suite**: Using the official `go.temporal.io/sdk/testsuite` ensures we can test the `FleetWakeUpWorkflow` in memory, simulating activity failures and ensuring `RollbackAMFActivity` is called correctly.
- **BDD with Godog**: `godog` is the standard BDD framework for Go. It allows us to parse Gherkin-style `Feature`/`Scenario` text and execute Go functions. We will write a `.feature` file that reflects our OpenSpec scenarios.
- **Agent Unit Testing**: The Agent tools (`SearchSkill`, `ExecuteSkill`) rely on HTTP calls. We will use Go's `httptest.Server` to mock the ACRF and IGW responses, ensuring the tools handle various JSON responses and errors correctly without requiring the full stack to be running.

## Risks / Trade-offs

- **[Risk] Slower test execution**: Adding more layers of testing (especially BDD and Workflow tests) can increase `go test` duration.
  - *Mitigation*: We will use `t.Parallel()` where appropriate and rely on in-memory mocks (like Temporal's testsuite and `httptest`) rather than spinning up full network ports whenever possible.
- **[Trade-off] Maintenance Overhead**: BDD tests require maintaining both the plain text `.feature` files and their Go glue code.
  - *Rationale*: The clarity provided by mapping directly to OpenSpec `spec.md` files outweighs the cost, especially for complex intent-to-skill resolution logic.
