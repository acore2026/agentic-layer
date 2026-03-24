package temporal_skills

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func FleetWakeUpWorkflow(ctx workflow.Context, input FleetUpdateInput) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("FleetWakeUpWorkflow started", "Action", input.Action)

	var amfResult string
	err := workflow.ExecuteActivity(ctx, CallAMFActivity, input).Get(ctx, &amfResult)
	if err != nil {
		logger.Error("CallAMFActivity failed", "Error", err)
		return "", err
	}

	var smfResult string
	err = workflow.ExecuteActivity(ctx, CallSMFActivity, input).Get(ctx, &smfResult)
	if err != nil {
		logger.Error("CallSMFActivity failed, initiating rollback", "Error", err)
		
		// Rollback AMF
		var rollbackResult string
		rollbackErr := workflow.ExecuteActivity(ctx, RollbackAMFActivity, input).Get(ctx, &rollbackResult)
		if rollbackErr != nil {
			logger.Error("RollbackAMFActivity failed", "Error", rollbackErr)
			return "", err // Return original error, but note rollback failed in logs
		}
		
		return "", err
	}

	var nefResult string
	err = workflow.ExecuteActivity(ctx, CallNEFActivity, input).Get(ctx, &nefResult)
	if err != nil {
		logger.Error("CallNEFActivity failed", "Error", err)
		// We could add SMF/AMF rollback here too, but following prompt instructions for now
		return "", err
	}

	logger.Info("FleetWakeUpWorkflow completed successfully")
	return "Workflow completed successfully", nil
}
