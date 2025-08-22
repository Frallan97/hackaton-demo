package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationService handles database migrations
type MigrationService struct {
	db *sql.DB
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *sql.DB) *MigrationService {
	return &MigrationService{db: db}
}

// RunMigrations automatically runs all pending migrations
func (ms *MigrationService) RunMigrations() error {
	log.Println("🔄 Starting database migrations...")

	// Create postgres driver instance
	driver, err := postgres.WithInstance(ms.db, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("⚠️  Current migration version: %d (dirty: %v)", version, dirty)
	} else if err == migrate.ErrNilVersion {
		log.Println("📝 No migrations have been run yet")
	} else {
		log.Printf("📝 Current migration version: %d", version)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("✅ Database is up to date - no new migrations to run")
	} else {
		log.Println("✅ All migrations completed successfully")
	}

	// Get final version
	finalVersion, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("⚠️  Could not get final version: %v", err)
	} else if err == migrate.ErrNilVersion {
		log.Println("📝 Database is now at version 0 (no migrations)")
	} else {
		log.Printf("📝 Database is now at version %d", finalVersion)
	}

	return nil
}
