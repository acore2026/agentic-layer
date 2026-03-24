## MODIFIED Requirements

### Requirement: Skill Invocation
The A-IGW SHALL expose an HTTP POST `/invoke` endpoint that triggers the deterministic sequence of 3GPP service operations asynchronously using a Temporal workflow.

#### Scenario: Successful Fleet Wake-Up Invocation
- **WHEN** a POST request is made to `/invoke` with a payload identifying the "Fleet Wake-Up" skill
- **THEN** the system SHALL start a Temporal workflow, return a `202 Accepted` response with the workflow ID, and asynchronously log the execution of three downstream API calls: `Namf_MT_EnableUEReachability`, `Nsmf_PDUSession_UpdateSMContext`, and `Nnef_AFSessionWithQoS_Create`.
