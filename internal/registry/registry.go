package registry

import (
	"strings"
	"sync"

	"github.com/google/6g-agentic-core/pkg/models"
)

type Registry interface {
	Register(profile models.SkillProfile)
	Discover(skillID string) (models.SkillProfile, bool)
}

type InMemoryRegistry struct {
	mu     sync.RWMutex
	skills map[string]models.SkillProfile
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		skills: make(map[string]models.SkillProfile),
	}
}

func (r *InMemoryRegistry) Register(profile models.SkillProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skills[profile.SkillID] = profile
}

func (r *InMemoryRegistry) Discover(skillID string) (models.SkillProfile, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1. Exact match
	if profile, ok := r.skills[skillID]; ok {
		return profile, true
	}

	// 2. Partial/Keyword match (Simple Semantic Engine simulation)
	for id, profile := range r.skills {
		if strings.Contains(strings.ToLower(id), strings.ToLower(skillID)) {
			return profile, true
		}
	}

	return models.SkillProfile{}, false
}
