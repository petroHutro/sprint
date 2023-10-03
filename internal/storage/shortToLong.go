package storage

import "errors"

func ShortToLong(shortLink string) (string, error) {
	sm.Lock()
	defer sm.Unlock()
	var link string
	for key, value := range db {
		if value == shortLink {
			link = key
			break
		}
	}
	if link == "" {
		return "", errors.New("no short link")
	}
	return link, nil
}
