package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/6g-agentic-core/internal/events"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func NewHandler(r *runner.Runner, sessionService session.Service, appName string, broker *events.Broker) http.Handler {
	mux := http.NewServeMux()

	// Mount the SSE stream if broker is provided
	if broker != nil {
		mux.Handle("/stream", broker)
	}

	mux.HandleFunc("/intent", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var intentReq struct {
			Prompt string `json:"prompt"`
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(req.Body).Decode(&intentReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if intentReq.UserID == "" {
			intentReq.UserID = "default-user"
		}
		sessionID := "current-session"

		ctx := req.Context()

		// Ensure session exists
		_, err := sessionService.Get(ctx, &session.GetRequest{
			AppName:   appName,
			UserID:    intentReq.UserID,
			SessionID: sessionID,
		})
		if err != nil {
			log.Printf("Creating new session for user %s", intentReq.UserID)
			_, err = sessionService.Create(ctx, &session.CreateRequest{
				AppName:   appName,
				UserID:    intentReq.UserID,
				SessionID: sessionID,
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
				return
			}
		}

		msg := &genai.Content{
			Parts: []*genai.Part{{Text: intentReq.Prompt}},
			Role:  "user",
		}

		log.Printf("Received intent from user %s: %s", intentReq.UserID, intentReq.Prompt)

		// Run the agent
		events := r.Run(ctx, intentReq.UserID, sessionID, msg, agent.RunConfig{})

		var finalResponse string
		for event, err := range events {
			if err != nil {
				log.Printf("Error during agent run: %v", err)
				http.Error(w, fmt.Sprintf("Reasoning error: %v", err), http.StatusInternalServerError)
				return
			}

			// Capture the model's text response
			if event.Content != nil {
				for _, part := range event.Content.Parts {
					if part.Text != "" {
						finalResponse = part.Text
					}
				}
			}

			if event.IsFinalResponse() {
				log.Printf("Agent reached final response: %s", finalResponse)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"response": finalResponse,
		})
	})

	return mux
}
