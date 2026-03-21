## Why

Standard 3GPP discovery (NRF) relies on rigid, identity-based matching that breaks when node identities change. This project implements a **Skill-Based Agentic Architecture** for 6G to decouple service logic from physical topology, enabling intent-driven networking where AI agents can discover capabilities (Skills) via semantic URIs rather than static IPs.

## What Changes

- **Bootstrap 6G Agentic Core Module:** Initialize the Go module and directory structure.
- **Implement ACRF (Agentic Capability Repository Function):** Create a dynamic skill registry with an in-memory map for the MVP.
- **Implement A-IGW (Interworking Gateway):** Create a "Universal Translator" that maps the `fleet-update` skill to deterministic 3GPP API calls.
- **Define Unified Skill Profile:** Establish a polymorphic data model for skills across Device, Network, and App domains.

## Capabilities

### New Capabilities
- `skill-registry`: Core ACRF capability for registering and discovering agentic skills via HTTP.
- `fleet-wake-up-translation`: Specific IGW capability that translates the `mcp://skill/device/fleet-update` skill into a sequence of AMF, SMF, and NEF service operations.
- `agentic-models`: Shared polymorphic data structures for Unified Agentic Skill Profiles.

### Modified Capabilities
- None (Initial bootstrap)

## Impact

- **New Services:** `acrf` (port 8080), `igw-fleet` (port 8081).
- **Dependencies:** `github.com/google/adk-go` (for future AAIHF integration).
- **Architecture:** Introduces a decoupled AI layer sitting on top of the legacy 5G Core (free5gc).
