package storage

import (
	"encoding/json"
	"io"
	"os"
)

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
	}
	return nil
}
