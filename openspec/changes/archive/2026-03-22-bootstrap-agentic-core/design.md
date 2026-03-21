## Context

The current 5G SBA (Service-Based Architecture) lacks a mechanism for AI-driven intent resolution. This design introduces the `ACRF` and `A-IGW` as the foundation for a 6G agentic layer. We are building a Go-based system that abstracts 3GPP network functions into "Skills."

## Goals / Non-Goals

**Goals:**
- **Decoupled Architecture:** Ensure the Registry (ACRF) and Translator (IGW) are independent microservices.
- **Polymorphic Modeling:** Support the `Unified Agentic Skill Profile` as defined in 3GPP S2-2600222.
- **Bootstrapping Registration:** Enable the IGW to self-register its skills with the ACRF on startup.
- **Deterministic Translation:** Ensure skill invocation results in a logged sequence of standard 3GPP service operations.

**Non-Goals:**
- **Semantic Matching:** Vector embeddings and LLM reasoning (AAIHF) are out of scope for this bootstrap phase.
- **Persistence:** The ACRF will use an in-memory map instead of a database for the MVP.
- **Real 3GPP Integration:** We will use log stubs instead of actual OpenAPI calls to `free5gc`.

## Decisions

- **Go Standard Layout:** Follow `/cmd`, `/internal`, `/pkg` to maintain scalability.
- **Polymorphic Skill Profiles:** Use a `SkillProfile` struct with optional pointers to domain-specific containers (`DeviceContainer`, `NetworkContainer`, `AppContainer`).
  - *Rationale:* This matches the 3GPP proposal's data model and allows different NFs/UEs to register diverse capabilities under a common header.
- **In-Memory Thread-Safe Registry:** Use `sync.Map` or a protected `map` with `RWMutex` for the ACRF registry.
  - *Rationale:* Simple to implement while ensuring correctness during parallel registration/discovery.
- **HTTP as Inter-Service Communication:** Use standard `net/http` for simplicity in the prototype.
  - *Rationale:* Alignment with 3GPP SBA principles and easy to test with `curl`.

## Risks / Trade-offs

- **[Risk] State Loss:** ACRF state is lost on restart.
  - *Mitigation:* Ensure IGW re-registers on startup.
- **[Trade-off] String-Based Discovery:** The MVP uses exact `SkillID` matching.
  - *Mitigation:* The registry logic is encapsulated behind an interface to allow for future semantic matching.
