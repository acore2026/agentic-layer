package registry

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/6g-agentic-core/pkg/models"
)

func NewHandler(reg Registry) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var profile models.SkillProfile
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		reg.Register(profile)
		log.Printf("Registered skill: %s (Entity: %s)", profile.SkillID, profile.EntityType)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/discover", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		skillID := r.URL.Query().Get("skill_id")
		if skillID == "" {
			http.Error(w, "Missing skill_id parameter", http.StatusBadRequest)
			return
		}

		profile, ok := reg.Discover(skillID)
		if !ok {
			http.Error(w, "Skill not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	})

	return mux
}
