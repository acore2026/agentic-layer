package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/6g-agentic-core/internal/agent"
	"github.com/google/6g-agentic-core/internal/testutil"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Load .env from project root
	_ = godotenv.Load("../.env")
	os.Exit(m.Run())
}

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

	// Give a small window for everything to be ready (registration takes time now because of embedding call)
	time.Sleep(2 * time.Second)

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

func TestSystem_ACRF_SemanticDiscovery(t *testing.T) {
	acrfURL, acrfCloser, err := testutil.SetupACRF()
	if err != nil {
		t.Fatalf("failed to setup ACRF: %v", err)
	}
	defer acrfCloser()

	// 1. Register a skill with a clear description
	skill := map[string]string{
		"skill_id":            "mcp://skill/network/slice-optimize",
		"description":         "Optimize network slices for high bandwidth and low latency video streaming.",
		"entity_type":         "NF",
		"agentic_service_uri": "http://localhost:9999",
	}
	skillJSON, _ := json.Marshal(skill)
	_, err = http.Post(acrfURL+"/register", "application/json", bytes.NewBuffer(skillJSON))
	if err != nil {
		t.Fatalf("failed to register skill: %v", err)
	}

	// 2. Discover it using a semantically similar query (not exact words)
	query := "improve video quality on the mobile network"
	resp, err := http.Get(fmt.Sprintf("%s/discover?skill_id=%s", acrfURL, url.QueryEscape(query)))
	if err != nil {
		t.Fatalf("failed to discover skill: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("discovery failed with status %d: %s", resp.StatusCode, string(body))
	}

	var discovered map[string]any
	json.NewDecoder(resp.Body).Decode(&discovered)
	if discovered["skill_id"] != "mcp://skill/network/slice-optimize" {
		t.Errorf("expected skill_id mcp://skill/network/slice-optimize, got %v", discovered["skill_id"])
	}
}

func TestSystem_ACRF_LowConfidence_404(t *testing.T) {
	acrfURL, acrfCloser, err := testutil.SetupACRF()
	if err != nil {
		t.Fatalf("failed to setup ACRF: %v", err)
	}
	defer acrfCloser()

	// Register a specific skill
	skill := map[string]string{
		"skill_id":    "mcp://skill/device/reboot",
		"description": "Remotely reboot a specific hardware device.",
	}
	skillJSON, _ := json.Marshal(skill)
	http.Post(acrfURL+"/register", "application/json", bytes.NewBuffer(skillJSON))

	// Query for something completely unrelated
	query := "how to cook a pizza"
	resp, err := http.Get(fmt.Sprintf("%s/discover?skill_id=%s", acrfURL, url.QueryEscape(query)))
	if err != nil {
		t.Fatalf("failed to call discover: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for unrelated query, got %d", resp.StatusCode)
	}
}

func TestSystem_EndToEnd_Hallucination(t *testing.T) {
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

	time.Sleep(1 * time.Second)

	// 4. Trigger Intent with a query that causes MockAgent to fail discovery
	intentReq := map[string]string{
		"prompt":  "Cook a pizza",
		"user_id": "test-operator",
	}
	jsonData, _ := json.Marshal(intentReq)

	resp, err := http.Post(aaihfURL+"/intent", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("failed to send intent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("intent request failed with status %d", resp.StatusCode)
	}

	var intentResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&intentResp); err != nil {
		t.Fatalf("failed to decode intent response: %v", err)
	}

	response := intentResp["response"]
	if !strings.Contains(response, "I couldn't find a skill") {
		t.Errorf("expected graceful hallucination handling, got: %s", response)
	}
}

func TestSystem_EndToEnd_IGWDown(t *testing.T) {
	// 1. Setup ACRF
	acrfURL, acrfCloser, err := testutil.SetupACRF()
	if err != nil {
		t.Fatalf("failed to setup ACRF: %v", err)
	}
	defer acrfCloser()
	os.Setenv("AGENTIC_ACRF_URL", acrfURL)

	// 2. Setup A-IGW - deliberately pointing to a dead port to simulate IGW down
	os.Setenv("AGENTIC_IGW_URL", "http://localhost:12345") // Dead port

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

	time.Sleep(1 * time.Second)

	// Note: We need to register the skill in ACRF manually since IGW is "down"
	skill := map[string]string{
		"skill_id":            "mcp://skill/device/fleet-update",
		"description":         "Trigger a wake-up",
		"entity_type":         "NF",
		"agentic_service_uri": "http://localhost:12345/invoke",
	}
	skillJSON, _ := json.Marshal(skill)
	http.Post(acrfURL+"/register", "application/json", bytes.NewBuffer(skillJSON))

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

	// The reasoning engine should bubble up the 500 error from the tool execution failure
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 status when IGW is down, got: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "failed to call A-IGW: Post \"http://localhost:12345/invoke\"") {
		t.Errorf("expected error containing A-IGW connection failure, got: %s", string(body))
	}
}
