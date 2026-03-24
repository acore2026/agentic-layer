## Context

The AAIHF (Agentic AI Host Function) is the core reasoning component of the 6G Agentic Layer. It currently processes user intents in a black-box fashion from the perspective of external clients. To enable real-time observability, we need a mechanism to stream internal events (reasoning steps, tool calls, execution results) to a frontend dashboard. Server-Sent Events (SSE) is chosen for its simplicity and native support in web browsers.

## Goals / Non-Goals

**Goals:**
- Provide a persistent streaming endpoint (`/stream`) on the AAIHF.
- Broadcast key lifecycle events of the intent resolution pipeline.
- Support multiple concurrent listeners (pub/sub pattern).
- Ensure CORS compliance for cross-domain dashboard integration.

**Non-Goals:**
- Real-time command/control via the stream (one-way only).
- Historical event playback (stream starts from current time).
- Per-session private streams (events are broadcast to all subscribers for this MVP).

## Decisions

- **Event Broker Pattern**: A centralized `Broker` in `internal/events` will manage client registration and message broadcasting. This keeps the HTTP and Agent logic decoupled.
- **Dependency Injection**: The `Broker` will be injected into the `Agent` or `Server` constructors. This allows for easier testing and prevents circular dependencies.
- **JSON Event Format**: All events will be sent as JSON objects with `type` and `data` fields for consistent parsing by the frontend.
- **Non-blocking Emissions**: Event emissions from the Agent logic will be non-blocking (using a buffered channel or select with default) to ensure that slow or disconnected SSE clients do not impact the core reasoning performance.

## Risks / Trade-offs

- **[Risk] High Memory Usage with many clients**: Each SSE client holds a connection and a buffer.
  - *Mitigation*: We will use a standard Go channel-based broker which is efficient. For an MVP with limited dashboard users, this is negligible.
- **[Trade-off] Global Broadcast**: All events are seen by all connected clients.
  - *Rationale*: For a prototype dashboard, global visibility is desired. Future versions could filter by `user_id`.
