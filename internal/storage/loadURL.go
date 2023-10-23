package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

func (m *memeryBase) LoadURL(fname string) error {
	file, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file with url: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var url URL
		if err := decoder.Decode(&url); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("cannot read file: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(time.Millisecond*500))
		defer cancel()
		m.SetDB(ctx, url.LongURL, url.ShortURL, url.UserID)
	}
	return nil
}
