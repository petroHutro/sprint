package storage

import (
	"context"
	"errors"
)

func (s *StorageBase) ShortToLong(ctx context.Context, shortLink string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("context cansel")
	default:
		if el, err := s.GetLong(ctx, shortLink); el != "" && err == nil {
			return el, nil
		}
		return "", errors.New("no short link")
	}
}
