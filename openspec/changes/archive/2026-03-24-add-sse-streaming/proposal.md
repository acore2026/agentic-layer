## Why

Currently, the AAIHF (Agentic AI Host Function) processes intents synchronously, and there is no way for external observers (like a frontend dashboard) to see the intermediate reasoning steps or real-time progress. Adding Server-Sent Events (SSE) streaming will enable real-time visualization of the AI's thought process, tool calls, and execution status, making the system more transparent and interactive.

## What Changes

- **SSE Broker Implementation**: Create a centralized event broker in `internal/events` to manage SSE client connections and broadcast messages.
- **Agent Instrumentation**: Update the `Agent` and `Server` in `internal/agent` to emit events at key stages of the intent resolution lifecycle (e.g., reasoning started, skill discovered, execution completed).
- **API Endpoint**: Expose a new `/stream` endpoint on the AAIHF server to allow clients to subscribe to the event stream.
- **CORS Support**: Ensure the SSE stream supports Cross-Origin Resource Sharing (CORS) for frontend integration.

## Capabilities

### New Capabilities
- `real-time-observability`: Provides a streaming event interface (SSE) for observing system-wide events and AI reasoning steps.

### Modified Capabilities
- `intent-resolution`: Instrumented to emit events during the reasoning and execution pipeline.

## Impact

- **AAIHF API**: Adds a new GET `/stream` endpoint.
- **Internal Architecture**: Introduces a dependency on the new `events` package within the `agent` package.
- **Dependencies**: Uses standard Go libraries; no new external dependencies required for the broker itself.
