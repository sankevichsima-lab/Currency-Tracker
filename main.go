package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sankevichsima-lab/Currency-Tracker/internal"
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
		log.Fatal("error: %v", err)
	}

	fmt.Println("Port:", cfg.Server.Port)
}
