package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SearchSkillInput struct {
	SkillID string `json:"skill_id" description:"The unique URI of the skill to discover (e.g., mcp://skill/device/fleet-update)"`
}

func SearchSkill(ctx context.Context, input SearchSkillInput) (string, error) {
	acrfURL := os.Getenv("ACRF_URL")
	if acrfURL == "" {
		acrfURL = "http://localhost:8080"
	}

	url := fmt.Sprintf("%s/discover?skill_id=%s", acrfURL, input.SkillID)
	resp, err := http.Get(url)
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
	return string(body), nil
}

type ExecuteSkillInput struct {
	SkillID string `json:"skill_id" description:"The URI of the skill to execute (must have been discovered first)"`
}

func ExecuteSkill(ctx context.Context, input ExecuteSkillInput) (string, error) {
	igwURL := os.Getenv("IGW_URL")
	if igwURL == "" {
		igwURL = "http://localhost:8081"
	}

	payload := map[string]string{"skill_id": input.SkillID}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/invoke", igwURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call A-IGW: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("A-IGW execution failed (Status %d): %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
