package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/6g-agentic-core/internal/translator"
	"github.com/google/6g-agentic-core/pkg/models"
)

var fleetTranslator = translator.NewFleetTranslator()

const (
	acrfURL = "http://localhost:8080/register"
	skillID = "mcp://skill/device/fleet-update"
)

func registerSkill() {
	profile := models.SkillProfile{
		SkillID:           skillID,
		EntityType:        "NF", // Interworking Gateway is an NF
		ServiceClass:      models.ServiceClassSilver,
		AgenticServiceURI: "http://localhost:8081/invoke",
		Network: &models.NetworkContainer{
			NetworkLocality: "Transport-Edge-West",
			ServiceArea:     []string{"TAI-Range-50-to-60"},
		},
	}

	jsonData, err := json.Marshal(profile)
	if err != nil {
		log.Fatalf("Failed to marshal skill profile: %v", err)
	}

	// Retry logic for ACRF registration (wait for ACRF to be up)
	for i := 0; i < 5; i++ {
		resp, err := http.Post(acrfURL, "application/json", bytes.NewBuffer(jsonData))
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Printf("Successfully registered skill %s with ACRF", skillID)
			return
		}
		log.Printf("ACRF not ready (attempt %d/5), retrying in 2s...", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("Failed to register skill %s with ACRF after 5 attempts", skillID)
}

func invokeHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := fleetTranslator.Translate(req.SkillID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Skill successfully translated and executed\n"))
}

func main() {
	// Register skill with ACRF in background
	go registerSkill()

	http.HandleFunc("/invoke", invokeHandler)

	log.Println("A-IGW (Interworking Gateway) starting on :8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
