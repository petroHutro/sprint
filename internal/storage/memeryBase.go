package storage

import (
	"context"
	"errors"
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
		return errors.New("key already DB")
	}
}
