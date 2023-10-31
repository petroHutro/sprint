package storage

import (
	"context"
	"errors"
	"fmt"
	"sprint/internal/utils"
	"sync"
)

type memeryBase struct {
	dbSL  map[string]string
	dbLS  map[string]string
	dbLID map[string]string
	dbLF  map[string]bool
	sm    sync.Mutex
}

func newMemeryBase() *memeryBase {
	return &memeryBase{
		dbSL:  make(map[string]string),
		dbLS:  make(map[string]string),
		dbLID: make(map[string]string),
		dbLF:  make(map[string]bool),
	}
}

func (m *memeryBase) GetShort(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()

		value, ok := m.dbLS[key]
		if ok {
			return value, nil
		}
		return "", errors.New("no long")
	}
}

func (m *memeryBase) GetLong(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()

		if long, ok := m.dbSL[key]; ok {
			if m.dbLF[long] {
				return "", ErrorDeleteURL
			}
			return long, nil
		}
		return "", errors.New("not key")
	}
}

func (m *memeryBase) Set(ctx context.Context, key, val string, id string, flag bool) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		if _, ok := m.dbLS[key]; !ok {
			m.dbLS[key] = val
			m.dbSL[val] = key
			m.dbLID[key] = id
			m.dbLF[key] = flag
			return nil
		}
		return &RepError{Err: errors.New("key already DB"), Repetition: true}
	}
}

func (m *memeryBase) SetAll(ctx context.Context, data []string, id string) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
		repetition := false

		for _, v := range data {
			shortLink := utils.GenerateString()
			err := m.Set(ctx, v, shortLink, id, false)
			if err != nil {
				var repErr *RepError
				if errors.As(err, &repErr) {
					repetition = true
				} else {
					return fmt.Errorf("cannot set: %w", err)
				}
			}
		}
		if repetition {
			return NewErrorRep(errors.New("long already db"), repetition)
		}
		return nil
	}
}

func (m *memeryBase) GetAllID(ctx context.Context, id string) ([]Urls, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()

		var urls []Urls
		for key, val := range m.dbLID {
			if val == id {
				short, ok := m.dbLS[key]
				if !ok {
					return nil, fmt.Errorf("cannot get: %w", errors.New("no long"))
				}
				if !m.dbLF[key] {
					urls = append(urls, Urls{Long: key, Short: short})
				}
			}
		}
		if urls == nil {
			return nil, errors.New("no data on id")
		}
		return urls, nil
	}
}

func (m *memeryBase) GetAll(ctx context.Context) ([]URL, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()

		var urls []URL
		for key, el := range m.dbLS {
			urls = append(urls, URL{LongURL: key, ShortURL: el, UserID: m.dbLID[key], FlagDel: m.dbLF[key]})
		}
		if urls == nil {
			return nil, errors.New("storage is empty")
		}
		return urls, nil
	}
}

func (m *memeryBase) delete(ctx context.Context, id []string, shorts []string) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		for i, value := range id {
			if long, ok := m.dbSL[shorts[i]]; ok {
				if m.dbLID[long] == value {
					m.dbLF[long] = true
				}
			}
		}
	}
	return nil
}
