package storage

import (
	"sprint/internal/utils"
)

func LongToShort(link, fname string) {
	if err := GetDB(string(link)); err == "" {
		shortLink := utils.LinkShortening()
		SetDB(string(link), shortLink)
		saveURL(string(link), shortLink, fname)
	}
}
