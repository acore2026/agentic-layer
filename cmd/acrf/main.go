package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/6g-agentic-core/internal/registry"
	"github.com/google/6g-agentic-core/pkg/models"
	"github.com/joho/godotenv"
)

var reg = registry.NewInMemoryRegistry()

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var profile models.SkillProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := reg.Register(profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func discoverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("skill_id")
	if query == "" {
		http.Error(w, "Missing skill_id (query) parameter", http.StatusBadRequest)
		return
	}

	profile, ok := reg.Discover(query)
	if !ok {
		http.Error(w, "No semantically similar skill found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func bootstrapSkills() {
	skills := []models.SkillProfile{
		{
			SkillID:           "mcp://skill/device/fleet-update",
			Description:       "Trigger a wake-up and firmware update sequence for a device fleet. Ensures reachability and sets up SM context.",
			EntityType:        "NF",
			AgenticServiceURI: "http://localhost:18081/invoke",
		},
		{
			SkillID:           "mcp://skill/qos/turbo-mode",
			Description:       "Enable Turbo Mode for Gaming Session.",
			EntityType:        "NF",
			AgenticServiceURI: "http://localhost:18081/invoke",
		},
		{
			SkillID:           "mcp://skill/reliability/path-diversity",
			Description:       "Ensure Zero-Interruption for V2X Feed",
			EntityType:        "NF",
			AgenticServiceURI: "http://localhost:18081/invoke",
		},
		{
			SkillID:           "mcp://skill/edge/secure-flight",
			Description:       "Secure Drone Corridor.",
			EntityType:        "NF",
			AgenticServiceURI: "http://localhost:18081/invoke",
		},
	}

	for _, skill := range skills {
		if err := reg.Register(skill); err != nil {
			log.Printf("Warning: failed to bootstrap skill %s: %v", skill.SkillID, err)
		}
	}
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Bootstrap default skills
	bootstrapSkills()

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/discover", discoverHandler)

	log.Println("ACRF (Agentic Capability Repository Function) starting on :18080...")
	if err := http.ListenAndServe(":18080", nil); err != nil {
		log.Fatal(err)
	}
}
