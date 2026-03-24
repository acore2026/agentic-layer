## Why

Currently, the A-IGW executes deterministic 3GPP signaling sequences synchronously. This approach is brittle and lacks retries, state management, and reliable rollbacks in case of partial failures (e.g., SMF update fails after AMF wake-up). Integrating Temporal.io provides a robust, stateful workflow engine to execute these complex network capabilities reliably.

## What Changes

- **Add Temporal.io SDK Dependency**: Introduce `go.temporal.io/sdk` to the `A-IGW`.
- **Implement Temporal Skills**: Create a new `temporal_skills` package containing the `FleetWakeUpWorkflow` and associated activities (`CallAMFActivity`, `CallSMFActivity`, `CallNEFActivity`, `RollbackAMFActivity`).
- **Update A-IGW Server**: Initialize a Temporal Worker and Client within the `igw-fleet` entry point.
- **Refactor /invoke Handler**: Change the `A-IGW` invocation handler to trigger the Temporal workflow asynchronously and return a `202 Accepted` with the workflow ID instead of blocking.

## Capabilities

### New Capabilities
- `temporal-skill-execution`: Execution of network capabilities using Temporal.io for deterministic, reliable state management and rollbacks.

### Modified Capabilities
- `fleet-wake-up-translation`: Update the execution strategy to be asynchronous and handle rollbacks on failure.

## Impact

- **Architecture**: Introduces Temporal.io as a required infrastructure component (for production) and SDK dependency for the A-IGW.
- **API**: Changes the `/invoke` API behavior from synchronous blocking execution to an asynchronous pattern returning `202 Accepted`.
- **Reliability**: Greatly improves system resilience during complex core network signaling.
