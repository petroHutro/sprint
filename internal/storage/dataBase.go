package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

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

func (d *dataBase) GetLong(ctx context.Context, key string) string {
	row := d.db.QueryRowContext(ctx, `
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

func (d *dataBase) GetShort(ctx context.Context, key string) string {
	row := d.db.QueryRowContext(ctx, `
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

func (d *dataBase) SetDB(ctx context.Context, key, val string) error {
	if d.GetShort(ctx, key) == "" {
		_, err := d.db.ExecContext(ctx, `
			INSERT INTO links
			(long, short)
			VALUES
			($1, $2);
    	`, key, val)
		return fmt.Errorf("cannot set database: %w", err)
	}
	return nil
}
