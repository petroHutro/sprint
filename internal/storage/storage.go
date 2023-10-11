package storage

import (
	"errors"
	"sync"
)

type URL struct {
	LongURL  string `json:"long"`
	ShortURL string `json:"short"`
}

var dbSL map[string]string = make(map[string]string)
var dbLS map[string]string = make(map[string]string)
var sm sync.Mutex

func GetDB(key string) string {
	sm.Lock()
	defer sm.Unlock()
	return dbLS[key]
}

func SetDB(key, val string) error {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := dbLS[key]; !ok {
		dbLS[key] = val
		dbSL[val] = key
		return nil
	}
	return errors.New("key already DB")
}
