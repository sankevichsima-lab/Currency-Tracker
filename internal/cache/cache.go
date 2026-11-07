package cache

import (
	"sync"
	"time"
)

// CachedPrice stores the coin price and its last update timestamp.
// CachedPrice хранит цену монеты и время её последнего обновления.
type CachedPrice struct {
	Price     float64
	UpdatedAt time.Time
}

// Cache provides thread-safe in-memory storage for prices.
// Cache представляет собой потокобезопасное хранилище цен в оперативной памяти.
type Cache struct {
	mu    sync.RWMutex
	items map[string]CachedPrice
}

// NewCache creates and initializes a new Cache instance.
// NewCache создает и инициализирует новый экземпляр кеша.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CachedPrice),
	}
}

// Set saves or updates a coin price using a write lock.
// Set сохраняет или обновляет цену монеты, блокируя доступ на запись.
func (c *Cache) Set(ticker string, price float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[ticker] = CachedPrice{
		Price:     price,
		UpdatedAt: time.Now(),
	}
}

// Get retrieves a coin price from cache using a read lock.
// Get получает цену монеты из кеша с блокировкой только на чтение.
func (c *Cache) Get(ticker string) (CachedPrice, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[ticker]
	return item, ok
}
