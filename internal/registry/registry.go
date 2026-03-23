package registry

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/6g-agentic-core/pkg/models"
)

type Registry interface {
	Register(profile models.SkillProfile) error
	Discover(query string) (models.SkillProfile, bool)
}

type skillEntry struct {
	Profile   models.SkillProfile
	Embedding []float32
}

type InMemoryRegistry struct {
	mu     sync.RWMutex
	skills []skillEntry
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		skills: []skillEntry{},
	}
}

func (r *InMemoryRegistry) Register(profile models.SkillProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Use SkillID and Description for embedding
	textToEmbed := fmt.Sprintf("%s: %s", profile.SkillID, profile.Description)
	embedding, err := getEmbedding(textToEmbed)
	if err != nil {
		// No fallback: Registration MUST have a semantic vector to be valid
		return fmt.Errorf("failed to generate mandatory embedding for skill %s: %v", profile.SkillID, err)
	}

	r.skills = append(r.skills, skillEntry{
		Profile:   profile,
		Embedding: embedding,
	})
	log.Printf("Registered skill: %s (Semantic: true)", profile.SkillID)
	return nil
}

func (r *InMemoryRegistry) Discover(query string) (models.SkillProfile, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.skills) == 0 {
		return models.SkillProfile{}, false
	}

	queryLower := strings.ToLower(query)
	queryWords := strings.Fields(queryLower)

	// HYBRID SEARCH:
	// 1. Identity-based fallback (Exact ID or ID contains query)
	for _, entry := range r.skills {
		idLower := strings.ToLower(entry.Profile.SkillID)
		if idLower == queryLower || strings.Contains(idLower, queryLower) {
			log.Printf("Identity match found for '%s': %s", query, entry.Profile.SkillID)
			return entry.Profile, true
		}
		
		// Word-based match
		for _, word := range queryWords {
			if len(word) < 3 { continue } // Skip short words
			if strings.Contains(idLower, word) {
				log.Printf("Word match found for '%s' (word: %s): %s", query, word, entry.Profile.SkillID)
				return entry.Profile, true
			}
		}
	}

	// 2. Semantic Search (Vector similarity)
	queryEmbedding, err := getEmbedding(query)
	if err != nil {
		log.Printf("Failed to generate embedding for query '%s' (falling back to identity only): %v", query, err)
		return models.SkillProfile{}, false
	}

	var bestMatch models.SkillProfile
	var maxScore float32 = -1.0
	threshold := float32(0.60) // Lowered threshold for better recall in MVP

	foundAnySemantic := false
	for _, entry := range r.skills {
		if entry.Embedding == nil {
			continue
		}
		score := cosineSimilarity(queryEmbedding, entry.Embedding)
		if score > maxScore {
			maxScore = score
			bestMatch = entry.Profile
			foundAnySemantic = true
		}
	}

	if foundAnySemantic {
		log.Printf("Semantic search for '%s' found best match '%s' with score %f", query, bestMatch.SkillID, maxScore)
		if maxScore >= threshold {
			return bestMatch, true
		}
	}

	return models.SkillProfile{}, false
}
