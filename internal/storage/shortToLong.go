package storage

import "errors"

func ShortToLong(shortLink string) (string, error) {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := dbSL[shortLink]; ok {
		return dbSL[shortLink], nil
	}
	return "", errors.New("no short link")
}
