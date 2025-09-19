package database

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
	_ "github.com/lib/pq"
)

// DBManager manages database connections and status
type DBManager struct {
	DB        *sql.DB
	Connected atomic.Bool
	Config    *config.Config
}

// NewDBManager creates a new database manager
func NewDBManager(cfg *config.Config) *DBManager {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		log.Printf("failed to open DB: %v", err)
		return nil
	}

	// Optimized connection pool settings for faster startup
	db.SetMaxOpenConns(10)                  // Reduced from 25 for faster startup
	db.SetMaxIdleConns(3)                   // Reduced from 5 for faster startup
	db.SetConnMaxLifetime(10 * time.Minute) // Increased from 5 minutes
	db.SetConnMaxIdleTime(2 * time.Minute)  // Increased from 1 minute

	manager := &DBManager{
		DB:     db,
		Config: cfg,
	}

	// Start connection monitoring with reduced frequency in production
	go manager.monitorConnection()

	return manager
}

// monitorConnection continuously monitors the database connection
func (dm *DBManager) monitorConnection() {
	// Use different monitoring intervals based on environment
	interval := 5 * time.Second // Default for development
	if os.Getenv("ENVIRONMENT") == "production" {
		interval = 30 * time.Second // Less frequent in production
	}

	for {
		if dm.DB != nil {
			err := dm.DB.Ping()
			if err == nil {
				if !dm.Connected.Load() {
					log.Println("✅ Connected to database successfully")
					dm.Connected.Store(true)
				}
			} else {
				if dm.Connected.Load() {
					log.Printf("❌ Lost connection to database: %v", err)
					dm.Connected.Store(false)
				} else {
					// Only log in debug mode to reduce noise
					if os.Getenv("DEBUG") == "true" {
						log.Printf("⚠️  Unable to ping database: %v", err)
					}
				}
			}
		}
		time.Sleep(interval)
	}
}

// IsConnected returns whether the database is currently connected
func (dm *DBManager) IsConnected() bool {
	return dm.Connected.Load()
}

// Close closes the database connection
func (dm *DBManager) Close() error {
	if dm.DB != nil {
		return dm.DB.Close()
	}
	return nil
}
