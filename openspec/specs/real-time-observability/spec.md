## ADDED Requirements

### Requirement: SSE Streaming Endpoint
The AAIHF SHALL provide a GET `/stream` endpoint that uses Server-Sent Events (SSE) to broadcast system events to connected clients.

#### Scenario: Successful Subscription
- **WHEN** a client makes a GET request to `/stream`
- **THEN** the server SHALL respond with `Content-Type: text/event-stream` and keep the connection open.

### Requirement: Event Broadcast
The system SHALL broadcast events in JSON format containing a `type` and a `data` payload.

#### Scenario: Event Emission
- **WHEN** an internal component (e.g., Agent) emits an event
- **THEN** all clients connected to `/stream` SHALL receive the data in the format `data: {"type": "...", "data": ...}\n\n`.

### Requirement: CORS Compliance
The streaming endpoint SHALL allow requests from any origin to support browser-based dashboards.

#### Scenario: Cross-Origin Request
- **WHEN** a client from a different origin requests `/stream`
- **THEN** the server SHALL include the `Access-Control-Allow-Origin: *` header in the response.
