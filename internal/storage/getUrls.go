package storage

import (
	"context"
	"fmt"
)

type Urls struct {
	Short string
	Long  string
}

func (s *StorageBase) GetUrls(ctx context.Context, userID int) ([]Urls, error) {
	urls, err := s.GetAllId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get urls: %v", err)
	}
	return urls, nil
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
