package storage

import (
	"context"
	"sync"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/store"
)

type cacheEntry struct {
	incident  models.IncidentData
	expiresAt time.Time
}

type CacheStore struct {
	cache map[string]cacheEntry
	next  store.IncidentStore
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCacheStore(store store.IncidentStore, ttl time.Duration) *CacheStore {
	return &CacheStore{
		cache: make(map[string]cacheEntry),
		next:  store,
		ttl:   ttl,
	}
}

func (c *CacheStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	c.mu.RLock()
	cache, ok := c.cache[id]
	c.mu.RUnlock()
	if ok {
		return cache.incident, nil
	}
	incident, err := c.next.Get(ctx, id)
	if err != nil {
		return incident, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[id] = cacheEntry{
		incident:  incident,
		expiresAt: time.Now().Add(c.ttl),
	}

	return incident, nil
}
func (c *CacheStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	incidents, err := c.next.GetAll(ctx)
	return incidents, err
}
func (c *CacheStore) Save(ctx context.Context, incident *models.IncidentData) error {
	err := c.next.Save(ctx, incident)
	c.mu.Lock()
	delete(c.cache, incident.Id)
	c.mu.Unlock()
	return err
}
func (c *CacheStore) Delete(ctx context.Context, id string) error {
	err := c.next.Delete(ctx, id)
	if err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.cache, id)
	c.mu.Unlock()
	return err
}
