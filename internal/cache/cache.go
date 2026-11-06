package cache

import (
	"sync"
	"time"
)

type CachedPrice struct {
	Price     float64
	UpdatedAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]CachedPrice
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CachedPrice),
	}
}

func (c *Cache) Set(ticker string, price float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[ticker] = CachedPrice{
		Price:     price,
		UpdatedAt: time.Now(),
	}
}

func (c *Cache) Get(ticker string) (CachedPrice, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[ticker]
	return item, ok
}
