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

type Worker struct {
	client *client.CoinGeckoClient
	db     *storage.Storage
	time   time.Duration
	cache  *cache.Cache
}

type PriceResult struct {
	ticker string
	price  float64
	err    error
}

func NewWorker(client *client.CoinGeckoClient, db *storage.Storage, interval time.Duration, cache *cache.Cache) *Worker {
	worker := Worker{
		client: client,
		db:     db,
		time:   interval,
		cache:  cache,
	}
	return &worker
}

func (w *Worker) fetchAndSave() {
	coins := []string{"bitcoin", "ethereum"}

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
	wg.Wait()
	close(resultsChan)

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

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.time)
	defer ticker.Stop()

	w.fetchAndSave()

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped")
			return
		case <-ticker.C:
			w.fetchAndSave()
		}
	}
}
