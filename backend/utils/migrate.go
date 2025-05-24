package utils

import (
	"fmt"
	"io/fs"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// RunMigrationsWithFS runs all pending database migrations using an embedded filesystem
func RunMigrationsWithFS(dbURL string, migrations fs.FS) error {
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	Info("Database migrations completed successfully")
	return nil
} 