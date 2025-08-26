// @title        Hackaton Demo API
// @version      0.1.0
// @description  Auto‑generated Swagger docs
// @host         localhost:8080
// @BasePath     /
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/frallan97/hackaton-demo-backend/database"
	_ "github.com/frallan97/hackaton-demo-backend/docs"
	"github.com/frallan97/hackaton-demo-backend/events"
	"github.com/frallan97/hackaton-demo-backend/handlers"
	"github.com/frallan97/hackaton-demo-backend/services"
)

func main() {
	// Print all environment variables for debugging
	log.Println("Environment variables at startup:")
	for _, e := range os.Environ() {
		log.Println(e)
	}

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Attempting to connect to Postgres with DSN: %s", cfg.GetDSN())

	// Initialize database manager
	dbManager := database.NewDBManager(cfg)
	if dbManager == nil {
		log.Fatal("Failed to create database manager")
	}
	defer dbManager.Close()

	// Wait for database connection to be established
	log.Println("Waiting for database connection...")
	for i := 0; i < 30; i++ { // Wait up to 30 seconds
		if dbManager.IsConnected() {
			log.Println("Database connection established successfully")
			break
		}
		if i == 29 {
			log.Fatal("Failed to establish database connection after 30 seconds")
		}
		log.Println("Waiting for database connection...")
		time.Sleep(1 * time.Second)
	}

	// Test database connection stability
	log.Println("Testing database connection stability...")
	for i := 0; i < 3; i++ {
		if err := dbManager.DB.Ping(); err != nil {
			log.Fatalf("Database connection test failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	log.Println("Database connection is stable")

	// Initialize migration service and run migrations
	migrationService := services.NewMigrationService(dbManager.DB)
	log.Println("Running database migrations...")
	if err := migrationService.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Small delay to ensure database connection is stable after migrations
	log.Println("Waiting for database connection to stabilize after migrations...")
	time.Sleep(2 * time.Second)

	// Final connection test after migrations
	if err := dbManager.DB.Ping(); err != nil {
		log.Fatalf("Database connection lost after migrations: %v", err)
	}
	log.Println("Database connection confirmed stable after migrations")

	// Initialize NATS event bus and services
	var eventBus events.EventBus
	var err error

	// Get NATS URL from environment, fallback to localhost for development
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	eventBus, err = events.NewNATSEventBus(natsURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS at %s, falling back to custom event bus: %v", natsURL, err)
		eventBus = events.NewEventBus()
	}
	_ = events.NewEventHandlerManager(eventBus) // Initialize handlers
	eventService := events.NewEventService(eventBus)

	// Publish system startup event
	if err := eventService.PublishSystemStartup(); err != nil {
		log.Printf("Warning: Failed to publish system startup event: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(dbManager.DB)
	jwtService := services.NewJWTService(cfg.JWTSecretKey)
	googleOAuthService := services.NewGoogleOAuthService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	// Initialize router with all controllers and services
	router := handlers.NewRouter(dbManager, userService, jwtService, googleOAuthService, eventService, cfg)
	handler := router.SetupRoutes()

	log.Printf("listening on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
