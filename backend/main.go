// @title        React-Go-App API
// @version      0.1.0
// @description  Auto‑generated Swagger docs
// @host         localhost:8080
// @BasePath     /
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/frallan97/react-go-app-backend/config"
	"github.com/frallan97/react-go-app-backend/database"
	_ "github.com/frallan97/react-go-app-backend/docs"
	"github.com/frallan97/react-go-app-backend/handlers"
	"github.com/frallan97/react-go-app-backend/services"
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
	defer dbManager.Close()

	// Initialize services
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

	log.Printf("listening on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
