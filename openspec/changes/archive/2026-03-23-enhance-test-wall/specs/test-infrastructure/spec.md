## ADDED Requirements

### Requirement: BDD Integration
The test infrastructure SHALL support Behavior-Driven Development (BDD) using the `godog` framework to map text-based scenarios to executable Go tests.

#### Scenario: Running BDD Tests
- **WHEN** the `go test` command is run for the BDD suite
- **THEN** it SHALL parse feature files and execute the corresponding step definitions, reporting success or failure.

### Requirement: Temporal Workflow Isolation Testing
The test infrastructure SHALL support testing Temporal workflows in isolation using the Temporal test suite without requiring a live Temporal server.

#### Scenario: Simulating Activity Failure
- **WHEN** a Temporal workflow test simulates a failure in `CallSMFActivity`
- **THEN** the test environment SHALL verify that `RollbackAMFActivity` is subsequently invoked.

### Requirement: Agent Tool Unit Testing
The test infrastructure SHALL include unit tests for the agent tools (`SearchSkill`, `ExecuteSkill`) using mock HTTP servers to simulate backend responses.

#### Scenario: Mocking Backend Errors
- **WHEN** an agent tool is tested against a mock server returning a 404
- **THEN** the tool SHALL return the appropriate error string without causing a panic or unhandled error.
