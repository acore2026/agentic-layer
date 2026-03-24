package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSearchSkill(t *testing.T) {
	// Mock ACRF Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/discover" {
			t.Errorf("Expected to request '/discover', got: %s", r.URL.Path)
		}
		
		skillID := r.URL.Query().Get("skill_id")
		
		if skillID == "valid-skill" || skillID == "mcp://skill/test" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"skill_id":"` + skillID + `","description":"test"}`))
			return
		}
		
		if skillID == "not-found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	os.Setenv("AGENTIC_ACRF_URL", server.URL)
	defer os.Unsetenv("AGENTIC_ACRF_URL")

	tests := []struct {
		name        string
		input       SearchSkillInput
		wantOutput  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Successful Discovery",
			input:      SearchSkillInput{SkillID: "valid-skill"},
			wantOutput: `{"skill_id":"valid-skill","description":"test"}`,
			wantErr:    false,
		},
		{
			name:       "URL Encoded Skill ID",
			input:      SearchSkillInput{SkillID: "mcp://skill/test"},
			wantOutput: `{"skill_id":"mcp://skill/test","description":"test"}`,
			wantErr:    false,
		},
		{
			name:       "Not Found",
			input:      SearchSkillInput{SkillID: "not-found"},
			wantOutput: "Skill not-found not found in ACRF registry.",
			wantErr:    false,
		},
		{
			name:        "Server Error",
			input:       SearchSkillInput{SkillID: "error-skill"},
			wantErr:     true,
			errContains: "ACRF returned error status: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchSkill(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchSkill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("SearchSkill() error = %v, expected it to contain %v", err, tt.errContains)
			}
			if !tt.wantErr && got != tt.wantOutput {
				t.Errorf("SearchSkill() got = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}

func TestExecuteSkill(t *testing.T) {
	// Mock IGW Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/invoke" {
			t.Errorf("Expected to request '/invoke', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got: %s", r.Method)
		}

		// Simple request validation (could parse JSON but keeping it simple)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Workflow execution started asynchronously"))
	}))
	defer server.Close()

	os.Setenv("AGENTIC_IGW_URL", server.URL)
	defer os.Unsetenv("AGENTIC_IGW_URL")

	tests := []struct {
		name       string
		input      ExecuteSkillInput
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "Successful Execution",
			input:      ExecuteSkillInput{SkillID: "mcp://skill/device/fleet-update"},
			wantOutput: "Workflow execution started asynchronously",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteSkill(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteSkill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantOutput {
				t.Errorf("ExecuteSkill() got = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}

func TestExecuteSkill_ErrorHandling(t *testing.T) {
	// Mock IGW Server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Temporal client is not initialized"))
	}))
	defer server.Close()

	os.Setenv("AGENTIC_IGW_URL", server.URL)
	defer os.Unsetenv("AGENTIC_IGW_URL")

	_, err := ExecuteSkill(context.Background(), ExecuteSkillInput{SkillID: "test"})
	if err == nil {
		t.Fatal("ExecuteSkill() expected error, got nil")
	}
	
	if !strings.Contains(err.Error(), "A-IGW execution failed (Status 500)") {
		t.Errorf("ExecuteSkill() error = %v, expected it to contain Status 500", err)
	}
	if !strings.Contains(err.Error(), "Temporal client is not initialized") {
		t.Errorf("ExecuteSkill() error = %v, expected it to contain Temporal client message", err)
	}
}
