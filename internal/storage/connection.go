package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDB(db *sql.DB) error {
	// defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func Connection(databaseDSN string) *sql.DB {
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	`localhost`, `url`, `1234`, `url`)

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	return db
}
