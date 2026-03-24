## Why

Currently, our testing relies heavily on a single manual integration bash script (`scripts/verify_system.sh`) and a couple of basic Go tests. While this has helped us identify bugs, it's brittle and insufficient as the project grows. We need a comprehensive "test wall" that includes Temporal Workflow tests, expanded Integration Tests, Behavior-Driven Development (BDD) tests using `godog`, and Unit Tests for the Agent tools. This will prevent regressions and ensure system reliability.

## What Changes

- **Temporal Workflow Test Suite**: Implement `go.temporal.io/sdk/testsuite` to test `FleetWakeUpWorkflow` in isolation, including rollback logic.
- **Expanded Go Integration Tests**: Enhance `tests/system_integration_test.go` to cover edge cases like LLM hallucinations, missing skills, and simulated backend failures.
- **BDD Integration (`godog`)**: Introduce the `cucumber/godog` framework to map OpenSpec specifications directly to executable Go tests.
- **Agent Unit Testing**: Add unit tests for `internal/agent/tools.go` with mock LLM providers to verify parameter injection and JSON parsing securely.

## Capabilities

### New Capabilities
- `test-infrastructure`: Establishes the testing frameworks and suites across the project (Temporal testsuite, godog, agent unit tests).

### Modified Capabilities
- `automated-system-verification`: Expand the scope of automated testing beyond simple positive integration paths to include edge cases, BDD, and isolated component testing.

## Impact

- **Dependencies**: Introduces `github.com/cucumber/godog` and `go.temporal.io/sdk/testsuite` to the project module.
- **Project Structure**: Adds new test files in `tests/`, `internal/agent/`, and potentially `internal/translator/temporal_skills/`.
- **Reliability**: Massively improves confidence in code changes by establishing a multi-layered testing strategy.