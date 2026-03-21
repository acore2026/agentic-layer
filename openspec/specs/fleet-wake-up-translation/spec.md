## ADDED Requirements

### Requirement: Bootstrapping Registration
The A-IGW SHALL register the `mcp://skill/device/fleet-update` skill profile with the ACRF on startup.

#### Scenario: Successful Self-Registration
- **WHEN** the A-IGW service starts
- **THEN** it SHALL send a POST `/register` request to the configured ACRF endpoint.

### Requirement: Skill Invocation
The A-IGW SHALL expose an HTTP POST `/invoke` endpoint that triggers the deterministic sequence of 3GPP service operations.

#### Scenario: Successful Fleet Wake-Up Invocation
- **WHEN** a POST request is made to `/invoke` with a payload identifying the "Fleet Wake-Up" skill
- **THEN** the system SHALL log the execution of three downstream API calls: `Namf_MT_EnableUEReachability`, `Nsmf_PDUSession_UpdateSMContext`, and `Nnef_AFSessionWithQoS_Create`.
