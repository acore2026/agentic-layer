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

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/discover", discoverHandler)

	log.Println("ACRF (Agentic Capability Repository Function) starting on :18080...")
	if err := http.ListenAndServe(":18080", nil); err != nil {
		log.Fatal(err)
	}
}
