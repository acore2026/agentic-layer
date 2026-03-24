package temporal_skills

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_FleetWakeUpWorkflow_Success() {
	s.env.OnActivity(CallAMFActivity, mock.Anything, mock.Anything).Return("AMF Success", nil)
	s.env.OnActivity(CallSMFActivity, mock.Anything, mock.Anything).Return("SMF Success", nil)
	s.env.OnActivity(CallNEFActivity, mock.Anything, mock.Anything).Return("NEF Success", nil)

	s.env.ExecuteWorkflow(FleetWakeUpWorkflow, FleetUpdateInput{Action: "test"})

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	
	var result string
	s.env.GetWorkflowResult(&result)
	s.Equal("Workflow completed successfully", result)
}

func (s *UnitTestSuite) Test_FleetWakeUpWorkflow_SMFFailureRollsbackAMF() {
	// AMF succeeds
	s.env.OnActivity(CallAMFActivity, mock.Anything, mock.Anything).Return("AMF Success", nil)
	
	// SMF fails
	smfError := errors.New("SMF totally failed")
	s.env.OnActivity(CallSMFActivity, mock.Anything, mock.Anything).Return("", smfError)
	
	// Because SMF failed, it should trigger RollbackAMFActivity
	s.env.OnActivity(RollbackAMFActivity, mock.Anything, mock.Anything).Return("Rollback Success", nil)

	s.env.ExecuteWorkflow(FleetWakeUpWorkflow, FleetUpdateInput{Action: "test"})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "SMF totally failed")
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
