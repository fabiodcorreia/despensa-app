package storage

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "modernc.org/sqlite"
)

type Store struct {
	file string
	db   *sql.DB
}

func NewStore(databaseFile string) *Store {
	return &Store{
		file: databaseFile,
	}
}

func NewStoreWithMigrations(databaseFile string, migrations embed.FS) (*Store, error) {
	s := NewStore(databaseFile)
	if err := s.runMigrations(migrations); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Connect() error {
	db, err := sql.Open("sqlite", s.file)
	if err != nil {
		return fmt.Errorf("store connection fail: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("store connection check fail: %w", err)
	}

	return nil
}

func (s *Store) Disconnect() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("store disconnect fail: %w", err)
	}
	return nil
}

func (s *Store) runMigrations(migrations embed.FS) error {
	d, err := iofs.New(migrations, "database/migration")
	if err != nil {
		return fmt.Errorf("store with migrations fail loading migration files: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, fmt.Sprintf("sqlite://%s", s.file))
	if err != nil {
		return fmt.Errorf("store with migrations fail connecting with db: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("store with migrations fail running migrations: %w", err)
		}
	}

	return nil
}
