package cache

import (
	"sync"
	"testing"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache()
	c.Set("bitcoin", 75000.0)

	item, ok := c.Get("bitcoin")
	if !ok || item.Price != 75000.0 {
		t.Fatal("incorrect value")
	}

	item, ok = c.Get("ethereum")
	if ok != false {
		t.Fatal("returned incorrect result")
	}
}

func TestCache_Concurrent(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()

			c.Set("bitcoin", val)
		}(float64(i) * 1000.0)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Get("bitcoin")
		}()
	}
	wg.Wait()
}
