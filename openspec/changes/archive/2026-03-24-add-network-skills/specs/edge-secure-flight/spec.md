## ADDED Requirements

### Requirement: Secure Flight Skill Execution
The A-IGW SHALL support the `mcp://skill/edge/secure-flight` skill by executing a sequence of edge and location-related 3GPP service operations.

#### Scenario: Successful Secure Flight Trigger
- **WHEN** the `secure-flight` skill is invoked
- **THEN** the system SHALL mock the following calls in order: `Nnef_TrafficInfluence_Create`, `Nnef_EventExposure_Subscribe`, and `Ngmlc_Location_ProvideLocation`.
