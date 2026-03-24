## Context

The A-IGW and ACRF currently only support a single fleet-update skill. We need to expand this to support QoS, Reliability, and Edge-focused skills as defined in the 3GPP 6G Agentic Architecture. This expansion will demonstrate the flexibility of the Skill-Based Agentic Architecture in handling diverse network requirements.

## Goals / Non-Goals

**Goals:**
- Implement handlers for 3 new network skills in A-IGW.
- Automate the registration of these skills in ACRF on startup.
- Use semantic descriptions to enable natural language discovery via AAIHF.

**Non-Goals:**
- Actual integration with a live 5G Core (Free5GC). We will continue to mock the 3GPP service operations.
- Dynamic skill registration via UI (this is still bootstrap-based).

## Decisions

- **Routing Logic**: Use a `switch` statement in `internal/translator/translator.go` to dispatch skill IDs to their respective mock signaling sequences.
- **Mock Sequences**: Each skill will have a hardcoded sequence of log statements representing 3GPP API calls (e.g., `Nnef_AFSessionWithQoS_Create`).
- **Bootstrap Registration**: ACRF `main.go` will be updated to perform POST requests to its own `/register` endpoint for all default skills.

## Risks / Trade-offs

- **[Risk] Scalability of switch statement**: As we add more skills, the `switch` statement might become bloated. 
  - *Mitigation*: For the MVP, it's sufficient. If we scale beyond 10-20 skills, we should move to a map-based handler registry.
- **[Trade-off] Mocking vs Reality**: The mock sequences don't reflect the true state of a 5G core.
  - *Rationale*: The goal is to prove the Agentic Layer's ability to discover and trigger these sequences, not the sequences themselves.
