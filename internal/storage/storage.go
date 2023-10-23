package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type URL struct {
	LongURL  string `json:"long"`
	ShortURL string `json:"short"`
	UserID   int    `json:"user"`
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
	SetDB(ctx context.Context, key, val string, id int) error
	SetAllDB(ctx context.Context, data []string, id int) error
	GetAllDB(ctx context.Context, id int) ([]Urls, error)
	GetLastID(ctx context.Context) int
}

func (s *StorageBase) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if s != nil {
		if _, ok := s.base.(*dataBase); ok {
			db := s.get()
			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("cannot ping: %w", err)
			}
			return nil
		}
	}
	return errors.New("not flag database, database empty")
}

func Connection(databaseDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot open DataBase: %w", err)
	}
	return db, nil
}
