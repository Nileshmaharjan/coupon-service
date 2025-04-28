package store

import (
    "errors"
    "sync"
    "time"

    "github.com/Nileshmaharjan/coupon-service/internal/coupon"
)

var _ coupon.Store = (*MemoryStore)(nil)

type MemoryStore struct {
    mu        sync.RWMutex
    campaigns map[string]*coupon.Campaign
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{campaigns: make(map[string]*coupon.Campaign)}
}

func (m *MemoryStore) Create(c *coupon.Campaign) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, exists := m.campaigns[c.ID]; exists {
        return errors.New("exists")
    }
    m.campaigns[c.ID] = c
    return nil
}

func (m *MemoryStore) Get(id string) (*coupon.Campaign, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    c, ok := m.campaigns[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return c, nil
}

func (m *MemoryStore) Issue(id string, now time.Time) (string, error) {
    c, err := m.Get(id)
    if err != nil {
        return "", err
    }
    return c.Issue(now)
}