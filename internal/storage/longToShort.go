package storage

import (
	"context"
	"fmt"
	"sprint/internal/utils"
)

func (s *StorageBase) LongToShort(ctx context.Context, link, fname string) error {
	if shortLink := s.GetShort(ctx, link); shortLink == "" {
		shortLink := utils.GetShortLink()
		s.SetDB(ctx, link, shortLink)
		if fname != "" {
			if err := saveURL(link, shortLink, fname); err != nil {
				return fmt.Errorf("cannot save url in file: %w", err)
			}
		}
	}
	return nil
}
