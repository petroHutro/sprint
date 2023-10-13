package storage

import (
	"context"
	"errors"
)

func (s *StorageBase) ShortToLong(ctx context.Context, shortLink string) (string, error) {
	if el := s.GetLong(ctx, shortLink); el != "" {
		return el, nil
	}
	// if _, ok := dbSL[shortLink]; ok {
	// 	return dbSL[shortLink], nil
	// }
	return "", errors.New("no short link")
}
