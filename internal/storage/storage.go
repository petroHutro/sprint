package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sprint/internal/config"
	"sync"
	"time"
)

type URL struct {
	LongURL  string `json:"long"`
	ShortURL string `json:"short"`
}

type StorageBase struct {
	base
	baseWithPointer
}

type baseWithPointer interface {
	get() *sql.DB
}

type base interface {
	GetShort(ctx context.Context, key string) string
	GetLong(ctx context.Context, key string) string
	SetDB(ctx context.Context, key, val string) error
}

type memeryBase struct {
	dbSL map[string]string
	dbLS map[string]string
	sm   sync.Mutex
}

type dataBase struct {
	db *sql.DB
}

func (d *dataBase) get() *sql.DB {
	return d.db
}

func newDataBase(databaseDSN string) (*dataBase, error) {
	db, err := Connection(databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot connection database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}
	return &dataBase{db: db}, nil
}

func newMemeryBase() *memeryBase {
	return &memeryBase{dbSL: make(map[string]string), dbLS: make(map[string]string)}
}

func (m *memeryBase) GetShort(ctx context.Context, key string) string {
	select {
	case <-ctx.Done():
		return ""
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		return m.dbLS[key]
	}
}

func (m *memeryBase) GetLong(ctx context.Context, key string) string {
	select {
	case <-ctx.Done():
		return ""
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		return m.dbSL[key]
	}
}

func (m *memeryBase) SetDB(ctx context.Context, key, val string) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		if _, ok := m.dbLS[key]; !ok {
			m.dbLS[key] = val
			m.dbSL[val] = key
			return nil
		}
		return errors.New("key already DB")
	}

}

func (m *dataBase) GetLong(ctx context.Context, key string) string {
	row := m.db.QueryRowContext(ctx, `
	SELECT long FROM links WHERE short = $1`,
		key,
	)

	var long string
	err := row.Scan(&long)
	if err != nil {
		return ""
	}
	return long
}

func (m *dataBase) GetShort(ctx context.Context, key string) string {
	row := m.db.QueryRowContext(ctx, `
		SELECT short FROM links WHERE long = $1`,
		key,
	)

	var short string
	err := row.Scan(&short)
	if err != nil {
		return ""
	}
	return short
}

func (m *dataBase) SetDB(ctx context.Context, key, val string) error {
	if m.GetShort(ctx, key) == "" {
		_, err := m.db.ExecContext(ctx, `
			INSERT INTO links
			(long, short)
			VALUES
			($1, $2);
    	`, key, val)
		return fmt.Errorf("cannot set database: %w", err)
	}
	return nil
}

func InitStorage(conf *config.Storage) (*StorageBase, error) {
	if conf.DatabaseDSN != "" {
		db, err := newDataBase(conf.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("cannot create data base: %w", err)
		}

		base := base(db)
		basePointer := baseWithPointer(db)
		storageBase := StorageBase{base, basePointer}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := storageBase.createTable(ctx); err != nil {
			return nil, fmt.Errorf("cannot create table: %w", err)
		}

		return &storageBase, nil
	} else {
		db := newMemeryBase()
		if conf.FileStoragePath != "" {
			if err := db.LoadURL(conf.FileStoragePath); err != nil {
				return nil, fmt.Errorf("cannot create memery base: %w", err)
			}
		}
		memery := base(db)
		basePointer := baseWithPointer(nil)
		storageBase := StorageBase{memery, basePointer}
		return &storageBase, nil
	}
}

func (s *StorageBase) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if _, ok := s.base.(*dataBase); ok {
		db := s.get()
		if err := db.PingContext(ctx); err != nil {
			return fmt.Errorf("cannot ping: %w", err)
		}
		return nil
	}
	return errors.New("not flag database, database empty")
}

func (s *StorageBase) createTable(ctx context.Context) error {
	tx, err := s.get().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'links') THEN
				CREATE TABLE links (
					id SERIAL PRIMARY KEY,
					long  TEXT NOT NULL,
					short TEXT NOT NULL
				);
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("cannot request create table: %w", err)
	}
	return tx.Commit()
}
