package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("record not found")

type Store struct {
	file string
	ctx  context.Context
	db   *sql.DB
}

var _ LocationStore = (*Store)(nil)
var _ ItemStore = (*Store)(nil)

func NewStore(ctx context.Context, databaseFile string) *Store {
	return &Store{
		file: databaseFile,
		ctx:  ctx,
	}
}

func NewStoreWithMigrations(ctx context.Context, databaseFile string, migrations fs.FS) (*Store, error) {
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
	return NewStore(ctx, databaseFile), nil
}

func (s *Store) Connect() error {
	db, err := sql.Open("sqlite", s.file)
	if err != nil {
		return fmt.Errorf("store connection fail: %w", err)
	}

	if err = db.Ping(); err != nil {
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
