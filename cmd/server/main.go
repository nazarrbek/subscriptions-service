package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/nazarrbek/subscriptions-service/docs"
	"github.com/nazarrbek/subscriptions-service/internal/config"
	"github.com/nazarrbek/subscriptions-service/internal/handler"
	"github.com/nazarrbek/subscriptions-service/internal/repository"
	"github.com/nazarrbek/subscriptions-service/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для управления подписками пользователей.

// @host localhost:8080
// @BasePath /
func main() {
	r := chi.NewRouter()
	//r.Get("/", checkResponse)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewSubscriptionRepository(db)

	subscriptionService := service.NewSubscriptionService(repo)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	r.Post("/subscriptions", subscriptionHandler.Create)
	r.Get("/subscriptions/{id}", subscriptionHandler.GetByID)
	r.Get("/subscriptions", subscriptionHandler.List)
	r.Put("/subscriptions/{id}", subscriptionHandler.Update)
	r.Delete("/subscriptions/{id}", subscriptionHandler.Delete)
	r.Get("/subscriptions/total", subscriptionHandler.CalculateTotal)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
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
