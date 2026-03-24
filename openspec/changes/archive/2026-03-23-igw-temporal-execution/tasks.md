## 1. Setup

- [x] 1.1 Run `go get go.temporal.io/sdk` to add Temporal dependencies.

## 2. Temporal Skills Package

- [x] 2.1 Create directory `internal/translator/temporal_skills`.
- [x] 2.2 Create `activities.go` with mock functions: `CallAMFActivity`, `CallSMFActivity`, `CallNEFActivity`, and `RollbackAMFActivity`.
- [x] 2.3 Create `workflows.go` implementing `FleetWakeUpWorkflow` with AMF -> SMF -> NEF execution, RetryPolicy, and rollback logic.

## 3. Server Integration

- [x] 3.1 Update `cmd/igw-fleet/main.go` to initialize a Temporal Client and start a Temporal Worker.
- [x] 3.2 Register the workflow and activities with the Temporal Worker in `main.go`.
- [x] 3.3 Refactor the `/invoke` HTTP handler in `internal/translator/server.go` to trigger `FleetWakeUpWorkflow` asynchronously and return `202 Accepted` with the WorkflowID.

## 4. Cleanup
- [x] 4.1 Remove files under `tmp-prompts` directory.