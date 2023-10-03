package storage

import (
	"errors"
	"sync"
)

type URL struct {
	LongURL  string `json:"long"`
	ShortURL string `json:"short"`
}

var db map[string]string = make(map[string]string)
var sm sync.Mutex

func GetDB(key string) string {
	sm.Lock()
	defer sm.Unlock()
	return db[key]
}

func SetDB(key, val string) error {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := db[key]; !ok {
		db[key] = val
		return nil
	}
	return errors.New("no key in DB")
}
