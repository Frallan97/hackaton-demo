package database

import (
	"database/sql"
	"log"
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

	// Configure optimized connection pool settings
	db.SetMaxOpenConns(10)                  // Reduced for better performance
	db.SetMaxIdleConns(5)                   // Keep idle connections ready
	db.SetConnMaxLifetime(10 * time.Minute) // Longer lifetime reduces reconnections
	db.SetConnMaxIdleTime(2 * time.Minute)  // Reasonable idle time

	manager := &DBManager{
		DB:     db,
		Config: cfg,
	}

	// Start connection monitoring
	go manager.monitorConnection()

	return manager
}

// monitorConnection continuously monitors the database connection
func (dm *DBManager) monitorConnection() {
	for {
		if dm.DB != nil {
			err := dm.DB.Ping()
			if err == nil {
				if !dm.Connected.Load() {
					log.Println("connected to Postgres successfully")
					dm.Connected.Store(true)
				}
			} else {
				if dm.Connected.Load() {
					log.Printf("lost connection to DB: %v", err)
					dm.Connected.Store(false)
				} else {
					// Only log connection failures if we're not already connected
					// This reduces noise during startup
					log.Printf("unable to ping DB: %v", err)
				}
			}
		}
		time.Sleep(30 * time.Second) // Optimized interval for production
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
