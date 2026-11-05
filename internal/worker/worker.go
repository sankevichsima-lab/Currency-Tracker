package worker

import (
	"context"
	"log"
	"time"

	"github.com/sankevichsima-lab/Currency-Tracker/internal/client"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
)

type Worker struct {
	client *client.CoinGeckoClient
	db     *storage.Storage
	time   time.Duration
}

func NewWorker(client client.CoinGeckoClient, db storage.Storage, time time.Duration) *Worker {
	worker := Worker{
		client: &client,
		db:     &db,
		time:   time,
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
