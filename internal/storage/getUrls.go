package storage

import (
	"context"
)

type Urls struct {
	Short string
	Long  string
}

func (s *StorageBase) GetUrls(ctx context.Context, userID int) ([]Urls, error) {
	urls, err := s.GetAllDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}
