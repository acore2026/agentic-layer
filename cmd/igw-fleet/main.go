package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/6g-agentic-core/internal/translator"
	"github.com/google/6g-agentic-core/internal/translator/temporal_skills"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	temporalHost := os.Getenv("AGENTIC_TEMPORAL_HOST")
	if temporalHost == "" {
		temporalHost = client.DefaultHostPort
	}

	log.Printf("Connecting to Temporal at %s...", temporalHost)

	var c client.Client
	var err error
	
	// Retry connection to Temporal
	for i := 0; i < 5; i++ {
		c, err = client.Dial(client.Options{
			HostPort: temporalHost,
		})
		if err == nil {
			break
		}
		log.Printf("Warning: Temporal not ready (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		log.Printf("Warning: Failed to connect to Temporal after 5 attempts: %v", err)
	} else {
		defer c.Close()

		// Create and start Temporal Worker
		w := worker.New(c, "FleetManagementTaskQueue", worker.Options{})
		
		// Register Workflows and Activities
		w.RegisterWorkflow(temporal_skills.FleetWakeUpWorkflow)
		w.RegisterActivity(temporal_skills.CallAMFActivity)
		w.RegisterActivity(temporal_skills.CallSMFActivity)
		w.RegisterActivity(temporal_skills.CallNEFActivity)
		w.RegisterActivity(temporal_skills.RollbackAMFActivity)

		// Start worker in a background goroutine
		go func() {
			err = w.Run(worker.InterruptCh())
			if err != nil {
				log.Println("Warning: Temporal worker stopped", err)
			}
		}()
	}

	// The HTTP Handler now uses the Temporal client to start workflows
	handler := translator.NewHandler(c)

	log.Println("A-IGW (Interworking Gateway) starting on :18081...")
	if err := http.ListenAndServe(":18081", handler); err != nil {
		log.Fatal(err)
	}
}
