## MODIFIED Requirements

### Requirement: Bootstrapping Registration
The system SHALL automatically register a default set of network skills with the ACRF on startup.

#### Scenario: Successful Bootstrap
- **WHEN** the ACRF starts up
- **THEN** it SHALL register `fleet-update`, `turbo-mode`, `path-diversity`, and `secure-flight` skills.
