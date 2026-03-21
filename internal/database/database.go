package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Ability to read .sql files
	_ "github.com/jackc/pgx/v5/stdlib"                         // Standard Postgres driver
	"github.com/jmoiron/sqlx"
)

// NewDatabase initializes the Postgres connection and runs migrations
func NewDatabase(connStr string) (*sqlx.DB, error) {
	// DSN = Data Source Name. In production, move these to environment variables!//
	dsn := connStr

	// 1. Open the connection using 'pgx' driver
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// 2. Configure Connection Pool (Java equivalent: HikariCP settings)
	db.SetMaxOpenConns(25)                 // Max active connections
	db.SetMaxIdleConns(25)                 // Max idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Recycle connections to avoid memory leaks

	// 3. Verify connection is alive
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// 4. Run Migrations automatically
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("✅ Database connection established and migrations synced.")
	return db, nil
}

// runMigrations looks at the /migrations folder and updates the DB schema
func runMigrations(dsn string) error {
	// We use "file://migrations" because our folder is at the project root
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}

	// Apply all 'Up' migrations/
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("🚀 Migrations applied successfully!")
	return nil
}
