package storage

import (
	"context"
	"fmt"
	"sprint/internal/utils"
)

func (s *StorageBase) LongToShort(ctx context.Context, link, fname string) error {
	shortLink := utils.GetShortLink()
	if err := s.SetDB(ctx, link, shortLink); err != nil {
		return fmt.Errorf("cannot set: %w", err)
	}

	if fname != "" {
		if err := saveURL(link, shortLink, fname); err != nil {
			return fmt.Errorf("cannot save url in file: %w", err)
		}
	}
	return nil
}
