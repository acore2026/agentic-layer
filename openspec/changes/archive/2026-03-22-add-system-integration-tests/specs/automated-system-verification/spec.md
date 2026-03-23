## ADDED Requirements

### Requirement: End-to-End Flow Verification
The system SHALL provide an automated test that validates the path from a natural language intent to a deterministic service operation.

#### Scenario: Successful Fleet Wake-Up Integration
- **WHEN** all services are running and an intent "Wake up the fleet" is posted to AAIHF
- **THEN** the test SHALL confirm that ACRF received a discovery request and A-IGW received an invocation request.

### Requirement: Service Orchestration
The test framework SHALL automate the lifecycle (start/stop) of all dependent microservices.

#### Scenario: Clean Test Environment
- **WHEN** an integration test finishes
- **THEN** all ports used by the microservices SHALL be released.

### Requirement: Mock Agent Verification
The system SHALL support a `MockCoreAgent` that enables testing the microservice plumbing without external LLM API calls.

#### Scenario: Deterministic Reasoning Test
- **WHEN** the `MockCoreAgent` is used in a test
- **THEN** it SHALL consistently produce the same sequence of tool calls for a given intent.
