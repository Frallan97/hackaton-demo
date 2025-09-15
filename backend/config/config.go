package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBURL      string
	ServerPort string
	Environment string

	// JWT Configuration
	JWTSecretKey string

	// Google OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// Stripe Configuration
	StripeSecretKey      string
	StripePublishableKey string
	StripeWebhookSecret  string
	StripeEndpointSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Check if we're in production (have environment variables set)
	hasEnvVars := os.Getenv("GOOGLE_CLIENT_ID") != "" && os.Getenv("GOOGLE_CLIENT_SECRET") != ""

	if !hasEnvVars {
		// Only load .env file if we don't have environment variables (local development)
		if err := godotenv.Load("../.env"); err != nil {
			// Try current directory as fallback
			if err := godotenv.Load(); err != nil {
				// Don't fail if .env doesn't exist, just log it
				log.Println("No .env file found, using system environment variables")
			}
		}
	} else {
		log.Println("Using environment variables (production mode)")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBURL:      getEnv("DB_URL", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "production"),

		// JWT Configuration
		JWTSecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),

		// Google OAuth Configuration
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/login"),

		// Stripe Configuration
		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripeEndpointSecret: getEnv("STRIPE_ENDPOINT_SECRET", ""),
	}

	// Debug logging for OAuth configuration
	log.Printf("OAuth Configuration - Client ID: %s, Redirect URL: %s",
		config.GoogleClientID, config.GoogleRedirectURL)

	return config
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if c.DBURL != "" {
		return c.DBURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
