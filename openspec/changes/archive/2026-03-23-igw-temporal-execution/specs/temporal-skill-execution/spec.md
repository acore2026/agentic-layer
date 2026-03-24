## ADDED Requirements

### Requirement: Temporal Workflow Orchestration
The A-IGW SHALL utilize Temporal.io to orchestrate complex network skill executions as stateful workflows, ensuring resilience against partial failures.

#### Scenario: Workflow Execution
- **WHEN** a skill is triggered
- **THEN** a Temporal workflow SHALL coordinate the execution of the required sequence of activities.

### Requirement: Activity Retries and Rollbacks
Network activities orchestrated by Temporal SHALL support automated retries for transient failures and explicit rollbacks (compensation actions) for terminal failures in subsequent steps.

#### Scenario: SMF Failure Rollback
- **WHEN** the SMF activity fails permanently within the Fleet Wake-Up workflow
- **THEN** the workflow SHALL explicitly execute a rollback activity for the preceding AMF step.
