## Context

The A-IGW acts as a translator between abstract 6G skill intents (like "Wake up the fleet") and deterministic 3GPP OpenAPI operations. Currently, this execution is handled via synchronous blocking logic inside the `/invoke` HTTP handler. If a complex sequence (e.g., AMF -> SMF -> NEF) fails midway, the system lacks state tracking to retry or perform compensations reliably. Integrating Temporal.io into A-IGW provides a stateful, fault-tolerant execution environment.

## Goals / Non-Goals

**Goals:**
- Add `go.temporal.io/sdk` dependency.
- Define Temporal Activities (`CallAMFActivity`, `CallSMFActivity`, `CallNEFActivity`, `RollbackAMFActivity`) and a Temporal Workflow (`FleetWakeUpWorkflow`).
- Start a Temporal Worker process inside the A-IGW service.
- Refactor the A-IGW `/invoke` HTTP handler to execute the workflow asynchronously and return `202 Accepted`.

**Non-Goals:**
- Complete deployment of Temporal Cluster (we assume a local or existing Temporal server is accessible or will be documented).
- Rewriting all other mock services (like AAIHF or ACRF) for Temporal—this change is scoped specifically to A-IGW execution.

## Decisions

- **Temporal Worker Co-location:** The Temporal Worker will be co-located with the A-IGW HTTP server within the `cmd/igw-fleet/main.go` entry point. 
  - *Rationale*: Reduces operational complexity for the MVP by keeping A-IGW as a single deployable binary.
- **Workflow and Activity Separation:** A new package `internal/translator/temporal_skills` will be created to house `workflows.go` and `activities.go`.
  - *Rationale*: Cleanly isolates the Temporal orchestration logic from the HTTP transport layer.
- **Retry and Rollback Strategy:** The SMF activity will utilize a standard Temporal RetryPolicy (e.g., 3 attempts, 2s intervals). If it ultimately fails, the Workflow will explicitly invoke `RollbackAMFActivity` for compensation.

## Risks / Trade-offs

- **[Risk] Added Infrastructure Dependency:** Temporal requires its own cluster (Temporal Server + DB + visibility store). 
  - *Mitigation*: For MVP testing, the Temporal CLI or a local docker-compose setup can be used.
- **[Risk] Workflow Blocking on Missing Temporal Server:** If the Temporal server is down, the IGW HTTP handler will fail to start the workflow.
  - *Mitigation*: The `ExecuteWorkflow` call will return an error which the HTTP handler will translate into an appropriate HTTP 500 error.
