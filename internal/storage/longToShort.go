package storage

import (
	"encoding/json"
	// "fmt"
	"io"
	"os"
	"sprint/internal/utils"
)

type URL struct {
	LongURL  string `json:"long"`
	ShortURL string `json:"short"`
}

func saveURL(long, short, fname string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	dataURL := URL{LongURL: long, ShortURL: short}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(dataURL)
	if err != nil {
		return err
	}
	return nil
}

func LoadURL(fname string) error {
	file, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for {
		var url URL
		if err := decoder.Decode(&url); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		SetDB(url.LongURL, url.ShortURL)
		// fmt.Println(url.LongURL, url.ShortURL)
	}
	return nil
}

func LongToShort(link, fname string) {
	if err := GetDB(string(link)); err == "" {
		shortLink := utils.LinkShortening()
		SetDB(string(link), shortLink)
		saveURL(string(link), shortLink, fname)
	}
}
