## 1. Project Initialization

- [x] 1.1 Initialize Go module `github.com/google/6g-agentic-core`
- [x] 1.2 Create directory structure (`cmd/`, `internal/`, `pkg/`)

## 2. Core Models

- [x] 2.1 Define `ServiceClass` enumeration in `pkg/models/models.go`
- [x] 2.2 Implement polymorphic `SkillProfile` and domain containers in `pkg/models/models.go`

## 3. ACRF (Skill Registry) Implementation

- [x] 3.1 Implement thread-safe in-memory registry in `internal/registry/registry.go`
- [x] 3.2 Create HTTP server with `/register` and `/discover` endpoints in `cmd/acrf/main.go`
- [x] 3.3 Verify registry functionality with basic unit tests or manual `curl` tests

## 4. A-IGW (Interworking Gateway) Implementation

- [x] 4.1 Implement skill-to-API translation logic for "Fleet Wake-Up" in `internal/translator/translator.go`
- [x] 4.2 Create HTTP server with `/invoke` endpoint in `cmd/igw-fleet/main.go`
- [x] 4.3 Implement boot-time registration logic to POST `fleet-update` profile to ACRF
- [x] 4.4 Verify translation logs upon skill invocation
