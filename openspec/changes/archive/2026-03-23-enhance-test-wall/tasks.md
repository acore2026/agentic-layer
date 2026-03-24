## 1. Setup

- [x] 1.1 Run `go get github.com/cucumber/godog/cmd/godog@latest` to add BDD dependency.

## 2. Unit Testing the Agent

- [x] 2.1 Add `internal/agent/tools_test.go` with table-driven unit tests for `SearchSkill` and `ExecuteSkill` using `httptest.NewServer`.

## 3. Temporal Workflow Test Suite

- [x] 3.1 Create `internal/translator/temporal_skills/workflows_test.go`.
- [x] 3.2 Implement a successful isolated workflow test using `testsuite.WorkflowTestSuite`.
- [x] 3.3 Implement an error-path test verifying `RollbackAMFActivity` is called when `CallSMFActivity` fails.

## 4. BDD Testing with godog

- [x] 4.1 Create `tests/features/system.feature` with a BDD scenario for the full system integration.
- [x] 4.2 Create `tests/godog_test.go` to parse the `.feature` file and execute steps.

## 5. Expanded Go Integration Tests

- [x] 5.1 Add an edge case test in `tests/system_integration_test.go` to handle LLM intent hallucinations gracefully.
- [x] 5.2 Add a failure path test verifying the 500 error propagation if the A-IGW temporal client is down.