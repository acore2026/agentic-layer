## MODIFIED Requirements

### Requirement: Skill Discovery
The ACRF SHALL provide an HTTP GET `/discover?skill_id=...` endpoint to retrieve a `SkillProfile`. It SHALL support both exact `SkillID` matching and semantic matching based on the provided query string.

#### Scenario: Successful Semantic Discovery
- **WHEN** a GET request is made to `/discover` with a natural language query (e.g., "wake up devices")
- **THEN** the system SHALL return the most semantically relevant `SkillProfile` JSON with a `200 OK`, provided it exceeds the similarity threshold.

#### Scenario: Skill Not Found
- **WHEN** a GET request is made to `/discover` with a query that has no semantically similar matches above the threshold
- **THEN** the system SHALL return a `404 Not Found`.
