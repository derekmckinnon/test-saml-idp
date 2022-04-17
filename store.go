package main

import (
	"sync"
)

type UserStore interface {
	PutUser(user *User)
	GetUser(id string) (*User, bool)
}

type ServiceProviderStore interface {
	PutServiceProvider(provider *ServiceProvider)
	GetServiceProvider(id string) (*ServiceProvider, bool)
}

type Store interface {
	UserStore
	ServiceProviderStore
}

type MemoryStore struct {
	data sync.Map
}

func getUserKey(id string) string {
	return "/user/" + id
}

func (m *MemoryStore) PutUser(user *User) {
	m.data.Store(getUserKey(user.Id), user)
}

func (m *MemoryStore) GetUser(id string) (*User, bool) {
	user, ok := m.data.Load(getUserKey(id))

	if !ok {
		return nil, ok
	}

	return user.(*User), ok
}

func getServiceProviderKey(id string) string {
	return "/service/" + id
}

func (m *MemoryStore) PutServiceProvider(provider *ServiceProvider) {
	m.data.Store(getServiceProviderKey(provider.Id), provider)
}

func (m *MemoryStore) GetServiceProvider(id string) (*ServiceProvider, bool) {
	sp, ok := m.data.Load(getServiceProviderKey(id))

	if !ok {
		return nil, ok
	}

	return sp.(*ServiceProvider), ok
}
