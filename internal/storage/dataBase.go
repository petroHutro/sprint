package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sprint/internal/utils"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type RepError struct {
	Err        error
	Repetition bool
}

func (e *RepError) Error() string {
	return e.Err.Error()
}

func NewErrorRep(err error, repetition bool) *RepError {
	return &RepError{Err: err, Repetition: repetition}
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
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO links
		(long, short)
		VALUES
		($1, $2);
	`, key, val)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return &RepError{Err: err, Repetition: true}
		}
		return fmt.Errorf("cannot set database: %w", err)
	}
	return nil
}

func (d *dataBase) SetAllDB(ctx context.Context, data []string) error {
	repetition := false
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin: %w", err)
	}

	for _, v := range data {
		shortLink := utils.GetShortLink()
		_, err := tx.ExecContext(ctx, `
			INSERT INTO links
			(long, short)
			VALUES
			($1, $2);
    	`, v, shortLink)

		if err != nil {
			if repErr, ok := err.(*RepError); ok && repErr.Repetition {
				repetition = true
			} else {
				tx.Rollback()
				return fmt.Errorf("cannot exec: %w", err)
			}

			// var pgErr *pgconn.PgError
			// if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			// 	repetition = true
			// } else {
			// 	tx.Rollback()
			// 	return fmt.Errorf("cannot exec: %w", err)
			// }
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	if repetition {
		return NewErrorRep(errors.New("long already db"), repetition)
	}
	return nil
}
