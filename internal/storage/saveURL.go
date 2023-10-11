package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveURL(long, short, fname string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()
	dataURL := URL{LongURL: long, ShortURL: short}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(dataURL)
	if err != nil {
		return fmt.Errorf("cannot encoding dataURL in json: %w", err)
	}
	return nil
}
