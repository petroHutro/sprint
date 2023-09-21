package storage

import "errors"

var db map[string]string = make(map[string]string)

func GetDB(key string) string {
	return db[key]
}

func SetDB(key, val string) error {
	if _, ok := db[key]; !ok {
		db[key] = val
		return nil
	}
	return errors.New("no key in DB")
}
