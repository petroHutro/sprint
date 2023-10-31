package storage

import (
	"context"
	"database/sql"
	"errors"
)

var ErrorDeleteURL = errors.New("url delete")
var ErrorContext = errors.New("context cansel")
var ErrorSet = errors.New("already DB")
var ErrorNoGet = errors.New("no key")

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
