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

func (s *Store) AddUser(user *samlidp.User) error {
	return s.Put(usersPrefix+user.Name, user)
}

func (s *Store) GetServiceProvider(id string) (service *samlidp.Service, err error) {
	err = s.Get(servicesPrefix+id, &service)
	return
}

func (s *Store) AddServiceProvider(service *samlidp.Service) error {
	return s.Put(servicesPrefix+service.Metadata.EntityID, service)
}

func (s *Store) GetSession(id string) (session *saml.Session, err error) {
	err = s.Get(sessionsPrefix+id, &session)
	return
}

func (s *Store) AddSession(session *saml.Session) error {
	return s.Put(sessionsPrefix+session.ID, session)
}
