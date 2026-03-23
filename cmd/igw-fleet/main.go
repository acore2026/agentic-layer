package main

import (
	"log"
	"net/http"

	"github.com/google/6g-agentic-core/internal/translator"
)

const (
	acrfURL   = "http://localhost:18080/register"
	skillID   = "mcp://skill/device/fleet-update"
	invokeURL = "http://localhost:18081/invoke"
)

func main() {
	fleetTranslator := translator.NewFleetTranslator()
	handler := translator.NewHandler(fleetTranslator)

	// Register skill with ACRF in background
	go translator.RegisterSkillWithACRF(acrfURL, skillID, invokeURL)

	log.Println("A-IGW (Interworking Gateway) starting on :18081...")
	if err := http.ListenAndServe(":18081", handler); err != nil {
		log.Fatal(err)
	}
}
