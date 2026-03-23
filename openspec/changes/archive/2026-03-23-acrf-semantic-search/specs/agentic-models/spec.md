## MODIFIED Requirements

### Requirement: Unified Agentic Skill Profile Structure
The system SHALL define a `SkillProfile` struct that includes a common header (including a mandatory `Description` field) and domain-specific containers for Device, Network, and App domains.

#### Scenario: Profile with Description
- **WHEN** a `SkillProfile` is initialized
- **THEN** it MUST include a `Description` string that provides a natural language summary of the skill's capability for semantic embedding.
