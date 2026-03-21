package registry

import (
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
	profile, ok := r.skills[skillID]
	return profile, ok
}
