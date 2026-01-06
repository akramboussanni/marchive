//go:build debug
// +build debug

package db

import (
	"embed"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var DB *sqlx.DB

func Init(dsn string) {
	log.Println("using sqlite db")

	var err error
	// Add WAL mode and busy timeout to DSN for better concurrency
	DB, err = sqlx.Open("sqlite", "marchive.db?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL")
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	// Set connection pool settings for SQLite
	DB.SetMaxOpenConns(1) // SQLite works best with a single writer
	DB.SetMaxIdleConns(1)

	if err = DB.Ping(); err != nil {
		log.Fatalf("cannot ping database: %v", err)
	}
}

func RunMigrations() {
	driver, err := sqlite.WithInstance(DB.DB, &sqlite.Config{})
	if err != nil {
		log.Fatalf("failed to create sqlite driver: %v", err)
	}

	migrationsSub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		log.Fatalf("failed to get migrations subdir: %v", err)
	}

	d, err := iofs.New(migrationsSub, ".")
	if err != nil {
		log.Fatalf("failed to create iofs driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("sqlite migrations applied successfully")
}
