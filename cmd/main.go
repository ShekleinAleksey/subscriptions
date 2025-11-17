package main

import (
	"log"
	"net/http"

	"github.com/ShekleinAleksey/subscriptions/config"
	"github.com/ShekleinAleksey/subscriptions/internal/handler"
	"github.com/ShekleinAleksey/subscriptions/internal/repository"
	"github.com/ShekleinAleksey/subscriptions/internal/service"
	"github.com/ShekleinAleksey/subscriptions/pkg/logger"
	"github.com/ShekleinAleksey/subscriptions/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для управления онлайн-подписками
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error get config: %v", err)
	}

	logger.SetLogrus(cfg.Log.Level)

	db, err := postgres.NewDB(cfg)
	if err != nil {
		logrus.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	logrus.Info("Initializing repository...")
	repo := repository.NewRepository(db)
	logrus.Info("Initializing service...")
	service := service.NewService(repo)
	logrus.Info("Initializing handler...")
	handlers := handler.NewHandler(service)

	router := handlers.InitRoutes()

	logrus.Info("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
