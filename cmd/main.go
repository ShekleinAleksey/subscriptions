package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ShekleinAleksey/subscriptions/config"
	"github.com/ShekleinAleksey/subscriptions/internal/handler"
	"github.com/ShekleinAleksey/subscriptions/internal/repository"
	"github.com/ShekleinAleksey/subscriptions/internal/service"
	"github.com/ShekleinAleksey/subscriptions/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для управления онлайн-подписками
// @host localhost:8080
// @BasePath /api/v1
func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// databaseURL := "postgres://admin:root123@localhost:5432/subscriptiondb" //fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	// // 	cfg.User,
	// // 	cfg.Password,
	// // 	cfg.Host,
	// // 	cfg.Port,
	// // 	cfg.DBName,
	// // 	cfg.SSLMode,
	// // )

	// logrus.Info("Running database migrations...")
	// if err := migrate.RunMigrations(databaseURL); err != nil {
	// 	logrus.Fatalf("Failed to run migrations: %v", err)
	// }

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	logrus.Info("Initializing repository...")
	repo := repository.NewRepository(db)
	logrus.Info("Initializing service...")
	service := service.NewService(repo)
	logrus.Info("Initializing handler...")
	handlers := handler.NewHandler(service)

	router := handlers.InitRoutes()

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
