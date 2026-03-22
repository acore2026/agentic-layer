## ADDED Requirements

### Requirement: SearchSkill Tool
The AAIHF SHALL define a `SearchSkill` tool that wraps the ACRF `/discover` endpoint to find skill profiles by ID.

#### Scenario: Tool Returns Skill Profile
- **WHEN** the `SearchSkill` tool is called with `skill_id="mcp://skill/device/fleet-update"`
- **THEN** it SHALL return the JSON SkillProfile from the ACRF.

### Requirement: ExecuteSkill Tool
The AAIHF SHALL define an `ExecuteSkill` tool that wraps the A-IGW `/invoke` endpoint to execute discovered skills.

#### Scenario: Tool Triggers Invocation
- **WHEN** the `ExecuteSkill` tool is called with a valid `skill_id`
- **THEN** it SHALL return the execution status from the A-IGW.
