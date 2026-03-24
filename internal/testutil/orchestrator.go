package testutil

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	coreagent "github.com/google/6g-agentic-core/internal/agent"
	"github.com/google/6g-agentic-core/internal/registry"
	"github.com/google/6g-agentic-core/internal/translator"
	"github.com/google/6g-agentic-core/internal/translator/temporal_skills"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"go.temporal.io/sdk/client"
)

const TestAppName = "6G-Agentic-Core-Test"

// mockWorkflowRun implements client.WorkflowRun for testing
type mockWorkflowRun struct{}

func (m mockWorkflowRun) GetID() string { return "mock-workflow-id" }
func (m mockWorkflowRun) GetRunID() string { return "mock-run-id" }
func (m mockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error { return nil }
func (m mockWorkflowRun) GetWithOptions(ctx context.Context, valuePtr interface{}, options client.WorkflowRunGetOptions) error { return nil }

// mockWorkflowStarter implements translator.WorkflowStarter for testing
type mockWorkflowStarter struct {
	t *translator.FleetTranslator
}

func (m mockWorkflowStarter) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	fmt.Println("[MockWorkflowStarter] Executing workflow via real translator handlers...")
	
	// Extract SkillID from FleetUpdateInput (first argument)
	if len(args) > 0 {
		if input, ok := args[0].(temporal_skills.FleetUpdateInput); ok {
			// Actually execute the handler mapping
			m.t.Translate(input.SkillID)
		} else {
			fmt.Printf("[MockWorkflowStarter] Warning: unexpected argument type: %T\n", args[0])
			m.t.Translate("mcp://skill/device/fleet-update")
		}
	} else {
		m.t.Translate("mcp://skill/device/fleet-update")
	}
	
	return mockWorkflowRun{}, nil
}

// StartService starts a given handler on a random port and returns the URL and a closer function.
func StartService(handler http.Handler) (string, func(), error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	server := &http.Server{Handler: handler}
	go func() {
		server.Serve(listener)
	}()

	closer := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}

	return url, closer, nil
}

func SetupACRF() (string, func(), error) {
	reg := registry.NewInMemoryRegistry()
	handler := registry.NewHandler(reg)
	return StartService(handler)
}

func SetupIGW(acrfURL string) (string, func(), error) {
	trans := translator.NewFleetTranslator()
	handler := translator.NewHandler(mockWorkflowStarter{t: trans})
	return StartService(handler)
}

func SetupAAIHF(coreAgent agent.Agent) (string, func(), error) {
	sessionService := session.InMemoryService()
	r, err := runner.New(runner.Config{
		AppName:        TestAppName,
		Agent:          coreAgent,
		SessionService: sessionService,
	})
	if err != nil {
		return "", nil, err
	}

	handler := coreagent.NewHandler(r, sessionService, TestAppName)
	return StartService(handler)
}
