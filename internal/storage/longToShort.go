package storage

import (
	"context"
	"errors"
	"fmt"
	"sprint/internal/utils"
)

func (s *StorageBase) LongToShort(ctx context.Context, link, fname string) error {
	shortLink := utils.GetShortLink()
	if err := s.SetDB(ctx, link, shortLink); err != nil {
		var repErr *RepError
		if errors.As(err, &repErr) && repErr.Repetition {
			return fmt.Errorf("key already DB: %w", err)
		} else {
			return fmt.Errorf("cannot set: %w", err)
		}
	}

	if fname != "" {
		if err := saveURL(link, shortLink, fname); err != nil {
			return fmt.Errorf("cannot save url in file: %w", err)
		}
	}
	return nil
}
