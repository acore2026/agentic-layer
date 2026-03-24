## ADDED Requirements

### Requirement: Path Diversity Skill Execution
The A-IGW SHALL support the `mcp://skill/reliability/path-diversity` skill by executing a sequence of reliability-related 3GPP service operations.

#### Scenario: Successful Path Diversity Trigger
- **WHEN** the `path-diversity` skill is invoked
- **THEN** the system SHALL mock the following calls in order: `NNF_Generic_Control`, `Nsmf_PDUSession_UpdateSMContext`, and `Nnef_TrafficInfluence_Create`.
