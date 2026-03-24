# Role
You are a Senior Go Developer and Telecommunications Architect. You are helping me build a prototype for a "6G Skill-Based Agentic Core Network" based on a recent 3GPP architectural proposal.

# Project Context
I already have a working 5G foundation running on a cloud server using `free5gc` (for the core NFs) and `ueransim-go` (for the RAN/UE). 
We are building a new, decoupled AI layer that sits on top of this 5G core to enable intent-driven networking. We are strictly separating the new AI logic from the legacy 3GPP deterministic state machines.

We are creating a new Go repository called `6g-agentic-core` which will contain three new microservices:
1. **AAIHF (Agentic AI Host Function):** The brain. It receives natural language intents and uses the `github.com/google/adk-go` framework to process them and discover skills.
2. **ACRF (Agentic Capability Repository Function):** A dynamic skill registry. It stores "Agentic Skill URIs" and maps them to HTTP execution endpoints. For this MVP, it will just be an in-memory Go map exposed via REST.
3. **Interworking Gateway (A-IGW):** The "Universal Translator". It registers a high-level skill with the ACRF, and when invoked, it translates that skill into a deterministic sequence of 3GPP OpenAPI calls to the `free5gc` NFs.

# The Target Use Case: "Fleet Wake-Up"
We are implementing one specific skill flow to prove the architecture:
- **Intent:** "Wake up fleet for firmware update."
- [cite_start]**Skill URI:** `mcp://skill/device/fleet-update` [cite: 106]
- [cite_start]**Service Directives (The Translation):** When the IGW receives this skill invocation, it must trigger 3 downstream API stubs[cite: 106]:
  1. [cite_start]`Namf_MT_EnableUEReachability` (to AMF) [cite: 106]
  2. [cite_start]`Nsmf_PDUSession_UpdateSMContext` (to SMF) [cite: 106]
  3. [cite_start]`Nnef_AFSessionWithQoS_Create` (to NEF) [cite: 106]

# Directory Structure Requirement
Please structure the Go code using standard layouts:
```text
agentic-layer/
├── cmd/
│   ├── aaihf/           # main.go
│   ├── acrf/            # main.go
│   └── igw-fleet/       # main.go
├── internal/
│   ├── agent/           # adk-go logic for AAIHF
│   ├── registry/        # In-memory map logic for ACRF
│   └── translator/      # Skill-to-API mapping logic for IGW
├── pkg/
│   └── models/          # Shared structs (e.g., SkillProfile)
├── go.mod
└── go.sum
```

# Your Task (Step 1)
I want to build this iteratively. Please do the following:
1. Provide the bash commands to initialize the `6g-agentic-core` Go module and create the directory structure above.
2. Write the Go code for `pkg/models/models.go` to define a `SkillProfile` struct containing `SkillID` (string), `EntityType` (string), and `AgenticServiceURI` (string).
3. Write the Go code for the `ACRF` (`cmd/acrf/main.go` and `internal/registry/registry.go`). It needs a simple HTTP server with a POST `/register` endpoint to save a `SkillProfile` to memory, and a GET `/discover?skill_id=...` endpoint to retrieve it.
4. Write the Go code for the `IGW` (`cmd/igw-fleet/main.go` and `internal/translator/translator.go`). On startup, it must POST its `mcp://skill/device/fleet-update` profile to the ACRF. It must expose an HTTP POST `/invoke` endpoint. When hit, it should print log messages simulating the 3 downstream AMF, SMF, and NEF API calls.

Write clean, idiomatic Go code. Do not implement the `AAIHF` (adk-go part) yet; let's get the Registry and Gateway running first.