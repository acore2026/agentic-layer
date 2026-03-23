package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/6g-agentic-core/internal/agent"
	"github.com/google/6g-agentic-core/internal/testutil"
)

func TestSystem_EndToEnd_FleetWakeUp(t *testing.T) {
	// 1. Setup ACRF
	acrfURL, acrfCloser, err := testutil.SetupACRF()
	if err != nil {
		t.Fatalf("failed to setup ACRF: %v", err)
	}
	defer acrfCloser()
	os.Setenv("AGENTIC_ACRF_URL", acrfURL)

	// 2. Setup A-IGW
	igwURL, igwCloser, err := testutil.SetupIGW(acrfURL)
	if err != nil {
		t.Fatalf("failed to setup IGW: %v", err)
	}
	defer igwCloser()
	os.Setenv("AGENTIC_IGW_URL", igwURL)

	// 3. Setup AAIHF with Mock Agent
	mockAgent, err := agent.NewMockCoreAgent()
	if err != nil {
		t.Fatalf("failed to create mock agent: %v", err)
	}

	aaihfURL, aaihfCloser, err := testutil.SetupAAIHF(mockAgent)
	if err != nil {
		t.Fatalf("failed to setup AAIHF: %v", err)
	}
	defer aaihfCloser()

	// Give a small window for everything to be ready
	time.Sleep(200 * time.Millisecond)

	// 4. Trigger Intent
	intentReq := map[string]string{
		"prompt":  "Wake up the fleet for updates",
		"user_id": "test-operator",
	}
	jsonData, _ := json.Marshal(intentReq)

	resp, err := http.Post(aaihfURL+"/intent", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("failed to send intent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("intent request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 5. Assert Response
	var intentResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&intentResp); err != nil {
		t.Fatalf("failed to decode intent response: %v", err)
	}

	response := intentResp["response"]
	t.Logf("Received response: %s", response)

	if !strings.Contains(response, "successfully triggered mcp://skill/device/fleet-update") {
		t.Errorf("unexpected response content: %s", response)
	}
}

func TestSystem_ACRF_Discovery(t *testing.T) {
	acrfURL, acrfCloser, err := testutil.SetupACRF()
	if err != nil {
		t.Fatalf("failed to setup ACRF: %v", err)
	}
	defer acrfCloser()

	// Register a dummy skill
	skill := map[string]string{
		"skill_id":            "mcp://skill/test",
		"entity_type":         "UE",
		"agentic_service_uri": "http://localhost:9999",
	}
	skillJSON, _ := json.Marshal(skill)
	_, err = http.Post(acrfURL+"/register", "application/json", bytes.NewBuffer(skillJSON))
	if err != nil {
		t.Fatalf("failed to register skill: %v", err)
	}

	// Discover it
	resp, err := http.Get(fmt.Sprintf("%s/discover?skill_id=%s", acrfURL, "mcp://skill/test"))
	if err != nil {
		t.Fatalf("failed to discover skill: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("discovery failed with status %d", resp.StatusCode)
	}

	var discovered map[string]any
	json.NewDecoder(resp.Body).Decode(&discovered)
	if discovered["skill_id"] != "mcp://skill/test" {
		t.Errorf("expected skill_id mcp://skill/test, got %v", discovered["skill_id"])
	}
}
