## ADDED Requirements

### Requirement: Turbo Mode Skill Execution
The A-IGW SHALL support the `mcp://skill/qos/turbo-mode` skill by executing a sequence of QoS-related 3GPP service operations.

#### Scenario: Successful Turbo Mode Trigger
- **WHEN** the `turbo-mode` skill is invoked
- **THEN** the system SHALL mock the following calls in order: `Nnef_AFSessionWithQoS_Create`, `Nnef_ChargeableParty_Create`, and `Npcf_PolicyAuthorization_Update`.
