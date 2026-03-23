## 1. Setup

- [x] 1.1 Create `tests/` directory
- [x] 1.2 Implement `internal/testutil` package for service orchestration

## 2. Mock Agent Implementation

- [x] 2.1 Refactor `internal/agent/agent.go` to support injecting a `MockCoreAgent`
- [x] 2.2 Implement the `MockCoreAgent` logic for deterministic testing

## 3. Integration Test Suite

- [x] 3.1 Implement `tests/system_integration_test.go`
- [x] 3.2 Add test cases for "Fleet Wake-Up" end-to-end flow
- [x] 3.3 Verify that tests pass using `go test ./tests/...`

## 4. Documentation

- [x] 4.1 Update `GEMINI.md` with instructions on running tests
- [x] 4.2 Document the "Test-First" implementation policy
