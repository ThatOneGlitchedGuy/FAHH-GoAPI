package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type CacheStore struct {
	data map[string]*CacheItem
	mu   sync.RWMutex
}

func NewCacheStore() *CacheStore {
	cs := &CacheStore{
		data: make(map[string]*CacheItem),
	}
	go cs.cleanupExpired()
	return cs
}

func (cs *CacheStore) Set(key string, value interface{}, duration time.Duration) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.data[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (cs *CacheStore) Get(key string) (interface{}, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	item, exists := cs.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

func (cs *CacheStore) Delete(key string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.data, key)
}

func (cs *CacheStore) Clear() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.data = make(map[string]*CacheItem)
}

func (cs *CacheStore) Exists(key string) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	item, exists := cs.data[key]
	if !exists {
		return false
	}

	if time.Now().After(item.ExpiresAt) {
		return false
	}

	return true
}

func (cs *CacheStore) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cs.mu.Lock()
		now := time.Now()
		for key, item := range cs.data {
			if now.After(item.ExpiresAt) {
				delete(cs.data, key)
			}
		}
		cs.mu.Unlock()
	}
}

func (cs *CacheStore) GetSize() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.data)
}

type CacheKey string

const (
	UserCachePrefix    CacheKey = "user:"
	ProductCachePrefix CacheKey = "product:"
	OrderCachePrefix   CacheKey = "order:"
)

func (ck CacheKey) With(id interface{}) string {
	return string(ck) + toString(id)
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return string(rune(val))
	case uint:
		return string(rune(val))
	default:
		return ""
	}
}
