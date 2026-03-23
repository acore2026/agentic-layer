package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/6g-agentic-core/internal/agent"
	"github.com/google/6g-agentic-core/internal/config"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

const AppName = "6G-Agentic-Core"

func main() {
	// 1. Load configuration (reads .env automatically)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	coreAgent, err := agent.NewCoreAgent(ctx)
	if err != nil {
		log.Fatalf("Failed to create core agent: %v", err)
	}

	sessionService := session.InMemoryService()
	r, err := runner.New(runner.Config{
		AppName:        AppName,
		Agent:          coreAgent,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	handler := agent.NewHandler(r, sessionService, AppName)

	log.Printf("AAIHF (Agentic AI Host Function) starting on :%s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
