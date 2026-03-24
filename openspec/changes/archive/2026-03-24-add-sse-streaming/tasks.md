## 1. SSE Broker Implementation

- [x] 1.1 Create `internal/events/sse.go` with `Broker` struct and methods.
- [x] 1.2 Implement client registration and non-blocking broadcast logic in `Broker`.
- [x] 1.3 Implement `ServeHTTP` with mandatory SSE and CORS headers.

## 2. Agent Instrumentation

- [x] 2.1 Update `internal/agent/agent.go` to accept `*events.Broker` in constructors.
- [x] 2.2 Implement `emitEvent` helper method in `internal/agent/agent.go`.
- [x] 2.3 Add `reasoning_started` event emission at the start of intent processing.
- [x] 2.4 Add `reasoning_completed` event emission with final response.

## 3. Server Integration

- [x] 3.1 Update `internal/agent/server.go` to accept and store the `Broker`.
- [x] 3.2 Update `cmd/aaihf/main.go` to instantiate the `Broker` and mount it on `/stream`.
- [x] 3.3 Pass the `Broker` into the Agent/Server initialization in `main.go`.

## 4. Verification

- [x] 4.1 Verify GET `/stream` returns correct headers and stays open.
- [x] 4.2 Verify real-time event delivery during an intent resolution using `curl`.
