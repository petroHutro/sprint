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
	FlagDel  bool   `json:"del"`
}

type QueryDelete struct {
	ID   int
	Data string
}

type StorageBase struct {
	base
	baseWithPointer
}

type baseWithPointer interface {
	get() *sql.DB
}

type base interface {
	GetShort(ctx context.Context, key string) (string, error)
	GetLong(ctx context.Context, key string) (string, error)
	SetDB(ctx context.Context, key, val string, id int, flag bool) error
	SetAllDB(ctx context.Context, data []string, id int) error
	GetAllDB(ctx context.Context, id int) ([]Urls, error)
	GetLastID(ctx context.Context) int
	DeleteS(ctx context.Context, id []int, shorts []string) error
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
