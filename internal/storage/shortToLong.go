package storage

import (
	"context"
	"errors"
	"fmt"
)

func (s *StorageBase) ShortToLong(ctx context.Context, shortLink string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("context cansel")
	default:
		el, err := s.GetLong(ctx, shortLink)
		if el != "" && err == nil {
			return el, nil
		} else if err != nil {
			return "", fmt.Errorf("cannot get: %w", err)
		}
		return "", errors.New("no short link")
	}
}
