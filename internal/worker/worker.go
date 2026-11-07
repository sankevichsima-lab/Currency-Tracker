package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/sankevichsima-lab/Currency-Tracker/internal/cache"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/client"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
)

// Worker periodically fetches coin prices and updates storage and cache.
// Worker периодически запрашивает котировки и обновляет базу данных и кеш.
type Worker struct {
	client *client.CoinGeckoClient
	db     *storage.Storage
	time   time.Duration
	cache  *cache.Cache
}

// PriceResult represents the result of a single currency fetch operation.
// PriceResult хранит результат одного запроса котировки по монете.
type PriceResult struct {
	ticker string
	price  float64
	err    error
}

// NewWorker creates and initializes a background Worker instance.
// NewWorker создает и инициализирует фоновый воркер.
func NewWorker(client *client.CoinGeckoClient, db *storage.Storage, interval time.Duration, cache *cache.Cache) *Worker {
	worker := Worker{
		client: client,
		db:     db,
		time:   interval,
		cache:  cache,
	}
	return &worker
}

// fetchAndSave runs concurrent requests for all coins and saves the results.
// fetchAndSave запускает параллельные запросы ко всем монетам и сохраняет результаты.
func (w *Worker) fetchAndSave() {
	coins := []string{"bitcoin", "ethereum"}

	// Buffered channel and WaitGroup for parallel data collection (Fan-Out).
	// Буферизированный канал и WaitGroup для параллельного сбора данных (Fan-Out).
	resultsChan := make(chan PriceResult, len(coins))
	var wg sync.WaitGroup

	for _, coin := range coins {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			price, err := w.client.GetPrice(c)
			resultsChan <- PriceResult{
				ticker: c,
				price:  price,
				err:    err,
			}
		}(coin)
	}

	// Wait for all HTTP requests to complete, then close the channel.
	// Дожидаемся завершения всех HTTP-запросов и закрываем канал.
	wg.Wait()
	close(resultsChan)

	// Process incoming results and persist them (Fan-In).
	// Обрабатываем полученные результаты и сохраняем их (Fan-In).
	for res := range resultsChan {
		if res.err != nil {
			log.Printf("Worker: failed to fetch %s: %v", res.ticker, res.err)
			continue
		}

		if err := w.db.SavePrice(res.ticker, res.price); err != nil {
			log.Printf("Worker: failed to save %s to db: %v", res.ticker, err)
			continue
		}

		w.cache.Set(res.ticker, res.price)
		log.Printf("Worker: successfully saved %s price: $%.2f", res.ticker, res.price)
	}
}

// Start launches the background worker loop with graceful shutdown support.
// Start запускает цикл работы воркера с поддержкой корректного завершения.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.time)
	defer ticker.Stop()

	// Initial fetch on application startup.
	// Первый опрос сразу при запуске приложения.
	w.fetchAndSave()

	for {
		select {
		case <-ctx.Done():
			// Exit worker loop when cancel signal is received.
			// Выход из цикла воркера при получении сигнала остановки.
			log.Println("Worker stopped")
			return
		case <-ticker.C:
			w.fetchAndSave()
		}
	}
}
