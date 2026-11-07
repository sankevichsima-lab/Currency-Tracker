package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sankevichsima-lab/Currency-Tracker/internal"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/api"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/cache"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/client"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/worker"
)

type PingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := PingResponse{
		Status:  "ok",
		Message: "json mess",
	}
	json.NewEncoder(w).Encode(data)
}

func main() {
	fmt.Println("Трекер валют запускается...")

	priceCache := cache.NewCache()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := internal.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	dbStorage, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}
	defer dbStorage.Close()

	fmt.Println("Port:", cfg.Server.Port)

	cgClient := client.NewCoinGeckoClient(cfg.API.URL, 10*time.Second)

	wrk := worker.NewWorker(cgClient, dbStorage, cfg.API.Interval, priceCache)
	go wrk.Start(ctx)

	http.HandleFunc("/ping", pingHandler)

	apiHandler := api.NewHandler(dbStorage, priceCache)
	http.HandleFunc("/api/history", apiHandler.GetHistoryHandler)
	http.HandleFunc("/api/price", apiHandler.GetPriceHandler)

	fmt.Println("Server is active on port:", cfg.Server.Port)
	go func() {
		err := http.ListenAndServe(":"+cfg.Server.Port, nil)
		if err != nil {
			log.Fatal("problem with listen server: ", err)
		}
	}()
	<-ctx.Done()
	fmt.Println("Приложение успешно остановлено.")
}
