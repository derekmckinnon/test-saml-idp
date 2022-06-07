package idp

import (
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
)

const (
	usersPrefix    = "/users/"
	servicesPrefix = "/services/"
	sessionsPrefix = "/sessions/"
)

type Store struct {
	samlidp.MemoryStore
}

func (s *Store) GetUser(name string) (user *samlidp.User, err error) {
	err = s.Get(usersPrefix+name, &user)
	return
}

func (s *Store) GetUsers() ([]*samlidp.User, error) {
	return getResources[samlidp.User](s, usersPrefix, s.GetUser)
}

func (s *Store) AddUser(user *samlidp.User) error {
	return s.Put(usersPrefix+user.Name, user)
}

func (s *Store) GetServiceProvider(id string) (service *samlidp.Service, err error) {
	err = s.Get(servicesPrefix+id, &service)
	return
}

func (s *Store) GetServiceProviders() ([]*samlidp.Service, error) {
	return getResources[samlidp.Service](s, servicesPrefix, s.GetServiceProvider)
}

func (s *Store) AddServiceProvider(service *samlidp.Service) error {
	return s.Put(servicesPrefix+service.Metadata.EntityID, service)
}

func (s *Store) GetSession(id string) (session *saml.Session, err error) {
	err = s.Get(sessionsPrefix+id, &session)
	return
}

func (s *Store) GetSessions() ([]*saml.Session, error) {
	return getResources[saml.Session](s, sessionsPrefix, s.GetSession)
}

func (s *Store) AddSession(session *saml.Session) error {
	return s.Put(sessionsPrefix+session.ID, session)
}

func getResources[T any](store *Store, prefix string, getter func(string) (*T, error)) ([]*T, error) {
	keys, _ := store.List(prefix)

	resources := make([]*T, len(keys))

	for i, key := range keys {
		resource, err := getter(key)
		if err != nil {
			return nil, err
		}

		resources[i] = resource
	}

	return resources, nil
}
