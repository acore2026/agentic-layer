## ADDED Requirements

### Requirement: Skill Registration
The ACRF SHALL provide an HTTP POST `/register` endpoint to save a `SkillProfile` to an in-memory registry.

#### Scenario: Successful Registration
- **WHEN** a POST request is made to `/register` with a valid `SkillProfile` JSON payload
- **THEN** the system SHALL return a `200 OK` and store the profile.

### Requirement: Bootstrapping Registration
The system SHALL automatically register a default set of network skills with the ACRF on startup.

#### Scenario: Successful Bootstrap
- **WHEN** the ACRF starts up
- **THEN** it SHALL register `fleet-update`, `turbo-mode`, `path-diversity`, and `secure-flight` skills.

### Requirement: Skill Discovery
The ACRF SHALL provide an HTTP GET `/discover?skill_id=...` endpoint to retrieve a `SkillProfile`. It SHALL support both exact `SkillID` matching and semantic matching based on the provided query string.

#### Scenario: Successful Semantic Discovery
- **WHEN** a GET request is made to `/discover` with a natural language query (e.g., "wake up devices")
- **THEN** the system SHALL return the most semantically relevant `SkillProfile` JSON with a `200 OK`, provided it exceeds the similarity threshold.

#### Scenario: Successful Identity Discovery
- **WHEN** a GET request is made to `/discover` with a matching `skill_id` query parameter (exact ID)
- **THEN** the system SHALL return the corresponding `SkillProfile` JSON with a `200 OK`.

#### Scenario: Skill Not Found
- **WHEN** a GET request is made to `/discover` with a query that has no semantically similar matches above the threshold and no exact ID match
- **THEN** the system SHALL return a `404 Not Found`.
