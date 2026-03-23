package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type SearchSkillInput struct {
	SkillID string `json:"skill_id" jsonschema:"The unique URI of the skill to discover (e.g. mcp://skill/device/fleet-update)"`
}

func SearchSkill(ctx context.Context, input SearchSkillInput) (string, error) {
	acrfURL := os.Getenv("AGENTIC_ACRF_URL")
	if acrfURL == "" {
		acrfURL = "http://localhost:18080"
	}

	log.Printf("[Tool:SearchSkill] Searching for: %s at %s", input.SkillID, acrfURL)

	url := fmt.Sprintf("%s/discover?skill_id=%s", acrfURL, input.SkillID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call ACRF: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Sprintf("Skill %s not found in ACRF registry.", input.SkillID), nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ACRF returned error status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[Tool:SearchSkill] Found profile: %s", string(body))
	return string(body), nil
}

type ExecuteSkillInput struct {
	SkillID string `json:"skill_id" jsonschema:"The URI of the skill to execute (must have been discovered first)"`
}

func ExecuteSkill(ctx context.Context, input ExecuteSkillInput) (string, error) {
	igwURL := os.Getenv("AGENTIC_IGW_URL")
	if igwURL == "" {
		igwURL = "http://localhost:18081"
	}

	log.Printf("[Tool:ExecuteSkill] Executing skill: %s at %s", input.SkillID, igwURL)

	payload := map[string]string{"skill_id": input.SkillID}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/invoke", igwURL)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call A-IGW: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("A-IGW execution failed (Status %d): %s", resp.StatusCode, string(body))
	}

	log.Printf("[Tool:ExecuteSkill] Execution successful: %s", string(body))
	return string(body), nil
}
