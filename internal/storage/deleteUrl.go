package storage

import (
	"context"
	"errors"
	"fmt"
)

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
