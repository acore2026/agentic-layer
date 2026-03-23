package translator

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/6g-agentic-core/pkg/models"
)

func RegisterSkillWithACRF(acrfURL string, skillID string, invokeURL string) {
	profile := models.SkillProfile{
		SkillID:           skillID,
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
		log.Printf("ACRF not ready (attempt %d/5), retrying in 2s...", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Printf("Failed to register skill %s with ACRF after 5 attempts", skillID)
}

func NewHandler(t Translator) http.Handler {
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

		if err := t.Translate(req.SkillID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Skill successfully translated and executed\n"))
	})

	return mux
}
