package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrorDeleteURL = errors.New("url delete")

type URL struct {
	LongURL  string `json:"long"`
	ShortURL string `json:"short"`
	UserID   string `json:"user"`
	FlagDel  bool   `json:"del"`
}

type QueryDelete struct {
	ID   string
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
	Set(ctx context.Context, key, val string, id string, flag bool) error
	SetAll(ctx context.Context, data []string, id string) error
	GetAllID(ctx context.Context, id string) ([]Urls, error)
	GetAll(ctx context.Context) ([]URL, error)
	delete(ctx context.Context, id []string, shorts []string) error
}

func (s *StorageBase) PingDB(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
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
}

func Connection(databaseDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot open DataBase: %w", err)
	}
	return db, nil
}

func (s *StorageBase) DeleteURL(ctx context.Context, fname string, id []string, shorts []string) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
		err := s.delete(ctx, id, shorts)
		if err != nil {
			return fmt.Errorf("cannot deleta: %w", err)
		}

		if fname != "" {
			urls, err := s.GetAll(ctx)
			if err != nil {
				return fmt.Errorf("cannot get all: %w", err)
			}
			saveURLs(urls, fname)
		}
		return nil
	}
}
