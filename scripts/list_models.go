package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("AGENTIC_GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("AGENTIC_GEMINI_API_KEY not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	page, err := client.Models.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range page.Items {
		fmt.Printf("Model: %s\n", m.Name)
	}
}
