package storage

import (
	"context"
	"errors"
	"fmt"
	"sprint/internal/utils"
)

func (s *StorageBase) LongToShort(ctx context.Context, link, fname string, id string) error {
	shortLink := utils.GenerateString()
	if err := s.Set(ctx, link, shortLink, id, false); err != nil {
		var repErr *RepError
		if errors.As(err, &repErr) && repErr.Repetition {
			return fmt.Errorf("key already DB: %w", err)
		} else {
			return fmt.Errorf("cannot set: %w", err)
		}
	}

	if fname != "" {
		if err := saveURL(link, shortLink, fname, id, false); err != nil {
			return fmt.Errorf("cannot save url in file: %w", err)
		}
	}
	return nil
}
