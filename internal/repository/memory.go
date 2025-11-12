package repository

import (
	"sync"

	"github.com/includingByMeAndMyself/urlshortening/internal/model"
)

type MemoryRepo struct {
	urls map[string]*model.URLPair
	mu   sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		urls: make(map[string]*model.URLPair),
	}
}

func (r *MemoryRepo) Save(pair *model.URLPair) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.urls[pair.ID] = pair
}

func (r *MemoryRepo) Get(id string) (*model.URLPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pair, ok := r.urls[id]
	return pair, ok
}
