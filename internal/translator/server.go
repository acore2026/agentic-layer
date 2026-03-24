package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/6g-agentic-core/internal/translator/temporal_skills"
	"github.com/google/6g-agentic-core/pkg/models"
	"go.temporal.io/sdk/client"
)

func RegisterSkillWithACRF(acrfURL string, skillID string, invokeURL string) {
	profile := models.SkillProfile{
		SkillID:           skillID,
		Description:       "Trigger a wake-up and firmware update sequence for a device fleet. Ensures reachability and sets up SM context.",
		EntityType:        "NF",
		ServiceClass:      models.ServiceClassSilver,
		AgenticServiceURI: invokeURL,
		Network: &models.NetworkContainer{
			NetworkLocality: "Transport-Edge-West",
			ServiceArea:     []string{"TAI-Range-50-to-60"},
		},
	}

	jsonData, err := json.Marshal(profile)
	if err != nil {
		log.Fatalf("Failed to marshal skill profile: %v", err)
	}

	for i := 0; i < 5; i++ {
		resp, err := http.Post(acrfURL, "application/json", bytes.NewBuffer(jsonData))
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Printf("Successfully registered skill %s with ACRF", skillID)
			return
		}
		log.Printf("ACRF not ready or error (attempt %d/5), retrying in 2s...", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Printf("Failed to register skill %s with ACRF after 5 attempts", skillID)
}

type WorkflowStarter interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

func NewHandler(tc WorkflowStarter) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/invoke", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			SkillID string `json:"skill_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Prepare Workflow Execution
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("skill-exec-%s", time.Now().Format("20060102150405")),
			TaskQueue: "FleetManagementTaskQueue",
		}
		
		// Map SkillID to appropriate action name for logs
		action := "NetworkSkillExecution"
		if req.SkillID == "mcp://skill/device/fleet-update" {
			action = "FirmwareUpdateWakeup"
		}

		input := temporal_skills.FleetUpdateInput{
			SkillID:   req.SkillID,
			TargetUEs: []string{"default-target"}, 
			Action:    action,
		}

		if tc == nil {
			http.Error(w, "Temporal client is not initialized (check IGW logs)", http.StatusInternalServerError)
			return
		}

		we, err := tc.ExecuteWorkflow(context.Background(), workflowOptions, temporal_skills.FleetWakeUpWorkflow, input)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to start Temporal workflow: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("Workflow execution started asynchronously. WorkflowID: %s, RunID: %s\n", we.GetID(), we.GetRunID())))
	})

	return mux
}
