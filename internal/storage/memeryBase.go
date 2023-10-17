package storage

import (
	"context"
	"errors"
	"fmt"
	"sprint/internal/utils"
	"sync"
)

type memeryBase struct {
	dbSL map[string]string
	dbLS map[string]string
	sm   sync.Mutex
}

func newMemeryBase() *memeryBase {
	return &memeryBase{dbSL: make(map[string]string), dbLS: make(map[string]string)}
}

func (m *memeryBase) GetShort(ctx context.Context, key string) string {
	select {
	case <-ctx.Done():
		return ""
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		return m.dbLS[key]
	}
}

func (m *memeryBase) GetLong(ctx context.Context, key string) string {
	select {
	case <-ctx.Done():
		return ""
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		return m.dbSL[key]
	}
}

func (m *memeryBase) SetDB(ctx context.Context, key, val string) error {
	select {
	case <-ctx.Done():
		return errors.New("context cansel")
	default:
		m.sm.Lock()
		defer m.sm.Unlock()
		if _, ok := m.dbLS[key]; !ok {
			m.dbLS[key] = val
			m.dbSL[val] = key
			return nil
		}
		return &RepError{Err: errors.New("key already DB"), Repetition: true}
	}
}

func (m *memeryBase) SetAllDB(ctx context.Context, data []string) error {
	repetition := false
	for _, v := range data {
		shortLink := utils.GetShortLink()
		err := m.SetDB(ctx, v, shortLink)
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
