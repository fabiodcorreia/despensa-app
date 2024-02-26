package storage

import (
	"database/sql"
	"fmt"
	"io/fs"

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

func NewStoreWithMigrations(databaseFile string, migrations fs.FS) (*Store, error) {
	d, err := iofs.New(migrations, "database/migration")
	if err != nil {
		return nil, fmt.Errorf("store with migrations - loading migration files: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, "sqlite://"+databaseFile)
	if err != nil {
		return nil, fmt.Errorf("store with migrations - connecting with db: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return nil, fmt.Errorf("store with migrations - running migrations: %w", err)
		}
	}
	return NewStore(databaseFile), nil
}

func (s *Store) Connect() error {
	db, err := sql.Open("sqlite", s.file)
	if err != nil {
		return fmt.Errorf("store connection fail: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("store connection check fail: %w", err)
	}
	s.db = db
	return nil
}

func (s *Store) Disconnect() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("store disconnect fail: %w", err)
	}
	return nil
}
