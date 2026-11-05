package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sankevichsima-lab/Currency-Tracker/internal"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/client"
	"github.com/sankevichsima-lab/Currency-Tracker/internal/storage"
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
	price, err := cgClient.GetBitcoinPrice()
	if err != nil {
		log.Printf("Failed to retrieve Bitcoin exchange rate: %v", err)
	} else {
		fmt.Printf("Current Bitcoin price: $%.2f USD\n", price)

		if err = dbStorage.SavePrice("bitcoin", price); err != nil {
			log.Printf("Failed to save price to DB: %v", err)
		} else {
			fmt.Println("Price successfully saved to SQLite database!")
		}
	}

	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Server is active on port:", cfg.Server.Port)
	err = http.ListenAndServe(":"+cfg.Server.Port, nil)
	if err != nil {
		log.Fatal("problem with listen server: ", err)
	}
}
