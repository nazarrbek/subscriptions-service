package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nazarrbek/subscriptions-service/internal/config"
	"github.com/nazarrbek/subscriptions-service/internal/repository"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", checkResponse)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(context.Background())

	log.Printf("Server started on :%s", cfg.AppPort)

	if err := http.ListenAndServe(":"+cfg.AppPort, r); err != nil {
		log.Fatal(err)
	}
}

func checkResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}
