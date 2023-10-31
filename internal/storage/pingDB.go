package storage

import (
	"context"
	"errors"
	"fmt"
)

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
