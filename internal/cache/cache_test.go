package cache

import (
	"sync"
	"testing"
)

// TestCache_SetAndGet verifies saving and retrieving values from cache.
// TestCache_SetAndGet проверяет сохранение и получение значений из кеша.
func TestCache_SetAndGet(t *testing.T) {
	c := NewCache()
	c.Set("bitcoin", 75000.0)

	// Verify that the existing key returns correct price data.
	// Проверяем, что существующий ключ возвращает правильную цену.
	item, ok := c.Get("bitcoin")
	if !ok || item.Price != 75000.0 {
		t.Fatal("incorrect value")
	}

	// Verify that a missing key returns false.
	// Проверяем, что отсутствующий ключ возвращает false.
	item, ok = c.Get("ethereum")
	if ok != false {
		t.Fatal("returned incorrect result")
	}
}

// TestCache_Concurrent tests thread safety under simultaneous reads and writes.
// TestCache_Concurrent проверяет потокобезопасность при одновременном чтении и записи.
func TestCache_Concurrent(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup

	// Run concurrent writers.
	// Запускаем параллельные горутины на запись.
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()

			c.Set("bitcoin", val)
		}(float64(i) * 1000.0)
	}

	// Run concurrent readers.
	// Запускаем параллельные горутины на чтение.
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Get("bitcoin")
		}()
	}

	// Wait for all goroutines to finish.
	// Ожидаем завершения всех горутин.
	wg.Wait()
}
