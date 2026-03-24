## MODIFIED Requirements

### Requirement: End-to-End Flow Verification
The system SHALL provide an automated test that validates the path from a natural language intent to a deterministic service operation. The test suite SHALL also cover failure paths, edge cases, and graceful degradation.

#### Scenario: Successful Fleet Wake-Up Integration
- **WHEN** all services are running and an intent "Wake up the fleet" is posted to AAIHF
- **THEN** the test SHALL confirm that ACRF received a discovery request and A-IGW received an invocation request.

#### Scenario: Handling Missing Skills
- **WHEN** an intent is posted for a capability that does not exist in the ACRF registry
- **THEN** the system SHALL return a graceful failure message to the user without crashing.

#### Scenario: Handling Infrastructure Failures
- **WHEN** an intent is posted but the A-IGW cannot execute the workflow (e.g., Temporal server is down)
- **THEN** the system SHALL report the execution failure back to the user via the LLM agent.
