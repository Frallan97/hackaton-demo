// @title        Go-React-Stripe Template API
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
	"sync"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/frallan97/hackaton-demo-backend/database"
	_ "github.com/frallan97/hackaton-demo-backend/docs"
	"github.com/frallan97/hackaton-demo-backend/events"
	"github.com/frallan97/hackaton-demo-backend/handlers"
	"github.com/frallan97/hackaton-demo-backend/services"
)

func main() {
	log.Println("🚀 Starting Go-React-Stripe Template API...")
	startTime := time.Now()

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("📡 Environment: %s", cfg.Environment)

	// Initialize database manager with optimized connection
	log.Println("🗄️  Connecting to database...")
	dbManager := database.NewDBManager(cfg)
	if dbManager == nil {
		log.Fatal("❌ Failed to create database manager")
	}
	defer dbManager.Close()

	// Fast database connection check with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connected := make(chan bool, 1)
	go func() {
		for {
			if dbManager.IsConnected() {
				connected <- true
				return
			}
			select {
			case <-ctx.Done():
				connected <- false
				return
			default:
				time.Sleep(100 * time.Millisecond) // Much faster polling
			}
		}
	}()

	select {
	case success := <-connected:
		if !success {
			log.Fatal("❌ Database connection timeout (10s)")
		}
		log.Println("✅ Database connected")
	case <-ctx.Done():
		log.Fatal("❌ Database connection timeout (10s)")
	}

	// Initialize services in parallel
	log.Println("🔧 Initializing services...")

	var wg sync.WaitGroup
	var migrationService *services.MigrationService
	var eventBus events.EventBus
	var eventService *events.EventService

	// Start migration in goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		migrationService = services.NewMigrationService(dbManager.DB)
		if err := migrationService.RunMigrations(); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
		log.Println("✅ Migrations completed")
	}()

	// Initialize event system in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		natsURL := os.Getenv("NATS_URL")
		if natsURL == "" {
			natsURL = "nats://localhost:4222"
		}

		// Try NATS with short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		natsChan := make(chan events.EventBus, 1)
		go func() {
			if bus, err := events.NewNATSEventBus(natsURL); err == nil {
				natsChan <- bus
			}
		}()

		select {
		case eventBus = <-natsChan:
			log.Println("✅ NATS connected")
		case <-ctx.Done():
			log.Println("⚠️  NATS unavailable, using fallback event bus")
			eventBus = events.NewEventBus()
		}

		_ = events.NewEventHandlerManager(eventBus) // Initialize handlers
		eventService = events.NewEventService(eventBus)
	}()

	// Wait for parallel initialization to complete
	wg.Wait()

	// Create services (fast, no I/O operations)
	userService := services.NewUserService(dbManager.DB)
	jwtService := services.NewJWTService(cfg.JWTSecretKey)
	googleOAuthService := services.NewGoogleOAuthService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	// Initialize router (fast, no I/O operations)
	log.Println("🌐 Setting up routes...")
	router := handlers.NewRouter(dbManager, userService, jwtService, googleOAuthService, eventService, cfg)
	handler := router.SetupRoutes()

	// Publish system startup event (non-blocking)
	go func() {
		if err := eventService.PublishSystemStartup(); err != nil {
			log.Printf("⚠️  Failed to publish startup event: %v", err)
		}
	}()

	elapsed := time.Since(startTime)
	log.Printf("🚀 Server listening on port %s", cfg.ServerPort)
	log.Printf("✨ Go-React-Stripe Template API ready! (startup: %v)", elapsed)

	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}
