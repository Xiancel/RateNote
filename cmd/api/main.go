package main

import (
	"RateNote/internal/db"
	httpHand "RateNote/internal/handler/http"
	repository "RateNote/internal/repository/postgres"
	itemService "RateNote/internal/service/item"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no found file .env")
	}

	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "1234!")
	dbName := getEnv("DB_NAME", "ratenote_db")
	dbSSLMode := getEnv("DB_SSLMODE", "")
	ServerPort := getEnv("APP_PORT", "8080")

	dbConf := db.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
	}

	database, err := db.NewDB(dbConf)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.Close()
	log.Println("Database connected")

	itemRepo := repository.NewItemRepository(database)
	log.Println("Repository initialized")

	itemSrv := itemService.NewService(itemRepo)
	log.Println("Service initialized")

	router := httpHand.NewRouter(httpHand.RouteConfig{
		ItemService: itemSrv,
	})
	log.Println("Http router initialized")

	server := http.Server{
		Addr:         ":" + ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on http://localhost:%s", ServerPort)
		log.Printf("API documentation: http://localhost:%s/api/v1", ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠️ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func getEnv(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}
