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
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

const TestAppName = "6G-Agentic-Core-Test"

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
	handler := translator.NewHandler(trans)
	url, closer, err := StartService(handler)
	if err != nil {
		return "", nil, err
	}

	// Register in background
	go translator.RegisterSkillWithACRF(acrfURL+"/register", "mcp://skill/device/fleet-update", url+"/invoke")
	// Wait a bit for registration
	time.Sleep(100 * time.Millisecond)

	return url, closer, nil
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
