package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveURL(long, short, fname string, id int, flagDel bool) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	dataURL := URL{LongURL: long, ShortURL: short, UserID: id, FlagDel: flagDel}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(dataURL)
	if err != nil {
		return fmt.Errorf("cannot encoding dataURL in json: %w", err)
	}
	return nil
}

func saveURLs(urls []URL, fname string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for url := range urls {
		err = encoder.Encode(url)
		if err != nil {
			return fmt.Errorf("cannot encoding dataURL in json: %w", err)
		}
	}
	return nil
}
