## 1. Setup and Dependencies

- [x] 1.1 Add `github.com/google/adk-go` and dependencies to `go.mod`
- [x] 1.2 Initialize the `internal/agent` package

## 2. Implement Agent Tools

- [x] 2.1 Implement `SearchSkill` tool (ACRF discovery) in `internal/agent/tools.go`
- [x] 2.2 Implement `ExecuteSkill` tool (A-IGW invocation) in `internal/agent/tools.go`

## 3. Implement AAIHF Reasoner

- [x] 3.1 Define the `GeneralCoreAgent` using `adk-go` in `internal/agent/agent.go`
- [x] 3.2 Provide system instructions to the agent explaining the 6G skill architecture

## 4. Implement AAIHF Service

- [x] 4.1 Create HTTP server with `/intent` endpoint in `cmd/aaihf/main.go`
- [x] 4.2 Integrate the agent logic into the intent handler
- [x] 4.3 Verify the end-to-end "Fleet Wake-Up" flow via natural language prompt
