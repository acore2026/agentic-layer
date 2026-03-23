package main

import (
	"log"
	"net/http"

	"github.com/google/6g-agentic-core/internal/registry"
)

func main() {
	reg := registry.NewInMemoryRegistry()
	handler := registry.NewHandler(reg)

	log.Println("ACRF (Agentic Capability Repository Function) starting on :18080...")
	if err := http.ListenAndServe(":18080", handler); err != nil {
		log.Fatal(err)
	}
}
