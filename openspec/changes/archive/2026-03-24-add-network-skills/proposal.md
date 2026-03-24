## Why

The current `6g-agentic-core` MVP only implements a single Agentic Skill (`mcp://skill/device/fleet-update`). To fully demonstrate the 3GPP 6G Agentic Architecture proposal, we need to expand the system's capabilities to support more diverse network operations such as QoS optimization, reliability enhancement, and edge computing management.

## What Changes

- **New Skill Implementations**: Add three new Agentic Skills to the ACRF registry and A-IGW executor:
    - **Turbo Mode (QoS)**: `mcp://skill/qos/turbo-mode`
    - **Resiliency (Reliability)**: `mcp://skill/reliability/path-diversity`
    - **Secure Drone Corridor (Edge)**: `mcp://skill/edge/secure-flight`
- **Executor Refactor**: Update the A-IGW `translator.go` to route these new skill IDs to specific 3GPP mock signaling sequences.
- **Bootstrap Registry**: Update ACRF `main.go` to automatically register these new skills on startup with semantic descriptions for better discovery.

## Capabilities

### New Capabilities
- `qos-optimization`: Support for dynamic Quality of Service adjustments via Agentic Skills.
- `reliability-enhancement`: Support for path diversity and zero-interruption network configurations.
- `edge-secure-flight`: Support for traffic influence and location-based drone corridor security.

### Modified Capabilities
- `skill-registry`: Expanded to bootstrap more default skills.
- `temporal-skill-execution`: (If using Temporal) Workflows will be added for the new signaling sequences. For this proposal, we focus on the mock execution.

## Impact

- **ACRF**: Will have 4 skills registered on startup instead of 1.
- **A-IGW**: New mock handlers for 3GPP signaling sequences.
- **AAIHF**: Can now resolve intents like "Enable Turbo Mode" or "Secure Drone Corridor" using the Semantic Matching Engine.
