package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/nazarrbek/subscriptions-service/docs"
	"github.com/nazarrbek/subscriptions-service/internal/config"
	"github.com/nazarrbek/subscriptions-service/internal/handler"
	"github.com/nazarrbek/subscriptions-service/internal/middleware"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	r := chi.NewRouter()

	//r.Get("/", checkResponse)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		logger.Error("connect database", "error", err)
		os.Exit(1)
	}

	repo := repository.NewSubscriptionRepository(db)

	subscriptionService := service.NewSubscriptionService(repo)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	r.Use(middleware.Logger)
	r.Post("/subscriptions", subscriptionHandler.Create)
	r.Get("/subscriptions/{id}", subscriptionHandler.GetByID)
	r.Get("/subscriptions", subscriptionHandler.List)
	r.Put("/subscriptions/{id}", subscriptionHandler.Update)
	r.Delete("/subscriptions/{id}", subscriptionHandler.Delete)
	r.Get("/subscriptions/total", subscriptionHandler.CalculateTotal)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server started", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	case err := <-serverErr:
		logger.Error("server error", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown server", "error", err)
	}

	if err := db.Close(context.Background()); err != nil {
		logger.Error("close database", "error", err)
	}
}

func checkResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
}
