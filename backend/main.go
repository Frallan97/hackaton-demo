// @title        Hackaton Demo API
// @version      0.1.0
// @description  Auto‑generated Swagger docs
// @host         localhost:8080
// @BasePath     /
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/frallan97/hackaton-demo-backend/database"
	_ "github.com/frallan97/hackaton-demo-backend/docs"
	"github.com/frallan97/hackaton-demo-backend/handlers"
	"github.com/frallan97/hackaton-demo-backend/services"
)

func main() {
	// Only log environment variables in debug mode
	if os.Getenv("DEBUG") == "true" {
		log.Println("Environment variables at startup:")
		for _, e := range os.Environ() {
			log.Println(e)
		}
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database manager with faster connection timeout
	dbManager := database.NewDBManager(cfg)
	if dbManager == nil {
		log.Fatal("Failed to create database manager")
	}
	defer dbManager.Close()

	// Wait for database connection with shorter timeout and exponential backoff
	log.Println("Establishing database connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	connected := make(chan bool, 1)
	go func() {
		backoff := 100 * time.Millisecond
		maxBackoff := 2 * time.Second

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if dbManager.IsConnected() {
					connected <- true
					return
				}
				time.Sleep(backoff)
				if backoff < maxBackoff {
					backoff *= 2
				}
			}
		}
	}()

	select {
	case <-connected:
		log.Println("Database connection established")
	case <-ctx.Done():
		log.Fatal("Failed to establish database connection within timeout")
	}

	// Single database ping test instead of multiple
	if err := dbManager.DB.Ping(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	// Initialize migration service and run migrations
	migrationService := services.NewMigrationService(dbManager.DB)
	log.Println("Running database migrations...")
	if err := migrationService.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services concurrently
	userService := services.NewUserService(dbManager.DB)
	jwtService := services.NewJWTService(cfg.JWTSecretKey)
	googleOAuthService := services.NewGoogleOAuthService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	// Initialize router with all controllers and services
	router := handlers.NewRouter(dbManager, userService, jwtService, googleOAuthService)
	handler := router.SetupRoutes()

	log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
