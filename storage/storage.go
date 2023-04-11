package storage

import (
	"sync"

	"github.com/vitalis-virtus/rest-api-metamask/model"
	"github.com/vitalis-virtus/rest-api-metamask/utils/fail"
)

type MemoryStorage struct {
	lock  sync.RWMutex
	users map[string]model.User
}

func (m *MemoryStorage) CreateIfNotExists(u model.User) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.users[u.Address]; exists {
		return fail.ErrUserExists
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

func (m *MemoryStorage) Get(address string) (model.User, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	u, exists := m.users[address]
	if !exists {
		return u, fail.ErrUserNotExists
	}
	return u, nil
}

func (m *MemoryStorage) Update(user model.User) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.users[user.Address] = user
	return nil
}
