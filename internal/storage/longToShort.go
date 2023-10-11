package storage

import (
	"fmt"
	"sprint/internal/utils"
)

func LongToShort(link, fname string) error {
	if shortLink := GetDB(link); shortLink == "" {
		shortLink := utils.GetShortLink()
		SetDB(link, shortLink)
		if err := saveURL(link, shortLink, fname); err != nil {
			return fmt.Errorf("cannot save url in file: %w", err)
		}
	}
	return nil
}
