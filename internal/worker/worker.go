package worker

import (
	"context"
	"log"
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
	price, err := w.client.GetBitcoinPrice()
	if err != nil {
		log.Printf("Worker: failed to fetch price: %v", err)
		return
	}

	err = w.db.SavePrice("bitcoin", price)
	if err != nil {
		log.Printf("Worker: failed to save to db: %v", err)
		return
	}

	w.cache.Set("bitcoin", price)

	log.Printf("Worker: successfully saved bitcoin price: $%.2f", price)
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
