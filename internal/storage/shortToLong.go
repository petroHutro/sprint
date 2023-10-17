package storage

import (
	"context"
	"errors"
)

func (s *StorageBase) ShortToLong(ctx context.Context, shortLink string) (string, error) {
	if el := s.GetLong(ctx, shortLink); el != "" {
		return el, nil
	}
	return "", errors.New("no short link")
}
