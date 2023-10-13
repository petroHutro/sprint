package storage

import (
	"context"
	"fmt"
	"sprint/internal/config"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

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