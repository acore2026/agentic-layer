## ADDED Requirements

### Requirement: Skill Registration
The ACRF SHALL provide an HTTP POST `/register` endpoint to save a `SkillProfile` to an in-memory registry.

#### Scenario: Successful Registration
- **WHEN** a POST request is made to `/register` with a valid `SkillProfile` JSON payload
- **THEN** the system SHALL return a `200 OK` and store the profile.

### Requirement: Skill Discovery
The ACRF SHALL provide an HTTP GET `/discover?skill_id=...` endpoint to retrieve a `SkillProfile` by its `SkillID`.

#### Scenario: Successful Discovery
- **WHEN** a GET request is made to `/discover` with a matching `skill_id` query parameter
- **THEN** the system SHALL return the corresponding `SkillProfile` JSON with a `200 OK`.

#### Scenario: Skill Not Found
- **WHEN** a GET request is made to `/discover` with a `skill_id` that does not exist
- **THEN** the system SHALL return a `404 Not Found`.
