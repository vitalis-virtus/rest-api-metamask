package storage

import (
	"errors"
	"sync"

	"github.com/vitalis-virtus/rest-api-metamask/model"
)

var ErrUserExists = errors.New("user already exists")

type MemoryStorage struct {
	lock  sync.Mutex
	users map[string]model.User
}

func (m *MemoryStorage) CreateIfNotExists(u model.User) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.users[u.Address]; exists {
		return ErrUserExists
	}
	m.users[u.Address] = u
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	ans := MemoryStorage{
		users: make(map[string]model.User),
	}
	return &ans
}
