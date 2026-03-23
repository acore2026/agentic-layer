package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

type embeddingRequest struct {
	Model   string           `json:"model"`
	Content embeddingContent `json:"content"`
}

type embeddingContent struct {
	Parts []embeddingPart `json:"parts"`
}

type embeddingPart struct {
	Text string `json:"text"`
}

type embeddingResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

func getEmbedding(text string) ([]float32, error) {
	apiKey := os.Getenv("AGENTIC_GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("AGENTIC_GEMINI_API_KEY not set")
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent?key=" + apiKey
	reqPayload := embeddingRequest{
		Model: "models/gemini-embedding-001",
		Content: embeddingContent{
			Parts: []embeddingPart{{Text: text}},
		},
	}

	jsonData, _ := json.Marshal(reqPayload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[Embedding] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Embedding] API error (%d): %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("embedding API error (%d): %s", resp.StatusCode, string(body))
	}

	var res embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("[Embedding] Decode error: %v", err)
		return nil, err
	}

	return res.Embedding.Values, nil
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}
