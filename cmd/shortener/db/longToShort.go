package db

import (
	"sprint/cmd/shortener/utils"
)

func LongToShort(link string) {
	if _, ok := db[string(link)]; !ok {
		shortLink := utils.LinkShortening()
		db[string(link)] = shortLink
	}
}
