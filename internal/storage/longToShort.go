package storage

import (
	"sprint/internal/utils"
)

func LongToShort(link string) {
	if err := GetDB(string(link)); err == "" {
		shortLink := utils.LinkShortening()
		SetDB(string(link), shortLink)
	}
}
